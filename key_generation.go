/*
 * Genarate rsa keys.
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/tyler-smith/go-bip39"
)

func main() {
	reader := rand.Reader
	bitSize := 144

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	savePEMKey("private.pem", key)

	asn1Bytes, err := asn1.Marshal(publicKey)
	entropy := asn1Bytes

	fmt.Println(entropy)
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(mnemonic)
	// converting mnemonic back to a byte array + checksum
	arr, _ := bip39.MnemonicToByteArray(mnemonic)

	fmt.Print(arr)
	fmt.Print("ASDF")
	savePublicPEMKey("public.pem", publicKey)
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
