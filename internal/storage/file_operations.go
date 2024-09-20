package storage

import (
	"os"
	"path/filepath"
)

type Filename string

func GetPath(dbDir string, filename Filename) string {
	if dbDir == "" {
		panic("dbDir is empty")
	}
	if filename == "" {
		panic("filename is empty")
	}

	return filepath.Join(dbDir, string(filename))
}

func OpenFile(dbDir string, filename Filename) (*os.File, error) {
	path := GetPath(dbDir, filename)
	return os.Open(path)
}

func CreateFile(dbDir string, filename Filename) (*os.File, error) {
	path := GetPath(dbDir, filename)
	return os.Create(path)
}

func FileExists(dbDir string, filename Filename) bool {
	path := GetPath(dbDir, filename)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
