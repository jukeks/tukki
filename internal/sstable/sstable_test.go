package sstable_test

import (
	"os"
	"testing"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/sstable"
	testutil "github.com/jukeks/tukki/testutil"
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

	tmpDir := t.TempDir()
	f := testutil.CreateTempFile(tmpDir, "sstable-test-*")
	tmpfile := f.Name()

	sstw := sstable.NewSSTableWriter(f)
	err := sstw.WriteFromIterator(mt.Iterate())
	if err != nil {
		t.Fatal(err)
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
	found := 0
	for entry, err := sstr.Next(); err == nil; entry, err = sstr.Next() {
		found++

		key := entry.Key
		value := entry.Value
		expectedValue, found := mt.Get(entry.Key)
		if !found {
			t.Fatalf("key %s not found", key)
		}

		if value != expectedValue.Value {
			t.Fatalf("value for key %s does not match", key)
		}
	}

	if found != len {
		t.Fatalf("not all keys found: %d", found)
	}
}
