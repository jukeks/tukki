package journal_test

import (
	"os"
	"testing"

	"github.com/jukeks/tukki/journal"
)

func TestJournal(t *testing.T) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	journalWriter := journal.NewJournalWriter(f)

	err = journalWriter.Write(&journal.JournalEntry{
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

	journalEntry, err := journalReader.Read()
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
