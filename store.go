package main

import (
	"github.com/asdine/storm"
)

type Store struct {
	db *storm.DB
}

func NewStore() (*Store, error) {
	db, err := storm.Open("mercury.db")
	if err != nil {
		return nil, err
	}

	// initialize buckets for each struct we'll store
	for _, t := range []interface{}{
		&Contact{},
	} {
		if err := db.Init(t); err != nil {
			return nil, err
		}
	}

	s := &Store{
		db: db,
	}

	return s, nil
}

// Contacts

type Contact struct {
	ID int
	Name string
	PublicKey []byte
}

func (s *Store) Contacts() ([]*Contact, error) {
	var contacts []*Contact
	err := s.db.All(&contacts)
	return contacts, err
}

func (s *Store) Contact(id int) (*Contact, error) {
	var c *Contact
	err := s.db.One("ID", id, c)
	return c, err
}

func (s *Store) CreateContact(c *Contact) error {
	return s.db.Save(c)
}
