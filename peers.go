package main

import (
	"context"
	"net"
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"time"
	"strconv"
	// "encoding/pem"

	log "github.com/sirupsen/logrus"
    "github.com/go-chi/chi"
	"github.com/grandcat/zeroconf"
	"github.com/gofrs/uuid"
)

type PeerManager struct {
	s *Store
	handler http.Handler
	myUUID string
	m *sync.Mutex
	peers map[string]*Peer
}

type Peer struct {
	ID string
	Addresses []net.IP
	Port int
}

func NewPeerManager(s *Store) *PeerManager {
	pm := &PeerManager{
		s: s,
		m: &sync.Mutex{},
		peers: map[string]*Peer{},
	}

	r := chi.NewRouter()
	r.Get("/messages", pm.messagesHandler)
	r.Get("/peers", pm.peersHandler)
	r.Get("/pubkey", pm.pubkeyHandler)
	pm.handler = r

	return pm
}

func (pm *PeerManager) Run() {
	log.Debug("PeerManager run")

	// set up our own listener so we can figure out what port we get
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.WithError(err).Error("unable to create listener")
		return
	}
	port := listener.Addr().(*net.TCPAddr).Port
	log.Debugf("PeerManager listening at %s", listener.Addr())
	go func() {
		err := http.Serve(listener, pm.handler)
		log.WithError(err).Error("PeerManager unable to listen and serve")
	}()

	// register us
	u, err := uuid.NewV4()
	if err != nil {
		log.WithError(err).Error("unable to generate UUID")
		return
	}
	pm.myUUID = u.String()

	go func() {
		// wait a sec for our server to start
		for {
			server, err := zeroconf.Register(pm.myUUID, "_mercury._tcp", "local.", port, nil, nil)
			if err != nil {
				log.WithError(err).Error("unable to register zeroconf service")
				continue
			}
			<-time.After(time.Second*5)
			server.Shutdown()
		}
	}()

	// find others
	entries := make(chan *zeroconf.ServiceEntry)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.WithError(err).Error("unable to create resolver")
		return
	}
	ctx := context.Background()
	err = resolver.Browse(ctx, "_mercury._tcp", "local.", entries)
	if err != nil {
		log.WithError(err).Error("unable to browse")
		return
	}

	ticker := time.Tick(time.Second * 1)
	for {
		select {
		case entry := <-entries:
			pm.handleEntry(entry)
		case <- ticker:
			pm.fetchMessages()
		}
	}
}

func (pm *PeerManager) handleEntry(entry *zeroconf.ServiceEntry) {
	if entry.Instance == pm.myUUID {
		// ignore ourself
		return
	}
	p := &Peer{
		ID: entry.Instance,
		Port: entry.Port,
	}
	p.Addresses = append(p.Addresses, entry.AddrIPv6...)
	p.Addresses = append(p.Addresses, entry.AddrIPv4...)

	pm.m.Lock()
	defer pm.m.Unlock()
	if _, ok := pm.peers[p.ID]; ok {
		// already have this peer
		log.Debugf("peer %s already known", p.ID)
		return
	}

	log.Debugf("new peer: %+v", p)
	pm.peers[p.ID] = p
}

func (pm *PeerManager) messagesHandler(w http.ResponseWriter, r *http.Request) {
	msgs, err := pm.s.EncryptedMessages()
	if err != nil {
		http.Error(w, "unable to get messages", http.StatusInternalServerError)
		log.WithError(err).Error("unable to get messages")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(msgs); err != nil {
		http.Error(w, "unable to encode messages", http.StatusInternalServerError)
		log.WithError(err).Error("unable to encode messages")
		return
	}
}

func (pm *PeerManager) peersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	pm.m.Lock()
	defer pm.m.Unlock()
	if err := enc.Encode(pm.peers); err != nil {
		http.Error(w, "unable to encode peers", http.StatusInternalServerError)
		log.WithError(err).Error("unable to encode peers")
		return
	}
}

func (pm *PeerManager) fetchMessages() {
	c := http.Client{
		Timeout: time.Second * 5,
	}

	pm.m.Lock()
	defer pm.m.Unlock()
	for peerID, peer := range pm.peers {
		// try all addrs
		retrieved := false
		for _, addr := range peer.Addresses {
			endpoint := url.URL{
				Scheme: "http",
				Path: "/messages",
			}

			// what kind of address is it?
			if addr.To4() != nil {
				endpoint.Host = addr.String() + ":" + strconv.Itoa(peer.Port)
			} else if addr.To16() != nil {
				// idk if this is right but whatever
				endpoint.Host = "[" + addr.String() + "]:" + strconv.Itoa(peer.Port)
			} else {
				log.Errorf("peer %s addr type is unknown", peerID)
				continue
			}

			req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
			if err != nil {
				// log.WithError(err).Error("unable to create request")
				continue
			}
			resp, err := c.Do(req)
			if err != nil {
				// log.WithError(err).Error("unable to get response")
				continue
			}
			retrieved = true

			// parse messages
			dec := json.NewDecoder(resp.Body)
			var msgs []*EncryptedMessage
			err = dec.Decode(&msgs)
			if err != nil {
				log.WithError(err).Error("unable to decode messages")
				continue
			}
			pm.handleMessages(msgs)

			log.Debugf("got messages from peer %s", peerID)
			break
		}
		if !retrieved {
			// remove this peer since we can't reach it
			delete(pm.peers, peerID)
			log.Debugf("removed dead peer %s", peerID)
		}
	}
}

func (pm *PeerManager) handleMessages(msgs []*EncryptedMessage) {
	// get our key
	i, err := pm.s.MyInfo()
	if err != nil {
		log.WithError(err).Error("unable to get my key")
		return
	}
	myKey, _ := KeyPairFromBytes(i.PrivateKey)

	msgIDs, err := pm.s.ProcessedMessageIDs()
	if err != nil {
		log.WithError(err).Error("unable to get processed message IDs")
		return
	}

	for _, msg := range msgs {
		if _, ok := msgIDs[msg.ID]; ok {
			continue
		}
		ret, err := myKey.UnSign("msg", string(msg.Contents))
		if err != nil {
			// not for us
			err := pm.s.AddEncryptedMessage(msg)
			if err != nil {
				log.WithError(err).Error("unable to create encrypted message")
				continue
			}
			log.Debug("added message for other peer")
		} else {
			// this is for us
			dmsg := &DecryptedMessage{
				ID: msg.ID,
				Sent: msg.Sent,
				Contents: *ret,
			} 
			err := pm.s.AddDecryptedMessage(dmsg)
			if err != nil {
				log.WithError(err).Error("unable to create decrypted message")
				continue
			}
			log.Debug("added message for us")
		}
	}
}

type PeerKey struct {
	PubKey string
}

func (pm *PeerManager) pubkeyHandler(w http.ResponseWriter, r *http.Request) {
	i, err := pm.s.MyInfo()
	if err != nil {
		log.WithError(err).Error("unable to get my key")
		http.Error(w, "unable to get my key", http.StatusInternalServerError)
		return
	}

	// make this into a contact so privkey isn't exposed
	c := &Contact{
		Name: i.Name,
		PublicKey: i.PublicKey,
	}
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(c)
	if err != nil {
		log.WithError(err).Error("unable to encode contact")
		http.Error(w, "unable to encode contact", http.StatusInternalServerError)
		return
	}
}
