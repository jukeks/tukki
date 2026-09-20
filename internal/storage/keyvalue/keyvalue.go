package keyvalue

type Value struct {
	Value   []byte
	Deleted bool
}

type IteratorEntry struct {
	Key     []byte
	Value   []byte
	Deleted bool
}

type KeyValueIterator interface {
	Next() (IteratorEntry, error)
}
