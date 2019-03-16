package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("starting Mercury...")
	log.SetLevel(log.DebugLevel)

	pm := NewPeerManager()

	// go this async in the future
	pm.Run()
}
