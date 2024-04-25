package operation

import (
	"log"
	"os"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/sstable"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type AddSegmentOperation struct {
	id       uint64
	segment  Segment
	memtable memtable.Memtable
}

func NewAddSegmentOperation(segment Segment, memtable memtable.Memtable) *AddSegmentOperation {
	return &AddSegmentOperation{
		segment:  segment,
		memtable: memtable,
	}
}

func (o *AddSegmentOperation) Id() uint64 {
	return o.id
}

func (o *AddSegmentOperation) StartJournalEntry() *segmentsv1.SegmentOperationJournalEntry {
	entry := &segmentsv1.SegmentOperationJournalEntry{
		Entry: &segmentsv1.SegmentOperationJournalEntry_Started{
			Started: &segmentsv1.SegmentOperation{
				Id: o.id,
				Operation: &segmentsv1.SegmentOperation_Add{
					Add: &segmentsv1.AddSegment{
						Segment: &segmentsv1.Segment{
							Id:       o.segment.Id,
							Filename: o.segment.Filename,
						},
					},
				},
			},
		},
	}
	return entry
}

func (o *AddSegmentOperation) CompletedJournalEntry() *segmentsv1.SegmentOperationJournalEntry {
	entry := &segmentsv1.SegmentOperationJournalEntry{
		Entry: &segmentsv1.SegmentOperationJournalEntry_Completed{
			Completed: o.id,
		},
	}
	return entry
}

func (o *AddSegmentOperation) Execute() error {
	f, err := os.Create(o.segment.Filename)
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return err
	}

	w := sstable.NewSSTableWriter(f)
	err = w.WriteFromIterator(o.memtable.Iterate())
	if err != nil {
		log.Printf("failed to write sstable from memtable: %v", err)
		// best effort cleanup
		f.Close()
		os.Remove(o.segment.Filename)
		return err
	}

	return f.Close()
}
