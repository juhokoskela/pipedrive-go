package v2

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestNewClient_ConfiguresRawClient(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ping" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	var out struct {
		Ok bool `json:"ok"`
	}
	if err := client.Raw.Do(context.Background(), http.MethodGet, "/ping", nil, nil, &out); err != nil {
		t.Fatalf("raw do error: %v", err)
	}
	if !out.Ok {
		t.Fatalf("expected ok=true")
	}
}

type stubRoundTripper func(*http.Request) (*http.Response, error)

func (f stubRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// The default base URL must still pin credentials: NewClient backfills the
// resolved base URL into the config so an empty BaseURL does not leave the
// auth middleware unpinned.
func TestNewClient_DefaultBaseURLPinsCredentials(t *testing.T) {
	t.Parallel()

	var attackerHeaders http.Header
	stub := stubRoundTripper(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "attacker.example" {
			attackerHeaders = r.Header.Clone()
			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{}`)),
				Request:    r,
			}, nil
		}
		h := make(http.Header)
		h.Set("Location", "https://attacker.example/steal")
		return &http.Response{
			StatusCode: http.StatusFound,
			Header:     h,
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    r,
		}, nil
	})

	client, err := NewClient(pipedrive.Config{
		HTTPClient: &http.Client{Transport: stub},
		Auth:       pipedrive.APITokenAuth("secret-token"),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_ = client.Raw.Do(context.Background(), http.MethodGet, "/deals", nil, nil, nil)

	if attackerHeaders == nil {
		t.Fatal("expected the redirect to be followed to the attacker host")
	}
	if got := attackerHeaders.Get("x-api-token"); got != "" {
		t.Fatalf("credentials leaked with default base URL: %q", got)
	}
}
