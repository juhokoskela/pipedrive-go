package pipedrive

import (
	"errors"
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

	err := APIErrorFromResponse(resp, body)
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
	err := RateLimitErrorFromResponse(resp, body, now)

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

func TestAPIError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  *APIError
		want string
	}{
		{
			name: "nil",
			err:  nil,
			want: "pipedrive: <nil>",
		},
		{
			name: "code and message",
			err:  &APIError{Status: 400, Code: "bad_request", Message: "Bad request"},
			want: "pipedrive: http 400 bad_request: Bad request",
		},
		{
			name: "message only",
			err:  &APIError{Status: 404, Message: "Not found"},
			want: "pipedrive: http 404: Not found",
		},
		{
			name: "code only",
			err:  &APIError{Status: 500, Code: "server_error"},
			want: "pipedrive: http 500 server_error",
		},
		{
			name: "status only",
			err:  &APIError{Status: 502},
			want: "pipedrive: http 502",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRateLimitError_ErrorAndUnwrap(t *testing.T) {
	t.Parallel()

	if got := (*RateLimitError)(nil).Error(); got != "pipedrive: rate limit" {
		t.Fatalf("unexpected nil error string: %q", got)
	}
	if got := (&RateLimitError{}).Error(); got != "pipedrive: rate limit" {
		t.Fatalf("unexpected empty error string: %q", got)
	}

	apiErr := &APIError{Status: 429, Code: "rate_limit", Message: "Slow down"}
	err := &RateLimitError{APIError: apiErr, RetryAfter: 3 * time.Second}
	if got := err.Error(); got != "pipedrive: http 429 rate_limit: Slow down (retry after 3s)" {
		t.Fatalf("unexpected error string: %q", got)
	}
	if !errors.Is(err, apiErr) {
		t.Fatalf("expected errors.Is to match wrapped APIError")
	}
}

func TestParseRetryAfter(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	future := now.Add(5 * time.Second).Format(http.TimeFormat)
	past := now.Add(-5 * time.Second).Format(http.TimeFormat)

	tests := []struct {
		name  string
		value string
		want  time.Duration
	}{
		{name: "empty", value: "", want: 0},
		{name: "seconds", value: "7", want: 7 * time.Second},
		{name: "negative seconds", value: "-1", want: 0},
		{name: "future http date", value: future, want: 5 * time.Second},
		{name: "past http date", value: past, want: 0},
		{name: "invalid", value: "later", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := parseRetryAfter(tt.value, now); got != tt.want {
				t.Fatalf("parseRetryAfter(%q) = %s, want %s", tt.value, got, tt.want)
			}
		})
	}
}

func TestParseResetHeader(t *testing.T) {
	t.Parallel()

	httpDate := time.Date(2025, 1, 1, 1, 2, 3, 0, time.UTC)

	tests := []struct {
		name  string
		value string
		want  time.Time
	}{
		{name: "empty", value: "", want: time.Time{}},
		{name: "unix seconds", value: "1735689605", want: time.Date(2025, 1, 1, 0, 0, 5, 0, time.UTC)},
		{name: "unix millis", value: "1735689605123", want: time.Date(2025, 1, 1, 0, 0, 5, 123000000, time.UTC)},
		{name: "http date", value: httpDate.Format(http.TimeFormat), want: httpDate},
		{name: "invalid", value: "not-a-time", want: time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := parseResetHeader(tt.value); !got.Equal(tt.want) {
				t.Fatalf("parseResetHeader(%q) = %s, want %s", tt.value, got, tt.want)
			}
		})
	}
}
