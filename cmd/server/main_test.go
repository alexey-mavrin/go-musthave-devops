package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "gauge OK",
			args: "/update/gauge/Alloc/1",
			want: want{code: 200},
		},
		{
			name: "counter OK",
			args: "/update/counter/RandomValue/1",
			want: want{code: 200},
		},
		{
			name: "counter no name 1",
			args: "/update/counter/",
			want: want{code: 404},
		},
		{
			name: "counter no name 2",
			args: "/update/counter//100",
			want: want{code: 404},
		},
		{
			name: "counter bad op",
			args: "/renew/counter/x/1",
			want: want{code: 404},
		},
		{
			name: "bad type",
			args: "/update/integer/x/1",
			want: want{code: 501},
		},
		{
			name: "counter bad counter val 1",
			args: "/update/counter/x/",
			want: want{code: 400},
		},
		{
			name: "counter bad counter val 2",
			args: "/update/counter/x/str",
			want: want{code: 400},
		},
		{
			name: "counter bad gauge val",
			args: "/update/gauge/x/1.2.3",
			want: want{code: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.args, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Handler)
			h.ServeHTTP(w, req)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()

			_, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

		})
	}
}
