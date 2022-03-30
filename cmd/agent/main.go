package main

import (
	"flag"
	"log"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/agent"
	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
	"github.com/caarlos0/env/v6"
)

type config struct {
	Address        *string        `env:"ADDRESS"`
	PollInterval   *time.Duration `env:"POLL_INTERVAL"`
	ReportInterval *time.Duration `env:"REPORT_INTERVAL"`
	Key            *string        `env:"KEY"`
	CryptoKey      *string        `env:"CRYPTO_KEY"`
}

const (
	defaultAddress        = "localhost:8080"
	defaultScheme         = "http"
	defaultPollInterval   = time.Second * 2
	defaultReportInterval = time.Second * 10
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func setAgentArgs() error {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}

	addressFlag := flag.String("a", defaultAddress, "server address")
	pollIntervalFlag := flag.Duration("p", defaultPollInterval, "poll interval")
	reportIntervalFlag := flag.Duration("r", defaultReportInterval, "report interval")
	keyFlag := flag.String("k", "", "key")
	cryptoKeyFlag := flag.String("c", "", "crypto key")

	flag.Parse()

	agent.Config.ServerAddr = defaultScheme + "://" + *addressFlag
	if cfg.Address != nil {
		agent.Config.ServerAddr = defaultScheme + "://" + *cfg.Address
	}

	agent.Config.PollInterval = *pollIntervalFlag
	if cfg.PollInterval != nil {
		agent.Config.PollInterval = *cfg.PollInterval
	}

	agent.Config.ReportInterval = *reportIntervalFlag
	if cfg.ReportInterval != nil {
		agent.Config.ReportInterval = *cfg.ReportInterval
	}

	agent.Config.Key = *keyFlag
	if cfg.Key != nil {
		agent.Config.Key = *cfg.Key
	}

	agent.Config.CryptoKey = *cryptoKeyFlag
	if cfg.CryptoKey != nil {
		agent.Config.CryptoKey = *cfg.CryptoKey
	}

	return nil
}

func main() {
	if err := setAgentArgs(); err != nil {
		log.Fatal(err)
	}
	common.PrintBuildInfo(buildVersion, buildDate, buildCommit)
	if err := agent.ReadServerKey(); err != nil {
		log.Fatal(err)
	}

	// we don't need \n as log.Printf do is automatically
	log.Printf("agent started with %+v", agent.Config)

	agent.RunAgent()
}
