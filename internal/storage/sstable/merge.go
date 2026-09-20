package sstable

import (
	"bytes"
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
			if _, err := writer.Write(entryB); err != nil {
				return nil, err
			}
			members.Add(entryB.Key)
			entryB, errB = b.Next()
			continue
		}

		// b is completely read
		if errB == io.EOF {
			if _, err := writer.Write(entryA); err != nil {
				return nil, err
			}
			members.Add(entryA.Key)
			entryA, errA = a.Next()
			continue
		}

		// merge sorted entries by key
		if bytes.Compare(entryA.Key, entryB.Key) < 0 {
			if _, err := writer.Write(entryA); err != nil {
				return nil, err
			}
			members.Add(entryA.Key)
			entryA, errA = a.Next()
			continue
		}
		if bytes.Compare(entryA.Key, entryB.Key) > 0 {
			if _, err := writer.Write(entryB); err != nil {
				return nil, err
			}
			members.Add(entryB.Key)
			entryB, errB = b.Next()
			continue
		}

		if bytes.Equal(entryA.Key, entryB.Key) {
			// b is newer segment
			if _, err := writer.Write(entryB); err != nil {
				return nil, err
			}
			members.Add(entryB.Key)

			entryA, errA = a.Next()
			entryB, errB = b.Next()
			continue
		}
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}

	return writer.WrittenOffsets(), nil
}
