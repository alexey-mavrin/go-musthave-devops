package config

import (
	"net"
	"os"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
	"github.com/caarlos0/env/v6"
)

type envVarConfig struct {
	Address          *string        `env:"ADDRESS"`
	StoreInterval    *time.Duration `env:"STORE_INTERVAL"`
	StoreFile        *string        `env:"STORE_FILE"`
	ConfigFile       *string        `env:"CONFIG"`
	Restore          *bool          `env:"RESTORE"`
	Key              *string        `env:"KEY"`
	CryptoKey        *string        `env:"CRYPTO_KEY"`
	DatabaseDSN      *string        `env:"DATABASE_DSN"`
	TrustedSubnetStr *string        `env:"TRUSTED_SUBNET"`
}

// ProcessEnvVars scans environment variables and store them in temporal struct
func (b *Builder) ProcessEnvVars() *Builder {
	err := env.Parse(&b.envVars)
	if err != nil {
		// one need to check if b.err is nill after processing
		b.err = err
		return b
	}

	// distinguish between unset and set to ""
	val, ok := os.LookupEnv("STORE_FILE")
	if ok && val == "" {
		empty := ""
		b.envVars.StoreFile = &empty
	}

	return b
}

// MergeEnvVars merges values from env variables into the partial config
func (b *Builder) MergeEnvVars() *Builder {
	common.CopyIfNotNil(&b.partial.Address, b.envVars.Address)
	common.CopyIfNotNil(&b.partial.StoreFile, b.envVars.StoreFile)
	common.CopyIfNotNil(&b.partial.Key, b.envVars.Key)
	common.CopyIfNotNil(&b.partial.CryptoKey, b.envVars.CryptoKey)
	common.CopyIfNotNil(&b.partial.DatabaseDSN, b.envVars.DatabaseDSN)

	if b.envVars.StoreInterval != nil {
		b.partial.StoreInterval = *b.envVars.StoreInterval
	}

	if b.envVars.Restore != nil {
		b.partial.Restore = *b.envVars.Restore
	}

	if b.envVars.TrustedSubnetStr != nil {
		_, subnet, err := net.ParseCIDR(*b.envVars.TrustedSubnetStr)
		if err != nil {
			b.err = err
			return b
		}
		b.partial.TrustedSubnet = subnet
	}

	return b
}
