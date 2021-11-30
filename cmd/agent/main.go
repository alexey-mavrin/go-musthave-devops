package main

import (
	"github.com/alexey-mavrin/go-musthave-devops/internal/agent"
)

func main() {
	// agent.Config.Server = "http://localhost:8080"
	// agent.Config.PollInterval = 2 * time.Second
	// agent.Config.ReportInterval = 10 * time.Second
	agent.RunAgent()
}
