package memtable_test

import (
	"testing"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/thanhpk/randstr"
)

func TestMemtable(t *testing.T) {
	mt := memtable.NewMemtable()

	len := 10000
	keys := make([]memtable.KeyType, len)
	values := make([]string, len)
	for i := 0; i < len; i++ {
		keys[i] = memtable.KeyType(randstr.String(16))
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

		if value != expected {
			t.Errorf("%s was expect but %s was found", expected, value)
		}
	}

	for i := 0; i < len; i++ {
		key := keys[i]
		mt.Delete(key)
		_, found := mt.Get(key)
		if found {
			t.Errorf("%v found even though deleted", key)
		}
	}
}
