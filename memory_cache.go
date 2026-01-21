package bnm

import (
	"container/list"
	"context"
	"errors"
	"sync"
)

// MemoryCache is an in-memory LRU cache implementation.
// It is safe for concurrent use by multiple goroutines.
type MemoryCache struct {
	capacity int
	data     map[string]*list.Element
	ll       *list.List
	mu       sync.Mutex
}

type cacheEntry struct {
	key   string
	value Response
}

var _ Cache = (*MemoryCache)(nil)

// NewMemoryCache creates and returns a new MemoryCache instance.
func NewMemoryCache(capacity int) (*MemoryCache, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}

	return &MemoryCache{
		capacity: capacity,
		data:     make(map[string]*list.Element, capacity),
		ll:       list.New(),
	}, nil
}

// Set stores a Response in the memory cache under the specified key.
// It overwrites any existing value for that key.
func (c *MemoryCache) Set(ctx context.Context, key string, res Response) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.data[key]; found {
		c.ll.MoveToFront(elem)
		elem.Value.(*cacheEntry).value = res
		return nil
	}

	elem := c.ll.PushFront(&cacheEntry{key, res})
	c.data[key] = elem

	if c.ll.Len() > c.capacity {
		c.evictOldest()
	}

	return nil
}

// Get retrieves a Response from the memory cache by key.
// If the key does not exist, it returns ErrNotFound.
// Note: This updates the LRU order (moves item to front).
func (c *MemoryCache) Get(ctx context.Context, key string) (Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, found := c.data[key]
	if !found {
		return Response{}, ErrNotFound
	}

	c.ll.MoveToFront(elem)
	return elem.Value.(*cacheEntry).value, nil
}

// evictOldest removes the least recently used item.
// Must be called with mutex held.
func (c *MemoryCache) evictOldest() {
	backElem := c.ll.Back()
	if backElem != nil {
		backEntry := backElem.Value.(*cacheEntry)
		delete(c.data, backEntry.key)
		c.ll.Remove(backElem)
	}
}
