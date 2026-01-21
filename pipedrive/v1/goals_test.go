package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestGoalsService_List(t *testing.T) {
	t.Parallel()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/goals/find" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("type.name"); got != "deals_won" {
			t.Fatalf("unexpected type.name: %q", got)
		}
		if got := q.Get("assignee.id"); got != "12" {
			t.Fatalf("unexpected assignee.id: %q", got)
		}
		if got := q.Get("assignee.type"); got != "team" {
			t.Fatalf("unexpected assignee.type: %q", got)
		}
		if got := q.Get("is_active"); got != "true" {
			t.Fatalf("unexpected is_active: %q", got)
		}
		if got := q.Get("expected_outcome.target"); got != "42" {
			t.Fatalf("unexpected expected_outcome.target: %q", got)
		}
		if got := q.Get("period.start"); got != "2024-01-01" {
			t.Fatalf("unexpected period.start: %q", got)
		}
		if got := q.Get("period.end"); got != "2024-01-31" {
			t.Fatalf("unexpected period.end: %q", got)
		}
		if got := q["type.params.pipeline_id"]; len(got) != 2 || got[0] != "3" || got[1] != "5" {
			t.Fatalf("unexpected pipeline ids: %#v", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"goals":[{"id":"goal-1","owner_id":1,"title":"Pipeline goals","type":{"name":"deals_won","params":{"pipeline_id":[3,5]}},"assignee":{"id":12,"type":"team"},"interval":"weekly","duration":{"start":"2024-01-01","end":"2024-12-31"},"expected_outcome":{"target":42,"tracking_metric":"quantity"},"is_active":true,"report_ids":["report-1"]}]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	goals, err := client.Goals.List(
		context.Background(),
		WithGoalsTypeName(GoalTypeNameDealsWon),
		WithGoalsAssigneeID(12),
		WithGoalsAssigneeType(GoalAssigneeTypeTeam),
		WithGoalsActive(true),
		WithGoalsExpectedOutcomeTarget(42),
		WithGoalsTypePipelineIDs(PipelineID(3), PipelineID(5)),
		WithGoalsPeriodStart(start),
		WithGoalsPeriodEnd(end),
		WithGoalsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(goals) != 1 || goals[0].ID != GoalID("goal-1") {
		t.Fatalf("unexpected goals: %#v", goals)
	}
	if goals[0].Type == nil || goals[0].Type.Name != GoalTypeNameDealsWon {
		t.Fatalf("unexpected goal type: %#v", goals[0].Type)
	}
}

func TestGoalsService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/goals" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["title"] != "New Goal" {
			t.Fatalf("unexpected title: %#v", payload["title"])
		}

		assignee, ok := payload["assignee"].(map[string]interface{})
		if !ok {
			t.Fatalf("missing assignee payload")
		}
		if assignee["id"] != float64(12) || assignee["type"] != "person" {
			t.Fatalf("unexpected assignee: %#v", assignee)
		}

		goalType, ok := payload["type"].(map[string]interface{})
		if !ok {
			t.Fatalf("missing type payload")
		}
		if goalType["name"] != "deals_started" {
			t.Fatalf("unexpected type name: %#v", goalType["name"])
		}
		params, ok := goalType["params"].(map[string]interface{})
		if !ok {
			t.Fatalf("missing params payload")
		}
		pipelines, ok := params["pipeline_id"].([]interface{})
		if !ok || len(pipelines) != 2 || pipelines[0] != float64(1) || pipelines[1] != float64(2) {
			t.Fatalf("unexpected pipeline ids: %#v", params["pipeline_id"])
		}

		expected, ok := payload["expected_outcome"].(map[string]interface{})
		if !ok {
			t.Fatalf("missing expected_outcome payload")
		}
		if expected["target"] != float64(100) || expected["tracking_metric"] != "quantity" {
			t.Fatalf("unexpected expected_outcome: %#v", expected)
		}

		duration, ok := payload["duration"].(map[string]interface{})
		if !ok {
			t.Fatalf("missing duration payload")
		}
		if duration["start"] != "2024-01-01" || duration["end"] != "2024-12-31" {
			t.Fatalf("unexpected duration: %#v", duration)
		}
		if payload["interval"] != "monthly" {
			t.Fatalf("unexpected interval: %#v", payload["interval"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"goal":{"id":"goal-2","title":"New Goal","interval":"monthly","assignee":{"id":12,"type":"person"}}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	goal, err := client.Goals.Create(
		context.Background(),
		WithGoalTitle("New Goal"),
		WithGoalAssigneeID(12),
		WithGoalAssigneeType(GoalAssigneeTypePerson),
		WithGoalTypeName(GoalTypeNameDealsStarted),
		WithGoalTypePipelineIDs(PipelineID(1), PipelineID(2)),
		WithGoalExpectedOutcomeTarget(100),
		WithGoalExpectedOutcomeTrackingMetric(GoalTrackingMetricQuantity),
		WithGoalDurationStart("2024-01-01"),
		WithGoalDurationEnd("2024-12-31"),
		WithGoalInterval(GoalIntervalMonthly),
		WithGoalsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if goal.ID != GoalID("goal-2") || goal.Interval != GoalIntervalMonthly {
		t.Fatalf("unexpected goal: %#v", goal)
	}
}

func TestGoalsService_Update(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/goals/goal-3" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(payload) != 2 {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if payload["title"] != "Updated Goal" {
			t.Fatalf("unexpected title: %#v", payload["title"])
		}
		if payload["interval"] != "quarterly" {
			t.Fatalf("unexpected interval: %#v", payload["interval"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"goal":{"id":"goal-3","title":"Updated Goal","interval":"quarterly"}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	goal, err := client.Goals.Update(
		context.Background(),
		GoalID("goal-3"),
		WithGoalTitle("Updated Goal"),
		WithGoalInterval(GoalIntervalQuarterly),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if goal.ID != GoalID("goal-3") || goal.Interval != GoalIntervalQuarterly {
		t.Fatalf("unexpected goal: %#v", goal)
	}
}

func TestGoalsService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/goals/goal-4" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ok, err := client.Goals.Delete(context.Background(), GoalID("goal-4"))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}

func TestGoalsService_GetResult(t *testing.T) {
	t.Parallel()

	start := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/goals/goal-5/results" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("period.start"); got != "2024-02-01" {
			t.Fatalf("unexpected period.start: %q", got)
		}
		if got := q.Get("period.end"); got != "2024-02-28" {
			t.Fatalf("unexpected period.end: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"progress":3,"goal":{"id":"goal-5","title":"Result goal"}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Goals.GetResult(
		context.Background(),
		GoalID("goal-5"),
		WithGoalResultStartDate(start),
		WithGoalResultEndDate(end),
	)
	if err != nil {
		t.Fatalf("GetResult error: %v", err)
	}
	if result.Progress != 3 || result.Goal == nil || result.Goal.ID != GoalID("goal-5") {
		t.Fatalf("unexpected result: %#v", result)
	}
}
