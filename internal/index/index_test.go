package index

import (
	"os"
	"testing"

	"github.com/jukeks/tukki/testutil"
)

func TestIndex(t *testing.T) {
	tmpDir := t.TempDir()
	f := testutil.CreateTempFile(tmpDir, "sstable-test-*")

	entries := make(map[string]int64)
	entries["key1"] = 0
	entries["key2"] = 1
	entries["key3"] = 2

	w := NewIndexWriter(f)
	for key, offset := range entries {
		if err := w.Add(key, offset); err != nil {
			t.Fatal(err)
		}
	}

	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	f, err = os.Open(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	index, err := OpenIndex(f)
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
