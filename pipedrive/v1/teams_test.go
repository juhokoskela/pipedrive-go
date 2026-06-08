package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestTeamsService_List(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/legacyTeams" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("status"); got != "active" {
			t.Fatalf("unexpected status: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "list" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Team"}]}`))
	})

	query := url.Values{}
	query.Set("status", "active")
	teams, err := client.Teams.List(
		context.Background(),
		WithTeamsQuery(query),
		WithTeamsRequestOptions(pipedrive.WithHeader("X-Test", "list")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(teams) != 1 || teams[0].ID != 1 {
		t.Fatalf("unexpected teams: %#v", teams)
	}
}

func TestTeamsService_Get(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/legacyTeams/2" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "get" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":2,"name":"Team"}}`))
	})

	team, err := client.Teams.Get(
		context.Background(),
		TeamID(2),
		WithTeamsRequestOptions(pipedrive.WithHeader("X-Test", "get")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if team.ID != 2 {
		t.Fatalf("unexpected team: %#v", team)
	}
}

func TestTeamsService_Create(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/legacyTeams" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "create" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["name"] != "Sales" {
			t.Fatalf("unexpected name: %#v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":3,"name":"Sales"}}`))
	})

	team, err := client.Teams.Create(
		context.Background(),
		map[string]any{"name": "Sales"},
		WithTeamsRequestOptions(pipedrive.WithHeader("X-Test", "create")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if team.ID != 3 {
		t.Fatalf("unexpected team: %#v", team)
	}
}

func TestTeamsService_Update(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/legacyTeams/4" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "update" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["name"] != "Updated" {
			t.Fatalf("unexpected name: %#v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":4,"name":"Updated"}}`))
	})

	team, err := client.Teams.Update(
		context.Background(),
		TeamID(4),
		map[string]any{"name": "Updated"},
		WithTeamsRequestOptions(pipedrive.WithHeader("X-Test", "update")),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if team.ID != 4 {
		t.Fatalf("unexpected team: %#v", team)
	}
}

func TestTeamsService_ListUsers(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/legacyTeams/5/users" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "list-users" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[1,2]}`))
	})

	users, err := client.Teams.ListUsers(
		context.Background(),
		TeamID(5),
		WithTeamsRequestOptions(pipedrive.WithHeader("X-Test", "list-users")),
	)
	if err != nil {
		t.Fatalf("ListUsers error: %v", err)
	}
	if len(users) != 2 || users[0] != 1 {
		t.Fatalf("unexpected users: %#v", users)
	}
}

func TestTeamsService_AddUsers(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/legacyTeams/6/users" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "add-users" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var body map[string][]int
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(body["users"]) != 2 {
			t.Fatalf("unexpected users: %#v", body["users"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[1,2]}`))
	})

	users, err := client.Teams.AddUsers(
		context.Background(),
		TeamID(6),
		[]UserID{1, 2},
		WithTeamsRequestOptions(pipedrive.WithHeader("X-Test", "add-users")),
	)
	if err != nil {
		t.Fatalf("AddUsers error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("unexpected users: %#v", users)
	}
}

func TestTeamsService_DeleteUsers(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/legacyTeams/7/users" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete-users" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var body map[string][]int
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(body["users"]) != 1 || body["users"][0] != 9 {
			t.Fatalf("unexpected users: %#v", body["users"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[9]}`))
	})

	users, err := client.Teams.DeleteUsers(
		context.Background(),
		TeamID(7),
		[]UserID{9},
		WithTeamsRequestOptions(pipedrive.WithHeader("X-Test", "delete-users")),
	)
	if err != nil {
		t.Fatalf("DeleteUsers error: %v", err)
	}
	if len(users) != 1 || users[0] != 9 {
		t.Fatalf("unexpected users: %#v", users)
	}
}
