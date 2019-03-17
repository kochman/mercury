package main

import (
	"net/http"
	"time"

	// "reflect"

	"github.com/gobuffalo/packr/v2"
	log "github.com/sirupsen/logrus"

	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gofrs/uuid"
)

type API struct {
	box   *packr.Box
	r     http.Handler
	store *Store
}

func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	s, err := a.box.Find("static/home/home.html")
	w.Header().Set("content-type", "text/html")
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
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
	}
	if err != nil {
		w.WriteHeader(400)
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

type FrontendMessage struct {
	TargetUserID int
	Message      string
}

func (a *API) SendMessage(w http.ResponseWriter, r *http.Request) {
	v := &FrontendMessage{}
	err := json.NewDecoder(r.Body).Decode(v)

	if err != nil {
		w.WriteHeader(400)
		log.WithError(err).Error("Bad Request")

		w.Write([]byte("malformed request"))
		return
	}

	c, err := a.store.Contacts()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var sendTo *Contact
	for _, contact := range c {
		if contact.ID == v.TargetUserID {
			sendTo = contact
			break
		}
	}

	if sendTo == nil {
		w.WriteHeader(500)
		w.Write([]byte("invalid contact"))
		return
	}

	u2, err := uuid.NewV4()
	if err != nil {
		// fmt.Fatalf("failed to generate UUID: %v", err)
	}
	fmt.Printf("generated Version 4 UUID %v", u2)
	keyPair, err := PublicKeyFromBytes(sendTo.PublicKey)
	if err != nil {
		w.Write([]byte("Unable to get public key"))
		w.WriteHeader(500)
		return
	}
	content, err := keyPair.Sign("msg", v.Message)
	if err != nil {
		w.Write([]byte("unable to sign message"))
		w.WriteHeader(500)
		return
	}

	msg := &EncryptedMessage{
		ID:       u2.String(),
		Sent:     time.Now(),
		Contents: []byte(*content),
	}

	// fmt.Println("type of ", reflect.TypeOf(u2))
	// nm.NewMessage(msg)
	a.store.AddEncryptedMessage(msg)
	//prints to web page
	// fmt.Fprintln(w,v)
	fmt.Fprintln(w, msg)

}

func NewAPI(store *Store, box *packr.Box) *API {

	a := &API{
		box:   box,
		store: store,
	}
	// gets users own info
	i, _ := store.MyInfo()

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)

	// MOCK DATA
	// TO DELETE

	msg := &EncryptedMessage{
		ID:			"1",
		Sent:		time.Now(),
		Contents:	[]byte("test"),
	}
	msg2 := &EncryptedMessage{
		ID:			"2",
		Sent:		time.Now(),
		Contents:	[]byte("test"),
	}

	store.AddEncryptedMessage(msg)
	store.AddEncryptedMessage(msg2)

	//DELETE ABOVE MOCK DATA

	r.Method("GET", "/static/*", http.FileServer(box))

	r.Route("/api", func(r chi.Router) {
		r.Get("/self", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")

			//create the pem object to perform encoding
			block := &pem.Block{
				Type:  "PUBLIC KEY",
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
		r.Route("/contacts", func(r chi.Router) {
			r.Get("/all", a.GetContactsHandler)
			r.Post("/create", a.CreateContactHandler)
		})
	})

	//DO POST REQUEST HERE FOR SENDING MESSAGES
	r.Post("/send", a.SendMessage)

	// listen
	r.Get("/", a.IndexHandler)
	r.Route("/contacts", func(r chi.Router) {
		r.Get("/", a.AddressBookHandler)

	})

	a.r = r

	return a
}

func (a *API) Run() {
	http.ListenAndServe(":3000", a.r)
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
