package testutil

import (
	"os"
	"path/filepath"

	"github.com/thanhpk/randstr"
)

func ensureTempDirectory(dir string) string {
	tmpDir := os.TempDir()
	fullPath := filepath.Join(tmpDir, dir)
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		panic(err)
	}

	return fullPath
}

func EnsureTempDirectory() (string, func()) {
	directory := "test-tukki-" + randstr.String(10)
	path := ensureTempDirectory(directory)
	return path, func() {
		os.RemoveAll(path)
	}
}

func CreateTempFile(tempDir string, prefix string) *os.File {
	dirPath := ensureTempDirectory(tempDir)
	f, err := os.CreateTemp(dirPath, prefix)
	if err != nil {
		panic(err)
	}

	return f
}
