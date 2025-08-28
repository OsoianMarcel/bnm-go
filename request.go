package bnm

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type clientDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func getRequest(ctx context.Context, client clientDoer, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("create request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("read body: %w", err)
	}

	return body, nil
}

func getRequestWithDefaultClient(ctx context.Context, url string) ([]byte, error) {
	return getRequest(ctx, http.DefaultClient, url)
}
