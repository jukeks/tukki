package segments

import (
	"log"
	"os"

	"github.com/jukeks/tukki/internal/segmentmembers"
	"github.com/jukeks/tukki/internal/sstable"
	"github.com/jukeks/tukki/internal/storage"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type MergeSegmentsOperation struct {
	id              OperationId
	dbDir           string
	segmentsToMerge []SegmentMetadata
	mergedSegment   SegmentMetadata
}

func NewMergeSegmentsOperation(id OperationId, dbDir string, segmentsToMerge []SegmentMetadata, mergedSegment SegmentMetadata) *MergeSegmentsOperation {
	return &MergeSegmentsOperation{
		id:              id,
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
		NewSegment: segmentMetadataToPb(&o.mergedSegment),
	}

	for _, segment := range o.segmentsToMerge {
		mergeOperation.SegmentsToMerge = append(
			mergeOperation.SegmentsToMerge,
			segmentMetadataToPb(&segment),
		)
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

func getEstimatedElementCount(dbDir string, segments []SegmentMetadata) (uint, error) {
	var size uint
	for _, segment := range segments {
		members, err := segmentmembers.OpenSegmentMembers(dbDir, segment.MembersFile)
		if err != nil {
			return 0, err
		}
		size += members.Size()
	}
	return size, nil
}

func (o *MergeSegmentsOperation) Execute() error {
	mergedPath := storage.GetPath(o.dbDir, o.mergedSegment.SegmentFile)
	mergedFile, err := os.Create(mergedPath)
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return err
	}

	aPath := storage.GetPath(o.dbDir, o.segmentsToMerge[0].SegmentFile)
	aFile, err := os.Open(aPath)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return err
	}
	defer aFile.Close()
	aReader := sstable.NewSSTableReader(aFile)

	bPath := storage.GetPath(o.dbDir, o.segmentsToMerge[1].SegmentFile)
	bFile, err := os.Open(bPath)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return err
	}
	defer bFile.Close()
	bReader := sstable.NewSSTableReader(bFile)

	totalMembers, err := getEstimatedElementCount(o.dbDir, o.segmentsToMerge)
	if err != nil {
		log.Printf("failed to get estimated element count: %v", err)
		return err
	}
	members := segmentmembers.NewSegmentMembers(totalMembers)

	err = sstable.MergeSSTables(mergedFile, aReader, bReader, members)
	if err != nil {
		log.Printf("failed to merge sstables: %v", err)
		return err
	}
	err = members.Save(o.dbDir, o.mergedSegment.MembersFile)
	if err != nil {
		log.Printf("failed to members: %v", err)
		return err
	}

	return mergedFile.Close()
}
