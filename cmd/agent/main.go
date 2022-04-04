package main

import (
	"log"

	"github.com/alexey-mavrin/go-musthave-devops/cmd/agent/internal/config"
	"github.com/alexey-mavrin/go-musthave-devops/internal/agent"
	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func setAgentArgs() error {
	builder := config.NewBuilder()

	agent.Config = builder.
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
