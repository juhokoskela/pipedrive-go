package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestTasksService_List(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("cursor"); got != "c0" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("done"); got != "1" {
			t.Fatalf("unexpected done: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "list" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"title":"Task","assignee_id":7,"done":1,"add_time":"2024-01-01 10:00:00","marked_as_done_time":"2024-01-02 10:00:00"}],"additional_data":{"next_cursor":"c1"}}`))
	})

	query := url.Values{}
	query.Set("cursor", "c0")
	query.Set("limit", "2")
	query.Set("done", "1")
	tasks, page, err := client.Tasks.List(
		context.Background(),
		WithTasksQuery(query),
		WithTasksRequestOptions(pipedrive.WithHeader("X-Test", "list")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(tasks) != 1 || tasks[0].ID != 1 || !tasks[0].Done {
		t.Fatalf("unexpected tasks: %#v", tasks)
	}
	if tasks[0].AddTime == nil || tasks[0].MarkedAsDoneTime == nil {
		t.Fatalf("expected timestamps")
	}
	if page == nil || page.NextCursor == nil || *page.NextCursor != "c1" {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestTasksService_Get(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/tasks/3" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "get" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":3,"title":"Task"}}`))
	})

	task, err := client.Tasks.Get(
		context.Background(),
		TaskID(3),
		WithTasksRequestOptions(pipedrive.WithHeader("X-Test", "get")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if task.ID != 3 {
		t.Fatalf("unexpected task: %#v", task)
	}
}

func TestTasksService_Create(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "create" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["title"] != "Task" {
			t.Fatalf("unexpected title: %#v", body["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":4,"title":"Task"}}`))
	})

	task, err := client.Tasks.Create(
		context.Background(),
		map[string]any{"title": "Task"},
		WithTasksRequestOptions(pipedrive.WithHeader("X-Test", "create")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if task.ID != 4 {
		t.Fatalf("unexpected task: %#v", task)
	}
}

func TestTasksService_Update(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/tasks/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "update" {
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
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5,"title":"Updated"}}`))
	})

	task, err := client.Tasks.Update(
		context.Background(),
		TaskID(5),
		map[string]any{"title": "Updated"},
		WithTasksRequestOptions(pipedrive.WithHeader("X-Test", "update")),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if task.ID != 5 {
		t.Fatalf("unexpected task: %#v", task)
	}
}

func TestTasksService_Delete(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/tasks/6" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	ok, err := client.Tasks.Delete(
		context.Background(),
		TaskID(6),
		WithTasksRequestOptions(pipedrive.WithHeader("X-Test", "delete")),
	)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok")
	}
}
