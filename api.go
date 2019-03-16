package main

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/packr/v2"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"encoding/pem"
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

	r.Get("/self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//create the pem object to perform encoding
		block := &pem.Block{
			Type:  "MESSAGE",
			Bytes: []byte(i.PublicKey),
		}

		// writes human readable public key to page
		w.Write(pem.EncodeToMemory(block))
		fmt.Println()
	})

	//parse through messages that pertain to certain user
	r.Get("/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("your messages go here"))
	})

	// listen
	r.Get("/", a.IndexHandler)

	http.ListenAndServe(":3000", r)

}
