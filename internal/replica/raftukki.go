package replica

import (
	"encoding/base64"
	"fmt"
	"io"
	"strconv"

	"github.com/hashicorp/raft"
	"github.com/jukeks/tukki/internal/db"
	raftlogv1 "github.com/jukeks/tukki/proto/gen/tukki/replication/raftlog/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Raftukki struct {
	db *db.Database
}

func tukkiWrapNotFound(err error) error {
	if err == nil {
		return nil
	}
	if err == db.ErrKeyNotFound {
		return ErrKeyNotFound
	}
	return err
}

// Get implements raft.StableStore.
func (t *Raftukki) Get(key []byte) ([]byte, error) {
	val, err := t.db.Get(string(key))
	return []byte(val), tukkiWrapNotFound(err)
}

// Set implements raft.StableStore.
func (t *Raftukki) Set(key []byte, val []byte) error {
	return t.db.Set(string(key), string(val))
}

// GetUint64 implements raft.StableStore.
func (t *Raftukki) GetUint64(key []byte) (uint64, error) {
	val, err := t.db.Get(string(key))
	if err != nil {
		if err == db.ErrKeyNotFound {
			return 0, ErrKeyNotFound
		}
		return 0, tukkiWrapNotFound(err)
	}
	number, err := strconv.ParseUint(val, 10, 64)
	return number, tukkiWrapNotFound(err)
}

// SetUint64 implements raft.StableStore.
func (t *Raftukki) SetUint64(key []byte, val uint64) error {
	return t.db.Set(string(key), fmt.Sprintf("%d", val))
}

const logPrefix = "raft-log-"

func logKey(index uint64) string {
	return fmt.Sprintf("%s%d", logPrefix, index)
}

// DeleteRange implements raft.LogStore.
func (t *Raftukki) DeleteRange(min uint64, max uint64) error {
	_, err := t.db.DeleteRange(logKey(min), logKey(max))
	return tukkiWrapNotFound(err)
}

// FirstIndex implements raft.LogStore.
func (t *Raftukki) FirstIndex() (uint64, error) {
	c, err := t.db.GetCursorWithRange(logKey(0), "")
	if err != nil {
		return 0, tukkiWrapNotFound(err)
	}
	defer c.Close()

	val, err := c.Next()
	if err != nil {
		return 0, tukkiWrapNotFound(err)
	}

	indexStr := val.Key[len(logPrefix):]
	return strconv.ParseUint(indexStr, 10, 64)
}

// LastIndex implements raft.LogStore.
func (t *Raftukki) LastIndex() (uint64, error) {
	// this is super slow, need to have reverse cursor to make it faster
	c, err := t.db.GetCursorWithRange(logKey(0), "")
	if err != nil {
		return 0, tukkiWrapNotFound(err)
	}
	defer c.Close()

	var lastIndex uint64
	for {
		val, err := c.Next()
		if err != nil {
			if err == io.EOF {
				return lastIndex, nil
			}
			return 0, tukkiWrapNotFound(err)
		}

		indexStr := val.Key[len(logPrefix):]
		index, err := strconv.ParseUint(indexStr, 10, 64)
		if err != nil {
			return 0, err
		}

		if index > lastIndex {
			lastIndex = index
		}
	}
}

// StoreLog implements raft.LogStore.
func (t *Raftukki) StoreLog(log *raft.Log) error {
	return t.StoreLogs([]*raft.Log{log})
}

func logFromRaft(r *raft.Log) *raftlogv1.Log {
	return &raftlogv1.Log{
		Index:      r.Index,
		Term:       r.Term,
		Type:       uint32(r.Type),
		Data:       r.Data,
		Extensions: r.Extensions,
		AppendedAt: timestamppb.New(r.AppendedAt),
	}
}

func logToRaft(l *raftlogv1.Log) *raft.Log {
	return &raft.Log{
		Index:      l.Index,
		Term:       l.Term,
		Type:       raft.LogType(l.Type),
		Data:       l.Data,
		Extensions: l.Extensions,
		AppendedAt: l.AppendedAt.AsTime(),
	}
}

// StoreLogs implements raft.LogStore.
func (t *Raftukki) StoreLogs(logs []*raft.Log) error {
	for _, log := range logs {
		protoLog := logFromRaft(log)
		raw, err := proto.Marshal(protoLog)
		if err != nil {
			return err
		}
		b64str := base64.StdEncoding.EncodeToString(raw)
		err = t.db.Set(logKey(log.Index), b64str)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetLog implements raft.LogStore.
func (t *Raftukki) GetLog(index uint64, log *raft.Log) error {
	val, err := t.db.Get(logKey(index))
	if err != nil {
		return tukkiWrapNotFound(err)
	}
	b, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return err
	}
	protoLog := &raftlogv1.Log{}
	err = proto.Unmarshal(b, protoLog)
	if err != nil {
		return err
	}
	*log = *logToRaft(protoLog)
	return nil
}

func NewRaftukki(db *db.Database) *Raftukki {
	return &Raftukki{
		db: db,
	}
}
