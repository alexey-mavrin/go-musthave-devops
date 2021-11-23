package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type statType int

var statistics struct {
	mu       sync.Mutex
	counters map[string]int64
	gauges   map[string]float64
}

const (
	statTypeGauge = iota
	statTypeCounter
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

func parseReq(r *http.Request) (statReq, error) {
	var stat statReq
	op := chi.URLParam(r, "op")
	typ := chi.URLParam(r, "typ")
	name := chi.URLParam(r, "name")
	rawVal := chi.URLParam(r, "rawVal")

	if op != "update" {
		return stat, errWrongOp
	}

	if typ != "counter" && typ != "gauge" {
		return stat, errWrongType
	}

	if len(name) == 0 {
		return stat, errNoName
	}

	if len(rawVal) == 0 {
		return stat, errBadValue
	}

	if typ == "counter" {
		stat.statType = statTypeCounter
		val, err := strconv.Atoi(rawVal)
		if err != nil {
			return stat, errBadValue
		}
		stat.valueCounter = int64(val)
	} else if typ == "gauge" {
		stat.statType = statTypeGauge
		val, err := strconv.ParseFloat(rawVal, 64)
		if err != nil {
			return stat, errBadValue
		}
		stat.valueGauge = val
	}

	stat.name = name

	return stat, nil
}

// Handler400 — return 400
func Handler400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("OK"))
}

// MetricHandler prints all available metrics
func MetricHandler(w http.ResponseWriter, r *http.Request) {
	typ := chi.URLParam(r, "typ")
	name := chi.URLParam(r, "name")
	fmt.Println("GET", typ, name)

	statistics.mu.Lock()
	defer statistics.mu.Unlock()

	if typ == "counter" {
		val, ok := statistics.counters[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte(fmt.Sprint(val)))
	} else if typ == "gauge" {
		val, ok := statistics.gauges[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte(fmt.Sprint(val)))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("OK"))
	}
}

// DumpHandler prints all available metrics
func DumpHandler(w http.ResponseWriter, r *http.Request) {
	str := ""
	statistics.mu.Lock()
	for k, v := range statistics.counters {
		str = str + fmt.Sprintf("%s %v\n", k, v)
	}
	for k, v := range statistics.gauges {
		str = str + fmt.Sprintf("%s %v\n", k, v)
	}
	statistics.mu.Unlock()
	w.Write([]byte(str))
}

// UpdateHandler — обработчик запроса.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	stat, err := parseReq(r)

	if err == errWrongOp {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err == errWrongType {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if err == errNoName {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err == errBadValue {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statistics.mu.Lock()
	if stat.statType == statTypeCounter {
		statistics.counters[stat.name] += stat.valueCounter
	}
	if stat.statType == statTypeGauge {
		statistics.gauges[stat.name] = stat.valueGauge
	}
	statistics.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Router return chi.Router for testing and actual work
func Router() chi.Router {
	statistics.counters = make(map[string]int64)
	statistics.gauges = make(map[string]float64)
	r := chi.NewRouter()
	r.Get("/", DumpHandler)
	r.Get("/value/{typ}/{name}", MetricHandler)
	r.Post("/{op}/{typ}/{name}/", Handler400)
	r.Post("/{op}/{typ}/{name}/{rawVal}", UpdateHandler)
	return r
}

func main() {
	r := Router()
	http.ListenAndServe(":8080", r)
}
