// Package agent is the package with the agent code.
package agent

import (
	"bytes"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
	"github.com/alexey-mavrin/go-musthave-devops/internal/crypt"
	"github.com/alexey-mavrin/go-musthave-devops/internal/iproute"
)

type statData struct {
	CPUutilLastTime time.Time
	CPUtime         []float64
	CPUutilization  []float64
	mu              sync.Mutex
	memStats        runtime.MemStats
	PollCount       int64
	RandomValue     int
	TotalMemory     float64
	FreeMemory      float64
}

func init() {
	cpuStat, err := cpu.Info()
	if err != nil {
		log.Println(err)
		return
	}
	numCPU := len(cpuStat)
	myStatData.CPUtime = make([]float64, numCPU)
	myStatData.CPUutilization = make([]float64, numCPU)
}

var myStatData statData

const (
	defaultServer = "http://localhost:8080"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

// ConfigType contains config options for the agent
type ConfigType struct {
	ServerAddr     string
	Key            string
	CryptoKey      string
	PollInterval   time.Duration
	ReportInterval time.Duration
	useJSON        bool
	useBatch       bool
}

var publicServerKey *rsa.PublicKey

// Config holds configuration parameters for the package
var Config ConfigType = ConfigType{
	ServerAddr:     defaultServer,
	PollInterval:   pollInterval,
	ReportInterval: reportInterval,
	useBatch:       true,
}

// ReadServerKey reads server public key if provided
func ReadServerKey() error {
	if Config.CryptoKey == "" {
		return nil
	}
	var err error
	publicServerKey, err = crypt.ReadPublicKey(Config.CryptoKey)
	log.Printf("key %s read: %v", Config.CryptoKey, publicServerKey)
	return err
}

func collectPSStats() {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
	}
	c, err := cpu.Times(true)
	timeNow := time.Now()
	if err != nil {
		log.Println(err)
	}
	myStatData.mu.Lock()
	timeDiff := timeNow.Sub(myStatData.CPUutilLastTime)
	myStatData.CPUutilLastTime = timeNow
	myStatData.TotalMemory = float64(m.Total)
	myStatData.FreeMemory = float64(m.Free)
	for n := range c {
		newCPUTime := c[n].User + c[n].System
		cpuUtil := (newCPUTime - myStatData.CPUtime[n]) * 1000 / float64(timeDiff.Milliseconds())
		myStatData.CPUutilization[n] = cpuUtil
		myStatData.CPUtime[n] = newCPUTime
	}
	myStatData.mu.Unlock()

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

// RunCollectPSStats collects data from psutils periodically
func RunCollectPSStats() {
	ticker := time.NewTicker(Config.PollInterval)
	for {
		<-ticker.C
		collectPSStats()
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
	case []float64:
		for item := range v {
			value := v[item]
			m := common.Metrics{
				ID:    fmt.Sprintf("%s%d", name, item),
				MType: common.NameGauge,
				Value: &value,
			}
			m.StoreHash(Config.Key)
			initial = append(initial, m)
		}
		return initial
	}
	return initial
}

func sendBatch(mm []common.Metrics) error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(mm); err != nil {
		return err
	}
	url := Config.ServerAddr + "/updates/"

	if publicServerKey != nil {
		encryptedBytes, err := crypt.EncryptOAEP(
			sha256.New(),
			crand.Reader,
			publicServerKey,
			body.Bytes(),
			nil)
		if err != nil {
			return err
		}
		body.Reset()
		body.Write(encryptedBytes)
	}

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	ip, err := iproute.GetSrcIPURL(url)
	if err != nil {
		log.Printf("error getting source IP address")
	}
	req.Header.Set("X-Real-IP", ip)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Sending %s, http status %d", url, resp.StatusCode)
	}

	return nil
}

func sendStatsBatch() error {
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
	bm = appendBatch(bm, "TotalMemory", myStatData.TotalMemory)
	bm = appendBatch(bm, "FreeMemory", myStatData.FreeMemory)
	bm = appendBatch(bm, "CPUutilization", myStatData.CPUutilization)
	myStatData.mu.Unlock()

	return sendBatch(bm)
}

// RunSendStats periodically sends statistics to a collector
func RunSendStats() {
	ticker := time.NewTicker(Config.ReportInterval)
	for {
		<-ticker.C
		err := sendStatsBatch()
		if err != nil {
			log.Printf("error sending update: %v", err)
		}
	}
}

// RunAgent is the function to start agent operation
func RunAgent() {
	rand.Seed(time.Now().UnixNano())
	go RunCollectStats()
	go RunCollectPSStats()
	RunSendStats()

}
