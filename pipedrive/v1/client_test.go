package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
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
