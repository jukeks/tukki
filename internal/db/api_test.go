package db

import (
	"testing"

	"github.com/jukeks/tukki/internal/segments"
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
	db, err := OpenDatabase(t.TempDir())
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

func TestGetSSTableReader(t *testing.T) {
	db, err := OpenDatabase(t.TempDir())
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	key := randstr.String(10)
	value := randstr.String(10)
	err = db.Set(key, value)
	if err != nil {
		t.Fatalf("failed to set key-value pair: %v", err)
	}

	_, err = db.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal segment: %v", err)
	}

	reader, cleanup, err := db.GetSSTableReader(0)
	if err != nil {
		t.Fatalf("failed to get segment reader: %v", err)
	}
	defer cleanup()

	val, err := reader.Next()
	if err != nil {
		t.Fatalf("failed to read value: %v", err)
	}

	if val.Key != key {
		t.Fatalf("expected key %s, got %s", key, val.Key)
	}
}

func TestGetSegmentMetadata(t *testing.T) {
	db, err := OpenDatabase(t.TempDir())
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	key := randstr.String(10)
	value := randstr.String(10)
	err = db.Set(key, value)
	if err != nil {
		t.Fatalf("failed to set key-value pair: %v", err)
	}

	_, err = db.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal segment: %v", err)
	}

	segments := db.GetSegmentMetadata()
	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}

	segment, ok := segments[0]
	if !ok {
		t.Fatalf("segment not found")
	}

	if segment.Id != 0 {
		t.Fatalf("expected segment id 0, got %d", segment.Id)
	}
}
