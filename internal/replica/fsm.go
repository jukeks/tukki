package replica

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/raft"
	"github.com/jukeks/tukki/internal/db"
	"github.com/jukeks/tukki/internal/storage/files"
	"github.com/jukeks/tukki/internal/storage/keyvalue"
	"github.com/jukeks/tukki/internal/storage/segments"
	"github.com/jukeks/tukki/internal/storage/sstable"
	sstablev1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/sstable/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type fsm Node

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case "set":
		return f.applySet(c.Key, c.Value)
	case "delete":
		return f.applyDelete(c.Key)
	default:
		panic(fmt.Sprintf("unrecognized command op: %s", c.Op))
	}
}

func (f *fsm) applySet(key, value string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.db.Set(key, value); err != nil {
		f.logger.Printf("failed to apply set key %s: %s", key, err)
		return err
	}
	return nil
}

func (f *fsm) applyDelete(key string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.db.Delete(key); err != nil {
		f.logger.Printf("failed to apply delete key %s: %s", key, err)
		return err
	}
	return nil
}

type snapshot struct {
	snapshot *db.Snapshot
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	ss := f.db.Snapshot()
	return &snapshot{snapshot: ss}, nil
}

func (f *fsm) Restore(rc io.ReadCloser) error {
	buff, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	snapshot, err := db.UnmarshalSnapshot(buff)
	if err != nil {
		return err
	}

	/*
		currentSegments := f.db.GetSegmentMetadata()
		ongoingSegment := f.db.GetOnGoingSegment()
	*/

	f.db.Close()

	result, err := f.db.Restore(snapshot)
	if err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	if err := f.handleMissingSegments(result.MissingSegments); err != nil {
		return fmt.Errorf("failed to handle missing segments: %w", err)
	}

	db, err := db.OpenDatabase(f.dbDir)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	f.db = db

	return err
}

func (f *fsm) handleMissingSegments(missingSegments []segments.SegmentMetadata) error {
	if len(missingSegments) == 0 {
		return nil
	}

	for _, peer := range f.peers {
		conn, err := grpc.Dial(peer.Addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("can not connect with server %v", err)
		}
		client := sstablev1.NewSstableServiceClient(conn)
		for _, segment := range missingSegments {
			req := &sstablev1.GetSstableRequest{
				Id: uint64(segment.Id),
			}
			stream, err := client.GetSstable(context.Background(), req)
			if err != nil {
				log.Fatalf("can not get sstable from server %v", err)
			}

			f, err := files.CreateFile(f.dbDir, segment.SegmentFile)
			if err != nil {
				return fmt.Errorf("can not create file %w", err)
			}

			bw := bufio.NewWriter(f)
			writer := sstable.NewSSTableWriter(bw)
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("can not receive sstable from server %w", err)
				}
				_, err = writer.Write(keyvalue.IteratorEntry{
					Key:     resp.Record.Key,
					Value:   resp.Record.Value,
					Deleted: resp.Record.Deleted,
				})
				if err != nil {
					return fmt.Errorf("can not add segment to db: %w", err)
				}
			}
			if err := bw.Flush(); err != nil {
				return fmt.Errorf("can not flush writer %w", err)
			}
			if err := f.Close(); err != nil {
				return fmt.Errorf("can not close file %w", err)
			}

			// TODO members and indexes
		}
	}

	return nil
}

func (f *snapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := f.snapshot.Marshal()
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *snapshot) Release() {}
