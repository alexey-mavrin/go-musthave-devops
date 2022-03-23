package main

import (
	"testing"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
	"github.com/stretchr/testify/assert"
)

func Test_setServerArgs(t *testing.T) {
	type want struct {
		Address       string
		StoreInterval time.Duration
		StoreFile     string
		Restore       bool
		Key           string
		DatabaseDSN   string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "setServerArgs sets default server config",
			want: want{
				Address:       "localhost:8080",
				StoreInterval: 300 * time.Second,
				StoreFile:     "/tmp/devops-metrics-db.json",
				Restore:       true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setServerArgs()
		})
		assert.EqualValues(t, tt.want, server.Config)
	}
}
