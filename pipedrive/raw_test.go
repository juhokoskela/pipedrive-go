package pipedrive

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRawClient_Do_DecodesJSONResponse(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(srv.Close)

	raw, err := NewRawClient(srv.URL, srv.Client())
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}

	var out struct {
		Ok bool `json:"ok"`
	}
	if err := raw.Do(context.Background(), http.MethodGet, "/test", nil, nil, &out); err != nil {
		t.Fatalf("Do error: %v", err)
	}
	if !out.Ok {
		t.Fatalf("expected ok=true")
	}
}

func TestRawClient_Do_AppliesRequestOptions(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(srv.Close)

	raw, err := NewRawClient(srv.URL, srv.Client())
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}

	var out struct {
		Ok bool `json:"ok"`
	}
	if err := raw.Do(context.Background(), http.MethodGet, "/", nil, nil, &out, WithHeader("X-Test", "1")); err != nil {
		t.Fatalf("Do error: %v", err)
	}
	if !out.Ok {
		t.Fatalf("expected ok=true")
	}
}

func TestRawClient_Do_ReturnsAPIErrorOnNon2xx(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req_abc")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad request"}`))
	}))
	t.Cleanup(srv.Close)

	raw, err := NewRawClient(srv.URL, srv.Client())
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}

	err = raw.Do(context.Background(), http.MethodGet, "/x", nil, nil, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T (%v)", err, err)
	}
	if apiErr.Status != http.StatusBadRequest {
		t.Fatalf("expected status=400, got %d", apiErr.Status)
	}
	if apiErr.RequestID != "req_abc" {
		t.Fatalf("expected request id req_abc, got %q", apiErr.RequestID)
	}
}

func TestRawClient_Do_ReturnsRateLimitErrorOn429(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.Header().Set("X-RateLimit-Limit", "10")
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.Header().Set("X-RateLimit-Reset", "1735689605")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"code":"rate_limit","message":"Too many requests"}`))
	}))
	t.Cleanup(srv.Close)

	raw, err := NewRawClient(srv.URL, srv.Client())
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}

	err = raw.Do(context.Background(), http.MethodGet, "/x", nil, nil, nil)
	var rlErr *RateLimitError
	if !errors.As(err, &rlErr) {
		t.Fatalf("expected RateLimitError, got %T (%v)", err, err)
	}
	if rlErr.RetryAfter != 1*time.Second {
		t.Fatalf("expected retry-after 1s, got %s", rlErr.RetryAfter)
	}
	if rlErr.Limit != 10 || rlErr.Remaining != 0 {
		t.Fatalf("unexpected limit headers: limit=%d remaining=%d", rlErr.Limit, rlErr.Remaining)
	}
}
