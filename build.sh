# #!/bin/sh


sh build_assets.sh

PROJECT_NAME="Stone's MC Clone"

TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
COMMIT=$(git rev-parse --short HEAD)
PROJECT_GIT=$(awk '/module/ {print $2}' go.mod)
BUILT=$(date -u +%Y-%m-%d-%H%M%S)
go build -ldflags " \
  -X \"${PROJECT_GIT}/game.PROJECT_NAME=${PROJECT_NAME}\" \
  -X \"${PROJECT_GIT}/game.PROJECT_TAG=${TAG}\" \
  -X \"${PROJECT_GIT}/game.PROJECT_COMMIT=${COMMIT}\" \
  -X \"${PROJECT_GIT}/game.PROJECT_BUILT=${BUILT}\"" \
  -o ./_out/main.exe

