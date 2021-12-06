package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
)

type statType int

type serverConfig struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

// Config stores server configuration
var Config serverConfig = serverConfig{}

var mu sync.Mutex

var statistics struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

const (
	statTypeGauge = iota
	statTypeCounter
)

const (
	strTypGauge   = "gauge"
	strTypCounter = "counter"
)

var (
	errWrongOp   = fmt.Errorf("unknown operation")
	errWrongType = fmt.Errorf("unknown type")
	errNoName    = fmt.Errorf("no stat name")
	errBadValue  = fmt.Errorf("bad value")
)

type statReq struct {
	statType     statType
	name         string
	valueCounter int64
	valueGauge   float64
}

func init() {
	statistics.Counters = make(map[string]int64)
	statistics.Gauges = make(map[string]float64)
}

// StartServer starts server
func StartServer() {
	if Config.Restore && Config.StoreFile != "" {
		loadStats()
	}

	if Config.StoreInterval > 0 && Config.StoreFile != "" {
		go statSaver()
	}
	r := Router()

	go http.ListenAndServe(Config.Address, r)

	signalChannel := make(chan os.Signal, 2)
	// Сервер должен штатно завершаться по сигналам: syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-signalChannel
	switch sig {
	case os.Interrupt:
		log.Println("sigint")
	case syscall.SIGTERM:
		log.Println("sigterm")
	case syscall.SIGINT:
		log.Println("sigint")
	case syscall.SIGQUIT:
		log.Println("sigquit")
	}

	mu.Lock()
	log.Print("server finished, storint stats")
	storeStats()
	mu.Unlock()

}

func statSaver() {
	ticker := time.NewTicker(Config.StoreInterval)
	for {
		<-ticker.C
		mu.Lock()
		storeStats()
		mu.Unlock()
	}

}

func storeStats() {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(Config.StoreFile, flags, 0644)
	if err != nil {
		log.Fatal("cannot open file for writing: ", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(statistics); err != nil {
		log.Fatal("cannot encode statistics: ", err)
	}
}

func loadStats() {
	flags := os.O_RDONLY
	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(Config.StoreFile, flags, 0)
	if err != nil {
		log.Print("cannot open file for reading ", err)
		return
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&statistics); err != nil {
		log.Fatal("cannot decode statistics ", err)
	}
}

func parseReq(r *http.Request) (statReq, error) {
	var stat statReq
	typ := chi.URLParam(r, "typ")
	name := chi.URLParam(r, "name")
	rawVal := chi.URLParam(r, "rawVal")

	if len(name) == 0 {
		return stat, errNoName
	}

	switch typ {
	case strTypCounter:
		stat.statType = statTypeCounter
		val, err := strconv.Atoi(rawVal)
		if err != nil {
			return stat, errBadValue
		}
		stat.valueCounter = int64(val)
	case strTypGauge:
		stat.statType = statTypeGauge
		val, err := strconv.ParseFloat(rawVal, 64)
		if err != nil {
			return stat, errBadValue
		}
		stat.valueGauge = val
	default:
		return stat, errWrongType
	}

	stat.name = name

	return stat, nil
}

// Handler400 — return 400
func Handler400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
}

// JSONMetricHandler prints all available metrics
func JSONMetricHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.Method, " ", r.URL)
	defer fmt.Println("")
	body, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Status":"Internal Server Error"}`))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Status":"Bad Request"}`))
		return
	}

	var m common.Metrics

	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Status":"Bad Request"}`))
		return
	}

	if m.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Status":"Bad Request"}`))
		return
	}

	fmt.Print(" type: ", m.MType, ", id: ", m.ID)
	mu.Lock()
	defer mu.Unlock()

	switch m.MType {
	case strTypCounter:
		val, ok := statistics.Counters[m.ID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"Status":"Not Found"}`))
			return
		}
		m.Delta = &val
	case strTypGauge:
		val, ok := statistics.Gauges[m.ID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"Status":"Not Found"}`))
			return
		}
		m.Value = &val
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Status":"Bad Request"}`))
		return
	}
	ret, _ := json.Marshal(m)
	w.Write(ret)
}

// MetricHandler prints all available metrics
func MetricHandler(w http.ResponseWriter, r *http.Request) {
	typ := chi.URLParam(r, "typ")
	name := chi.URLParam(r, "name")
	fmt.Println("GET", typ, name)

	mu.Lock()
	defer mu.Unlock()

	if typ == strTypCounter {
		val, ok := statistics.Counters[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
		w.Write([]byte(fmt.Sprint(val)))
	} else if typ == strTypGauge {
		val, ok := statistics.Gauges[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
		w.Write([]byte(fmt.Sprint(val)))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}
}

// DumpHandler prints all available metrics
func DumpHandler(w http.ResponseWriter, r *http.Request) {
	str := ""
	mu.Lock()

	cNames := make([]string, 0, len(statistics.Counters))
	for k := range statistics.Counters {
		cNames = append(cNames, k)
	}
	sort.Strings(cNames)

	gNames := make([]string, 0, len(statistics.Gauges))
	for k := range statistics.Gauges {
		gNames = append(gNames, k)
	}
	sort.Strings(gNames)

	for _, n := range cNames {
		str = str + fmt.Sprintf("%s %v\n", n, statistics.Counters[n])
	}
	for _, n := range gNames {
		str = str + fmt.Sprintf("%s %v\n", n, statistics.Gauges[n])
	}

	mu.Unlock()
	w.Write([]byte(str))
}

// JSONUpdateHandler — stores metrics in server from json updates
func JSONUpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.Method, " ", r.URL)
	defer fmt.Println("")

	body, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Status":"Internal Server Error"}`))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Status":"Bad Request"}`))
		return
	}

	var m common.Metrics

	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Status":"Bad Request"}`))
		return
	}

	fmt.Print(" type: ", m.MType, ", id: ", m.ID)
	var stat statReq
	switch m.MType {
	case strTypCounter:
		stat.statType = statTypeCounter
		stat.valueCounter = *m.Delta
		fmt.Print(", delta: ", *m.Delta)
	case strTypGauge:
		stat.statType = statTypeGauge
		stat.valueGauge = *m.Value
		fmt.Print(", value: ", *m.Value)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"Status":"Not Implemented"}`))
		return
	}

	if m.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Status":"Bad Request"}`))
		return
	}

	stat.name = m.ID

	updateStatStorage(stat)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"Status":"OK"}`))
}

// UpdateHandler — stores metrics in server
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	stat, err := parseReq(r)

	switch err {
	case errWrongOp, errNoName:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	case errWrongType:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not Implemented"))
		return
	case errBadValue:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	updateStatStorage(stat)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func updateStatStorage(stat statReq) {
	mu.Lock()
	switch stat.statType {
	case statTypeCounter:
		statistics.Counters[stat.name] += stat.valueCounter
	case statTypeGauge:
		statistics.Gauges[stat.name] = stat.valueGauge
	}

	if Config.StoreInterval == 0 && Config.StoreFile != "" {
		storeStats()
	}

	mu.Unlock()
}

// Router return chi.Router for testing and actual work
func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", DumpHandler)
	r.Get("/value/{typ}/{name}", MetricHandler)
	r.Post("/value/", JSONMetricHandler)
	r.Post("/update/", JSONUpdateHandler)
	r.Post("/update/{typ}/{name}/", Handler400)
	r.Post("/update/{typ}/{name}/{rawVal}", UpdateHandler)
	return r
}
