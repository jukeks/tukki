package kv

import (
	"context"
	"io"
	"sort"
	"testing"

	"google.golang.org/grpc"

	"github.com/jukeks/tukki/internal/db"
	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
	"github.com/jukeks/tukki/testutil"
)

type Store struct {
	store map[string]string
}

func (s *Store) Get(key string) (string, error) {
	value, ok := s.store[key]
	if !ok {
		return "", nil
	}

	return value, nil
}

func (s *Store) Set(key, value string) error {
	s.store[key] = value
	return nil
}

func (s *Store) Delete(key string) error {
	delete(s.store, key)
	return nil
}

type iterator struct {
	data []db.Pair
	pos  int
}

func (i *iterator) Next() (db.Pair, error) {
	if i.pos >= len(i.data) {
		return db.Pair{}, io.EOF
	}

	pair := i.data[i.pos]
	i.pos++
	return pair, nil
}

func (s *Store) GetRange(min, max string) (db.KeyValueIterator, error) {
	resp := make([]db.Pair, 0, len(s.store))
	for k, v := range s.store {
		if k < min || k > max {
			continue
		}
		resp = append(resp, db.Pair{Key: k, Value: v})
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Key < resp[j].Key
	})

	return &iterator{data: resp}, nil
}

func (s *Store) DeleteRange(min, max string) (uint64, error) {
	toDelete, err := s.GetRange(min, max)
	if err != nil {
		return 0, err
	}

	deleted := 0
	for pair, err := toDelete.Next(); err != io.EOF; pair, err = toDelete.Next() {
		delete(s.store, pair.Key)
		deleted++
	}

	return uint64(deleted), nil
}

func TestKvServer(t *testing.T) {
	store := &Store{store: make(map[string]string)}

	conn, cleanup, err := testutil.RunServicer(func(s *grpc.Server) {
		kvv1.RegisterKvServiceServer(s, NewKVServer(store))
	})
	if err != nil {
		t.Fatalf("failed to run servicer: %v", err)
	}
	defer cleanup()

	client := kvv1.NewKvServiceClient(conn)
	ctx := context.Background()

	key := "key"
	value := "value"

	_, err = client.Set(ctx, &kvv1.SetRequest{
		Pair: &kvv1.KvPair{
			Key:   key,
			Value: value,
		},
	})

	if err != nil {
		t.Fatalf("failed to set key-value pair: %v", err)
	}

	resp, err := client.Query(ctx, &kvv1.QueryRequest{
		Key: key,
	})
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}

	if resp.GetPair().Value != value {
		t.Fatalf("expected value %s, got %s", value, resp.GetPair().Value)
	}

	_, err = client.Delete(ctx, &kvv1.DeleteRequest{
		Key: key,
	})
	if err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}

	resp, err = client.Query(ctx, &kvv1.QueryRequest{
		Key: key,
	})
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}

	if resp.GetPair().Value != "" {
		t.Fatalf("expected empty value, got %s", resp.GetPair().Value)
	}
}

func readStream(stream kvv1.KvService_QueryRangeClient) ([]Pair, error) {
	var resp []Pair
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		resp = append(resp, Pair{Key: msg.Pair.Key, Value: msg.Pair.Value})
	}

	return resp, nil
}

func TestKvServerRanges(t *testing.T) {
	store := &Store{store: make(map[string]string)}

	conn, cleanup, err := testutil.RunServicer(func(s *grpc.Server) {
		kvv1.RegisterKvServiceServer(s, NewKVServer(store))
	})
	if err != nil {
		t.Fatalf("failed to run servicer: %v", err)
	}
	defer cleanup()

	client := kvv1.NewKvServiceClient(conn)
	ctx := context.Background()

	pairs := []Pair{
		{Key: "a", Value: "1"},
		{Key: "b", Value: "2"},
		{Key: "c", Value: "3"},
	}

	for _, pair := range pairs {
		_, err = client.Set(ctx, &kvv1.SetRequest{
			Pair: &kvv1.KvPair{
				Key:   pair.Key,
				Value: pair.Value,
			},
		})
		if err != nil {
			t.Fatalf("failed to set key-value pair: %v", err)
		}
	}

	stream, err := client.QueryRange(ctx, &kvv1.QueryRangeRequest{
		Min: "a",
		Max: "c",
	})
	if err != nil {
		t.Fatalf("failed to query range: %v", err)
	}

	resp, err := readStream(stream)
	if err != nil {
		t.Fatalf("failed to read stream: %v", err)
	}

	if len(resp) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(resp))
	}

	for i, pair := range pairs {
		if resp[i].Key != pair.Key {
			t.Fatalf("expected key %s, got %s", pair.Key, resp[i].Key)
		}
		if resp[i].Value != pair.Value {
			t.Fatalf("expected value %s, got %s", pair.Value, resp[i].Value)
		}
	}

	_, err = client.DeleteRange(ctx, &kvv1.DeleteRangeRequest{
		Min: "a",
		Max: "b",
	})
	if err != nil {
		t.Fatalf("failed to delete range: %v", err)
	}

	stream, err = client.QueryRange(ctx, &kvv1.QueryRangeRequest{
		Min: "a",
		Max: "c",
	})
	if err != nil {
		t.Fatalf("failed to query range: %v", err)
	}

	resp, err = readStream(stream)
	if err != nil {
		t.Fatalf("failed to read stream: %v", err)
	}

	if len(resp) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(resp))
	}

	if resp[0].Key != "c" {
		t.Fatalf("expected key c, got %s", resp[0].Key)
	}
}
