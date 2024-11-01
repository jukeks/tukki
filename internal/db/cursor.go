package db

import (
	"errors"
	"fmt"
	"io"
	"log"
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
	openedSegments []subIterator
	indexes        map[segments.SegmentId]*index.Index
	opened         bool
	memtable       memtable.Memtable

	start string
	end   string
}

type subIterator interface {
	Get() (keyvalue.IteratorEntry, error)
	Progress()
	Close()
}

type sstableSubIterator struct {
	reader *sstable.SSTableReader
	file   *os.File

	current keyvalue.IteratorEntry
	err     error
}

func (s *sstableSubIterator) Close() {
	s.file.Close()
}

func (s *sstableSubIterator) Get() (keyvalue.IteratorEntry, error) {
	return s.current, s.err
}

func (s *sstableSubIterator) Progress() {
	s.current, s.err = s.reader.Next()
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

func NewCursor(dbDir,
	start, end string, mt memtable.Memtable,
	sortedSegments []segments.SegmentMetadata,
	indexes map[segments.SegmentId]*index.Index) (*Cursor, error) {
	cursor := &Cursor{
		dbDir: dbDir, memtable: mt,
		allSegments: sortedSegments, indexes: indexes, start: start, end: end}

	if start != "" {
		err := cursor.seek()
		if err != nil {
			return nil, fmt.Errorf("failed to seek: %w", err)
		}
	}

	return cursor, nil
}

func (c *Cursor) open() error {
	c.openedSegments = make([]subIterator, 0)
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
			file:    segmentFile,
			reader:  reader,
			current: next,
			err:     err,
		}
		c.openedSegments = append(c.openedSegments, opened)
	}
	c.opened = true

	return nil
}

func (c *Cursor) Next() (keyvalue.IteratorEntry, error) {
	if !c.opened {
		if err := c.open(); err != nil {
			return keyvalue.IteratorEntry{},
				fmt.Errorf("failed to open iterator: %w", err)
		}
	}

	for {
		// find first error free segment that is not past the end
		result := c.openedSegments[0]
		canProceed := false
		for _, segment := range c.openedSegments {
			log.Printf("checking segment %+v", segment)
			current, err := segment.Get()
			if err == nil && (current.Key <= c.end || c.end == "") {
				result = segment
				canProceed = true
				break
			}
		}
		if !canProceed {
			return keyvalue.IteratorEntry{}, io.EOF
		}

		for _, segment := range c.openedSegments {
			current, err := segment.Get()
			if err != nil {
				if err != io.EOF {
					return keyvalue.IteratorEntry{},
						fmt.Errorf("failed to read next entry: %w", err)
				}
				continue
			}

			if c.end != "" && current.Key > c.end {
				continue
			}

			currentResult, _ := result.Get()
			if current.Key < currentResult.Key {
				log.Printf("current key %s is less than result key %s", current.Key, currentResult.Key)
				result = segment
			}
		}

		ret, _ := result.Get()

		// there can be earlier entries in other segments
		// we need to advance them to the next entry
		// to avoid exposing old value
		for _, segment := range c.openedSegments {
			current, _ := segment.Get()
			if current.Key == ret.Key {
				segment.Progress()
			}
		}

		if ret.Deleted {
			log.Printf("deleted key %s, retry", ret.Key)
			continue
		}

		return ret, nil
	}
}

func (c *Cursor) Close() {
	for _, opened := range c.openedSegments {
		opened.Close()
	}
	c.opened = false
}

var errIteratorAlreadyOpened = errors.New("iterator already opened")
var errStartNotSet = errors.New("start key not set")

func (c *Cursor) seek() error {
	if c.opened {
		// not handled yet
		return errIteratorAlreadyOpened
	}

	if c.start == "" {
		return errStartNotSet
	}

	c.openedSegments = make([]subIterator, 0)

	mtIterator := c.memtable.Iterate()
	mtState := &memtableSubIterator{
		memtable: c.memtable,
		iterator: mtIterator,
	}
	mtState.Progress()
	// seek in memtable
	for {
		current, _ := mtState.Get()
		if current.Key >= c.start {
			break
		}
		mtState.Progress()
	}
	c.openedSegments = append(c.openedSegments, mtState)

	for _, segment := range c.allSegments {
		index := c.indexes[segment.Id]
		found := false
		offset := uint64(0)
		for _, entry := range index.EntryList {
			if entry.Key >= c.start {
				offset = entry.Offset
				found = true
				break
			}
		}
		if !found {
			continue
		}

		segmentFile, err := files.OpenFile(c.dbDir, segment.SegmentFile)
		if err != nil {
			return fmt.Errorf("failed to open segment file: %w", err)
		}
		_, err = segmentFile.Seek(int64(offset), 0)
		if err != nil {
			segmentFile.Close()
			return fmt.Errorf("failed to seek in segment file: %w", err)
		}

		reader := sstable.NewSSTableReader(segmentFile)
		next, err := reader.Next()
		opened := &sstableSubIterator{
			file:    segmentFile,
			reader:  reader,
			current: next,
			err:     err,
		}
		c.openedSegments = append(c.openedSegments, opened)
	}
	c.opened = true

	if len(c.openedSegments) == 0 {
		return io.EOF
	}

	return nil
}
