package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
)

type statData struct {
	mu          sync.Mutex
	memStats    runtime.MemStats
	PollCount   int64
	RandomValue int
}

var myStatData statData

const (
	defaultServer = "http://localhost:8080"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

type agentConfig struct {
	Server         string
	PollInterval   time.Duration
	ReportInterval time.Duration
	Key            string
	useJSON        bool
	useBatch       bool
}

// Config holds configuration parameters for the package
var Config agentConfig = agentConfig{
	Server:         defaultServer,
	PollInterval:   pollInterval,
	ReportInterval: reportInterval,
	useBatch:       true,
}

func collectStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	myStatData.mu.Lock()
	myStatData.PollCount++
	myStatData.RandomValue = rand.Int()
	myStatData.memStats = memStats
	myStatData.mu.Unlock()
}

// RunCollectStats collects data periodically
func RunCollectStats() {
	ticker := time.NewTicker(Config.PollInterval)
	for {
		<-ticker.C
		collectStats()
	}
}

func appendBatch(mm *[]common.Metrics, name string, data interface{}) {
	if mm == nil {
		log.Print("addBatch: trying to add to nil slice")
	}

	switch v := data.(type) {
	case int64:
		delta := v
		m := common.Metrics{
			ID:    name,
			MType: common.NameCounter,
			Delta: &delta,
		}
		m.StoreHash(Config.Key)
		*mm = append(*mm, m)

	case float64:
		value := v
		m := common.Metrics{
			ID:    name,
			MType: common.NameGauge,
			Value: &value,
		}
		m.StoreHash(Config.Key)
		*mm = append(*mm, m)
	}
}

func sendBatch(mm []common.Metrics) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(mm); err != nil {
		log.Fatal(err)
	}
	url := Config.Server + "/updates/"
	resp, err := http.Post(url, "application/json", &body)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Sending %s, http status %d", url, resp.StatusCode)
	}
}

func sendStatsBatch() {
	bm := make([]common.Metrics, 0, 100)
	myStatData.mu.Lock()
	appendBatch(&bm, "PollCount", myStatData.PollCount)
	appendBatch(&bm, "RandomValue", float64(myStatData.RandomValue))
	appendBatch(&bm, "Alloc", float64(myStatData.memStats.Alloc))
	appendBatch(&bm, "BuckHashSys", float64(myStatData.memStats.BuckHashSys))
	appendBatch(&bm, "Frees", float64(myStatData.memStats.Frees))
	appendBatch(&bm, "GCCPUFraction", float64(myStatData.memStats.GCCPUFraction))
	appendBatch(&bm, "GCSys", float64(myStatData.memStats.GCSys))
	appendBatch(&bm, "HeapAlloc", float64(myStatData.memStats.HeapAlloc))
	appendBatch(&bm, "HeapIdle", float64(myStatData.memStats.HeapIdle))
	appendBatch(&bm, "HeapInuse", float64(myStatData.memStats.HeapInuse))
	appendBatch(&bm, "HeapObjects", float64(myStatData.memStats.HeapObjects))
	appendBatch(&bm, "HeapReleased", float64(myStatData.memStats.HeapReleased))
	appendBatch(&bm, "HeapSys", float64(myStatData.memStats.HeapSys))
	appendBatch(&bm, "LastGC", float64(myStatData.memStats.LastGC))
	appendBatch(&bm, "Lookups", float64(myStatData.memStats.Lookups))
	appendBatch(&bm, "MCacheInuse", float64(myStatData.memStats.MCacheInuse))
	appendBatch(&bm, "MCacheSys", float64(myStatData.memStats.MCacheSys))
	appendBatch(&bm, "MSpanInuse", float64(myStatData.memStats.MSpanInuse))
	appendBatch(&bm, "MSpanSys", float64(myStatData.memStats.MSpanSys))
	appendBatch(&bm, "Mallocs", float64(myStatData.memStats.Mallocs))
	appendBatch(&bm, "NextGC", float64(myStatData.memStats.NextGC))
	appendBatch(&bm, "NumForcedGC", float64(myStatData.memStats.NumForcedGC))
	appendBatch(&bm, "NumGC", float64(myStatData.memStats.NumGC))
	appendBatch(&bm, "OtherSys", float64(myStatData.memStats.OtherSys))
	appendBatch(&bm, "PauseTotalNs", float64(myStatData.memStats.PauseTotalNs))
	appendBatch(&bm, "StackInuse", float64(myStatData.memStats.StackInuse))
	appendBatch(&bm, "StackSys", float64(myStatData.memStats.StackSys))
	appendBatch(&bm, "TotalAlloc", float64(myStatData.memStats.TotalAlloc))
	appendBatch(&bm, "Sys", float64(myStatData.memStats.Sys))
	myStatData.mu.Unlock()

	sendBatch(bm)
}

// RunSendStats periodically sends statistics to a collector
func RunSendStats() {
	ticker := time.NewTicker(Config.ReportInterval)
	for {
		<-ticker.C
		sendStatsBatch()
	}
}

// RunAgent is the function to start agent operation
func RunAgent() {
	rand.Seed(time.Now().UnixNano())
	go RunCollectStats()
	RunSendStats()

}
