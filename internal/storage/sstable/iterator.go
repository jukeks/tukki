package sstable

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jukeks/tukki/internal/storage/index"
	"github.com/jukeks/tukki/internal/storage/keyvalue"
)

func NewSSTableIterator(segmentFile *os.File, index *index.Index) keyvalue.SubIterator {
	iter := &sstableSubIterator{
		segmentFile: segmentFile,
		index:       index,
		reader:      NewSSTableReader(segmentFile),
	}

	iter.Progress()
	return iter
}

type sstableSubIterator struct {
	segmentFile *os.File
	index       *index.Index

	reader *SSTableReader

	current keyvalue.IteratorEntry
	err     error
}

func (s *sstableSubIterator) Close() error {
	return s.segmentFile.Close()
}

func (s *sstableSubIterator) Get() (keyvalue.IteratorEntry, error) {
	return s.current, s.err
}

func (s *sstableSubIterator) Progress() {
	s.current, s.err = s.reader.Next()
	if s.err != nil && s.err != io.EOF {
		log.Printf("failed to read next entry from file %s: %v", s.segmentFile.Name(), s.err)
	}
}

var ErrIndexNotProvided = errors.New("index not provided")

func (s *sstableSubIterator) Seek(key string) error {
	found := false
	offset := uint64(0)
	if s.index == nil {
		return ErrIndexNotProvided
	}

	for _, entry := range s.index.EntryList {
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

	s.reader = NewSSTableReader(s.segmentFile)
	s.current, s.err = s.reader.Next()
	return nil
}
