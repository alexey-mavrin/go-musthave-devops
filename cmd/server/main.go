package main

import (
	"log"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
	"github.com/caarlos0/env/v6"
)

type config struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	server.Config.Address = cfg.Address
	server.Config.StoreInterval = cfg.StoreInterval
	server.Config.StoreFile = cfg.StoreFile
	server.Config.Restore = cfg.Restore

	server.StartServer()
}
