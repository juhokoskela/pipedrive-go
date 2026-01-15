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
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if stage.ID != 11 || stage.Name != "Qualified" {
		t.Fatalf("unexpected stage: %#v", stage)
	}
}
