package operation

import (
	"io"

	"github.com/jukeks/tukki/internal/journal"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type SegmentOperationJournal struct {
	journal *journal.Journal
}

func OpenSegmentOperationJournal(dbDir string) (*SegmentOperationJournal, map[uint64]SegmentOperation, error) {
	var segmentOperationsMap map[uint64]SegmentOperation
	handle := func(r *journal.JournalReader) error {
		var err error
		segmentOperationsMap, err = readJournal(r)
		return err
	}

	j, err := journal.OpenJournal(dbDir, "segment_operations.journal", handle)
	if err != nil {
		return nil, nil, err
	}

	return &SegmentOperationJournal{j}, segmentOperationsMap, nil
}

func readJournal(r *journal.JournalReader) (map[uint64]SegmentOperation, error) {
	operations := make(map[uint64]SegmentOperation)
	for {
		journalEntry := &segmentsv1.SegmentOperationJournalEntry{}
		err := r.Read(journalEntry)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch journalEntry.Entry.(type) {
		case *segmentsv1.SegmentOperationJournalEntry_Started:
			started := journalEntry.GetStarted()
			operation := segmentOperationFromProto(started)
			operations[operation.Id()] = operation
		case *segmentsv1.SegmentOperationJournalEntry_Completed:
			completedId := journalEntry.GetCompleted()
			delete(operations, completedId)
		}
	}

	return operations, nil
}

func segmentOperationFromProto(proto *segmentsv1.SegmentOperation) SegmentOperation {
	switch proto.Operation.(type) {
	case *segmentsv1.SegmentOperation_Add:
		addOperation := proto.GetAdd()
		return &AddSegmentOperation{
			id: proto.Id,
			segment: Segment{
				Id:       addOperation.Segment.Id,
				Filename: addOperation.Segment.Filename,
			},
		}
	case *segmentsv1.SegmentOperation_Merge:
		mergeOperation := proto.GetMerge()
		segmentsToMerge := make([]Segment, len(mergeOperation.SegmentsToMerge))
		for i, segmentProto := range mergeOperation.SegmentsToMerge {
			segmentsToMerge[i] = Segment{
				Id:       segmentProto.Id,
				Filename: segmentProto.Filename,
			}
		}
		return &MergeSegmentsOperation{
			id:              proto.Id,
			segmentsToMerge: segmentsToMerge,
			mergedSegment: Segment{
				Id:       mergeOperation.NewSegment.Id,
				Filename: mergeOperation.NewSegment.Filename,
			},
		}
	}

	return nil
}
