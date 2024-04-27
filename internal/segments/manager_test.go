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

	mt := memtable.NewMemtable()
	wal, err := memtable.OpenWal(dbDir, ongoing.WalFilename, mt)
	if err != nil {
		t.Fatalf("failed to open wal: %v", err)
	}
	writeToWalAndMemtable(t, wal, mt, "key1", "value1")

	err = sm.SealCurrentSegment(mt)
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	if len(sm.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", sm.segments)
	}

	ongoing = sm.GetOnGoingSegment()
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

func writeToWalAndMemtable(t *testing.T, wal *memtable.MembtableJournal, mt memtable.Memtable, key, value string) {
	err := wal.Set(key, value)
	if err != nil {
		t.Fatalf("failed to write to wal: %v", err)
	}

	mt.Insert(key, value)
}

func TestMergeSegments(t *testing.T) {
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
	sm, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	ongoing := sm.GetOnGoingSegment()
	mt := memtable.NewMemtable()
	wal, err := memtable.OpenWal(dbDir, ongoing.WalFilename, mt)
	if err != nil {
		t.Fatalf("failed to open wal: %v", err)
	}
	writeToWalAndMemtable(t, wal, mt, "key1", "value1")
	writeToWalAndMemtable(t, wal, mt, "key2", "value2")
	wal.Close()

	err = sm.SealCurrentSegment(mt)
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	ongoing = sm.GetOnGoingSegment()
	if ongoing.Segment.Id != 1 {
		t.Fatalf("expected ongoing segment id to be 1, got %d", ongoing.Segment.Id)
	}

	mt = memtable.NewMemtable()
	wal, err = memtable.OpenWal(dbDir, ongoing.WalFilename, mt)
	if err != nil {
		t.Fatalf("failed to open wal: %v", err)
	}
	writeToWalAndMemtable(t, wal, mt, "key1", "value1")
	writeToWalAndMemtable(t, wal, mt, "key2", "value2")
	wal.Close()

	err = sm.SealCurrentSegment(mt)
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	if len(sm.segments) != 2 {
		t.Fatalf("expected segments map to have 2 element, got %v", sm.segments)
	}

	ongoing = sm.GetOnGoingSegment()
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
