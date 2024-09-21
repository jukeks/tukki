package sstable_test

import (
	"os"
	"testing"

	"github.com/jukeks/tukki/internal/storage/memtable"
	"github.com/jukeks/tukki/internal/storage/segmentmembers"
	"github.com/jukeks/tukki/internal/storage/sstable"
	testutil "github.com/jukeks/tukki/testutil"
	"github.com/thanhpk/randstr"
)

func createSSTable(dbDir string, length int) string {
	mt := memtable.NewMemtable()
	keys := make([]string, length)
	values := make([]string, length)
	for i := 0; i < length; i++ {
		keys[i] = randstr.String(16)
		values[i] = randstr.String(16)
		mt.Insert(keys[i], values[i])
	}

	f := testutil.CreateTempFile(dbDir, "sstable-test-*")
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
		for entry, err := sstr.Next(); err == nil; entry, err = sstr.Next() {
			mt.Insert(entry.Key, entry.Value)
		}
	}

	return mt
}

func checkMemtableAreEqual(mt1, mt2 memtable.Memtable, expectedLen int) bool {
	mt1Iter := mt1.Iterate()
	mt2Iter := mt2.Iterate()

	for i := 0; i < expectedLen; i++ {
		entry1, err := mt1Iter.Next()
		if err != nil {
			return false
		}
		entry2, err := mt2Iter.Next()
		if err != nil {
			return false
		}

		if entry1.Key != entry2.Key {
			return false
		}
		if entry1.Value != entry2.Value {
			return false
		}
	}

	_, err := mt1Iter.Next()
	if err == nil {
		return false
	}

	_, err = mt2Iter.Next()
	if err == nil {
		return false
	}

	if mt1.MemberCount() != mt2.MemberCount() {
		return false
	}

	return true
}

func testMerge(t *testing.T, table1Len, table2Len int) {
	dbDir := t.TempDir()
	filename1 := createSSTable(dbDir, table1Len)
	defer os.Remove(filename1)
	filename2 := createSSTable(dbDir, table2Len)
	defer os.Remove(filename2)

	mt := readToMemtable(filename1, filename2)

	f1, _ := os.Open(filename1)
	defer f1.Close()

	f2, _ := os.Open(filename2)
	defer f2.Close()

	f := testutil.CreateTempFile(dbDir, "sstable-test-*")
	outfile := f.Name()
	defer f.Close()
	defer os.Remove(outfile)

	reader1 := sstable.NewSSTableReader(f1)
	reader2 := sstable.NewSSTableReader(f2)

	members := segmentmembers.NewSegmentMembers(uint(table1Len + table2Len))
	sstable.MergeSSTables(f, reader1, reader2, members)

	f, _ = os.Open(outfile)

	sstr := sstable.NewSSTableReader(f)

	found := 0
	prevKey := ""
	for entry, err := sstr.Next(); err == nil; entry, err = sstr.Next() {
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
	areEqual := checkMemtableAreEqual(mt, mt2, table1Len+table2Len)
	if !areEqual {
		t.Fatalf("memtables are not equal")
	}
}

func TestMerge(t *testing.T) {
	testMerge(t, 1000, 1000)
	testMerge(t, 1000, 0)
	testMerge(t, 0, 1000)
}

func TestMergeUpdates(t *testing.T) {
	mt1 := memtable.NewMemtable()
	mt1.Insert("a", "a")

	mt2 := memtable.NewMemtable()
	mt2.Insert("a", "b")

	f1 := testutil.CreateTempFile("test-tukki", "sstable-test-*")
	defer os.Remove(f1.Name())
	f2 := testutil.CreateTempFile("test-tukki", "sstable-test-*")
	defer os.Remove(f2.Name())

	sstw1 := sstable.NewSSTableWriter(f1)
	err := sstw1.WriteFromIterator(mt1.Iterate())
	if err != nil {
		t.Fatalf("failed to write to sstable: %v", err)
	}
	sstw2 := sstable.NewSSTableWriter(f2)
	err = sstw2.WriteFromIterator(mt2.Iterate())
	if err != nil {
		t.Fatalf("failed to write to sstable: %v", err)
	}

	f1, _ = os.Open(f1.Name())
	defer f1.Close()
	f2, _ = os.Open(f2.Name())

	f := testutil.CreateTempFile("test-tukki", "sstable-test-*")
	outfile := f.Name()
	defer f.Close()
	defer os.Remove(outfile)

	reader1 := sstable.NewSSTableReader(f1)
	reader2 := sstable.NewSSTableReader(f2)

	members := segmentmembers.NewSegmentMembers(100000)
	sstable.MergeSSTables(f, reader1, reader2, members)

	f, _ = os.Open(outfile)
	mergeReader := sstable.NewSSTableReader(f)

	m := make(map[string]string)
	for entry, err := mergeReader.Next(); err == nil; entry, err = mergeReader.Next() {
		m[entry.Key] = entry.Value
	}

	if len(m) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(m))
	}

	if m["a"] != "b" {
		t.Fatalf("expected value b, got %s", m["a"])
	}
}
