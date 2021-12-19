package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
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
}

type stat struct {
	url  string
	data []byte
}

// Config holds configuration parameters for the package
var Config agentConfig = agentConfig{
	Server:         defaultServer,
	PollInterval:   pollInterval,
	ReportInterval: reportInterval,
	useJSON:        true,
}

func sendStat(s stat) {
	body := bytes.NewBuffer(s.data)
	contentType := "text/plain"
	if Config.useJSON {
		contentType = "application/json"
	}
	resp, err := http.Post(Config.Server+s.url, contentType, body)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Sending %s, http status %d", s.url, resp.StatusCode)
	}
}

func makeStatGauge(name string, value float64, useJSON bool) (stat, error) {
	var (
		s   stat
		err error
	)
	if useJSON {
		s.url = "/update/"
		s.data, err = makeStatJSONGauge(name, value)
		if err != nil {
			return s, err
		}
	} else {
		s.url = makeStatStringGauge(name, value)
	}
	return s, nil
}

func makeStatCounter(name string, value int64, useJSON bool) (stat, error) {
	var (
		s   stat
		err error
	)
	if useJSON {
		s.url = "/update/"
		s.data, err = makeStatJSONCounter(name, value)
		if err != nil {
			return s, err
		}
	} else {
		s.url = makeStatStringCounter(name, value)
	}
	return s, nil
}

func makeStatStringGauge(name string, value float64) string {
	return fmt.Sprintf("/update/gauge/%s/%G", name, value)
}

func makeStatStringCounter(name string, value int64) string {
	return fmt.Sprintf("/update/counter/%s/%d", name, value)
}

func makeStatJSONGauge(name string, value float64) ([]byte, error) {
	var m = common.Metrics{
		ID:    name,
		MType: common.NameGauge,
		Value: &value,
	}
	m.StoreHash(Config.Key)
	return json.Marshal(m)
}

func makeStatJSONCounter(name string, delta int64) ([]byte, error) {
	var m = common.Metrics{
		ID:    name,
		MType: common.NameCounter,
		Delta: &delta,
	}
	m.StoreHash(Config.Key)
	return json.Marshal(m)
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

func sendStats() {
	myStatData.mu.Lock()
	PollCount, _ := makeStatCounter("PollCount", myStatData.PollCount, Config.useJSON)
	myStatData.PollCount = 0

	RandomValueGauge, _ := makeStatGauge("RandomValue", float64(myStatData.RandomValue), Config.useJSON)
	Alloc, _ := makeStatGauge("Alloc", float64(myStatData.memStats.Alloc), Config.useJSON)
	BuckHashSys, _ := makeStatGauge("BuckHashSys", float64(myStatData.memStats.BuckHashSys), Config.useJSON)
	Frees, _ := makeStatGauge("Frees", float64(myStatData.memStats.Frees), Config.useJSON)
	GCCPUFraction, _ := makeStatGauge("GCCPUFraction", float64(myStatData.memStats.GCCPUFraction), Config.useJSON)
	GCSys, _ := makeStatGauge("GCSys", float64(myStatData.memStats.GCSys), Config.useJSON)
	HeapAlloc, _ := makeStatGauge("HeapAlloc", float64(myStatData.memStats.HeapAlloc), Config.useJSON)
	HeapIdle, _ := makeStatGauge("HeapIdle", float64(myStatData.memStats.HeapIdle), Config.useJSON)
	HeapInuse, _ := makeStatGauge("HeapInuse", float64(myStatData.memStats.HeapInuse), Config.useJSON)
	HeapObjects, _ := makeStatGauge("HeapObjects", float64(myStatData.memStats.HeapObjects), Config.useJSON)
	HeapReleased, _ := makeStatGauge("HeapReleased", float64(myStatData.memStats.HeapReleased), Config.useJSON)
	HeapSys, _ := makeStatGauge("HeapSys", float64(myStatData.memStats.HeapSys), Config.useJSON)
	LastGC, _ := makeStatGauge("LastGC", float64(myStatData.memStats.LastGC), Config.useJSON)
	Lookups, _ := makeStatGauge("Lookups", float64(myStatData.memStats.Lookups), Config.useJSON)
	MCacheInuse, _ := makeStatGauge("MCacheInuse", float64(myStatData.memStats.MCacheInuse), Config.useJSON)
	MCacheSys, _ := makeStatGauge("MCacheSys", float64(myStatData.memStats.MCacheSys), Config.useJSON)
	MSpanInuse, _ := makeStatGauge("MSpanInuse", float64(myStatData.memStats.MSpanInuse), Config.useJSON)
	MSpanSys, _ := makeStatGauge("MSpanSys", float64(myStatData.memStats.MSpanSys), Config.useJSON)
	Mallocs, _ := makeStatGauge("Mallocs", float64(myStatData.memStats.Mallocs), Config.useJSON)
	NextGC, _ := makeStatGauge("NextGC", float64(myStatData.memStats.NextGC), Config.useJSON)
	NumForcedGC, _ := makeStatGauge("NumForcedGC", float64(myStatData.memStats.NumForcedGC), Config.useJSON)
	NumGC, _ := makeStatGauge("NumGC", float64(myStatData.memStats.NumGC), Config.useJSON)
	OtherSys, _ := makeStatGauge("OtherSys", float64(myStatData.memStats.OtherSys), Config.useJSON)
	PauseTotalNs, _ := makeStatGauge("PauseTotalNs", float64(myStatData.memStats.PauseTotalNs), Config.useJSON)
	StackInuse, _ := makeStatGauge("StackInuse", float64(myStatData.memStats.StackInuse), Config.useJSON)
	StackSys, _ := makeStatGauge("StackSys", float64(myStatData.memStats.StackSys), Config.useJSON)
	TotalAlloc, _ := makeStatGauge("TotalAlloc", float64(myStatData.memStats.TotalAlloc), Config.useJSON)
	Sys, _ := makeStatGauge("Sys", float64(myStatData.memStats.Sys), Config.useJSON)
	myStatData.mu.Unlock()

	sendStat(PollCount)
	sendStat(RandomValueGauge)
	sendStat(Alloc)
	sendStat(BuckHashSys)
	sendStat(Frees)
	sendStat(GCCPUFraction)
	sendStat(GCSys)
	sendStat(HeapAlloc)
	sendStat(HeapIdle)
	sendStat(HeapInuse)
	sendStat(HeapObjects)
	sendStat(HeapReleased)
	sendStat(HeapSys)
	sendStat(LastGC)
	sendStat(Lookups)
	sendStat(MCacheInuse)
	sendStat(MCacheSys)
	sendStat(MSpanInuse)
	sendStat(MSpanSys)
	sendStat(Mallocs)
	sendStat(NextGC)
	sendStat(NumForcedGC)
	sendStat(NumGC)
	sendStat(OtherSys)
	sendStat(PauseTotalNs)
	sendStat(StackInuse)
	sendStat(StackSys)
	sendStat(TotalAlloc)
	sendStat(Sys)
}

// RunSendStats periodically sends statistics to a collector
func RunSendStats() {
	ticker := time.NewTicker(Config.ReportInterval)
	for {
		<-ticker.C
		sendStats()
	}
}

// RunAgent is the function to start agent operation
func RunAgent() {
	rand.Seed(time.Now().UnixNano())
	go RunCollectStats()
	RunSendStats()

}
