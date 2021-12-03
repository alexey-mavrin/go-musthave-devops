package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
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
	statistics.mu.Lock()
	defer statistics.mu.Unlock()

	switch m.MType {
	case strTypCounter:
		val, ok := statistics.counters[m.ID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"Status":"Not Found"}`))
			return
		}
		m.Delta = &val
	case strTypGauge:
		val, ok := statistics.gauges[m.ID]
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

	statistics.mu.Lock()
	defer statistics.mu.Unlock()

	if typ == strTypCounter {
		val, ok := statistics.counters[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
		w.Write([]byte(fmt.Sprint(val)))
	} else if typ == strTypGauge {
		val, ok := statistics.gauges[name]
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
	statistics.mu.Lock()

	cNames := make([]string, 0, len(statistics.counters))
	for k := range statistics.counters {
		cNames = append(cNames, k)
	}
	sort.Strings(cNames)

	gNames := make([]string, 0, len(statistics.gauges))
	for k := range statistics.gauges {
		gNames = append(gNames, k)
	}
	sort.Strings(gNames)

	for _, n := range cNames {
		str = str + fmt.Sprintf("%s %v\n", n, statistics.counters[n])
	}
	for _, n := range gNames {
		str = str + fmt.Sprintf("%s %v\n", n, statistics.gauges[n])
	}

	statistics.mu.Unlock()
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
	statistics.mu.Lock()
	switch stat.statType {
	case statTypeCounter:
		statistics.counters[stat.name] += stat.valueCounter
	case statTypeGauge:
		statistics.gauges[stat.name] = stat.valueGauge
	}
	statistics.mu.Unlock()
}

// Router return chi.Router for testing and actual work
func Router() chi.Router {
	statistics.counters = make(map[string]int64)
	statistics.gauges = make(map[string]float64)
	r := chi.NewRouter()
	r.Get("/", DumpHandler)
	r.Get("/value/{typ}/{name}", MetricHandler)
	r.Post("/value/", JSONMetricHandler)
	r.Post("/update/", JSONUpdateHandler)
	r.Post("/update/{typ}/{name}/", Handler400)
	r.Post("/update/{typ}/{name}/{rawVal}", UpdateHandler)
	return r
}
