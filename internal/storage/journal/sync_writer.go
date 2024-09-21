package journal

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/jukeks/tukki/internal/storage"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type SynchronousJournalWriter struct {
	w           WriteSyncer
	b           *bufio.Writer
	journalCopy *bytes.Buffer
}

func NewSynchronousJournalWriter(w WriteSyncer, head []byte) *SynchronousJournalWriter {
	return &SynchronousJournalWriter{
		w:           w,
		b:           bufio.NewWriter(w),
		journalCopy: bytes.NewBuffer(head),
	}
}

func (j *SynchronousJournalWriter) Write(journalEntry protoreflect.ProtoMessage) error {
	_, err := storage.WriteLengthPrefixedProtobufMessage(j.journalCopy, journalEntry)
	if err != nil {
		return fmt.Errorf("failed to write journal entry to copy: %w", err)
	}

	_, err = storage.WriteLengthPrefixedProtobufMessage(j.b, journalEntry)
	if err != nil {
		return fmt.Errorf("failed to write journal entry: %w", err)
	}

	err = j.b.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	err = j.w.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	return nil
}

func (j *SynchronousJournalWriter) Close() error {
	return nil
}

func (j *SynchronousJournalWriter) Snapshot() []byte {
	b := j.journalCopy.Bytes()
	buff := make([]byte, len(b))
	copy(buff, b)
	return buff
}
