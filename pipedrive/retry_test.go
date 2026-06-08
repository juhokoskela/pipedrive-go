package pipedrive

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRetryTransport_RetriesOnRetryableStatus(t *testing.T) {
	t.Parallel()

	var calls int
	next := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		status := 502
		if calls == 2 {
			status = 200
		}
		return &http.Response{
			StatusCode: status,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	var sleeps []time.Duration
	sleep := func(_ context.Context, d time.Duration) error {
		sleeps = append(sleeps, d)
		return nil
	}

	policy := DefaultRetryPolicy()
	policy.BaseDelay = 1 * time.Millisecond
	policy.MaxDelay = 1 * time.Millisecond
	policy.Jitter = func(d time.Duration) time.Duration { return d }

	rt := newRetryTransport(next, policy, retryTransportOptions{
		sleep: sleep,
		now:   func() time.Time { return time.Unix(0, 0).UTC() },
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()

	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(sleeps) != 1 {
		t.Fatalf("expected 1 sleep, got %d", len(sleeps))
	}
	if sleeps[0] != 1*time.Millisecond {
		t.Fatalf("expected sleep 1ms, got %s", sleeps[0])
	}
}

func TestRetryTransport_RespectsRetryAfter(t *testing.T) {
	t.Parallel()

	var calls int
	next := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		h := make(http.Header)
		if calls == 1 {
			h.Set("Retry-After", "2")
			return &http.Response{
				StatusCode: 429,
				Header:     h,
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			}, nil
		}
		return &http.Response{
			StatusCode: 200,
			Header:     h,
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	var sleeps []time.Duration
	sleep := func(_ context.Context, d time.Duration) error {
		sleeps = append(sleeps, d)
		return nil
	}

	policy := DefaultRetryPolicy()
	policy.Jitter = func(d time.Duration) time.Duration { return d }

	rt := newRetryTransport(next, policy, retryTransportOptions{
		sleep: sleep,
		now:   func() time.Time { return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) },
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.test", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()

	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
	if len(sleeps) != 1 || sleeps[0] != 2*time.Second {
		t.Fatalf("expected sleep 2s, got %v", sleeps)
	}
}

func TestRetryTransport_CanBeDisabledPerRequest(t *testing.T) {
	t.Parallel()

	var calls int
	next := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: 503,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	rt := newRetryTransport(next, DefaultRetryPolicy(), retryTransportOptions{
		sleep: func(context.Context, time.Duration) error { return nil },
		now:   time.Now,
	})

	req, _ := http.NewRequestWithContext(withNoRetry(context.Background()), http.MethodGet, "https://example.test", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()

	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
	if resp.StatusCode != 503 {
		t.Fatalf("expected status 503, got %d", resp.StatusCode)
	}
}

func TestRetryTransport_DoesNotRetryWhenBodyNotReplayable(t *testing.T) {
	t.Parallel()

	var calls int
	next := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: 502,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	rt := newRetryTransport(next, DefaultRetryPolicy(), retryTransportOptions{
		sleep: func(context.Context, time.Duration) error { return nil },
		now:   time.Now,
	})

	body := io.NopCloser(strings.NewReader("payload"))
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.test", body)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()

	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetryTransport_UsesPerRequestRetryPolicyOverride(t *testing.T) {
	t.Parallel()

	var calls int
	next := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		status := 502
		if calls == 3 {
			status = 200
		}
		return &http.Response{
			StatusCode: status,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	sleep := func(_ context.Context, _ time.Duration) error { return nil }
	policy := RetryPolicy{MaxAttempts: 2, BaseDelay: 0, Jitter: func(d time.Duration) time.Duration { return d }}
	rt := newRetryTransport(next, policy, retryTransportOptions{
		sleep: sleep,
		now:   time.Now,
	})

	override := RetryPolicy{MaxAttempts: 3, BaseDelay: 0, Jitter: func(d time.Duration) time.Duration { return d }}
	ctx := withRetryPolicy(context.Background(), override)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.test", nil)

	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()

	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestRetryTransport_StopsAfterExactMaxAttempts(t *testing.T) {
	t.Parallel()

	var calls int
	next := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: 502,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	rt := newRetryTransport(next, RetryPolicy{MaxAttempts: 2}, retryTransportOptions{
		sleep: func(context.Context, time.Duration) error { return nil },
		now:   time.Now,
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()

	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
	if resp.StatusCode != 502 {
		t.Fatalf("expected status 502, got %d", resp.StatusCode)
	}
}

func TestRetryTransport_NilRequestReturnsError(t *testing.T) {
	t.Parallel()

	rt := newRetryTransport(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		t.Fatalf("unexpected call to next transport")
		return nil, nil
	}), DefaultRetryPolicy(), retryTransportOptions{})

	_, err := rt.RoundTrip(nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "pipedrive: nil request" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRetryHelpers(t *testing.T) {
	t.Parallel()

	for _, method := range []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodTrace,
	} {
		if !isIdempotentMethod(method) {
			t.Fatalf("expected %s to be idempotent", method)
		}
	}
	for _, method := range []string{http.MethodPost, http.MethodPatch, "CUSTOM"} {
		if isIdempotentMethod(method) {
			t.Fatalf("expected %s to be non-idempotent", method)
		}
	}

	if got := fullJitter(0); got != 0 {
		t.Fatalf("fullJitter(0) = %s, want 0", got)
	}
	if got := fullJitter(-time.Second); got != 0 {
		t.Fatalf("fullJitter(-1s) = %s, want 0", got)
	}
	for i := 0; i < 20; i++ {
		if got := fullJitter(10 * time.Millisecond); got < 0 || got >= 10*time.Millisecond {
			t.Fatalf("fullJitter out of range: %s", got)
		}
	}
}

func TestSleepWithContext_ReturnsOnCancel(t *testing.T) {
	t.Parallel()

	if err := sleepWithContext(context.Background(), 0); err != nil {
		t.Fatalf("zero sleep error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := sleepWithContext(ctx, time.Hour)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

func TestSanitizeRetryPolicy_ClampsInvalidValues(t *testing.T) {
	t.Parallel()

	policy := sanitizeRetryPolicy(RetryPolicy{
		MaxAttempts: -3,
		BaseDelay:   -time.Second,
		MaxDelay:    0,
		Jitter:      nil,
	})

	if policy.MaxAttempts != 1 {
		t.Fatalf("expected max attempts 1, got %d", policy.MaxAttempts)
	}
	if policy.BaseDelay != 0 {
		t.Fatalf("expected base delay 0, got %s", policy.BaseDelay)
	}
	if policy.MaxDelay != 0 {
		t.Fatalf("expected max delay 0, got %s", policy.MaxDelay)
	}
	if policy.Jitter == nil {
		t.Fatalf("expected default jitter")
	}
	if got := policy.Jitter(3 * time.Second); got != 3*time.Second {
		t.Fatalf("unexpected default jitter result: %s", got)
	}

	policy = sanitizeRetryPolicy(RetryPolicy{
		MaxAttempts: 2,
		BaseDelay:   time.Second,
		MaxDelay:    -time.Second,
		Jitter:      func(d time.Duration) time.Duration { return d / 2 },
	})
	if policy.MaxDelay != time.Second {
		t.Fatalf("expected max delay to default to base delay, got %s", policy.MaxDelay)
	}
	if got := policy.Jitter(4 * time.Second); got != 2*time.Second {
		t.Fatalf("unexpected custom jitter result: %s", got)
	}
}
