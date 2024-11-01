package kv

import (
	"context"

	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
)

type Pair struct {
	Key   string
	Value string
}

type DB interface {
	Get(key string) (string, error)
	GetRange(min, max string) ([]Pair, error)
	Set(key, value string) error
	Delete(key string) error
	DeleteRange(min, max string) (uint64, error)
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

func (s *KVServer) QueryRange(ctx context.Context, req *kvv1.QueryRangeRequest) (*kvv1.QueryRangeResponse, error) {
	pairs, err := s.db.GetRange(req.Min, req.Max)
	if err != nil {
		return nil, err
	}

	var kvPairs []*kvv1.KvPair
	for _, pair := range pairs {
		kvPairs = append(kvPairs, &kvv1.KvPair{
			Key:   pair.Key,
			Value: pair.Value,
		})
	}

	return &kvv1.QueryRangeResponse{Pairs: kvPairs}, nil
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

func (s *KVServer) DeleteRange(ctx context.Context, req *kvv1.DeleteRangeRequest) (*kvv1.DeleteRangeResponse, error) {
	deleted, err := s.db.DeleteRange(req.Min, req.Max)
	if err != nil {
		return nil, err
	}

	return &kvv1.DeleteRangeResponse{Deleted: deleted}, nil
}
