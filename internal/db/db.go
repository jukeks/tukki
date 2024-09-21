package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/index"
	"github.com/jukeks/tukki/internal/storage/segmentmembers"
	"github.com/jukeks/tukki/internal/storage/segments"
)

type Database struct {
	dbDir            string
	operationJournal *segments.SegmentOperationJournal

	mu         sync.Mutex
	segments   map[segments.SegmentId]segments.SegmentMetadata
	members    map[segments.SegmentId]*segmentmembers.SegmentMembers
	indexes    map[segments.SegmentId]*index.Index
	operations map[segments.OperationId]segments.SegmentOperation

	ongoing *LiveSegment

	walSizeLimit uint64
}

func OpenDatabase(dbDir string) (*Database, error) {
	operationJournal, currentSegments, err := segments.OpenSegmentOperationJournal(dbDir)
	if err != nil {
		return nil, err
	}

	bootstrapped := true
	if currentSegments == nil {
		bootstrapped = false
		currentSegments = &segments.CurrentSegments{
			Segments:   make(map[segments.SegmentId]segments.SegmentMetadata),
			Operations: make(map[segments.OperationId]segments.SegmentOperation),
		}
	}

	var ongoing *LiveSegment
	if currentSegments.Ongoing != nil {
		ongoing = &LiveSegment{
			Segment:     currentSegments.Ongoing.Segment,
			WalFilename: currentSegments.Ongoing.WalFilename,
		}
	}

	db := &Database{
		dbDir:            dbDir,
		operationJournal: operationJournal,
		segments:         currentSegments.Segments,
		members:          make(map[segments.SegmentId]*segmentmembers.SegmentMembers),
		indexes:          make(map[segments.SegmentId]*index.Index),
		operations:       currentSegments.Operations,
		ongoing:          ongoing,
		walSizeLimit:     2 * 1024 * 1024,
	}

	if !bootstrapped {
		err = db.Initialize()
		if err != nil {
			log.Printf("failed to initialize segment manager: %v", err)
			return nil, err
		}
	}

	err = db.ongoing.Open(dbDir)
	if err != nil {
		log.Printf("failed to open wal: %v", err)
		return nil, err
	}

	for _, segment := range db.segments {
		members, err := segmentmembers.OpenSegmentMembers(dbDir, segment.MembersFile)
		if err != nil {
			log.Printf("failed to open segment members: %v", err)
			return nil, err
		}
		db.members[segment.Id] = members
	}

	for _, segment := range db.segments {
		idx, err := index.OpenIndex(dbDir, segment.IndexFile)
		if err != nil {
			log.Printf("failed to open index: %v", err)
			return nil, err
		}
		db.indexes[segment.Id] = idx
	}

	return db, nil
}

func getWalFilename(id segments.SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("wal-%d.journal", id))
}

func getSegmentFilename(id segments.SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("segment-%d", id))
}

func getMergedSegmentFilename(a, b segments.SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("segment-%d-%d", a, b))
}

func getMembersFilename(id segments.SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("members-%d", id))
}

func getMergedMembersFilename(a, b segments.SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("members-%d-%d", a, b))
}

func getIndexFilename(id segments.SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("index-%d", id))
}

func getMergedIndexFilename(a, b segments.SegmentId) files.Filename {
	return files.Filename(fmt.Sprintf("index-%d-%d", a, b))
}

func (db *Database) GetOnGoingSegment() *LiveSegment {
	return db.ongoing
}

func (db *Database) getNextSegmentId() segments.SegmentId {
	return db.ongoing.Segment.Id + 1
}

func (db *Database) getNextOperationId() segments.OperationId {
	var maxId segments.OperationId
	for id := range db.operations {
		if id > maxId {
			maxId = id
		}
	}

	return maxId + 1
}

func lsToOs(ls *LiveSegment) *segments.OpenSegment {
	return &segments.OpenSegment{
		Segment:     ls.Segment,
		WalFilename: ls.WalFilename,
		Memtable:    ls.Memtable,
	}
}

