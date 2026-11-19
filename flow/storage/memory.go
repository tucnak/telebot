package storage

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

// Memory defines in-memory storage with use [sync.Map].
type Memory struct {
	data sync.Map
}

// NewMemory constructor [Memory].
func NewMemory() Storage {
	return &Memory{}
}

// Get gets value by key from memory storage.
func (m *Memory) Get(key string) (interface{}, error) {
	value, ok := m.data.Load(key)
	if !ok {
		return nil, ErrKeyNotFound
	}
	return value, nil
}

// Set sets value by key to memory storage.
func (m *Memory) Set(_ context.Context, key string, value interface{}, _ time.Duration) error {
	m.data.Store(key, value)
	return nil
}

// Delete deletes value by key from memory storage.
func (m *Memory) Delete(key string) error {
	m.data.Delete(key)
	return nil
}
