package index

import (
	"bufio"
	"io"

	"github.com/jukeks/tukki/internal/storage"
	indexv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/index/v1"
)

type Index struct {
	Entries map[string]int64
}

func OpenIndex(reader io.Reader) (*Index, error) {
	br := bufio.NewReader(reader)

	entries := make(map[string]int64)
	for {
		var record indexv1.IndexEntry
		err := storage.ReadLengthPrefixedProtobufMessage(br, &record)
		if err != nil {
			if err == io.EOF {
				break
			}
			return &Index{}, err
		}
		entries[record.Key] = record.Offset
	}

	return &Index{
		Entries: entries,
	}, nil
}

func NewIndex() *Index {
	return &Index{
		Entries: make(map[string]int64),
	}
}

func NewIndexWriter(writer io.WriteCloser) *IndexWriter {
	return &IndexWriter{
		writer: writer,
		index:  NewIndex(),
	}
}

type IndexWriter struct {
	writer io.WriteCloser
	index  *Index
}

func (w *IndexWriter) Close() error {
	for key, offset := range w.index.Entries {
		record := indexv1.IndexEntry{
			Key:    key,
			Offset: offset,
		}
		_, err := storage.WriteLengthPrefixedProtobufMessage(w.writer, &record)
		if err != nil {
			return err
		}
	}

	return w.writer.Close()
}

func (w *IndexWriter) Add(key string, offset int64) error {
	w.index.Entries[key] = offset
	return nil
}
