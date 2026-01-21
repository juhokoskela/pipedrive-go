package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestUsersService_GetPermissions(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/users/7/permissions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"can_add_products":true,"can_use_api":false,"can_merge_people":true}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	perms, err := client.Users.GetPermissions(
		context.Background(),
		UserID(7),
		WithUsersRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("GetPermissions error: %v", err)
	}
	if !perms.CanAddProducts || perms.CanUseAPI || !perms.CanMergePeople {
		t.Fatalf("unexpected permissions: %#v", perms)
	}
}

func TestUsersService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/users" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Jane Doe","email":"jane@example.com","last_login":"2024-01-01 10:00:00","active_flag":true,"access":[{"app":"sales","admin":true,"permission_set_id":"perm-1"}]}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	users, err := client.Users.List(context.Background(), WithUsersRequestOptions(pipedrive.WithHeader("X-Test", "1")))
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(users) != 1 || users[0].ID != 1 || users[0].Email != "jane@example.com" {
		t.Fatalf("unexpected users: %#v", users)
	}
	if users[0].LastLogin == nil || users[0].LastLogin.IsZero() {
		t.Fatalf("expected last_login to be parsed")
	}
	if len(users[0].Access) != 1 || users[0].Access[0].PermissionSetID != "perm-1" {
		t.Fatalf("unexpected access: %#v", users[0].Access)
	}
}

func TestUsersService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/users/9" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":9,"name":"Sam","email":"sam@example.com","active_flag":true}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	user, err := client.Users.Get(context.Background(), UserID(9))
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if user.ID != 9 || user.Name != "Sam" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestUsersService_GetCurrent(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/users/me" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":1,"name":"Me","email":"me@example.com","company_id":42,"company_name":"Acme","language":{"language_code":"en","country_code":"US"}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	user, err := client.Users.GetCurrent(context.Background())
	if err != nil {
		t.Fatalf("GetCurrent error: %v", err)
	}
	if user.ID != 1 || user.CompanyID != 42 || user.Language == nil || user.Language.LanguageCode != "en" {
		t.Fatalf("unexpected user: %#v", user)
	}
}
