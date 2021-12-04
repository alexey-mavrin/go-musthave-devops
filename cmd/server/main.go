package main

import (
	"log"
	"net/http"

	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
	"github.com/caarlos0/env/v6"
)

type config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := server.Router()
	http.ListenAndServe(cfg.Address, r)
}
