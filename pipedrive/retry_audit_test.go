package pipedrive

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRetryTransport_ServiceUnavailableRetryAfter(t *testing.T) {
	t.Parallel()
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.Header().Set("Retry-After", "7")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	t.Cleanup(srv.Close)
	var delays []time.Duration
	client := srv.Client()
	client.Transport = newRetryTransport(client.Transport, DefaultRetryPolicy(), retryTransportOptions{
		sleep: func(_ context.Context, delay time.Duration) error { delays = append(delays, delay); return nil },
	})
	resp, err := client.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if len(delays) != 1 || delays[0] != 7*time.Second {
		t.Fatalf("delays = %v, want [7s]", delays)
	}
}

func TestRetryTransport_BackoffSaturatesBeforeOverflow(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name      string
		base, max time.Duration
		attempt   int
	}{
		{"many attempts", time.Second, 5 * time.Second, 64},
		{"large base", time.Duration(math.MaxInt64/2 + 1), time.Duration(math.MaxInt64), 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			policy := sanitizeRetryPolicy(RetryPolicy{BaseDelay: tc.base, MaxDelay: tc.max})
			rt := &retryTransport{}
			if got := rt.nextDelay(tc.attempt, nil, policy); got != tc.max {
				t.Fatalf("delay = %v, want %v", got, tc.max)
			}
		})
	}
}
