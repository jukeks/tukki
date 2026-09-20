package db

import (
	"bytes"
	"io"
	"testing"
)

func setupDB(t *testing.T) (*Database, []Pair) {
	dbDir := t.TempDir()
	db, err := OpenDatabase(dbDir)
	if err != nil {
		t.Fatalf("failed to open segment manager: %v", err)
	}

	// segment 1
	pairs := []Pair{
		{[]byte("key1"), []byte("value1")},
		{[]byte("key5"), []byte("value5")},
		{[]byte("key2"), []byte("value2")},
		{[]byte("key4"), []byte("value4")},
		{[]byte("key9"), []byte("value9")},
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
		{[]byte("key1"), []byte("value1-new")},
		{[]byte("key3"), []byte("value3")},
		{[]byte("key5"), []byte("value5-new")},
	}
	for _, pair := range pairs {
		if err := db.Set(pair.Key, pair.Value); err != nil {
			t.Fatalf("failed to set key: %v", err)
		}
	}
	if err := db.Delete([]byte("key4")); err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}

	_, err = db.SealCurrentSegment()
	if err != nil {
		t.Fatalf("failed to seal current segment: %v", err)
	}

	// memtable
	pairs = []Pair{
		{[]byte("key6"), []byte("value6")},
		{[]byte("key7"), []byte("value7")},
		{[]byte("key1"), []byte("value1-new-new")},
	}
	for _, pair := range pairs {
		if err := db.Set(pair.Key, pair.Value); err != nil {
			t.Fatalf("failed to set key: %v", err)
		}
	}

	expected := []Pair{
		{[]byte("key1"), []byte("value1-new-new")},
		{[]byte("key2"), []byte("value2")},
		{[]byte("key3"), []byte("value3")},
		{[]byte("key5"), []byte("value5-new")},
		{[]byte("key6"), []byte("value6")},
		{[]byte("key7"), []byte("value7")},
		{[]byte("key9"), []byte("value9")},
	}

	return db, expected
}

func TestCursorWithoutRange(t *testing.T) {
	db, expected := setupDB(t)
	defer db.Close()

	iterator, err := db.GetCursor()
	if err != nil {
		t.Fatalf("failed to get iterator: %v", err)
	}
	defer iterator.Close()

	for _, pair := range expected {
		entry, err := iterator.Next()
		if err != nil {
			t.Fatalf("failed to get key: %v", err)
		}
		t.Logf("entry: %v", entry)
		if !bytes.Equal(entry.Key, pair.Key) {
			t.Fatalf("expected key %s, got %s", pair.Key, entry.Key)
		}
		if !bytes.Equal(entry.Value, pair.Value) {
			t.Fatalf("expected value %s, got %s", pair.Value, entry.Value)
		}
	}

	_, err = iterator.Next()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestCursorWithRange(t *testing.T) {
	db, expected := setupDB(t)
	defer db.Close()

	iterator, err := db.GetCursorWithRange([]byte("key2"), []byte("key5"))
	if err != nil {
		t.Fatalf("failed to get iterator: %v", err)
	}
	defer iterator.Close()

	expected = expected[1:4]
	for _, pair := range expected {
		entry, err := iterator.Next()
		if err != nil {
			t.Fatalf("failed to get key: %v", err)
		}
		t.Logf("entry: %v", entry)
		if !bytes.Equal(entry.Key, pair.Key) {
			t.Fatalf("expected key %s, got %s", pair.Key, entry.Key)
		}
		if !bytes.Equal(entry.Value, pair.Value) {
			t.Fatalf("expected value %s, got %s", pair.Value, entry.Value)
		}
	}

	_, err = iterator.Next()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestCursorWithRangeAndEmptyStart(t *testing.T) {
	db, expected := setupDB(t)
	defer db.Close()

	iterator, err := db.GetCursorWithRange([]byte(""), []byte("key5"))
	if err != nil {
		t.Fatalf("failed to get iterator: %v", err)
	}
	defer iterator.Close()

	expected = expected[:4]
	for _, pair := range expected {
		entry, err := iterator.Next()
		if err != nil {
			t.Fatalf("failed to get key: %v", err)
		}
		t.Logf("entry: %v", entry)
		if !bytes.Equal(entry.Key, pair.Key) {
			t.Fatalf("expected key %s, got %s", pair.Key, entry.Key)
		}
		if !bytes.Equal(entry.Value, pair.Value) {
			t.Fatalf("expected value %s, got %s", pair.Value, entry.Value)
		}
	}

	_, err = iterator.Next()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestCursorWithRangeAndEmptyEnd(t *testing.T) {
	db, expected := setupDB(t)
	defer db.Close()

	iterator, err := db.GetCursorWithRange([]byte("key2"), []byte(""))
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
		if !bytes.Equal(entry.Key, pair.Key) {
			t.Fatalf("expected key %s, got %s", pair.Key, entry.Key)
		}
		if !bytes.Equal(entry.Value, pair.Value) {
			t.Fatalf("expected value %s, got %s", pair.Value, entry.Value)
		}
	}
}
