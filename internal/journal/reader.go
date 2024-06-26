package journal

import (
	"fmt"
	"io"

	"github.com/jukeks/tukki/internal/storage"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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
