package main

import (
	"context"
	"net"
	"encoding/json"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
    "github.com/go-chi/chi"
	"github.com/grandcat/zeroconf"
	"github.com/gofrs/uuid"
)

const MulticastGroupAddr = "[ff12::9316]:9316"

type PeerManager struct {
	s *Store
	handler http.Handler
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
	// r.Post("/peers", pm.addPeersHandler)
	r.Post("/notify", pm.notifyHandler)
	pm.handler = r

	return pm
}

func (pm *PeerManager) Run() {
	log.Debug("PeerManager run")

	go func() {
		err := http.ListenAndServe(":3001", pm.handler)
		log.WithError(err).Error("PeerManager unable to listen and serve")
	}()

	// register us
	u, err := uuid.NewV4()
	if err != nil {
		log.WithError(err).Error("unable to get hostname")
		return
	}
	server, err := zeroconf.Register(u.String(), "_mercury._tcp", "local.", 9316, nil, nil)
	if err != nil {
		log.WithError(err).Error("unable to register zeroconf service")
		return
	}
	defer server.Shutdown()

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

	for {
		select {
		case entry := <-entries:
			pm.handleEntry(entry)
		}
	}
}

func (pm *PeerManager) handleEntry(entry *zeroconf.ServiceEntry) {
	p := &Peer{
		ID: entry.Instance,
		Port: entry.Port,
	}
	p.Addresses = append(p.Addresses, entry.AddrIPv6...)
	p.Addresses = append(p.Addresses, entry.AddrIPv4...)

	log.Debugf("new peer: %+v", p)

	pm.m.Lock()
	defer pm.m.Unlock()
	if _, ok := pm.peers[p.ID]; ok {
		// already have this peer
		log.Debugf("peer %s already known", p.ID)
		return
	}
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
	var peers []*Peer
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(peers); err != nil {
		http.Error(w, "unable to decode peers", http.StatusInternalServerError)
		log.WithError(err).Error("unable to decode peers")
		return
	}

	log.Debugf("got peers %+v", peers)
}

func (pm *PeerManager) notifyHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("notifyHandler")
	// pm.notify <- 
}


// func oldpeermanagerstuff() {
// 	// join multicast group, send broadcast messages, handle incoming messages...
// 	groupAddr, err := net.ResolveUDPAddr("udp6", MulticastGroupAddr)
// 	if err != nil {
// 		log.WithError(err).Error("unable to resolve multicast group addr")
// 		return
// 	}

// 	c, err := net.ListenPacket("udp6", ":0")
// 	if err != nil {
// 		log.WithError(err).Error("unable to announce")
// 		return
// 	}
// 	defer c.Close()

// 	intfs, err := net.Interfaces()
// 	if err != nil {
// 		log.WithError(err).Error("unable to get interfaces")
// 		return
// 	}

// 	p := ipv6.NewPacketConn(c)

// 	joined := false
// 	for _, intf := range intfs {
// 		if err := p.JoinGroup(&intf, groupAddr); err != nil {
// 			continue
// 		}
// 		joined = true
// 	}
// 	if !joined {
// 		log.Error("unable to join multicast group on any interface")
// 		return
// 	}

// 	go func() {
// 		wcm := &ipv6.ControlMessage{
// 			HopLimit: 1,
// 		}
// 		ticker := time.Tick(time.Second)
// 		for range ticker {
// 			for _, intf := range intfs {
// 				wcm.IfIndex = intf.Index
// 				_, err := p.WriteTo([]byte("hello world"), wcm, groupAddr)
// 				if err != nil {
// 					// this will fail on a lot of interfaces
// 					continue
// 				}
// 				log.Debugf("wrote announcement to group on %s", intf.Name)
// 			}
// 		}
// 	}()

// 	b := make([]byte, 1500)
// 	for {
// 		n, _, src, err := p.ReadFrom(b)
// 		if err != nil {
// 			log.WithError(err).Error("unable to read")
// 			continue
// 		}
// 		// if rcm.Dst.IsMulticast() {
// 		// 	if rcm.Dst.Equal(groupAddr) {
// 		// 		log.Debug("message to group")
// 		// 	}
// 		// }
// 		log.Debugf("got %d-byte message from %s: \"%s\"", n, src.String(), b)
// 	}
// }
