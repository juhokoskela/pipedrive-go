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

func TestNewHTTPClient_AppliesUserAgentAndAPIToken(t *testing.T) {
	t.Parallel()

	base := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("User-Agent"); got != "pipedrive-go/test" {
			t.Fatalf("unexpected user-agent: %q", got)
		}
		if got := req.Header.Get("x-api-token"); got != "token123" {
			t.Fatalf("unexpected x-api-token: %q", got)
		}
		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	httpClient := NewHTTPClient(Config{
		HTTPClient: &http.Client{Transport: base},
		UserAgent:  "pipedrive-go/test",
		Auth:       APITokenAuth("token123"),
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	resp, err := httpClient.Transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
}

func TestAPITokenAuth_Apply(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	if err := APITokenAuth("token123").Apply(req); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if got := req.Header.Get("x-api-token"); got != "token123" {
		t.Fatalf("unexpected token header: %q", got)
	}

	if err := APITokenAuth("other").Apply(req); err != nil {
		t.Fatalf("Apply existing header error: %v", err)
	}
	if got := req.Header.Get("x-api-token"); got != "token123" {
		t.Fatalf("expected existing token to be preserved, got %q", got)
	}

	if err := APITokenAuth("").Apply(req); err != nil {
		t.Fatalf("empty auth should not error: %v", err)
	}
	if err := APITokenAuth("token").Apply(nil); err != nil {
		t.Fatalf("nil request should not error: %v", err)
	}
}

func TestOAuth2Auth_Apply(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	auth := OAuth2Auth{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "access123"})}
	if err := auth.Apply(req); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer access123" {
		t.Fatalf("unexpected authorization header: %q", got)
	}

	req.Header.Set("Authorization", "Bearer existing")
	if err := auth.Apply(req); err != nil {
		t.Fatalf("Apply existing header error: %v", err)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer existing" {
		t.Fatalf("expected existing authorization to be preserved, got %q", got)
	}

	if err := (OAuth2Auth{}).Apply(req); err != nil {
		t.Fatalf("empty auth should not error: %v", err)
	}
	if err := auth.Apply(nil); err != nil {
		t.Fatalf("nil request should not error: %v", err)
	}
}

func TestMultiAuth_ApplySkipsNilAndStopsOnError(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	wantErr := errors.New("auth failed")
	auth := MultiAuth{
		nil,
		APITokenAuth("token123"),
		errAuth{err: wantErr},
		APITokenAuth("after-error"),
	}
	err = auth.Apply(req)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected auth error, got %v", err)
	}
	if got := req.Header.Get("x-api-token"); got != "token123" {
		t.Fatalf("unexpected token header: %q", got)
	}
}

type errAuth struct {
	err error
}

func (a errAuth) Apply(*http.Request) error {
	return a.err
}
