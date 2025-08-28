package bnm

import (
	"context"
	"sync"
)

// Cache represents a simple key-value store for storing Response objects.
type Cache interface {
	// Set stores a Response in the cache associated with the given key.
	// It returns an error if the operation fails.
	Set(ctx context.Context, key string, res Response) error

	// Get retrieves a Response from the cache by its key.
	// It returns ErrNotFound if the key does not exist.
	Get(ctx context.Context, key string) (Response, error)
}

// MemoryCache is an in-memory implementation of the Cache interface.
// It is safe for concurrent use by multiple goroutines.
type MemoryCache struct {
	data map[string]Response
	mu   sync.RWMutex
}

var _ Cache = (*MemoryCache)(nil)

// NewMemoryCache creates and returns a new MemoryCache instance.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]Response),
	}
}

// Set stores a Response in the memory cache under the specified key.
// It overwrites any existing value for that key.
func (c *MemoryCache) Set(ctx context.Context, key string, res Response) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = res

	return nil
}

// Get retrieves a Response from the memory cache by key.
// If the key does not exist, it returns ErrNotFound.
func (c *MemoryCache) Get(ctx context.Context, key string) (Response, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if v, ok := c.data[key]; ok {
		return v, nil
	}

	return Response{}, ErrNotFound
}
