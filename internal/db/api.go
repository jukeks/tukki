package db

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/index"
	"github.com/jukeks/tukki/internal/storage/memtable"
	"github.com/jukeks/tukki/internal/storage/segments"
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

		segmentFile, err := files.OpenFile(db.dbDir, segment.SegmentFile)
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

	return db.getSegmentsSortedUnlocked()
}

func (db *Database) getSegmentsSortedUnlocked() []segments.SegmentMetadata {
	ids := make([]segments.SegmentId, len(db.segments))
	i := 0
	for segmentId := range db.segments {
		ids[i] = segmentId
		i++
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] > ids[j]
	})

	segments := make([]segments.SegmentMetadata, len(ids))
	for i, segmentId := range ids {
		segments[i] = db.segments[segmentId]
	}

	return segments
}

func (db *Database) getIndexesCopyUnlocked() map[segments.SegmentId]*index.Index {
	copy := make(map[segments.SegmentId]*index.Index)
	for k, v := range db.indexes {
		copy[k] = v
	}
	return copy
}

func (db *Database) getStateCopy() (
	memtable.Memtable, []segments.SegmentMetadata, map[segments.SegmentId]*index.Index,
) {
	db.mu.Lock()
	defer db.mu.Unlock()

	mt := db.ongoing.Memtable.Copy()
	segments := db.getSegmentsSortedUnlocked()
	index := db.getIndexesCopyUnlocked()

	return mt, segments, index

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

func (db *Database) GetSSTableReader(segmentId segments.SegmentId) (
	*sstable.SSTableReader, Cleanup, error,
) {
	db.mu.Lock()
	segmentMetadata, ok := db.segments[segmentId]
	db.mu.Unlock()

	if !ok {
		return nil, nil, fmt.Errorf("segment not found: %d", segmentId)
	}

	f, err := files.OpenFile(db.dbDir, segmentMetadata.SegmentFile)
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
	return db.getCurrentWalSize() > db.config.WalSizeLimit
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

func (db *Database) GetCursor() (*Cursor, error) {
	return db.GetCursorWithRange("", "")
}

func (db *Database) GetCursorWithRange(start, end string) (*Cursor, error) {
	mt, segments, indexes := db.getStateCopy()
	return NewCursor(db.dbDir, start, end, mt, segments, indexes)
}

type Pair struct {
	Key   string
	Value string
}

type KeyValueIterator interface {
	Next() (Pair, error)
}

func (db *Database) DeleteRange(start, end string) (int, error) {
	cursor, err := db.GetCursorWithRange(start, end)
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	var deleted int
	for {
		entry, err := cursor.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		err = db.Delete(entry.Key)
		if err != nil {
			return 0, err
		}
		deleted++
	}

	return deleted, nil
}
