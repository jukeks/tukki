package journal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/storage"
	journalv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/journal/v1"
)

type Journal struct {
	journalFile *os.File
	w           *JournalWriter
}

func NewJournal(dbDir string, mt memtable.Memtable) (*Journal, error) {
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

		journalReader := NewJournalReader(journalFile)
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
		w:           NewJournalWriter(journalFile),
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

func (j *JournalWriter) Write(journalEntry *journalv1.JournalEntry) error {
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

func (j *JournalReader) Read() (*journalv1.JournalEntry, error) {
	journalEntry := &journalv1.JournalEntry{}
	err := storage.ReadLengthPrefixedProtobufMessage(j.r, journalEntry)
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("failed to read journal entry: %w", err)
	}

	return journalEntry, nil
}

func readJournal(journalReader *JournalReader, mt memtable.Memtable) error {
	for {
		journalEntry, err := journalReader.Read()
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
