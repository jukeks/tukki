package sstable

import (
	"io"

	"github.com/jukeks/tukki/internal/storage/index"
	"github.com/jukeks/tukki/internal/storage/keyvalue"
	"github.com/jukeks/tukki/internal/storage/segmentmembers"
)

func MergeSSTables(sstableWriter io.Writer, a, b keyvalue.KeyValueIterator,
	members *segmentmembers.SegmentMembers) (index.OffsetMap, error) {

	writer := NewSSTableWriter(sstableWriter)

	entryA, errA := a.Next()
	entryB, errB := b.Next()
	for {
		if errA == io.EOF && errB == io.EOF {
			break
		}
		if errA != nil && errA != io.EOF {
			return nil, errA
		}
		if errB != nil && errB != io.EOF {
			return nil, errB
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

	return writer.WrittenOffsets(), nil
}
