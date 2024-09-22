package replica

import (
	"errors"

	"github.com/hashicorp/raft"
	raftbadger "github.com/rfyiamcool/raft-badger"
)

/*
	github.com/rfyiamcool/raft-badger fails initialization as it returns wrong
	types of errors on empty database. This is a wrapper around the library to
	convert the errors to the correct type. I should sent a PR to the original.
*/

type BadgerWrapper struct {
	db *raftbadger.Storage
}

var (
	ErrKeyNotFound = errors.New("not found")
)

func wrapNotFound(err error) error {
	if err == nil {
		return nil
	}
	if err == raftbadger.ErrNotFoundKey {
		return ErrKeyNotFound
	}
	return err
}

// Get implements raft.StableStore.
func (b *BadgerWrapper) Get(key []byte) ([]byte, error) {
	val, err := b.db.Get(key)
	return val, wrapNotFound(err)
}

// GetUint64 implements raft.StableStore.
func (b *BadgerWrapper) GetUint64(key []byte) (uint64, error) {
	val, err := b.db.GetUint64(key)
	return val, wrapNotFound(err)
}

// Set implements raft.StableStore.
func (b *BadgerWrapper) Set(key []byte, val []byte) error {
	return b.db.Set(key, val)
}

// SetUint64 implements raft.StableStore.
func (b *BadgerWrapper) SetUint64(key []byte, val uint64) error {
	return b.db.SetUint64(key, val)
}

// DeleteRange implements raft.LogStore.
func (b *BadgerWrapper) DeleteRange(min uint64, max uint64) error {
	return b.db.DeleteRange(min, max)
}

// FirstIndex implements raft.LogStore.
func (b *BadgerWrapper) FirstIndex() (uint64, error) {
	val, err := b.db.FirstIndex()
	return val, wrapNotFound(err)
}

// GetLog implements raft.LogStore.
func (b *BadgerWrapper) GetLog(index uint64, log *raft.Log) error {
	return b.db.GetLog(index, log)
}

// LastIndex implements raft.LogStore.
func (b *BadgerWrapper) LastIndex() (uint64, error) {
	idx, err := b.db.LastIndex()
	if err != nil {
		if err == raftbadger.ErrNotFoundLastIndex {
			return 0, nil
		}
	}

	return idx, err
}

// StoreLog implements raft.LogStore.
func (b *BadgerWrapper) StoreLog(log *raft.Log) error {
	return b.db.StoreLog(log)
}

// StoreLogs implements raft.LogStore.
func (b *BadgerWrapper) StoreLogs(logs []*raft.Log) error {
	return b.db.StoreLogs(logs)
}

func NewBadgerWrapper(db *raftbadger.Storage) *BadgerWrapper {
	return &BadgerWrapper{
		db: db,
	}
}
