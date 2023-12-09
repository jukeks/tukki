package sstable

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/jukeks/tukki/internal/memtable"
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

func (w *SSTableWriter) Write(iterator memtable.KeyValueIterator) (int, error) {
	writer := bufio.NewWriter(w.writer)

	written := 0
	for iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()

		payload, err := proto.Marshal(&sstablev1.SSTableRecord{
			Key:   key,
			Value: value,
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

func (r *SSTableReader) Read() (memtable.KeyValueIterator, error) {
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

func (i *sstableIterator) Value() string {
	return i.current.Value
}

func (i *sstableIterator) Next() bool {
	length := uint32(0)
	err := binary.Read(i.reader, binary.LittleEndian, &length)
	if err != nil {
		if err == io.EOF {
			return false
		}

		log.Fatalf("failed to read payload len: %v", err)
		return false
	}

	payload := make([]byte, length)
	n, err := io.ReadFull(i.reader, payload)
	if err != nil {
		log.Fatalf("failed to read payload: %v", err)
		return false
	}

	if n != int(length) {
		log.Fatalf("failed to read payload of len %d vs %d", length, n)
		return false
	}

	record := &sstablev1.SSTableRecord{}
	err = proto.Unmarshal(payload, record)
	if err != nil {
		log.Fatalf("failed to unmarshal payload of len %d vs %d: %v: %v", length, len(payload), err, payload)
		return false
	}

	i.current = record
	return true
}
