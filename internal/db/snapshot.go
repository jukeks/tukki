package db

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/index"
	"github.com/jukeks/tukki/internal/storage/keyvalue"
	"github.com/jukeks/tukki/internal/storage/segmentmembers"
	"github.com/jukeks/tukki/internal/storage/segments"
	"github.com/jukeks/tukki/internal/storage/sstable"
	snapshotv1 "github.com/jukeks/tukki/proto/gen/tukki/replication/snapshot/v1"
	"google.golang.org/protobuf/proto"
)

type Snapshot struct {
	Wal        []byte
	Operations []byte
}

func NewSnapshot(wal []byte, operations []byte) *Snapshot {
	return &Snapshot{
		Wal:        wal,
		Operations: operations,
	}
}

func UnmarshalSnapshot(data []byte) (*Snapshot, error) {
	snapshot := &snapshotv1.Snapshot{}
	err := proto.Unmarshal(data, snapshot)
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		Wal:        snapshot.Wal,
		Operations: snapshot.Operations,
	}, nil
}

func (s *Snapshot) Marshal() ([]byte, error) {
	snapshot := &snapshotv1.Snapshot{
		Wal:        s.Wal,
		Operations: s.Operations,
	}

	return proto.Marshal(snapshot)
}

func (db *Database) Snapshot() (*Snapshot, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	currentSegments := db.getSegmentsSortedUnlocked()
	journalEntry := segments.NewSnapshotJournalEntry(currentSegments)
	if err := db.operationJournal.Write(journalEntry); err != nil {
		return nil, fmt.Errorf("failed to write snapshot journal entry: %w", err)
	}

	wal := db.ongoing.Wal.Snapshot()
	operations := db.operationJournal.Snapshot()

	db.lastSnapshotSegments = currentSegments

	newPinnedSegments := make([]segments.SegmentMetadata, 0)
	okToRemove := make([]segments.SegmentMetadata, 0)
	// check if we can remove any segments after new snapshot
	for _, pinnedSegment := range db.freedButNotRemoved {
		isPinned := false
		for _, currentSegment := range currentSegments {
			if currentSegment.SegmentFile == pinnedSegment.SegmentFile {
				// segment is pinned, so we need to keep it still
				isPinned = true
				continue
			}
		}
		if isPinned {
			newPinnedSegments = append(newPinnedSegments, pinnedSegment)
			continue
		}

		log.Printf("snapshot: segment %s is not pinned anymore, removing", pinnedSegment.SegmentFile)
		okToRemove = append(okToRemove, pinnedSegment)
	}

	for _, segment := range okToRemove {
		files.RemoveFile(db.dbDir, segment.SegmentFile)
		files.RemoveFile(db.dbDir, segment.IndexFile)
		files.RemoveFile(db.dbDir, segment.MembersFile)
	}

	db.freedButNotRemoved = newPinnedSegments

	return NewSnapshot(wal, operations), nil
}

type RestoreResult struct {
	AllSegments     []segments.SegmentMetadata
	MissingSegments []segments.SegmentMetadata
}

func (db *Database) Restore(snapshot *Snapshot) (*RestoreResult, error) {
	// remove current WAL
	if err := files.RemoveFile(db.dbDir, db.ongoing.WalFilename); err != nil {
		return nil, fmt.Errorf("failed to remove wal file: %w", err)
	}

	// overwrite segment journal
	f, err := files.CreateFile(db.dbDir, segments.SegmentJournalFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to create segment journal file: %w", err)
	}
	opBuff := bytes.NewBuffer(snapshot.Operations)
	_, err = io.Copy(f, opBuff)
	if err != nil {
		return nil, fmt.Errorf("failed to write operations to segment journal: %w", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close segment journal file: %w", err)
	}

	// read current segments
	opJournal, currentSegments, err := segments.OpenSegmentOperationJournal(db.dbDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open segment journal: %w", err)
	}
	if err := opJournal.Close(); err != nil {
		return nil, fmt.Errorf("failed to close segment journal: %w", err)
	}

	if currentSegments == nil {
		return nil, errors.New("failed to read current segments")
	}

	// recreate WAL
	walBuff := bytes.NewBuffer(snapshot.Wal)
	f, err = files.CreateFile(db.dbDir, currentSegments.Ongoing.WalFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to create segment journal file: %w", err)
	}
	_, err = io.Copy(f, walBuff)
	if err != nil {
		return nil, fmt.Errorf("failed to write WAL to file: %w", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close WAL file: %w", err)
	}

	// check for missing segments
	missingSegments := make([]segments.SegmentMetadata, 0)
	for _, metadata := range currentSegments.Segments {
		err := checkFilesExist(db.dbDir, metadata.SegmentFile, metadata.MembersFile, metadata.IndexFile)
		if err != nil {
			missingSegments = append(missingSegments, metadata)
		}
	}

	return &RestoreResult{MissingSegments: missingSegments}, nil
}

func checkFilesExist(dbDir string, filenames ...files.Filename) error {
	for _, file := range filenames {
		if !files.FileExists(dbDir, file) {
			return fmt.Errorf("file %s does not exist", file)
		}
	}
	return nil
}

func (db *Database) RestoreSegment(segment segments.SegmentMetadata, iterator keyvalue.KeyValueIterator) error {
	f, err := files.CreateFile(db.dbDir, segment.SegmentFile)
	if err != nil {
		return fmt.Errorf("failed to create segment file: %w", err)
	}

	writer := sstable.NewSSTableWriter(f)
	if err := writer.WriteFromIterator(iterator); err != nil {
		return fmt.Errorf("failed to write segment: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close segment file: %w", err)
	}

	offsets := writer.WrittenOffsets()
	f, err = files.CreateFile(db.dbDir, segment.IndexFile)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	indexWriter := index.NewIndexWriter(f)
	if err := indexWriter.WriteFromOffsets(offsets); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close index file: %w", err)
	}

	members := segmentmembers.NewSegmentMembers(uint(len(offsets)))
	for key := range offsets {
		members.Add(key)
	}
	if err := members.Save(db.dbDir, segment.MembersFile); err != nil {
		return fmt.Errorf("failed to save members: %w", err)
	}

	return nil
}
