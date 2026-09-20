package db

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/jukeks/tukki/internal/storage/journal"
)

func TestBinaryDataLifecycle(t *testing.T) {
	for _, mode := range []journal.WriteMode{journal.WriteModeSync, journal.WriteModeAsync} {
		t.Run(map[journal.WriteMode]string{journal.WriteModeSync: "sync", journal.WriteModeAsync: "async"}[mode], func(t *testing.T) {
			dir := t.TempDir()
			config := GetDefaultConfig()
			config.JournalMode = mode
			d, err := OpenDatabaseWithConfig(dir, config)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if d != nil {
					d.Close()
				}
			})
			keys := [][]byte{{}, {0}, {0, 255}, {127}, {128}, {255}, {255, 0}}
			value := []byte{255, 0, 128, 195, 40}
			for _, key := range keys {
				inputKey, inputValue := bytes.Clone(key), bytes.Clone(value)
				if err := d.Set(inputKey, inputValue); err != nil {
					t.Fatal(err)
				}
				if len(inputKey) > 0 {
					inputKey[0] ^= 255
				}
				inputValue[0] = 1
			}
			verify := func(wantKeys [][]byte) {
				t.Helper()
				cursor, err := d.GetCursor()
				if err != nil {
					t.Fatal(err)
				}
				defer cursor.Close()
				for _, key := range wantKeys {
					got, err := d.Get(key)
					if err != nil || !bytes.Equal(got, value) {
						t.Fatalf("Get(%x) = %x, %v", key, got, err)
					}
					got[0] = 1
					pair, err := cursor.Next()
					if err != nil || !bytes.Equal(pair.Key, key) || !bytes.Equal(pair.Value, value) {
						t.Fatalf("cursor = %x:%x, %v; want %x", pair.Key, pair.Value, err, key)
					}
					pair.Value[0] = 2
					if len(pair.Key) > 0 {
						pair.Key[0] ^= 255
					}
					got, err = d.Get(key)
					if err != nil || !bytes.Equal(got, value) {
						t.Fatalf("result mutation changed Get(%x): %x, %v", key, got, err)
					}
				}
				if _, err := cursor.Next(); err != io.EOF {
					t.Fatalf("end of cursor: %v", err)
				}
			}
			reopen := func() {
				t.Helper()
				if err := d.Close(); err != nil {
					t.Fatal(err)
				}
				d = nil
				d, err = OpenDatabaseWithConfig(dir, config)
				if err != nil {
					t.Fatal(err)
				}
			}
			seal := func() {
				t.Helper()
				if err := d.ongoing.Close(); err != nil {
					t.Fatal(err)
				}
				if _, err := d.SealCurrentSegment(); err != nil {
					t.Fatal(err)
				}
			}
			verify(keys)
			reopen() // Read binary keys and values from the WAL.
			verify(keys)
			seal()
			verify(keys) // Read through the SSTable index and bloom filter.
			count, err := d.DeleteRange([]byte{0, 255}, []byte{128})
			if err != nil || count != 3 {
				t.Fatalf("DeleteRange = %d, %v", count, err)
			}
			// A queued tombstone must own its key.
			deletedKey := []byte{255, 0}
			if err := d.Delete(deletedKey); err != nil {
				t.Fatal(err)
			}
			deletedKey[0] = 1
			seal()
			if err := d.CompactSegments(1024*1024, 1, 0); err != nil {
				t.Fatal(err)
			}
			remaining := [][]byte{{}, {0}, {255}}
			verify(remaining)
			reopen()
			verify(remaining)
			for _, key := range [][]byte{{0, 255}, {127}, {128}, {255, 0}} {
				if _, err := d.Get(key); !errors.Is(err, ErrKeyNotFound) {
					t.Fatalf("deleted key %x: %v", key, err)
				}
			}
		})
	}
}

func TestBinarySnapshotAndBounds(t *testing.T) {
	dir := t.TempDir()
	d, err := OpenDatabase(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if d != nil {
			d.Close()
		}
	})
	key, value := []byte{255, 0}, []byte{128, 0, 255}
	if err := d.Set(key, value); err != nil {
		t.Fatal(err)
	}
	max := bytes.Clone(key)
	cursor, err := d.GetCursorWithRange([]byte{255}, max)
	if err != nil {
		t.Fatal(err)
	}
	max[0] = 0
	pair, err := cursor.Next()
	cursor.Close()
	if err != nil || !bytes.Equal(pair.Key, key) {
		t.Fatalf("range bound mutation: %x, %v", pair.Key, err)
	}
	snapshot, err := d.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	raw, err := snapshot.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err = UnmarshalSnapshot(raw)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	closed := d
	d = nil
	if _, err := closed.Restore(snapshot); err != nil {
		t.Fatal(err)
	}
	d, err = OpenDatabase(dir)
	if err != nil {
		t.Fatal(err)
	}
	got, err := d.Get(key)
	if err != nil || !bytes.Equal(got, value) {
		t.Fatalf("restored value = %x, %v", got, err)
	}
	if err := d.Set(nil, nil); err != nil {
		t.Fatal(err)
	}
	got, err = d.Get([]byte{})
	if err != nil || len(got) != 0 {
		t.Fatalf("empty key/value = %x, %v", got, err)
	}
}
