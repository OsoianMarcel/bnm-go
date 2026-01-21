package bnm

import (
	"context"
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
