package grpc

import (
	"context"
	v1 "github.com/dimuska139/cacher/internal/api/grpc/gen/cacher/cache/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Logger interface {
	Error(msg string, args ...interface{})
}

type Storage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
	Delete(key string) error
}

type CacheServer struct {
	logger  Logger
	storage Storage
}

func NewCacheServer(
	logger Logger,
	storage Storage,
) *CacheServer {
	return &CacheServer{
		logger:  logger,
		storage: storage,
	}
}

// Get возвращает данные по ключу из кеша
func (s *CacheServer) Get(ctx context.Context, request *v1.GetRequest) (*v1.GetResponse, error) {
	data, err := s.storage.Get(request.GetKey())
	if err != nil {
		s.logger.Error("Can't get data from storage",
			"err", err,
			"key", request.GetKey())
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}

	return &v1.GetResponse{
		Value: data,
	}, nil
}

// Set записывает данные в кеш
func (s *CacheServer) Set(ctx context.Context, request *v1.SetRequest) (*v1.SetResponse, error) {
	err := s.storage.Set(request.GetKey(), request.GetValue(), time.Second*time.Duration(request.GetTtl()))
	if err != nil {
		s.logger.Error("Can't get save data to storage", "err", err)
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}

	return &v1.SetResponse{}, nil
}

// Delete удаляет данные из кеша
func (s *CacheServer) Delete(ctx context.Context, request *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	err := s.storage.Delete(request.GetKey())
	if err != nil {
		s.logger.Error("Can't delete data from storage",
			"err", err,
			"key", request.GetKey())
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}

	return &v1.DeleteResponse{}, nil
}
