package memtable

import (
	"github.com/emirpasic/gods/trees/redblacktree"
)

type KeyType int

type Memtable interface {
	Get(key KeyType) (string, bool)
	Insert(key KeyType, value string)
	Delete(key KeyType)
}

func NewMemtable() Memtable {
	t := redblacktree.NewWithIntComparator()
	return &memtableRedBlackTree{
		t: t,
	}
}

type memtableRedBlackTree struct {
	t *redblacktree.Tree
}

func (m *memtableRedBlackTree) Get(key KeyType) (string, bool) {
	value, found := m.t.Get(int(key))
	if !found {
		return "", false
	}

	return value.(string), true
}

func (m *memtableRedBlackTree) Delete(key KeyType) {
	m.t.Remove(int(key))
}

func (m *memtableRedBlackTree) Insert(key KeyType, value string) {
	intKey := int(key)
	m.t.Put(intKey, value)
}

func (m *memtableRedBlackTree) Flush() {
	for iter := m.t.Iterator(); iter.Next(); {
		key := iter.Key()
		value := iter.Value()
	}
}
