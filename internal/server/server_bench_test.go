package server

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
)

func supressLog() {
	log.SetOutput(ioutil.Discard)
}

func resumeLog() {
	log.SetOutput(os.Stdout)
}

func BenchmarkUpdate(b *testing.B) {
	router := Router()
	ts := httptest.NewServer(router)

	b.Run("updateCounter", func(b *testing.B) {
		method := "POST"
		path := "/update/counter/RandomValue/1"
		j := false
		l := false
		for i := 0; i < b.N; i++ {
			r := strings.NewReader("")
			resp, body := testRequestBench(ts, method, path, r, j, l)
			if resp.StatusCode != 200 {
				log.Fatal("request status is not 200, ", body)
			}
			resp.Body.Close()
		}
	})

	b.Run("updataJSONCounter", func(b *testing.B) {
		method := "POST"
		path := "/update/"
		j := true
		l := false
		for i := 0; i < b.N; i++ {
			r := strings.NewReader(`{"id":"xyz","type":"counter","delta":10}`)
			resp, body := testRequestBench(ts, method, path, r, j, l)
			if resp.StatusCode != 200 {
				log.Fatal("request status is not 200, ", body)
			}
			resp.Body.Close()
		}
	})
}

func testRequestBench(ts *httptest.Server, method, path string, r io.Reader, useJSON bool, doLog bool) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, r)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	if useJSON {
		req.Header.Add("Content-Type", "application/json")
	}

	if !doLog {
		supressLog()
	}
	resp, err := http.DefaultClient.Do(req)
	if !doLog {
		resumeLog()
	}

	if err != nil {
		log.Fatal("http.DefaultClient.Do: ", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll: ", err)
	}

	return resp, string(respBody)
}
