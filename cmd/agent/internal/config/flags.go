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
	useGRPC        common.BoolFlag
	gRPCServer     common.StringFlag
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

	b.flags.useGRPC.Option = "use-grpc"
	b.flags.useGRPC.Value = flag.Bool(b.flags.useGRPC.Option, false, "use gRPC")

	b.flags.gRPCServer.Option = "grpc-server"
	b.flags.gRPCServer.Value = flag.String(b.flags.gRPCServer.Option, "", "gRPC server")

	flag.Parse()

	b.flags.configFile.Set = common.IsFlagPassed(b.flags.configFile.Option)
	b.flags.address.Set = common.IsFlagPassed(b.flags.address.Option)
	b.flags.pollInterval.Set = common.IsFlagPassed(b.flags.pollInterval.Option)
	b.flags.reportInterval.Set = common.IsFlagPassed(b.flags.reportInterval.Option)
	b.flags.key.Set = common.IsFlagPassed(b.flags.key.Option)
	b.flags.cryptoKey.Set = common.IsFlagPassed(b.flags.cryptoKey.Option)
	b.flags.useGRPC.Set = common.IsFlagPassed(b.flags.useGRPC.Option)
	b.flags.gRPCServer.Set = common.IsFlagPassed(b.flags.gRPCServer.Option)

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
	if b.flags.useGRPC.Set {
		b.partial.UseGRPC = *b.flags.useGRPC.Value
	}
	if b.flags.gRPCServer.Set {
		b.partial.GRPCServer = *b.flags.gRPCServer.Value
	}
	return b
}
