package embedded

import (
	"runtime"
	"sync"
	"time"
)

// item элемент кеша
type item struct {
	Value      []byte
	Expiration int64
}

// IsExpired проверяет, не истекло ли время жизни элемента кеша
func (i *item) IsExpired() bool {
	return i.Expiration > 0 && i.Expiration < time.Now().UnixNano()
}

// EmbeddedStorage кеш внутри памяти приложения
type EmbeddedStorage struct {
	items           map[string]item
	cleanupInterval time.Duration
	mx              sync.RWMutex
	stopCleaning    chan bool
}

// NewEmbeddedStorage создаёт кеш внутри памяти приложения
func NewEmbeddedStorage(cleanupInterval time.Duration) *EmbeddedStorage {
	cache := &EmbeddedStorage{
		items:           make(map[string]item),
		cleanupInterval: cleanupInterval,
		stopCleaning:    make(chan bool),
	}

	go cache.cleaner()

	runtime.SetFinalizer(cache, finalizer)
	return cache
}

// finalizer корректно завершает работу EmbedStorage, останавливая функцию очистки кеша
func finalizer(c *EmbeddedStorage) {
	c.stopCleaning <- true
}

// Get возвращает закешированные данные
func (s *EmbeddedStorage) Get(key string) ([]byte, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	i, found := s.items[key]
	if !found || i.IsExpired() {
		return nil, nil
	}

	return i.Value, nil
}

// Set записывает информацию в кеш. Если запись в кеше уже есть, то она обновится
func (s *EmbeddedStorage) Set(key string, value []byte, ttl time.Duration) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	var expireAt int64
	if ttl > 0 {
		expireAt = time.Now().Add(ttl).UnixNano()
	}

	s.items[key] = item{
		Value:      value,
		Expiration: expireAt,
	}
	return nil
}

// Delete удаляет запись из кеша по ключу
func (s *EmbeddedStorage) Delete(key string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.items, key)
	return nil
}

// deleteExpired удаляет из кеша записи с истёкшим временем жизни
func (s *EmbeddedStorage) deleteExpired() {
	s.mx.Lock()
	defer s.mx.Unlock()
	for k, v := range s.items {
		if v.IsExpired() {
			delete(s.items, k)
		}
	}
}

// cleaner следит за актуальностью данных в кеше
func (s *EmbeddedStorage) cleaner() {
	ticker := time.NewTicker(s.cleanupInterval)
	for {
		select {
		case <-s.stopCleaning:
			ticker.Stop()
			return
		case <-ticker.C:
			s.deleteExpired()
		}
	}
}
