package memtable

import "github.com/jukeks/tukki/internal/storage/keyvalue"

type memtableSubIterator struct {
	memtable Memtable
	iterator keyvalue.KeyValueIterator
	current  keyvalue.IteratorEntry
	err      error
}

func NewMemtableIterator(memtable Memtable) keyvalue.SubIterator {
	iter := &memtableSubIterator{
		memtable: memtable,
		iterator: memtable.Iterate(),
	}
	iter.Progress()

	return iter
}

func (m *memtableSubIterator) Close() error {
	return nil
}

func (m *memtableSubIterator) Get() (keyvalue.IteratorEntry, error) {
	return m.current, m.err
}

func (m *memtableSubIterator) Progress() {
	m.current, m.err = m.iterator.Next()
}

func (m *memtableSubIterator) Seek(key string) error {
	m.iterator = m.memtable.Iterate()
	for {
		m.Progress()
		if m.err != nil {
			return m.err
		}

		if m.current.Key >= key {
			return nil
		}
	}
}