func (db *Database) Initialize() error {
	firstSegment := NewLiveSegment(0)

	op := segments.NewAddSegmentOperation(
		db.getNextOperationId(),
		db.dbDir,
		nil,
		lsToOs(firstSegment),
	)
	startEntry := op.StartJournalEntry()
	err := db.operationJournal.Write(startEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}
	db.operations[op.Id()] = op

	err = op.Execute()
	if err != nil {
		log.Printf("failed to execute operation: %v", err)
		return err
	}
	db.ongoing = firstSegment

	completedEntry := op.CompletedJournalEntry()
	err = db.operationJournal.Write(completedEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}
	delete(db.operations, op.Id())

	return nil
}

func (db *Database) SealCurrentSegment() (*LiveSegment, error) {
	ongoingSegment := db.ongoing
	nextSegmentId := db.getNextSegmentId()
	nextSegment := NewLiveSegment(nextSegmentId)

	op := segments.NewAddSegmentOperation(
		db.getNextOperationId(),
		db.dbDir,
		lsToOs(ongoingSegment),
		lsToOs(nextSegment),
	)
	startEntry := op.StartJournalEntry()
	err := db.operationJournal.Write(startEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return nil, err
	}
	db.operations[op.Id()] = op

	err = op.Execute()
	if err != nil {
		log.Printf("failed to execute operation: %v", err)
		return nil, err
	}

	db.mu.Lock()
	db.segments[ongoingSegment.Segment.Id] = ongoingSegment.Segment
	db.mu.Unlock()

	db.ongoing = nextSegment
	err = db.ongoing.Open(db.dbDir)
	if err != nil {
		log.Printf("failed to open wal: %v", err)
		return nil, err
	}

	members, err := segmentmembers.OpenSegmentMembers(db.dbDir,
		ongoingSegment.Segment.MembersFile)
	if err != nil {
		log.Printf("failed to open segment members: %v", err)
		return nil, err
	}
	db.members[ongoingSegment.Segment.Id] = members

	index, err := index.OpenIndex(db.dbDir, ongoingSegment.Segment.IndexFile)
	if err != nil {
		log.Printf("failed to open index: %v", err)
		return nil, err
	}
	db.indexes[ongoingSegment.Segment.Id] = index

	completedEntry := op.CompletedJournalEntry()
	err = db.operationJournal.Write(completedEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return nil, err
	}
	delete(db.operations, op.Id())

	return db.ongoing, nil
}

func (db *Database) Close() error {
	err := db.ongoing.Close()
	if err != nil {
		log.Printf("failed to close wal: %v", err)
		return err
	}
	return db.operationJournal.Close()
}

func (db *Database) MergeSegments(a, b segments.SegmentId) error {
	db.mu.Lock()
	segmentA := db.segments[a]
	segmentB := db.segments[b]
	db.mu.Unlock()

	mergedSegment := segments.SegmentMetadata{
		Id:          segmentB.Id,
		SegmentFile: getMergedSegmentFilename(segmentA.Id, segmentB.Id),
		MembersFile: getMergedMembersFilename(segmentA.Id, segmentB.Id),
		IndexFile:   getMergedIndexFilename(segmentA.Id, segmentB.Id),
	}

	op := segments.NewMergeSegmentsOperation(db.getNextOperationId(), db.dbDir, []segments.SegmentMetadata{segmentA, segmentB}, mergedSegment)
	startEntry := op.StartJournalEntry()
	err := db.operationJournal.Write(startEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}
	db.operations[op.Id()] = op

	err = op.Execute()
	if err != nil {
		log.Printf("failed to execute operation: %v", err)
		return err
	}

	members, err := segmentmembers.OpenSegmentMembers(db.dbDir,
		mergedSegment.MembersFile)
	if err != nil {
		log.Printf("failed to open segment members: %v", err)
		return err
	}
	db.mu.Lock()
	delete(db.members, a)
	delete(db.members, b)
	db.members[mergedSegment.Id] = members

	delete(db.segments, a)
	delete(db.segments, b)
	db.segments[mergedSegment.Id] = mergedSegment
	db.mu.Unlock()

	completedEntry := op.CompletedJournalEntry()
	err = db.operationJournal.Write(completedEntry)
	if err != nil {
		log.Printf("failed to write journal entry: %v", err)
		return err
	}
	delete(db.operations, op.Id())

	return nil
}
