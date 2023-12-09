package db

import (
	"fmt"
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
	err := d.journal.Delete(key)
	if err != nil {
		return err
	}

	d.memtable.Delete(memtable.KeyType(key))
	return nil
}

func (d *Database) Close() error {
	return d.journal.Close()
}
