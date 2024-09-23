package db

import (
	"testing"

	"github.com/thanhpk/randstr"
)

func TestDB(t *testing.T) {
	dbDir := t.TempDir()
	database, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	key := randstr.String(10)
	value := randstr.String(16 * 1024)
	database.Set(key, value)

	storedValue, err := database.Get(key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if storedValue != value {
		t.Fatalf("stored value does not match: %s != %s", storedValue, value)
	}

	database.Close()

	// reopen database
	database, err = OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	storedValue, err = database.Get(key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if storedValue != value {
		t.Fatalf("stored value does not match: %s != %s", storedValue, value)
	}

	database.Close()
}

func TestDelete(t *testing.T) {
	database, err := OpenDatabase(t.TempDir())
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	key := randstr.String(10)
	value := randstr.String(16 * 1024)
	database.Set(key, value)

	storedValue, err := database.Get(key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if storedValue != value {
		t.Fatalf("stored value does not match: %s != %s", storedValue, value)
	}

	database.Delete(key)

	_, err = database.Get(key)
	if err == nil {
		t.Fatalf("key should not exist anymore")
	}

	database.Close()
}

func TestSegmentManager(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	ongoing := db.GetOnGoingSegment()
	if ongoing.Segment.Id != 0 {
		t.Fatalf("expected ongoing segment id to be 0, got %d", ongoing.Segment.Id)
	}

	if ongoing.WalFilename != "wal-0.journal" {
		t.Fatalf("expected ongoing segment journal filename to be 'wal-0.journal', got '%s'", ongoing.WalFilename)
	}

	if len(db.segments) != 0 {
		t.Fatalf("expected segments map to be empty, got %v", db.segments)
	}

	if len(db.operations) != 0 {
		t.Fatalf("expected operations map to be empty, got %v", db.operations)
	}

	writeLiveSegment(t, ongoing, "key1", "value1")
	nextSegment, err := db.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	if len(db.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", db.segments)
	}

	ongoing = nextSegment
	if ongoing.Segment.Id != 1 {
		t.Fatalf("expected ongoing segment id to be 1, got %d", ongoing.Segment.Id)
	}
	if ongoing.WalFilename != "wal-1.journal" {
		t.Fatalf("expected ongoing segment journal filename to be 'wal-1.journal', got '%s'", ongoing.WalFilename)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close segment manager: %v", err)
	}

	db, err = OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager again: %v", err)
	}

	if len(db.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", db.segments)
	}

	ongoing = db.GetOnGoingSegment()
	if ongoing.Segment.Id != 1 {
		t.Fatalf("expected ongoing segment id to be 1, got %d", ongoing.Segment.Id)
	}
	if ongoing.WalFilename != "wal-1.journal" {
		t.Fatalf("expected ongoing segment journal filename to be 'wal-1.journal', got '%s'",
			ongoing.WalFilename)
	}
}

func writeLiveSegment(t *testing.T, liveSegment *LiveSegment, key, value string) {
	err := liveSegment.Set(key, value)
	if err != nil {
		t.Fatalf("failed to set key-value pair: %v", err)
	}
}

func TestMergeSegments(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	ongoing := db.GetOnGoingSegment()
	writeLiveSegment(t, ongoing, "key1", "value1")
	writeLiveSegment(t, ongoing, "key2", "value2")
	if err := ongoing.Close(); err != nil {
		t.Fatalf("failed to close ongoing segment: %v", err)
	}

	nextSegment, err := db.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	ongoing = nextSegment
	if ongoing.Segment.Id != 1 {
		t.Fatalf("expected ongoing segment id to be 1, got %d", ongoing.Segment.Id)
	}

	writeLiveSegment(t, ongoing, "key3", "value3")
	writeLiveSegment(t, ongoing, "key4", "value4")
	err = ongoing.Delete("key2")
	if err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}
	if err := ongoing.Close(); err != nil {
		t.Fatalf("failed to close ongoing segment: %v", err)
	}

	nextSegment, err = db.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	if len(db.segments) != 2 {
		t.Fatalf("expected segments map to have 2 element, got %v", db.segments)
	}

	ongoing = nextSegment
	if ongoing.Segment.Id != 2 {
		t.Fatalf("expected ongoing segment id to be 2, got %d", ongoing.Segment.Id)
	}

	err = db.MergeSegments(0, 1)
	if err != nil {
		t.Fatalf("failed to merge segments: %v", err)
	}

	if len(db.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", db.segments)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close segment manager: %v", err)
	}

	db, err = OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager again: %v", err)
	}

	if len(db.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", db.segments)
	}
}

func TestSegmentRotated(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	written := 0
	kvMap := make(map[string]string)
	for written < int(db.config.WalSizeLimit)*3 {
		key := randstr.String(10)
		value := randstr.String(16 * 1024)
		err = db.Set(key, value)
		if err != nil {
			t.Fatalf("failed to set key-value pair: %v", err)
		}
		written += len(key) + len(value)
		kvMap[key] = value
	}

	if db.ongoing.Segment.Id != 2 {
		t.Fatalf("expected ongoing segment id to be 2, got %d", db.ongoing.Segment.Id)
	}

	for k, v := range kvMap {
		value, err := db.Get(k)
		if err != nil {
			t.Fatalf("failed to get key-value pair: %v", err)
		}
		if value != v {
			t.Fatalf("expected value to be %s, got %s", v, value)
		}
	}
}

func BenchmarkWrite(b *testing.B) {
	dbDir := b.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		b.Fatalf("failed to open database: %v", err)
	}
	b.Cleanup(func() {
		db.Close()
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		key := randstr.String(10)
		value := randstr.String(16 * 1024)
		b.StartTimer()
		db.Set(key, value)
	}
}
