package main

import (
	
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"encoding/pem"

	"os"

	log "github.com/sirupsen/logrus"
)


func CreateRoutes(store *Store){

	// gets users own info
	i, _ := store.MyInfo()

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)

	r.Get("/self", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")


		//create the pem object to perform encoding
		block := &pem.Block{
			Type: "MESSAGE",
			Bytes: []byte(i.PublicKey),
		}
	
		// writes human readable public key to page
		w.Write(pem.EncodeToMemory(block))
		fmt.Println()
	})
	
	//parse through messages that pertain to certain user
	r.Get("/messages", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("your messages go here"))
	})

	// listen
	http.ListenAndServe(":3000", r)

}
