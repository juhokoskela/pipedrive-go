package v2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestActivityFieldsService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activityFields" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("include_fields"); got != "ui_visibility,important_fields" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"field_code":"cf_1","field_name":"Stage"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	fields, next, err := client.ActivityFields.List(
		context.Background(),
		WithActivityFieldsIncludeFields(FieldIncludeField("ui_visibility"), FieldIncludeField("important_fields")),
		WithActivityFieldsPageSize(2),
		WithActivityFieldsCursor("c1"),
		WithActivityFieldRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(fields) != 1 || fields[0].FieldCode != "cf_1" {
		t.Fatalf("unexpected fields: %#v", fields)
	}
}

func TestActivityFieldsService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/activityFields/cf_2" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("include_fields"); got != "ui_visibility" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "2" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"field_code":"cf_2","field_name":"Note"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	field, err := client.ActivityFields.Get(
		context.Background(),
		"cf_2",
		WithActivityFieldIncludeFields(FieldIncludeField("ui_visibility")),
		WithActivityFieldRequestOptions(pipedrive.WithHeader("X-Test", "2")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if field.FieldCode != "cf_2" || field.FieldName != "Note" {
		t.Fatalf("unexpected field: %#v", field)
	}
}
