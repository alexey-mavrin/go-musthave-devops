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
	configFile     stringFlag
	address        stringFlag
	pollInterval   timeFlag
	reportInterval timeFlag
	key            stringFlag
	cryptoKey      stringFlag
}

// ProcessFlags sets command-line flags to use
func (b *Builder) ProcessFlags() *Builder {
	// FIXME: using `C` temporarily
	b.flags.configFile.option = "C"
	b.flags.configFile.value = flag.String(b.flags.configFile.option, "", "server config file")

	b.flags.address.option = "a"
	b.flags.address.value = flag.String(b.flags.address.option, b.defaultConfig.ServerAddr, "server address")

	b.flags.pollInterval.option = "p"
	b.flags.pollInterval.value = flag.Duration(b.flags.pollInterval.option, b.defaultConfig.PollInterval, "poll interval")

	b.flags.reportInterval.option = "r"
	b.flags.reportInterval.value = flag.Duration(b.flags.reportInterval.option, b.defaultConfig.ReportInterval, "report interval")

	b.flags.key.option = "k"
	b.flags.key.value = flag.String(b.flags.key.option, "", "key")

	b.flags.cryptoKey.option = "c"
	b.flags.cryptoKey.value = flag.String(b.flags.cryptoKey.option, "", "crypto key")

	flag.Parse()

	b.flags.configFile.set = common.IsFlagPassed(b.flags.configFile.option)
	b.flags.address.set = common.IsFlagPassed(b.flags.address.option)
	b.flags.pollInterval.set = common.IsFlagPassed(b.flags.pollInterval.option)
	b.flags.reportInterval.set = common.IsFlagPassed(b.flags.reportInterval.option)
	b.flags.key.set = common.IsFlagPassed(b.flags.key.option)
	b.flags.cryptoKey.set = common.IsFlagPassed(b.flags.cryptoKey.option)

	return b
}

// MergeFlags merges values set with flags into the partial
func (b *Builder) MergeFlags() *Builder {
	if b.flags.address.set {
		b.partial.ServerAddr = *b.flags.address.value
	}
	if b.flags.pollInterval.set {
		b.partial.PollInterval = *b.flags.pollInterval.value
	}
	if b.flags.reportInterval.set {
		b.partial.ReportInterval = *b.flags.reportInterval.value
	}
	if b.flags.key.set {
		b.partial.Key = *b.flags.key.value
	}
	if b.flags.cryptoKey.set {
		b.partial.CryptoKey = *b.flags.cryptoKey.value
	}
	return b
}
