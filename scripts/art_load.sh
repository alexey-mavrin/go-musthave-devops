#!/bin/bash

set -ue
pushd "$(dirname "$0")/.."

log() {
    echo "$@" > /dev/stderr
}

build_server() {
    log "Building server..."
    go build -o server_main cmd/server/main.go
}

build_agent() {
    log "Building agent..."
    go build -o agent_main cmd/agent/main.go
}

run_server() {
    log "Starting server..."
    ./server_main > /dev/null 2>&1 &
}

kill_server() {
    log "Stopping server..."
    pkill server_main
}

run_agent() {
    log "Loading data into server..."
    ./agent_main -p 100ms -r 200ms > /dev/null 2>&1 &
    sleep 3
    pkill agent_main
}

run_read() {
    log "Running read load..."
    bombardier -c 500 -d 100s  http://localhost:8080/
}

kill_server || true
build_server
build_agent
run_server
run_agent
run_read

popd
