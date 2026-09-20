package memtable_test

import (
	"testing"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/journal"
	"github.com/jukeks/tukki/internal/storage/memtable"
)

func TestWal(t *testing.T) {
	dbDir := t.TempDir()
	filename := files.Filename("wal")
	mt := memtable.NewMemtable()
	wal, err := memtable.OpenWal(dbDir, filename, journal.WriteModeAsync, mt)
	if err != nil {
		t.Fatalf("failed to open wal: %v", err)
	}

	if err := wal.Set([]byte("key1"), []byte("value1")); err != nil {
		t.Fatalf("failed to set key1: %v", err)
	}
	if err := wal.Set([]byte("key2"), []byte("value2")); err != nil {
		t.Fatalf("failed to set key2: %v", err)
	}
	if err := wal.Set([]byte("key3"), []byte("value3")); err != nil {
		t.Fatalf("failed to set key3: %v", err)
	}
	if err := wal.Delete([]byte("key2")); err != nil {
		t.Fatalf("failed to delete key2: %v", err)
	}

	if err := wal.Close(); err != nil {
		t.Fatalf("failed to close wal: %v", err)
	}

	mt = memtable.NewMemtable()
	wal, err = memtable.OpenWal(dbDir, filename, journal.WriteModeAsync, mt)
	if err != nil {
		t.Fatalf("failed to open wal: %v", err)
	}
	defer wal.Close()

	value1, found := mt.Get([]byte("key1"))
	if !found {
		t.Fatalf("key1 not found")
	}
	if string(value1.Value) != "value1" {
		t.Fatalf("expected value1 to be 'value1', got '%s'", value1.Value)
	}
	value3, found := mt.Get([]byte("key3"))
	if !found {
		t.Fatalf("key3 not found")
	}
	if string(value3.Value) != "value3" {
		t.Fatalf("expected value3 to be 'value3', got '%s'", value3.Value)
	}

	value2, found := mt.Get([]byte("key2"))
	if !found {
		t.Fatalf("key2 not found")
	}
	if !value2.Deleted {
		t.Fatalf("expected value2 to be deleted, got '%s'", value2.Value)
	}
}
