package sstable_test

import (
	"os"
	"testing"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/sstable"
	testutil "github.com/jukeks/tukki/tests/util"
	"github.com/thanhpk/randstr"
)

func createSSTable(length int) string {
	mt := memtable.NewMemtable()
	keys := make([]string, length)
	values := make([]string, length)
	for i := 0; i < length; i++ {
		keys[i] = randstr.String(16)
		values[i] = randstr.String(16)
		mt.Insert(keys[i], values[i])
	}

	f := testutil.CreateTempFile("test-tukki", "sstable-test-*")
	tmpfile := f.Name()

	sstw := sstable.NewSSTableWriter(f)
	err := sstw.WriteFromIterator(mt.Iterate())
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	return tmpfile
}

func TestMerge(t *testing.T) {
	testutil.EnsureTempDirectory("test-tukki")
	filename1 := createSSTable(1200)
	defer os.Remove(filename1)
	filename2 := createSSTable(3000)
	defer os.Remove(filename2)

	f1, err := os.Open(filename1)
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	f2, err := os.Open(filename2)
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	f := testutil.CreateTempFile("test-tukki", "sstable-test-*")
	outfile := f.Name()
	defer f.Close()
	defer os.Remove(outfile)

	reader1 := sstable.NewSSTableReader(f1)
	reader2 := sstable.NewSSTableReader(f2)
	iter1, err := reader1.Read()
	if err != nil {
		panic(err)
	}
	iter2, err := reader2.Read()
	if err != nil {
		panic(err)
	}

	sstable.MergeSSTables(f, iter1, iter2)
	err = f1.Close()
	if err != nil {
		panic(err)
	}
	f2.Close()
	if err != nil {
		panic(err)
	}

	f, err = os.Open(outfile)
	if err != nil {
		panic(err)
	}

	sstr := sstable.NewSSTableReader(f)
	ssti, err := sstr.Read()
	if err != nil {
		panic(err)
	}

	found := 0
	prevKey := ""
	for entry, err := ssti.Next(); err == nil; entry, err = ssti.Next() {
		found++

		if prevKey != "" && prevKey > entry.Key {
			t.Fatalf("keys not sorted")
		}
		prevKey = entry.Key
	}

	if found != 4200 {
		t.Fatalf("not all keys found: %d", found)
	}
}
