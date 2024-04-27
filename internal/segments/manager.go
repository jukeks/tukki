package segments

import (
	"fmt"
	"log"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/storage"
)

type SegmentManager struct {
	dbDir            string
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
		dbDir:            dbDir,
		segmentJournal:   segmentJournal,
		operationJournal: operationJournal,
		segments:         currentSegments.Segments,
		operations:       operations,
		ongoing:          currentSegments.Ongoing,
	}, nil
}

func getWalFilename(id SegmentId) storage.Filename {
	return storage.Filename(fmt.Sprintf("wal-%d.journal", id))
}

func getSegmentFilename(id SegmentId) storage.Filename {
	return storage.Filename(fmt.Sprintf("segment-%d", id))
}

func (sm *SegmentManager) GetOnGoingSegment() OngoingSegment {
	return sm.ongoing
}

func (sm *SegmentManager) SealCurrentSegment(mt memtable.Memtable) error {
	segment := sm.ongoing

	op := NewAddSegmentOperation(sm.dbDir, Segment{
		Id:       segment.Id,
		Filename: getSegmentFilename(segment.Id),
	}, mt)
	startEntry := op.StartJournalEntry()
	err := sm.operationJournal.Write(startEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}

	err = op.Execute()
	if err != nil {
		log.Printf("failed to execute operation: %v", err)
		return err
	}

	err = sm.segmentJournal.AddSegment(op.segment)
	if err != nil {
		return err
	}
	sm.segments[segment.Id] = op.segment

	nextSegmentId := sm.getNextSegmentId()
	newOngoing := OngoingSegment{
		Id:              nextSegmentId,
		JournalFilename: getWalFilename(nextSegmentId),
	}

	err = sm.segmentJournal.StartSegment(newOngoing.Id, newOngoing.JournalFilename)
	if err != nil {
		return err
	}
	sm.ongoing = newOngoing

	completedEntry := op.CompletedJournalEntry()
	err = sm.operationJournal.Write(completedEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}

	return nil
}

func (sm *SegmentManager) Close() error {
	err := sm.segmentJournal.Close()
	if err != nil {
		return err
	}

	return sm.operationJournal.Close()
}

func (sm *SegmentManager) getNextSegmentId() SegmentId {
	return sm.ongoing.Id + 1
}

func getMergedSegmentFilename(a, b SegmentId) storage.Filename {
	return storage.Filename(fmt.Sprintf("segment-%d-%d", a, b))
}

func (sm *SegmentManager) MergeSegments(a, b SegmentId) error {
	segmentA := sm.segments[a]
	segmentB := sm.segments[b]

	mergedSegment := Segment{
		Id:       segmentA.Id,
		Filename: getMergedSegmentFilename(segmentA.Id, segmentB.Id),
	}

	op := NewMergeSegmentsOperation(sm.dbDir, []Segment{segmentA, segmentB}, mergedSegment)
	startEntry := op.StartJournalEntry()
	err := sm.operationJournal.Write(startEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}

	err = op.Execute()
	if err != nil {
		log.Printf("failed to execute operation: %v", err)
		return err
	}

	// these should be done in a transaction
	err = sm.segmentJournal.AddSegment(mergedSegment)
	if err != nil {
		return err
	}
	err = sm.segmentJournal.RemoveSegment(b)
	if err != nil {
		return err
	}

	delete(sm.segments, a)
	delete(sm.segments, b)
	sm.segments[mergedSegment.Id] = mergedSegment

	completedEntry := op.CompletedJournalEntry()
	err = sm.operationJournal.Write(completedEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}

	return nil
}
