package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestRolesService_List(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("include"); got != "settings" {
			t.Fatalf("unexpected include: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "list" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Role"}]}`))
	})

	query := url.Values{}
	query.Set("include", "settings")
	roles, err := client.Roles.List(
		context.Background(),
		WithRolesQuery(query),
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "list")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(roles) != 1 || roles[0].ID != 1 {
		t.Fatalf("unexpected roles: %#v", roles)
	}
}

func TestRolesService_Get(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/2" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "get" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":2,"name":"Role"}}`))
	})

	role, err := client.Roles.Get(
		context.Background(),
		RoleID(2),
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "get")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if role.ID != 2 {
		t.Fatalf("unexpected role: %#v", role)
	}
}

func TestRolesService_Create(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles" {
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

	role, err := client.Roles.Create(
		context.Background(),
		map[string]any{"name": "Sales"},
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "create")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if role.ID != 3 {
		t.Fatalf("unexpected role: %#v", role)
	}
}

func TestRolesService_Update(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/4" {
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

	role, err := client.Roles.Update(
		context.Background(),
		RoleID(4),
		map[string]any{"name": "Updated"},
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "update")),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if role.ID != 4 {
		t.Fatalf("unexpected role: %#v", role)
	}
}

func TestRolesService_Delete(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	ok, err := client.Roles.Delete(
		context.Background(),
		RoleID(5),
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "delete")),
	)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok")
	}
}

func TestRolesService_ListAssignments(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/6/assignments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "assignments" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"user_id":1,"role_id":6}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	assignments, page, err := client.Roles.ListAssignments(
		context.Background(),
		RoleID(6),
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "assignments")),
	)
	if err != nil {
		t.Fatalf("ListAssignments error: %v", err)
	}
	if len(assignments) != 1 || assignments[0].UserID != 1 {
		t.Fatalf("unexpected assignments: %#v", assignments)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestRolesService_AddAssignment(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/7/assignments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "add-assignment" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["user_id"] != float64(2) {
			t.Fatalf("unexpected user_id: %#v", body["user_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"user_id":2,"role_id":7}}`))
	})

	assignment, err := client.Roles.AddAssignment(
		context.Background(),
		RoleID(7),
		UserID(2),
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "add-assignment")),
	)
	if err != nil {
		t.Fatalf("AddAssignment error: %v", err)
	}
	if assignment.UserID != 2 {
		t.Fatalf("unexpected assignment: %#v", assignment)
	}
}

func TestRolesService_DeleteAssignment(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/7/assignments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete-assignment" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["user_id"] != float64(2) {
			t.Fatalf("unexpected user_id: %#v", body["user_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	ok, err := client.Roles.DeleteAssignment(
		context.Background(),
		RoleID(7),
		UserID(2),
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "delete-assignment")),
	)
	if err != nil {
		t.Fatalf("DeleteAssignment error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok")
	}
}

func TestRolesService_ListPipelines(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/8/pipelines" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "pipelines" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"pipeline_id":1}]}`))
	})

	pipelines, err := client.Roles.ListPipelines(
		context.Background(),
		RoleID(8),
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "pipelines")),
	)
	if err != nil {
		t.Fatalf("ListPipelines error: %v", err)
	}
	if len(pipelines) != 1 {
		t.Fatalf("unexpected pipelines: %#v", pipelines)
	}
}

func TestRolesService_UpdatePipelines(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/9/pipelines" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "update-pipelines" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["pipelines"] == nil {
			t.Fatalf("expected pipelines payload")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"updated":true}}`))
	})

	data, err := client.Roles.UpdatePipelines(
		context.Background(),
		RoleID(9),
		map[string]any{"pipelines": []int{1}},
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "update-pipelines")),
	)
	if err != nil {
		t.Fatalf("UpdatePipelines error: %v", err)
	}
	if data["updated"] != true {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestRolesService_ListSettings(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/10/settings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "settings" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"key":"deal_default_visibility"}]}`))
	})

	settings, err := client.Roles.ListSettings(
		context.Background(),
		RoleID(10),
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "settings")),
	)
	if err != nil {
		t.Fatalf("ListSettings error: %v", err)
	}
	if len(settings) != 1 {
		t.Fatalf("unexpected settings: %#v", settings)
	}
}

func TestRolesService_UpsertSetting(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/roles/11/settings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "upsert-setting" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["value"] == nil {
			t.Fatalf("expected value payload")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":11}}`))
	})

	data, err := client.Roles.UpsertSetting(
		context.Background(),
		RoleID(11),
		map[string]any{"value": 1},
		WithRolesRequestOptions(pipedrive.WithHeader("X-Test", "upsert-setting")),
	)
	if err != nil {
		t.Fatalf("UpsertSetting error: %v", err)
	}
	if data["id"] != float64(11) {
		t.Fatalf("unexpected data: %#v", data)
	}
}
