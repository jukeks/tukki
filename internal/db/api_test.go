package db

import (
	"testing"

	"github.com/jukeks/tukki/internal/segments"
	testutil "github.com/jukeks/tukki/testutil"
	"github.com/thanhpk/randstr"
)

func TestSegmentManagerGetSegmentsSorted(t *testing.T) {
	db := Database{
		segments: map[segments.SegmentId]segments.SegmentMetadata{
			1: {
				Id:          1,
				SegmentFile: "segment1",
			},
			3: {
				Id:          3,
				SegmentFile: "segment3",
			},
			2: {
				Id:          2,
				SegmentFile: "segment2",
			},
		},
	}

	segments := db.getSegmentsSorted()
	if len(segments) != 3 {
		t.Errorf("expected 3 segments, got %d", len(segments))
	}

	if segments[0].Id != 3 {
		t.Errorf("expected segment id 3, got %d", segments[0].Id)
	}
	if segments[1].Id != 2 {
		t.Errorf("expected segment id 2, got %d", segments[1].Id)
	}
	if segments[2].Id != 1 {
		t.Errorf("expected segment id 1, got %d", segments[2].Id)
	}
}

func TestGetFromSegments(t *testing.T) {
	dbDir, cleanup := testutil.EnsureTempDirectory()
	defer cleanup()

	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	key := randstr.String(10)
	value := randstr.String(10)
	err = db.Set(key, value)
	if err != nil {
		t.Fatalf("failed to set key-value pair: %v", err)
	}

	_, found := db.ongoing.Memtable.Get(key)
	if !found {
		t.Fatalf("key not found in memtable")
	}

	_, err = db.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal segment: %v", err)
	}

	_, found = db.ongoing.Memtable.Get(key)
	if found {
		t.Fatalf("key found in memtable")
	}

	val, err := db.Get(key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}
	if val != value {
		t.Fatalf("expected value %s, got %s", value, val)
	}
}
