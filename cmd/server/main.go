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

func parseURL(url string) (statReq, error) {
	var stat statReq
	parts := strings.Split(url, "/")

	var op, typ, name, rawVal string

	if len(parts) > 1 {
		op = parts[1]
	}
	if len(parts) > 2 {
		typ = parts[2]
	}
	if len(parts) > 3 {
		name = parts[3]
	}
	if len(parts) > 4 {
		rawVal = parts[4]
	}

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

	return stat, nil
}

// Handler — обработчик запроса.
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	url := string(r.URL.Path)
	_, err := parseURL(url)

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
