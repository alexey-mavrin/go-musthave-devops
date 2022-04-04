package config

import (
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
	"github.com/caarlos0/env/v6"
)

type envVarConfig struct {
	ServerAddr     *string        `env:"ADDRESS"`
	PollInterval   *time.Duration `env:"POLL_INTERVAL"`
	ReportInterval *time.Duration `env:"REPORT_INTERVAL"`
	ConfigFile     *string        `env:"CONFIG"`
	Key            *string        `env:"KEY"`
	CryptoKey      *string        `env:"CRYPTO_KEY"`
}

// ProcessEnvVars scans environment variables and store them in temporal struct
func (b *Builder) ProcessEnvVars() *Builder {
	err := env.Parse(&b.envVars)
	if err != nil {
		// one need to check if b.err is nill after processing
		b.err = err
		return b
	}

	return b
}

// MergeEnvVars merges values from env variables into the partial config
func (b *Builder) MergeEnvVars() *Builder {
	common.CopyIfNotNil(&b.partial.ServerAddr, b.envVars.ServerAddr)
	common.CopyIfNotNil(&b.partial.Key, b.envVars.Key)
	common.CopyIfNotNil(&b.partial.CryptoKey, b.envVars.CryptoKey)

	if b.envVars.PollInterval != nil {
		b.partial.PollInterval = *b.envVars.PollInterval
	}

	if b.envVars.ReportInterval != nil {
		b.partial.ReportInterval = *b.envVars.ReportInterval
	}

	return b
}
