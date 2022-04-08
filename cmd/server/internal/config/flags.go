package config

import (
	"flag"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
)

type stringFlag struct {
	value  *string
	option string
	set    bool
}

type boolFlag struct {
	value  *bool
	option string
	set    bool
}

type timeFlag struct {
	value  *time.Duration
	option string
	set    bool
}

type flags struct {
	configFile    stringFlag
	address       stringFlag
	storeInterval timeFlag
	storeFile     stringFlag
	restore       boolFlag
	key           stringFlag
	cryptoKey     stringFlag
	databaseDSN   stringFlag
}

// ProcessFlags sets command-line flags to use
func (b *Builder) ProcessFlags() *Builder {
	// FIXME: using `C` temporarily
	b.flags.configFile.option = "C"
	b.flags.configFile.value = flag.String(b.flags.configFile.option, "", "server config file")

	b.flags.address.option = "a"
	b.flags.address.value = flag.String(b.flags.address.option, b.defaultConfig.Address, "bind address")

	b.flags.storeInterval.option = "i"
	b.flags.storeInterval.value = flag.Duration(b.flags.storeInterval.option, b.defaultConfig.StoreInterval, "store interval")

	b.flags.storeFile.option = "f"
	b.flags.storeFile.value = flag.String(b.flags.storeFile.option, b.defaultConfig.StoreFile, "store file")

	b.flags.restore.option = "r"
	b.flags.restore.value = flag.Bool(b.flags.restore.option, b.defaultConfig.Restore, "restore")

	b.flags.key.option = "k"
	b.flags.key.value = flag.String(b.flags.key.option, "", "key")

	b.flags.cryptoKey.option = "c"
	b.flags.cryptoKey.value = flag.String(b.flags.cryptoKey.option, "", "crypto key")

	b.flags.databaseDSN.option = "d"
	b.flags.databaseDSN.value = flag.String(b.flags.databaseDSN.option, "", "database dsn")

	flag.Parse()

	b.flags.configFile.set = common.IsFlagPassed(b.flags.configFile.option)
	b.flags.address.set = common.IsFlagPassed(b.flags.address.option)
	b.flags.storeInterval.set = common.IsFlagPassed(b.flags.storeInterval.option)
	b.flags.storeFile.set = common.IsFlagPassed(b.flags.storeFile.option)
	b.flags.restore.set = common.IsFlagPassed(b.flags.restore.option)
	b.flags.key.set = common.IsFlagPassed(b.flags.key.option)
	b.flags.cryptoKey.set = common.IsFlagPassed(b.flags.cryptoKey.option)
	b.flags.databaseDSN.set = common.IsFlagPassed(b.flags.databaseDSN.option)

	return b
}

// MergeFlags merges values set with flags into the partial
func (b *Builder) MergeFlags() *Builder {
	if b.flags.address.set {
		b.partial.Address = *b.flags.address.value
	}
	if b.flags.storeInterval.set {
		b.partial.StoreInterval = *b.flags.storeInterval.value
	}
	if b.flags.storeFile.set {
		b.partial.StoreFile = *b.flags.storeFile.value
	}
	if b.flags.restore.set {
		b.partial.Restore = *b.flags.restore.value
	}
	if b.flags.key.set {
		b.partial.Key = *b.flags.key.value
	}
	if b.flags.cryptoKey.set {
		b.partial.CryptoKey = *b.flags.cryptoKey.value
	}
	if b.flags.databaseDSN.set {
		b.partial.DatabaseDSN = *b.flags.databaseDSN.value
	}
	return b
}
