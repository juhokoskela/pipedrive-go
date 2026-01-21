package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func TestUsersService_Create(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/users" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["name"] != "Jane" {
			t.Fatalf("unexpected name: %#v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":1,"name":"Jane"}}`))
	})

	user, err := client.Users.Create(context.Background(), map[string]any{"name": "Jane"})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if user.ID != 1 {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestUsersService_Update(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/users/2" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["name"] != "Updated" {
			t.Fatalf("unexpected name: %#v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":2,"name":"Updated"}}`))
	})

	user, err := client.Users.Update(context.Background(), UserID(2), map[string]any{"name": "Updated"})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if user.ID != 2 {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestUsersService_FindByName(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/users/find" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("term"); got != "Jane" {
			t.Fatalf("unexpected term: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Jane"}]}`))
	})

	query := url.Values{}
	query.Set("term", "Jane")
	users, err := client.Users.FindByName(context.Background(), query)
	if err != nil {
		t.Fatalf("FindByName error: %v", err)
	}
	if len(users) != 1 || users[0].ID != 1 {
		t.Fatalf("unexpected users: %#v", users)
	}
}

func TestUsersService_ListRoleAssignments(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/users/3/roleAssignments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("start"); got != "1" {
			t.Fatalf("unexpected start: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"role_id":2}]}`))
	})

	query := url.Values{}
	query.Set("start", "1")
	assignments, err := client.Users.ListRoleAssignments(context.Background(), UserID(3), query)
	if err != nil {
		t.Fatalf("ListRoleAssignments error: %v", err)
	}
	if len(assignments) != 1 {
		t.Fatalf("unexpected assignments: %#v", assignments)
	}
}

func TestUsersService_ListRoleSettings(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/users/4/roleSettings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"key":"deal_default_visibility"}]}`))
	})

	settings, err := client.Users.ListRoleSettings(context.Background(), UserID(4))
	if err != nil {
		t.Fatalf("ListRoleSettings error: %v", err)
	}
	if len(settings) != 1 {
		t.Fatalf("unexpected settings: %#v", settings)
	}
}

func TestUsersService_ListTeams(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/legacyTeams/user/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Team"}]}`))
	})

	teams, err := client.Users.ListTeams(context.Background(), UserID(5), nil)
	if err != nil {
		t.Fatalf("ListTeams error: %v", err)
	}
	if len(teams) != 1 || teams[0].ID != 1 {
		t.Fatalf("unexpected teams: %#v", teams)
	}
}
