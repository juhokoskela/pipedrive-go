package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestPipelinesService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/pipelines" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("sort_by"); got != "update_time" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "desc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1,"name":"Pipeline"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	pipelines, next, err := client.Pipelines.List(
		context.Background(),
		WithPipelinesSortBy(PipelineSortByUpdateTime),
		WithPipelinesSortDirection(SortDesc),
		WithPipelinesPageSize(2),
		WithPipelinesCursor("c1"),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(pipelines) != 1 || pipelines[0].ID != 1 {
		t.Fatalf("unexpected pipelines: %#v", pipelines)
	}
}

func TestPipelinesService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/pipelines" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Sales" {
			t.Fatalf("unexpected name: %v", payload["name"])
		}
		if payload["is_deal_probability_enabled"] != true {
			t.Fatalf("unexpected is_deal_probability_enabled: %v", payload["is_deal_probability_enabled"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":5,"name":"Sales"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	pipeline, err := client.Pipelines.Create(
		context.Background(),
		WithPipelineName("Sales"),
		WithPipelineDealProbabilityEnabled(true),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if pipeline.ID != 5 || pipeline.Name != "Sales" {
		t.Fatalf("unexpected pipeline: %#v", pipeline)
	}
}
