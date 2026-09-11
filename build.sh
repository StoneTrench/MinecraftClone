#!/bin/bash
set -e

# ===============================================
#	COLLECT
# ===============================================

echo \[Build\] Collecting info

__BUILD_NAME="Stone's MC Clone"
__BUILD_API_VERSION="0"

TARGET_OS="${1:-$(go env GOOS)}"
TARGET_ARCH="${2:-$(go env GOARCH)}"
__BUILD_MODE="${3:-"DEBUG"}"
__BUILD_TARGET="${TARGET_OS}/${TARGET_ARCH}"

# Validate OS (basic check)
case "$TARGET_OS" in
    windows|linux|darwin) ;;
    *)
        echo "ERROR: Unsupported OS '$TARGET_OS'. Supported: windows, linux, darwin"
        exit 1
        ;;
esac

OUTPUT_NAME="main"
if [ "$TARGET_OS" = "windows" ]; then
    OUTPUT_NAME="main.exe"
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

GOOS="$TARGET_OS" GOARCH="$TARGET_ARCH" go build -ldflags " \
  -X \"${PROJECT_GIT}/src/engine.BUILD_NAME=${__BUILD_NAME}\" \
  -X \"${PROJECT_GIT}/src/engine.BUILD_TAG=${__BUILD_TAG}\" \
  -X \"${PROJECT_GIT}/src/engine.BUILD_BUILT=${__BUILD_BUILT}\" \
  -X \"${PROJECT_GIT}/src/engine.BUILD_COMMIT=${__BUILD_COMMIT}\" \
  -X \"${PROJECT_GIT}/src/engine.BUILD_MODE=${__BUILD_MODE}\" \
  -X \"${PROJECT_GIT}/src/engine.BUILD_TARGET=${__BUILD_TARGET}\" \
  -X \"${PROJECT_GIT}/src/engine.BUILD_API_VERSION=${__BUILD_API_VERSION}\"" \
  -o "./_out/$OUTPUT_NAME"

# ===============================================
#	ZIP
# ===============================================

exit
echo \[Build\] Compressing project

if [ ! -f "./_out/$OUTPUT_NAME" ]; then
    echo "ERROR: Build failed - binary not found at ./_out/$OUTPUT_NAME"
    exit 1
fi

mkdir -p ./_dist/
tar -czf "./_dist/${__BUILD_NAME}_${__BUILD_TAG}-${TARGET_OS}-${TARGET_ARCH}.tar.gz" -C ./_out .