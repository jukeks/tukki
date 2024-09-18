package db

import (
	"errors"

	"github.com/jukeks/tukki/internal/segments"
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
	MissingSegments []segments.SegmentId
}

func (db *Database) Restore(snapshot *Snapshot) (RestoreResult, error) {
	return RestoreResult{}, errors.New("not implemented")
}
