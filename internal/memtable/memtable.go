package memtable

import (
	"io"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/jukeks/tukki/internal/keyvalue"
)

type Memtable interface {
	Get(key string) (keyvalue.Value, bool)
	Insert(key string, value string)
	Delete(key string)
	Iterate() keyvalue.KeyValueIterator
}

func NewMemtable() Memtable {
	t := redblacktree.NewWithStringComparator()
	return &memtableRedBlackTree{
		t: t,
	}
}

type memtableRedBlackTree struct {
	t *redblacktree.Tree
}

func (m *memtableRedBlackTree) Get(key string) (keyvalue.Value, bool) {
	value, found := m.t.Get(string(key))
	if !found {
		return keyvalue.Value{}, false
	}

	return value.(keyvalue.Value), true
}

func (m *memtableRedBlackTree) Delete(key string) {
	m.t.Put(string(key), keyvalue.Value{
		Deleted: true,
	})
}

func (m *memtableRedBlackTree) Insert(key, value string) {
	m.t.Put(string(key), keyvalue.Value{
		Value: value,
	})
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

func (i *memtableRedBlackTreeIterator) Key() string {
	return i.iter.Key().(string)
}

func (i *memtableRedBlackTreeIterator) Value() keyvalue.Value {
	return i.iter.Value().(keyvalue.Value)
}

func (m *memtableRedBlackTree) Iterate() keyvalue.KeyValueIterator {
	iter := m.t.Iterator()
	return &memtableRedBlackTreeIterator{
		iter: &iter,
	}
}
