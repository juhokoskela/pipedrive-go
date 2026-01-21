package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestPermissionSetsService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/permissionSets" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("app"); got != "sales" {
			t.Fatalf("unexpected app: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":"perm-1","name":"Sales admin","description":"desc","app":"sales","type":"admin","assignment_count":2}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sets, err := client.PermissionSets.List(
		context.Background(),
		WithPermissionSetsApp(PermissionSetAppSales),
		WithPermissionSetsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(sets) != 1 || sets[0].ID != "perm-1" || sets[0].AssignmentCount != 2 {
		t.Fatalf("unexpected permission sets: %#v", sets)
	}
}

func TestPermissionSetsService_Get(t *testing.T) {
	t.Parallel()

	id := "perm-1"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/permissionSets/"+id {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"` + id + `","name":"Sales admin","description":"desc","app":"sales","type":"admin","assignment_count":2,"contents":["can_add_deals"]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	set, err := client.PermissionSets.Get(
		context.Background(),
		PermissionSetID(id),
		WithPermissionSetsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if set.ID != PermissionSetID(id) || len(set.Contents) != 1 {
		t.Fatalf("unexpected permission set: %#v", set)
	}
}

func TestPermissionSetsService_ListAssignments(t *testing.T) {
	t.Parallel()

	id := "perm-1"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/permissionSets/"+id+"/assignments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("start"); got != "2" {
			t.Fatalf("unexpected start: %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"user_id":10,"permission_set_id":"` + id + `","name":"Sales admin"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	assignments, err := client.PermissionSets.ListAssignments(
		context.Background(),
		PermissionSetID(id),
		WithPermissionSetAssignmentsStart(2),
		WithPermissionSetAssignmentsLimit(1),
	)
	if err != nil {
		t.Fatalf("ListAssignments error: %v", err)
	}
	if len(assignments) != 1 || assignments[0].UserID != 10 {
		t.Fatalf("unexpected assignments: %#v", assignments)
	}
}
