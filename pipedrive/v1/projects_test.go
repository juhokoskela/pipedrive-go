package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/juhokoskela/pipedrive-go/v2/pipedrive"
)

func TestProjectsService_ListActivities(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projects/10/activities" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "activities" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"subject":"Call"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	activities, page, err := client.Projects.ListActivities(
		context.Background(),
		ProjectID(10),
		WithProjectsRequestOptions(pipedrive.WithHeader("X-Test", "activities")),
	)
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
		_, _ = w.Write([]byte(`{"success":true,"data":[{"item_id":12,"item_type":"task","phase_id":2,"group_id":3}]}`))
	})

	plan, err := client.Projects.GetPlan(context.Background(), ProjectID(11))
	if err != nil {
		t.Fatalf("GetPlan error: %v", err)
	}
	if len(plan) != 1 || plan[0]["item_id"] != float64(12) || plan[0]["item_type"] != "task" {
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
		if got := r.Header.Get("X-Test"); got != "plan-activity" {
			t.Fatalf("unexpected header X-Test: %q", got)
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

	plan, err := client.Projects.UpdatePlanActivity(
		context.Background(),
		ProjectID(12),
		ProjectPlanActivityID(3),
		map[string]any{"status": "done"},
		WithProjectsRequestOptions(pipedrive.WithHeader("X-Test", "plan-activity")),
	)
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
		if got := r.Header.Get("X-Test"); got != "plan-task" {
			t.Fatalf("unexpected header X-Test: %q", got)
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

	plan, err := client.Projects.UpdatePlanTask(
		context.Background(),
		ProjectID(12),
		ProjectPlanTaskID(4),
		map[string]any{"title": "Updated"},
		WithProjectsRequestOptions(pipedrive.WithHeader("X-Test", "plan-task")),
	)
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

func TestProjectsService_ActivityOutcome(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, outcome string
		want          bool
	}{
		{name: "omitted"},
		{name: "null", outcome: `,"outcome":null`},
		{name: "value", outcome: `,"outcome":7`, want: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[{"id":1` + tt.outcome + `}]}`))
			})
			activities, _, err := client.Projects.ListActivities(t.Context(), 1)
			if err != nil {
				t.Fatal(err)
			}
			if len(activities) != 1 {
				t.Fatalf("activities = %v", activities)
			}
			if tt.want {
				if activities[0].OutcomeID == nil || *activities[0].OutcomeID != 7 {
					t.Errorf("outcome ID = %v", activities[0].OutcomeID)
				}
			} else if activities[0].OutcomeID != nil {
				t.Errorf("outcome ID = %v, want nil", activities[0].OutcomeID)
			}
		})
	}
}
