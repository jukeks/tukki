package journal

import (
	"io"
	"os"
	"testing"

	"github.com/thanhpk/randstr"

	"github.com/jukeks/tukki/internal/storage/files"
	testutil "github.com/jukeks/tukki/testutil"

	walv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/wal/v1"
)

func TestJournalWriter(t *testing.T) {
	f := testutil.CreateTempFile(t.TempDir(), "")
	defer f.Close()

	journalWriter := NewJournalWriter(f, WriteModeSync, nil)
	err := journalWriter.Write(&walv1.WalEntry{
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
	journalReader := NewJournalReader(readerFile)

	journalEntry := walv1.WalEntry{}
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

func TestOpenJournal(t *testing.T) {
	dbDir := t.TempDir()
	filename := files.Filename(randstr.String(10))

	j, err := OpenJournal(dbDir, filename, WriteModeSync, func(r *JournalReader) error {
		return nil
	})
	if err != nil {
		t.Fatalf("failed to create journal: %v", err)
	}
	if j == nil {
		t.Fatalf("journal is nil")
	}

	err = j.Writer.Write(&walv1.WalEntry{
		Key:   "key1",
		Value: "value1",
	})
	if err != nil {
		t.Fatalf("failed to write journal entry: %v", err)
	}
	err = j.Writer.Write(&walv1.WalEntry{
		Key:   "key2",
		Value: "value2",
	})
	if err != nil {
		t.Fatalf("failed to write journal entry: %v", err)
	}

	err = j.Close()
	if err != nil {
		t.Fatalf("failed to close journal: %v", err)
	}

	var entries []*walv1.WalEntry
	j, err = OpenJournal(dbDir, filename, WriteModeSync, func(r *JournalReader) error {
		for {
			journalEntry := &walv1.WalEntry{}
			err := r.Read(journalEntry)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			entries = append(entries, journalEntry)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to open journal: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "key1" {
		t.Fatalf("expected key1, got %s", entries[0].Key)
	}
	if entries[0].Value != "value1" {
		t.Fatalf("expected value1, got %s", entries[0].Value)
	}
	if entries[1].Key != "key2" {
		t.Fatalf("expected key2, got %s", entries[1].Key)
	}
	if entries[1].Value != "value2" {
		t.Fatalf("expected value2, got %s", entries[1].Value)
	}

	j.Close()
}

func TestJournalWriterSnapshotSync(t *testing.T) {
	testJournalWriterSnapshot(t, WriteModeSync)
}

func TestJournalWriterSnapshotAsync(t *testing.T) {
	testJournalWriterSnapshot(t, WriteModeAsync)
}

func testJournalWriterSnapshot(t *testing.T, writeMode WriteMode) {
	dbDir := t.TempDir()
	filename := files.Filename(randstr.String(10))

	j, err := OpenJournal(dbDir, filename, writeMode, func(r *JournalReader) error {
		return nil
	})
	if err != nil {
		t.Fatalf("failed to create journal: %v", err)
	}
	if j == nil {
		t.Fatalf("journal is nil")
	}

	err = j.Writer.Write(&walv1.WalEntry{
		Key:   "key1",
		Value: "value1",
	})
	if err != nil {
		t.Fatalf("failed to write journal entry: %v", err)
	}

	snapshot := j.Writer.Snapshot()
	if len(snapshot) == 0 {
		t.Fatalf("snapshot is empty")
	}

	j.Close()

	filename2 := files.Filename(randstr.String(10))
	f, err := files.CreateFile(dbDir, filename2)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	f.Close()

	if err := os.WriteFile(f.Name(), snapshot, 0644); err != nil {
		t.Fatalf("failed to write snapshot to file: %v", err)
	}

	var entries []*walv1.WalEntry
	j2, err := OpenJournal(dbDir, filename, writeMode, func(r *JournalReader) error {
		for {
			journalEntry := &walv1.WalEntry{}
			err := r.Read(journalEntry)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			entries = append(entries, journalEntry)
		}
		return nil
	})
	j2.Close()

	if entries[0].Key != "key1" {
		t.Fatalf("expected key1, got %s", entries[0].Key)
	}
}
