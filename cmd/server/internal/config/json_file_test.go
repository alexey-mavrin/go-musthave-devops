package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func addrStr(t string) *string {
	return &t
}

func addrBool(t bool) *bool {
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
					Address:          addrStr("l:1"),
					StoreFile:        addrStr("/path/to/file.db"),
					Key:              addrStr("qwerty"),
					CryptoKey:        addrStr("/path/to/key.pem"),
					DatabaseDSN:      addrStr("psql:l:1234"),
					StoreIntervalStr: addrStr("33s"),
					Restore:          addrBool(true),
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
			builder.envVars.ConfigFile = addrStr(tt.jsonConfig)
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
