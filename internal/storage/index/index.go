package index

import (
	"bufio"
	"io"

	"github.com/jukeks/tukki/internal/storage"
	indexv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/index/v1"
)

type Index struct {
	Entries map[string]uint64
}

func OpenIndex(dbDir string, filename storage.Filename) (*Index, error) {
	f, err := storage.OpenFile(dbDir, filename)
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

type OffsetMap map[string]uint64

func (w *IndexWriter) WriteFromOffsets(offsets OffsetMap) error {
	bw := bufio.NewWriter(w.writer)
	for key, offset := range offsets {
		record := indexv1.IndexEntry{
			Key:    key,
			Offset: offset,
		}
		_, err := storage.WriteLengthPrefixedProtobufMessage(bw, &record)
		if err != nil {
			return err
		}
	}

	return bw.Flush()
}

func (w *IndexWriter) Close() error {
	return w.writer.Close()
}
