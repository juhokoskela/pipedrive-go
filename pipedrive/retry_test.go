package pipedrive

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

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
