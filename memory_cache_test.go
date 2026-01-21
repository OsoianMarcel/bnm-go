package bnm_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/OsoianMarcel/bnm-go/v2"
)

func TestNewMemoryCache(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		capacity    int
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid capacity",
			capacity: 10,
			wantErr:  false,
		},
		{
			name:        "zero capacity",
			capacity:    0,
			wantErr:     true,
			errContains: "capacity must be positive",
		},
		{
			name:        "negative capacity",
			capacity:    -5,
			wantErr:     true,
			errContains: "capacity must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache, err := bnm.NewMemoryCache(tt.capacity)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && err.Error() != tt.errContains {
					t.Errorf("error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cache == nil {
				t.Fatal("expected non-nil cache")
			}
		})
	}
}

func TestMemoryCache_SetGet(t *testing.T) {
	t.Parallel()

	cache, err := bnm.NewMemoryCache(5)
	if err != nil {
		t.Fatalf("failed to create memory cache: %v", err)
	}

	ctx := t.Context()

	tests := []struct {
		name     string
		key      string
		response bnm.Response
	}{
		{
			name: "simple set and get",
			key:  "key1",
			response: bnm.Response{
				Name: "Test Currency",
			},
		},
		{
			name: "another entry",
			key:  "key2",
			response: bnm.Response{
				Name: "Another Currency",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cache.Set(ctx, tt.key, tt.response); err != nil {
				t.Fatalf("Set() error = %v", err)
			}

			got, err := cache.Get(ctx, tt.key)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}

			if got.Name != tt.response.Name {
				t.Errorf("Get() Name = %v, want %v", got.Name, tt.response.Name)
			}
		})
	}
}

func TestMemoryCache_Update(t *testing.T) {
	t.Parallel()

	cache, err := bnm.NewMemoryCache(3)
	if err != nil {
		t.Fatalf("failed to create memory cache: %v", err)
	}

	ctx := t.Context()
	key := "update-key"

	// Set initial value
	initial := bnm.Response{Name: "initial"}
	if err := cache.Set(ctx, key, initial); err != nil {
		t.Fatalf("Set() initial error = %v", err)
	}

	// Update value
	updated := bnm.Response{Name: "updated"}
	if err := cache.Set(ctx, key, updated); err != nil {
		t.Fatalf("Set() update error = %v", err)
	}

	// Verify updated value
	got, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.Name != "updated" {
		t.Errorf("Get() Name = %v, want %v", got.Name, "updated")
	}
}

func TestMemoryCache_GetNotFound(t *testing.T) {
	t.Parallel()

	cache, err := bnm.NewMemoryCache(3)
	if err != nil {
		t.Fatalf("failed to create memory cache: %v", err)
	}

	ctx := t.Context()

	_, err = cache.Get(ctx, "nonexistent-key")
	if !errors.Is(err, bnm.ErrNotFound) {
		t.Errorf("Get() error = %v, want %v", err, bnm.ErrNotFound)
	}
}

func TestMemoryCache_LRUEviction(t *testing.T) {
	t.Parallel()

	cache, err := bnm.NewMemoryCache(2)
	if err != nil {
		t.Fatalf("failed to create memory cache: %v", err)
	}

	ctx := t.Context()

	// Fill cache to capacity
	if err := cache.Set(ctx, "key1", bnm.Response{Name: "one"}); err != nil {
		t.Fatalf("Set() key1 error = %v", err)
	}
	if err := cache.Set(ctx, "key2", bnm.Response{Name: "two"}); err != nil {
		t.Fatalf("Set() key2 error = %v", err)
	}

	// Update key1 to make it most recently used (key2 becomes LRU)
	if err := cache.Set(ctx, "key1", bnm.Response{Name: "one-updated"}); err != nil {
		t.Fatalf("Set() key1 update error = %v", err)
	}

	// Add key3, should evict key2 (LRU)
	if err := cache.Set(ctx, "key3", bnm.Response{Name: "three"}); err != nil {
		t.Fatalf("Set() key3 error = %v", err)
	}

	// Verify key2 was evicted
	_, err = cache.Get(ctx, "key2")
	if !errors.Is(err, bnm.ErrNotFound) {
		t.Errorf("Get() key2 error = %v, want %v", err, bnm.ErrNotFound)
	}

	// Verify key1 still exists with updated value
	got, err := cache.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Get() key1 error = %v", err)
	}
	if got.Name != "one-updated" {
		t.Errorf("Get() key1 Name = %v, want %v", got.Name, "one-updated")
	}

	// Verify key3 exists
	got, err = cache.Get(ctx, "key3")
	if err != nil {
		t.Fatalf("Get() key3 error = %v", err)
	}
	if got.Name != "three" {
		t.Errorf("Get() key3 Name = %v, want %v", got.Name, "three")
	}
}

