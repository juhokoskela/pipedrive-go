package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
		if got := r.URL.Query().Get("include_fields"); got != "activities_count,products_count" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := r.URL.Query().Get("custom_fields"); got != "cf_1,cf_2" {
			t.Fatalf("unexpected custom_fields: %q", got)
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

	deal, err := client.Deals.Get(
		context.Background(),
		1,
		WithDealIncludeFields(DealIncludeFieldActivitiesCount, DealIncludeFieldProductsCount),
		WithDealCustomFields("cf_1", "cf_2"),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if deal.ID != 1 || deal.Title != "Test deal" {
		t.Fatalf("unexpected deal: %#v", deal)
	}
}

func TestDealsService_List(t *testing.T) {
	t.Parallel()

	since := time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)
	until := time.Date(2024, 6, 2, 11, 0, 0, 0, time.UTC)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("filter_id"); got != "3" {
			t.Fatalf("unexpected filter_id: %q", got)
		}
		if got := q.Get("ids"); got != "7,8" {
			t.Fatalf("unexpected ids: %q", got)
		}
		if got := q.Get("owner_id"); got != "4" {
			t.Fatalf("unexpected owner_id: %q", got)
		}
		if got := q.Get("person_id"); got != "5" {
			t.Fatalf("unexpected person_id: %q", got)
		}
		if got := q.Get("org_id"); got != "6" {
			t.Fatalf("unexpected org_id: %q", got)
		}
		if got := q.Get("pipeline_id"); got != "2" {
			t.Fatalf("unexpected pipeline_id: %q", got)
		}
		if got := q.Get("stage_id"); got != "9" {
			t.Fatalf("unexpected stage_id: %q", got)
		}
		if got := q.Get("status"); got != "open,won" {
			t.Fatalf("unexpected status: %q", got)
		}
		if got := q.Get("updated_since"); got != since.Format(time.RFC3339) {
			t.Fatalf("unexpected updated_since: %q", got)
		}
		if got := q.Get("updated_until"); got != until.Format(time.RFC3339) {
			t.Fatalf("unexpected updated_until: %q", got)
		}
		if got := q.Get("sort_by"); got != "update_time" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "desc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		if got := q.Get("include_fields"); got != "activities_count,products_count" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("custom_fields"); got != "cf_1,cf_2" {
			t.Fatalf("unexpected custom_fields: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "list" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":7,"title":"Deal"}],"additional_data":{"next_cursor":null}}`))
	})

	deals, next, err := client.Deals.List(
		context.Background(),
		WithDealsFilterID(3),
		WithDealsIDs(DealID(7), DealID(8)),
		WithDealsOwnerID(UserID(4)),
		WithDealsPersonID(PersonID(5)),
		WithDealsOrganizationID(OrganizationID(6)),
		WithDealsPipelineID(PipelineID(2)),
		WithDealsStageID(StageID(9)),
		WithDealsStatus(DealStatusOpen, DealStatusWon),
		WithDealsUpdatedSince(since),
		WithDealsUpdatedUntil(until),
		WithDealsSortBy(DealSortByUpdateTime),
		WithDealsSortDirection(SortDesc),
		WithDealsIncludeFields(DealIncludeFieldActivitiesCount, DealIncludeFieldProductsCount),
		WithDealsCustomFields("cf_1", "cf_2"),
		WithDealsPageSize(2),
		WithDealsCursor("c1"),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "list")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(deals) != 1 || deals[0].ID != 7 {
		t.Fatalf("unexpected deals: %#v", deals)
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
		if got := r.Header.Get("X-Test"); got != "create" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		if payload["owner_id"] != float64(3) {
			t.Fatalf("unexpected owner_id: %#v", payload["owner_id"])
		}
		if payload["person_id"] != float64(4) {
			t.Fatalf("unexpected person_id: %#v", payload["person_id"])
		}
		if payload["org_id"] != float64(5) {
			t.Fatalf("unexpected org_id: %#v", payload["org_id"])
		}
		if payload["stage_id"] != float64(6) {
			t.Fatalf("unexpected stage_id: %#v", payload["stage_id"])
		}
		if payload["pipeline_id"] != float64(2) {
			t.Fatalf("unexpected pipeline_id: %#v", payload["pipeline_id"])
		}
		if payload["status"] != "open" {
			t.Fatalf("unexpected status: %#v", payload["status"])
		}
		if payload["expected_close_date"] != "2024-06-10" {
			t.Fatalf("unexpected expected_close_date: %#v", payload["expected_close_date"])
		}
		if payload["probability"] != float64(55) {
			t.Fatalf("unexpected probability: %#v", payload["probability"])
		}
		if payload["lost_reason"] != "No budget" {
			t.Fatalf("unexpected lost_reason: %#v", payload["lost_reason"])
		}
		if payload["visible_to"] != float64(3) {
			t.Fatalf("unexpected visible_to: %#v", payload["visible_to"])
		}
		labels, ok := payload["label_ids"].([]interface{})
		if !ok || len(labels) != 2 || labels[0] != float64(10) || labels[1] != float64(11) {
			t.Fatalf("unexpected label_ids: %#v", payload["label_ids"])
		}
		customFields, ok := payload["custom_fields"].(map[string]interface{})
		if !ok || customFields["cf_key"] != "value" {
			t.Fatalf("unexpected custom_fields: %#v", payload["custom_fields"])
		}
		if payload["is_archived"] != true {
			t.Fatalf("unexpected is_archived: %#v", payload["is_archived"])
		}
		if payload["is_deleted"] != false {
			t.Fatalf("unexpected is_deleted: %#v", payload["is_deleted"])
		}
		if payload["archive_time"] != "2024-06-11T10:00:00Z" {
			t.Fatalf("unexpected archive_time: %#v", payload["archive_time"])
		}
		if payload["close_time"] != "2024-06-12T10:00:00Z" {
			t.Fatalf("unexpected close_time: %#v", payload["close_time"])
		}
		if payload["lost_time"] != "2024-06-13T10:00:00Z" {
			t.Fatalf("unexpected lost_time: %#v", payload["lost_time"])
		}
		if payload["won_time"] != "2024-06-14T10:00:00Z" {
			t.Fatalf("unexpected won_time: %#v", payload["won_time"])
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
		WithDealOwnerID(UserID(3)),
		WithDealPersonID(PersonID(4)),
		WithDealOrganizationID(OrganizationID(5)),
		WithDealStageID(StageID(6)),
		WithDealPipelineID(PipelineID(2)),
		WithDealStatus(DealStatusOpen),
		WithDealExpectedCloseDate("2024-06-10"),
		WithDealProbability(55),
		WithDealLostReason("No budget"),
		WithDealVisibleTo(3),
		WithDealLabelIDs(10, 11),
		WithDealCustomFieldsMap(map[string]interface{}{"cf_key": "value"}),
		WithDealArchived(true),
		WithDealDeleted(false),
		WithDealArchiveTime("2024-06-11T10:00:00Z"),
		WithDealCloseTime("2024-06-12T10:00:00Z"),
		WithDealLostTime("2024-06-13T10:00:00Z"),
		WithDealWonTime("2024-06-14T10:00:00Z"),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "create")),
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
		if got := r.Header.Get("X-Test"); got != "update" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "update")),
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
		if got := r.Header.Get("X-Test"); got != "delete" {
			t.Fatalf("unexpected header X-Test: %q", got)
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

	result, err := client.Deals.Delete(
		context.Background(),
		DealID(7),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "delete")),
	)
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
		if got := q.Get("person_id"); got != "5" {
			t.Fatalf("unexpected person_id: %q", got)
		}
		if got := q.Get("organization_id"); got != "6" {
			t.Fatalf("unexpected organization_id: %q", got)
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
		if got := r.Header.Get("X-Test"); got != "search" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealSearchPersonID(PersonID(5)),
		WithDealSearchOrganizationID(OrganizationID(6)),
		WithDealSearchIncludeFields(DealSearchIncludeFieldDealCCEmail),
		WithDealSearchPageSize(1),
		WithDealSearchCursor("c1"),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "search")),
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

	since := time.Date(2024, 4, 1, 9, 0, 0, 0, time.UTC)
	until := time.Date(2024, 4, 2, 9, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/archived" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("filter_id"); got != "4" {
			t.Fatalf("unexpected filter_id: %q", got)
		}
		if got := q.Get("ids"); got != "2,3" {
			t.Fatalf("unexpected ids: %q", got)
		}
		if got := q.Get("owner_id"); got != "3" {
			t.Fatalf("unexpected owner_id: %q", got)
		}
		if got := q.Get("person_id"); got != "5" {
			t.Fatalf("unexpected person_id: %q", got)
		}
		if got := q.Get("org_id"); got != "6" {
			t.Fatalf("unexpected org_id: %q", got)
		}
		if got := q.Get("pipeline_id"); got != "7" {
			t.Fatalf("unexpected pipeline_id: %q", got)
		}
		if got := q.Get("stage_id"); got != "8" {
			t.Fatalf("unexpected stage_id: %q", got)
		}
		if got := q.Get("status"); got != "won,lost" {
			t.Fatalf("unexpected status: %q", got)
		}
		if got := q.Get("updated_since"); got != since.Format(time.RFC3339) {
			t.Fatalf("unexpected updated_since: %q", got)
		}
		if got := q.Get("updated_until"); got != until.Format(time.RFC3339) {
			t.Fatalf("unexpected updated_until: %q", got)
		}
		if got := q.Get("sort_by"); got != "update_time" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "desc" {
			t.Fatalf("unexpected sort_direction: %q", got)
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
		if got := r.Header.Get("X-Test"); got != "archived" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithArchivedDealsFilterID(4),
		WithArchivedDealsIDs(DealID(2), DealID(3)),
		WithArchivedDealsOwnerID(UserID(3)),
		WithArchivedDealsPersonID(PersonID(5)),
		WithArchivedDealsOrganizationID(OrganizationID(6)),
		WithArchivedDealsPipelineID(PipelineID(7)),
		WithArchivedDealsStageID(StageID(8)),
		WithArchivedDealsStatus(DealStatusWon, DealStatusLost),
		WithArchivedDealsUpdatedSince(since),
		WithArchivedDealsUpdatedUntil(until),
		WithArchivedDealsSortBy(DealSortByUpdateTime),
		WithArchivedDealsSortDirection(SortDesc),
		WithArchivedDealsIncludeFields(DealIncludeFieldActivitiesCount),
		WithArchivedDealsCustomFields("cf_1"),
		WithArchivedDealsPageSize(2),
		WithArchivedDealsCursor("c1"),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "archived")),
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

func TestDealsService_ListArchivedPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/archived" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		switch calls {
		case 1:
			if got := r.URL.Query().Get("cursor"); got != "start" {
				t.Fatalf("unexpected first cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1}],"additional_data":{"next_cursor":"next"}}`))
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "next" {
				t.Fatalf("unexpected second cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":2}],"additional_data":{"next_cursor":null}}`))
		default:
			t.Fatalf("unexpected call count: %d", calls)
		}
	})

	pager := client.Deals.ListArchivedPager(WithArchivedDealsPageSize(2), WithArchivedDealsCursor("start"))
	var ids []DealID
	for pager.Next(context.Background()) {
		for _, deal := range pager.Items() {
			ids = append(ids, deal.ID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestDealsService_ForEachArchived(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/archived" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []DealID
	err := client.Deals.ForEachArchived(context.Background(), func(deal Deal) error {
		ids = append(ids, deal.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("ForEachArchived error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
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
		if got := r.Header.Get("X-Test"); got != "convert" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"conversion_id":"` + convID + `"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	job, err := client.Deals.ConvertToLead(
		context.Background(),
		DealID(7),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "convert")),
	)
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
		if got := r.Header.Get("X-Test"); got != "conversion-status" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"conversion_id":"` + convID + `","status":"completed","deal_id":9,"lead_id":"` + leadID + `"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	status, err := client.Deals.ConversionStatus(
		context.Background(),
		DealID(7),
		ConversionID(convID),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "conversion-status")),
	)
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
		if got := r.Header.Get("X-Test"); got != "followers" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "followers")),
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

func TestDealsService_ListFollowersPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/followers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		switch calls {
		case 1:
			if got := r.URL.Query().Get("cursor"); got != "start" {
				t.Fatalf("unexpected first cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"user_id":1}],"additional_data":{"next_cursor":"next"}}`))
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "next" {
				t.Fatalf("unexpected second cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"user_id":2}],"additional_data":{"next_cursor":null}}`))
		default:
			t.Fatalf("unexpected call count: %d", calls)
		}
	})

	pager := client.Deals.ListFollowersPager(DealID(7), WithDealFollowersCursor("start"))
	var ids []UserID
	for pager.Next(context.Background()) {
		for _, follower := range pager.Items() {
			ids = append(ids, follower.UserID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestDealsService_ForEachFollowers(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/followers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"user_id":1},{"user_id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []UserID
	err := client.Deals.ForEachFollowers(context.Background(), DealID(7), func(follower Follower) error {
		ids = append(ids, follower.UserID)
		return nil
	})
	if err != nil {
		t.Fatalf("ForEachFollowers error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
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
		if got := r.Header.Get("X-Test"); got != "add-follower" {
			t.Fatalf("unexpected header X-Test: %q", got)
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

	follower, err := client.Deals.AddFollower(
		context.Background(),
		DealID(7),
		UserID(6),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "add-follower")),
	)
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
		if got := r.Header.Get("X-Test"); got != "delete-follower" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user_id":6}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Deals.DeleteFollower(
		context.Background(),
		DealID(7),
		UserID(6),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "delete-follower")),
	)
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
		if got := r.Header.Get("X-Test"); got != "followers-changelog" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "followers-changelog")),
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

