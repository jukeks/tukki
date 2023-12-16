package keyvalue

type Value struct {
	Deleted bool
	Value   string
}

type KeyValueIterator interface {
	Next() bool
	Key() string
	Value() Value
}
