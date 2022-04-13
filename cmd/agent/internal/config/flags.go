package config

import (
	"flag"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
)

type flags struct {
	configFile     common.StringFlag
	address        common.StringFlag
	pollInterval   common.TimeFlag
	reportInterval common.TimeFlag
	key            common.StringFlag
	cryptoKey      common.StringFlag
}

// ProcessFlags sets command-line flags to use
func (b *Builder) ProcessFlags() *Builder {
	b.flags.configFile.Option = "c"
	b.flags.configFile.Value = flag.String(b.flags.configFile.Option, "", "server config file")

	b.flags.address.Option = "a"
	b.flags.address.Value = flag.String(b.flags.address.Option, b.defaultConfig.ServerAddr, "server address")

	b.flags.pollInterval.Option = "p"
	b.flags.pollInterval.Value = flag.Duration(b.flags.pollInterval.Option, b.defaultConfig.PollInterval, "poll interval")

	b.flags.reportInterval.Option = "r"
	b.flags.reportInterval.Value = flag.Duration(b.flags.reportInterval.Option, b.defaultConfig.ReportInterval, "report interval")

	b.flags.key.Option = "k"
	b.flags.key.Value = flag.String(b.flags.key.Option, "", "key")

	b.flags.cryptoKey.Option = "crypto-key"
	b.flags.cryptoKey.Value = flag.String(b.flags.cryptoKey.Option, "", "crypto key")

	flag.Parse()

	b.flags.configFile.Set = common.IsFlagPassed(b.flags.configFile.Option)
	b.flags.address.Set = common.IsFlagPassed(b.flags.address.Option)
	b.flags.pollInterval.Set = common.IsFlagPassed(b.flags.pollInterval.Option)
	b.flags.reportInterval.Set = common.IsFlagPassed(b.flags.reportInterval.Option)
	b.flags.key.Set = common.IsFlagPassed(b.flags.key.Option)
	b.flags.cryptoKey.Set = common.IsFlagPassed(b.flags.cryptoKey.Option)

	return b
}

// MergeFlags merges values set with flags into the partial
func (b *Builder) MergeFlags() *Builder {
	if b.flags.address.Set {
		b.partial.ServerAddr = *b.flags.address.Value
	}
	if b.flags.pollInterval.Set {
		b.partial.PollInterval = *b.flags.pollInterval.Value
	}
	if b.flags.reportInterval.Set {
		b.partial.ReportInterval = *b.flags.reportInterval.Value
	}
	if b.flags.key.Set {
		b.partial.Key = *b.flags.key.Value
	}
	if b.flags.cryptoKey.Set {
		b.partial.CryptoKey = *b.flags.cryptoKey.Value
	}
	return b
}
