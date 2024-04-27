package segments

import (
	"errors"
	"os"
	"sort"

	"github.com/jukeks/tukki/internal/sstable"
	"github.com/jukeks/tukki/internal/storage"
)

var ErrKeyNotFound = errors.New("key not found in segments")

func (sm *SegmentManager) Get(key string) (string, error) {
	for _, segment := range sm.GetSegmentsSorted() {
		segmentPath := storage.GetPath(sm.dbDir, segment.Filename)
		segmentFile, err := os.Open(segmentPath)
		if err != nil {
			return "", err
		}
		defer segmentFile.Close()

		reader := sstable.NewSSTableReader(segmentFile)
		for entry, err := reader.Next(); err == nil; entry, err = reader.Next() {
			if entry.Key == key {
				if entry.Deleted {
					return "", ErrKeyNotFound
				}
				return entry.Value, nil
			}
		}
	}

	return "", ErrKeyNotFound
}

func (sm *SegmentManager) GetSegmentsSorted() []Segment {
	keys := make([]SegmentId, len(sm.segments))
	i := 0
	for k := range sm.segments {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	segments := make([]Segment, len(keys))
	for i, k := range keys {
		segments[i] = sm.segments[k]
	}

	return segments
}
