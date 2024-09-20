package storage

import (
	"testing"
)

func TestGetPathPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetPath did not panic")
		}
	}()
	GetPath("", "test")
}

func TestFileExists(t *testing.T) {
	dbDir := t.TempDir()

	filename := Filename("test")
	if FileExists(dbDir, filename) {
		t.Errorf("file exists")
	}

	f, err := CreateFile(dbDir, filename)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	f.Close()

	if !FileExists(dbDir, filename) {
		t.Errorf("file does not exist")
	}
}
