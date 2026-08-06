package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTasksService_AllOperations(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /tasks":
			if r.URL.Query().Get("project_id") != "4" || r.URL.Query().Get("is_done") != "false" {
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":9,"title":"Ship","project_id":4,"assignee_ids":[2,3]}],"additional_data":{"next_cursor":null}}`))
		case "POST /tasks":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["title"] != "Ship" || body["project_id"] != float64(4) {
				t.Fatalf("unexpected body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"id":9,"title":"Ship","project_id":4}}`))
		case "GET /tasks/9":
			_, _ = w.Write([]byte(`{"data":{"id":9,"title":"Ship","project_id":4}}`))
		case "PATCH /tasks/9":
			_, _ = w.Write([]byte(`{"data":{"id":9,"title":"Shipped","project_id":4,"is_done":true}}`))
		case "DELETE /tasks/9":
			_, _ = w.Write([]byte(`{"data":{"id":9}}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	ctx := context.Background()
	tasks, _, err := client.Tasks.List(ctx, WithTasksProjectID(4), WithTasksDone(false))
	if err != nil || len(tasks) != 1 || len(tasks[0].AssigneeIDs) != 2 {
		t.Fatalf("List = %#v, %v", tasks, err)
	}
	created, err := client.Tasks.Create(ctx, WithTaskTitle("Ship"), WithTaskProjectID(4), WithTaskAssigneeIDs(2, 3))
	if err != nil || created.ID != 9 {
		t.Fatalf("Create = %#v, %v", created, err)
	}
	got, err := client.Tasks.Get(ctx, 9)
	if err != nil || got.ID != 9 {
		t.Fatalf("Get = %#v, %v", got, err)
	}
	updated, err := client.Tasks.Update(ctx, 9, WithTaskTitle("Shipped"), WithTaskDone(true), ClearTaskDescription())
	if err != nil || !updated.IsDone {
		t.Fatalf("Update = %#v, %v", updated, err)
	}
	deleted, err := client.Tasks.Delete(ctx, 9)
	if err != nil || deleted.ID != 9 {
		t.Fatalf("Delete = %#v, %v", deleted, err)
	}
}
