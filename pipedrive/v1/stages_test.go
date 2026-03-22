package v1

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

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
