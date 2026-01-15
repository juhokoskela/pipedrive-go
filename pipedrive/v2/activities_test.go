package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestActivitiesService_List(t *testing.T) {
	t.Parallel()

	since := time.Date(2024, 5, 6, 8, 0, 0, 0, time.UTC)
	until := time.Date(2024, 5, 7, 9, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activities" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("filter_id"); got != "2" {
			t.Fatalf("unexpected filter_id: %q", got)
		}
		if got := q.Get("owner_id"); got != "5" {
			t.Fatalf("unexpected owner_id: %q", got)
		}
		if got := q.Get("deal_id"); got != "9" {
			t.Fatalf("unexpected deal_id: %q", got)
		}
		if got := q.Get("lead_id"); got != "lead-123" {
			t.Fatalf("unexpected lead_id: %q", got)
		}
		if got := q.Get("person_id"); got != "7" {
			t.Fatalf("unexpected person_id: %q", got)
		}
		if got := q.Get("org_id"); got != "8" {
			t.Fatalf("unexpected org_id: %q", got)
		}
		if got := q.Get("done"); got != "true" {
			t.Fatalf("unexpected done: %q", got)
		}
		if got := q.Get("updated_since"); got != since.Format(time.RFC3339) {
			t.Fatalf("unexpected updated_since: %q", got)
		}
		if got := q.Get("updated_until"); got != until.Format(time.RFC3339) {
			t.Fatalf("unexpected updated_until: %q", got)
		}
		if got := q.Get("sort_by"); got != "due_date" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "desc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		if got := q.Get("include_fields"); got != "attendees" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("ids"); got != "1,2" {
			t.Fatalf("unexpected ids: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c3" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1,"subject":"Call"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	activities, next, err := client.Activities.List(
		context.Background(),
		WithActivitiesFilterID(2),
		WithActivitiesOwnerID(5),
		WithActivitiesDealID(DealID(9)),
		WithActivitiesLeadID(LeadID("lead-123")),
		WithActivitiesPersonID(PersonID(7)),
		WithActivitiesOrgID(OrganizationID(8)),
		WithActivitiesDone(true),
		WithActivitiesUpdatedSince(since),
		WithActivitiesUpdatedUntil(until),
		WithActivitiesSortBy(ActivitySortByDueDate),
		WithActivitiesSortDirection(SortDesc),
		WithActivitiesIncludeFields(ActivityIncludeFieldAttendees),
		WithActivitiesIDs(ActivityID(1), ActivityID(2)),
		WithActivitiesPageSize(2),
		WithActivitiesCursor("c3"),
		WithActivityRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(activities) != 1 || activities[0].ID != 1 {
		t.Fatalf("unexpected activities: %#v", activities)
	}
}

func TestActivitiesService_Create(t *testing.T) {
	t.Parallel()

	personID := PersonID(7)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activities" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "2" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["subject"] != "Call" {
			t.Fatalf("unexpected subject: %v", payload["subject"])
		}
		if payload["type"] != "call" {
			t.Fatalf("unexpected type: %v", payload["type"])
		}
		if payload["owner_id"] != float64(5) {
			t.Fatalf("unexpected owner_id: %v", payload["owner_id"])
		}
		if payload["lead_id"] != "lead-123" {
			t.Fatalf("unexpected lead_id: %v", payload["lead_id"])
		}
		location, ok := payload["location"].(map[string]interface{})
		if !ok || location["value"] != "HQ" {
			t.Fatalf("unexpected location: %#v", payload["location"])
		}
		participants, ok := payload["participants"].([]interface{})
		if !ok || len(participants) != 1 {
			t.Fatalf("unexpected participants: %#v", payload["participants"])
		}
		participant, ok := participants[0].(map[string]interface{})
		if !ok || participant["person_id"] != float64(7) {
			t.Fatalf("unexpected participant: %#v", participants[0])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":1,"subject":"Call"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	activity, err := client.Activities.Create(
		context.Background(),
		WithActivitySubject("Call"),
		WithActivityType("call"),
		WithActivityOwnerID(UserID(5)),
		WithActivityLeadID(LeadID("lead-123")),
		WithActivityLocation(ActivityLocation{Value: "HQ"}),
		WithActivityParticipants(ActivityParticipant{PersonID: &personID, Primary: true}),
		WithActivityRequestOptions(pipedrive.WithHeader("X-Test", "2")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if activity.ID != 1 || activity.Subject != "Call" {
		t.Fatalf("unexpected activity: %#v", activity)
	}
}

func TestActivitiesService_Update(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activities/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["done"] != true {
			t.Fatalf("unexpected done: %v", payload["done"])
		}
		if payload["duration"] != "01:00:00" {
			t.Fatalf("unexpected duration: %v", payload["duration"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":5,"subject":"Updated"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	activity, err := client.Activities.Update(
		context.Background(),
		ActivityID(5),
		WithActivityDone(true),
		WithActivityDuration("01:00:00"),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if activity.ID != 5 || activity.Subject != "Updated" {
		t.Fatalf("unexpected activity: %#v", activity)
	}
}

func TestActivitiesService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activities/4" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":4}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Activities.Delete(context.Background(), ActivityID(4))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if result.ID != 4 {
		t.Fatalf("unexpected delete result: %#v", result)
	}
}
