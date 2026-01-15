package v2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestPersonsService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/42" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("include_fields"); got != "activities_count,files_count" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := r.URL.Query().Get("custom_fields"); got != "cf_1,cf_2" {
			t.Fatalf("unexpected custom_fields: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":42,"name":"Test person"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	person, err := client.Persons.Get(
		context.Background(),
		42,
		WithPersonIncludeFields(PersonIncludeFieldActivitiesCount, PersonIncludeFieldFilesCount),
		WithPersonCustomFields("cf_1", "cf_2"),
		WithPersonRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if person.ID != 42 || person.Name != "Test person" {
		t.Fatalf("unexpected person: %#v", person)
	}
}

func TestPersonsService_List(t *testing.T) {
	t.Parallel()

	since := time.Date(2024, 2, 1, 10, 20, 0, 0, time.UTC)
	until := time.Date(2024, 2, 2, 11, 30, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("filter_id"); got != "7" {
			t.Fatalf("unexpected filter_id: %q", got)
		}
		if got := q.Get("owner_id"); got != "9" {
			t.Fatalf("unexpected owner_id: %q", got)
		}
		if got := q.Get("org_id"); got != "11" {
			t.Fatalf("unexpected org_id: %q", got)
		}
		if got := q.Get("deal_id"); got != "22" {
			t.Fatalf("unexpected deal_id: %q", got)
		}
		if got := q.Get("updated_since"); got != since.Format(time.RFC3339) {
			t.Fatalf("unexpected updated_since: %q", got)
		}
		if got := q.Get("updated_until"); got != until.Format(time.RFC3339) {
			t.Fatalf("unexpected updated_until: %q", got)
		}
		if got := q.Get("sort_by"); got != "add_time" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "desc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		if got := q.Get("include_fields"); got != "activities_count,files_count" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("custom_fields"); got != "cf_1,cf_2" {
			t.Fatalf("unexpected custom_fields: %q", got)
		}
		if got := q.Get("ids"); got != "1,2,3" {
			t.Fatalf("unexpected ids: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1,"name":"Person"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	people, next, err := client.Persons.List(
		context.Background(),
		WithPersonsFilterID(7),
		WithPersonsOwnerID(9),
		WithPersonsOrgID(OrganizationID(11)),
		WithPersonsDealID(DealID(22)),
		WithPersonsUpdatedSince(since),
		WithPersonsUpdatedUntil(until),
		WithPersonsSortBy(PersonSortByAddTime),
		WithPersonsSortDirection(SortDesc),
		WithPersonsIncludeFields(PersonIncludeFieldActivitiesCount, PersonIncludeFieldFilesCount),
		WithPersonsCustomFields("cf_1", "cf_2"),
		WithPersonsIDs(PersonID(1), PersonID(2), PersonID(3)),
		WithPersonsPageSize(2),
		WithPersonsCursor("c1"),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(people) != 1 || people[0].ID != 1 {
		t.Fatalf("unexpected people: %#v", people)
	}
}
