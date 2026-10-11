package filesystem

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

type DirEntry = fs.DirEntry

func Open(name string) (io.ReadCloser, error) {
	return open(name)
}
func Create(name string) (io.ReadWriteCloser, error) {
	return create(name)
}
func WriteFile(name string, data []byte) error {
	return writeFile(name, data)
}
func ReadFile(name string) ([]byte, error) {
	return readFile(name)
}
func ReadDir(name string) ([]fs.DirEntry, error) {
	return readDir(name)
}
func OpenOrCreate(name string) (io.ReadCloser, error) {
	file, err := Open(name)
	if errors.Is(err, os.ErrNotExist) {
		err = file.Close()
		if err != nil {
			return nil, err
		}
		file, err = Create(name)
	}
	return file, err
}
