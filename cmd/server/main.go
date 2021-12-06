package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
	"github.com/caarlos0/env/v6"
)

const (
	defaultAddress   = "localhost:8080"
	defaultInterval  = "300s"
	defaultStoreFile = "/tmp/devops-metrics-db.json"
	defaultRestore   = true
)

type config struct {
	Address       *string        `env:"ADDRESS"`
	StoreInterval *time.Duration `env:"STORE_INTERVAL"`
	StoreFile     *string        `env:"STORE_FILE"`
	Restore       *bool          `env:"RESTORE"`
}

func setServerArgs() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// distinguish between unset and set to ""
	val, ok := os.LookupEnv("STORE_FILE")
	if ok && val == "" {
		empty := ""
		cfg.StoreFile = &empty
	}

	addressFlag := flag.String("a", defaultAddress, "bind address")
	storeIntervalFlag := flag.String("i", defaultInterval, "store interval")
	fileFlag := flag.String("f", defaultStoreFile, "store file")
	restoreFlag := flag.Bool("r", defaultRestore, "restore")

	log.Printf("server is invoked with ENV address %v store interval %v store file %v restore %v", cfg.Address, cfg.StoreInterval, cfg.StoreFile, cfg.Restore)
	log.Printf("server is invoked with flags address %v store interval %v store file %v restore %v", *addressFlag, *storeIntervalFlag, *fileFlag, *restoreFlag)

	flag.Parse()

	if cfg.Address != nil {
		server.Config.Address = *cfg.Address
	} else {
		server.Config.Address = *addressFlag
	}

	if cfg.StoreInterval != nil {
		server.Config.StoreInterval = *cfg.StoreInterval
	} else {
		storeInterval, err := time.ParseDuration(*storeIntervalFlag)
		if err != nil {
			log.Fatal("cant parse duration ", *storeIntervalFlag)
		}
		server.Config.StoreInterval = storeInterval
	}

	if cfg.StoreFile != nil {
		server.Config.StoreFile = *cfg.StoreFile
	} else {
		server.Config.StoreFile = *fileFlag
	}

	if cfg.Restore != nil {
		server.Config.Restore = *cfg.Restore
	} else {
		server.Config.Restore = *restoreFlag
	}

}

func main() {
	setServerArgs()

	jsonConfig, _ := json.Marshal(server.Config)
	log.Print("server started with ", string(jsonConfig))

	server.StartServer()
}
