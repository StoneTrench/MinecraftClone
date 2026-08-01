#!/bin/bash
set -e


# sh build_mod.sh
# sh build_assets.sh

cd ./_out/mods/test/
GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o main.wasm main.go
cd ../../../

PROJECT_NAME="Stone's MC Clone"

TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
COMMIT=$(git rev-parse --short HEAD)
PROJECT_GIT=$(awk '/module/ {print $2}' go.mod)
BUILT=$(date -u +%Y-%m-%d\ %H\:%M\:%S)
MODE="DEBUG"

go build -ldflags " \
  -X \"${PROJECT_GIT}/metadata.BUILD_NAME=${PROJECT_NAME}\" \
  -X \"${PROJECT_GIT}/metadata.BUILD_TAG=${TAG}\" \
  -X \"${PROJECT_GIT}/metadata.BUILD_BUILT=${BUILT}\" \
  -X \"${PROJECT_GIT}/metadata.BUILD_COMMIT=${COMMIT}\" \
  -X \"${PROJECT_GIT}/metadata.BUILD_MODE=${MODE}\"" \
  -o ./_out/main.exe

