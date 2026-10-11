//go:build !js

package filesystem

import (
	"io"
	"os"
)

func open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}
func create(name string) (io.ReadWriteCloser, error) {
	return os.Create(name)
}
func writeFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}
func readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
func readDir(name string) ([]DirEntry, error) {
	return os.ReadDir(name)
}
