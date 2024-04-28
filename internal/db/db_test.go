package db

import (
	"testing"

	testutil "github.com/jukeks/tukki/tests/util"
	"github.com/thanhpk/randstr"
)

func TestDB(t *testing.T) {
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
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
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
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

	database.Delete(key)

	_, err = database.Get(key)
	if err == nil {
		t.Fatalf("key should not exist anymore")
	}

	database.Close()
}

func TestSegmentManager(t *testing.T) {
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
	sm, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	ongoing := sm.GetOnGoingSegment()
	if ongoing.Segment.Id != 0 {
		t.Fatalf("expected ongoing segment id to be 0, got %d", ongoing.Segment.Id)
	}

	if ongoing.WalFilename != "wal-0.journal" {
		t.Fatalf("expected ongoing segment journal filename to be 'wal-0.journal', got '%s'", ongoing.WalFilename)
	}

	if len(sm.segments) != 0 {
		t.Fatalf("expected segments map to be empty, got %v", sm.segments)
	}

	if len(sm.operations) != 0 {
		t.Fatalf("expected operations map to be empty, got %v", sm.operations)
	}

	writeLiveSegment(t, ongoing, "key1", "value1")
	nextSegment, err := sm.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	if len(sm.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", sm.segments)
	}

	ongoing = nextSegment
	if ongoing.Segment.Id != 1 {
		t.Fatalf("expected ongoing segment id to be 1, got %d", ongoing.Segment.Id)
	}
	if ongoing.WalFilename != "wal-1.journal" {
		t.Fatalf("expected ongoing segment journal filename to be 'wal-1.journal', got '%s'", ongoing.WalFilename)
	}

	err = sm.Close()
	if err != nil {
		t.Fatalf("failed to close segment manager: %v", err)
	}

	sm, err = OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager again: %v", err)
	}

	if len(sm.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", sm.segments)
	}

	ongoing = sm.GetOnGoingSegment()
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
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
	sm, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	ongoing := sm.GetOnGoingSegment()
	writeLiveSegment(t, ongoing, "key1", "value1")
	writeLiveSegment(t, ongoing, "key2", "value2")
	if err := ongoing.Close(); err != nil {
		t.Fatalf("failed to close ongoing segment: %v", err)
	}

	nextSegment, err := sm.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	ongoing = nextSegment
	if ongoing.Segment.Id != 1 {
		t.Fatalf("expected ongoing segment id to be 1, got %d", ongoing.Segment.Id)
	}

	writeLiveSegment(t, ongoing, "key1", "value1")
	writeLiveSegment(t, ongoing, "key2", "value2")
	if err := ongoing.Close(); err != nil {
		t.Fatalf("failed to close ongoing segment: %v", err)
	}

	nextSegment, err = sm.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	if len(sm.segments) != 2 {
		t.Fatalf("expected segments map to have 2 element, got %v", sm.segments)
	}

	ongoing = nextSegment
	if ongoing.Segment.Id != 2 {
		t.Fatalf("expected ongoing segment id to be 2, got %d", ongoing.Segment.Id)
	}

	err = sm.MergeSegments(0, 1)
	if err != nil {
		t.Fatalf("failed to merge segments: %v", err)
	}

	if len(sm.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", sm.segments)
	}

	err = sm.Close()
	if err != nil {
		t.Fatalf("failed to close segment manager: %v", err)
	}

	sm, err = OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager again: %v", err)
	}

	if len(sm.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", sm.segments)
	}
}
