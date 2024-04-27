package db

import (
	"log"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/segments"
)

type Database struct {
	memtable       memtable.Memtable
	wal            *memtable.MembtableJournal
	segmentManager *segments.SegmentManager
}

func OpenDatabase(dbDir string) *Database {
	segmentsManager, err := segments.OpenDatabase(dbDir)
	if err != nil {
		log.Fatalf("failed to open segments manager: %v", err)
	}

	ongoing := segmentsManager.GetOnGoingSegment()

	mt := memtable.NewMemtable()
	wal, err := memtable.OpenWal(dbDir, ongoing.WalFilename, mt)
	if err != nil {
		log.Fatalf("failed to create journal: %v", err)
	}

	return &Database{
		memtable:       mt,
		wal:            wal,
		segmentManager: segmentsManager,
	}
}

func (d *Database) Set(key, value string) error {
	err := d.wal.Set(key, value)
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
	err := d.wal.Delete(key)
	if err != nil {
		return err
	}

	d.memtable.Delete(key)
	return nil
}

func (d *Database) Close() error {
	err := d.wal.Close()
	if err != nil {
		return err
	}

	return d.segmentManager.Close()
}
