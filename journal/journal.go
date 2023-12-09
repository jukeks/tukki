package journal

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jukeks/tukki/lib"
)

type WriteSyncer interface {
	io.Writer
	Sync() error
}

type JournalWriter struct {
	w WriteSyncer
	b *bufio.Writer
}

func NewJournalWriter(w WriteSyncer) *JournalWriter {
	return &JournalWriter{w: w, b: bufio.NewWriter(w)}
}

func (j *JournalWriter) Write(journalEntry *JournalEntry) error {
	err := lib.WriteLengthPrefixedProtobufMessage(j.b, journalEntry)
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

type JournalReader struct {
	r io.Reader
}

func NewJournalReader(r io.Reader) *JournalReader {
	return &JournalReader{r: r}
}

func (j *JournalReader) Read() (*JournalEntry, error) {
	journalEntry := &JournalEntry{}
	err := lib.ReadLengthPrefixedProtobufMessage(j.r, journalEntry)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("failed to read journal entry: %w", err)
	}

	return journalEntry, nil
}
