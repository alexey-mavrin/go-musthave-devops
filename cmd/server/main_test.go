package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type want struct {
		code int
		body []string
	}
	r := Router()
	ts := httptest.NewServer(r)
	defer ts.Close()
	tests := []struct {
		name   string
		method string
		args   string
		want   want
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
			name:   "counter bad gauge val",
			method: "POST",
			args:   "/update/gauge/x/1.2.3",
			want:   want{code: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, body := testRequest(t, ts, tt.method, tt.args)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			for _, s := range tt.want.body {
				assert.Contains(t, body, s)
			}
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
