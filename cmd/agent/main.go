package main

import (
	"encoding/json"
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

func setAgentArgs() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	addressFlag := flag.String("a", defaultAddress, "server address")
	pollIntervalFlag := flag.Duration("p", defaultPollInterval, "poll interval")
	reportIntervalFlag := flag.Duration("r", defaultReportInterval, "report interval")

	flag.Parse()

	if cfg.Address != nil {
		agent.Config.Server = defaultScheme + "://" + *cfg.Address
	} else {
		agent.Config.Server = defaultScheme + "://" + *addressFlag
	}

	if cfg.PollInterval != nil {
		agent.Config.PollInterval = *cfg.PollInterval
	} else {
		agent.Config.PollInterval = *pollIntervalFlag
	}
	if cfg.ReportInterval != nil {
		agent.Config.ReportInterval = *cfg.ReportInterval
	} else {
		agent.Config.ReportInterval = *reportIntervalFlag
	}
}

func main() {
	setAgentArgs()

	jsonConfig, _ := json.Marshal(agent.Config)
	log.Print("agent started with ", string(jsonConfig))

	agent.RunAgent()
}
