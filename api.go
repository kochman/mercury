package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func CreateRoutes(store *Store) {

	// gets users own info
	i, _ := store.MyInfo()

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)

	r.Get("/self", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r)
		w.Write([]byte(i.Name))
	})
	fmt.Println("here")
	http.ListenAndServe(":3000", r)

}
