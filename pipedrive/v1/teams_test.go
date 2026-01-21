package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
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

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Team"}]}`))
	})

	teams, err := client.Teams.List(context.Background())
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

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":2,"name":"Team"}}`))
	})

	team, err := client.Teams.Get(context.Background(), TeamID(2))
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

	team, err := client.Teams.Create(context.Background(), map[string]any{"name": "Sales"})
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

	team, err := client.Teams.Update(context.Background(), TeamID(4), map[string]any{"name": "Updated"})
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

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[1,2]}`))
	})

	users, err := client.Teams.ListUsers(context.Background(), TeamID(5))
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

	users, err := client.Teams.AddUsers(context.Background(), TeamID(6), []UserID{1, 2})
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

	users, err := client.Teams.DeleteUsers(context.Background(), TeamID(7), []UserID{9})
	if err != nil {
		t.Fatalf("DeleteUsers error: %v", err)
	}
	if len(users) != 1 || users[0] != 9 {
		t.Fatalf("unexpected users: %#v", users)
	}
}
