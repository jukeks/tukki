package sstable_test

import (
	"os"
	"testing"

	"github.com/jukeks/tukki/lib/testhelpers"
	"github.com/jukeks/tukki/memtable"
	"github.com/jukeks/tukki/sstable"
	"github.com/thanhpk/randstr"
)

func TestSSTable(t *testing.T) {
	mt := memtable.NewMemtable()

	len := 10000
	keys := make([]memtable.KeyType, len)
	values := make([]string, len)
	for i := 0; i < len; i++ {
		keys[i] = memtable.KeyType(randstr.String(16))
		values[i] = randstr.String(16)
		mt.Insert(keys[i], values[i])
	}

	f := testhelpers.CreateTempFile("test-tukki", "sstable-test-*")
	tmpfile := f.Name()
	defer os.Remove(tmpfile)

	sstw := sstable.NewSSTableWriter(f)
	written, err := sstw.Write(mt.Iterate())
	if err != nil {
		t.Fatal(err)
	}

	if written != len {
		t.Fatalf("not all keys written: %d", written)
	}

	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	f, err = os.Open(tmpfile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	sstr := sstable.NewSSTableReader(f)
	ssti, err := sstr.Read()
	if err != nil {
		t.Fatal(err)
	}

	found := 0
	for ssti.Next() {
		found++
		key := ssti.Key()
		value := ssti.Value()

		expectedValue, found := mt.Get(memtable.KeyType(key))
		if !found {
			t.Fatalf("key %s not found", key)
		}

		if value != expectedValue {
			t.Fatalf("value for key %s does not match", key)
		}
	}

	if found != len {
		t.Fatalf("not all keys found: %d", found)
	}
}
