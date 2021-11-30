package main

import (
	"net/http"

	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
)

func main() {
	r := server.Router()
	http.ListenAndServe(":8080", r)
}
