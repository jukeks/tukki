package kv

import (
	"context"
	"sort"
	"testing"

	"google.golang.org/grpc"

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

func (s *Store) GetRange(min, max string) ([]Pair, error) {
	resp := make([]Pair, 0, len(s.store))
	for k, v := range s.store {
		if k < min || k > max {
			continue
		}
		resp = append(resp, Pair{Key: k, Value: v})
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Key < resp[j].Key
	})

	return resp, nil
}

func (s *Store) DeleteRange(min, max string) (uint64, error) {
	toDelete, err := s.GetRange(min, max)
	if err != nil {
		return 0, err
	}

	for _, pair := range toDelete {
		delete(s.store, pair.Key)
	}

	return uint64(len(toDelete)), nil
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
