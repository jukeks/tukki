package db

import (
	"errors"
	"os"
	"sort"

	"github.com/jukeks/tukki/internal/segments"
	"github.com/jukeks/tukki/internal/sstable"
	"github.com/jukeks/tukki/internal/storage"
)

var ErrKeyNotFound = errors.New("key not found in segments")

func (db *Database) Get(key string) (string, error) {
	value, found := db.ongoing.Memtable.Get(key)
	if found {
		if value.Deleted {
			return "", ErrKeyNotFound
		}
		return value.Value, nil
	}

	return db.getFromSegments(key)
}

func (db *Database) getFromSegments(key string) (string, error) {
	for _, segment := range db.getSegmentsSorted() {
		segmentPath := storage.GetPath(db.dbDir, segment.Filename)
		segmentFile, err := os.Open(segmentPath)
		if err != nil {
			return "", err
		}
		defer segmentFile.Close()

		reader := sstable.NewSSTableReader(segmentFile)
		for entry, err := reader.Next(); err == nil; entry, err = reader.Next() {
			if entry.Key == key {
				if entry.Deleted {
					return "", ErrKeyNotFound
				}
				return entry.Value, nil
			}
		}
	}

	return "", ErrKeyNotFound
}

func (db *Database) getSegmentsSorted() []segments.SegmentMetadata {
	keys := make([]segments.SegmentId, len(db.segments))
	i := 0
	for k := range db.segments {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	segments := make([]segments.SegmentMetadata, len(keys))
	for i, k := range keys {
		segments[i] = db.segments[k]
	}

	return segments
}

func (db *Database) Set(key, value string) error {
	return db.ongoing.Set(key, value)
}

func (db *Database) Delete(key string) error {
	return db.ongoing.Delete(key)
}
