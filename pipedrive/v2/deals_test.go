package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestDealsService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":1,"title":"Test deal"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	deal, err := client.Deals.Get(context.Background(), 1, WithDealRequestOptions(pipedrive.WithHeader("X-Test", "1")))
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if deal.ID != 1 || deal.Title != "Test deal" {
		t.Fatalf("unexpected deal: %#v", deal)
	}
}

func TestDealsService_ListPager(t *testing.T) {
	t.Parallel()

	var listCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		switch r.URL.Path {
		case "/deals":
			listCalls++
			cursor := r.URL.Query().Get("cursor")
			w.Header().Set("Content-Type", "application/json")
			if listCalls == 1 {
				if cursor != "" {
					t.Fatalf("expected no cursor on first page, got %q", cursor)
				}
				if got := r.URL.Query().Get("limit"); got != "2" {
					t.Fatalf("expected limit=2, got %q", got)
				}
				_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":"c2"}}`))
				return
			}
			if listCalls == 2 {
				if cursor != "c2" {
					t.Fatalf("expected cursor c2 on second page, got %q", cursor)
				}
				if got := r.URL.Query().Get("limit"); got != "2" {
					t.Fatalf("expected limit=2, got %q", got)
				}
				_, _ = w.Write([]byte(`{"data":[{"id":3}],"additional_data":{"next_cursor":null}}`))
				return
			}
			t.Fatalf("unexpected listCalls=%d", listCalls)
		default:
			if strings.HasPrefix(r.URL.Path, "/deals/") {
				t.Fatalf("unexpected deal path: %s", r.URL.Path)
			}
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	pager := client.Deals.ListPager(WithDealsPageSize(2))

	var ids []DealID
	for pager.Next(context.Background()) {
		for _, d := range pager.Items() {
			ids = append(ids, d.ID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if want := []DealID{1, 2, 3}; len(ids) != len(want) || ids[0] != want[0] || ids[1] != want[1] || ids[2] != want[2] {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestDealsService_ForEach(t *testing.T) {
	t.Parallel()

	var listCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals" {
			http.NotFound(w, r)
			return
		}
		listCalls++
		cursor := r.URL.Query().Get("cursor")
		w.Header().Set("Content-Type", "application/json")
		if listCalls == 1 {
			if cursor != "" {
				t.Fatalf("expected no cursor on first page, got %q", cursor)
			}
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("expected limit=2, got %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":"c2"}}`))
			return
		}
		if listCalls == 2 {
			if cursor != "c2" {
				t.Fatalf("expected cursor c2 on second page, got %q", cursor)
			}
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("expected limit=2, got %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":3}],"additional_data":{"next_cursor":null}}`))
			return
		}
		t.Fatalf("unexpected listCalls=%d", listCalls)
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	var ids []DealID
	err = client.Deals.ForEach(context.Background(), func(d Deal) error {
		ids = append(ids, d.ID)
		return nil
	}, WithDealsPageSize(2))
	if err != nil {
		t.Fatalf("ForEach error: %v", err)
	}
	if want := []DealID{1, 2, 3}; len(ids) != len(want) || ids[0] != want[0] || ids[1] != want[1] || ids[2] != want[2] {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestDealsService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["title"] != "New Deal" {
			t.Fatalf("unexpected title: %#v", payload["title"])
		}
		if payload["value"] != float64(1200) {
			t.Fatalf("unexpected value: %#v", payload["value"])
		}
		if payload["currency"] != "USD" {
			t.Fatalf("unexpected currency: %#v", payload["currency"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":7,"title":"New Deal"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	deal, err := client.Deals.Create(
		context.Background(),
		WithDealTitle("New Deal"),
		WithDealValue(1200),
		WithDealCurrency("USD"),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if deal.ID != 7 || deal.Title != "New Deal" {
		t.Fatalf("unexpected deal: %#v", deal)
	}
}

func TestDealsService_Update(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["title"] != "Updated Deal" {
			t.Fatalf("unexpected title: %#v", payload["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":7,"title":"Updated Deal"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	deal, err := client.Deals.Update(
		context.Background(),
		DealID(7),
		WithDealTitle("Updated Deal"),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if deal.ID != 7 || deal.Title != "Updated Deal" {
		t.Fatalf("unexpected deal: %#v", deal)
	}
}

func TestDealsService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":7}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Deals.Delete(context.Background(), DealID(7))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if result.ID != 7 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestDealsService_Search(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("term"); got != "deal" {
			t.Fatalf("unexpected term: %q", got)
		}
		if got := q.Get("fields"); got != "title" {
			t.Fatalf("unexpected fields: %q", got)
		}
		if got := q.Get("exact_match"); got != "true" {
			t.Fatalf("unexpected exact_match: %q", got)
		}
		if got := q.Get("status"); got != "open" {
			t.Fatalf("unexpected status: %q", got)
		}
		if got := q.Get("include_fields"); got != "deal.cc_email" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"items":[{"item":{"id":7,"title":"Deal"}}]},"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	results, next, err := client.Deals.Search(
		context.Background(),
		"deal",
		WithDealSearchFields(DealSearchFieldTitle),
		WithDealSearchExactMatch(true),
		WithDealSearchStatus(DealSearchStatusOpen),
		WithDealSearchIncludeFields(DealSearchIncludeFieldDealCCEmail),
		WithDealSearchPageSize(1),
		WithDealSearchCursor("c1"),
	)
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(results.Items) != 1 {
		t.Fatalf("unexpected results: %#v", results)
	}
}

func TestDealsService_ListArchived(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/archived" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("owner_id"); got != "3" {
			t.Fatalf("unexpected owner_id: %q", got)
		}
		if got := q.Get("status"); got != "won" {
			t.Fatalf("unexpected status: %q", got)
		}
		if got := q.Get("include_fields"); got != "activities_count" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("custom_fields"); got != "cf_1" {
			t.Fatalf("unexpected custom_fields: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":2,"title":"Archived"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	deals, next, err := client.Deals.ListArchived(
		context.Background(),
		WithArchivedDealsOwnerID(UserID(3)),
		WithArchivedDealsStatus(DealStatusWon),
		WithArchivedDealsIncludeFields(DealIncludeFieldActivitiesCount),
		WithArchivedDealsCustomFields("cf_1"),
		WithArchivedDealsPageSize(2),
		WithArchivedDealsCursor("c1"),
	)
	if err != nil {
		t.Fatalf("ListArchived error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(deals) != 1 || deals[0].ID != 2 {
		t.Fatalf("unexpected deals: %#v", deals)
	}
}

func TestDealsService_ConvertToLead(t *testing.T) {
	t.Parallel()

	convID := "11111111-1111-1111-1111-111111111111"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/convert/lead" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"conversion_id":"` + convID + `"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	job, err := client.Deals.ConvertToLead(context.Background(), DealID(7))
	if err != nil {
		t.Fatalf("ConvertToLead error: %v", err)
	}
	if string(job.ConversionID) != convID {
		t.Fatalf("unexpected conversion id: %#v", job)
	}
}

func TestDealsService_ConversionStatus(t *testing.T) {
	t.Parallel()

	convID := "11111111-1111-1111-1111-111111111111"
	leadID := "22222222-2222-2222-2222-222222222222"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/convert/status/"+convID {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"conversion_id":"` + convID + `","status":"completed","deal_id":9,"lead_id":"` + leadID + `"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	status, err := client.Deals.ConversionStatus(context.Background(), DealID(7), ConversionID(convID))
	if err != nil {
		t.Fatalf("ConversionStatus error: %v", err)
	}
	if status.Status != ConversionStatusCompleted {
		t.Fatalf("unexpected status: %#v", status)
	}
	if status.DealID == nil || *status.DealID != 9 {
		t.Fatalf("unexpected deal id: %#v", status)
	}
	if status.LeadID == nil || string(*status.LeadID) != leadID {
		t.Fatalf("unexpected lead id: %#v", status)
	}
}

func TestDealsService_ListFollowers(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/followers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"user_id":5}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	followers, next, err := client.Deals.ListFollowers(
		context.Background(),
		DealID(7),
		WithDealFollowersPageSize(2),
		WithDealFollowersCursor("c1"),
	)
	if err != nil {
		t.Fatalf("ListFollowers error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(followers) != 1 || followers[0].UserID != 5 {
		t.Fatalf("unexpected followers: %#v", followers)
	}
}

func TestDealsService_AddFollower(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/followers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["user_id"] != float64(6) {
			t.Fatalf("unexpected user_id: %#v", payload["user_id"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user_id":6}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	follower, err := client.Deals.AddFollower(context.Background(), DealID(7), UserID(6))
	if err != nil {
		t.Fatalf("AddFollower error: %v", err)
	}
	if follower.UserID != 6 {
		t.Fatalf("unexpected follower: %#v", follower)
	}
}

func TestDealsService_DeleteFollower(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/followers/6" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user_id":6}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Deals.DeleteFollower(context.Background(), DealID(7), UserID(6))
	if err != nil {
		t.Fatalf("DeleteFollower error: %v", err)
	}
	if result.UserID != 6 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestDealsService_FollowersChangelog(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/followers/changelog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"action":"added","actor_user_id":1,"follower_user_id":2,"time":"2024-01-01T10:00:00Z"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	changelog, next, err := client.Deals.FollowersChangelog(
		context.Background(),
		DealID(7),
		WithDealFollowersChangelogPageSize(1),
		WithDealFollowersChangelogCursor("c1"),
	)
	if err != nil {
		t.Fatalf("FollowersChangelog error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(changelog) != 1 || changelog[0].FollowerUserID != 2 {
		t.Fatalf("unexpected changelog: %#v", changelog)
	}
}

func TestDealsService_ListProducts(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := q.Get("sort_by"); got != "add_time" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "desc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":10,"deal_id":7,"product_id":3,"name":"Widget"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	products, next, err := client.Deals.ListProducts(
		context.Background(),
		DealID(7),
		WithDealProductsPageSize(2),
		WithDealProductsCursor("c1"),
		WithDealProductsSortBy(DealProductSortByAddTime),
		WithDealProductsSortDirection(SortDesc),
	)
	if err != nil {
		t.Fatalf("ListProducts error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(products) != 1 || products[0].ID != 10 {
		t.Fatalf("unexpected products: %#v", products)
	}
}

func TestDealsService_ListProductsAcrossDeals(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["deal_ids"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
			t.Fatalf("unexpected deal_ids: %#v", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":12,"deal_id":1,"product_id":9,"name":"Bulk"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	products, next, err := client.Deals.ListProductsAcrossDeals(
		context.Background(),
		[]DealID{1, 2},
		WithDealsProductsPageSize(1),
		WithDealsProductsCursor("c1"),
	)
	if err != nil {
		t.Fatalf("ListProductsAcrossDeals error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(products) != 1 || products[0].ID != 12 {
		t.Fatalf("unexpected products: %#v", products)
	}
}

func TestDealsService_AddProduct(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["product_id"] != float64(9) {
			t.Fatalf("unexpected product_id: %#v", payload["product_id"])
		}
		if payload["item_price"] != float64(199) {
			t.Fatalf("unexpected item_price: %#v", payload["item_price"])
		}
		if payload["quantity"] != float64(2) {
			t.Fatalf("unexpected quantity: %#v", payload["quantity"])
		}
		if payload["discount_type"] != "percentage" {
			t.Fatalf("unexpected discount_type: %#v", payload["discount_type"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":15,"deal_id":7,"product_id":9}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	product, err := client.Deals.AddProduct(
		context.Background(),
		DealID(7),
		WithDealProductProductID(ProductID(9)),
		WithDealProductItemPrice(199),
		WithDealProductQuantity(2),
		WithDealProductDiscountType(DealProductDiscountTypePercentage),
	)
	if err != nil {
		t.Fatalf("AddProduct error: %v", err)
	}
	if product.ID != 15 {
		t.Fatalf("unexpected product: %#v", product)
	}
}

func TestDealsService_AddProducts(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/products/bulk" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		items, ok := payload["data"].([]interface{})
		if !ok || len(items) != 2 {
			t.Fatalf("unexpected data: %#v", payload["data"])
		}
		first, ok := items[0].(map[string]interface{})
		if !ok || first["product_id"] != float64(9) {
			t.Fatalf("unexpected first product: %#v", items[0])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":21},{"id":22}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	products, err := client.Deals.AddProducts(
		context.Background(),
		DealID(7),
		[]DealProductInput{
			NewDealProductInput(
				WithDealProductProductID(ProductID(9)),
				WithDealProductItemPrice(100),
				WithDealProductQuantity(1),
			),
			NewDealProductInput(
				WithDealProductProductID(ProductID(10)),
				WithDealProductItemPrice(200),
				WithDealProductQuantity(2),
			),
		},
	)
	if err != nil {
		t.Fatalf("AddProducts error: %v", err)
	}
	if len(products) != 2 || products[0].ID != 21 || products[1].ID != 22 {
		t.Fatalf("unexpected products: %#v", products)
	}
}

func TestDealsService_UpdateProduct(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/products/15" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["item_price"] != float64(250) {
			t.Fatalf("unexpected item_price: %#v", payload["item_price"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":15,"deal_id":7,"product_id":9}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	product, err := client.Deals.UpdateProduct(
		context.Background(),
		DealID(7),
		DealProductAttachmentID(15),
		WithDealProductItemPrice(250),
	)
	if err != nil {
		t.Fatalf("UpdateProduct error: %v", err)
	}
	if product.ID != 15 {
		t.Fatalf("unexpected product: %#v", product)
	}
}

func TestDealsService_DeleteProduct(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/products/15" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":15}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Deals.DeleteProduct(context.Background(), DealID(7), DealProductAttachmentID(15))
	if err != nil {
		t.Fatalf("DeleteProduct error: %v", err)
	}
	if result.ID != 15 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestDealsService_DeleteProducts(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ids"); got != "15,16" {
			t.Fatalf("unexpected ids: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"ids":[15,16]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Deals.DeleteProducts(
		context.Background(),
		DealID(7),
		WithDealProductAttachmentIDs(DealProductAttachmentID(15), DealProductAttachmentID(16)),
	)
	if err != nil {
		t.Fatalf("DeleteProducts error: %v", err)
	}
	if len(result.IDs) != 2 || result.IDs[0] != 15 || result.IDs[1] != 16 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestDealsService_ListAdditionalDiscounts(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/discounts" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"11111111-1111-1111-1111-111111111111","deal_id":7,"amount":10,"type":"amount"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	discounts, err := client.Deals.ListAdditionalDiscounts(context.Background(), DealID(7))
	if err != nil {
		t.Fatalf("ListAdditionalDiscounts error: %v", err)
	}
	if len(discounts) != 1 || discounts[0].ID != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("unexpected discounts: %#v", discounts)
	}
}

func TestDealsService_AddAdditionalDiscount(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/discounts" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["amount"] != float64(10) {
			t.Fatalf("unexpected amount: %#v", payload["amount"])
		}
		if payload["type"] != "amount" {
			t.Fatalf("unexpected type: %#v", payload["type"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"11111111-1111-1111-1111-111111111111","deal_id":7,"amount":10,"type":"amount"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	discount, err := client.Deals.AddAdditionalDiscount(
		context.Background(),
		DealID(7),
		WithAdditionalDiscountAmount(10),
		WithAdditionalDiscountType(AdditionalDiscountTypeAmount),
		WithAdditionalDiscountDescription("Promo"),
	)
	if err != nil {
		t.Fatalf("AddAdditionalDiscount error: %v", err)
	}
	if string(discount.ID) != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("unexpected discount: %#v", discount)
	}
}

func TestDealsService_UpdateAdditionalDiscount(t *testing.T) {
	t.Parallel()

	discountID := "11111111-1111-1111-1111-111111111111"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/discounts/"+discountID {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["amount"] != float64(15) {
			t.Fatalf("unexpected amount: %#v", payload["amount"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"` + discountID + `","deal_id":7,"amount":15,"type":"amount"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	discount, err := client.Deals.UpdateAdditionalDiscount(
		context.Background(),
		DealID(7),
		AdditionalDiscountID(discountID),
		WithAdditionalDiscountAmount(15),
	)
	if err != nil {
		t.Fatalf("UpdateAdditionalDiscount error: %v", err)
	}
	if string(discount.ID) != discountID {
		t.Fatalf("unexpected discount: %#v", discount)
	}
}

func TestDealsService_DeleteAdditionalDiscount(t *testing.T) {
	t.Parallel()

	discountID := "11111111-1111-1111-1111-111111111111"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/discounts/"+discountID {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"` + discountID + `"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Deals.DeleteAdditionalDiscount(
		context.Background(),
		DealID(7),
		AdditionalDiscountID(discountID),
	)
	if err != nil {
		t.Fatalf("DeleteAdditionalDiscount error: %v", err)
	}
	if string(result.ID) != discountID {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestDealsService_ListInstallments(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/installments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["deal_ids"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
			t.Fatalf("unexpected deal_ids: %#v", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":5,"deal_id":1,"amount":20,"billing_date":"2024-01-10","description":"Installment"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	installments, next, err := client.Deals.ListInstallments(
		context.Background(),
		[]DealID{1, 2},
		WithInstallmentsPageSize(1),
	)
	if err != nil {
		t.Fatalf("ListInstallments error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(installments) != 1 || installments[0].ID != 5 {
		t.Fatalf("unexpected installments: %#v", installments)
	}
}

func TestDealsService_AddInstallment(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/installments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["amount"] != float64(20) {
			t.Fatalf("unexpected amount: %#v", payload["amount"])
		}
		if payload["billing_date"] != "2024-01-10" {
			t.Fatalf("unexpected billing_date: %#v", payload["billing_date"])
		}
		if payload["description"] != "Installment" {
			t.Fatalf("unexpected description: %#v", payload["description"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":5,"deal_id":7}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	installment, err := client.Deals.AddInstallment(
		context.Background(),
		DealID(7),
		WithInstallmentAmount(20),
		WithInstallmentBillingDate("2024-01-10"),
		WithInstallmentDescription("Installment"),
	)
	if err != nil {
		t.Fatalf("AddInstallment error: %v", err)
	}
	if installment.ID != 5 {
		t.Fatalf("unexpected installment: %#v", installment)
	}
}

func TestDealsService_UpdateInstallment(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/installments/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["amount"] != float64(25) {
			t.Fatalf("unexpected amount: %#v", payload["amount"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":5,"deal_id":7}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	installment, err := client.Deals.UpdateInstallment(
		context.Background(),
		DealID(7),
		InstallmentID(5),
		WithInstallmentAmount(25),
	)
	if err != nil {
		t.Fatalf("UpdateInstallment error: %v", err)
	}
	if installment.ID != 5 {
		t.Fatalf("unexpected installment: %#v", installment)
	}
}

func TestDealsService_DeleteInstallment(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/installments/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":5}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Deals.DeleteInstallment(
		context.Background(),
		DealID(7),
		InstallmentID(5),
	)
	if err != nil {
		t.Fatalf("DeleteInstallment error: %v", err)
	}
	if result.ID != 5 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
