package bnm_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/OsoianMarcel/bnm-go/v2"
)

type mockCache struct {
	getFunc func(context.Context, string) (bnm.Response, error)
	setFunc func(context.Context, string, bnm.Response) error
}

func (m *mockCache) Get(ctx context.Context, key string) (bnm.Response, error) {
	return m.getFunc(ctx, key)
}
func (m *mockCache) Set(ctx context.Context, key string, r bnm.Response) error {
	return m.setFunc(ctx, key, r)
}

func dummyQuery() bnm.Query {
	return bnm.NewQuery(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), bnm.LANG_EN)
}

func TestFetch_CacheHit(t *testing.T) {
	cache := &mockCache{
		getFunc: func(_ context.Context, _ string) (bnm.Response, error) {
			return bnm.Response{Date: "2025-01-01"}, nil
		},
	}
	client := bnm.NewClient(bnm.WithCache(cache))

	resp, err := client.Fetch(t.Context(), dummyQuery())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Date != "2025-01-01" {
		t.Fatalf("expected cached response, got %v", resp)
	}
}

func TestFetch_CacheMissThenFetch(t *testing.T) {
	cache := &mockCache{
		getFunc: func(_ context.Context, _ string) (bnm.Response, error) {
			return bnm.Response{}, bnm.ErrNotFound
		},
		setFunc: func(_ context.Context, _ string, _ bnm.Response) error {
			return nil
		},
	}
	client := bnm.NewClient(
		bnm.WithCache(cache),
		bnm.WithGetRequest(func(_ context.Context, _ string) ([]byte, error) {
			return []byte(`{"date":"2025-01-01"}`), nil
		}),
		bnm.WithUnmarshaler(func(b []byte) (bnm.Response, error) {
			return bnm.Response{Date: "2025-01-01"}, nil
		}),
	)

	resp, err := client.Fetch(t.Context(), dummyQuery())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Date != "2025-01-01" {
		t.Fatalf("expected response from API, got %v", resp)
	}
}

func TestFetch_CacheError(t *testing.T) {
	cache := &mockCache{
		getFunc: func(_ context.Context, _ string) (bnm.Response, error) {
			return bnm.Response{}, errors.New("db down")
		},
	}
	client := bnm.NewClient(bnm.WithCache(cache))

	_, err := client.Fetch(t.Context(), dummyQuery())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetch_GetRequestError(t *testing.T) {
	cache := &mockCache{
		getFunc: func(_ context.Context, _ string) (bnm.Response, error) {
			return bnm.Response{}, bnm.ErrNotFound
		},
	}
	client := bnm.NewClient(
		bnm.WithCache(cache),
		bnm.WithGetRequest(func(_ context.Context, _ string) ([]byte, error) {
			return nil, errors.New("network error")
		}),
	)

	_, err := client.Fetch(t.Context(), dummyQuery())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetch_UnmarshalError(t *testing.T) {
	cache := &mockCache{
		getFunc: func(_ context.Context, _ string) (bnm.Response, error) {
			return bnm.Response{}, bnm.ErrNotFound
		},
	}
	client := bnm.NewClient(
		bnm.WithCache(cache),
		bnm.WithGetRequest(func(_ context.Context, _ string) ([]byte, error) {
			return []byte("invalid"), nil
		}),
		bnm.WithUnmarshaler(func(_ []byte) (bnm.Response, error) {
			return bnm.Response{}, errors.New("bad json")
		}),
	)

	_, err := client.Fetch(t.Context(), dummyQuery())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetch_CacheSetErrorWarns(t *testing.T) {
	cache := &mockCache{
		getFunc: func(_ context.Context, _ string) (bnm.Response, error) {
			return bnm.Response{}, bnm.ErrNotFound
		},
		setFunc: func(_ context.Context, _ string, _ bnm.Response) error {
			return errors.New("redis timeout")
		},
	}

	warnCalled := false
	client := bnm.NewClient(
		bnm.WithCache(cache),
		bnm.WithGetRequest(func(_ context.Context, _ string) ([]byte, error) {
			return []byte(`{"date":"2025-01-01"}`), nil
		}),
		bnm.WithUnmarshaler(func(_ []byte) (bnm.Response, error) {
			return bnm.Response{Date: "2025-01-01"}, nil
		}),
		bnm.WithWarnError(func(err error) {
			warnCalled = true
		}),
	)

	resp, err := client.Fetch(t.Context(), dummyQuery())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Date != "2025-01-01" {
		t.Fatalf("expected valid response, got %v", resp)
	}
	if !warnCalled {
		t.Fatal("expected warnError to be called, but it wasn't")
	}
}
