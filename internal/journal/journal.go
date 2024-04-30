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

type WriteMode int

const (
	WriteModeSync WriteMode = iota
	WriteModeAsync
)

type WriteSyncer interface {
	io.Writer
	Sync() error
}

type JournalWriter interface {
	Write(journalEntry protoreflect.ProtoMessage) error
}

type Journal struct {
	File   *os.File
	Writer JournalWriter
}

type ExistingJournalHandler func(r *JournalReader) error

func OpenJournal(dbDir string, journalName storage.Filename, writemode WriteMode,
	existingHandler ExistingJournalHandler) (*Journal, error) {

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
		Writer: NewJournalWriter(journalFile, writemode),
	}, nil
}

func (j *Journal) Close() error {
	return j.File.Close()
}

func NewJournalWriter(w WriteSyncer, writeMode WriteMode) JournalWriter {
	return &SynchronousJournalWriter{w: w, b: bufio.NewWriter(w)}
}
