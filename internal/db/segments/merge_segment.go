package segments

import (
	"fmt"

	"github.com/jukeks/tukki/internal/storage"
	"github.com/jukeks/tukki/internal/storage/index"
	"github.com/jukeks/tukki/internal/storage/segmentmembers"
	"github.com/jukeks/tukki/internal/storage/sstable"
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
	mergedFile, err := storage.CreateFile(o.dbDir, o.mergedSegment.SegmentFile)
	if err != nil {
		return fmt.Errorf("failed to create merged segment file: %w", err)
	}

	aFile, err := storage.OpenFile(o.dbDir, o.segmentsToMerge[0].SegmentFile)
	if err != nil {
		return fmt.Errorf("failed to open a segment file: %w", err)
	}
	defer aFile.Close()
	aReader := sstable.NewSSTableReader(aFile)

	bFile, err := storage.OpenFile(o.dbDir, o.segmentsToMerge[1].SegmentFile)
	if err != nil {
		return fmt.Errorf("failed to open b segment file: %w", err)
	}
	defer bFile.Close()
	bReader := sstable.NewSSTableReader(bFile)

	totalMembers, err := getEstimatedElementCount(o.dbDir, o.segmentsToMerge)
	if err != nil {
		return fmt.Errorf("failed to get estimated element count: %w", err)
	}
	members := segmentmembers.NewSegmentMembers(totalMembers)

	indexFile, err := storage.CreateFile(o.dbDir, o.mergedSegment.IndexFile)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	indexWriter := index.NewIndexWriter(indexFile)

	offsets, err := sstable.MergeSSTables(mergedFile, aReader, bReader, members)
	if err != nil {
		return fmt.Errorf("failed to merge sstables: %w", err)
	}
	err = members.Save(o.dbDir, o.mergedSegment.MembersFile)
	if err != nil {
		return fmt.Errorf("failed to save members: %w", err)
	}
	if err := indexWriter.WriteFromOffsets(offsets); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}
	err = indexWriter.Close()
	if err != nil {
		return fmt.Errorf("failed to close index writer: %w", err)
	}

	return mergedFile.Close()
}
