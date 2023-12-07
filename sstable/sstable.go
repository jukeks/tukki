package sstable

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/jukeks/tukki/memtable"
	"google.golang.org/protobuf/proto"
)

type SSTable struct {
}

type SSTableWriter struct {
	writer io.Writer
}

func NewSSTableWriter(writer io.Writer) *SSTableWriter {
	return &SSTableWriter{
		writer: writer,
	}
}

func (w *SSTableWriter) Write(iterator memtable.KeyValueIterator) error {
	writer := bufio.NewWriter(w.writer)

	for iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()

		payload, err := proto.Marshal(&SSTableRecord{
			Key:   key,
			Value: value,
		})
		if err != nil {
			return fmt.Errorf("failed to serialize key value: %w", err)
		}

		err = binary.Write(writer, binary.LittleEndian, uint32(len(payload)))
		if err != nil {
			return fmt.Errorf("failed to write payload len: %w", err)
		}

		_, err = writer.Write(payload)
		if err != nil {
			return fmt.Errorf("failed to write payload: %w", err)
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

func (r *SSTableReader) Read() (memtable.KeyValueIterator, error) {
	reader := bufio.NewReader(r.reader)
	return newSSTableIterator(
		reader,
	), nil
}

type sstableIterator struct {
	reader  io.Reader
	current *SSTableRecord
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
	len := uint32(0)
	err := binary.Read(i.reader, binary.LittleEndian, &len)
	if err != nil {
		log.Printf("failed to read payload len: %v", err)
		return false
	}

	payload := make([]byte, len)
	_, err = i.reader.Read(payload)
	if err != nil {
		log.Printf("failed to read payload: %v", err)
		return false
	}

	record := &SSTableRecord{}
	err = proto.Unmarshal(payload, record)
	if err != nil {
		log.Printf("failed to unmarshal payload: %v", err)
		return false
	}

	i.current = record
	return true
}
