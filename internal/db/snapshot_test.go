package db

import (
	"testing"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/sstable"
)

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

	ss, err := db.Snapshot()
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}
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

	ss, err := db.Snapshot()
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}
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

func TestRestore(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	if err := db.Set("key1", "value1"); err != nil {
		t.Fatalf("failed to set key1: %v", err)
	}

	ss, err := db.Snapshot()
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}
	if ss == nil {
		t.Fatalf("snapshot is nil")
	}

	if err := db.Set("key2", "value2"); err != nil {
		t.Fatalf("failed to set key2: %v", err)
	}

	db.Close()

	if _, err := db.Restore(ss); err != nil {
		t.Fatalf("failed to restore snapshot: %v", err)
	}

	db, err = OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	if _, err := db.Get("key1"); err != nil {
		t.Fatalf("failed to get key1: %v", err)
	}

	if _, err := db.Get("key2"); err == nil {
		t.Fatalf("key2 should not exist")
	}
}

func TestRestoreMissingSegments(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	if err := db.Set("key1", "value1"); err != nil {
		t.Fatalf("failed to set key1: %v", err)
	}

	if _, err := db.SealCurrentSegment(); err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	ss, err := db.Snapshot()
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}
	if ss == nil {
		t.Fatalf("snapshot is nil")
	}

	dbDir2 := t.TempDir()
	db2, err := OpenDatabase(dbDir2)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	result, err := db2.Restore(ss)
	if err != nil {
		t.Fatalf("failed to restore snapshot: %v", err)
	}

	if len(result.MissingSegments) != 1 {
		t.Fatalf("expected 1 missing segment, got %d", len(result.MissingSegments))
	}
}

func TestRestoreSegments(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	if err := db.Set("key1", "value1"); err != nil {
		t.Fatalf("failed to set key1: %v", err)
	}

	if _, err := db.SealCurrentSegment(); err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	ss, err := db.Snapshot()
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}
	if ss == nil {
		t.Fatalf("snapshot is nil")
	}

	dbDir2 := t.TempDir()
	db2, err := OpenDatabase(dbDir2)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}
	db2.Close()

	result, err := db2.Restore(ss)
	if err != nil {
		t.Fatalf("failed to restore snapshot: %v", err)
	}

	if len(result.MissingSegments) != 1 {
		t.Fatalf("expected 1 missing segment, got %d", len(result.MissingSegments))
	}

	f, err := files.OpenFile(dbDir, result.MissingSegments[0].SegmentFile)
	if err != nil {
		t.Fatalf("failed to open segment file: %v", err)
	}
	defer f.Close()
	reader := sstable.NewSSTableReader(f)
	if err := db2.RestoreSegment(result.MissingSegments[0], reader); err != nil {
		t.Fatalf("failed to restore segment: %v", err)
	}

	db2, err = OpenDatabase(dbDir2)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	if _, err := db2.Get("key1"); err != nil {
		t.Fatalf("failed to get key1: %v", err)
	}
}

func TestSnapshotReleasesPinnedSegments(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}
	defer db.Close()

	writePairsAndSeal(t, db, Pair{"key1", "value1"}, Pair{"key2", "value2"})
	writePairsAndSeal(t, db, Pair{"key3", "value3"}, Pair{"key4", "value4"})
	writePairsAndSeal(t, db, Pair{"key5", "value5"}, Pair{"key6", "value6"})
	writePairsAndSeal(t, db, Pair{"key7", "value7"}, Pair{"key8", "value8"})

	_, err = db.Snapshot()
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}

	if len(db.lastSnapshotSegments) != 4 {
		t.Fatalf("expected 3 last snapshot segments, got %d", len(db.lastSnapshotSegments))
	}

	if err := db.Compact(); err != nil {
		t.Fatalf("failed to compact: %v", err)
	}

	if len(db.freedButNotRemoved) != 4 {
		t.Fatalf("expected 3 freed but not removed segments, got %d", len(db.freedButNotRemoved))
	}

	_, err = db.Snapshot()
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}

	if len(db.lastSnapshotSegments) != 1 {
		t.Fatalf("expected 1 last snapshot segments, got %d", len(db.lastSnapshotSegments))
	}
	if len(db.freedButNotRemoved) != 0 {
		t.Fatalf("expected 0 freed but not removed segments, got %d", len(db.freedButNotRemoved))
	}
}
