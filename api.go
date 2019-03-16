package main

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/packr/v2"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"encoding/pem"

	"os"

	log "github.com/sirupsen/logrus"
)

type API struct {
	box *packr.Box
}

func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	s, err := a.box.Find("static/index.html")
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Write(s)
}

func CreateRoutes(store *Store, box *packr.Box) {

	a := API{
		box: box,
	}
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
<<<<<<< HEAD
=======

	// listen
	r.Get("/", a.IndexHandler)

>>>>>>> 882fe05523e033503abd515a93c7ae0c437f7b59
	http.ListenAndServe(":3000", r)

}
