package testutil

import (
	"os"
	"path/filepath"
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

func CreateTempFile(tempDir string, filename string) *os.File {
	dirPath := ensureTempDirectory(tempDir)
	f, err := os.CreateTemp(dirPath, filename)
	if err != nil {
		panic(err)
	}

	return f
}
