package pipedrive

import (
	"context"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration

	Jitter func(time.Duration) time.Duration

	RetryAllMethods bool
}

func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts:     4,
		BaseDelay:       200 * time.Millisecond,
		MaxDelay:        5 * time.Second,
		Jitter:          fullJitter,
		RetryAllMethods: false,
	}
}

type retryTransport struct {
	next   http.RoundTripper
	policy RetryPolicy
	opts   retryTransportOptions
}

type retryTransportOptions struct {
	sleep func(context.Context, time.Duration) error
	now   func() time.Time
}

func newRetryTransport(next http.RoundTripper, policy RetryPolicy, opts retryTransportOptions) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	policy = sanitizeRetryPolicy(policy)
	if opts.sleep == nil {
		opts.sleep = sleepWithContext
	}
	if opts.now == nil {
		opts.now = time.Now
	}
	return &retryTransport{
		next:   next,
		policy: policy,
		opts:   opts,
	}
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, nil
	}

	policy := t.policy
	if override, ok := retryPolicyFromContext(req.Context()); ok {
		policy = sanitizeRetryPolicy(override)
	}

	if isNoRetry(req.Context()) || policy.MaxAttempts <= 1 {
		return t.next.RoundTrip(req)
	}

	origGetBody := req.GetBody
	canReplayBody := req.Body == nil || req.Body == http.NoBody || origGetBody != nil

	for attempt := 1; attempt <= t.policy.MaxAttempts; attempt++ {
		attemptReq := req
		if attempt > 1 {
			attemptReq = req.Clone(req.Context())
			if origGetBody != nil {
				body, err := origGetBody()
				if err != nil {
					return nil, err
				}
				attemptReq.Body = body
			}
		}

		resp, err := t.next.RoundTrip(attemptReq)
		if !t.shouldRetry(attempt, attemptReq, canReplayBody, resp, err, policy) {
			return resp, err
		}

		delay := t.nextDelay(attempt, resp, policy)
		if resp != nil && resp.Body != nil {
			_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
			_ = resp.Body.Close()
		}

		if delay > 0 {
			if err := t.opts.sleep(req.Context(), delay); err != nil {
				return nil, err
			}
		}
	}

	return t.next.RoundTrip(req)
}

func (t *retryTransport) shouldRetry(attempt int, req *http.Request, canReplayBody bool, resp *http.Response, err error, policy RetryPolicy) bool {
	if err != nil || resp == nil {
		return false
	}
	if attempt >= policy.MaxAttempts {
		return false
	}

	switch resp.StatusCode {
	case 429:
		return canReplayBody
	case 502, 503, 504:
		if !canReplayBody {
			return false
		}
		return policy.RetryAllMethods || isIdempotentMethod(req.Method)
	default:
		return false
	}
}

func (t *retryTransport) nextDelay(attempt int, resp *http.Response, policy RetryPolicy) time.Duration {
	if resp != nil && resp.StatusCode == 429 {
		if ra := parseRetryAfter(resp.Header.Get("Retry-After"), t.opts.now()); ra > 0 {
			return ra
		}
	}

	if policy.BaseDelay <= 0 {
		return 0
	}

	exp := float64(policy.BaseDelay) * math.Pow(2, float64(attempt-1))
	delay := time.Duration(exp)
	if delay > policy.MaxDelay {
		delay = policy.MaxDelay
	}
	return policy.Jitter(delay)
}

func isIdempotentMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

func fullJitter(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	return time.Duration(rand.Float64() * float64(d))
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type noRetryKey struct{}

func withNoRetry(ctx context.Context) context.Context {
	return context.WithValue(ctx, noRetryKey{}, true)
}

func isNoRetry(ctx context.Context) bool {
	disabled, _ := ctx.Value(noRetryKey{}).(bool)
	return disabled
}

type retryPolicyKey struct{}

func withRetryPolicy(ctx context.Context, policy RetryPolicy) context.Context {
	return context.WithValue(ctx, retryPolicyKey{}, policy)
}

func retryPolicyFromContext(ctx context.Context) (RetryPolicy, bool) {
	policy, ok := ctx.Value(retryPolicyKey{}).(RetryPolicy)
	return policy, ok
}

func sanitizeRetryPolicy(policy RetryPolicy) RetryPolicy {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 1
	}
	if policy.BaseDelay < 0 {
		policy.BaseDelay = 0
	}
	if policy.MaxDelay <= 0 {
		policy.MaxDelay = policy.BaseDelay
	}
	if policy.Jitter == nil {
		policy.Jitter = func(d time.Duration) time.Duration { return d }
	}
	return policy
}
