package db

import (
	"bytes"
	"errors"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/journal"
	"github.com/jukeks/tukki/internal/storage/memtable"
	"github.com/jukeks/tukki/internal/storage/segments"
)

type LiveSegment struct {
	WalFilename files.Filename
	Segment     segments.SegmentMetadata
	Memtable    memtable.Memtable
	Wal         *memtable.Wal
}

func NewLiveSegment(id segments.SegmentId) *LiveSegment {
	return &LiveSegment{
		Segment: segments.SegmentMetadata{
			Id:          id,
			SegmentFile: segments.GetSegmentFilename(id),
			MembersFile: segments.GetMembersFilename(id),
			IndexFile:   segments.GetIndexFilename(id),
		},
		WalFilename: segments.GetWalFilename(id),
	}
}

func (ls *LiveSegment) Open(dbDir string, mode journal.WriteMode) error {
	if ls.Memtable != nil {
		panic("live segment already opened")
	}

	ls.Memtable = memtable.NewMemtable()
	wal, err := memtable.OpenWal(dbDir, ls.WalFilename, mode, ls.Memtable)
	if err != nil {
		return err
	}
	ls.Wal = wal
	return nil
}

func (ls *LiveSegment) Close() error {
	return ls.Wal.Close()
}

func (d *LiveSegment) Set(key, value []byte) error {
	err := d.Wal.Set(key, value)
	if err != nil {
		return err
	}

	d.Memtable.Insert(key, value)
	return nil
}

var errTombstone = errors.New("tombstone")

func (d *LiveSegment) Get(key []byte) ([]byte, error) {
	value, found := d.Memtable.Get(key)
	if found {
		if value.Deleted {
			return nil, errTombstone
		}
		return bytes.Clone(value.Value), nil
	}

	return nil, ErrKeyNotFound
}

func (d *LiveSegment) Delete(key []byte) error {
	err := d.Wal.Delete(key)
	if err != nil {
		return err
	}

	d.Memtable.Delete(key)
	return nil
}
