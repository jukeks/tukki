package db

import (
	"fmt"
	"io"
	"os"

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

type sstableSubIterator struct {
	reader      *sstable.SSTableReader
	segmentFile *os.File
	index       *index.Index

	current keyvalue.IteratorEntry
	err     error
}

func (s *sstableSubIterator) Close() {
	s.segmentFile.Close()
}

func (s *sstableSubIterator) Get() (keyvalue.IteratorEntry, error) {
	return s.current, s.err
}

func (s *sstableSubIterator) Progress() {
	s.current, s.err = s.reader.Next()
}

func (s *sstableSubIterator) Seek(key string) error {
	index := s.index
	found := false
	offset := uint64(0)
	for _, entry := range index.EntryList {
		if entry.Key >= key {
			offset = entry.Offset
			found = true
			break
		}
	}
	if !found {
		return io.EOF
	}

	_, err := s.segmentFile.Seek(int64(offset), 0)
	if err != nil {
		s.segmentFile.Close()
		return fmt.Errorf("failed to seek in segment file: %w", err)
	}

	s.reader = sstable.NewSSTableReader(s.segmentFile)
	s.current, s.err = s.reader.Next()
	return nil
}

type memtableSubIterator struct {
	memtable memtable.Memtable
	iterator keyvalue.KeyValueIterator
	current  keyvalue.IteratorEntry
	err      error
}

func (m *memtableSubIterator) Close() {
}

func (m *memtableSubIterator) Get() (keyvalue.IteratorEntry, error) {
	return m.current, m.err
}

func (m *memtableSubIterator) Progress() {
	m.current, m.err = m.iterator.Next()
}

func (m *memtableSubIterator) Seek(key string) error {
	m.iterator = m.memtable.Iterate()
	for {
		m.Progress()
		if m.err != nil {
			return m.err
		}

		if m.current.Key >= key {
			return nil
		}
	}
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
	mtIterator := c.memtable.Iterate()
	mtState := &memtableSubIterator{
		memtable: c.memtable,
		iterator: mtIterator,
	}
	mtState.Progress()
	c.openedSegments = append(c.openedSegments, mtState)

	for _, segment := range c.allSegments {
		segmentFile, err := files.OpenFile(c.dbDir, segment.SegmentFile)
		if err != nil {
			return fmt.Errorf("failed to open segment file: %w", err)
		}
		reader := sstable.NewSSTableReader(segmentFile)
		next, err := reader.Next()
		opened := &sstableSubIterator{
			segmentFile: segmentFile,
			index:       c.indexes[segment.Id],
			reader:      reader,
			current:     next,
			err:         err,
		}
		c.openedSegments = append(c.openedSegments, opened)
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
