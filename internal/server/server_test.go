package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	/*
		type args struct {
			w http.ResponseWriter
			r *http.Request
		}
	*/
	type want struct {
		body []string
		code int
	}
	router := Router()
	ts := httptest.NewServer(router)
	defer ts.Close()
	tests := []struct {
		name    string
		method  string
		args    string
		body    string
		want    want
		useJSON bool
	}{
		{
			name:   "gauge OK",
			method: "POST",
			args:   "/update/gauge/Alloc/2128506",
			want:   want{code: 200},
		},
		{
			name:   "counter OK",
			method: "POST",
			args:   "/update/counter/RandomValue/100500",
			want:   want{code: 200},
		},
		{
			name:   "Get all",
			method: "GET",
			args:   "/",
			want: want{
				code: 200,
				body: []string{
					"Alloc 2.128506e+06",
					"RandomValue 100500",
				},
			},
		},
		{
			name:   "Get counter",
			method: "GET",
			args:   "/value/counter/RandomValue",
			want: want{
				code: 200,
				body: []string{"100500"},
			},
		},
		{
			name:   "counter OK",
			method: "POST",
			args:   "/update/counter/RandomValue/1",
			want:   want{code: 200},
		},
		{
			name:   "Get counter",
			method: "GET",
			args:   "/value/counter/RandomValue",
			want: want{
				code: 200,
				body: []string{"100501"},
			},
		},
		{
			name:   "Get gauge",
			method: "GET",
			args:   "/value/gauge/Alloc",
			want: want{
				code: 200,
				body: []string{"2.128506e+06"},
			},
		},
		{
			name:   "counter no name 1",
			method: "POST",
			args:   "/update/counter/",
			want:   want{code: 404},
		},
		{
			name:   "counter no name 2",
			method: "POST",
			args:   "/update/counter//100",
			want:   want{code: 404},
		},
		{
			name:   "counter bad op",
			method: "POST",
			args:   "/renew/counter/x/1",
			want:   want{code: 404},
		},
		{
			name:   "bad type",
			method: "POST",
			args:   "/update/integer/x/1",
			want:   want{code: 501},
		},
		{
			name:   "counter bad counter val 1",
			method: "POST",
			args:   "/update/counter/x/",
			want:   want{code: 400},
		},
		{
			name:   "counter bad counter val 2",
			method: "POST",
			args:   "/update/counter/x/str",
			want:   want{code: 400},
		},
		{
			name:   "gauge bad gauge val",
			method: "POST",
			args:   "/update/gauge/x/1.2.3",
			want:   want{code: 400},
		},
		{
			name:   "get unknown gauge",
			method: "GET",
			args:   "/value/gauge/nosuchname",
			want:   want{code: 404},
		},
		{
			name:   "get unknown counter",
			method: "GET",
			args:   "/value/counter/nosuchname",
			want:   want{code: 404},
		},
		{
			name:   "get unknown type",
			method: "GET",
			args:   "/value/float/name",
			want:   want{code: 400},
		},
		{
			name:    "update json counter",
			method:  "POST",
			args:    "/update/",
			useJSON: true,
			body:    `{"id":"xyz","type":"counter","delta":10}`,
			want:    want{code: 200},
		},
		{
			name:    "get json counter",
			method:  "POST",
			args:    "/value/",
			useJSON: true,
			body:    `{"id":"xyz","type":"counter"}`,
			want: want{
				code: 200,
				body: []string{
					`{"id":"xyz","type":"counter","delta":10}`,
				},
			},
		},
		{
			name:    "update json unknown type",
			method:  "POST",
			args:    "/update/",
			useJSON: true,
			body:    `{"id":"xyz","type":"nosuch","value":10}`,
			want:    want{code: 501},
		},
		{
			name:    "update json gauge",
			method:  "POST",
			args:    "/update/",
			useJSON: true,
			body:    `{"id":"xyz","type":"gauge","value":10}`,
			want:    want{code: 200},
		},
		{
			name:    "get json gauge",
			method:  "POST",
			args:    "/value/",
			useJSON: true,
			body:    `{"id":"xyz","type":"gauge"}`,
			want: want{
				code: 200,
				body: []string{
					`{"id":"xyz","type":"gauge","value":10}`,
				},
			},
		},

		{
			name:    "update json counter wrong content type",
			method:  "POST",
			args:    "/update/",
			useJSON: false,
			body:    `{"id":"xyz","type":"counter","delta":10}`,
			want:    want{code: 400},
		},
		{
			name:    "value json counter wrong content type",
			method:  "POST",
			args:    "/value/",
			useJSON: false,
			body:    `{"id":"xyz","type":"counter"}`,
			want:    want{code: 400},
		},
		{
			name:    "update json counter wrong json",
			method:  "POST",
			args:    "/update/",
			useJSON: true,
			body:    `{"id":"xyz","type":"counter","delta:10}`,
			want:    want{code: 400},
		},
		{
			name:    "value json counter wrong json",
			method:  "POST",
			args:    "/value/",
			useJSON: true,
			body:    `{"id":"xyz","type":}`,
			want:    want{code: 400},
		},
		{
			name:    "update json counter no ID",
			method:  "POST",
			args:    "/update/",
			useJSON: true,
			body:    `{"type":"counter","delta":10}`,
			want:    want{code: 400},
		},
		{
			name:    "value json counter no ID",
			method:  "POST",
			args:    "/value/",
			useJSON: true,
			body:    `{"type":"counter"}`,
			want:    want{code: 400},
		},
		{
			name:    "value json counter unknown ID",
			method:  "POST",
			args:    "/value/",
			useJSON: true,
			body:    `{"type":"counter","id":"qqqqqqq"}`,
			want:    want{code: 404},
		},
		{
			name:    "value json gauge unknown ID",
			method:  "POST",
			args:    "/value/",
			useJSON: true,
			body:    `{"type":"gauge","id":"qqqqqqq"}`,
			want:    want{code: 404},
		},
		{
			name:    "value json unknown type",
			method:  "POST",
			args:    "/value/",
			useJSON: true,
			body:    `{"type":"someother","id":"qqqqqqq"}`,
			want:    want{code: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := strings.NewReader(tt.body)
			resp, body := testRequest(t, ts, tt.method, tt.args, r, tt.useJSON)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			for _, s := range tt.want.body {
				assert.Contains(t, body, s)
			}
		})
	}

	tmpDir := os.TempDir()
	os.Mkdir(tmpDir, 0755)
	tmpFile := tmpDir + "/1.json"
	Config.StoreFile = tmpFile

	resp, _ := testRequest(t, ts, "POST", "/update/counter/c123/123", nil, false)
	resp.Body.Close()

	tmpF, _ := os.OpenFile(tmpFile, os.O_RDONLY, 0)

	tmpBuf, _ := io.ReadAll(tmpF)
	assert.Contains(t, string(tmpBuf), `"c123":123`)
	t.Logf(string(tmpBuf))

	statistics.Counters["c123"] = 0
	loadStats()
	assert.Equal(t, statistics.Counters["c123"], int64(123))
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, r io.Reader, useJSON bool) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, r)
	require.NoError(t, err)

	if useJSON {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
