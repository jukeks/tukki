package segments

import (
	"log"
	"os"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/sstable"
	"github.com/jukeks/tukki/internal/storage"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type AddSegmentOperation struct {
	id       OperationId
	dbDir    string
	segment  Segment
	memtable memtable.Memtable
}

func NewAddSegmentOperation(dbDir string, segment Segment, memtable memtable.Memtable) *AddSegmentOperation {
	return &AddSegmentOperation{
		dbDir:    dbDir,
		segment:  segment,
		memtable: memtable,
	}
}

func (o *AddSegmentOperation) Id() OperationId {
	return o.id
}

func (o *AddSegmentOperation) StartJournalEntry() *segmentsv1.SegmentOperationJournalEntry {
	entry := &segmentsv1.SegmentOperationJournalEntry{
		Entry: &segmentsv1.SegmentOperationJournalEntry_Started{
			Started: &segmentsv1.SegmentOperation{
				Id: uint64(o.id),
				Operation: &segmentsv1.SegmentOperation_Add{
					Add: &segmentsv1.AddSegment{
						Segment: &segmentsv1.Segment{
							Id:       uint64(o.segment.Id),
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
			Completed: uint64(o.id),
		},
	}
	return entry
}

func (o *AddSegmentOperation) Execute() error {
	path := storage.GetPath(o.dbDir, o.segment.Filename)
	f, err := os.Create(path)
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
		os.Remove(path)
		return err
	}

	return f.Close()
}
