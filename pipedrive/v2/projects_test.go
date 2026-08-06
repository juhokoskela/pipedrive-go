package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestProjectsService_AllOperations(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Test"); got != "projects" {
			t.Fatalf("unexpected request header: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /projects":
			query := r.URL.Query()
			if query.Get("filter_id") != "5" || query.Get("status") != "open,completed" ||
				query.Get("phase_id") != "3" || query.Get("deal_id") != "4" ||
				query.Get("person_id") != "5" || query.Get("org_id") != "6" ||
				query.Get("limit") != "25" || query.Get("cursor") != "projects-start" {
				t.Fatalf("unexpected projects query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1,"title":"Migration"}],"additional_data":{"next_cursor":"next"}}`))
		case "GET /projects/archived":
			query := r.URL.Query()
			if query.Get("filter_id") != "7" || query.Get("status") != "canceled" ||
				query.Get("phase_id") != "8" || query.Get("limit") != "10" ||
				query.Get("cursor") != "archived-start" {
				t.Fatalf("unexpected archived projects query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":2,"title":"Archived"}],"additional_data":{"next_cursor":null}}`))
		case "GET /projects/search":
			query := r.URL.Query()
			if query.Get("term") != "migration" || query.Get("fields") != "title,description" ||
				query.Get("exact_match") != "true" || query.Get("person_id") != "5" ||
				query.Get("organization_id") != "6" || query.Get("limit") != "15" ||
				query.Get("cursor") != "search-start" {
				t.Fatalf("unexpected search query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":{"items":[{"result_score":0.9,"item":{"id":1,"title":"Migration"}}]},"additional_data":{"next_cursor":null}}`))
		case "POST /projects":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["title"] != "Migration" || body["description"] != "Move systems" ||
				body["status"] != "open" || body["board_id"] != float64(3) ||
				body["phase_id"] != float64(4) || body["owner_id"] != float64(7) ||
				body["start_date"] != "2026-08-01" || body["end_date"] != "2026-08-31" ||
				body["health_status"] != float64(2) || body["template_id"] != float64(9) {
				t.Fatalf("unexpected body: %#v", body)
			}
			if len(body["deal_ids"].([]any)) != 2 || len(body["person_ids"].([]any)) != 1 ||
				len(body["org_ids"].([]any)) != 1 || len(body["label_ids"].([]any)) != 2 ||
				body["custom_fields"].(map[string]any)["risk"] != "high" {
				t.Fatalf("unexpected collection fields: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"id":1,"title":"Migration"}}`))
		case "GET /projects/1":
			_, _ = w.Write([]byte(`{"data":{"id":1,"title":"Migration"}}`))
		case "PATCH /projects/1":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["title"] != "Updated" || len(body["deal_ids"].([]any)) != 0 ||
				len(body["person_ids"].([]any)) != 0 || len(body["org_ids"].([]any)) != 0 ||
				len(body["label_ids"].([]any)) != 0 {
				t.Fatalf("unexpected update body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"id":1,"title":"Updated"}}`))
		case "DELETE /projects/1":
			_, _ = w.Write([]byte(`{"data":{"id":1}}`))
		case "POST /projects/1/archive":
			_, _ = w.Write([]byte(`{"data":{"id":1,"archive_time":"2026-08-06T10:00:00Z"}}`))
		case "GET /projects/1/changelog":
			if query := r.URL.Query(); query.Get("limit") != "20" || query.Get("cursor") != "changes-start" {
				t.Fatalf("unexpected changelog query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"time":"2026-08-06T10:00:00Z","actor_user_id":7,"new_values":{"title":"Updated"},"old_values":{"title":"Migration"}}],"additional_data":{"next_cursor":null}}`))
		case "GET /projects/1/permittedUsers":
			_, _ = w.Write([]byte(`{"data":[7,8]}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	ctx := context.Background()
	requestOpt := WithProjectRequestOptions(pipedrive.WithHeader("X-Test", "projects"))
	projects, next, err := client.Projects.List(ctx,
		WithProjectFilterID(5),
		WithProjectsStatus(ProjectStatusOpen, ProjectStatusCompleted),
		WithProjectsPhaseID(3),
		WithProjectsDealID(4),
		WithProjectsPersonID(5),
		WithProjectsOrganizationID(6),
		WithProjectsPageSize(25),
		WithProjectsCursor("projects-start"),
		requestOpt,
	)
	if err != nil || len(projects) != 1 || projects[0].ID != 1 || next == nil || *next != "next" {
		t.Fatalf("List = %#v, %v, %v", projects, next, err)
	}
	archived, _, err := client.Projects.ListArchived(ctx,
		WithProjectFilterID(7),
		WithProjectsStatus(ProjectStatusCanceled),
		WithProjectsPhaseID(8),
		WithProjectsPageSize(10),
		WithProjectsCursor("archived-start"),
		requestOpt,
	)
	if err != nil || len(archived) != 1 || archived[0].ID != 2 {
		t.Fatalf("ListArchived = %#v, %v", archived, err)
	}
	results, _, err := client.Projects.Search(ctx, "migration",
		WithProjectSearchFields(ProjectSearchFieldTitle, ProjectSearchFieldDescription),
		WithProjectSearchExactMatch(true),
		WithProjectSearchPersonID(5),
		WithProjectSearchOrganizationID(6),
		WithProjectSearchPageSize(15),
		WithProjectSearchCursor("search-start"),
		requestOpt,
	)
	if err != nil || len(results) != 1 || results[0].Item.ID != 1 {
		t.Fatalf("Search = %#v, %v", results, err)
	}
	created, err := client.Projects.Create(ctx,
		WithProjectTitle("Migration"),
		WithProjectDescription("Move systems"),
		WithProjectStatus(ProjectStatusOpen),
		WithProjectBoardID(3),
		WithProjectPhaseID(4),
		WithProjectOwnerID(7),
		WithProjectStartDate("2026-08-01"),
		WithProjectEndDate("2026-08-31"),
		WithProjectHealthStatus(2),
		WithProjectTemplateID(9),
		WithProjectDealIDs(10, 11),
		WithProjectPersonIDs(12),
		WithProjectOrganizationIDs(13),
		WithProjectLabelIDs(14, 15),
		WithProjectCustomFields(map[string]interface{}{"risk": "high"}),
		requestOpt,
	)
	if err != nil || created.ID != 1 {
		t.Fatalf("Create = %#v, %v", created, err)
	}
	got, err := client.Projects.Get(ctx, 1, requestOpt)
	if err != nil || got.ID != 1 {
		t.Fatalf("Get = %#v, %v", got, err)
	}
	updated, err := client.Projects.Update(ctx, 1,
		WithProjectTitle("Updated"),
		WithProjectDealIDs(),
		WithProjectPersonIDs(),
		WithProjectOrganizationIDs(),
		WithProjectLabelIDs(),
		requestOpt,
	)
	if err != nil || updated.Title != "Updated" {
		t.Fatalf("Update = %#v, %v", updated, err)
	}
	deleted, err := client.Projects.Delete(ctx, 1, requestOpt)
	if err != nil || deleted.ID != 1 {
		t.Fatalf("Delete = %#v, %v", deleted, err)
	}
	archivedProject, err := client.Projects.Archive(ctx, 1, requestOpt)
	if err != nil || archivedProject.ArchiveTime == nil {
		t.Fatalf("Archive = %#v, %v", archivedProject, err)
	}
	changes, _, err := client.Projects.ListChangelog(ctx, 1,
		WithProjectChangelogPageSize(20),
		WithProjectChangelogCursor("changes-start"),
		requestOpt,
	)
	if err != nil || len(changes) != 1 || changes[0].ActorUserID != 7 {
		t.Fatalf("ListChangelog = %#v, %v", changes, err)
	}
	users, err := client.Projects.ListPermittedUsers(ctx, 1, requestOpt)
	if err != nil || len(users) != 2 || users[1] != 8 {
		t.Fatalf("ListPermittedUsers = %#v, %v", users, err)
	}
}

func TestProjectsService_PagersAndIterators(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cursor := r.URL.Query().Get("cursor")
		id := 1
		next := `null`
		switch cursor {
		case "start":
			next = `"next"`
		case "next":
			id = 2
		}
		switch r.URL.Path {
		case "/projects", "/projects/archived":
			_, _ = w.Write([]byte(`{"data":[{"id":` + strconv.Itoa(id) + `}],"additional_data":{"next_cursor":` + next + `}}`))
		case "/projects/search":
			_, _ = w.Write([]byte(`{"data":{"items":[{"item":{"id":` + strconv.Itoa(id) + `}}]},"additional_data":{"next_cursor":` + next + `}}`))
		case "/projects/1/changelog":
			_, _ = w.Write([]byte(`{"data":[{"actor_user_id":` + strconv.Itoa(id) + `}],"additional_data":{"next_cursor":` + next + `}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})

	ctx := context.Background()
	projectPager := client.Projects.ListPager(WithProjectsPageSize(2), WithProjectsCursor("start"))
	var projectIDs []ProjectID
	for projectPager.Next(ctx) {
		for _, project := range projectPager.Items() {
			projectIDs = append(projectIDs, project.ID)
		}
	}
	if err := projectPager.Err(); err != nil || len(projectIDs) != 2 || projectIDs[1] != 2 {
		t.Fatalf("project pager = %v, %v", projectIDs, err)
	}

	archivedPager := client.Projects.ListArchivedPager(WithProjectsCursor("start"))
	var archivedIDs []ProjectID
	for archivedPager.Next(ctx) {
		for _, project := range archivedPager.Items() {
			archivedIDs = append(archivedIDs, project.ID)
		}
	}
	if err := archivedPager.Err(); err != nil || len(archivedIDs) != 2 {
		t.Fatalf("archived pager = %v, %v", archivedIDs, err)
	}

	searchPager := client.Projects.SearchPager("migration", WithProjectSearchCursor("start"))
	var searchIDs []ProjectID
	for searchPager.Next(ctx) {
		for _, result := range searchPager.Items() {
			searchIDs = append(searchIDs, result.Item.ID)
		}
	}
	if err := searchPager.Err(); err != nil || len(searchIDs) != 2 {
		t.Fatalf("search pager = %v, %v", searchIDs, err)
	}

	changelogPager := client.Projects.ChangelogPager(1, WithProjectChangelogCursor("start"))
	var actors []UserID
	for changelogPager.Next(ctx) {
		for _, change := range changelogPager.Items() {
			actors = append(actors, change.ActorUserID)
		}
	}
	if err := changelogPager.Err(); err != nil || len(actors) != 2 {
		t.Fatalf("changelog pager = %v, %v", actors, err)
	}

	var iterated int
	if err := client.Projects.ForEach(ctx, func(Project) error { iterated++; return nil }); err != nil {
		t.Fatalf("ForEach error: %v", err)
	}
	if err := client.Projects.ForEachArchived(ctx, func(Project) error { iterated++; return nil }); err != nil {
		t.Fatalf("ForEachArchived error: %v", err)
	}
	if err := client.Projects.ForEachSearch(ctx, "migration", func(ProjectSearchResult) error { iterated++; return nil }); err != nil {
		t.Fatalf("ForEachSearch error: %v", err)
	}
	if err := client.Projects.ForEachChangelog(ctx, 1, func(ProjectChangelogEntry) error { iterated++; return nil }); err != nil {
		t.Fatalf("ForEachChangelog error: %v", err)
	}
	if iterated != 4 {
		t.Fatalf("unexpected iteration count: %d", iterated)
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
