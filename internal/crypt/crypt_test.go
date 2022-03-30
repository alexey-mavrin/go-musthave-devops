package crypt

import (
	"crypto/rsa"
	"os"
	"reflect"
	"testing"
)

func generateTestKeyPair(t *testing.T) Keys {
	tmpDir, err := os.MkdirTemp("", "keys")
	if err != nil {
		t.Fatal(err)
	}
	keys, err := generateKeyPair(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(keys.privateKeyFile)
		os.Remove(keys.publicKeyFile)
	})
	return keys
}

func TestReadKeys(t *testing.T) {
	keys := generateTestKeyPair(t)

	type args struct {
		privateKeyFile string
		publicKeyFile  string
	}
	type want struct {
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		publicErr  error
		privateErr error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct scenario",
			args: args{
				privateKeyFile: keys.privateKeyFile,
				publicKeyFile:  keys.publicKeyFile,
			},
			want: want{
				publicKey:  &keys.privateKey.PublicKey,
				privateKey: keys.privateKey,
				publicErr:  nil,
				privateErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubKey, err := ReadPublicKey(tt.args.publicKeyFile)
			if err != tt.want.publicErr {
				t.Errorf("ReadPublicKey() error = %v, want %v",
					err, tt.want.publicErr)
				return
			}
			if !reflect.DeepEqual(pubKey, tt.want.publicKey) {
				t.Errorf("ReadPublicKey() = %v, want %v",
					pubKey, tt.want.publicKey)
			}

			privKey, err := ReadPrivateKey(tt.args.privateKeyFile)
			if err != tt.want.privateErr {
				t.Errorf("ReadPrivateKey() error = %v, want %v",
					err, tt.want.privateErr)
				return
			}
			if !reflect.DeepEqual(privKey, tt.want.privateKey) {
				t.Errorf("ReadPrivateKey() = %v, want %v",
					pubKey, tt.want.privateKey)
			}
		})
	}
}
