package main

import (
	"fmt"
	"net/http"
)

// Handler — обработчик запроса.
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
