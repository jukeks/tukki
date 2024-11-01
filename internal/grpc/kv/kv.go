package kv

import (
	"context"
	"io"

	"github.com/jukeks/tukki/internal/db"
	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
)

type Pair struct {
	Key   string
	Value string
}

type DB interface {
	Get(key string) (string, error)
	GetRange(min, max string) (db.KeyValueIterator, error)
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

func (s *KVServer) QueryRange(req *kvv1.QueryRangeRequest, resp kvv1.KvService_QueryRangeServer) error {
	cursor, err := s.db.GetRange(req.Min, req.Max)
	if err != nil {
		return err
	}

	for {
		pair, err := cursor.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		err = resp.Send(&kvv1.QueryRangeResponse{
			Pair: &kvv1.KvPair{
				Key:   pair.Key,
				Value: pair.Value,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
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
