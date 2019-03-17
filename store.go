package main

import (
	"errors"

	"github.com/asdine/storm"
)

type Store struct {
	db *storm.DB
}

var (
	ErrNotFound = errors.New("not found")
)

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
	ID        int `storm:"id,increment"`
	Name      string
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

// MyInfo stores the info for the current user
type MyInfo struct {
	ID         int
	Name       string
	PublicKey  []byte
	PrivateKey []byte
}

func (s *Store) SetMyInfo(info *MyInfo) error {
	return s.db.Save(info)
}

func (s *Store) MyInfo() (*MyInfo, error) {
	info := MyInfo{}
	err := s.db.One("ID", 1, &info)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &info, nil
}
