package segments

import (
	"fmt"
	"log"
	"os"

	"github.com/jukeks/tukki/internal/index"
	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/segmentmembers"
	"github.com/jukeks/tukki/internal/sstable"
	"github.com/jukeks/tukki/internal/storage"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type OpenSegment struct {
	Segment     SegmentMetadata
	WalFilename storage.Filename
	Memtable    memtable.Memtable
}

type AddSegmentOperation struct {
	id                OperationId
	dbDir             string
	completingSegment *OpenSegment
	newSegment        *OpenSegment
}

func NewAddSegmentOperation(
	id OperationId,
	dbDir string,
	completingSegment *OpenSegment,
	newSegment *OpenSegment) *AddSegmentOperation {

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
			Segment:     segmentMetadataToPb(&o.completingSegment.Segment),
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
							Segment:     segmentMetadataToPb(&o.newSegment.Segment),
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
		path := storage.GetPath(o.dbDir, completingSegment.Segment.SegmentFile)
		f, err := os.Create(path)
		if err != nil {
			log.Printf("failed to create file: %v", err)
			return err
		}

		indexPath := storage.GetPath(o.dbDir, completingSegment.Segment.IndexFile)
		indexF, err := os.Create(indexPath)
		if err != nil {
			return fmt.Errorf("failed to create index file: %w", err)
		}
		indexWriter := index.NewIndexWriter(indexF)
		// TODO USE INDEX WRITER

		w := sstable.NewSSTableWriter(f)
		err = w.WriteFromIterator(completingSegment.Memtable.Iterate())
		if err != nil {
			return fmt.Errorf("failed to write sstable from memtable: %w", err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}

		if err := indexWriter.Close(); err != nil {
			return fmt.Errorf("failed to close index file: %w", err)
		}

		// remove completing wal
		path = storage.GetPath(o.dbDir, completingSegment.WalFilename)
		err = os.Remove(path)
		if err != nil {
			log.Printf("failed to remove file: %v", err)
			return err
		}

		// populate segment members
		members := segmentmembers.NewSegmentMembers(uint(completingSegment.Memtable.Size()))
		members.Fill(completingSegment.Memtable.Iterate())
		err = members.Save(o.dbDir, completingSegment.Segment.MembersFile)
		if err != nil {
			log.Printf("failed to save segment members: %v", err)
			return err
		}
	}

	return nil
}
