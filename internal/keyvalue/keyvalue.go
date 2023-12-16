package keyvalue

type Value struct {
	Value   string
	Deleted bool
}

type IteratorEntry struct {
	Key     string
	Value   string
	Deleted bool
}

type KeyValueIterator interface {
	Next() (IteratorEntry, error)
}
