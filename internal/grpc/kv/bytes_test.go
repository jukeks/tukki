package kv

import (
	"bytes"
	"context"
	"testing"

	"google.golang.org/grpc"

	"github.com/jukeks/tukki/internal/db"
	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
	"github.com/jukeks/tukki/testutil"
)

// databaseStore adapts a real database to the server's DB interface.
type databaseStore struct {
	db *db.Database
}

func (s *databaseStore) Get(key []byte) ([]byte, error) {
	return s.db.Get(key)
}

func (s *databaseStore) Set(key, value []byte) error {
	return s.db.Set(key, value)
}

func (s *databaseStore) Delete(key []byte) error {
	return s.db.Delete(key)
}

func (s *databaseStore) GetRange(min, max []byte) (db.KeyValueIterator, error) {
	return s.db.GetCursorWithRange(min, max)
}

func (s *databaseStore) DeleteRange(min, max []byte) (uint64, error) {
	deleted, err := s.db.DeleteRange(min, max)
	return uint64(deleted), err
}

func TestKvServerBinaryData(t *testing.T) {
	d, err := db.OpenDatabase(t.TempDir())
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer d.Close()

	conn, cleanup, err := testutil.RunServicer(func(s *grpc.Server) {
		kvv1.RegisterKvServiceServer(s, NewKVServer(&databaseStore{db: d}))
	})
	if err != nil {
		t.Fatalf("failed to run servicer: %v", err)
	}
	defer cleanup()

	client := kvv1.NewKvServiceClient(conn)
	ctx := context.Background()

	// These keys and values are not valid UTF-8.
	pairs := []Pair{
		{Key: []byte{0x00}, Value: []byte{0xff, 0x00, 0x80}},
		{Key: []byte{0x80, 0x01}, Value: []byte{0xc3, 0x28}},
		{Key: []byte{0xff}, Value: []byte{0x00}},
	}

	for _, pair := range pairs {
		_, err = client.Set(ctx, &kvv1.SetRequest{
			Pair: &kvv1.KvPair{Key: pair.Key, Value: pair.Value},
		})
		if err != nil {
			t.Fatalf("failed to set %x: %v", pair.Key, err)
		}
	}

	for _, pair := range pairs {
		resp, err := client.Query(ctx, &kvv1.QueryRequest{Key: pair.Key})
		if err != nil {
			t.Fatalf("failed to query %x: %v", pair.Key, err)
		}
		if !bytes.Equal(resp.GetPair().Key, pair.Key) {
			t.Fatalf("expected key %x, got %x", pair.Key, resp.GetPair().Key)
		}
		if !bytes.Equal(resp.GetPair().Value, pair.Value) {
			t.Fatalf("expected value %x, got %x", pair.Value, resp.GetPair().Value)
		}
	}

	stream, err := client.QueryRange(ctx, &kvv1.QueryRangeRequest{
		Min: pairs[0].Key,
		Max: pairs[len(pairs)-1].Key,
	})
	if err != nil {
		t.Fatalf("failed to query range: %v", err)
	}

	got, err := readStream(stream)
	if err != nil {
		t.Fatalf("failed to read stream: %v", err)
	}

	if len(got) != len(pairs) {
		t.Fatalf("expected %d pairs, got %d", len(pairs), len(got))
	}
	for i, pair := range pairs {
		if !bytes.Equal(got[i].Key, pair.Key) {
			t.Fatalf("expected key %x, got %x", pair.Key, got[i].Key)
		}
		if !bytes.Equal(got[i].Value, pair.Value) {
			t.Fatalf("expected value %x, got %x", pair.Value, got[i].Value)
		}
	}

	resp, err := client.DeleteRange(ctx, &kvv1.DeleteRangeRequest{
		Min: pairs[0].Key,
		Max: pairs[len(pairs)-1].Key,
	})
	if err != nil {
		t.Fatalf("failed to delete range: %v", err)
	}
	if resp.Deleted != uint64(len(pairs)) {
		t.Fatalf("expected %d deleted, got %d", len(pairs), resp.Deleted)
	}

	for _, pair := range pairs {
		_, err := client.Query(ctx, &kvv1.QueryRequest{Key: pair.Key})
		if err == nil {
			t.Fatalf("expected an error for deleted key %x", pair.Key)
		}
	}
}
