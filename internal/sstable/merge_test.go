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

func readToMemtable(filenames ...string) memtable.Memtable {
	mt := memtable.NewMemtable()
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		sstr := sstable.NewSSTableReader(f)
		ssti := sstr.Iterate()

		for entry, err := ssti.Next(); err == nil; entry, err = ssti.Next() {
			mt.Insert(entry.Key, entry.Value)
		}
	}

	return mt
}

func TestMerge(t *testing.T) {
	table1Len := 1200
	table2Len := 3000
	testutil.EnsureTempDirectory("test-tukki")
	filename1 := createSSTable(table1Len)
	defer os.Remove(filename1)
	filename2 := createSSTable(table2Len)
	defer os.Remove(filename2)

	mt := readToMemtable(filename1, filename2)

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
	iter1 := reader1.Iterate()
	iter2 := reader2.Iterate()

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
	ssti := sstr.Iterate()

	found := 0
	prevKey := ""
	for entry, err := ssti.Next(); err == nil; entry, err = ssti.Next() {
		found++

		if prevKey != "" && prevKey > entry.Key {
			t.Fatalf("keys not sorted")
		}
		prevKey = entry.Key
	}

	if found != table1Len+table2Len {
		t.Fatalf("not all keys found: %d", found)
	}

	mt2 := readToMemtable(outfile)
	mt1Iter := mt.Iterate()
	mt2Iter := mt2.Iterate()

	for i := 0; i < table1Len+table2Len; i++ {
		entry1, err := mt1Iter.Next()
		if err != nil {
			panic(err)
		}
		entry2, err := mt2Iter.Next()
		if err != nil {
			panic(err)
		}

		if entry1.Key != entry2.Key {
			t.Fatalf("keys not equal")
		}
		if entry1.Value != entry2.Value {
			t.Fatalf("values not equal")
		}
	}

	_, err = mt1Iter.Next()
	if err == nil {
		t.Fatalf("expected EOF")
	}

	_, err = mt2Iter.Next()
	if err == nil {
		t.Fatalf("expected EOF")
	}
}
