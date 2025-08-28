package bnm

import (
	"context"
	"fmt"
)

// Option configures a Client.
type Option func(*Client)

// GetRequestFunc defines a function signature for performing HTTP GET requests.
// It takes a context and URL and returns the response body or an error.
type GetRequestFunc func(ctx context.Context, url string) ([]byte, error)

// UnmarshalerFunc defines a function signature for unmarshaling a byte slice
// into a Response object.
type UnmarshalerFunc func([]byte) (Response, error)

// WarnFunc defines a function signature for logging non-critical errors.
type WarnFunc func(error)

// Client is used to fetch exchange rates from the National Bank of Moldova (BNM) API.
type Client struct {
	cache       Cache
	getRequest  GetRequestFunc
	unmarshaler UnmarshalerFunc
	warnError   WarnFunc
}

// NewClient creates a new Client instance with optional configuration.
// By default, it uses the standard HTTP client and the default unmarshaler.
//
// Example:
//
//	client := NewClient(
//	    WithCache(myCache),
//	    WithWarnError(logError),
//	)
func NewClient(opts ...Option) *Client {
	c := &Client{
		getRequest:  getRequestWithDefaultClient,
		unmarshaler: unmarshalResponse,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithCache sets a Cache implementation on the Client.
func WithCache(cache Cache) Option {
	return func(c *Client) { c.cache = cache }
}

// WithGetRequest sets a custom GetRequestFunc on the Client.
func WithGetRequest(fn GetRequestFunc) Option {
	return func(c *Client) { c.getRequest = fn }
}

// WithUnmarshaler sets a custom UnmarshalerFunc on the Client.
func WithUnmarshaler(u UnmarshalerFunc) Option {
	return func(c *Client) { c.unmarshaler = u }
}

// WithWarnError sets a WarnFunc on the Client to log non-critical errors.
func WithWarnError(fn WarnFunc) Option {
	return func(c *Client) { c.warnError = fn }
}

// Fetch retrieves exchange rates for a given query.
// It first checks the cache (if configured), then fetches from the BNM API,
// and finally stores the result in the cache.
//
// Returns an error if the request fails, if unmarshaling fails, or
// if there is an unexpected cache error.
//
// Example:
//
//	resp, err := client.Fetch(ctx, NewQuery(time.Now(), bnm.LANG_EN)
func (c *Client) Fetch(ctx context.Context, query Query) (Response, error) {
	if c.cache != nil {
		if cache, err := c.cache.Get(ctx, query.ID()); err == nil {
			return cache, nil
		} else if err != ErrNotFound {
			return Response{}, fmt.Errorf("get cache: %w", err)
		}
	}

	data, err := c.getRequest(ctx, query.RequestURL())
	if err != nil {
		return Response{}, fmt.Errorf("get request: %w", err)
	}

	res, err := c.unmarshaler(data)
	if err != nil {
		return Response{}, fmt.Errorf("parse body: %w", err)
	}

	if c.cache != nil {
		if err := c.cache.Set(ctx, query.ID(), res); err != nil {
			c.warnError(fmt.Errorf("set cache: %w", err))
		}
	}

	return res, nil
}
