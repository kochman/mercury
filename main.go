package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("starting Mercury...")
	log.SetLevel(log.DebugLevel)

	store, err := NewStore()
	if err != nil {
		log.WithError(err).Error("unable to create store")
		os.Exit(1)
	}
	_ = store

	pm := NewPeerManager()

	// go this async in the future
	pm.Run()
}
