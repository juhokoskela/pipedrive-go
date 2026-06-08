package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestStagesService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/stages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("pipeline_id"); got != "3" {
			t.Fatalf("unexpected pipeline_id: %q", got)
		}
		if got := q.Get("sort_by"); got != "order_nr" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "asc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "list" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":10,"name":"Stage"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	stages, next, err := client.Stages.List(
		context.Background(),
		WithStagesPipelineID(3),
		WithStagesSortBy(StageSortByOrder),
		WithStagesSortDirection(SortAsc),
		WithStagesPageSize(1),
		WithStagesCursor("c1"),
		WithStageRequestOptions(pipedrive.WithHeader("X-Test", "list")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(stages) != 1 || stages[0].ID != 10 {
		t.Fatalf("unexpected stages: %#v", stages)
	}
}

func TestStagesService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/stages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "create" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Qualified" {
			t.Fatalf("unexpected name: %v", payload["name"])
		}
		if payload["pipeline_id"] != float64(2) {
			t.Fatalf("unexpected pipeline_id: %v", payload["pipeline_id"])
		}
		if payload["deal_probability"] != float64(25) {
			t.Fatalf("unexpected deal_probability: %v", payload["deal_probability"])
		}
		if payload["is_deal_rot_enabled"] != true {
			t.Fatalf("unexpected is_deal_rot_enabled: %v", payload["is_deal_rot_enabled"])
		}
		if payload["days_to_rotten"] != float64(14) {
			t.Fatalf("unexpected days_to_rotten: %v", payload["days_to_rotten"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":11,"name":"Qualified"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	stage, err := client.Stages.Create(
		context.Background(),
		WithStageName("Qualified"),
		WithStagePipelineID(2),
		WithStageDealProbability(25),
		WithStageDealRotEnabled(true),
		WithStageDaysToRotten(14),
		WithStageRequestOptions(pipedrive.WithHeader("X-Test", "create")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if stage.ID != 11 || stage.Name != "Qualified" {
		t.Fatalf("unexpected stage: %#v", stage)
	}
}

func TestStagesService_Get(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/stages/11" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "get" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":11,"name":"Qualified"}}`))
	})

	stage, err := client.Stages.Get(
		context.Background(),
		StageID(11),
		WithStageRequestOptions(pipedrive.WithHeader("X-Test", "get")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if stage.ID != 11 || stage.Name != "Qualified" {
		t.Fatalf("unexpected stage: %#v", stage)
	}
}

func TestStagesService_ListPager(t *testing.T) {
	t.Parallel()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/stages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		switch calls {
		case 1:
			if got := r.URL.Query().Get("cursor"); got != "start" {
				t.Fatalf("unexpected first cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1}],"additional_data":{"next_cursor":"next"}}`))
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "next" {
				t.Fatalf("unexpected second cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":2}],"additional_data":{"next_cursor":null}}`))
		default:
			t.Fatalf("unexpected call count: %d", calls)
		}
	})

	pager := client.Stages.ListPager(WithStagesPageSize(2), WithStagesCursor("start"))
	var ids []StageID
	for pager.Next(context.Background()) {
		for _, stage := range pager.Items() {
			ids = append(ids, stage.ID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestStagesService_ForEach(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/stages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":null}}`))
	})

	var ids []StageID
	err := client.Stages.ForEach(context.Background(), func(stage Stage) error {
		ids = append(ids, stage.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("ForEach error: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestStagesService_Update(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/stages/11" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "update" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Updated" {
			t.Fatalf("unexpected name: %#v", payload["name"])
		}
		if payload["pipeline_id"] != float64(2) {
			t.Fatalf("unexpected pipeline_id: %#v", payload["pipeline_id"])
		}
		if payload["deal_probability"] != float64(35) {
			t.Fatalf("unexpected deal_probability: %#v", payload["deal_probability"])
		}
		if payload["is_deal_rot_enabled"] != false {
			t.Fatalf("unexpected is_deal_rot_enabled: %#v", payload["is_deal_rot_enabled"])
		}
		if payload["days_to_rotten"] != float64(21) {
			t.Fatalf("unexpected days_to_rotten: %#v", payload["days_to_rotten"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":11,"name":"Updated"}}`))
	})

	stage, err := client.Stages.Update(
		context.Background(),
		StageID(11),
		WithStageName("Updated"),
		WithStagePipelineID(2),
		WithStageDealProbability(35),
		WithStageDealRotEnabled(false),
		WithStageDaysToRotten(21),
		WithStageRequestOptions(pipedrive.WithHeader("X-Test", "update")),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if stage.ID != 11 || stage.Name != "Updated" {
		t.Fatalf("unexpected stage: %#v", stage)
	}
}

func TestStagesService_Delete(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/stages/11" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":11}}`))
	})

	result, err := client.Stages.Delete(
		context.Background(),
		StageID(11),
		WithStageRequestOptions(pipedrive.WithHeader("X-Test", "delete")),
	)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if result.ID != 11 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
