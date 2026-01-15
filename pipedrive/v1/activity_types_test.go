package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestActivityTypesService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activityTypes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":4,"name":"Call","icon_key":"call","color":"FFFFFF"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	types, err := client.ActivityTypes.List(
		context.Background(),
		WithActivityTypesRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(types) != 1 || types[0].ID != 4 || types[0].Name != "Call" {
		t.Fatalf("unexpected activity types: %#v", types)
	}
}

func TestActivityTypesService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activityTypes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Video call" {
			t.Fatalf("unexpected name: %#v", payload["name"])
		}
		if payload["icon_key"] != "camera" {
			t.Fatalf("unexpected icon_key: %#v", payload["icon_key"])
		}
		if payload["color"] != "FFFFFF" {
			t.Fatalf("unexpected color: %#v", payload["color"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":12,"name":"Video call","icon_key":"camera","color":"FFFFFF","add_time":"2020-09-01 10:16:23"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	created, err := client.ActivityTypes.Create(
		context.Background(),
		WithActivityTypeName("Video call"),
		WithActivityTypeIconKey("camera"),
		WithActivityTypeColor("FFFFFF"),
		WithActivityTypesRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if created.ID != 12 || created.Name != "Video call" {
		t.Fatalf("unexpected created type: %#v", created)
	}
	if created.AddTime == nil || created.AddTime.IsZero() {
		t.Fatalf("expected add_time to be parsed")
	}
}

func TestActivityTypesService_Update(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activityTypes/12" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Updated" {
			t.Fatalf("unexpected name: %#v", payload["name"])
		}
		if payload["order_nr"] != float64(2) {
			t.Fatalf("unexpected order_nr: %#v", payload["order_nr"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":12,"name":"Updated","icon_key":"camera","order_nr":2}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	updated, err := client.ActivityTypes.Update(
		context.Background(),
		ActivityTypeID(12),
		WithActivityTypeName("Updated"),
		WithActivityTypeOrder(2),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if updated.ID != 12 || updated.Name != "Updated" {
		t.Fatalf("unexpected updated type: %#v", updated)
	}
}

func TestActivityTypesService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activityTypes/12" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":12,"name":"Video call","icon_key":"camera"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	deleted, err := client.ActivityTypes.Delete(context.Background(), ActivityTypeID(12))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if deleted.ID != 12 {
		t.Fatalf("unexpected deleted type: %#v", deleted)
	}
}

func TestActivityTypesService_DeleteBulk(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activityTypes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ids"); got != "7,8" {
			t.Fatalf("unexpected ids: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":[7,8]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.ActivityTypes.DeleteBulk(context.Background(), []ActivityTypeID{7, 8})
	if err != nil {
		t.Fatalf("DeleteBulk error: %v", err)
	}
	if len(result.IDs) != 2 || result.IDs[0] != 7 || result.IDs[1] != 8 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
