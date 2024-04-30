package segments

import (
	"io"

	"github.com/jukeks/tukki/internal/journal"
	"github.com/jukeks/tukki/internal/storage"
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

const SegmentJournalFilename = "segment_operations.journal"

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

	j, err := journal.OpenJournal(dbDir, SegmentJournalFilename,
		journal.WriteModeSync, handle)
	if err != nil {
		return nil, nil, err
	}

	return &SegmentOperationJournal{j}, currentSegments, nil
}

type CurrentSegments struct {
	Ongoing    *OpenSegment
	Segments   map[SegmentId]SegmentMetadata
	Operations map[OperationId]SegmentOperation
}

func readOperationJournal(r *journal.JournalReader) (
	*CurrentSegments,
	error) {

	operations := make(map[OperationId]SegmentOperation)
	segments := make(map[SegmentId]SegmentMetadata)
	var ongoing *OpenSegment

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
			switch op := operation.(type) {
			case *AddSegmentOperation:
				if op.completingSegment != nil {
					segments[op.completingSegment.Segment.Id] = op.completingSegment.Segment
				}
				ongoing = op.newSegment
			case *MergeSegmentsOperation:
				delete(segments, op.segmentsToMerge[0].Id)
				delete(segments, op.segmentsToMerge[1].Id)
				segments[op.mergedSegment.Id] = op.mergedSegment
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

func pbToSegmentMetadata(segmentPb *segmentsv1.Segment) *SegmentMetadata {
	return &SegmentMetadata{
		Id:          SegmentId(segmentPb.Id),
		SegmentFile: storage.Filename(segmentPb.Filename),
		MembersFile: storage.Filename(segmentPb.MembersFilename),
	}
}

func segmentMetadataToPb(segment *SegmentMetadata) *segmentsv1.Segment {
	return &segmentsv1.Segment{
		Id:              uint64(segment.Id),
		Filename:        string(segment.SegmentFile),
		MembersFilename: string(segment.MembersFile),
	}
}

func segmentOperationFromProto(proto *segmentsv1.SegmentOperation) SegmentOperation {
	switch proto.Operation.(type) {
	case *segmentsv1.SegmentOperation_Add:
		addOperation := proto.GetAdd()
		completingSegmentPb := addOperation.CompletingSegment
		var completingSegment *OpenSegment
		if completingSegmentPb != nil {
			completingSegment = &OpenSegment{
				WalFilename: storage.Filename(completingSegmentPb.WalFilename),
				Segment:     *pbToSegmentMetadata(completingSegmentPb.Segment),
			}
		}
		newSegmentPb := addOperation.NewSegment
		return &AddSegmentOperation{
			id:                OperationId(proto.Id),
			completingSegment: completingSegment,
			newSegment: &OpenSegment{
				WalFilename: storage.Filename(newSegmentPb.WalFilename),
				Segment:     *pbToSegmentMetadata(newSegmentPb.Segment),
			},
		}
	case *segmentsv1.SegmentOperation_Merge:
		mergeOperation := proto.GetMerge()
		segmentsToMerge := make([]SegmentMetadata, len(mergeOperation.SegmentsToMerge))
		for i, segmentProto := range mergeOperation.SegmentsToMerge {
			segmentsToMerge[i] = *pbToSegmentMetadata(segmentProto)
		}
		return &MergeSegmentsOperation{
			id:              OperationId(proto.Id),
			segmentsToMerge: segmentsToMerge,
			mergedSegment:   *pbToSegmentMetadata(mergeOperation.NewSegment),
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
