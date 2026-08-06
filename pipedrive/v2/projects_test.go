package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestProjectsService_AllOperations(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /projects":
			if got := r.URL.Query().Get("status"); got != "open,completed" {
				t.Fatalf("unexpected status: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1,"title":"Migration"}],"additional_data":{"next_cursor":"next"}}`))
		case "GET /projects/archived":
			_, _ = w.Write([]byte(`{"data":[{"id":2,"title":"Archived"}],"additional_data":{"next_cursor":null}}`))
		case "GET /projects/search":
			if got := r.URL.Query().Get("term"); got != "migration" {
				t.Fatalf("unexpected term: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":{"items":[{"result_score":0.9,"item":{"id":1,"title":"Migration"}}]},"additional_data":{"next_cursor":null}}`))
		case "POST /projects":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["title"] != "Migration" || body["board_id"] != float64(3) {
				t.Fatalf("unexpected body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"id":1,"title":"Migration"}}`))
		case "GET /projects/1":
			_, _ = w.Write([]byte(`{"data":{"id":1,"title":"Migration"}}`))
		case "PATCH /projects/1":
			_, _ = w.Write([]byte(`{"data":{"id":1,"title":"Updated"}}`))
		case "DELETE /projects/1":
			_, _ = w.Write([]byte(`{"data":{"id":1}}`))
		case "POST /projects/1/archive":
			_, _ = w.Write([]byte(`{"data":{"id":1,"archive_time":"2026-08-06T10:00:00Z"}}`))
		case "GET /projects/1/changelog":
			_, _ = w.Write([]byte(`{"data":[{"time":"2026-08-06T10:00:00Z","actor_user_id":7,"new_values":{"title":"Updated"},"old_values":{"title":"Migration"}}],"additional_data":{"next_cursor":null}}`))
		case "GET /projects/1/permittedUsers":
			_, _ = w.Write([]byte(`{"data":[7,8]}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	ctx := context.Background()
	projects, next, err := client.Projects.List(ctx, WithProjectsStatus(ProjectStatusOpen, ProjectStatusCompleted))
	if err != nil || len(projects) != 1 || projects[0].ID != 1 || next == nil || *next != "next" {
		t.Fatalf("List = %#v, %v, %v", projects, next, err)
	}
	archived, _, err := client.Projects.ListArchived(ctx)
	if err != nil || len(archived) != 1 || archived[0].ID != 2 {
		t.Fatalf("ListArchived = %#v, %v", archived, err)
	}
	results, _, err := client.Projects.Search(ctx, "migration", WithProjectSearchFields(ProjectSearchFieldTitle))
	if err != nil || len(results) != 1 || results[0].Item.ID != 1 {
		t.Fatalf("Search = %#v, %v", results, err)
	}
	created, err := client.Projects.Create(ctx, WithProjectTitle("Migration"), WithProjectBoardID(3))
	if err != nil || created.ID != 1 {
		t.Fatalf("Create = %#v, %v", created, err)
	}
	got, err := client.Projects.Get(ctx, 1)
	if err != nil || got.ID != 1 {
		t.Fatalf("Get = %#v, %v", got, err)
	}
	updated, err := client.Projects.Update(ctx, 1, WithProjectTitle("Updated"))
	if err != nil || updated.Title != "Updated" {
		t.Fatalf("Update = %#v, %v", updated, err)
	}
	deleted, err := client.Projects.Delete(ctx, 1)
	if err != nil || deleted.ID != 1 {
		t.Fatalf("Delete = %#v, %v", deleted, err)
	}
	archivedProject, err := client.Projects.Archive(ctx, 1)
	if err != nil || archivedProject.ArchiveTime == nil {
		t.Fatalf("Archive = %#v, %v", archivedProject, err)
	}
	changes, _, err := client.Projects.ListChangelog(ctx, 1)
	if err != nil || len(changes) != 1 || changes[0].ActorUserID != 7 {
		t.Fatalf("ListChangelog = %#v, %v", changes, err)
	}
	users, err := client.Projects.ListPermittedUsers(ctx, 1, WithProjectRequestOptions(pipedrive.WithHeader("X-Test", "users")))
	if err != nil || len(users) != 2 || users[1] != 8 {
		t.Fatalf("ListPermittedUsers = %#v, %v", users, err)
	}
}

func TestProjectsService_MissingData(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})
	if _, err := client.Projects.Get(context.Background(), 1); err == nil {
		t.Fatal("expected missing data error")
	}
}
