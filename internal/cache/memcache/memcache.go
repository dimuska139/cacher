package memcache

import (
	"fmt"
	"time"
)

//go:generate mockgen -source=memcache.go -destination=./memcache_mock.go -package=memcache

// Memcacher интерфейс для библиотеки-клиента Memcache
type Memcacher interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expiration int64) error
	Delete(key string) error
}

// MemcacheStorage реализация кеша через Memcache
type MemcacheStorage struct {
	memcacheClient Memcacher
}

// NewMemcacheStorage создаёт реализацию кеша через Memcache
func NewMemcacheStorage(memcacheClient Memcacher) *MemcacheStorage {
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
