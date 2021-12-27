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
	defaultAddress       = "localhost:8080"
	defaultStoreInterval = time.Second * 300
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultRestore       = true
)

type config struct {
	Address       *string        `env:"ADDRESS"`
	StoreInterval *time.Duration `env:"STORE_INTERVAL"`
	StoreFile     *string        `env:"STORE_FILE"`
	Restore       *bool          `env:"RESTORE"`
	Key           *string        `env:"KEY"`
	DatabaseDSN   *string        `env:"DATABASE_DSN"`
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
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
	storeIntervalFlag := flag.Duration("i", defaultStoreInterval, "store interval")
	fileFlag := flag.String("f", defaultStoreFile, "store file")
	restoreFlag := flag.Bool("r", defaultRestore, "restore")
	keyFlag := flag.String("k", "", "key")
	dbFlag := flag.String("d", "", "database dsn")

	flag.Parse()

	log.Printf("server is invoked with ENV %+v", cfg)
	log.Printf("server is invoked with flags address %v store interval %v store file %v restore %v database %v", *addressFlag, *storeIntervalFlag, *fileFlag, *restoreFlag, *dbFlag)

	server.Config.Address = *addressFlag
	if cfg.Address != nil {
		server.Config.Address = *cfg.Address
	}
	server.Config.StoreInterval = *storeIntervalFlag
	if cfg.StoreInterval != nil {
		server.Config.StoreInterval = *cfg.StoreInterval
	}

	// we need to distinguish between default string value and empty env var
	server.Config.StoreFile = defaultStoreFile
	if cfg.StoreFile != nil {
		server.Config.StoreFile = *cfg.StoreFile
	}
	if isFlagPassed("f") {
		server.Config.StoreFile = *fileFlag
	}

	server.Config.Restore = *restoreFlag
	if cfg.Restore != nil {
		server.Config.Restore = *cfg.Restore
	}

	server.Config.Key = *keyFlag
	if cfg.Key != nil {
		server.Config.Key = *cfg.Key
	}

	server.Config.DatabaseDSN = *dbFlag
	if cfg.DatabaseDSN != nil {
		server.Config.DatabaseDSN = *cfg.DatabaseDSN
	}
}

func main() {
	setServerArgs()

	prettyConfig, _ := json.Marshal(server.Config)
	log.Printf("server started with %v", string(prettyConfig))

	err := server.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
