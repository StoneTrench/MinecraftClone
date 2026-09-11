#!/bin/bash
set -e

GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o main.wasm main.go
# GOARCH="wasm" GOOS="wasip1" go build -o main.wasm main.go
# GOARCH="wasm" GOOS="wasip1" go build -o main.wasm -buildmode=c-shared -ldflags=-checklinkname=0 main.go
# componentize-go -d ../../../bindings/wit -w main build
