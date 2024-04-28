package segments

import (
	"testing"

	testutil "github.com/jukeks/tukki/tests/util"
	"github.com/thanhpk/randstr"
)

func TestSegmentManagerGetSegmentsSorted(t *testing.T) {
	sm := SegmentManager{
		segments: map[SegmentId]Segment{
			1: {
				Id:       1,
				Filename: "segment1",
			},
			3: {
				Id:       3,
				Filename: "segment3",
			},
			2: {
				Id:       2,
				Filename: "segment2",
			},
		},
	}

	segments := sm.GetSegmentsSorted()
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

func TestSegmentManagerGet(t *testing.T) {
	dbDir := testutil.EnsureTempDirectory("test-tukki-segments-" + randstr.String(10))
	sm, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	ongoing := sm.GetOnGoingSegment()

	key := randstr.String(10)
	value := randstr.String(10)
	writeToWalAndMemtable(t, ongoing, key, value)

	_, err = sm.Get(key)
	if err != ErrKeyNotFound {
		t.Fatalf("expected key not found error, got %v", err)
	}

	_, err = sm.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal segment: %v", err)
	}

	val, err := sm.Get(key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}
	if val != value {
		t.Fatalf("expected value %s, got %s", value, val)
	}
}
