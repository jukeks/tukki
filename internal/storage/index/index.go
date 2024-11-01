package index

import (
	"bufio"
	"io"
	"sort"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/marshalling"
	indexv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/index/v1"
)

type IndexEntry struct {
	Key    string
	Offset uint64
}

type Index struct {
	Entries   map[string]uint64
	EntryList []IndexEntry
}

func OpenIndex(dbDir string, filename files.Filename) (*Index, error) {
	f, err := files.OpenFile(dbDir, filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	entries := make(map[string]uint64)
	entryList := make([]IndexEntry, 0)
	for {
		var record indexv1.IndexEntry
		err := marshalling.ReadLengthPrefixedProtobufMessage(reader, &record)
		if err != nil {
			if err == io.EOF {
				break
			}
			return &Index{}, err
		}
		entries[record.Key] = record.Offset
		entryList = append(entryList, IndexEntry{
			Key:    record.Key,
			Offset: record.Offset,
		})
	}

	sort.Slice(entryList, func(i, j int) bool {
		return entryList[i].Key < entryList[j].Key
	})

	return &Index{
		Entries:   entries,
		EntryList: entryList,
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
		_, err := marshalling.WriteLengthPrefixedProtobufMessage(bw, &record)
		if err != nil {
			return err
		}
	}

	return bw.Flush()
}

func (w *IndexWriter) Close() error {
	return w.writer.Close()
}
