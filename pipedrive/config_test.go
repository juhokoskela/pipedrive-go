package pipedrive

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewHTTPClient_MiddlewareOrder(t *testing.T) {
	t.Parallel()

	var calls []string

	base := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calls = append(calls, "base")
		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	mw1 := func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			calls = append(calls, "mw1")
			return next.RoundTrip(req)
		})
	}
	mw2 := func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			calls = append(calls, "mw2")
			return next.RoundTrip(req)
		})
	}

	httpClient := NewHTTPClient(Config{
		HTTPClient:  &http.Client{Transport: base},
		Middleware:  []Middleware{mw1, mw2},
		RetryPolicy: &RetryPolicy{MaxAttempts: 1},
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	resp, err := httpClient.Transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()

	want := []string{"mw1", "mw2", "base"}
	if len(calls) != len(want) {
		t.Fatalf("unexpected call count: got %v want %v", calls, want)
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf("unexpected order: got %v want %v", calls, want)
		}
	}
}

func TestNewHTTPClient_DefaultRetryEnabled(t *testing.T) {
	t.Parallel()

	httpClient := NewHTTPClient(Config{
		HTTPClient: &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			}, nil
		})},
	})

	if _, ok := httpClient.Transport.(*retryTransport); !ok {
		t.Fatalf("expected default transport to include retry transport, got %T", httpClient.Transport)
	}
}
