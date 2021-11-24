package main

import (
	"internal/server"
	"net/http"
)

func main() {
	r := server.Router()
	http.ListenAndServe(":8080", r)
}
