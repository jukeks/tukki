package memtable

import (
	"bytes"
	"io"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/jukeks/tukki/internal/storage/keyvalue"
)

type Memtable interface {
	Get(key []byte) (keyvalue.Value, bool)
	Insert(key []byte, value []byte)
	Delete(key []byte)
	Iterate() keyvalue.KeyValueIterator
	MemberCount() int
	Size() uint64
	Copy() Memtable
}

func NewMemtable() Memtable {
	t := redblacktree.NewWithStringComparator()
	return &memtableRedBlackTree{
		t: t,
	}
}

type memtableRedBlackTree struct {
	t    *redblacktree.Tree
	size uint64
}

func (m *memtableRedBlackTree) Get(key []byte) (keyvalue.Value, bool) {
	value, found := m.t.Get(string(key))
	if !found {
		return keyvalue.Value{}, false
	}

	result := value.(keyvalue.Value)
	result.Value = bytes.Clone(result.Value)
	return result, true
}

func (m *memtableRedBlackTree) Delete(key []byte) {
	m.size += uint64(len(key))
	m.t.Put(string(key), keyvalue.Value{
		Deleted: true,
	})
}

func (m *memtableRedBlackTree) Insert(key, value []byte) {
	m.size += uint64(len(key) + len(value))
	m.t.Put(string(key), keyvalue.Value{
		Value: bytes.Clone(value),
	})
}

func (m *memtableRedBlackTree) Copy() Memtable {
	t := redblacktree.NewWithStringComparator()
	iter := m.t.Iterator()
	for iter.Next() {
		key := iter.Key().(string)
		value := iter.Value().(keyvalue.Value)
		t.Put(key, value)
	}

	return &memtableRedBlackTree{
		t:    t,
		size: m.size,
	}
}

type memtableRedBlackTreeIterator struct {
	iter *redblacktree.Iterator
}

func (i *memtableRedBlackTreeIterator) Next() (keyvalue.IteratorEntry, error) {
	if !i.iter.Next() {
		return keyvalue.IteratorEntry{}, io.EOF
	}

	key := i.iter.Key().(string)
	value := i.iter.Value().(keyvalue.Value)
	return keyvalue.IteratorEntry{
		Key:     []byte(key),
		Value:   bytes.Clone(value.Value),
		Deleted: value.Deleted,
	}, nil
}

func (m *memtableRedBlackTree) Iterate() keyvalue.KeyValueIterator {
	iter := m.t.Iterator()
	return &memtableRedBlackTreeIterator{
		iter: &iter,
	}
}

func (m *memtableRedBlackTree) MemberCount() int {
	return m.t.Size()
}

func (m *memtableRedBlackTree) Size() uint64 {
	return m.size
}
