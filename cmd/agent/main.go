package main

import (
	"math/rand"
	"time"

	"internal/agent"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	go agent.RunCollectStats()
	agent.RunSendStats()
}
