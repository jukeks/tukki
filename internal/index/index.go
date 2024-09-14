package index

import (
	"bufio"
	"io"
	"os"

	"github.com/jukeks/tukki/internal/sstable"
	"github.com/jukeks/tukki/internal/storage"
	indexv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/index/v1"
)

type Index struct {
	Entries map[string]uint64
}

func OpenIndex(dbDir string, filename storage.Filename) (*Index, error) {
	path := storage.GetPath(dbDir, filename)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	entries := make(map[string]uint64)
	for {
		var record indexv1.IndexEntry
		err := storage.ReadLengthPrefixedProtobufMessage(reader, &record)
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
		Entries: make(map[string]uint64),
	}
}

type IndexWriter struct {
	writer io.WriteCloser
}

func NewIndexWriter(writer io.WriteCloser) *IndexWriter {
	return &IndexWriter{
		writer: writer,
	}
}

func (w *IndexWriter) WriteFromOffsets(offsets sstable.KeyMap) error {
	for key, offset := range offsets {
		record := indexv1.IndexEntry{
			Key:    key,
			Offset: offset,
		}
		_, err := storage.WriteLengthPrefixedProtobufMessage(w.writer, &record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *IndexWriter) Close() error {
	return w.writer.Close()
}