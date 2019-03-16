package main

import (
	log "github.com/sirupsen/logrus"
)

type PeerManager struct {}

func NewPeerManager() *PeerManager {
    return &PeerManager{}
}

func (pm *PeerManager) Run() {
	log.Debug("PeerManager run")

	// join multicast group, send broadcast messages, handle incoming messages...
}
