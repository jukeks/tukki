package segments

import "fmt"

type SegmentManager struct {
	segmentJournal   *SegmentJournal
	operationJournal *SegmentOperationJournal

	segments   map[SegmentId]Segment
	operations map[OperationId]SegmentOperation

	ongoing OngoingSegment
}

func OpenDatabase(dbDir string) (*SegmentManager, error) {
	segmentJournal, currentSegments, err := OpenSegmentJournal(dbDir)
	if err != nil {
		return nil, err
	}

	if currentSegments == nil {
		ongoing := OngoingSegment{
			Id:              0,
			JournalFilename: getWalFilename(0),
		}
		currentSegments = &CurrentSegments{
			Ongoing:  ongoing,
			Segments: make(map[SegmentId]Segment),
		}

		err = segmentJournal.StartSegment(ongoing.Id, ongoing.JournalFilename)
		if err != nil {
			return nil, err
		}
	}

	operationJournal, operations, err := OpenSegmentOperationJournal(dbDir)
	if err != nil {
		return nil, err
	}

	return &SegmentManager{
		segmentJournal:   segmentJournal,
		operationJournal: operationJournal,
		segments:         currentSegments.Segments,
		operations:       operations,
		ongoing:          currentSegments.Ongoing,
	}, nil
}

func getWalFilename(id SegmentId) string {
	return fmt.Sprintf("wal-%d.journal", id)
}

func (sm *SegmentManager) GetOnGoingSegment() OngoingSegment {
	return sm.ongoing
}
