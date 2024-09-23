package journal

import (
	"bytes"
	"fmt"

	"github.com/jukeks/tukki/internal/storage/marshalling"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type InMemJournalWriter struct {
	journalCopy *bytes.Buffer
}

func NewInMemJournalWriter(head []byte) *InMemJournalWriter {
	buff := make([]byte, 0, 2*1024*1024)
	copy(buff, head)
	return &InMemJournalWriter{
		journalCopy: bytes.NewBuffer(buff),
	}
}

func (j *InMemJournalWriter) Write(journalEntry protoreflect.ProtoMessage) error {
	_, err := marshalling.WriteLengthPrefixedProtobufMessage(j.journalCopy, journalEntry)
	if err != nil {
		return fmt.Errorf("failed to write journal entry to copy: %w", err)
	}

	return nil
}

func (j *InMemJournalWriter) Close() error {
	return nil
}

func (j *InMemJournalWriter) Snapshot() []byte {
	b := j.journalCopy.Bytes()
	buff := make([]byte, len(b))
	copy(buff, b)
	return buff
}
