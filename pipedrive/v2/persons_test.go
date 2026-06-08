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

func TestPersonsService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "2" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Ada" {
			t.Fatalf("unexpected name: %v", payload["name"])
		}
		if payload["owner_id"] != float64(3) {
			t.Fatalf("unexpected owner_id: %v", payload["owner_id"])
		}
		if payload["org_id"] != float64(5) {
			t.Fatalf("unexpected org_id: %v", payload["org_id"])
		}
		if payload["marketing_status"] != "subscribed" {
			t.Fatalf("unexpected marketing_status: %v", payload["marketing_status"])
		}
		emails, ok := payload["emails"].([]interface{})
		if !ok || len(emails) != 1 {
			t.Fatalf("unexpected emails: %#v", payload["emails"])
		}
		email, ok := emails[0].(map[string]interface{})
		if !ok || email["value"] != "ada@example.com" {
			t.Fatalf("unexpected email: %#v", emails[0])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":12,"name":"Ada"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	person, err := client.Persons.Create(
		context.Background(),
		WithPersonName("Ada"),
		WithPersonOwnerID(UserID(3)),
		WithPersonOrgID(OrganizationID(5)),
		WithPersonEmails(LabeledValue{Value: "ada@example.com", Primary: true, Label: "work"}),
		WithPersonMarketingStatus(PersonMarketingStatusSubscribed),
		WithPersonRequestOptions(pipedrive.WithHeader("X-Test", "2")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if person.ID != 12 || person.Name != "Ada" {
		t.Fatalf("unexpected person: %#v", person)
	}
}

func TestPersonsService_Search(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("term"); got != "ada" {
			t.Fatalf("unexpected term: %q", got)
		}
		if got := q.Get("fields"); got != "name,email" {
			t.Fatalf("unexpected fields: %q", got)
		}
		if got := q.Get("exact_match"); got != "true" {
			t.Fatalf("unexpected exact_match: %q", got)
		}
		if got := q.Get("organization_id"); got != "15" {
			t.Fatalf("unexpected organization_id: %q", got)
		}
		if got := q.Get("include_fields"); got != "person.picture" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c2" {
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

	results, next, err := client.Persons.Search(
		context.Background(),
		"ada",
		WithPersonSearchFields(PersonSearchFieldName, PersonSearchFieldEmail),
		WithPersonSearchExactMatch(true),
		WithPersonSearchOrganizationID(OrganizationID(15)),
		WithPersonSearchIncludeFields(PersonSearchIncludeFieldPicture),
		WithPersonSearchPageSize(2),
		WithPersonSearchCursor("c2"),
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

func TestPersonsService_ListFollowers(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/5/followers" {
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

	followers, next, err := client.Persons.ListFollowers(
		context.Background(),
		PersonID(5),
		WithPersonFollowersPageSize(2),
		WithPersonFollowersCursor("c1"),
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

func TestPersonsService_ListFollowersPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/5/followers" {
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

	pager := client.Persons.ListFollowersPager(PersonID(5), WithPersonFollowersCursor("start"))
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

func TestPersonsService_ForEachFollowers(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/5/followers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"user_id":1},{"user_id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []UserID
	err := client.Persons.ForEachFollowers(context.Background(), PersonID(5), func(follower Follower) error {
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

func TestPersonsService_AddFollower(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/5/followers" {
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

	follower, err := client.Persons.AddFollower(context.Background(), PersonID(5), UserID(7))
	if err != nil {
		t.Fatalf("AddFollower error: %v", err)
	}
	if follower.UserID != 7 {
		t.Fatalf("unexpected follower: %#v", follower)
	}
}

func TestPersonsService_DeleteFollower(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/5/followers/7" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user_id":7}}`))
	})

	result, err := client.Persons.DeleteFollower(context.Background(), PersonID(5), UserID(7))
	if err != nil {
		t.Fatalf("DeleteFollower error: %v", err)
	}
	if result.UserID != 7 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestPersonsService_FollowersChangelog(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/5/followers/changelog" {
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

	changelog, next, err := client.Persons.FollowersChangelog(
		context.Background(),
		PersonID(5),
		WithPersonFollowersChangelogPageSize(1),
		WithPersonFollowersChangelogCursor("c1"),
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

func TestPersonsService_FollowersChangelogPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/5/followers/changelog" {
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

	pager := client.Persons.FollowersChangelogPager(PersonID(5), WithPersonFollowersChangelogCursor("start"))
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

func TestPersonsService_ForEachFollowersChangelog(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/5/followers/changelog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"follower_user_id":1},{"follower_user_id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []UserID
	err := client.Persons.ForEachFollowersChangelog(context.Background(), PersonID(5), func(entry FollowerChangelog) error {
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

func TestPersonsService_GetPicture(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/5/picture" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":9,"item_type":"person","item_id":5,"pictures":{"128":"https://example.com/128.png"}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	picture, err := client.Persons.GetPicture(context.Background(), PersonID(5))
	if err != nil {
		t.Fatalf("GetPicture error: %v", err)
	}
	if picture.ID != 9 || picture.ItemID == nil || *picture.ItemID != 5 {
		t.Fatalf("unexpected picture: %#v", picture)
	}
}

func TestPersonsService_ListPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons" {
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

	pager := client.Persons.ListPager(WithPersonsPageSize(2), WithPersonsCursor("start"))
	var ids []PersonID
	for pager.Next(context.Background()) {
		for _, person := range pager.Items() {
			ids = append(ids, person.ID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestPersonsService_ForEach(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []PersonID
	err := client.Persons.ForEach(context.Background(), func(person Person) error {
		ids = append(ids, person.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("ForEach error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestPersonsService_Update(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/12" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "update" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Ada Updated" {
			t.Fatalf("unexpected name: %#v", payload["name"])
		}
		if payload["visible_to"] != float64(3) {
			t.Fatalf("unexpected visible_to: %#v", payload["visible_to"])
		}
		phones, ok := payload["phones"].([]interface{})
		if !ok || len(phones) != 1 {
			t.Fatalf("unexpected phones: %#v", payload["phones"])
		}
		phone, ok := phones[0].(map[string]interface{})
		if !ok || phone["value"] != "+123" {
			t.Fatalf("unexpected phone: %#v", phones[0])
		}
		labels, ok := payload["label_ids"].([]interface{})
		if !ok || len(labels) != 2 || labels[0] != float64(10) || labels[1] != float64(11) {
			t.Fatalf("unexpected label_ids: %#v", payload["label_ids"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":12,"name":"Ada Updated"}}`))
	})

	person, err := client.Persons.Update(
		context.Background(),
		PersonID(12),
		WithPersonName("Ada Updated"),
		WithPersonPhones(LabeledValue{Value: "+123", Primary: true, Label: "mobile"}),
		WithPersonLabelIDs(10, 11),
		WithPersonVisibleTo(3),
		WithPersonRequestOptions(pipedrive.WithHeader("X-Test", "update")),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if person.ID != 12 || person.Name != "Ada Updated" {
		t.Fatalf("unexpected person: %#v", person)
	}
}

func TestPersonsService_Delete(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/12" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":12}}`))
	})

	result, err := client.Persons.Delete(
		context.Background(),
		PersonID(12),
		WithPersonRequestOptions(pipedrive.WithHeader("X-Test", "delete")),
	)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if result.ID != 12 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
