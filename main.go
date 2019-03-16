package main

import (
	"fmt"
	"os"

	"github.com/gobuffalo/packr/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("starting Mercury...")
	log.SetLevel(log.DebugLevel)

	box := packr.New("static", "./static")

	store, err := NewStore()
	if err != nil {
		log.WithError(err).Error("unable to create store")
		os.Exit(1)
	}
	_ = store

	i, err := store.MyInfo()
	_ = i

	if err != nil {
		if err == ErrNotFound {
			key, err := NewKeyPair()
			if err != nil {
				log.WithError(err)
				return
			}
			privKey, _ := key.PrivateKeyAsBytes()
			pubKey, _ := key.PublicKeyAsBytes()

			myInfo := MyInfo{
				ID:         1,
				Name:       "joey",
				PrivateKey: privKey,
				PublicKey:  pubKey,
			}
			store.SetMyInfo(&myInfo)
		}
	}
	i, err = store.MyInfo()
	if err != nil {
		fmt.Println(err)
	}

	api := NewAPI(store, box)
	go api.Run()
	// fmt.Println(i)

	msgs, err := store.EncryptedMessages()
	if err != nil {
		log.WithError(err).Error("unable to get messages")
	}
	fmt.Println(msgs)

	pm := NewPeerManager()

	mm := NewMessagesManager(store, pm)
	_ = mm

	pm.Run()
}
