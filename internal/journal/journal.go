package journal

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jukeks/tukki/internal/storage"
	"google.golang.org/protobuf/reflect/protoreflect"
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

func (j *JournalWriter) Write(journalEntry protoreflect.ProtoMessage) error {
	err := storage.WriteLengthPrefixedProtobufMessage(j.b, journalEntry)
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

func (j *JournalReader) Read(journalEntry protoreflect.ProtoMessage) error {
	err := storage.ReadLengthPrefixedProtobufMessage(j.r, journalEntry)
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("failed to read journal entry: %w", err)
	}

	return nil
}
