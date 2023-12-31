package db

import (
	"log"

	"github.com/jukeks/tukki/internal/journal"
	"github.com/jukeks/tukki/internal/memtable"
)

type Database struct {
	memtable memtable.Memtable
	journal  *journal.Journal
}

func NewDatabase(dbDir string) *Database {
	memtable := memtable.NewMemtable()

	journal, err := journal.NewJournal(dbDir, memtable)
	if err != nil {
		log.Fatalf("failed to create journal: %v", err)
	}

	return &Database{
		memtable: memtable,
		journal:  journal,
	}
}

func (d *Database) Set(key, value string) error {
	err := d.journal.Set(key, value)
	if err != nil {
		return err
	}

	d.memtable.Insert(key, value)
	return nil
}

func (d *Database) Get(key string) (string, error) {
	value, found := d.memtable.Get(key)
	if found {
		if value.Deleted {
			return "", ErrKeyNotFound
		}
		return value.Value, nil
	}

	// TODO check segments
	return "", ErrKeyNotFound
}

func (d *Database) Delete(key string) error {
	err := d.journal.Delete(key)
	if err != nil {
		return err
	}

	d.memtable.Delete(key)
	return nil
}

func (d *Database) Close() error {
	return d.journal.Close()
}
