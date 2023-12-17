package sstable

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jukeks/tukki/internal/keyvalue"
	"github.com/jukeks/tukki/internal/storage"
	sstablev1 "github.com/jukeks/tukki/proto/gen/tukki/storage/sstable/v1"
)

type SSTableWriter struct {
	writer io.Writer
}

func NewSSTableWriter(writer io.Writer) *SSTableWriter {
	return &SSTableWriter{
		writer: writer,
	}
}

func (w *SSTableWriter) Write(entry keyvalue.IteratorEntry) error {
	record := sstablev1.SSTableRecord{
		Key:     entry.Key,
		Value:   entry.Value,
		Deleted: entry.Deleted,
	}
	return storage.WriteLengthPrefixedProtobufMessage(w.writer, &record)
}

func (w *SSTableWriter) WriteFromIterator(iterator keyvalue.KeyValueIterator) error {
	writer := bufio.NewWriter(w.writer)

	for entry, err := iterator.Next(); err == nil; entry, err = iterator.Next() {
		err := w.Write(entry)
		if err != nil {
			return fmt.Errorf("failed to write entry: %w", err)
		}
	}

	err := writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	return nil
}

type SSTableReader struct {
	reader io.Reader
}

func NewSSTableReader(reader io.Reader) *SSTableReader {
	return &SSTableReader{
		reader: reader,
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
