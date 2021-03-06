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

func (s *Store) PublicKeyAdded(pubkey []byte) (bool, error) {
	c := &Contact{}
	err := s.db.One("PublicKey", pubkey, c)
	if err == storm.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
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
	ID       string
	Sent     time.Time
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
	ID       string
	Sent     time.Time
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

// return all messages IDs that have been saved so we don't have to process them again
func (s *Store) ProcessedMessageIDs() (map[string]struct{}, error) {
	emsgs, err := s.EncryptedMessages()
	if err != nil {
		return nil, err
	}

	dmsgs, err := s.DecryptedMessages()
	if err != nil {
		return nil, err
	}

	ids := map[string]struct{}{}
	for _, msg := range emsgs {
		ids[msg.ID] = struct{}{}
	}
	for _, msg := range dmsgs {
		ids[msg.ID] = struct{}{}
	}

	return ids, nil
}
