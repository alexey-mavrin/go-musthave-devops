package server

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func p(s string) *string {
	return &s
}

var theTruth = true

func TestReadJSONConfig(t *testing.T) {
	type args struct {
		cf string
	}
	tests := []struct {
		name    string
		args    args
		want    JSONConfig
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "fully configured file",
			args: args{
				cf: "testdata/1.json",
			},
			want: JSONConfig{
				Address:          p("l:1"),
				Restore:          &theTruth,
				StoreIntervalStr: p("33s"),
				StoreFile:        p("/path/to/file.db"),
				DatabaseDSN:      p("psql:l:1234"),
				CryptoKey:        p("/path/to/key.pem"),
				Key:              p("qwerty"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "partially configured file",
			args: args{
				cf: "testdata/2.json",
			},
			want: JSONConfig{
				Address:          p("l:1"),
				Restore:          &theTruth,
				StoreIntervalStr: p("33s"),
				StoreFile:        nil,
				DatabaseDSN:      nil,
				CryptoKey:        nil,
				Key:              nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "file with error",
			args: args{
				cf: "testdata/3.json",
			},
			want:    JSONConfig{},
			wantErr: assert.Error,
		},
		{
			name: "missed file",
			args: args{
				cf: "testdata/0.json",
			},
			want:    JSONConfig{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadJSONConfig(tt.args.cf)
			tt.wantErr(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadJSONConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
