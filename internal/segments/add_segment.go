package segments

import (
	"log"
	"os"

	"github.com/jukeks/tukki/internal/sstable"
	"github.com/jukeks/tukki/internal/storage"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type AddSegmentOperation struct {
	id                OperationId
	dbDir             string
	completingSegment *LiveSegment
	newSegment        *LiveSegment
}

func NewAddSegmentOperation(
	id OperationId,
	dbDir string,
	completingSegment *LiveSegment,
	newSegment *LiveSegment) *AddSegmentOperation {

	return &AddSegmentOperation{
		id:                id,
		dbDir:             dbDir,
		completingSegment: completingSegment,
		newSegment:        newSegment,
	}
}

func (o *AddSegmentOperation) Id() OperationId {
	return o.id
}

func (o *AddSegmentOperation) StartJournalEntry() *segmentsv1.SegmentOperationJournalEntry {
	var completingSegment *segmentsv1.LiveSegment
	if o.completingSegment != nil {
		completingSegment = &segmentsv1.LiveSegment{
			WalFilename: string(o.completingSegment.WalFilename),
			Segment: &segmentsv1.Segment{
				Id:       uint64(o.completingSegment.Segment.Id),
				Filename: string(o.completingSegment.Segment.Filename),
			},
		}
	}

	entry := &segmentsv1.SegmentOperationJournalEntry{
		Entry: &segmentsv1.SegmentOperationJournalEntry_Started{
			Started: &segmentsv1.SegmentOperation{
				Id: uint64(o.id),
				Operation: &segmentsv1.SegmentOperation_Add{
					Add: &segmentsv1.AddSegment{
						CompletingSegment: completingSegment,
						NewSegment: &segmentsv1.LiveSegment{
							WalFilename: string(o.newSegment.WalFilename),
							Segment: &segmentsv1.Segment{
								Id:       uint64(o.newSegment.Segment.Id),
								Filename: string(o.newSegment.Segment.Filename),
							},
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
	if o.completingSegment != nil {
		completingSegment := o.completingSegment
		// write completing segment to disk
		path := storage.GetPath(o.dbDir, completingSegment.Segment.Filename)
		f, err := os.Create(path)
		if err != nil {
			log.Printf("failed to create file: %v", err)
			return err
		}

		w := sstable.NewSSTableWriter(f)
		err = w.WriteFromIterator(completingSegment.Memtable.Iterate())
		if err != nil {
			log.Printf("failed to write sstable from memtable: %v", err)
			return err
		}

		err = f.Close()
		if err != nil {
			log.Printf("failed to close file: %v", err)
			return err
		}

		// remove completing wal
		path = storage.GetPath(o.dbDir, completingSegment.WalFilename)
		err = os.Remove(path)
		if err != nil {
			log.Printf("failed to remove file: %v", err)
			return err
		}
	}

	return nil
}
