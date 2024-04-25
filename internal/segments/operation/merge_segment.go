package operation

import (
	"log"
	"os"

	"github.com/jukeks/tukki/internal/sstable"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type MergeSegmentsOperation struct {
	id              uint64
	segmentsToMerge []Segment
	mergedSegment   Segment
}

func NewMergeSegmentsOperation(segmentsToMerge []Segment, mergedSegment Segment) *MergeSegmentsOperation {
	return &MergeSegmentsOperation{
		segmentsToMerge: segmentsToMerge,
		mergedSegment:   mergedSegment,
	}
}

func (o *MergeSegmentsOperation) Id() uint64 {
	return o.id
}

func (o *MergeSegmentsOperation) StartJournalEntry() *segmentsv1.SegmentJournalEntry {
	mergeOperation := &segmentsv1.MergeSegments{
		NewSegment: &segmentsv1.Segment{
			Id:       o.mergedSegment.Id,
			Filename: o.mergedSegment.Filename,
		},
	}

	for _, segment := range o.segmentsToMerge {
		mergeOperation.SegmentsToMerge = append(
			mergeOperation.SegmentsToMerge,
			&segmentsv1.Segment{
				Id:       segment.Id,
				Filename: segment.Filename,
			})
	}

	entry := &segmentsv1.SegmentJournalEntry{
		Entry: &segmentsv1.SegmentJournalEntry_Started{
			Started: &segmentsv1.SegmentOperation{
				Id: o.id,
				Operation: &segmentsv1.SegmentOperation_Merge{
					Merge: mergeOperation,
				},
			},
		},
	}

	return entry
}

func (o *MergeSegmentsOperation) CompletedJournalEntry() *segmentsv1.SegmentJournalEntry {
	entry := &segmentsv1.SegmentJournalEntry{
		Entry: &segmentsv1.SegmentJournalEntry_Completed{
			Completed: o.id,
		},
	}

	return entry
}

func (o *MergeSegmentsOperation) Execute() error {
	mergedFile, err := os.Create(o.mergedSegment.Filename)
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return err
	}

	aFile, err := os.Open(o.segmentsToMerge[0].Filename)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return err
	}
	defer aFile.Close()
	aReader := sstable.NewSSTableReader(aFile)

	bFile, err := os.Open(o.segmentsToMerge[1].Filename)
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
