package memtable

import (
	"io"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/jukeks/tukki/internal/storage/keyvalue"
)

type Memtable interface {
	Get(key string) (keyvalue.Value, bool)
	Insert(key string, value string)
	Delete(key string)
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

func (m *memtableRedBlackTree) Get(key string) (keyvalue.Value, bool) {
	value, found := m.t.Get(string(key))
	if !found {
		return keyvalue.Value{}, false
	}

	return value.(keyvalue.Value), true
}

func (m *memtableRedBlackTree) Delete(key string) {
	m.size += uint64(len(key))
	m.t.Put(string(key), keyvalue.Value{
		Deleted: true,
	})
}

func (m *memtableRedBlackTree) Insert(key, value string) {
	m.size += uint64(len(key) + len(value))
	m.t.Put(string(key), keyvalue.Value{
		Value: value,
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
		Key:     key,
		Value:   value.Value,
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
