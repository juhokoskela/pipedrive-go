package pipedrive

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
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

