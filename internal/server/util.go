package server

import (
	"encoding/json"
	"os"
)

// JSONConfig is used to parse json config file
type JSONConfig struct {
	Address          *string `json:"address"`
	StoreFile        *string `json:"store_file"`
	Key              *string `json:"key"`
	CryptoKey        *string `json:"crypto_key"`
	DatabaseDSN      *string `json:"database_dsn"`
	StoreIntervalStr *string `json:"store_interval"`
	Restore          *bool   `json:"restore"`
}

// ReadJSONConfig parses config file and returns parsed data in struct
func ReadJSONConfig(cf string) (JSONConfig, error) {
	buf, err := os.ReadFile(cf)
	if err != nil {
		return JSONConfig{}, err
	}

	var jcfg JSONConfig
	err = json.Unmarshal(buf, &jcfg)
	if err != nil {
		return JSONConfig{}, err
	}

	return jcfg, nil
}
