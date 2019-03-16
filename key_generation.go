/*
 * Genarate rsa keys.
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
)

// KeyPair generates a private key public key pair and allows signing
type KeyPair struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// NewKeyPair generates a new key pair
func NewKeyPair() (*KeyPair, error) {
	retVal := &KeyPair{}

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	retVal.privateKey = key
	retVal.PublicKey = &key.PublicKey

	return retVal, nil
}

// PublicKeyFromBytes returns a key pair with a public key from bytes
func PublicKeyFromBytes(key []byte) (*KeyPair, error) {
	var ret rsa.PublicKey
	_, err := asn1.Unmarshal(key, &ret)
	if err != nil {
		return nil, err
	}
	return &KeyPair{PublicKey: &ret}, nil
}

// KeyPairFromBytes returns a key pair with a public and private key from bytes
func KeyPairFromBytes(privatekey []byte) (*KeyPair, error) {
	var ret rsa.PrivateKey
	_, err := asn1.Unmarshal(privatekey, &ret)
	if err != nil {
		return nil, err
	}

	return &KeyPair{privateKey: &ret, PublicKey: &ret.PublicKey}, nil
}

// PublicKeyAsBytes returns the public key as bytes
func (key *KeyPair) PublicKeyAsBytes() ([]byte, error) {
	return asn1.Marshal(*key.PublicKey)
}

// PrivateKeyAsBytes returns the private key as a byte array
func (key *KeyPair) PrivateKeyAsBytes() ([]byte, error) {
	return asn1.Marshal(*key.privateKey)
}

// Sign allows you to sign text with the public key
func (key *KeyPair) Sign(label string, text string) (*string, error) {
	rng := rand.Reader
	ret, err := rsa.EncryptOAEP(sha256.New(), rng, key.PublicKey, []byte(text), []byte(label))
	if err != nil {
		return nil, err
	}
	x := string(ret)

	return &x, nil
}

// UnSign unsigns the text with the private key
func (key *KeyPair) UnSign(label string, text string) (*string, error) {
	if key.privateKey == nil {
		return nil, errors.New("no private key")
	}
	reader := rand.Reader

	ret, err := rsa.DecryptOAEP(sha256.New(), reader, key.privateKey, []byte(text), []byte(label))
	if err != nil {
		return nil, err
	}
	x := string(ret)
	return &x, nil

}
