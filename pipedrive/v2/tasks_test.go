package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestTasksService_AllOperations(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Test"); got != "tasks" {
			t.Fatalf("unexpected request header: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /tasks":
			query := r.URL.Query()
			if query.Get("parent_task_id") == "null" {
				_, _ = w.Write([]byte(`{"data":[],"additional_data":{"next_cursor":null}}`))
				return
			}
			if query.Get("project_id") != "4" || query.Get("is_done") != "false" ||
				query.Get("is_milestone") != "true" || query.Get("assignee_id") != "2" ||
				query.Get("parent_task_id") != "5" || query.Get("limit") != "25" ||
				query.Get("cursor") != "tasks-start" {
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":9,"title":"Ship","project_id":4,"assignee_ids":[2,3]}],"additional_data":{"next_cursor":null}}`))
		case "POST /tasks":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["title"] != "Ship" || body["project_id"] != float64(4) ||
				body["parent_task_id"] != float64(5) || body["description"] != "Release" ||
				body["done"] != float64(0) || body["milestone"] != float64(1) ||
				body["due_date"] != "2026-08-15" || body["start_date"] != "2026-08-10" ||
				body["assignee_id"] != float64(2) || body["priority"] != float64(3) ||
				len(body["assignee_ids"].([]any)) != 2 {
				t.Fatalf("unexpected body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"id":9,"title":"Ship","project_id":4}}`))
		case "GET /tasks/9":
			_, _ = w.Write([]byte(`{"data":{"id":9,"title":"Ship","project_id":4}}`))
		case "PATCH /tasks/9":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			for _, key := range []string{"parent_task_id", "description", "due_date", "start_date", "assignee_id", "priority"} {
				if value, ok := body[key]; !ok || value != nil {
					t.Fatalf("expected %s to be null in %#v", key, body)
				}
			}
			if body["done"] != float64(1) || body["milestone"] != float64(0) || len(body["assignee_ids"].([]any)) != 0 {
				t.Fatalf("unexpected update body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"id":9,"title":"Shipped","project_id":4,"is_done":true}}`))
		case "DELETE /tasks/9":
			_, _ = w.Write([]byte(`{"data":{"id":9}}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	ctx := context.Background()
	requestOpt := WithTaskRequestOptions(pipedrive.WithHeader("X-Test", "tasks"))
	tasks, _, err := client.Tasks.List(ctx,
		WithTasksPageSize(25),
		WithTasksCursor("tasks-start"),
		WithTasksDone(false),
		WithTasksMilestone(true),
		WithTasksAssigneeID(2),
		WithTasksProjectID(4),
		WithTasksParentTaskID(5),
		requestOpt,
	)
	if err != nil || len(tasks) != 1 || len(tasks[0].AssigneeIDs) != 2 {
		t.Fatalf("List = %#v, %v", tasks, err)
	}
	rootTasks, _, err := client.Tasks.List(ctx, WithTasksRootOnly(), requestOpt)
	if err != nil || len(rootTasks) != 0 {
		t.Fatalf("root List = %#v, %v", rootTasks, err)
	}
	created, err := client.Tasks.Create(ctx,
		WithTaskTitle("Ship"),
		WithTaskProjectID(4),
		WithTaskParentTaskID(5),
		WithTaskDescription("Release"),
		WithTaskDone(false),
		WithTaskMilestone(true),
		WithTaskDueDate("2026-08-15"),
		WithTaskStartDate("2026-08-10"),
		WithTaskAssigneeID(2),
		WithTaskAssigneeIDs(2, 3),
		WithTaskPriority(3),
		requestOpt,
	)
	if err != nil || created.ID != 9 {
		t.Fatalf("Create = %#v, %v", created, err)
	}
	got, err := client.Tasks.Get(ctx, 9, requestOpt)
	if err != nil || got.ID != 9 {
		t.Fatalf("Get = %#v, %v", got, err)
	}
	updated, err := client.Tasks.Update(ctx, 9,
		WithTaskTitle("Shipped"),
		ClearTaskParentTaskID(),
		ClearTaskDescription(),
		WithTaskDone(true),
		WithTaskMilestone(false),
		ClearTaskDueDate(),
		ClearTaskStartDate(),
		ClearTaskAssigneeID(),
		WithTaskAssigneeIDs(),
		ClearTaskPriority(),
		requestOpt,
	)
	if err != nil || !updated.IsDone {
		t.Fatalf("Update = %#v, %v", updated, err)
	}
	deleted, err := client.Tasks.Delete(ctx, 9, requestOpt)
	if err != nil || deleted.ID != 9 {
		t.Fatalf("Delete = %#v, %v", deleted, err)
	}
}

func TestTasksService_PagerAndIterator(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/tasks" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("cursor") {
		case "start":
			_, _ = w.Write([]byte(`{"data":[{"id":1}],"additional_data":{"next_cursor":"next"}}`))
		case "next":
			_, _ = w.Write([]byte(`{"data":[{"id":2}],"additional_data":{"next_cursor":null}}`))
		default:
			_, _ = w.Write([]byte(`{"data":[{"id":3}],"additional_data":{"next_cursor":null}}`))
		}
	})

	ctx := context.Background()
	pager := client.Tasks.ListPager(WithTasksPageSize(2), WithTasksCursor("start"))
	var ids []TaskID
	for pager.Next(ctx) {
		for _, task := range pager.Items() {
			ids = append(ids, task.ID)
		}
	}
	if err := pager.Err(); err != nil || len(ids) != 2 || ids[1] != 2 {
		t.Fatalf("pager = %v, %v", ids, err)
	}

	var iterated []TaskID
	if err := client.Tasks.ForEach(ctx, func(task Task) error {
		iterated = append(iterated, task.ID)
		return nil
	}); err != nil {
		t.Fatalf("ForEach error: %v", err)
	}
	if len(iterated) != 1 || iterated[0] != 3 {
		t.Fatalf("unexpected iterated tasks: %v", iterated)
	}
}
