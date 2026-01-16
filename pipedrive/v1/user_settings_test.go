package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestUserSettingsService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/userSettings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"beta_app":true,"list_limit":25,"callto_link_syntax":"callto"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	settings, err := client.UserSettings.Get(
		context.Background(),
		WithUserSettingsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if !settings.BetaApp {
		t.Fatalf("expected BetaApp to be true")
	}
	if settings.ListLimit != 25 {
		t.Fatalf("unexpected ListLimit: %d", settings.ListLimit)
	}
	if settings.CalltoLinkSyntax != "callto" {
		t.Fatalf("unexpected CalltoLinkSyntax: %q", settings.CalltoLinkSyntax)
	}
}
