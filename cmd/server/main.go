package main

import (
	"encoding/json"
	"log"

	"github.com/alexey-mavrin/go-musthave-devops/cmd/server/internal/config"
	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
	"github.com/alexey-mavrin/go-musthave-devops/internal/server"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func setServerArgs() error {
	builder := config.NewBuilder()

	server.Config = builder.
		MergeDefaults().
		ProcessFlags().
		ProcessEnvVars().
		ReadJSONConfig().
		MergeJSONConfig().
		MergeFlags().
		MergeEnvVars().
		ReportJSONConfig().
		ReportFlags().
		ReportEnvVars().
		Final()

	return builder.Err()
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
