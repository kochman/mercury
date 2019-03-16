package main

import (
	
	// "os"

	// log "github.com/sirupsen/logrus"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)


type API struct{
	handler http.Handler
}


func CreateRoutes(store *Store){



	i, _ := store.MyInfo()
	
	
	fmt.Print("this is a ",i.Name)

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)

	r.Get("/self", func(w http.ResponseWriter, r *http.Request){
		fmt.Print(r)
		w.Write([]byte(i.Name))
	})

	http.ListenAndServe(":3000", r)

}
