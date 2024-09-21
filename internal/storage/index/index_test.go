package index

import (
	"testing"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/thanhpk/randstr"
)

func TestIndex(t *testing.T) {
	dbDir := t.TempDir()
	filename := files.Filename(randstr.String(10))

	entries := make(OffsetMap)
	entries["key1"] = 0
	entries["key2"] = 1
	entries["key3"] = 2

	f, err := files.CreateFile(dbDir, filename)
	if err != nil {
		t.Fatal(err)
	}

	w := NewIndexWriter(f)
	err = w.WriteFromOffsets(entries)
	if err != nil {
		t.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	index, err := OpenIndex(dbDir, filename)
	if err != nil {
		t.Fatal(err)
	}

	if len(index.Entries) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(index.Entries))
	}
	for key, offset := range entries {
		if index.Entries[key] != offset {
			t.Fatalf("expected offset %d for key %s, got %d", offset, key, index.Entries[key])
		}
	}
}
