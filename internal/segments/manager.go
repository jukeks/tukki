package segments

import (
	"fmt"
	"log"

	"github.com/jukeks/tukki/internal/memtable"
	"github.com/jukeks/tukki/internal/storage"
)

type SegmentManager struct {
	dbDir            string
	operationJournal *SegmentOperationJournal

	segments   map[SegmentId]Segment
	operations map[OperationId]SegmentOperation

	ongoing LiveSegment
}

func OpenDatabase(dbDir string) (*SegmentManager, error) {
	operationJournal, currentSegments, err := OpenSegmentOperationJournal(dbDir)
	if err != nil {
		return nil, err
	}

	bootstrapped := true
	if currentSegments == nil {
		bootstrapped = false
		currentSegments = &CurrentSegments{
			Segments:   make(map[SegmentId]Segment),
			Operations: make(map[OperationId]SegmentOperation),
		}
	}

	sm := &SegmentManager{
		dbDir:            dbDir,
		operationJournal: operationJournal,
		segments:         currentSegments.Segments,
		operations:       currentSegments.Operations,
		ongoing:          currentSegments.Ongoing,
	}

	if !bootstrapped {
		err = sm.Initialize()
		if err != nil {
			log.Printf("failed to initialize segment manager: %v", err)
			return nil, err
		}
	}

	return sm, nil
}

func (sm *SegmentManager) getNextOperationId() OperationId {
	// get highest operation id in ongoing operations
	var maxId OperationId
	for id := range sm.operations {
		if id > maxId {
			maxId = id
		}
	}

	return maxId + 1
}

func getWalFilename(id SegmentId) storage.Filename {
	return storage.Filename(fmt.Sprintf("wal-%d.journal", id))
}

func getSegmentFilename(id SegmentId) storage.Filename {
	return storage.Filename(fmt.Sprintf("segment-%d", id))
}

func (sm *SegmentManager) GetOnGoingSegment() LiveSegment {
	return sm.ongoing
}

func (sm *SegmentManager) Initialize() error {
	firstSegment := &LiveSegment{
		WalFilename: getWalFilename(0),
		Segment: Segment{
			Id:       0,
			Filename: getSegmentFilename(0),
		},
	}
	op := NewAddSegmentOperation(sm.getNextOperationId(), sm.dbDir, nil, firstSegment)
	startEntry := op.StartJournalEntry()
	err := sm.operationJournal.Write(startEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}
	sm.operations[op.Id()] = op

	err = op.Execute()
	if err != nil {
		log.Printf("failed to execute operation: %v", err)
		return err
	}
	sm.ongoing = *firstSegment

	completedEntry := op.CompletedJournalEntry()
	err = sm.operationJournal.Write(completedEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}
	delete(sm.operations, op.Id())

	return nil
}

func (sm *SegmentManager) SealCurrentSegment(mt memtable.Memtable) error {
	ongoingSegment := sm.ongoing
	nextSegmentId := sm.getNextSegmentId()
	nextSegment := &LiveSegment{
		WalFilename: getWalFilename(nextSegmentId),
		Segment: Segment{
			Id:       nextSegmentId,
			Filename: getSegmentFilename(nextSegmentId),
		},
	}
	ongoingSegment.memtable = mt

	op := NewAddSegmentOperation(sm.getNextOperationId(), sm.dbDir, &ongoingSegment, nextSegment)
	startEntry := op.StartJournalEntry()
	err := sm.operationJournal.Write(startEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}
	sm.operations[op.Id()] = op

	err = op.Execute()
	if err != nil {
		log.Printf("failed to execute operation: %v", err)
		return err
	}
	sm.segments[ongoingSegment.Segment.Id] = ongoingSegment.Segment
	sm.ongoing = *nextSegment

	completedEntry := op.CompletedJournalEntry()
	err = sm.operationJournal.Write(completedEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}
	delete(sm.operations, op.Id())

	return nil
}

func (sm *SegmentManager) Close() error {
	return sm.operationJournal.Close()
}

func (sm *SegmentManager) getNextSegmentId() SegmentId {
	return sm.ongoing.Segment.Id + 1
}

func getMergedSegmentFilename(a, b SegmentId) storage.Filename {
	return storage.Filename(fmt.Sprintf("segment-%d-%d", a, b))
}

func (sm *SegmentManager) MergeSegments(a, b SegmentId) error {
	segmentA := sm.segments[a]
	segmentB := sm.segments[b]

	mergedSegment := Segment{
		Id:       segmentB.Id,
		Filename: getMergedSegmentFilename(segmentA.Id, segmentB.Id),
	}

	op := NewMergeSegmentsOperation(sm.getNextOperationId(), sm.dbDir, []Segment{segmentA, segmentB}, mergedSegment)
	startEntry := op.StartJournalEntry()
	err := sm.operationJournal.Write(startEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}
	sm.operations[op.Id()] = op

	err = op.Execute()
	if err != nil {
		log.Printf("failed to execute operation: %v", err)
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
	delete(sm.operations, op.Id())

	return nil
}
