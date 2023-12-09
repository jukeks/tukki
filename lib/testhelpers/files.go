package testhelpers

import (
	"os"
	"path/filepath"
)

func EnsureTempDirectory(dir string) string {
	tmpDir := os.TempDir()
	fullPath := filepath.Join(tmpDir, dir)
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		panic(err)
	}

	return fullPath
}

func CreateTempFile(tempDir string, prefix string) *os.File {
	dirPath := EnsureTempDirectory(tempDir)
	f, err := os.CreateTemp(dirPath, prefix)
	if err != nil {
		panic(err)
	}

	return f
}
