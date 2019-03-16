/*
 * Genarate rsa keys.
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
)

// KeyPair generates a private key public key pair and allows signing
type KeyPair struct {
	privateKey *rsa.PrivateKey
}

// NewKeyPair generates a new key pair
func NewKeyPair() (*KeyPair, error) {
	retVal := &KeyPair{}

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	retVal.privateKey = key

	return retVal, nil
}

// Sign allows you to sign text with the public key
func (key *KeyPair) Sign(label string, text string) (*string, error) {
	rng := rand.Reader
	ret, err := rsa.EncryptOAEP(sha256.New(), rng, &key.privateKey.PublicKey, []byte(text), []byte(label))
	if err != nil {
		return nil, err
	}
	x := string(ret)

	return &x, nil
}

// UnSign unsigns the text with the private key
func (key *KeyPair) UnSign(label string, text string) (*string, error) {
	reader := rand.Reader

	ret, err := rsa.DecryptOAEP(sha256.New(), reader, key.privateKey, []byte(text), []byte(label))
	if err != nil {
		return nil, err
	}
	x := string(ret)
	return &x, nil

}

func main() {
	key, _ := NewKeyPair()
	ret, err := key.Sign("msg", "This is a test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*ret)
	us, _ := key.UnSign("msg", *ret)
	fmt.Println(*us)

}

// 	reader := rand.Reader
// 	bitSize := 144

// 	key, err := rsa.GenerateKey(reader, bitSize)
// 	checkError(err)

// 	publicKey := key.PublicKey

// 	savePEMKey("private.pem", key)

// 	asn1Bytes, err := asn1.Marshal(publicKey)
// 	entropy := asn1Bytes

// 	fmt.Println(entropy)
// 	mnemonic, err := bip39.NewMnemonic(entropy)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(mnemonic)
// 	// converting mnemonic back to a byte array + checksum
// 	arr, _ := bip39.MnemonicToByteArray(mnemonic)

// 	fmt.Print(arr)
// 	fmt.Print("ASDF")
// 	savePublicPEMKey("public.pem", publicKey)
// }

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
