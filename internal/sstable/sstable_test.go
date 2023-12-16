package sstable_test

import (
	"os"
	"testing"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/sstable"
	testutil "github.com/jukeks/tukki/tests/util"
	"github.com/thanhpk/randstr"
)

func TestSSTable(t *testing.T) {
	mt := memtable.NewMemtable()

	len := 10000
	keys := make([]string, len)
	values := make([]string, len)
	for i := 0; i < len; i++ {
		keys[i] = randstr.String(16)
		values[i] = randstr.String(16)
		mt.Insert(keys[i], values[i])
	}

	f := testutil.CreateTempFile("test-tukki", "sstable-test-*")
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

		expectedValue, found := mt.Get(key)
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
