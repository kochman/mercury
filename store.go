package main

import (
	"errors"
	"time"

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
		&MyInfo{},
		&EncryptedMessage{},
		&DecryptedMessage{},
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
	ID        int
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
	info := &MyInfo{}
	err := s.db.One("ID", 1, info)
	if err == storm.ErrNotFound {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return info, nil
}

// Messages

type EncryptedMessage struct {
	ID []byte
	Sent time.Time
	Contents []byte
}

func (s *Store) EncryptedMessages() ([]*EncryptedMessage, error) {
	var msgs []*EncryptedMessage
	err := s.db.All(&msgs)
	return msgs, err
}

func (s *Store) AddEncryptedMessage(msg *EncryptedMessage) error {
	return s.db.Save(msg)
}

type DecryptedMessage struct {
	ID []byte
	Sent time.Time
	Contents string
}

func (s *Store) DecryptedMessages() ([]*DecryptedMessage, error) {
	var msgs []*DecryptedMessage
	err := s.db.All(&msgs)
	return msgs, err
}

func (s *Store) AddDecryptedMessage(msg *DecryptedMessage) error {
	return s.db.Save(msg)
}

