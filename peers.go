package main

import (
	"context"
	"os"
	"sync"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/grandcat/zeroconf"
)

const MulticastGroupAddr = "[ff12::9316]:9316"

type PeerManager struct {
	m *sync.Mutex
	callbacks []func(*Peer)
}

type Peer struct {
	Addresses []net.IP
	Port int
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		m: &sync.Mutex{},
		callbacks: []func(*Peer){},
	}
}

func (pm *PeerManager) RegisterNewPeerCallback(f func(*Peer)) {
	pm.m.Lock()
	pm.callbacks = append(pm.callbacks, f)
	pm.m.Unlock()
}

func (pm *PeerManager) Run() {
	log.Debug("PeerManager run")

	// register us
	name, err := os.Hostname()
	if err != nil {
		log.WithError(err).Error("unable to get hostname")
		return
	}
	server, err := zeroconf.Register(name, "_mercury._tcp", "local.", 9316, nil, nil)
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

	for entry := range entries {
		log.Debug(entry)
		p := &Peer{
			Port: entry.Port,
		}
		p.Addresses = append(p.Addresses, entry.AddrIPv6...)
		p.Addresses = append(p.Addresses, entry.AddrIPv4...)

		pm.m.Lock()
		for _, cb := range pm.callbacks {
			cb(p)
		}
		pm.m.Unlock()
	}
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
