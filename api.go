package main

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/packr/v2"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	r.Get("/", a.IndexHandler)
	r.Get("/self", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r)
		w.Write([]byte(i.Name))
	})
	fmt.Println("here")
	http.ListenAndServe(":3000", r)

}
