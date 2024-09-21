package db

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/jukeks/tukki/internal/db/segments"
	"github.com/jukeks/tukki/internal/storage"
	"github.com/jukeks/tukki/internal/storage/sstable"
)

var ErrKeyNotFound = errors.New("key not found in segments")

func (db *Database) Get(key string) (string, error) {
	value, err := db.ongoing.Get(key)
	if err == nil {
		return value, nil
	}

	return db.getFromSegments(key)
}

func (db *Database) getFromSegments(key string) (string, error) {
	for _, segment := range db.getSegmentsSorted() {
		contains := db.members[segment.Id].Contains(key)
		if !contains {
			// this looks unnecessary right now, but eventually all segment
			// indexes might not be in memory, so it's beneficial to check
			// if the key is in the segment before reading the index
			continue
		}

		offset, found := db.indexes[segment.Id].Entries[key]
		if !found {
			// false positive, key is not in segment
			continue
		}

		segmentFile, err := storage.OpenFile(db.dbDir, segment.SegmentFile)
		if err != nil {
			return "", err
		}
		defer segmentFile.Close()

		reader := sstable.NewSSTableSeeker(segmentFile)
		entry, err := reader.ReadAt(offset)
		if err != nil {
			return "", err
		}

		if entry.Key != key {
			return "", fmt.Errorf("expected key %s, got %s", key, entry.Key)
		}

		if entry.Deleted {
			return "", ErrKeyNotFound
		}

		return entry.Value, nil
	}

	return "", ErrKeyNotFound
}

func (db *Database) getSegmentsSorted() []segments.SegmentMetadata {
	db.mu.Lock()
	defer db.mu.Unlock()

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

func (db *Database) GetSegmentMetadata() map[segments.SegmentId]segments.SegmentMetadata {
	db.mu.Lock()
	defer db.mu.Unlock()

	segments := make(map[segments.SegmentId]segments.SegmentMetadata)
	for k, v := range db.segments {
		segments[k] = v
	}

	return segments
}

type Cleanup func()

func (db *Database) GetSSTableReader(segmentId segments.SegmentId) (*sstable.SSTableReader, Cleanup, error) {
	db.mu.Lock()
	segmentMetadata, ok := db.segments[segmentId]
	db.mu.Unlock()

	if !ok {
		return nil, nil, fmt.Errorf("segment not found: %d", segmentId)
	}

	f, err := storage.OpenFile(db.dbDir, segmentMetadata.SegmentFile)
	if err != nil {
		return nil, nil, err
	}

	return sstable.NewSSTableReader(f), func() { f.Close() }, nil
}

func (db *Database) Set(key, value string) error {
	if err := db.handleWalSizeLimit(); err != nil {
		return err
	}

	return db.ongoing.Set(key, value)
}

func (db *Database) Delete(key string) error {
	if err := db.handleWalSizeLimit(); err != nil {
		return err
	}

	return db.ongoing.Delete(key)
}

func (db *Database) getCurrentWalSize() uint64 {
	return db.ongoing.Memtable.Size()
}

func (db *Database) isOverWalSizeLimit() bool {
	return db.getCurrentWalSize() > db.walSizeLimit
}

func (db *Database) handleWalSizeLimit() error {
	if !db.isOverWalSizeLimit() {
		return nil
	}

	log.Printf("wal size limit reached, sealing current segment")
	err := db.ongoing.Close()
	if err != nil {
		return err
	}
	_, err = db.SealCurrentSegment()
	if err != nil {
		return err
	}
	return nil
}
