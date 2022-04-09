// Package config provides config types for server
package config

import (
	"log"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
)

// Builder is the base for building server config
type Builder struct {
	flags         flags
	jsonConfig    JSONConfig
	envVars       envVarConfig
	defaultConfig server.ConfigType
	err           error
	partial       server.ConfigType
}

// NewBuilder returns a pointer to the Builder struct filled with default values
func NewBuilder() *Builder {
	b := Builder{
		defaultConfig: server.ConfigType{
			Address:       "localhost:8080",
			StoreInterval: time.Second * 300,
			StoreFile:     "/tmp/devops-metrics-db.json",
			Restore:       true,
		},
	}
	return &b
}

// MergeDefaults copies parameter from default values to the partial config
func (b *Builder) MergeDefaults() *Builder {
	b.partial.Address = b.defaultConfig.Address
	b.partial.StoreInterval = b.defaultConfig.StoreInterval
	b.partial.StoreFile = b.defaultConfig.StoreFile
	b.partial.Restore = b.defaultConfig.Restore

	return b
}

// Err returns err
func (b Builder) Err() error {
	return b.err
}

// ReportEnvVars prints parsed env vars
func (b *Builder) ReportEnvVars() *Builder {
	log.Printf("server is invoked with ENV %+v", b.envVars)
	return b
}

// ReportFlags prints passed flags
func (b *Builder) ReportFlags() *Builder {
	log.Printf("server is invoked with flags address %v store interval %v store file %v restore %v database %v",
		b.flags.address,
		b.flags.storeInterval,
		b.flags.storeFile,
		b.flags.restore,
		b.flags.databaseDSN,
	)

	return b
}

// ReportJSONConfig prints values passed via config file
func (b *Builder) ReportJSONConfig() *Builder {
	log.Printf("config file values: %+v", b.jsonConfig)
	return b
}

// Final returns the finally built config
func (b Builder) Final() server.ConfigType {
	cfg := b.partial
	return cfg
}
