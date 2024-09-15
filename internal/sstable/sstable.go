package sstable

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jukeks/tukki/internal/index"
	"github.com/jukeks/tukki/internal/keyvalue"
	"github.com/jukeks/tukki/internal/storage"
	sstablev1 "github.com/jukeks/tukki/proto/gen/tukki/storage/sstable/v1"
)

type SSTableWriter struct {
	writer    io.Writer
	offsetMap index.OffsetMap
}

func NewSSTableWriter(writer io.Writer) *SSTableWriter {
	return &SSTableWriter{
		writer:    writer,
		offsetMap: make(index.OffsetMap),
	}
}

func (w *SSTableWriter) Write(entry keyvalue.IteratorEntry) (uint32, error) {
	record := sstablev1.SSTableRecord{
		Key:     entry.Key,
		Value:   entry.Value,
		Deleted: entry.Deleted,
	}
	return storage.WriteLengthPrefixedProtobufMessage(w.writer, &record)
}

func (w *SSTableWriter) WriteFromIterator(iterator keyvalue.KeyValueIterator) error {
	writer := bufio.NewWriter(w.writer)

	var offset uint64 = 0
	for entry, err := iterator.Next(); err == nil; entry, err = iterator.Next() {
		len, err := w.Write(entry)
		if err != nil {
			return fmt.Errorf("failed to write entry: %w", err)
		}
		w.offsetMap[entry.Key] = offset
		offset += uint64(len)
	}

	err := writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	return nil
}

func (w *SSTableWriter) WrittenOffsets() index.OffsetMap {
	return w.offsetMap
}

type SSTableReader struct {
	reader *bufio.Reader
}

func NewSSTableReader(reader io.Reader) *SSTableReader {
	return &SSTableReader{
		reader: bufio.NewReader(reader),
	}
}

func (i *SSTableReader) Next() (keyvalue.IteratorEntry, error) {
	var record sstablev1.SSTableRecord
	err := storage.ReadLengthPrefixedProtobufMessage(i.reader, &record)
	if err != nil {
		return keyvalue.IteratorEntry{}, err
	}

	return keyvalue.IteratorEntry{
		Key:     record.Key,
		Value:   record.Value,
		Deleted: record.Deleted,
	}, nil
}

type SSTableSeeker struct {
	reader io.ReadSeeker
}

func NewSSTableSeeker(reader io.ReadSeeker) *SSTableSeeker {
	return &SSTableSeeker{
		reader: reader,
	}
}

func (r *SSTableSeeker) ReadAt(offset uint64) (keyvalue.IteratorEntry, error) {
	_, err := r.reader.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return keyvalue.IteratorEntry{}, err
	}

	var record sstablev1.SSTableRecord
	err = storage.ReadLengthPrefixedProtobufMessage(r.reader, &record)
	if err != nil {
		return keyvalue.IteratorEntry{}, err
	}

	return keyvalue.IteratorEntry{
		Key:     record.Key,
		Value:   record.Value,
		Deleted: record.Deleted,
	}, nil
}
