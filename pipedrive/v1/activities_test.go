package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestActivitiesService_ListCollection(t *testing.T) {
	t.Parallel()

	since := time.Date(2022, 11, 1, 8, 55, 59, 0, time.UTC)
	until := time.Date(2022, 11, 2, 9, 10, 11, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activities/collection" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("cursor"); got != "cursor-1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := q.Get("limit"); got != "50" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("since"); got != "2022-11-01 08:55:59" {
			t.Fatalf("unexpected since: %q", got)
		}
		if got := q.Get("until"); got != "2022-11-02 09:10:11" {
			t.Fatalf("unexpected until: %q", got)
		}
		if got := q.Get("user_id"); got != "7" {
			t.Fatalf("unexpected user_id: %q", got)
		}
		if got := q.Get("done"); got != "true" {
			t.Fatalf("unexpected done: %q", got)
		}
		if got := q.Get("type"); got != "call,meeting" {
			t.Fatalf("unexpected type: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1,"subject":"Call","type":"call","done":false,"user_id":7,"add_time":"2022-11-01 08:55:59"}],"additional_data":{"next_cursor":"next"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	activities, next, err := client.Activities.ListCollection(
		context.Background(),
		WithActivitiesCollectionCursor("cursor-1"),
		WithActivitiesCollectionLimit(50),
		WithActivitiesCollectionSince(since),
		WithActivitiesCollectionUntil(until),
		WithActivitiesCollectionUserID(UserID(7)),
		WithActivitiesCollectionDone(true),
		WithActivitiesCollectionTypes("call", "meeting"),
		WithActivitiesRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("ListCollection error: %v", err)
	}
	if next == nil || *next != "next" {
		t.Fatalf("unexpected next cursor: %#v", next)
	}
	if len(activities) != 1 || activities[0].ID != 1 || activities[0].Subject != "Call" {
		t.Fatalf("unexpected activities: %#v", activities)
	}
	if activities[0].AddTime == nil || activities[0].AddTime.IsZero() {
		t.Fatalf("expected add_time to be parsed")
	}
}

func TestActivitiesService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activities" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ids"); got != "1,2" {
			t.Fatalf("unexpected ids: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":[1,2]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Activities.Delete(
		context.Background(),
		[]ActivityID{1, 2},
		WithActivitiesRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if len(result.IDs) != 2 || result.IDs[0] != 1 || result.IDs[1] != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
