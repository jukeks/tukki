package db

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/segments"
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

func (db *Database) Snapshot() *Snapshot {
	wal := db.ongoing.Wal.Snapshot()
	operations := db.operationJournal.Snapshot()

	return NewSnapshot(wal, operations)
}

type RestoreResult struct {
	MissingSegments []segments.SegmentMetadata
}

func (db *Database) Restore(snapshot *Snapshot) (*RestoreResult, error) {
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

	// overwrite current WAL
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
