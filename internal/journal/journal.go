package journal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

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

type Journal struct {
	File   *os.File
	Writer *JournalWriter
}

type ExistingJournalHandler func(r *JournalReader) error

func OpenJournal(dbDir string, journalName string, existingHandler ExistingJournalHandler) (*Journal, error) {
	journalPath := storage.GetPath(dbDir, journalName)

	var journalFile *os.File
	var err error

	if _, err = os.Stat(journalPath); err == nil {
		log.Printf("journal file exists %s, reading journal", journalPath)
		// read journal
		journalFile, err = os.Open(journalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open journal file %s: %w", journalPath, err)
		}

		journalReader := NewJournalReader(journalFile)
		err = existingHandler(journalReader)
		journalFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read journal file %s: %w", journalPath, err)
		}

		// open journal for appending
		log.Printf("opening journal for appending")
		journalFile, err = os.OpenFile(journalPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open journal file %s: %w", journalPath, err)
		}
	} else {
		log.Printf("journal file does not exist, creating %s", journalPath)
		journalFile, err = os.Create(journalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create journal file: %w", err)
		}
	}

	return &Journal{
		File:   journalFile,
		Writer: NewJournalWriter(journalFile),
	}, nil
}

func (j *Journal) Close() error {
	return j.File.Close()
}
