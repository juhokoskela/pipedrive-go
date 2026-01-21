package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestCurrenciesService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/currencies" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("term"); got != "us" {
			t.Fatalf("unexpected term: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"code":"USD","name":"US Dollar"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	currs, err := client.Currencies.List(context.Background(), ListCurrenciesRequest{Term: "us"}, pipedrive.WithHeader("X-Test", "1"))
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(currs) != 1 || currs[0].Code != "USD" {
		t.Fatalf("unexpected currencies: %#v", currs)
	}
}
