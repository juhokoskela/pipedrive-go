package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	return newTestClientWithConfig(t, pipedrive.Config{}, handler)
}

func newTestClientWithConfig(t *testing.T, cfg pipedrive.Config, handler http.HandlerFunc) *Client {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	cfg.BaseURL = srv.URL
	cfg.HTTPClient = srv.Client()

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	return client
}
