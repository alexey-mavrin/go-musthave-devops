package server

import (
	"bytes"
	crand "crypto/rand"
	"crypto/sha256"
	"io/ioutil"
	"net/http"

	"github.com/alexey-mavrin/go-musthave-devops/internal/crypt"
)

// DecryptBody is chi middleware function used to decrypt the received body
func DecryptBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		r2 := r.Clone(r.Context())
		if privateServerKey != nil {
			body, _ := ioutil.ReadAll(r.Body)
			decryptedBytes, err := crypt.DecryptOAEP(
				sha256.New(),
				crand.Reader,
				privateServerKey,
				body,
				nil)
			if err != nil {
				panic(err)
			}
			r2.Body = ioutil.NopCloser(bytes.NewReader(decryptedBytes))

		}
		next.ServeHTTP(rw, r2)
	})
}
