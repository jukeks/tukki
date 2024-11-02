package files

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

func RemoveFile(dbDir string, filename Filename) error {
	path := GetPath(dbDir, filename)
	return os.Remove(path)
}

func FileExists(dbDir string, filename Filename) bool {
	path := GetPath(dbDir, filename)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func FileSize(dbDir string, filename Filename) (int64, error) {
	path := GetPath(dbDir, filename)
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}
