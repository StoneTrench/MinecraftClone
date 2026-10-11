//go:build js && wasm

package main

import (
	"embed"

	"github.com/StoneTrench/go-mc-clone/src/engine/platform/filesystem"
)

//go:embed include
var modsFS embed.FS

func init() {
	filesystem.EmbedInclude(modsFS)
}
