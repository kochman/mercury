/*
 * Genarate rsa keys.
 */

package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// KeyPair generates a private key public key pair and allows signing
type KeyPair struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// NewKeyPair generates a new key pair
func NewKeyPair(file string) (*KeyPair, error) {
	retVal := &KeyPair{}

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	retVal.privateKey = key
	retVal.PublicKey = &key.PublicKey
	savePEMKey(file, retVal.privateKey)

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

// LoadKeyPair loads a key pair from a file
func LoadKeyPair(file string) (*KeyPair, error) {
	retVal := &KeyPair{}

	key, err := readPEMKey(file)
	if err != nil {
		return nil, err
	}
	retVal.privateKey = key
	retVal.PublicKey = &key.PublicKey
	return retVal, nil
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

func readPEMKey(fileName string) (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	size := pemfileinfo.Size()

	pembytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		return nil, err
	}

	return privateKeyImported, nil
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
