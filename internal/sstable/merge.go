package sstable

import (
	"io"

	"github.com/jukeks/tukki/internal/keyvalue"
)

func MergeSSTables(w io.Writer, a, b keyvalue.KeyValueIterator) error {
	writer := NewSSTableWriter(w)

	entryA, errA := a.Next()
	entryB, errB := b.Next()
	for {
		if errA == io.EOF && errB == io.EOF {
			break
		}
		if errA != nil && errA != io.EOF {
			return errA
		}
		if errB != nil && errB != io.EOF {
			return errB
		}

		// a is completely read
		if errA == io.EOF {
			writer.Write(entryB)
			entryB, errB = b.Next()
			continue
		}

		// b is completely read
		if errB == io.EOF {
			writer.Write(entryA)
			entryA, errA = a.Next()
			continue
		}

		// merge sorted entries by key
		if entryA.Key < entryB.Key {
			writer.Write(entryA)
			entryA, errA = a.Next()
			continue
		}
		if entryA.Key > entryB.Key {
			writer.Write(entryB)
			entryB, errB = b.Next()
			continue
		}

		if entryA.Key == entryB.Key {
			// b is newer segment
			writer.Write(entryB)

			entryA, errA = a.Next()
			entryB, errB = b.Next()
			continue
		}
	}

	return nil
}
