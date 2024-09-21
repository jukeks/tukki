package kv

import (
	"context"

	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
)

type DB interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type KVServer struct {
	kvv1.UnimplementedKvServiceServer
	db DB
}

func NewKVServer(db DB) *KVServer {
	return &KVServer{db: db}
}

func (s *KVServer) Query(ctx context.Context, req *kvv1.QueryRequest) (*kvv1.QueryResponse, error) {
	value, err := s.db.Get(req.Key)
	if err != nil {
		return nil, err
	}

	return &kvv1.QueryResponse{Value: &kvv1.QueryResponse_Pair{Pair: &kvv1.KvPair{
		Key:   req.Key,
		Value: value,
	}}}, nil
}

func (s *KVServer) Set(ctx context.Context, req *kvv1.SetRequest) (*kvv1.SetResponse, error) {
	err := s.db.Set(req.Pair.Key, req.Pair.Value)
	if err != nil {
		return nil, err
	}

	return &kvv1.SetResponse{}, nil
}

func (s *KVServer) Delete(ctx context.Context, req *kvv1.DeleteRequest) (*kvv1.DeleteResponse, error) {
	err := s.db.Delete(req.Key)
	if err != nil {
		return nil, err
	}

	return &kvv1.DeleteResponse{}, nil
}
