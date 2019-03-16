package main

import (
	log "github.com/sirupsen/logrus"
)

type MessagesManager struct {
	s *Store
}

func NewMessagesManager(s *Store, pm *PeerManager) *MessagesManager {
	mm := &MessagesManager{
		s: s,
	}

	pm.RegisterNewPeerCallback(mm.newPeerCallback)

	return mm
}

func (mm *MessagesManager) newPeerCallback(peer *Peer) {
	log.Debugf("newPeerCallback %+v", peer)
}
