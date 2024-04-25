package memtable

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jukeks/tukki/internal/journal"
	journalv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/journal/v1"
)

type Journal struct {
	journalFile *os.File
	w           *journal.JournalWriter
}

func NewJournal(dbDir string, mt Memtable) (*Journal, error) {
	journalPath := filepath.Join(dbDir, "journal")

	var journalFile *os.File
	var err error

	if _, err = os.Stat(journalPath); err == nil {
		log.Printf("journal file exists, reading journal")
		// read journal
		journalFile, err = os.Open(journalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open journal file: %w", err)
		}

		journalReader := journal.NewJournalReader(journalFile)
		err = readJournal(journalReader, mt)
		journalFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read journal: %w", err)
		}

		// open journal for appending
		log.Printf("opening journal for appending")
		journalFile, err = os.OpenFile(journalPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open journal file: %w", err)
		}
	} else {
		log.Printf("journal file does not exist, creating %s", journalPath)
		journalFile, err = os.Create(journalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create journal file: %w", err)
		}
	}

	return &Journal{
		journalFile: journalFile,
		w:           journal.NewJournalWriter(journalFile),
	}, nil
}

func (j *Journal) Set(key, value string) error {
	return j.w.Write(&journalv1.JournalEntry{
		Key:     key,
		Value:   value,
		Deleted: false,
	})
}

func (j *Journal) Delete(key string) error {
	return j.w.Write(&journalv1.JournalEntry{
		Key:     key,
		Deleted: true,
	})
}

func (j *Journal) Close() error {
	return j.journalFile.Close()
}

func readJournal(journalReader *journal.JournalReader, mt Memtable) error {
	for {
		journalEntry := &journalv1.JournalEntry{}
		err := journalReader.Read(journalEntry)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if journalEntry.Deleted {
			mt.Delete(journalEntry.Key)
		} else {
			mt.Insert(journalEntry.Key, journalEntry.Value)
		}
	}
}
