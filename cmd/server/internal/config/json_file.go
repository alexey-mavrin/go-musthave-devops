package config

import (
	"encoding/json"
	"net"
	"os"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
)

// JSONConfig is used to parse json config file
type JSONConfig struct {
	Address          *string `json:"address"`
	StoreFile        *string `json:"store_file"`
	Key              *string `json:"key"`
	CryptoKey        *string `json:"crypto_key"`
	DatabaseDSN      *string `json:"database_dsn"`
	StoreIntervalStr *string `json:"store_interval"`
	TrustedSubnetStr *string `json:"trusted_subnet"`
	Restore          *bool   `json:"restore"`
}

// ReadJSONConfig parses config file and returns parsed data in struct
func (b *Builder) ReadJSONConfig() *Builder {
	var cf string
	if b.flags.configFile.Set {
		cf = *b.flags.configFile.Value
	} else if b.envVars.ConfigFile != nil {
		cf = *b.envVars.ConfigFile
	} else {
		return b
	}
	buf, err := os.ReadFile(cf)
	if err != nil {
		b.err = err
		return b
	}

	err = json.Unmarshal(buf, &b.jsonConfig)
	if err != nil {
		b.err = err
		return b
	}

	return b
}

// MergeJSONConfig merges values set from config file into partial config
func (b *Builder) MergeJSONConfig() *Builder {
	common.CopyIfNotNil(&b.partial.Address, b.jsonConfig.Address)
	common.CopyIfNotNil(&b.partial.StoreFile, b.jsonConfig.StoreFile)
	common.CopyIfNotNil(&b.partial.Key, b.jsonConfig.Key)
	common.CopyIfNotNil(&b.partial.CryptoKey, b.jsonConfig.CryptoKey)
	common.CopyIfNotNil(&b.partial.DatabaseDSN, b.jsonConfig.DatabaseDSN)

	if b.jsonConfig.StoreIntervalStr != nil {
		storeInterval, err := time.ParseDuration(*b.jsonConfig.StoreIntervalStr)
		if err != nil {
			b.err = err
			return b
		}
		b.partial.StoreInterval = storeInterval
	}

	if b.jsonConfig.Restore != nil {
		b.partial.Restore = *b.jsonConfig.Restore
	}

	if b.jsonConfig.TrustedSubnetStr != nil {
		_, subnet, err := net.ParseCIDR(*b.jsonConfig.TrustedSubnetStr)
		if err != nil {
			b.err = err
			return b
		}
		b.partial.TrustedSubnet = subnet
	}

	return b
}
