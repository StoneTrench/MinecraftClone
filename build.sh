#!/bin/bash
set -e

# ===============================================
#	COLLECT
# ===============================================

echo \[Build\] Collecting info

__BUILD_NAME="Stone's MC Clone"
__BUILD_API_VERSION="0"

__BUILD_TARGET_OS="${1:-$(go env GOOS)}"
__BUILD_TARGET_ARCH="${2:-$(go env GOARCH)}"
__BUILD_MODE="${3:-"DEBUG"}"

# Validate OS (basic check)
case "$__BUILD_TARGET_OS" in
    windows|linux|darwin|js|wasip1) ;;
    *)
        echo "ERROR: Unsupported OS '$__BUILD_TARGET_OS'. Supported: windows, linux, darwin"
        exit 1
        ;;
esac

OUTPUT_NAME="main"
if [ "$__BUILD_TARGET_OS" = "windows" ]; then
    OUTPUT_NAME="main.exe"
fi
if [ "$__BUILD_TARGET_OS" = "js" ]; then
    OUTPUT_NAME="main.wasm"
fi
if [ "$__BUILD_TARGET_OS" = "wasip1" ]; then
    OUTPUT_NAME="main.wasm"
fi
if [ "$__BUILD_TARGET_ARCH" = "wasm" ]; then
  CGO_ENABLED=0
fi

PROJECT_GIT=$(awk '/module/ {print $2}' go.mod)
__BUILD_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
__BUILD_COMMIT=$(git rev-parse --short HEAD)
__BUILD_BUILT=$(date -u +%Y-%m-%d\ %H\:%M\:%S)


# ===============================================
#	CONSTRUCT
# ===============================================

echo \[Build\] Including content

rm -fr "./_out/"
mkdir -p "./_out/"

sh modbuild.sh

echo \[Build\] Building project

GOOS="$__BUILD_TARGET_OS" GOARCH="$__BUILD_TARGET_ARCH" go build -ldflags " \
  -X \"${PROJECT_GIT}/src/engine/platform.BUILD_NAME=${__BUILD_NAME}\" \
  -X \"${PROJECT_GIT}/src/engine/platform.BUILD_TAG=${__BUILD_TAG}\" \
  -X \"${PROJECT_GIT}/src/engine/platform.BUILD_BUILT=${__BUILD_BUILT}\" \
  -X \"${PROJECT_GIT}/src/engine/platform.BUILD_COMMIT=${__BUILD_COMMIT}\" \
  -X \"${PROJECT_GIT}/src/engine/platform.BUILD_MODE=${__BUILD_MODE}\" \
  -X \"${PROJECT_GIT}/src/engine/platform.BUILD_TARGET_OS=${__BUILD_TARGET_OS}\" \
  -X \"${PROJECT_GIT}/src/engine/platform.BUILD_TARGET_ARCH=${__BUILD_TARGET_ARCH}\" \
  -X \"${PROJECT_GIT}/src/engine/platform.BUILD_API_VERSION=${__BUILD_API_VERSION}\"" \
  -o "./_out/$OUTPUT_NAME"

# ===============================================
#	ZIP
# ===============================================

if [ $__BUILD_MODE != "RELEASE" ]; then
	exit
fi

echo \[Build\] Compressing project

if [ ! -f "./_out/$OUTPUT_NAME" ]; then
    echo "ERROR: Build failed - binary not found at ./_out/$OUTPUT_NAME"
    exit 1
fi

mkdir -p ./_dist/
tar -czf "./_dist/${__BUILD_NAME}_${__BUILD_TAG}-${__BUILD_TARGET_OS}-${__BUILD_TARGET_ARCH}.tar.gz" -C ./_out .