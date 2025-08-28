package bnm_test

import (
	"testing"

	"github.com/OsoianMarcel/bnm-go/v2"
)

func TestMemoryCache_SetGet(t *testing.T) {
	cache := bnm.NewMemoryCache()

	if err := cache.Set(t.Context(), "1", bnm.Response{Name: "one"}); err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	res, err := cache.Get(t.Context(), "1")
	if err != nil {
		t.Fatalf("failed to get cache after set: %v", err)
	}

	if res.Name != "one" {
		t.Errorf("expected Name 'one', got '%s'", res.Name)
	}
}

func TestMemoryCache_GetNotFound(t *testing.T) {
	cache := bnm.NewMemoryCache()

	_, err := cache.Get(t.Context(), "inexistent")
	if err != bnm.ErrNotFound {
		t.Errorf("expected error ErrNotFound, got %v", err)
	}
}
