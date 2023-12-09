package db

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jukeks/tukki/journal"
	"github.com/jukeks/tukki/memtable"
)

type Database struct {
	memtable memtable.Memtable
	journal  *journal.JournalWriter
}

func NewDatabase(dbDir string) *Database {
	memtable := memtable.NewMemtable()
	journalPath := filepath.Join(dbDir, "journal")

	var journalFile *os.File
	var err error
	// if file exists
	if _, err = os.Stat(journalPath); err == nil {
		log.Printf("journal file exists, reading journal")
		// read journal
		journalFile, err = os.Open(journalPath)
		if err != nil {
			log.Fatalf("failed to open journal file: %v", err)
		}

		journalReader := journal.NewJournalReader(journalFile)
		err = readJournal(journalReader, memtable)
		journalFile.Close()
		if err != nil {
			log.Fatalf("failed to read journal: %v", err)
		}

		// open journal for appending
		log.Printf("opening journal for appending")
		journalFile, err = os.OpenFile(journalPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open journal file: %v", err)
		}
	} else {
		log.Printf("journal file does not exist, creating %s", journalPath)
		journalFile, err = os.Create(journalPath)
		if err != nil {
			log.Fatalf("failed to create journal file: %v", err)
		}
	}

	journalWriter := journal.NewJournalWriter(journalFile)
	return &Database{
		memtable: memtable,
		journal:  journalWriter,
	}
}

func (d *Database) Put(key, value string) error {
	err := d.journal.Write(&journal.JournalEntry{
		Key:     key,
		Value:   value,
		Deleted: false,
	})
	if err != nil {
		return err
	}

	d.memtable.Insert(memtable.KeyType(key), value)
	return nil
}

func (d *Database) Get(key string) (string, error) {
	value, found := d.memtable.Get(memtable.KeyType(key))
	if found {
		return value, nil
	}

	// TODO check segments
	return "", fmt.Errorf("key not found: %s", key)
}

func (d *Database) Delete(key string) error {
	err := d.journal.Write(&journal.JournalEntry{
		Key:     key,
		Deleted: true,
	})
	if err != nil {
		return err
	}

	d.memtable.Delete(memtable.KeyType(key))
	return nil
}

func (d *Database) Close() error {
	return nil
}

func readJournal(journalReader *journal.JournalReader, mt memtable.Memtable) error {
	for {
		journalEntry, err := journalReader.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if journalEntry.Deleted {
			mt.Delete(memtable.KeyType(journalEntry.Key))
		} else {
			mt.Insert(memtable.KeyType(journalEntry.Key), journalEntry.Value)
		}
	}
}
