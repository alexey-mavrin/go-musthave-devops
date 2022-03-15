package server_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
)

func ExampleHandler400() {
	var body bytes.Buffer
	req, _ := http.NewRequest("GET", "/nosuch", &body)
	res := httptest.NewRecorder()
	server.Handler400(res, req)
	fmt.Println(res.Code)
	// Output:
	// 400
}

func ExampleJSONMetricHandler() {
	var body bytes.Buffer

	body.Write([]byte(`{"id":"x100","type":"counter","delta":100}`))
	req, _ := http.NewRequest("POST", "/update/", &body)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	// set the value to retrieve later.
	server.JSONUpdateHandler(res, req)

	body.Reset()
	body.Write([]byte(`{"id":"x100","type":"counter"}`))
	req, _ = http.NewRequest("POST", "/value/", &body)
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	// retrieve the value.
	server.JSONMetricHandler(res, req)
	fmt.Println(res.Code)
	log.Print(res)
	// Output:
	// 200
}
