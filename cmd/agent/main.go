package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
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
	contentType   = "text/plain"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func sendStat(statString string) {
	_, err := http.Post(defaultServer+statString, contentType, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func makeStatStringGauge(name string, value float64) string {
	str := fmt.Sprintf("/update/gauge/%s/%G", name, value)
	return str
}

func makeStatStringCounter(name string, value int64) string {
	str := fmt.Sprintf("/update/counter/%s/%d", name, value)
	return str
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

func runCollectStats() {
	ticker := time.NewTicker(pollInterval)
	for {
		select {
		case <-ticker.C:
			collectStats()
		}
	}
}

func sendStats() {
	myStatData.mu.Lock()
	PollCount := makeStatStringCounter("PollCount", myStatData.PollCount)
	RandomValue := makeStatStringCounter("RandomValue", int64(myStatData.RandomValue))
	Alloc := makeStatStringGauge("Alloc", float64(myStatData.memStats.Alloc))
	BuckHashSys := makeStatStringGauge("BuckHashSys", float64(myStatData.memStats.BuckHashSys))
	Frees := makeStatStringGauge("Frees", float64(myStatData.memStats.Frees))
	GCCPUFraction := makeStatStringGauge("GCCPUFraction", float64(myStatData.memStats.GCCPUFraction))
	GCSys := makeStatStringGauge("GCSys", float64(myStatData.memStats.GCSys))
	HeapAlloc := makeStatStringGauge("HeapAlloc", float64(myStatData.memStats.HeapAlloc))
	HeapIdle := makeStatStringGauge("HeapIdle", float64(myStatData.memStats.HeapIdle))
	HeapInuse := makeStatStringGauge("HeapInuse", float64(myStatData.memStats.HeapInuse))
	HeapObjects := makeStatStringGauge("HeapObjects", float64(myStatData.memStats.HeapObjects))
	HeapReleased := makeStatStringGauge("HeapReleased", float64(myStatData.memStats.HeapReleased))
	HeapSys := makeStatStringGauge("HeapSys", float64(myStatData.memStats.HeapSys))
	LastGC := makeStatStringGauge("LastGC", float64(myStatData.memStats.LastGC))
	Lookups := makeStatStringGauge("Lookups", float64(myStatData.memStats.Lookups))
	MCacheInuse := makeStatStringGauge("MCacheInuse", float64(myStatData.memStats.MCacheInuse))
	MCacheSys := makeStatStringGauge("MCacheSys", float64(myStatData.memStats.MCacheSys))
	MSpanInuse := makeStatStringGauge("MSpanInuse", float64(myStatData.memStats.MSpanInuse))
	MSpanSys := makeStatStringGauge("MSpanSys", float64(myStatData.memStats.MSpanSys))
	Mallocs := makeStatStringGauge("Mallocs", float64(myStatData.memStats.Mallocs))
	NextGC := makeStatStringGauge("NextGC", float64(myStatData.memStats.NextGC))
	NumForcedGC := makeStatStringGauge("NumForcedGC", float64(myStatData.memStats.NumForcedGC))
	NumGC := makeStatStringGauge("NumGC", float64(myStatData.memStats.NumGC))
	OtherSys := makeStatStringGauge("OtherSys", float64(myStatData.memStats.OtherSys))
	PauseTotalNs := makeStatStringGauge("PauseTotalNs", float64(myStatData.memStats.PauseTotalNs))
	StackInuse := makeStatStringGauge("StackInuse", float64(myStatData.memStats.StackInuse))
	StackSys := makeStatStringGauge("StackSys", float64(myStatData.memStats.StackSys))
	Sys := makeStatStringGauge("Sys", float64(myStatData.memStats.Sys))
	myStatData.mu.Unlock()

	sendStat(PollCount)
	sendStat(RandomValue)
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
	sendStat(Sys)
}

func runSendStats() {
	ticker := time.NewTicker(reportInterval)
	for {
		select {
		case <-ticker.C:
			sendStats()
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	go runCollectStats()
	runSendStats()
}
