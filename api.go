package main

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"encoding/json"
	"encoding/pem"
)

type API struct {
	box   *packr.Box
	store *Store
}

func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	s, err := a.box.Find("static/home/home.html")
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Write(s)
}

func (a *API) AddressBookHandler(w http.ResponseWriter, r *http.Request) {
	s, err := a.box.Find("static/addressbook/addressbook.html")
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Write(s)
}

func (a *API) GetContactsHandler(w http.ResponseWriter, r *http.Request) {
	ret, err := a.store.Contacts()
	if err != nil {
		log.WithError(err)
		w.WriteHeader(500)
		w.Write([]byte("Unable to read contacts"))
		return
	}

	WriteJSON(w, ret)

}

func (a *API) CreateContactHandler(w http.ResponseWriter, r *http.Request) {
	v := &Contact{}
	err := json.NewDecoder(r.Body).Decode(v)
	if v.Name == "" {
		w.WriteHeader(403)
		w.Write([]byte("Bad Request"))
	}
	if err != nil {
		w.WriteHeader(403)
		log.Error("Bad Request")
		w.Write([]byte("Bad request"))
		return
	}
	v.ID = 0
	err = a.store.CreateContact(v)
	if err != nil {
		w.WriteHeader(500)
		log.WithError(err)
		w.Write([]byte("Failed to create contact"))
		return
	}
}

func CreateRoutes(store *Store, box *packr.Box) {

	a := API{
		box:   box,
		store: store,
	}
	// gets users own info
	i, _ := store.MyInfo()

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)

	r.Get("/self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//create the pem object to perform encoding
		block := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: []byte(i.PublicKey),
		}

		// writes human readable public key to page
		w.Write(pem.EncodeToMemory(block))
	})

	r.Method("GET", "/static/*", http.FileServer(box))

	//parse through messages that pertain to certain user
	r.Get("/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("your messages go here"))
	})

	// listen
	r.Get("/", a.IndexHandler)
	r.Route("/contacts", func(r chi.Router) {
		r.Get("/", a.AddressBookHandler)
		r.Get("/all", a.GetContactsHandler)
		r.Post("/create", a.CreateContactHandler)

	})

	http.ListenAndServe(":3000", r)

}

// WriteJSON writes the data as JSON.
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}
