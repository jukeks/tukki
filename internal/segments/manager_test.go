package segments

import (
	"testing"

	"github.com/jukeks/tukki/internal/memtable"
	testutil "github.com/jukeks/tukki/tests/util"
	"github.com/thanhpk/randstr"
)

func TestSegmentManager(t *testing.T) {
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
	sm, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	ongoing := sm.GetOnGoingSegment()
	if ongoing.Id != 0 {
		t.Fatalf("expected ongoing segment id to be 0, got %d", ongoing.Id)
	}

	if ongoing.JournalFilename != "wal-0.journal" {
		t.Fatalf("expected ongoing segment journal filename to be 'wal-0.journal', got '%s'", ongoing.JournalFilename)
	}

	if len(sm.segments) != 0 {
		t.Fatalf("expected segments map to be empty, got %v", sm.segments)
	}

	if len(sm.operations) != 0 {
		t.Fatalf("expected operations map to be empty, got %v", sm.operations)
	}

	mt := memtable.NewMemtable()
	mt.Insert("key1", "value1")

	err = sm.SealCurrentSegment(mt)
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	if len(sm.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", sm.segments)
	}

	ongoing = sm.GetOnGoingSegment()
	if ongoing.Id != 1 {
		t.Fatalf("expected ongoing segment id to be 1, got %d", ongoing.Id)
	}
	if ongoing.JournalFilename != "wal-1.journal" {
		t.Fatalf("expected ongoing segment journal filename to be 'wal-1.journal', got '%s'", ongoing.JournalFilename)
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
	if ongoing.Id != 1 {
		t.Fatalf("expected ongoing segment id to be 1, got %d", ongoing.Id)
	}
	if ongoing.JournalFilename != "wal-1.journal" {
		t.Fatalf("expected ongoing segment journal filename to be 'wal-1.journal', got '%s'", ongoing.JournalFilename)
	}
}

func TestMergeSegments(t *testing.T) {
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
	sm, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	mt := memtable.NewMemtable()
	mt.Insert("key1", "value1")
	mt.Insert("key2", "value2")

	err = sm.SealCurrentSegment(mt)
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	mt = memtable.NewMemtable()
	mt.Insert("key3", "value3")
	mt.Insert("key4", "value4")

	err = sm.SealCurrentSegment(mt)
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	if len(sm.segments) != 2 {
		t.Fatalf("expected segments map to have 2 element, got %v", sm.segments)
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
