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
	ServerAddr     string
	PollInterval   time.Duration
	ReportInterval time.Duration
	Key            string
	useJSON        bool
	useBatch       bool
}

// Config holds configuration parameters for the package
var Config agentConfig = agentConfig{
	ServerAddr:     defaultServer,
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

func appendBatch(initial []common.Metrics, name string, data interface{}) []common.Metrics {
	if initial == nil {
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
		return append(initial, m)

	case float64:
		value := v
		m := common.Metrics{
			ID:    name,
			MType: common.NameGauge,
			Value: &value,
		}
		m.StoreHash(Config.Key)
		return append(initial, m)
	}
	return initial
}

func sendBatch(mm []common.Metrics) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(mm); err != nil {
		log.Print(err)
		return
	}
	url := Config.ServerAddr + "/updates/"
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
	bm = appendBatch(bm, "PollCount", myStatData.PollCount)
	bm = appendBatch(bm, "RandomValue", float64(myStatData.RandomValue))
	bm = appendBatch(bm, "Alloc", float64(myStatData.memStats.Alloc))
	bm = appendBatch(bm, "BuckHashSys", float64(myStatData.memStats.BuckHashSys))
	bm = appendBatch(bm, "Frees", float64(myStatData.memStats.Frees))
	bm = appendBatch(bm, "GCCPUFraction", float64(myStatData.memStats.GCCPUFraction))
	bm = appendBatch(bm, "GCSys", float64(myStatData.memStats.GCSys))
	bm = appendBatch(bm, "HeapAlloc", float64(myStatData.memStats.HeapAlloc))
	bm = appendBatch(bm, "HeapIdle", float64(myStatData.memStats.HeapIdle))
	bm = appendBatch(bm, "HeapInuse", float64(myStatData.memStats.HeapInuse))
	bm = appendBatch(bm, "HeapObjects", float64(myStatData.memStats.HeapObjects))
	bm = appendBatch(bm, "HeapReleased", float64(myStatData.memStats.HeapReleased))
	bm = appendBatch(bm, "HeapSys", float64(myStatData.memStats.HeapSys))
	bm = appendBatch(bm, "LastGC", float64(myStatData.memStats.LastGC))
	bm = appendBatch(bm, "Lookups", float64(myStatData.memStats.Lookups))
	bm = appendBatch(bm, "MCacheInuse", float64(myStatData.memStats.MCacheInuse))
	bm = appendBatch(bm, "MCacheSys", float64(myStatData.memStats.MCacheSys))
	bm = appendBatch(bm, "MSpanInuse", float64(myStatData.memStats.MSpanInuse))
	bm = appendBatch(bm, "MSpanSys", float64(myStatData.memStats.MSpanSys))
	bm = appendBatch(bm, "Mallocs", float64(myStatData.memStats.Mallocs))
	bm = appendBatch(bm, "NextGC", float64(myStatData.memStats.NextGC))
	bm = appendBatch(bm, "NumForcedGC", float64(myStatData.memStats.NumForcedGC))
	bm = appendBatch(bm, "NumGC", float64(myStatData.memStats.NumGC))
	bm = appendBatch(bm, "OtherSys", float64(myStatData.memStats.OtherSys))
	bm = appendBatch(bm, "PauseTotalNs", float64(myStatData.memStats.PauseTotalNs))
	bm = appendBatch(bm, "StackInuse", float64(myStatData.memStats.StackInuse))
	bm = appendBatch(bm, "StackSys", float64(myStatData.memStats.StackSys))
	bm = appendBatch(bm, "TotalAlloc", float64(myStatData.memStats.TotalAlloc))
	bm = appendBatch(bm, "Sys", float64(myStatData.memStats.Sys))
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
