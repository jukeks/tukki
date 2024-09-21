package db

import (
	"github.com/jukeks/tukki/internal/db/segments"
	"github.com/jukeks/tukki/internal/storage"
	"github.com/jukeks/tukki/internal/storage/memtable"
)

type LiveSegment struct {
	WalFilename storage.Filename
	Segment     segments.SegmentMetadata
	Memtable    memtable.Memtable
	Wal         *memtable.Wal
}

func NewLiveSegment(id segments.SegmentId) *LiveSegment {
	return &LiveSegment{
		Segment: segments.SegmentMetadata{
			Id:          id,
			SegmentFile: getSegmentFilename(id),
			MembersFile: getMembersFilename(id),
			IndexFile:   getIndexFilename(id),
		},
		WalFilename: getWalFilename(id),
	}
}

func (ls *LiveSegment) Open(dbDir string) error {
	if ls.Memtable != nil {
		panic("live segment already opened")
	}

	ls.Memtable = memtable.NewMemtable()
	wal, err := memtable.OpenWal(dbDir, ls.WalFilename, ls.Memtable)
	if err != nil {
		return err
	}
	ls.Wal = wal
	return nil
}

func (ls *LiveSegment) Close() error {
	return ls.Wal.Close()
}

func (d *LiveSegment) Set(key, value string) error {
	err := d.Wal.Set(key, value)
	if err != nil {
		return err
	}

	d.Memtable.Insert(key, value)
	return nil
}

func (d *LiveSegment) Get(key string) (string, error) {
	value, found := d.Memtable.Get(key)
	if found {
		if value.Deleted {
			return "", ErrKeyNotFound
		}
		return value.Value, nil
	}

	return "", ErrKeyNotFound
}

func (d *LiveSegment) Delete(key string) error {
	err := d.Wal.Delete(key)
	if err != nil {
		return err
	}

	d.Memtable.Delete(key)
	return nil
}
