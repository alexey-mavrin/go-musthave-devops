package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type statType int

const (
	statTypeGauge = iota
	statTypeCounter
)

var (
	errNotEngoughParts = fmt.Errorf("Too few parts in URL")
	errWrongOp         = fmt.Errorf("Unknown Operation")
	errWrongType       = fmt.Errorf("Unknown Type")
	errBadName         = fmt.Errorf("Bad Stat Name")
	errBadValue        = fmt.Errorf("Bad Value")
)

type statReq struct {
	statType     statType
	name         string
	valueCounter int64
	valueGauge   float64
}

func parseURL(url string) (statReq, error) {
	var stat statReq
	parts := strings.Split(url, "/")

	if len(parts) != 5 {
		return stat, errNotEngoughParts
	}

	op := parts[1]
	typ := parts[2]
	name := parts[3]
	rawVal := parts[4]

	if op != "update" {
		return stat, errWrongOp
	}

	if typ != "counter" && typ != "gauge" {
		return stat, errWrongType
	}

	if len(name) == 0 {
		return stat, errBadName
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

	return stat, nil
}

// Handler — обработчик запроса.
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	url := string(r.URL.Path)
	_, err := parseURL(url)

	if err == errNotEngoughParts {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err == errWrongOp {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err == errWrongType {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err == errBadName {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err == errBadValue {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
