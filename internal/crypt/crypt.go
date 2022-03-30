package crypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"hash"
	"io"
	"os"
)

// ReadPublicKey reads public key from file
func ReadPublicKey(file string) (*rsa.PublicKey, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(buf)
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// ReadPrivateKey read private key from file
func ReadPrivateKey(file string) (*rsa.PrivateKey, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(buf)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return key, nil
}

// EncryptOAEP encrypts long message
func EncryptOAEP(hash hash.Hash,
	random io.Reader,
	public *rsa.PublicKey,
	msg []byte,
	label []byte,
) ([]byte, error) {
	msgLen := len(msg)
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

// DecryptOAEP descrypts long message
func DecryptOAEP(hash hash.Hash,
	random io.Reader,
	private *rsa.PrivateKey,
	msg []byte,
	label []byte,
) ([]byte, error) {
	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, private, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
