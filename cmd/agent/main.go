package main

import (
	"flag"
	"log"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/agent"
	"github.com/caarlos0/env/v6"
)

type config struct {
	Address        *string        `env:"ADDRESS"`
	PollInterval   *time.Duration `env:"POLL_INTERVAL"`
	ReportInterval *time.Duration `env:"REPORT_INTERVAL"`
}

const (
	defaultAddress        = "localhost:8080"
	defaultScheme         = "http"
	defaultPollInterval   = time.Second * 2
	defaultReportInterval = time.Second * 10
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

	flag.Parse()

	agent.Config.Server = defaultScheme + "://" + *addressFlag
	if cfg.Address != nil {
		agent.Config.Server = defaultScheme + "://" + *cfg.Address
	}

	agent.Config.PollInterval = *pollIntervalFlag
	if cfg.PollInterval != nil {
		agent.Config.PollInterval = *cfg.PollInterval
	}

	agent.Config.ReportInterval = *reportIntervalFlag
	if cfg.ReportInterval != nil {
		agent.Config.ReportInterval = *cfg.ReportInterval
	}

	return nil
}

func main() {
	if err := setAgentArgs(); err != nil {
		log.Fatal(err)
	}

	// we don't need \n as log.Printf do is automatically
	log.Printf("agent started with %+v", agent.Config)

	agent.RunAgent()
}
