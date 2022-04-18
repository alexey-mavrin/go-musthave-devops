package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
)

// JSONConfig is used to parse json config file
type JSONConfig struct {
	ServerAddr        *string `json:"address"`
	Key               *string `json:"key"`
	CryptoKey         *string `json:"crypto_key"`
	PollIntervalStr   *string `json:"poll_interval"`
	ReportIntervalStr *string `json:"report_interval"`
	UseGRPC           *bool   `json:"use_grpc"`
	GRPCServer        *string `json:"grpc_server"`
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
	common.CopyIfNotNil(&b.partial.ServerAddr, b.jsonConfig.ServerAddr)
	common.CopyIfNotNil(&b.partial.Key, b.jsonConfig.Key)
	common.CopyIfNotNil(&b.partial.CryptoKey, b.jsonConfig.CryptoKey)
	common.CopyIfNotNil(&b.partial.GRPCServer, b.jsonConfig.GRPCServer)

	if b.jsonConfig.PollIntervalStr != nil {
		pollInterval, err := time.ParseDuration(*b.jsonConfig.PollIntervalStr)
		if err != nil {
			b.err = err
			return b
		}
		b.partial.PollInterval = pollInterval
	}

	if b.jsonConfig.ReportIntervalStr != nil {
		reportInterval, err := time.ParseDuration(*b.jsonConfig.ReportIntervalStr)
		if err != nil {
			b.err = err
			return b
		}
		b.partial.ReportInterval = reportInterval
	}

	if b.jsonConfig.UseGRPC != nil {
		b.partial.UseGRPC = *b.jsonConfig.UseGRPC
	}

	return b
}
