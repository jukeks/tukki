package journal

import (
	"testing"

	"github.com/jukeks/tukki/internal/storage/files"
	walv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/wal/v1"
	"github.com/thanhpk/randstr"
)

func BenchmarkSyncWriter(b *testing.B) {
	dbDir := b.TempDir()
	filename := files.Filename("journal")
	j, err := OpenJournal(dbDir, filename, WriteModeSync, func(r *JournalReader) error {
		return nil
	})
	if err != nil {
		b.Fatalf("failed to open wal: %v", err)
	}

	entry := &walv1.WalEntry{
		Key:   "key",
		Value: randstr.String(16 * 1024),
	}

	b.SetBytes(int64(len(entry.Key) + len(entry.Value)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := j.Writer.Write(entry); err != nil {
			b.Fatalf("failed to set key: %v", err)
		}
	}

	err = j.Close()
	if err != nil {
		b.Fatalf("failed to close journal: %v", err)
	}
}

func BenchmarkAsyncWriter(b *testing.B) {
	dbDir := b.TempDir()
	filename := files.Filename("journal")
	j, err := OpenJournal(dbDir, filename, WriteModeAsync, func(r *JournalReader) error {
		return nil
	})
	if err != nil {
		b.Fatalf("failed to open wal: %v", err)
	}

	entry := &walv1.WalEntry{
		Key:   "key",
		Value: randstr.String(16 * 1024),
	}

	b.SetBytes(int64(len(entry.Key) + len(entry.Value)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := j.Writer.Write(entry); err != nil {
			b.Fatalf("failed to set key: %v", err)
		}
	}
	err = j.Close()
	if err != nil {
		b.Fatalf("failed to close journal: %v", err)
	}
}
