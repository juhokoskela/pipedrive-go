package pipedrive

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/oauth2"
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

type failingTokenSource struct{ err error }

func (s failingTokenSource) Token() (*oauth2.Token, error) { return nil, s.err }

type trackingRequestBody struct {
	io.Reader
	closed bool
}

func (b *trackingRequestBody) Close() error { b.closed = true; return nil }

func TestNewHTTPClient_ClosesBodyWhenAuthenticationFails(t *testing.T) {
	t.Parallel()
	wantErr := errors.New("token refresh failed")
	client := NewHTTPClient(Config{Auth: OAuth2Auth{TokenSource: failingTokenSource{wantErr}}})
	body := &trackingRequestBody{Reader: strings.NewReader("upload")}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.test", body)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Do(req)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want token error", err)
	}
	if !body.closed {
		t.Fatal("request body was not closed")
	}
}

type trackingIdleTransport struct{ closed int }

func (t *trackingIdleTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("unexpected request")
}
func (t *trackingIdleTransport) CloseIdleConnections() { t.closed++ }

func TestNewHTTPClient_ForwardsCloseIdleConnections(t *testing.T) {
	t.Parallel()
	base := &trackingIdleTransport{}
	client := NewHTTPClient(Config{
		HTTPClient: &http.Client{Transport: base},
		Auth:       APITokenAuth("test"), UserAgent: "test-agent",
		Middleware: []Middleware{func(next http.RoundTripper) http.RoundTripper { return roundTripperFunc(next.RoundTrip) }},
	})
	client.CloseIdleConnections()
	if base.closed != 1 {
		t.Fatalf("underlying CloseIdleConnections calls = %d, want 1", base.closed)
	}
}
