package kv

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
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

func TestKvServer(t *testing.T) {
	store := &Store{store: make(map[string]string)}

	server := NewKVServer(store)
	s := grpc.NewServer()
	kvv1.RegisterKvServiceServer(s, server)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("failed to serve: %v", err)
		}
	}()
	defer s.Stop()

	conn, err := grpc.Dial(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

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
