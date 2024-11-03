package segments

import (
	segmentsv1 "github.com/jukeks/tukki/proto/gen/tukki/storage/segments/v1"
)

func NewSnapshotJournalEntry(sgs []SegmentMetadata) *segmentsv1.SegmentOperationJournalEntry {
	pbSegments := make([]*segmentsv1.Segment, 0, len(sgs))
	for _, sg := range sgs {
		pbSegments = append(pbSegments, segmentMetadataToPb(&sg))
	}

	return &segmentsv1.SegmentOperationJournalEntry{
		Entry: &segmentsv1.SegmentOperationJournalEntry_Snapshot{
			Snapshot: &segmentsv1.Snapshot{
				Segments: pbSegments,
			},
		},
	}
}
