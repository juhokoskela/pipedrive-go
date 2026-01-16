package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestBillingService_ListAddons(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/billing/subscriptions/addons" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"code":"leadbooster_v2"},{"code":"prospector"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	addons, err := client.Billing.ListAddons(context.Background(), WithBillingRequestOptions(pipedrive.WithHeader("X-Test", "1")))
	if err != nil {
		t.Fatalf("ListAddons error: %v", err)
	}
	if len(addons) != 2 || addons[0].Code != "leadbooster_v2" {
		t.Fatalf("unexpected addons: %#v", addons)
	}
}
