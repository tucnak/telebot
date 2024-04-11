package telebot

import (
	"sync"
)

// Cache is a provider for cache
//
// All cache providers must implement methods, which work with storage
// to store a cache of user state, inline results
// for a Telegram bot.
type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error

	Keys() []string
}

// InMemoryCache is a provider of in memory cache.
//
// Would be enabled by default
type inMemoryCache struct {
	lock sync.RWMutex
	data map[string]interface{}

	keys         []string // To maintain order of keys
	currentIndex int      // To keep track of current index
}

func NewInMemoryCache() Cache {
	return &inMemoryCache{}
}

func (m *inMemoryCache) Get(key string) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.data[key], nil
}

func (m *inMemoryCache) Set(key string, value interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	m.data[key] = value
	m.keys = append(m.keys, key) // Update the keys slice

	return nil
}

func (m *inMemoryCache) Keys() []string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.keys
}
