package memcache

import (
	"fmt"
	"github.com/dimuska139/cacher/libs/memcache"
	"time"
)

type MemcacheStorage struct {
	memcacheClient *memcache.Client
}

// NewMemcacheStorage создаёт реализацию кеша через Memcache
func NewMemcacheStorage(memcacheClient *memcache.Client) *MemcacheStorage {
	return &MemcacheStorage{
		memcacheClient: memcacheClient,
	}
}

// Get возвращает закешированные данные
func (s *MemcacheStorage) Get(key string) ([]byte, error) {
	value, err := s.memcacheClient.Get(key)
	if err != nil {
		return nil, fmt.Errorf("can't get data from memcache: %w", err)
	}

	return value, nil
}

// Set записывает информацию в кеш. Если запись в кеше уже есть, то она обновится
func (s *MemcacheStorage) Set(key string, value []byte, ttl time.Duration) error {
	err := s.memcacheClient.Set(key, value, int64(ttl.Seconds()))
	if err != nil {
		return fmt.Errorf("can't get write data to memcache: %w", err)
	}

	return nil
}

// Delete удаляет запись из кеша по ключу
func (s *MemcacheStorage) Delete(key string) error {
	err := s.memcacheClient.Delete(key)
	if err != nil {
		return fmt.Errorf("can't delete data from memcache: %w", err)
	}

	return nil
}
