package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestRecentsService_List(t *testing.T) {
	t.Parallel()

	since := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/recents" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("since_timestamp"); got != "2024-01-01 00:00:00" {
			t.Fatalf("unexpected since_timestamp: %q", got)
		}
		if got := q.Get("items"); got != "activity,deal" {
			t.Fatalf("unexpected items: %q", got)
		}
		if got := q.Get("start"); got != "1" {
			t.Fatalf("unexpected start: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"item":"activity","id":12,"data":{"id":12,"subject":"Call"}}],"additional_data":{"since_timestamp":"2024-01-01 00:00:00","last_timestamp_on_page":"2024-01-01 00:00:01","pagination":{"start":1,"limit":2,"more_items_in_collection":false}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	items, additional, err := client.Recents.List(
		context.Background(),
		WithRecentsSince(since),
		WithRecentsItems(RecentsItemActivity, RecentsItemDeal),
		WithRecentsStart(1),
		WithRecentsLimit(2),
		WithRecentsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(items) != 1 || items[0].Item != RecentsItemActivity || items[0].ID != 12 {
		t.Fatalf("unexpected items: %#v", items)
	}
	if !json.Valid(items[0].Data) {
		t.Fatalf("expected valid data JSON")
	}
	var payload struct {
		Subject string `json:"subject"`
	}
	if err := json.Unmarshal(items[0].Data, &payload); err != nil {
		t.Fatalf("decode data: %v", err)
	}
	if payload.Subject != "Call" {
		t.Fatalf("unexpected data: %#v", payload)
	}
	if additional == nil || additional.Pagination == nil || additional.Pagination.MoreItemsInCollection {
		t.Fatalf("unexpected additional data: %#v", additional)
	}
}

func TestRecentsService_ListRequiresSince(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.Recents.List(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
}
