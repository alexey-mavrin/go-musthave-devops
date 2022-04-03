package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
	"github.com/caarlos0/env/v6"
)

const (
	defaultAddress       = "localhost:8080"
	defaultStoreInterval = time.Second * 300
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultRestore       = true
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

type config struct {
	Address       *string        `env:"ADDRESS"`
	StoreInterval *time.Duration `env:"STORE_INTERVAL"`
	StoreFile     *string        `env:"STORE_FILE"`
	Restore       *bool          `env:"RESTORE"`
	Key           *string        `env:"KEY"`
	CryptoKey     *string        `env:"CRYPTO_KEY"`
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

func setServerArgs() error {
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

	// FIXME: using `C` temporarily
	configFileFlag := flag.String("C", "", "server config file")
	addressFlag := flag.String("a", defaultAddress, "bind address")
	storeIntervalFlag := flag.Duration("i", defaultStoreInterval, "store interval")
	fileFlag := flag.String("f", defaultStoreFile, "store file")
	restoreFlag := flag.Bool("r", defaultRestore, "restore")
	keyFlag := flag.String("k", "", "key")
	cryptoKeyFlag := flag.String("c", "", "crypto key")
	dbFlag := flag.String("d", "", "database dsn")

	flag.Parse()

	log.Printf("server is invoked with ENV %+v", cfg)
	log.Printf("server is invoked with flags address %v store interval %v store file %v restore %v database %v", *addressFlag, *storeIntervalFlag, *fileFlag, *restoreFlag, *dbFlag)

	var fileConfig server.JSONConfig
	if isFlagPassed("C") {
		fileConfig, err = server.ReadJSONConfig(*configFileFlag)
		if err != nil {
			return err
		}
	}

	if theAddress := common.FirstSet(
		fileConfig.Address,
		addressFlag,
		cfg.Address,
	); theAddress != nil {
		server.Config.Address = *theAddress
	}

	server.Config.StoreInterval = defaultStoreInterval
	if fileConfig.StoreIntervalStr != nil {
		server.Config.StoreInterval, err = time.ParseDuration(*fileConfig.StoreIntervalStr)
		if err != nil {
			return err
		}
	}
	if isFlagPassed("i") {
		server.Config.StoreInterval = *storeIntervalFlag
	}
	if cfg.StoreInterval != nil {
		server.Config.StoreInterval = *cfg.StoreInterval
	}

	// we need to distinguish between default string value and empty env var
	server.Config.StoreFile = defaultStoreFile
	if fileConfig.StoreFile != nil {
		server.Config.StoreFile = *fileConfig.StoreFile
	}
	if cfg.StoreFile != nil {
		server.Config.StoreFile = *cfg.StoreFile
	}
	if isFlagPassed("f") {
		server.Config.StoreFile = *fileFlag
	}

	server.Config.Restore = true
	if fileConfig.Restore != nil {
		server.Config.Restore = *fileConfig.Restore
	}
	if isFlagPassed("r") {
		server.Config.Restore = *restoreFlag
	}
	if cfg.Restore != nil {
		server.Config.Restore = *cfg.Restore
	}

	if theKey := common.FirstSet(
		fileConfig.Key,
		keyFlag,
		cfg.Key,
	); theKey != nil {
		server.Config.Key = *theKey
	}

	if theCryptoKey := common.FirstSet(
		fileConfig.CryptoKey,
		cryptoKeyFlag,
		cfg.CryptoKey,
	); theCryptoKey != nil {
		server.Config.CryptoKey = *theCryptoKey
	}

	if theDatabaseDSN := common.FirstSet(
		fileConfig.DatabaseDSN,
		dbFlag,
		cfg.DatabaseDSN,
	); theDatabaseDSN != nil {
		server.Config.DatabaseDSN = *theDatabaseDSN
	}

	return nil
}

func main() {
	if err := setServerArgs(); err != nil {
		log.Fatal(err)
	}

	common.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	prettyConfig, err := json.Marshal(server.Config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server started with %v", string(prettyConfig))
	if err = server.ReadServerKey(); err != nil {
		log.Fatal(err)
	}

	err = server.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
