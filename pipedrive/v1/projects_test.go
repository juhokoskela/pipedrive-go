package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func TestProjectsService_List(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"title":"Project"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	query := url.Values{}
	query.Set("limit", "1")
	projects, page, err := client.Projects.List(context.Background(), WithProjectsQuery(query))
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(projects) != 1 || projects[0].ID != 1 {
		t.Fatalf("unexpected projects: %#v", projects)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestProjectsService_Create(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["title"] != "New" {
			t.Fatalf("unexpected title: %#v", body["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":2,"title":"New"}}`))
	})

	project, err := client.Projects.Create(context.Background(), map[string]any{"title": "New"})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if project.ID != 2 || project.Title != "New" {
		t.Fatalf("unexpected project: %#v", project)
	}
}

func TestProjectsService_Get(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/3" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":3,"title":"Project"}}`))
	})

	project, err := client.Projects.Get(context.Background(), ProjectID(3))
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if project.ID != 3 {
		t.Fatalf("unexpected project: %#v", project)
	}
}

func TestProjectsService_Update(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/4" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["title"] != "Updated" {
			t.Fatalf("unexpected title: %#v", body["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":4,"title":"Updated"}}`))
	})

	project, err := client.Projects.Update(context.Background(), ProjectID(4), map[string]any{"title": "Updated"})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if project.ID != 4 || project.Title != "Updated" {
		t.Fatalf("unexpected project: %#v", project)
	}
}

func TestProjectsService_Delete(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5}}`))
	})

	id, err := client.Projects.Delete(context.Background(), ProjectID(5))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if id != 5 {
		t.Fatalf("unexpected id: %d", id)
	}
}

func TestProjectsService_Archive(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/6/archive" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":6,"status":"archived"}}`))
	})

	project, err := client.Projects.Archive(context.Background(), ProjectID(6))
	if err != nil {
		t.Fatalf("Archive error: %v", err)
	}
	if project.ID != 6 || project.Status != "archived" {
		t.Fatalf("unexpected project: %#v", project)
	}
}

func TestProjectsService_ListBoards(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/boards" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Board"}]}`))
	})

	boards, err := client.Projects.ListBoards(context.Background())
	if err != nil {
		t.Fatalf("ListBoards error: %v", err)
	}
	if len(boards) != 1 || boards[0].ID != 1 {
		t.Fatalf("unexpected boards: %#v", boards)
	}
}

func TestProjectsService_GetBoard(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/boards/2" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":2,"name":"Board"}}`))
	})

	board, err := client.Projects.GetBoard(context.Background(), ProjectBoardID(2))
	if err != nil {
		t.Fatalf("GetBoard error: %v", err)
	}
	if board.ID != 2 {
		t.Fatalf("unexpected board: %#v", board)
	}
}

func TestProjectsService_ListPhases(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/phases" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":3,"name":"Phase"}]}`))
	})

	phases, err := client.Projects.ListPhases(context.Background())
	if err != nil {
		t.Fatalf("ListPhases error: %v", err)
	}
	if len(phases) != 1 || phases[0].ID != 3 {
		t.Fatalf("unexpected phases: %#v", phases)
	}
}

func TestProjectsService_GetPhase(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/phases/4" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":4,"name":"Phase"}}`))
	})

	phase, err := client.Projects.GetPhase(context.Background(), ProjectPhaseID(4))
	if err != nil {
		t.Fatalf("GetPhase error: %v", err)
	}
	if phase.ID != 4 {
		t.Fatalf("unexpected phase: %#v", phase)
	}
}

func TestProjectsService_ListActivities(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/10/activities" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"subject":"Call"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	activities, page, err := client.Projects.ListActivities(context.Background(), ProjectID(10))
	if err != nil {
		t.Fatalf("ListActivities error: %v", err)
	}
	if len(activities) != 1 || activities[0].ID != 1 {
		t.Fatalf("unexpected activities: %#v", activities)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestProjectsService_ListGroups(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/10/groups" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":2,"name":"Group"}]}`))
	})

	groups, err := client.Projects.ListGroups(context.Background(), ProjectID(10))
	if err != nil {
		t.Fatalf("ListGroups error: %v", err)
	}
	if len(groups) != 1 || groups[0].ID != 2 {
		t.Fatalf("unexpected groups: %#v", groups)
	}
}

func TestProjectsService_GetPlan(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/11/plan" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"status":"ok"}}`))
	})

	plan, err := client.Projects.GetPlan(context.Background(), ProjectID(11))
	if err != nil {
		t.Fatalf("GetPlan error: %v", err)
	}
	if plan["status"] != "ok" {
		t.Fatalf("unexpected plan: %#v", plan)
	}
}

func TestProjectsService_UpdatePlanActivity(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/12/plan/activities/3" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["status"] != "done" {
			t.Fatalf("unexpected status: %#v", body["status"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"status":"done"}}`))
	})

	plan, err := client.Projects.UpdatePlanActivity(context.Background(), ProjectID(12), ProjectPlanActivityID(3), map[string]any{"status": "done"})
	if err != nil {
		t.Fatalf("UpdatePlanActivity error: %v", err)
	}
	if plan["status"] != "done" {
		t.Fatalf("unexpected plan: %#v", plan)
	}
}

func TestProjectsService_UpdatePlanTask(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/12/plan/tasks/4" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["title"] != "Updated" {
			t.Fatalf("unexpected title: %#v", body["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"title":"Updated"}}`))
	})

	plan, err := client.Projects.UpdatePlanTask(context.Background(), ProjectID(12), ProjectPlanTaskID(4), map[string]any{"title": "Updated"})
	if err != nil {
		t.Fatalf("UpdatePlanTask error: %v", err)
	}
	if plan["title"] != "Updated" {
		t.Fatalf("unexpected plan: %#v", plan)
	}
}

func TestProjectsService_ListTasks(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/13/tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"title":"Task"}]}`))
	})

	tasks, err := client.Projects.ListTasks(context.Background(), ProjectID(13))
	if err != nil {
		t.Fatalf("ListTasks error: %v", err)
	}
	if len(tasks) != 1 || tasks[0].ID != 1 {
		t.Fatalf("unexpected tasks: %#v", tasks)
	}
}
