//go:build js && wasm

package filesystem

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"sync"
)

// TODO: this doesn't work at all, please fix for wasm compat on the web!

var (
	embed_include embed.FS
	memMu         sync.RWMutex
	memFiles      = map[string][]byte{}
	memDirs       = map[string]struct{}{}
)

func embedPath(name string) string {
	rel := strings.TrimPrefix(path.Clean("/"+name), "/")
	return "include/" + rel
}

func open(name string) (io.ReadCloser, error) {
	return os.Open(embedPath(name))
}
func create(name string) (io.ReadWriteCloser, error) {
	return os.Create(embedPath(name))
}
func readFile(name string) ([]byte, error) {
	key := embedPath(name)
	memMu.RLock()
	if b, ok := memFiles[key]; ok {
		memMu.RUnlock()
		return append([]byte(nil), b...), nil
	}
	memMu.RUnlock()
	return embed_include.ReadFile(key)
}

func writeFile(name string, data []byte) error {
	key := embedPath(name)
	memMu.Lock()
	memFiles[key] = append([]byte(nil), data...)
	memMu.Unlock()
	return nil
}
func readDir(name string) ([]fs.DirEntry, error) {
	return embed_include.ReadDir(embedPath(name))
}

func EmbedInclude(e embed.FS) {
	embed_include = e
}
