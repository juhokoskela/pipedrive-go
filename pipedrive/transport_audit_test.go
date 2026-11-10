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

type failingAuditTokenSource struct{ err error }

func (s failingAuditTokenSource) Token() (*oauth2.Token, error) { return nil, s.err }

type auditRequestBody struct {
	io.Reader
	closed bool
}

func (b *auditRequestBody) Close() error { b.closed = true; return nil }

func TestNewHTTPClient_ClosesBodyWhenAuthenticationFails(t *testing.T) {
	t.Parallel()
	wantErr := errors.New("token refresh failed")
	client := NewHTTPClient(Config{Auth: OAuth2Auth{TokenSource: failingAuditTokenSource{wantErr}}})
	body := &auditRequestBody{Reader: strings.NewReader("upload")}
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

type auditIdleTransport struct{ closed int }

func (t *auditIdleTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("unexpected request")
}
func (t *auditIdleTransport) CloseIdleConnections() { t.closed++ }

func TestNewHTTPClient_ForwardsCloseIdleConnections(t *testing.T) {
	t.Parallel()
	base := &auditIdleTransport{}
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
