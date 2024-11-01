package db

import (
	"io"
	"testing"
)

type Pair struct {
	Key   string
	Value string
}

func setupDB(t *testing.T) (*Database, []Pair) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	// segment 1
	pairs := []Pair{
		{"key5", "value5"},
		{"key2", "value2"},
		{"key4", "value4"},
	}
	for _, pair := range pairs {
		if err := db.Set(pair.Key, pair.Value); err != nil {
			t.Fatalf("failed to set key: %v", err)
		}
	}
	_, err = db.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	// segment 2
	pairs = []Pair{
		{"key1", "value1"},
		{"key3", "value3"},
	}
	for _, pair := range pairs {
		if err := db.Set(pair.Key, pair.Value); err != nil {
			t.Fatalf("failed to set key: %v", err)
		}
	}

	_, err = db.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	// memtable
	pairs = []Pair{
		{"key6", "value6"},
		{"key7", "value7"},
	}
	for _, pair := range pairs {
		if err := db.Set(pair.Key, pair.Value); err != nil {
			t.Fatalf("failed to set key: %v", err)
		}
	}

	expected := []Pair{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
		{"key4", "value4"},
		{"key5", "value5"},
		{"key6", "value6"},
		{"key7", "value7"},
	}

	return db, expected
}

func TestIteratorWithoutRange(t *testing.T) {
	db, expected := setupDB(t)
	defer db.Close()

	iterator, err := db.GetIterator("", "")
	if err != nil {
		t.Fatalf("failed to get iterator: %v", err)
	}
	defer iterator.Close()

	for _, pair := range expected {
		entry, err := iterator.Next()
		t.Logf("entry: %v", entry)
		if err != nil {
			t.Fatalf("failed to get key: %v", err)
		}
		if entry.Key != pair.Key {
			t.Fatalf("expected key %s, got %s", pair.Key, entry.Key)
		}
		if entry.Value != pair.Value {
			t.Fatalf("expected value %s, got %s", pair.Value, entry.Value)
		}
	}

	_, err = iterator.Next()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestIteratorWithRange(t *testing.T) {
	db, expected := setupDB(t)
	defer db.Close()

	iterator, err := db.GetIterator("key2", "key5")
	if err != nil {
		t.Fatalf("failed to get iterator: %v", err)
	}
	defer iterator.Close()

	expected = expected[1:4]
	for _, pair := range expected {
		entry, err := iterator.Next()
		t.Logf("entry: %v", entry)
		if err != nil {
			t.Fatalf("failed to get key: %v", err)
		}
		if entry.Key != pair.Key {
			t.Fatalf("expected key %s, got %s", pair.Key, entry.Key)
		}
		if entry.Value != pair.Value {
			t.Fatalf("expected value %s, got %s", pair.Value, entry.Value)
		}
	}

	_, err = iterator.Next()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestIteratorWithRangeAndEmptyStart(t *testing.T) {
	db, expected := setupDB(t)
	defer db.Close()

	iterator, err := db.GetIterator("", "key5")
	if err != nil {
		t.Fatalf("failed to get iterator: %v", err)
	}
	defer iterator.Close()

	expected = expected[:4]
	for _, pair := range expected {
		entry, err := iterator.Next()
		t.Logf("entry: %v", entry)
		if err != nil {
			t.Fatalf("failed to get key: %v", err)
		}
		if entry.Key != pair.Key {
			t.Fatalf("expected key %s, got %s", pair.Key, entry.Key)
		}
		if entry.Value != pair.Value {
			t.Fatalf("expected value %s, got %s", pair.Value, entry.Value)
		}
	}

	_, err = iterator.Next()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestIteratorWithRangeAndEmptyEnd(t *testing.T) {
	db, expected := setupDB(t)
	defer db.Close()

	iterator, err := db.GetIterator("key2", "")
	if err != nil {
		t.Fatalf("failed to get iterator: %v", err)
	}

	expected = expected[1:]
	for _, pair := range expected {
		entry, err := iterator.Next()
		t.Logf("entry: %v", entry)
		if err != nil {
			t.Fatalf("failed to get key: %v", err)
		}
		if entry.Key != pair.Key {
			t.Fatalf("expected key %s, got %s", pair.Key, entry.Key)
		}
		if entry.Value != pair.Value {
			t.Fatalf("expected value %s, got %s", pair.Value, entry.Value)
		}
	}
}
