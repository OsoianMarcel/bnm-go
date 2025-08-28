package bnm

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeHttpClient struct {
	res *http.Response
	err error
}

func (f *fakeHttpClient) Do(_ *http.Request) (*http.Response, error) {
	return f.res, f.err
}

type errReader struct{}

func (errReader) Read(p []byte) (int, error) { return 0, errors.New("read error") }

func TestGetRequest(t *testing.T) {
	tests := []struct {
		name    string
		client  clientDoer
		url     string
		want    string
		wantErr string
	}{
		{
			name:    "invalid url",
			client:  &fakeHttpClient{},
			url:     "http://%41:8080/",
			wantErr: "create request",
		},
		{
			name:    "do request error",
			client:  &fakeHttpClient{err: errors.New("boom")},
			url:     "http://example.com",
			wantErr: "do request",
		},
		{
			name: "non-200 status",
			client: &fakeHttpClient{res: &http.Response{
				StatusCode: http.StatusTeapot,
				Body:       io.NopCloser(strings.NewReader("")),
			}},
			url:     "http://example.com",
			wantErr: "status code",
		},
		{
			name: "body read error",
			client: &fakeHttpClient{res: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(errReader{}),
			}},
			url:     "http://example.com",
			wantErr: "read body",
		},
		{
			name: "success",
			client: &fakeHttpClient{res: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("hello")),
			}},
			url:  "http://example.com",
			want: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRequest(t.Context(), tt.client, tt.url)

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(got) != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGetRequestWithDefaultClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	defer ts.Close()

	body, err := getRequestWithDefaultClient(t.Context(), ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(body) != "ok" {
		t.Errorf("expected %q, got %q", "ok", string(body))
	}
}
