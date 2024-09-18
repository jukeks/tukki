package db

import "testing"

func TestSnapshots(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}
	defer db.Close()

	if err := db.Set("key1", "value1"); err != nil {
		t.Fatalf("failed to set key1: %v", err)
	}

	ss := db.Snapshot()
	if ss == nil {
		t.Fatalf("snapshot is nil")
	}
}

func TestMarshalling(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	if err := db.Set("key1", "value1"); err != nil {
		t.Fatalf("failed to set key1: %v", err)
	}

	ss := db.Snapshot()
	if ss == nil {
		t.Fatalf("snapshot is nil")
	}

	buff, err := ss.Marshal()
	if err != nil {
		t.Fatalf("failed to marshal snapshot: %v", err)
	}

	ss2, err := UnmarshalSnapshot(buff)
	if err != nil {
		t.Fatalf("failed to unmarshal snapshot: %v", err)
	}

	if ss2 == nil {
		t.Fatalf("snapshot is nil")
	}
}
