package db

import (
	"errors"
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

type Iterator struct {
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

type sstableState struct {
	reader *sstable.SSTableReader
	file   *os.File

	current keyvalue.IteratorEntry
	err     error
}

func (s *sstableState) Close() {
	s.file.Close()
}

func (s *sstableState) Get() (keyvalue.IteratorEntry, error) {
	return s.current, s.err
}

func (s *sstableState) Progress() {
	s.current, s.err = s.reader.Next()
}

type memtableState struct {
	memtable memtable.Memtable
	iterator keyvalue.KeyValueIterator
	current  keyvalue.IteratorEntry
	err      error
}

func (m *memtableState) Close() {
}

func (m *memtableState) Get() (keyvalue.IteratorEntry, error) {
	return m.current, m.err
}

func (m *memtableState) Progress() {
	m.current, m.err = m.iterator.Next()
}

func NewIterator(dbDir,
	start, end string, mt memtable.Memtable,
	sortedSegments []segments.SegmentMetadata,
	indexes map[segments.SegmentId]*index.Index) (*Iterator, error) {
	iterator := &Iterator{
		dbDir: dbDir, memtable: mt,
		allSegments: sortedSegments, indexes: indexes, start: start, end: end}

	if start != "" {
		err := iterator.seek()
		if err != nil {
			return nil, fmt.Errorf("failed to seek: %w", err)
		}
	}

	return iterator, nil
}

func (i *Iterator) open() error {
	i.openedSegments = make([]subIterator, 0)
	mtIterator := i.memtable.Iterate()
	mtState := &memtableState{
		memtable: i.memtable,
		iterator: mtIterator,
	}
	mtState.Progress()
	i.openedSegments = append(i.openedSegments, mtState)

	for _, segment := range i.allSegments {
		segmentFile, err := files.OpenFile(i.dbDir, segment.SegmentFile)
		if err != nil {
			return fmt.Errorf("failed to open segment file: %w", err)
		}
		reader := sstable.NewSSTableReader(segmentFile)
		next, err := reader.Next()
		opened := &sstableState{
			file:    segmentFile,
			reader:  reader,
			current: next,
			err:     err,
		}
		i.openedSegments = append(i.openedSegments, opened)
	}
	i.opened = true

	return nil
}

func (i *Iterator) Next() (keyvalue.IteratorEntry, error) {
	if !i.opened {
		if err := i.open(); err != nil {
			return keyvalue.IteratorEntry{},
				fmt.Errorf("failed to open iterator: %w", err)
		}
	}

	// find first error free segment that is not past the end
	result := i.openedSegments[0]
	canProceed := false
	for _, segment := range i.openedSegments {
		current, err := segment.Get()
		if err == nil && (current.Key < i.end || i.end == "") {
			result = segment
			canProceed = true
			break
		}
	}
	if !canProceed {
		return keyvalue.IteratorEntry{}, io.EOF
	}

	for _, segment := range i.openedSegments {
		current, err := segment.Get()
		if err != nil {
			if err != io.EOF {
				return keyvalue.IteratorEntry{},
					fmt.Errorf("failed to read next entry: %w", err)
			}
			continue
		}

		if i.end != "" && current.Key > i.end {
			continue
		}

		currentResult, _ := result.Get()
		if current.Key < currentResult.Key {
			result = segment
		}
	}

	ret, _ := result.Get()

	// there can be earlier entries in other segments
	// we need to advance them to the next entry
	// to avoid exposing old value
	for _, segment := range i.openedSegments {
		current, _ := segment.Get()
		if current.Key == ret.Key {
			segment.Progress()
		}
	}

	return ret, nil
}

func (i *Iterator) Close() {
	for _, opened := range i.openedSegments {
		opened.Close()
	}
	i.opened = false
}

var errIteratorAlreadyOpened = errors.New("iterator already opened")
var errStartNotSet = errors.New("start key not set")

func (i *Iterator) seek() error {
	if i.opened {
		// not handled yet
		return errIteratorAlreadyOpened
	}

	if i.start == "" {
		return errStartNotSet
	}

	i.openedSegments = make([]subIterator, 0)

	mtIterator := i.memtable.Iterate()
	mtState := &memtableState{
		memtable: i.memtable,
		iterator: mtIterator,
	}
	mtState.Progress()
	// seek in memtable
	for {
		current, _ := mtState.Get()
		if current.Key >= i.start {
			break
		}
		mtState.Progress()
	}
	i.openedSegments = append(i.openedSegments, mtState)

	for _, segment := range i.allSegments {
		index := i.indexes[segment.Id]
		found := false
		offset := uint64(0)
		for _, entry := range index.EntryList {
			if entry.Key >= i.start {
				offset = entry.Offset
				found = true
				break
			}
		}
		if !found {
			continue
		}

		segmentFile, err := files.OpenFile(i.dbDir, segment.SegmentFile)
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
		opened := &sstableState{
			file:    segmentFile,
			reader:  reader,
			current: next,
			err:     err,
		}
		i.openedSegments = append(i.openedSegments, opened)
	}
	i.opened = true

	if len(i.openedSegments) == 0 {
		return io.EOF
	}

	return nil
}
