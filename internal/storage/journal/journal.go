package journal

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jukeks/tukki/internal/storage/files"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type WriteMode int

const (
	WriteModeSync WriteMode = iota
	WriteModeAsync
	WriteModeInMemory
)

type WriteSyncer interface {
	io.Writer
	Sync() error
}

type JournalWriter interface {
	Write(journalEntry protoreflect.ProtoMessage) error
	Close() error
	Snapshot() []byte
}

type Journal struct {
	File   *os.File
	Writer JournalWriter
}

type ExistingJournalHandler func(r *JournalReader) error

func OpenJournal(dbDir string, journalName files.Filename, writemode WriteMode,
	existingHandler ExistingJournalHandler) (*Journal, error) {

	journalPath := files.GetPath(dbDir, journalName)

	var journalFile *os.File
	var journalCopy []byte
	var err error

	if _, err = os.Stat(journalPath); err == nil {
		log.Printf("journal file exists %s, reading journal", journalPath)
		// read journal
		journalFile, err = os.Open(journalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open journal file %s: %w", journalPath, err)
		}
		journalCopy, err = os.ReadFile(journalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read journal file %s: %w", journalPath, err)
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
		Writer: NewJournalWriter(journalFile, writemode, journalCopy),
	}, nil
}

func (j *Journal) Close() error {
	j.Writer.Close()
	return j.File.Close()
}

func NewJournalWriter(w WriteSyncer, writeMode WriteMode, head []byte) JournalWriter {
	if writeMode == WriteModeSync {
		return NewSynchronousJournalWriter(w, head)
	}
	if writeMode == WriteModeInMemory {
		return NewInMemJournalWriter(head)
	}

	return NewAsynchronousJournalWriter(w, head)
}
