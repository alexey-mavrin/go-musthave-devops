package server_test

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
)

type testValues struct {
	useJSON    bool
	supressLog bool
	body       string
	method     string
	path       string
	server     *httptest.Server
}

func newBenchTest(router http.Handler) *testValues {
	return &testValues{
		useJSON:    false,
		supressLog: true,
		method:     http.MethodGet,
		server:     httptest.NewServer(router),
	}
}

func (tv *testValues) setJSON(j bool) *testValues {
	tv.useJSON = j
	return tv
}

func (tv *testValues) supressLogging(l bool) *testValues {
	tv.supressLog = l
	return tv
}

func (tv *testValues) setMethod(m string) *testValues {
	tv.method = m
	return tv
}

func (tv *testValues) setPath(p string) *testValues {
	tv.path = p
	return tv
}

func (tv *testValues) setBody(b string) *testValues {
	tv.body = b
	return tv
}

func BenchmarkUpdate(b *testing.B) {
	router := server.Router()

	b.Run("updateCounter", func(b *testing.B) {
		tv := newBenchTest(router).
			setMethod(http.MethodPost).
			setPath("/update/counter/RandomValue/1").
			setJSON(false).
			supressLogging(true)
		for i := 0; i < b.N; i++ {
			statusCode, body, err := tv.doTestRequest()
			if err != nil {
				b.Fatal(err)
			}
			if statusCode != 200 {
				b.Fatal("request status is not 200, ", body)
			}
		}
	})

	b.Run("updataJSONCounter", func(b *testing.B) {
		tv := newBenchTest(router).
			setMethod(http.MethodPost).
			setPath("/update/").
			setJSON(true).
			supressLogging(true).
			setBody(`{"id":"xyz","type":"counter","delta":10}`)
		for i := 0; i < b.N; i++ {
			statusCode, body, err := tv.doTestRequest()
			if err != nil {
				b.Fatal(err)
			}
			if statusCode != 200 {
				b.Fatal("request status is not 200, ", body)
			}
		}
	})
}

func (tv testValues) doTestRequest() (int, string, error) {
	r := strings.NewReader(tv.body)
	req, err := http.NewRequest(tv.method, tv.server.URL+tv.path, r)
	if err != nil {
		return 0, "", err
	}

	if tv.useJSON {
		req.Header.Add("Content-Type", "application/json")
	}

	if tv.supressLog {
		log.SetOutput(ioutil.Discard)
	}

	resp, err := http.DefaultClient.Do(req)

	if tv.supressLog {
		log.SetOutput(os.Stdout)
	}

	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, string(respBody), nil
}
