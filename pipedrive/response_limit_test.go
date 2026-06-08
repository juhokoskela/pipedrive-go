package pipedrive

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestResponseLimitTransport_DefaultLimitApplied(t *testing.T) {
	t.Parallel()

	rt := newResponseLimitTransport(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       &fixedSizeReadCloser{remaining: defaultMaxResponseSize + 1},
			Request:    req,
		}, nil
	}), 0)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected round trip error: %v", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	var tooLarge *ResponseTooLargeError
	if !errors.As(err, &tooLarge) {
		t.Fatalf("expected ResponseTooLargeError, got %T (%v)", err, err)
	}
	if tooLarge.Limit != defaultMaxResponseSize {
		t.Fatalf("expected limit=%d, got %d", defaultMaxResponseSize, tooLarge.Limit)
	}
}

func TestResponseLimitTransport_SkipsLimitForNon2xxResponses(t *testing.T) {
	t.Parallel()

	rt := newResponseLimitTransport(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
			Body:       &fixedSizeReadCloser{remaining: defaultMaxResponseSize + 1},
			Request:    req,
		}, nil
	}), 1)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected round trip error: %v", err)
	}
	defer resp.Body.Close()

	if _, ok := resp.Body.(*responseLimitReadCloser); ok {
		t.Fatalf("expected original body for non-2xx responses")
	}
}

func TestResponseLimitReadCloser_ReportsOverflowOnProbeRead(t *testing.T) {
	t.Parallel()

	rc := &responseLimitReadCloser{
		body:      io.NopCloser(strings.NewReader("abc")),
		limit:     2,
		remaining: 2,
	}

	buf := make([]byte, 4)
	n, err := rc.Read(buf)
	if err != nil {
		t.Fatalf("unexpected first read error: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected first read to return 2 bytes, got %d", n)
	}
	if got := string(buf[:n]); got != "ab" {
		t.Fatalf("unexpected first read payload: %q", got)
	}

	n, err = rc.Read(buf)
	var tooLarge *ResponseTooLargeError
	if !errors.As(err, &tooLarge) {
		t.Fatalf("expected ResponseTooLargeError, got %T (%v)", err, err)
	}
	if n != 0 {
		t.Fatalf("expected probe read to return 0 bytes, got %d", n)
	}
	if tooLarge.Limit != 2 {
		t.Fatalf("expected limit=2, got %d", tooLarge.Limit)
	}
}

func TestResponseLimitReadCloser_AllowsBodyExactlyAtLimit(t *testing.T) {
	t.Parallel()

	rc := &responseLimitReadCloser{
		body:      io.NopCloser(strings.NewReader("ab")),
		limit:     2,
		remaining: 2,
	}

	body, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if got := string(body); got != "ab" {
		t.Fatalf("unexpected body: %q", got)
	}
}

func TestRawClient_Do_ResponseLimitCanBeOverriddenPerRequest(t *testing.T) {
	t.Parallel()

	httpClient := NewHTTPClient(Config{
		HTTPClient: &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
				Request:    req,
			}, nil
		})},
		MaxResponseSize: 4,
		RetryPolicy:     &RetryPolicy{MaxAttempts: 1},
	})

	raw, err := NewRawClient("https://example.test", httpClient)
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}

	var out struct {
		Ok bool `json:"ok"`
	}
	if err := raw.Do(context.Background(), http.MethodGet, "/x", nil, nil, &out, WithResponseSizeLimit(32)); err != nil {
		t.Fatalf("Do error: %v", err)
	}
	if !out.Ok {
		t.Fatalf("expected ok=true")
	}
}

func TestRawClient_Do_ResponseLimitCanBeDisabledPerRequest(t *testing.T) {
	t.Parallel()

	httpClient := NewHTTPClient(Config{
		HTTPClient: &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
				Request:    req,
			}, nil
		})},
		MaxResponseSize: 4,
		RetryPolicy:     &RetryPolicy{MaxAttempts: 1},
	})

	raw, err := NewRawClient("https://example.test", httpClient)
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}

	var out struct {
		Ok bool `json:"ok"`
	}
	if err := raw.Do(context.Background(), http.MethodGet, "/x", nil, nil, &out, WithNoResponseSizeLimit()); err != nil {
		t.Fatalf("Do error: %v", err)
	}
	if !out.Ok {
		t.Fatalf("expected ok=true")
	}
}

type fixedSizeReadCloser struct {
	remaining int64
}

func (r *fixedSizeReadCloser) Read(p []byte) (int, error) {
	if r.remaining == 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > r.remaining {
		p = p[:int(r.remaining)]
	}
	for i := range p {
		p[i] = 'a'
	}
	r.remaining -= int64(len(p))
	return len(p), nil
}

func (r *fixedSizeReadCloser) Close() error {
	return nil
}
