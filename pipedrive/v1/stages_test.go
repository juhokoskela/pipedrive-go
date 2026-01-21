package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestStagesService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/stages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ids"); got != "1,2" {
			t.Fatalf("unexpected ids: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":[1,2]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Stages.Delete(
		context.Background(),
		[]StageID{1, 2},
		WithStagesRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if result == nil || len(result.IDs) != 2 || result.IDs[0] != 1 || result.IDs[1] != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestStagesService_ListDeals(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/stages/3/deals" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"title":"Deal"}],"additional_data":{"start":0,"limit":1,"more_items_in_collection":false}}`))
	})

	query := url.Values{}
	query.Set("limit", "1")
	deals, page, err := client.Stages.ListDeals(context.Background(), StageID(3), WithStageDealsQuery(query))
	if err != nil {
		t.Fatalf("ListDeals error: %v", err)
	}
	if len(deals) != 1 || deals[0].ID != 1 {
		t.Fatalf("unexpected deals: %#v", deals)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}
