package crypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path"
)

// Keys is used to generate key pair
type Keys struct {
	privateKey     *rsa.PrivateKey
	privateKeyFile string
	publicKeyFile  string
}

// GenerateStoreKeys generates and stores key pair in the given directory.
// Directory must exist and have write access
func GenerateStoreKeys(keyDir string) error {
	_, err := generateKeyPair(keyDir)
	return err
}

func generateKeyPair(keyDir string) (Keys, error) {
	var keys Keys
	keys.privateKeyFile = path.Join(keyDir, "private.key")
	keys.publicKeyFile = path.Join(keyDir, "public.key")

	var err error
	keys.privateKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return keys, err
	}
	publicKey := keys.privateKey.PublicKey

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keys.privateKey),
	})

	err = os.WriteFile(keys.privateKeyFile, privateKeyPEM.Bytes(), 0644)
	if err != nil {
		return keys, err
	}

	var publicKeyPEM bytes.Buffer
	pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&publicKey),
	})
	err = os.WriteFile(keys.publicKeyFile, publicKeyPEM.Bytes(), 0644)
	if err != nil {
		return keys, err
	}
	return keys, nil
}
