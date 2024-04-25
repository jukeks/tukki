package memtable_test

import (
	"os"
	"testing"

	"github.com/jukeks/tukki/internal/journal"
	"github.com/jukeks/tukki/internal/memtable"
	journalv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/journal/v1"
	testutil "github.com/jukeks/tukki/tests/util"
	"github.com/thanhpk/randstr"
)

func TestJournal(t *testing.T) {
	f := testutil.CreateTempFile("test-tukki", "")
	defer f.Close()
	defer os.Remove(f.Name())

	journalWriter := journal.NewJournalWriter(f)
	err := journalWriter.Write(&journalv1.JournalEntry{
		Key:     "key",
		Value:   "value",
		Deleted: false,
	})
	if err != nil {
		t.Fatalf("failed to write journal entry: %v", err)
	}

	readerFile, err := os.Open(f.Name())
	if err != nil {
		t.Fatalf("failed to open journal file: %v", err)
	}
	defer readerFile.Close()
	journalReader := journal.NewJournalReader(readerFile)

	journalEntry := journalv1.JournalEntry{}
	err = journalReader.Read(&journalEntry)
	if err != nil {
		t.Fatalf("failed to read journal entry: %v", err)
	}

	if journalEntry.Key != "key" {
		t.Fatalf("expected key to be 'key', got '%s'", journalEntry.Key)
	}

	if journalEntry.Value != "value" {
		t.Fatalf("expected value to be 'value', got '%s'", journalEntry.Value)
	}

	if journalEntry.Deleted != false {
		t.Fatalf("expected deleted to be false, got '%v'", journalEntry.Deleted)
	}
}

func TestNewJournal(t *testing.T) {
	mt := memtable.NewMemtable()
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
	j, err := memtable.NewJournal(dbDir, mt)
	if err != nil {
		t.Fatalf("failed to create journal: %v", err)
	}

	err = j.Set("key", "value")
	if err != nil {
		t.Fatalf("failed to write journal entry: %v", err)
	}
	err = j.Set("key2", "value2")
	if err != nil {
		t.Fatalf("failed to write journal entry: %v", err)
	}
	err = j.Delete("key2")
	if err != nil {
		t.Fatalf("failed to write journal entry: %v", err)
	}

	j.Close()

	mt = memtable.NewMemtable()
	journal, err := memtable.NewJournal(dbDir, mt)
	if err != nil {
		t.Fatalf("failed to create journal reader: %v", err)
	}

	value, found := mt.Get("key")
	if !found {
		t.Fatalf("expected key to be found in memtable")
	}

	if value.Value != "value" {
		t.Fatalf("expected value to be 'value', got %v", value)
	}

	value, found = mt.Get("key2")
	if !found {
		t.Fatalf("expected key to be found in memtable")
	}
	if value.Value != "" {
		t.Fatalf("expected value to be '', got %v", value)
	}
	if value.Deleted != true {
		t.Fatalf("expected deleted to be true, got %v", value)
	}

	err = journal.Close()
	if err != nil {
		t.Fatalf("failed to close journal: %v", err)
	}
}