func TestDealsService_FollowersChangelogPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/followers/changelog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		switch calls {
		case 1:
			if got := r.URL.Query().Get("cursor"); got != "start" {
				t.Fatalf("unexpected first cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"follower_user_id":1}],"additional_data":{"next_cursor":"next"}}`))
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "next" {
				t.Fatalf("unexpected second cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"follower_user_id":2}],"additional_data":{"next_cursor":null}}`))
		default:
			t.Fatalf("unexpected call count: %d", calls)
		}
	})

	pager := client.Deals.FollowersChangelogPager(DealID(7), WithDealFollowersChangelogCursor("start"))
	var ids []UserID
	for pager.Next(context.Background()) {
		for _, entry := range pager.Items() {
			ids = append(ids, entry.FollowerUserID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestDealsService_ForEachFollowersChangelog(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/followers/changelog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"follower_user_id":1},{"follower_user_id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []UserID
	err := client.Deals.ForEachFollowersChangelog(context.Background(), DealID(7), func(entry FollowerChangelog) error {
		ids = append(ids, entry.FollowerUserID)
		return nil
	})
	if err != nil {
		t.Fatalf("ForEachFollowersChangelog error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
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
		if got := r.Header.Get("X-Test"); got != "deal-products" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "deal-products")),
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

func TestDealsService_ListProductsPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		switch calls {
		case 1:
			if got := r.URL.Query().Get("cursor"); got != "start" {
				t.Fatalf("unexpected first cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1}],"additional_data":{"next_cursor":"next"}}`))
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "next" {
				t.Fatalf("unexpected second cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":2}],"additional_data":{"next_cursor":null}}`))
		default:
			t.Fatalf("unexpected call count: %d", calls)
		}
	})

	pager := client.Deals.ListProductsPager(DealID(7), WithDealProductsCursor("start"))
	var ids []DealProductAttachmentID
	for pager.Next(context.Background()) {
		for _, product := range pager.Items() {
			ids = append(ids, product.ID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestDealsService_ForEachProducts(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []DealProductAttachmentID
	err := client.Deals.ForEachProducts(context.Background(), DealID(7), func(product DealProduct) error {
		ids = append(ids, product.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("ForEachProducts error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
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
		if got := q.Get("sort_by"); got != "deal_id" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "asc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "deals-products" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealsProductsSortBy(DealProductSortByDealID),
		WithDealsProductsSortDirection(SortAsc),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "deals-products")),
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

func TestDealsService_ListProductsAcrossDealsPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["deal_ids"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
			t.Fatalf("unexpected deal_ids: %#v", got)
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		switch calls {
		case 1:
			if got := r.URL.Query().Get("cursor"); got != "start" {
				t.Fatalf("unexpected first cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1}],"additional_data":{"next_cursor":"next"}}`))
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "next" {
				t.Fatalf("unexpected second cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":2}],"additional_data":{"next_cursor":null}}`))
		default:
			t.Fatalf("unexpected call count: %d", calls)
		}
	})

	pager := client.Deals.ListProductsAcrossDealsPager([]DealID{1, 2}, WithDealsProductsCursor("start"))
	var ids []DealProductAttachmentID
	for pager.Next(context.Background()) {
		for _, product := range pager.Items() {
			ids = append(ids, product.ID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestDealsService_ForEachProductsAcrossDeals(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["deal_ids"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
			t.Fatalf("unexpected deal_ids: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []DealProductAttachmentID
	err := client.Deals.ForEachProductsAcrossDeals(context.Background(), []DealID{1, 2}, func(product DealProduct) error {
		ids = append(ids, product.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("ForEachProductsAcrossDeals error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
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
		if got := r.Header.Get("X-Test"); got != "add-product" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["product_id"] != float64(9) {
			t.Fatalf("unexpected product_id: %#v", payload["product_id"])
		}
		if payload["product_variation_id"] != float64(10) {
			t.Fatalf("unexpected product_variation_id: %#v", payload["product_variation_id"])
		}
		if payload["item_price"] != float64(199) {
			t.Fatalf("unexpected item_price: %#v", payload["item_price"])
		}
		if payload["quantity"] != float64(2) {
			t.Fatalf("unexpected quantity: %#v", payload["quantity"])
		}
		if payload["discount"] != float64(5) {
			t.Fatalf("unexpected discount: %#v", payload["discount"])
		}
		if payload["discount_type"] != "percentage" {
			t.Fatalf("unexpected discount_type: %#v", payload["discount_type"])
		}
		if payload["comments"] != "Launch bundle" {
			t.Fatalf("unexpected comments: %#v", payload["comments"])
		}
		if payload["tax"] != float64(20) {
			t.Fatalf("unexpected tax: %#v", payload["tax"])
		}
		if payload["tax_method"] != "exclusive" {
			t.Fatalf("unexpected tax_method: %#v", payload["tax_method"])
		}
		if payload["is_enabled"] != true {
			t.Fatalf("unexpected is_enabled: %#v", payload["is_enabled"])
		}
		if payload["billing_frequency"] != "monthly" {
			t.Fatalf("unexpected billing_frequency: %#v", payload["billing_frequency"])
		}
		if payload["billing_frequency_cycles"] != float64(3) {
			t.Fatalf("unexpected billing_frequency_cycles: %#v", payload["billing_frequency_cycles"])
		}
		if payload["billing_start_date"] != "2024-07-01" {
			t.Fatalf("unexpected billing_start_date: %#v", payload["billing_start_date"])
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
		WithDealProductVariationID(ProductVariationID(10)),
		WithDealProductItemPrice(199),
		WithDealProductQuantity(2),
		WithDealProductDiscount(5),
		WithDealProductDiscountType(DealProductDiscountTypePercentage),
		WithDealProductComments("Launch bundle"),
		WithDealProductTax(20),
		WithDealProductTaxMethod(DealProductTaxMethodExclusive),
		WithDealProductIsEnabled(true),
		WithDealProductBillingFrequency(BillingFrequencyMonthly),
		WithDealProductBillingFrequencyCycles(3),
		WithDealProductBillingStartDate("2024-07-01"),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "add-product")),
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
		if got := r.Header.Get("X-Test"); got != "add-products" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "add-products")),
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
		if got := r.Header.Get("X-Test"); got != "update-product" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "update-product")),
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
		if got := r.Header.Get("X-Test"); got != "delete-product" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":15}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Deals.DeleteProduct(
		context.Background(),
		DealID(7),
		DealProductAttachmentID(15),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "delete-product")),
	)
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
		if got := r.Header.Get("X-Test"); got != "delete-products" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "delete-products")),
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
		if got := r.Header.Get("X-Test"); got != "discounts" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"11111111-1111-1111-1111-111111111111","deal_id":7,"amount":10,"type":"amount"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	discounts, err := client.Deals.ListAdditionalDiscounts(
		context.Background(),
		DealID(7),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "discounts")),
	)
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
		if got := r.Header.Get("X-Test"); got != "add-discount" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "add-discount")),
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
		if got := r.Header.Get("X-Test"); got != "update-discount" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "update-discount")),
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
		if got := r.Header.Get("X-Test"); got != "delete-discount" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "delete-discount")),
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
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := q.Get("sort_by"); got != "billing_date" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "desc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "installments" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithInstallmentsCursor("c1"),
		WithInstallmentsSortBy(InstallmentSortByBillingDate),
		WithInstallmentsSortDirection(SortDesc),
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "installments")),
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

func TestDealsService_ListInstallmentsPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/installments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["deal_ids"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
			t.Fatalf("unexpected deal_ids: %#v", got)
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		switch calls {
		case 1:
			if got := r.URL.Query().Get("cursor"); got != "start" {
				t.Fatalf("unexpected first cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1}],"additional_data":{"next_cursor":"next"}}`))
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "next" {
				t.Fatalf("unexpected second cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":2}],"additional_data":{"next_cursor":null}}`))
		default:
			t.Fatalf("unexpected call count: %d", calls)
		}
	})

	pager := client.Deals.ListInstallmentsPager([]DealID{1, 2}, WithInstallmentsCursor("start"))
	var ids []InstallmentID
	for pager.Next(context.Background()) {
		for _, installment := range pager.Items() {
			ids = append(ids, installment.ID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestDealsService_ForEachInstallments(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/installments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["deal_ids"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
			t.Fatalf("unexpected deal_ids: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []InstallmentID
	err := client.Deals.ForEachInstallments(context.Background(), []DealID{1, 2}, func(installment Installment) error {
		ids = append(ids, installment.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("ForEachInstallments error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
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
		if got := r.Header.Get("X-Test"); got != "add-installment" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "add-installment")),
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
		if got := r.Header.Get("X-Test"); got != "update-installment" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "update-installment")),
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
		if got := r.Header.Get("X-Test"); got != "delete-installment" {
			t.Fatalf("unexpected header X-Test: %q", got)
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
		WithDealRequestOptions(pipedrive.WithHeader("X-Test", "delete-installment")),
	)
	if err != nil {
		t.Fatalf("DeleteInstallment error: %v", err)
	}
	if result.ID != 5 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
