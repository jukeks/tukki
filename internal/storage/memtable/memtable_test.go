package memtable_test

import (
	"testing"

	"github.com/jukeks/tukki/internal/storage/memtable"
	"github.com/thanhpk/randstr"
)

func TestMemtable(t *testing.T) {
	mt := memtable.NewMemtable()

	len := 10000
	keys := make([]string, len)
	values := make([]string, len)
	for i := 0; i < len; i++ {
		keys[i] = randstr.String(16)
		values[i] = randstr.String(16)
		mt.Insert(keys[i], values[i])
	}

	for i := 0; i < len; i++ {
		key := keys[i]
		expected := values[i]

		value, found := mt.Get(key)
		if !found {
			t.Errorf("%v not found", key)
		}

		if value.Value != expected {
			t.Errorf("%s was expect but %s was found", expected, value.Value)
		}
	}

	for i := 0; i < len; i++ {
		key := keys[i]
		mt.Delete(key)
		value, found := mt.Get(key)
		if found && !value.Deleted {
			t.Errorf("%v found even though deleted", key)
		}
	}
}

func TestMemtableIterator(t *testing.T) {
	mt := memtable.NewMemtable()

	len := 10000
	keys := make([]string, len)
	values := make([]string, len)
	for i := 0; i < len; i++ {
		keys[i] = randstr.String(16)
		values[i] = randstr.String(16)
		mt.Insert(keys[i], values[i])
	}

	iter := mt.Iterate()
	lastKey := ""
	for entry, err := iter.Next(); err == nil; entry, err = iter.Next() {
		if lastKey != "" && entry.Key < lastKey {
			t.Errorf("iterator not sorted")
		}
		lastKey = entry.Key
	}

}
