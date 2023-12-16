package sstable

import (
	"io"

	"github.com/jukeks/tukki/internal/keyvalue"
)

func MergeSSTables(w io.Writer, a, b keyvalue.KeyValueIterator) error {
	writer := NewSSTableWriter(w)

	for {
		entryA, errA := a.Next()
		entryB, errB := b.Next()

		if errA == io.EOF && errB == io.EOF {
			break
		}

		if errA == io.EOF {
			writer.Write(entryB)
			continue
		}

		if errB == io.EOF {
			writer.Write(entryA)
			continue
		}

		if entryA.Key < entryB.Key {
			writer.Write(entryA)
			continue
		}

		if entryA.Key > entryB.Key {
			writer.Write(entryB)
			continue
		}

		if entryA.Deleted {
			writer.Write(entryA)
			continue
		}

		writer.Write(entryB)
	}

	return nil
}
