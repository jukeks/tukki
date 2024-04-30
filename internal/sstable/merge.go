package sstable

import (
	"io"

	"github.com/jukeks/tukki/internal/keyvalue"
	"github.com/jukeks/tukki/internal/segmentmembers"
)

func MergeSSTables(w io.Writer, a, b keyvalue.KeyValueIterator,
	members *segmentmembers.SegmentMembers) error {
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
			members.Add(entryB.Key)
			entryB, errB = b.Next()
			continue
		}

		// b is completely read
		if errB == io.EOF {
			writer.Write(entryA)
			members.Add(entryA.Key)
			entryA, errA = a.Next()
			continue
		}

		// merge sorted entries by key
		if entryA.Key < entryB.Key {
			writer.Write(entryA)
			members.Add(entryA.Key)
			entryA, errA = a.Next()
			continue
		}
		if entryA.Key > entryB.Key {
			writer.Write(entryB)
			members.Add(entryB.Key)
			entryB, errB = b.Next()
			continue
		}

		if entryA.Key == entryB.Key {
			// b is newer segment
			writer.Write(entryB)
			members.Add(entryB.Key)

			entryA, errA = a.Next()
			entryB, errB = b.Next()
			continue
		}
	}

	return nil
}
