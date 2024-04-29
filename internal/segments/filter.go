package segments

import (
	"bytes"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/jukeks/tukki/internal/storage"
)

type SegmentBlooms struct {
	filter *bloom.BloomFilter
}

func NewSegmentBlooms(n uint) *SegmentBlooms {
	return &SegmentBlooms{
		filter: bloom.NewWithEstimates(n, 0.01),
	}
}

func OpenSegmentBlooms(dbDir string, filename storage.Filename) (*SegmentBlooms, error) {
	path := storage.GetPath(dbDir, filename)
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var g bloom.BloomFilter
	buff := bytes.NewBuffer(b)
	_, err = g.ReadFrom(buff)
	if err != nil {
		return nil, err
	}

	return &SegmentBlooms{
		filter: &g,
	}, nil
}

func (sb *SegmentBlooms) Save(dbDir string, filename storage.Filename) error {
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

func (sb *SegmentBlooms) Add(key string) {
	sb.filter.Add([]byte(key))
}

func (sb *SegmentBlooms) Test(key string) bool {
	return sb.filter.Test([]byte(key))
}
