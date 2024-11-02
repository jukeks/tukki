package db

import (
	"fmt"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/index"
	"github.com/jukeks/tukki/internal/storage/keyvalue"
	"github.com/jukeks/tukki/internal/storage/memtable"
	"github.com/jukeks/tukki/internal/storage/segments"
	"github.com/jukeks/tukki/internal/storage/sstable"
)

type Cursor struct {
	dbDir          string
	allSegments    []segments.SegmentMetadata
	openedSegments []keyvalue.SubIterator
	indexes        map[segments.SegmentId]*index.Index
	opened         bool
	memtable       memtable.Memtable
	iterator       keyvalue.Iterator

	start string
	end   string
}

func NewCursor(dbDir,
	start, end string, mt memtable.Memtable,
	sortedSegments []segments.SegmentMetadata,
	indexes map[segments.SegmentId]*index.Index) (*Cursor, error) {
	cursor := &Cursor{
		dbDir: dbDir, memtable: mt,
		allSegments: sortedSegments, indexes: indexes, start: start, end: end}

	if err := cursor.open(); err != nil {
		return nil, fmt.Errorf("failed to open iterator: %w", err)
	}

	iter, err := keyvalue.NewIterator(start, end, false, cursor.openedSegments...)
	if err != nil {
		return nil, fmt.Errorf("failed to create iterator: %w", err)
	}

	cursor.iterator = iter

	return cursor, nil
}

func (c *Cursor) open() error {
	c.openedSegments = make([]keyvalue.SubIterator, 0)
	c.openedSegments = append(c.openedSegments, memtable.NewMemtableIterator(c.memtable))

	for _, segment := range c.allSegments {
		segmentFile, err := files.OpenFile(c.dbDir, segment.SegmentFile)
		if err != nil {
			return fmt.Errorf("failed to open segment file: %w", err)
		}
		subIter := sstable.NewSSTableIterator(segmentFile, c.indexes[segment.Id])
		c.openedSegments = append(c.openedSegments, subIter)
	}
	c.opened = true

	return nil
}

func (c *Cursor) Next() (Pair, error) {
	next, err := c.iterator.Next()
	if err != nil {
		return Pair{}, err
	}

	return Pair{Key: next.Key, Value: next.Value}, nil
}

func (c *Cursor) Close() {
	for _, opened := range c.openedSegments {
		opened.Close()
	}
	c.opened = false
}
