package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestOrganizationsService_List(t *testing.T) {
	t.Parallel()

	since := time.Date(2024, 3, 4, 9, 15, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("filter_id"); got != "3" {
			t.Fatalf("unexpected filter_id: %q", got)
		}
		if got := q.Get("owner_id"); got != "4" {
			t.Fatalf("unexpected owner_id: %q", got)
		}
		if got := q.Get("updated_since"); got != since.Format(time.RFC3339) {
			t.Fatalf("unexpected updated_since: %q", got)
		}
		if got := q.Get("sort_by"); got != "update_time" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "asc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		if got := q.Get("include_fields"); got != "people_count,notes_count" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("custom_fields"); got != "cf_1" {
			t.Fatalf("unexpected custom_fields: %q", got)
		}
		if got := q.Get("ids"); got != "10,11" {
			t.Fatalf("unexpected ids: %q", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c2" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":10,"name":"Org"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	orgs, next, err := client.Organizations.List(
		context.Background(),
		WithOrganizationsFilterID(3),
		WithOrganizationsOwnerID(4),
		WithOrganizationsUpdatedSince(since),
		WithOrganizationsSortBy(OrganizationSortByUpdateTime),
		WithOrganizationsSortDirection(SortAsc),
		WithOrganizationsIncludeFields(OrganizationIncludeFieldPeopleCount, OrganizationIncludeFieldNotesCount),
		WithOrganizationsCustomFields("cf_1"),
		WithOrganizationsIDs(OrganizationID(10), OrganizationID(11)),
		WithOrganizationsPageSize(1),
		WithOrganizationsCursor("c2"),
		WithOrganizationRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(orgs) != 1 || orgs[0].ID != 10 {
		t.Fatalf("unexpected orgs: %#v", orgs)
	}
}

func TestOrganizationsService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "2" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Acme" {
			t.Fatalf("unexpected name: %v", payload["name"])
		}
		if payload["owner_id"] != float64(3) {
			t.Fatalf("unexpected owner_id: %v", payload["owner_id"])
		}
		address, ok := payload["address"].(map[string]interface{})
		if !ok || address["value"] != "HQ" {
			t.Fatalf("unexpected address: %#v", payload["address"])
		}
		customFields, ok := payload["custom_fields"].(map[string]interface{})
		if !ok || customFields["cf_key"] != "value" {
			t.Fatalf("unexpected custom_fields: %#v", payload["custom_fields"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":10,"name":"Acme"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	org, err := client.Organizations.Create(
		context.Background(),
		WithOrganizationName("Acme"),
		WithOrganizationOwnerID(UserID(3)),
		WithOrganizationAddress(OrganizationAddress{Value: "HQ"}),
		WithOrganizationCustomFieldsMap(map[string]interface{}{"cf_key": "value"}),
		WithOrganizationRequestOptions(pipedrive.WithHeader("X-Test", "2")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if org.ID != 10 || org.Name != "Acme" {
		t.Fatalf("unexpected org: %#v", org)
	}
}

func TestOrganizationsService_Search(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("term"); got != "acme" {
			t.Fatalf("unexpected term: %q", got)
		}
		if got := q.Get("fields"); got != "name,address" {
			t.Fatalf("unexpected fields: %q", got)
		}
		if got := q.Get("exact_match"); got != "true" {
			t.Fatalf("unexpected exact_match: %q", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"items":[{"result_score":0.9,"item":{"id":1}}]},"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	results, next, err := client.Organizations.Search(
		context.Background(),
		"acme",
		WithOrganizationSearchFields(OrganizationSearchFieldName, OrganizationSearchFieldAddress),
		WithOrganizationSearchExactMatch(true),
		WithOrganizationSearchPageSize(1),
		WithOrganizationSearchCursor("c1"),
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

func TestOrganizationsService_ListFollowers(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/5/followers" {
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
		_, _ = w.Write([]byte(`{"data":[{"user_id":3}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	followers, next, err := client.Organizations.ListFollowers(
		context.Background(),
		OrganizationID(5),
		WithOrganizationFollowersPageSize(2),
		WithOrganizationFollowersCursor("c1"),
	)
	if err != nil {
		t.Fatalf("ListFollowers error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(followers) != 1 || followers[0].UserID != 3 {
		t.Fatalf("unexpected followers: %#v", followers)
	}
}

func TestOrganizationsService_AddFollower(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/5/followers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["user_id"] != float64(7) {
			t.Fatalf("unexpected user_id: %v", payload["user_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user_id":7}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	follower, err := client.Organizations.AddFollower(context.Background(), OrganizationID(5), UserID(7))
	if err != nil {
		t.Fatalf("AddFollower error: %v", err)
	}
	if follower.UserID != 7 {
		t.Fatalf("unexpected follower: %#v", follower)
	}
}

func TestOrganizationsService_FollowersChangelog(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/5/followers/changelog" {
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

	changelog, next, err := client.Organizations.FollowersChangelog(
		context.Background(),
		OrganizationID(5),
		WithOrganizationFollowersChangelogPageSize(1),
		WithOrganizationFollowersChangelogCursor("c1"),
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
