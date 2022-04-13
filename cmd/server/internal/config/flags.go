package config

import (
	"flag"
	"net"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
)

type flags struct {
	configFile       common.StringFlag
	address          common.StringFlag
	storeInterval    common.TimeFlag
	storeFile        common.StringFlag
	restore          common.BoolFlag
	key              common.StringFlag
	cryptoKey        common.StringFlag
	databaseDSN      common.StringFlag
	trustedSubnetStr common.StringFlag
}

// ProcessFlags sets command-line flags to use
func (b *Builder) ProcessFlags() *Builder {
	b.flags.configFile.Option = "c"
	b.flags.configFile.Value = flag.String(b.flags.configFile.Option, "", "server config file")

	b.flags.address.Option = "a"
	b.flags.address.Value = flag.String(b.flags.address.Option, b.defaultConfig.Address, "bind address")

	b.flags.storeInterval.Option = "i"
	b.flags.storeInterval.Value = flag.Duration(b.flags.storeInterval.Option, b.defaultConfig.StoreInterval, "store interval")

	b.flags.storeFile.Option = "f"
	b.flags.storeFile.Value = flag.String(b.flags.storeFile.Option, b.defaultConfig.StoreFile, "store file")

	b.flags.restore.Option = "r"
	b.flags.restore.Value = flag.Bool(b.flags.restore.Option, b.defaultConfig.Restore, "restore")

	b.flags.key.Option = "k"
	b.flags.key.Value = flag.String(b.flags.key.Option, "", "key")

	b.flags.cryptoKey.Option = "crypto-key"
	b.flags.cryptoKey.Value = flag.String(b.flags.cryptoKey.Option, "", "crypto key")

	b.flags.databaseDSN.Option = "d"
	b.flags.databaseDSN.Value = flag.String(b.flags.databaseDSN.Option, "", "database dsn")

	b.flags.trustedSubnetStr.Option = "t"
	b.flags.trustedSubnetStr.Value = flag.String(b.flags.trustedSubnetStr.Option, "", "trusted subnet")

	flag.Parse()

	b.flags.configFile.Set = common.IsFlagPassed(b.flags.configFile.Option)
	b.flags.address.Set = common.IsFlagPassed(b.flags.address.Option)
	b.flags.storeInterval.Set = common.IsFlagPassed(b.flags.storeInterval.Option)
	b.flags.storeFile.Set = common.IsFlagPassed(b.flags.storeFile.Option)
	b.flags.restore.Set = common.IsFlagPassed(b.flags.restore.Option)
	b.flags.key.Set = common.IsFlagPassed(b.flags.key.Option)
	b.flags.cryptoKey.Set = common.IsFlagPassed(b.flags.cryptoKey.Option)
	b.flags.databaseDSN.Set = common.IsFlagPassed(b.flags.databaseDSN.Option)
	b.flags.trustedSubnetStr.Set = common.IsFlagPassed(b.flags.trustedSubnetStr.Option)

	return b
}

// MergeFlags merges values set with flags into the partial
func (b *Builder) MergeFlags() *Builder {
	if b.flags.address.Set {
		b.partial.Address = *b.flags.address.Value
	}
	if b.flags.storeInterval.Set {
		b.partial.StoreInterval = *b.flags.storeInterval.Value
	}
	if b.flags.storeFile.Set {
		b.partial.StoreFile = *b.flags.storeFile.Value
	}
	if b.flags.restore.Set {
		b.partial.Restore = *b.flags.restore.Value
	}
	if b.flags.key.Set {
		b.partial.Key = *b.flags.key.Value
	}
	if b.flags.cryptoKey.Set {
		b.partial.CryptoKey = *b.flags.cryptoKey.Value
	}
	if b.flags.databaseDSN.Set {
		b.partial.DatabaseDSN = *b.flags.databaseDSN.Value
	}
	if b.flags.trustedSubnetStr.Set {
		_, subnet, err := net.ParseCIDR(*b.flags.trustedSubnetStr.Value)
		if err != nil {
			b.err = err
			return b
		}
		b.partial.TrustedSubnet = subnet
	}

	return b
}
