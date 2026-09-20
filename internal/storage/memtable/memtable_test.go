package memtable_test

import (
	"bytes"
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
		mt.Insert([]byte(keys[i]), []byte(values[i]))
	}

	for i := 0; i < len; i++ {
		key := keys[i]
		expected := values[i]

		value, found := mt.Get([]byte(key))
		if !found {
			t.Errorf("%v not found", key)
		}

		if !bytes.Equal(value.Value, []byte(expected)) {
			t.Errorf("%s was expect but %s was found", expected, value.Value)
		}
	}

	for i := 0; i < len; i++ {
		key := keys[i]
		mt.Delete([]byte(key))
		value, found := mt.Get([]byte(key))
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
		mt.Insert([]byte(keys[i]), []byte(values[i]))
	}

	iter := mt.Iterate()
	lastKey := ""
	for entry, err := iter.Next(); err == nil; entry, err = iter.Next() {
		if lastKey != "" && string(entry.Key) < string(lastKey) {
			t.Errorf("iterator not sorted")
		}
		lastKey = string(entry.Key)
	}

}
