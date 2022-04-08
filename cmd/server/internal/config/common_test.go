package config

import (
	"testing"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	tests := []struct {
		want    *Builder
		wantErr assert.ErrorAssertionFunc
		name    string
	}{
		{
			name: "get new builder struct with defaults",
			want: &Builder{
				defaultConfig: server.ConfigType{
					Address:       "localhost:8080",
					StoreFile:     "/tmp/devops-metrics-db.json",
					Restore:       true,
					StoreInterval: 300 * time.Second,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBuilder()
			assert.EqualValues(t,
				tt.want.defaultConfig,
				got.defaultConfig,
			)
			tt.wantErr(t, got.err)
		})
	}
}

func TestBuilder_MergeDefaults(t *testing.T) {
	tests := []struct {
		want    *Builder
		wantErr assert.ErrorAssertionFunc
		name    string
	}{
		{
			name: "merge default fields",
			want: &Builder{
				partial: server.ConfigType{
					Address:       "localhost:8080",
					StoreFile:     "/tmp/devops-metrics-db.json",
					Restore:       true,
					StoreInterval: 300 * time.Second,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBuilder().MergeDefaults()
			assert.EqualValues(t,
				tt.want.partial,
				got.partial,
			)
			tt.wantErr(t, got.err)
		})
	}
}

func TestBuilder_Final(t *testing.T) {
	tests := []struct {
		want       *server.ConfigType
		wantErr    assert.ErrorAssertionFunc
		name       string
		jsonConfig string
	}{
		{
			name: "simple test with defaults only",
			want: &server.ConfigType{
				Address:       "localhost:8080",
				StoreFile:     "/tmp/devops-metrics-db.json",
				Restore:       true,
				StoreInterval: 300 * time.Second,
			},
			wantErr: assert.NoError,
		},
		{
			name:       "some values from defaults, others from json",
			jsonConfig: "testdata/2.json",
			want: &server.ConfigType{
				Address:       "l:1",
				StoreFile:     "/tmp/devops-metrics-db.json",
				Key:           "",
				CryptoKey:     "",
				DatabaseDSN:   "",
				StoreInterval: 33 * time.Second,
				Restore:       false,
			},
			wantErr: assert.NoError,
		},
		{
			name:       "missed json file",
			jsonConfig: "testdata/100500.json",
			wantErr:    assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder()
			if tt.jsonConfig != "" {
				builder.envVars.ConfigFile = addrStr(tt.jsonConfig)
			}
			builder.MergeDefaults()
			builder.ReadJSONConfig().MergeJSONConfig()
			got := builder.Final()
			if tt.want != nil {
				assert.EqualValues(t,
					*tt.want,
					got,
				)
			}
			tt.wantErr(t, builder.err)
		})
	}
}
