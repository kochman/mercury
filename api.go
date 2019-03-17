package main

import (
	"fmt"
	"encoding/pem"
	"encoding/json"
	"net/http"
	"time"
	// "reflect"

	"github.com/gofrs/uuid"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
)

type API struct {
	box *packr.Box
	r http.Handler
}

func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	s, err := a.box.Find("static/index.html")
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Write(s)
}

func NewAPI(store *Store, box *packr.Box) *API {
	a := &API{
		box: box,
	}
	// gets users own info
	i, _ := store.MyInfo()

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)

	r.Get("/self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		//create the pem object to perform encoding
		block := &pem.Block{
			Type:  "MESSAGE",
			Bytes: []byte(i.PublicKey),
		}

		// writes human readable public key to page
		w.Write(pem.EncodeToMemory(block))
	})

	// MOCK DATA
	// TO DELETE
	
	msg := &EncryptedMessage{

		ID:			[]byte("1"),
		Sent:		time.Now(),
		Contents:	[]byte("test"),
	}
	msg2 := &EncryptedMessage{

		ID:			[]byte("2"),
		Sent:		time.Now(),
		Contents:	[]byte("test"),
	}



	store.AddEncryptedMessage(msg)
	store.AddEncryptedMessage(msg2)

	//DELETE ABOVE MOCK DATA


	r.Route("/api", func(r chi.Router){
		r.Get("/self", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")

			//create the pem object to perform encoding
			block := &pem.Block{
				Type:  "MESSAGE",
				Bytes: []byte(i.PublicKey),
			}

			// writes human readable public key to page
			w.Write(pem.EncodeToMemory(block))
		})
		// display messages that user has in decrypted store
		r.Get("/messages", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// get all the decrypted messages this person has
			a, _ := store.EncryptedMessages()
	
			// spacing out the json data
			output, err := json.MarshalIndent(a, "", " ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	
			//write the output to the page
			w.Write(output)
		})
	})

	

	//DO POST REQUEST HERE FOR SENDING MESSAGES
	r.Post("/send", func(w http.ResponseWriter, r *http.Request){
		err := r.ParseForm()
		if err != nil {
			// panic(err)
		}
		// v := r.Form.Get("to")
		to := r.Form.Get("message")
		
		u2, err := uuid.NewV4()
		if err != nil {
			// fmt.Fatalf("failed to generate UUID: %v", err)
		}
		fmt.Printf("generated Version 4 UUID %v", u2)

		msg := &EncryptedMessage{
			ID:		 	[]byte(u2.String()),
			Sent:	 	time.Now(),
			Contents:	[]byte(to),
		}

		// fmt.Println("type of ", reflect.TypeOf(u2))
		// nm.NewMessage(msg)
		store.AddEncryptedMessage(msg)
		//prints to web page
		// fmt.Fprintln(w,v)
		fmt.Fprintln(w,to)

		
	})


	// listen
	r.Get("/", a.IndexHandler)

	a.r = r

	return a
}

func (a *API) Run() {
	http.ListenAndServe(":3000", a.r)
}
