package main

import (
	"log"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/agent"
	"github.com/caarlos0/env/v6"
)

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	Scheme         string        `env:"SCHEME" envDefault:"http"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	agent.Config.Server = cfg.Scheme + "://" + cfg.Address
	agent.Config.PollInterval = cfg.PollInterval
	agent.Config.ReportInterval = cfg.ReportInterval

	agent.RunAgent()
}
