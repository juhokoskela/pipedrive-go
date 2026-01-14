package pipedrive

import (
	"net/http"
	"testing"
	"time"
)

func TestAPIErrorFromResponse_ExtractsRequestIDAndMessage(t *testing.T) {
	t.Parallel()

	h := make(http.Header)
	h.Set("X-Request-Id", "req_123")
	resp := &http.Response{
		StatusCode: 400,
		Header:     h,
	}
	body := []byte(`{"success":false,"error":"bad request","error_info":"details"}`)

	err := apiErrorFromResponse(resp, body)
	if err.Status != 400 {
		t.Fatalf("expected status=400, got %d", err.Status)
	}
	if err.RequestID != "req_123" {
		t.Fatalf("expected request id req_123, got %q", err.RequestID)
	}
	if err.Message != "bad request" {
		t.Fatalf("expected message %q, got %q", "bad request", err.Message)
	}
}

func TestRateLimitErrorFromResponse_ParsesHeaders(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	h := make(http.Header)
	h.Set("Retry-After", "2")
	h.Set("X-RateLimit-Limit", "10")
	h.Set("X-RateLimit-Remaining", "0")
	h.Set("X-RateLimit-Reset", "1735689605") // 2025-01-01T00:00:05Z
	resp := &http.Response{
		StatusCode: 429,
		Header:     h,
	}

	body := []byte(`{"code":"rate_limit","message":"Too many requests"}`)
	err := rateLimitErrorFromResponse(resp, body, now)

	if err.Status != 429 {
		t.Fatalf("expected status=429, got %d", err.Status)
	}
	if err.Code != "rate_limit" {
		t.Fatalf("expected code=rate_limit, got %q", err.Code)
	}
	if err.Message != "Too many requests" {
		t.Fatalf("expected message %q, got %q", "Too many requests", err.Message)
	}
	if err.RetryAfter != 2*time.Second {
		t.Fatalf("expected retry-after 2s, got %s", err.RetryAfter)
	}
	if err.Limit != 10 {
		t.Fatalf("expected limit=10, got %d", err.Limit)
	}
	if err.Remaining != 0 {
		t.Fatalf("expected remaining=0, got %d", err.Remaining)
	}
	if !err.Reset.Equal(time.Date(2025, 1, 1, 0, 0, 5, 0, time.UTC)) {
		t.Fatalf("unexpected reset time: %s", err.Reset)
	}
}
