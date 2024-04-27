package segments

import (
	"log"
	"os"

	"github.com/jukeks/tukki/internal/sstable"
	"github.com/jukeks/tukki/internal/storage"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type MergeSegmentsOperation struct {
	id              OperationId
	dbDir           string
	segmentsToMerge []Segment
	mergedSegment   Segment
}

func NewMergeSegmentsOperation(dbDir string, segmentsToMerge []Segment, mergedSegment Segment) *MergeSegmentsOperation {
	return &MergeSegmentsOperation{
		dbDir:           dbDir,
		segmentsToMerge: segmentsToMerge,
		mergedSegment:   mergedSegment,
	}
}

func (o *MergeSegmentsOperation) Id() OperationId {
	return o.id
}

func (o *MergeSegmentsOperation) StartJournalEntry() *segmentsv1.SegmentOperationJournalEntry {
	mergeOperation := &segmentsv1.MergeSegments{
		NewSegment: &segmentsv1.Segment{
			Id:       uint64(o.mergedSegment.Id),
			Filename: string(o.mergedSegment.Filename),
		},
	}

	for _, segment := range o.segmentsToMerge {
		mergeOperation.SegmentsToMerge = append(
			mergeOperation.SegmentsToMerge,
			&segmentsv1.Segment{
				Id:       uint64(segment.Id),
				Filename: string(segment.Filename),
			})
	}

	entry := &segmentsv1.SegmentOperationJournalEntry{
		Entry: &segmentsv1.SegmentOperationJournalEntry_Started{
			Started: &segmentsv1.SegmentOperation{
				Id: uint64(o.id),
				Operation: &segmentsv1.SegmentOperation_Merge{
					Merge: mergeOperation,
				},
			},
		},
	}

	return entry
}

func (o *MergeSegmentsOperation) CompletedJournalEntry() *segmentsv1.SegmentOperationJournalEntry {
	entry := &segmentsv1.SegmentOperationJournalEntry{
		Entry: &segmentsv1.SegmentOperationJournalEntry_Completed{
			Completed: uint64(o.id),
		},
	}

	return entry
}

func (o *MergeSegmentsOperation) Execute() error {
	mergedPath := storage.GetPath(o.dbDir, o.mergedSegment.Filename)
	mergedFile, err := os.Create(mergedPath)
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return err
	}

	aPath := storage.GetPath(o.dbDir, o.segmentsToMerge[0].Filename)
	aFile, err := os.Open(aPath)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return err
	}
	defer aFile.Close()
	aReader := sstable.NewSSTableReader(aFile)

	bPath := storage.GetPath(o.dbDir, o.segmentsToMerge[1].Filename)
	bFile, err := os.Open(bPath)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return err
	}
	defer bFile.Close()
	bReader := sstable.NewSSTableReader(bFile)

	err = sstable.MergeSSTables(mergedFile, aReader, bReader)
	if err != nil {
		log.Printf("failed to merge sstables: %v", err)
		return err
	}

	return mergedFile.Close()
}
