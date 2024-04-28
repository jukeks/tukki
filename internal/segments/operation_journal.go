package segments

import (
	"io"

	"github.com/jukeks/tukki/internal/journal"
	"github.com/jukeks/tukki/internal/storage"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

type SegmentOperationJournal struct {
	journal *journal.Journal
}

func OpenSegmentOperationJournal(dbDir string) (
	*SegmentOperationJournal,
	*CurrentSegments,
	error,
) {
	var currentSegments *CurrentSegments
	handle := func(r *journal.JournalReader) error {
		var err error
		currentSegments, err = readOperationJournal(r)
		return err
	}

	j, err := journal.OpenJournal(dbDir, "segment_operations.journal", handle)
	if err != nil {
		return nil, nil, err
	}

	return &SegmentOperationJournal{j}, currentSegments, nil
}

type CurrentSegments struct {
	Ongoing    *LiveSegment
	Segments   map[SegmentId]Segment
	Operations map[OperationId]SegmentOperation
}

func readOperationJournal(r *journal.JournalReader) (
	*CurrentSegments,
	error) {

	operations := make(map[OperationId]SegmentOperation)
	segments := make(map[SegmentId]Segment)
	var ongoing *LiveSegment

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
			completedId := OperationId(journalEntry.GetCompleted())
			operation := operations[completedId]
			switch operation.(type) {
			case *AddSegmentOperation:
				add := operation.(*AddSegmentOperation)
				if add.completingSegment != nil {
					segments[add.completingSegment.Segment.Id] = add.completingSegment.Segment
				}
				ongoing = add.newSegment
			case *MergeSegmentsOperation:
				merge := operation.(*MergeSegmentsOperation)
				delete(segments, merge.segmentsToMerge[0].Id)
				delete(segments, merge.segmentsToMerge[1].Id)
				segments[merge.mergedSegment.Id] = merge.mergedSegment
			}

			delete(operations, completedId)
		}
	}

	return &CurrentSegments{
		Ongoing:    ongoing,
		Segments:   segments,
		Operations: operations,
	}, nil
}

func segmentOperationFromProto(proto *segmentsv1.SegmentOperation) SegmentOperation {
	switch proto.Operation.(type) {
	case *segmentsv1.SegmentOperation_Add:
		addOperation := proto.GetAdd()
		completingSegmentPb := addOperation.CompletingSegment
		var completingSegment *LiveSegment
		if completingSegmentPb != nil {
			completingSegment = &LiveSegment{
				WalFilename: storage.Filename(completingSegmentPb.WalFilename),
				Segment: Segment{
					Id:       SegmentId(completingSegmentPb.Segment.Id),
					Filename: storage.Filename(completingSegmentPb.Segment.Filename),
				},
			}
		}
		newSegmentPb := addOperation.NewSegment
		return &AddSegmentOperation{
			id:                OperationId(proto.Id),
			completingSegment: completingSegment,
			newSegment: &LiveSegment{
				WalFilename: storage.Filename(newSegmentPb.WalFilename),
				Segment: Segment{
					Id:       SegmentId(newSegmentPb.Segment.Id),
					Filename: storage.Filename(newSegmentPb.Segment.Filename),
				},
			},
		}
	case *segmentsv1.SegmentOperation_Merge:
		mergeOperation := proto.GetMerge()
		segmentsToMerge := make([]Segment, len(mergeOperation.SegmentsToMerge))
		for i, segmentProto := range mergeOperation.SegmentsToMerge {
			segmentsToMerge[i] = Segment{
				Id:       SegmentId(segmentProto.Id),
				Filename: storage.Filename(segmentProto.Filename),
			}
		}
		return &MergeSegmentsOperation{
			id:              OperationId(proto.Id),
			segmentsToMerge: segmentsToMerge,
			mergedSegment: Segment{
				Id:       SegmentId(mergeOperation.NewSegment.Id),
				Filename: storage.Filename(mergeOperation.NewSegment.Filename),
			},
		}
	}

	return nil
}

func (soj *SegmentOperationJournal) Write(entry *segmentsv1.SegmentOperationJournalEntry) error {
	return soj.journal.Writer.Write(entry)
}

func (soj *SegmentOperationJournal) Close() error {
	return soj.journal.Close()
}
