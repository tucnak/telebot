package telebot

import (
	"errors"
	"sync"
)

type CacheKind string

const (
	CacheUserContext CacheKind = "user_context"
)

// Cache is a provider for cache
//
// All cache providers must implement methods, which work with storage
// to store a cache of user state, inline results
// for a Telegram bot.
type Cache interface {
	Get(kind CacheKind, key string) (interface{}, error)
	Put(kind CacheKind, key string, value interface{}) error
	Clear(kind CacheKind, key string) error

	Keys(kind CacheKind) []string
}

// InMemoryCache is a provider of in memory cache.
//
// Would be enabled by default
type inMemoryCache struct {
	lock sync.RWMutex
	data map[CacheKind]map[string]interface{}

	keys map[CacheKind][]string // To maintain order of keys
}

func NewInMemoryCache() Cache {
	return &inMemoryCache{}
}

func (m *inMemoryCache) Get(kind CacheKind, key string) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if _, ok := m.data[kind]; !ok {
		return nil, errors.New("telebot: cache kind is not found")
	}

	if _, ok := m.data[kind][key]; !ok {
		return nil, errors.New("telebot: cache key is not found")
	}

	return m.data[kind][key], nil
}

func (m *inMemoryCache) Put(kind CacheKind, key string, value interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.data[kind]; !ok {
		m.data[kind] = make(map[string]interface{})
		m.keys[kind] = make([]string, 1)
	}

	m.data[kind][key] = value
	m.keys[kind] = append(m.keys[kind], key) // Update the keys slice

	return nil
}

func (m *inMemoryCache) Clear(kind CacheKind, key string) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if _, ok := m.data[kind]; !ok {
		return errors.New("telebot: cache kind is not found")
	}

	if _, ok := m.data[kind][key]; !ok {
		return errors.New("telebot: cache key not found")
	}

	delete(m.data[kind], key)     // delete value
	keysRemove(m.keys[kind], key) // delete key

	return nil
}

func (m *inMemoryCache) Keys(kind CacheKind) []string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.keys[kind]
}

func keysRemove(keys []string, remove string) []string {
	for i, k := range keys {
		if k == remove {
			return append(keys[:i], keys[i+1:]...)
		}
	}

	return keys
}
