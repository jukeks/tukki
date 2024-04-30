package segmentmembers

import (
	"bytes"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/jukeks/tukki/internal/keyvalue"
	"github.com/jukeks/tukki/internal/storage"
)

type SegmentMembers struct {
	filter *bloom.BloomFilter
}

func NewSegmentMembers(n uint) *SegmentMembers {
	return &SegmentMembers{
		filter: bloom.NewWithEstimates(n, 0.01),
	}
}

func OpenSegmentMembers(dbDir string, filename storage.Filename) (*SegmentMembers, error) {
	path := storage.GetPath(dbDir, filename)
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var f bloom.BloomFilter
	buff := bytes.NewBuffer(b)
	_, err = f.ReadFrom(buff)
	if err != nil {
		return nil, err
	}

	return &SegmentMembers{
		filter: &f,
	}, nil
}

func (sb *SegmentMembers) Save(dbDir string, filename storage.Filename) error {
	path := storage.GetPath(dbDir, filename)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = sb.filter.WriteTo(f)
	if err != nil {
		return err
	}

	return nil
}

func (sb *SegmentMembers) Add(key string) {
	sb.filter.AddString(key)
}

func (sb *SegmentMembers) Contains(key string) bool {
	return sb.filter.Test([]byte(key))
}

func (sb *SegmentMembers) Fill(iterator keyvalue.KeyValueIterator) {
	for {
		entry, err := iterator.Next()
		if err != nil {
			break
		}

		sb.Add(entry.Key)
	}
}

func (sb *SegmentMembers) Size() uint {
	return uint(sb.filter.ApproximatedSize())
}
