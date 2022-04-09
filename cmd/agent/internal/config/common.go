// Package config provides config types for agent
package config

import (
	"log"
	"strings"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/agent"
)

// Builder is the base for building server config
type Builder struct {
	flags         flags
	jsonConfig    JSONConfig
	envVars       envVarConfig
	defaultConfig agent.ConfigType
	err           error
	partial       agent.ConfigType
}

// NewBuilder returns a pointer to the Builder struct filled with default values
func NewBuilder() *Builder {
	b := Builder{
		defaultConfig: agent.ConfigType{
			ServerAddr:     "http://localhost:8080",
			PollInterval:   time.Second * 2,
			ReportInterval: time.Second * 10,
		},
	}
	return &b
}

// MergeDefaults copies parameter from default values to the partial config
func (b *Builder) MergeDefaults() *Builder {
	b.partial.ServerAddr = b.defaultConfig.ServerAddr
	b.partial.PollInterval = b.defaultConfig.PollInterval
	b.partial.ReportInterval = b.defaultConfig.ReportInterval

	return b
}

// Err returns err
func (b Builder) Err() error {
	return b.err
}

// ReportEnvVars prints parsed env vars
func (b *Builder) ReportEnvVars() *Builder {
	log.Printf("agent is invoked with ENV %+v", b.envVars)
	return b
}

// ReportFlags prints passed flags
func (b *Builder) ReportFlags() *Builder {
	log.Printf("agent is invoked with flags address %v poll interval %v report interval %v key file %v",
		b.flags.address,
		b.flags.pollInterval,
		b.flags.reportInterval,
		b.flags.cryptoKey,
	)

	return b
}

// ReportJSONConfig prints values passed via config file
func (b *Builder) ReportJSONConfig() *Builder {
	log.Printf("config file values: %+v", b.jsonConfig)
	return b
}

// Final returns the finally built config
func (b Builder) Final() agent.ConfigType {
	// check if server address was given without "http://" or "https://"
	if !strings.HasPrefix(b.partial.ServerAddr, "http") {
		b.partial.ServerAddr = "http://" + b.partial.ServerAddr
	}
	cfg := b.partial
	return cfg
}
