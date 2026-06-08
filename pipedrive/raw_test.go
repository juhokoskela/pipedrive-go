package pipedrive

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestRawClient_Do_ReturnsResponseTooLargeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(srv.Close)

	httpClient := NewHTTPClient(Config{
		HTTPClient:      srv.Client(),
		MaxResponseSize: 4,
		RetryPolicy:     &RetryPolicy{MaxAttempts: 1},
	})

	raw, err := NewRawClient(srv.URL, httpClient)
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}

	err = raw.Do(context.Background(), http.MethodGet, "/x", nil, nil, nil)
	var tooLarge *ResponseTooLargeError
	if !errors.As(err, &tooLarge) {
		t.Fatalf("expected ResponseTooLargeError, got %T (%v)", err, err)
	}
	if tooLarge.Limit != 4 {
		t.Fatalf("expected limit=4, got %d", tooLarge.Limit)
	}
}

func TestRawClient_Do_ReturnsAPIErrorWhenLargeErrorBodyExceedsLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req_big")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"payload too large for the response cap"}`))
	}))
	t.Cleanup(srv.Close)

	httpClient := NewHTTPClient(Config{
		HTTPClient:      srv.Client(),
		MaxResponseSize: 4,
		RetryPolicy:     &RetryPolicy{MaxAttempts: 1},
	})

	raw, err := NewRawClient(srv.URL, httpClient)
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
	if apiErr.RequestID != "req_big" {
		t.Fatalf("expected request id req_big, got %q", apiErr.RequestID)
	}
}

func TestRawClient_Do_ReturnsRateLimitErrorWhenLarge429BodyExceedsLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.Header().Set("X-RateLimit-Limit", "10")
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message":"payload too large for the response cap"}`))
	}))
	t.Cleanup(srv.Close)

	httpClient := NewHTTPClient(Config{
		HTTPClient:      srv.Client(),
		MaxResponseSize: 4,
		RetryPolicy:     &RetryPolicy{MaxAttempts: 1},
	})

	raw, err := NewRawClient(srv.URL, httpClient)
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
		t.Fatalf("unexpected rate limit headers: limit=%d remaining=%d", rlErr.Limit, rlErr.Remaining)
	}
}

func TestNewRawClient_RejectsInvalidBaseURL(t *testing.T) {
	t.Parallel()

	if _, err := NewRawClient("://bad", nil); err == nil {
		t.Fatalf("expected parse error")
	}
	if _, err := NewRawClient("/relative", nil); err == nil {
		t.Fatalf("expected missing scheme/host error")
	}
}

func TestRawClient_Do_NilClient(t *testing.T) {
	t.Parallel()

	var raw *RawClient
	if err := raw.Do(context.Background(), http.MethodGet, "/x", nil, nil, nil); err == nil {
		t.Fatalf("expected nil client error")
	}
}

func TestRawClient_Do_SendsQueryAndJSONBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/base/items" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["tag"]; len(got) != 2 || got[0] != "one" || got[1] != "two" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("unexpected accept: %q", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}

		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["name"] != "item" {
			t.Fatalf("unexpected payload: %#v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(srv.Close)

	raw, err := NewRawClient(srv.URL+"/base/", srv.Client())
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}

	var out struct {
		Ok bool `json:"ok"`
	}
	query := url.Values{"tag": {"one", "two"}}
	if err := raw.Do(context.Background(), http.MethodPost, "/items", query, map[string]string{"name": "item"}, &out); err != nil {
		t.Fatalf("Do error: %v", err)
	}
	if !out.Ok {
		t.Fatalf("expected ok")
	}
}