func TestMemoryCache_LRUEvictionOnGet(t *testing.T) {
	t.Parallel()

	cache, err := bnm.NewMemoryCache(2)
	if err != nil {
		t.Fatalf("failed to create memory cache: %v", err)
	}

	ctx := t.Context()

	// Fill cache
	if err := cache.Set(ctx, "key1", bnm.Response{Name: "one"}); err != nil {
		t.Fatalf("Set() key1 error = %v", err)
	}
	if err := cache.Set(ctx, "key2", bnm.Response{Name: "two"}); err != nil {
		t.Fatalf("Set() key2 error = %v", err)
	}

	// Access key1 to make it most recently used
	if _, err := cache.Get(ctx, "key1"); err != nil {
		t.Fatalf("Get() key1 error = %v", err)
	}

	// Add key3, should evict key2 (now LRU)
	if err := cache.Set(ctx, "key3", bnm.Response{Name: "three"}); err != nil {
		t.Fatalf("Set() key3 error = %v", err)
	}

	// Verify key2 was evicted
	_, err = cache.Get(ctx, "key2")
	if !errors.Is(err, bnm.ErrNotFound) {
		t.Errorf("Get() key2 error = %v, want %v", err, bnm.ErrNotFound)
	}

	// Verify key1 and key3 still exist
	if _, err := cache.Get(ctx, "key1"); err != nil {
		t.Errorf("Get() key1 error = %v", err)
	}
	if _, err := cache.Get(ctx, "key3"); err != nil {
		t.Errorf("Get() key3 error = %v", err)
	}
}

func TestMemoryCache_Concurrent(t *testing.T) {
	t.Parallel()

	cache, err := bnm.NewMemoryCache(100)
	if err != nil {
		t.Fatalf("failed to create memory cache: %v", err)
	}

	ctx := t.Context()
	var wg sync.WaitGroup
	goroutines := 50
	opsPerGoroutine := 100

	// Concurrent writes
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := string(rune('A' + (id+j)%26))
				resp := bnm.Response{Name: key}
				if err := cache.Set(ctx, key, resp); err != nil {
					t.Errorf("Set() error = %v", err)
				}
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := string(rune('A' + (id+j)%26))
				cache.Get(ctx, key) // Error is expected for some keys
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkMemoryCache_Set(b *testing.B) {
	cache, _ := bnm.NewMemoryCache(1000)
	ctx := b.Context()
	resp := bnm.Response{Name: "benchmark"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(ctx, "key", resp)
	}
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	cache, _ := bnm.NewMemoryCache(1000)
	ctx := b.Context()
	resp := bnm.Response{Name: "benchmark"}
	cache.Set(ctx, "key", resp)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(ctx, "key")
	}
}

func BenchmarkMemoryCache_SetParallel(b *testing.B) {
	cache, _ := bnm.NewMemoryCache(1000)
	ctx := b.Context()
	resp := bnm.Response{Name: "benchmark"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(ctx, string(rune('A'+(i%26))), resp)
			i++
		}
	})
}

func BenchmarkMemoryCache_GetParallel(b *testing.B) {
	cache, _ := bnm.NewMemoryCache(1000)
	ctx := b.Context()

	// Pre-populate cache
	for i := 0; i < 26; i++ {
		cache.Set(ctx, string(rune('A'+i)), bnm.Response{Name: string(rune('A' + i))})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Get(ctx, string(rune('A'+(i%26))))
			i++
		}
	})
}
