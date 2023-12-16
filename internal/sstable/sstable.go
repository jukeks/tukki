package sstable

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/jukeks/tukki/internal/keyvalue"
	"github.com/jukeks/tukki/internal/storage"
	sstablev1 "github.com/jukeks/tukki/proto/gen/tukki/storage/sstable/v1"
	"google.golang.org/protobuf/proto"
)

type SSTableWriter struct {
	writer io.Writer
}

func NewSSTableWriter(writer io.Writer) *SSTableWriter {
	return &SSTableWriter{
		writer: writer,
	}
}

func (w *SSTableWriter) Write(iterator keyvalue.KeyValueIterator) (int, error) {
	writer := bufio.NewWriter(w.writer)

	written := 0
	for iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()

		payload, err := proto.Marshal(&sstablev1.SSTableRecord{
			Key:     key,
			Value:   value.Value,
			Deleted: value.Deleted,
		})
		if err != nil {
			return written, fmt.Errorf("failed to serialize key value: %w", err)
		}

		err = binary.Write(writer, binary.LittleEndian, uint32(len(payload)))
		if err != nil {
			return written, fmt.Errorf("failed to write payload len: %w", err)
		}

		_, err = writer.Write(payload)
		if err != nil {
			return written, fmt.Errorf("failed to write payload: %w", err)
		}
		written++
	}

	err := writer.Flush()
	if err != nil {
		return written, fmt.Errorf("failed to flush: %w", err)
	}

	return written, nil
}

type SSTableReader struct {
	reader io.Reader
}

func NewSSTableReader(reader io.Reader) *SSTableReader {
	return &SSTableReader{
		reader: reader,
	}
}

func (r *SSTableReader) Read() (keyvalue.KeyValueIterator, error) {
	reader := bufio.NewReader(r.reader)
	return newSSTableIterator(
		reader,
	), nil
}

type sstableIterator struct {
	reader  io.Reader
	current *sstablev1.SSTableRecord
}

func newSSTableIterator(reader io.Reader) *sstableIterator {
	return &sstableIterator{
		reader: reader,
	}
}

func (i *sstableIterator) Key() string {
	return i.current.Key
}

func (i *sstableIterator) Value() keyvalue.Value {
	return keyvalue.Value{
		Value:   i.current.Value,
		Deleted: i.current.Deleted,
	}
}

func (i *sstableIterator) Next() bool {
	var record sstablev1.SSTableRecord
	err := storage.ReadLengthPrefixedProtobufMessage(i.reader, &record)
	if err != nil {
		if err == io.EOF {
			return false
		}

		log.Fatalf("failed to read record: %v", err)
	}

	i.current = &record
	return true
}
