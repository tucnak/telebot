package storage

import (
	"context"
	"time"
)

// Storage defines a global interface for key-value storage systems.
// Uses in-memory storage by default.
type Storage interface {
	// Get gets a value by key.
	Get(key string) (interface{}, error)

	// Set sets a value by key with an expiration time.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Delete deletes a value by key.
	Delete(key string) error
}
