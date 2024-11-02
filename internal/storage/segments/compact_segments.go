package segments

import (
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/index"
	"github.com/jukeks/tukki/internal/storage/keyvalue"
	"github.com/jukeks/tukki/internal/storage/segmentmembers"
	"github.com/jukeks/tukki/internal/storage/sstable"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type CompactSegmentsOperation struct {
	id                OperationId
	dbDir             string
	segmentsToCompact []SegmentMetadata
	newSegments       []SegmentMetadata
	targetSegmentSize uint64
}

func NewCompactSegmentsOperation(id OperationId, dbDir string, segmentsToCompact []SegmentMetadata, targetSegmentSize uint64) *CompactSegmentsOperation {
	return &CompactSegmentsOperation{
		id:                id,
		dbDir:             dbDir,
		segmentsToCompact: segmentsToCompact,
		targetSegmentSize: targetSegmentSize,
	}
}

func (o *CompactSegmentsOperation) Id() OperationId {
	return o.id
}

func (o *CompactSegmentsOperation) SegmentsToCompact() []SegmentMetadata {
	return o.segmentsToCompact
}

func (o *CompactSegmentsOperation) NewSegments() []SegmentMetadata {
	return o.newSegments
}

func (o *CompactSegmentsOperation) StartJournalEntry() *segmentsv1.SegmentOperationJournalEntry {
	compactOperation := &segmentsv1.CompactSegments{
		TargetSegmentSize: o.targetSegmentSize,
	}
	for _, segment := range o.segmentsToCompact {
		compactOperation.SegmentsToCompact = append(
			compactOperation.SegmentsToCompact,
			segmentMetadataToPb(&segment),
		)
	}

	entry := &segmentsv1.SegmentOperationJournalEntry{
		Entry: &segmentsv1.SegmentOperationJournalEntry_Started{
			Started: &segmentsv1.SegmentOperation{
				Id: uint64(o.id),
				Operation: &segmentsv1.SegmentOperation_Compact{
					Compact: compactOperation,
				},
			},
		},
	}

	return entry
}

func (o *CompactSegmentsOperation) CompletedJournalEntry() *segmentsv1.SegmentOperationJournalEntry {
	added := make([]*segmentsv1.Segment, len(o.newSegments))
	for i, segment := range o.newSegments {
		added[i] = segmentMetadataToPb(&segment)
	}
	freed := make([]*segmentsv1.Segment, len(o.segmentsToCompact))
	for i, segment := range o.segmentsToCompact {
		freed[i] = segmentMetadataToPb(&segment)
	}

	entry := &segmentsv1.SegmentOperationJournalEntry{
		Entry: &segmentsv1.SegmentOperationJournalEntry_CompletedV2{
			CompletedV2: &segmentsv1.SegmentOperationCompleted{
				Id:    uint64(o.id),
				Added: added,
				Freed: freed,
			},
		},
	}

	return entry
}

func (o *CompactSegmentsOperation) getNewSegmentIds() []SegmentId {
	ids := make([]SegmentId, len(o.segmentsToCompact))
	for i, segment := range o.segmentsToCompact {
		ids[i] = segment.Id
	}

	// reverse
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	return ids
}

func (o *CompactSegmentsOperation) Execute() error {
	openedSegments := make([]keyvalue.SubIterator, 0)

	for _, segment := range o.segmentsToCompact {
		segmentFile, err := files.OpenFile(o.dbDir, segment.SegmentFile)
		if err != nil {
			return fmt.Errorf("failed to open segment file: %w", err)
		}
		subIter := sstable.NewSSTableIterator(segmentFile, nil)
		openedSegments = append(openedSegments, subIter)
	}

	iterator, err := keyvalue.NewIterator("", "", true, openedSegments...)
	if err != nil {
		return fmt.Errorf("failed to create iterator: %w", err)
	}
	defer iterator.Close()

	usedSegments := 0
	newSegmentIds := o.getNewSegmentIds()
	sequence := o.Id()
	for {
		newId := newSegmentIds[usedSegments]
		usedSegments++
		newSegment := SegmentMetadata{
			Id:          newId,
			SegmentFile: GetCompactedSegmentFilename(newId, sequence),
			MembersFile: GetCompactedMembersFilename(newId, sequence),
			IndexFile:   GetCompactedIndexFilename(newId, sequence),
		}

		o.newSegments = append(o.newSegments, newSegment)

		segmentFile, err := files.CreateFile(o.dbDir, newSegment.SegmentFile)
		if err != nil {
			return fmt.Errorf("failed to create segment file: %w", err)
		}

		log.Printf("Compacting to %s", newSegment.SegmentFile)

		done := false
		writer := sstable.NewSSTableWriter(segmentFile)
		err = writer.WriteFromIteratorUntil(iterator, o.targetSegmentSize)
		if err != nil {
			if err == io.EOF {
				done = true
			} else {
				segmentFile.Close()
				return fmt.Errorf("failed to write segment: %w", err)
			}
		}
		if err := segmentFile.Close(); err != nil {
			return fmt.Errorf("failed to close segment file: %w", err)
		}

		offsets := writer.WrittenOffsets()

		members := segmentmembers.NewSegmentMembers(uint(len(offsets)))
		for key := range offsets {
			log.Printf("Adding key %s to segment %d", key, newId)
			members.Add(key)
		}
		err = members.Save(o.dbDir, newSegment.MembersFile)
		if err != nil {
			return fmt.Errorf("failed to save members: %w", err)
		}

		indexFile, err := files.CreateFile(o.dbDir, newSegment.IndexFile)
		if err != nil {
			return fmt.Errorf("failed to create index file: %w", err)
		}
		indexWriter := index.NewIndexWriter(indexFile)
		if err := indexWriter.WriteFromOffsets(offsets); err != nil {
			indexFile.Close()
			return fmt.Errorf("failed to write index: %w", err)
		}
		if err := indexWriter.Close(); err != nil {
			return fmt.Errorf("failed to close index writer: %w", err)
		}

		if done {
			break
		}
	}

	return nil
}
