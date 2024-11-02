package db

import (
	"testing"

	"github.com/jukeks/tukki/internal/storage/segments"
)

func writePairsAndSeal(t *testing.T, db *Database, pairs ...Pair) {
	for _, pair := range pairs {
		if err := db.Set(pair.Key, pair.Value); err != nil {
			t.Fatalf("failed to set key: %v", err)
		}
	}

	if _, err := db.SealCurrentSegment(); err != nil {
		t.Fatalf("failed to seal database: %v", err)
	}
}

func verifyPairs(t *testing.T, db *Database, pairs ...Pair) {
	for _, pair := range pairs {
		value, err := db.Get(pair.Key)
		if err != nil {
			t.Fatalf("failed to get key %s: %v", pair.Key, err)
		}

		if value != pair.Value {
			t.Fatalf("expected value %s, got %s", pair.Value, value)
		}
	}
}

func TestCompactSegments(t *testing.T) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}
	defer db.Close()

	writePairsAndSeal(t, db, Pair{Key: "a", Value: "1"}, Pair{Key: "b", Value: "2"})
	if len(db.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", db.segments)
	}

	writePairsAndSeal(t, db, Pair{Key: "c", Value: "3"}, Pair{Key: "d", Value: "4"})
	if len(db.segments) != 2 {
		t.Fatalf("expected segments map to have 2 elements, got %v", db.segments)
	}

	writePairsAndSeal(t, db, Pair{Key: "e", Value: "5"}, Pair{Key: "f", Value: "6"})
	if len(db.segments) != 3 {
		t.Fatalf("expected segments map to have 3 elements, got %v", db.segments)
	}

	writePairsAndSeal(t, db, Pair{Key: "b", Value: "7"})

	segmentMetadata := db.getSegmentsSorted()

	segmentIds := make([]segments.SegmentId, len(segmentMetadata))
	for i, segment := range segmentMetadata {
		segmentIds[i] = segment.Id
	}

	if err := db.CompactSegments(160*1024*1024, segmentIds...); err != nil {
		t.Fatalf("failed to compact segments: %v", err)
	}

	if len(db.segments) != 1 {
		t.Fatalf("expected segments map to have 1 element, got %v", db.segments)
	}

	compactedSegment := db.segments[segmentIds[0]]
	if compactedSegment.Id != 0 {
		t.Fatalf("expected compacted segment id to be 0, got %d", compactedSegment.Id)
	}

	verifyPairs(t, db,
		Pair{Key: "a", Value: "1"},
		Pair{Key: "b", Value: "7"},
		Pair{Key: "c", Value: "3"},
		Pair{Key: "d", Value: "4"},
		Pair{Key: "e", Value: "5"},
		Pair{Key: "f", Value: "6"})
}
