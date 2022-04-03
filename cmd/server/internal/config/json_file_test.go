package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func addr[T any](t T) *T {
	return &t
}

func TestBuilder_ReadJSONConfig(t *testing.T) {
	tests := []struct {
		want       *Builder
		wantErr    assert.ErrorAssertionFunc
		name       string
		jsonConfig string
	}{
		{
			name:       "simple json config",
			jsonConfig: "testdata/1.json",
			want: &Builder{
				jsonConfig: JSONConfig{
					Address:          addr("l:1"),
					StoreFile:        addr("/path/to/file.db"),
					Key:              addr("qwerty"),
					CryptoKey:        addr("/path/to/key.pem"),
					DatabaseDSN:      addr("psql:l:1234"),
					StoreIntervalStr: addr("33s"),
					Restore:          addr(true),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:       "invalid json file",
			jsonConfig: "testdata/3.json",
			wantErr:    assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder()
			builder.MergeDefaults()
			builder.envVars.ConfigFile = addr(tt.jsonConfig)
			got := builder.ReadJSONConfig()
			if tt.want != nil {
				assert.EqualValues(t,
					tt.want.jsonConfig,
					got.jsonConfig,
				)
			}
			tt.wantErr(t, got.err)
		})
	}
}
