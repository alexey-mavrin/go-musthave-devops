package main

import (
	"math/rand"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/agent"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	go agent.RunCollectStats()
	agent.RunSendStats()
}
