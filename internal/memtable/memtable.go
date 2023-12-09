package memtable

import (
	"github.com/emirpasic/gods/trees/redblacktree"
)

type KeyType string

type KeyValueIterator interface {
	Next() bool
	Key() string
	Value() string
}

type Memtable interface {
	Get(key KeyType) (string, bool)
	Insert(key KeyType, value string)
	Delete(key KeyType)
	Iterate() KeyValueIterator
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

func (m *memtableRedBlackTree) Get(key KeyType) (string, bool) {
	value, found := m.t.Get(string(key))
	if !found {
		return "", false
	}

	return value.(string), true
}

func (m *memtableRedBlackTree) Delete(key KeyType) {
	m.t.Remove(string(key))
}

func (m *memtableRedBlackTree) Insert(key KeyType, value string) {
	m.t.Put(string(key), value)
}

type memtableRedBlackTreeIterator struct {
	iter *redblacktree.Iterator
}

func (i *memtableRedBlackTreeIterator) Next() bool {
	return i.iter.Next()
}

func (i *memtableRedBlackTreeIterator) Key() string {
	return i.iter.Key().(string)
}

func (i *memtableRedBlackTreeIterator) Value() string {
	return i.iter.Value().(string)
}

func (m *memtableRedBlackTree) Iterate() KeyValueIterator {
	iter := m.t.Iterator()
	return &memtableRedBlackTreeIterator{
		iter: &iter,
	}
}
