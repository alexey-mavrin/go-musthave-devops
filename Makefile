include .env


PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

MAKEFLAGS += --silent

exec:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(run)


test:
	go test -v ./...

build: build_agent build_server build_staticlint

build_agent:
	echo "Building agent..."
	go build -o agent_main cmd/agent/main.go

build_server:
	echo "Building server..."
	go build -o server_main cmd/server/main.go

build_staticlint:
	echo "Building staticlint..."
	go build -o staticlint_main cmd/staticlint/main.go

kill_server:
	echo "Stopping server..."
	pkill server_main || true

run_server: build_server kill_server
	echo "Starting server..."
	./server_main > /dev/null 2>&1 &

run_agent: build_agent build_server run_server
	echo "Loading data into server..."
	./agent_main -p 100ms -r 200ms > /dev/null 2>&1 &
	sleep 3
	pkill agent_main

run_load:
	echo "Running read load..."
	bombardier -c 500 -d 100s  http://localhost:8080/

test_load: build_agent build_server run_server run_agent run_load

custom_check:
	go run ./cmd/staticlint/main.go ./...

server_keys:
	mkdir -p keys
	go run ./cmd/genkeys/main.go keys

