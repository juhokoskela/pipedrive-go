package v2

import (
	"context"
	"encoding/json"
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

func TestActivityFieldsService_ListPager(t *testing.T) {
	t.Parallel()

	runFieldListPagerTest(t, "/activityFields", func(client *Client) *pipedrive.CursorPager[Field] {
		return client.ActivityFields.ListPager(WithActivityFieldsPageSize(2), WithActivityFieldsCursor("start"))
	})
}

func TestActivityFieldsService_ForEach(t *testing.T) {
	t.Parallel()

	runFieldForEachTest(t, "/activityFields", func(ctx context.Context, client *Client, fn func(Field) error) error {
		return client.ActivityFields.ForEach(ctx, fn)
	})
}

func runFieldGetTest(t *testing.T, path string, get func(context.Context, *Client) (*Field, error)) {
	t.Helper()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != path {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("include_fields"); got != "ui_visibility,important_fields" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "get" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"field_code":"cf_1","field_name":"Priority"}}`))
	})

	field, err := get(context.Background(), client)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if field.FieldCode != "cf_1" || field.FieldName != "Priority" {
		t.Fatalf("unexpected field: %#v", field)
	}
}

func runFieldListPagerTest(t *testing.T, path string, pagerFor func(*Client) *pipedrive.CursorPager[Field]) {
	t.Helper()

	var calls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != path {
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
			_, _ = w.Write([]byte(`{"data":[{"field_code":"cf_1"}],"additional_data":{"next_cursor":"next"}}`))
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "next" {
				t.Fatalf("unexpected second cursor: %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"field_code":"cf_2"}],"additional_data":{"next_cursor":null}}`))
		default:
			t.Fatalf("unexpected call count: %d", calls)
		}
	})

	pager := pagerFor(client)
	var codes []string
	for pager.Next(context.Background()) {
		for _, field := range pager.Items() {
			codes = append(codes, field.FieldCode)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if len(codes) != 2 || codes[0] != "cf_1" || codes[1] != "cf_2" {
		t.Fatalf("unexpected field codes: %v", codes)
	}
}

func runFieldForEachTest(t *testing.T, path string, forEach func(context.Context, *Client, func(Field) error) error) {
	t.Helper()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != path {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"field_code":"cf_1"},{"field_code":"cf_2"}],"additional_data":{"next_cursor":null}}`))
	})

	var codes []string
	err := forEach(context.Background(), client, func(field Field) error {
		codes = append(codes, field.FieldCode)
		return nil
	})
	if err != nil {
		t.Fatalf("ForEach error: %v", err)
	}
	if len(codes) != 2 || codes[0] != "cf_1" || codes[1] != "cf_2" {
		t.Fatalf("unexpected field codes: %v", codes)
	}
}

func runFieldUpdateTest(t *testing.T, path string, update func(context.Context, *Client) (*Field, error)) {
	t.Helper()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != path {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "update" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["field_name"] != "Priority updated" {
			t.Fatalf("unexpected field_name: %#v", payload["field_name"])
		}
		if payload["field_type"] != "varchar" {
			t.Fatalf("unexpected field_type: %#v", payload["field_type"])
		}
		if payload["description"] != "Updated description" {
			t.Fatalf("unexpected description: %#v", payload["description"])
		}
		if got := payload["ui_visibility"].(map[string]interface{})["add"]; got != true {
			t.Fatalf("unexpected ui_visibility: %#v", payload["ui_visibility"])
		}
		if got := payload["important_fields"].(map[string]interface{})["pipeline_id"]; got != float64(1) {
			t.Fatalf("unexpected important_fields: %#v", payload["important_fields"])
		}
		if got := payload["required_fields"].(map[string]interface{})["stage_id"]; got != float64(2) {
			t.Fatalf("unexpected required_fields: %#v", payload["required_fields"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"field_code":"cf_1","field_name":"Priority updated"}}`))
	})

	field, err := update(context.Background(), client)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if field.FieldName != "Priority updated" {
		t.Fatalf("unexpected field: %#v", field)
	}
}

func runFieldDeleteTest(t *testing.T, path string, deleteField func(context.Context, *Client) (*Field, error)) {
	t.Helper()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != path {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"field_code":"cf_1","field_name":"Deleted"}}`))
	})

	field, err := deleteField(context.Background(), client)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if field.FieldCode != "cf_1" {
		t.Fatalf("unexpected field: %#v", field)
	}
}

func runFieldAddOptionsRequestOptionsTest(t *testing.T, path string, addOptions func(context.Context, *Client) ([]FieldOption, error)) {
	t.Helper()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != path {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "add-options" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1,"label":"Critical"}]}`))
	})

	options, err := addOptions(context.Background(), client)
	if err != nil {
		t.Fatalf("AddOptions error: %v", err)
	}
	if len(options) != 1 || options[0].Label != "Critical" {
		t.Fatalf("unexpected options: %#v", options)
	}
}

func runFieldUpdateOptionsTest(t *testing.T, path string, updateOptions func(context.Context, *Client) ([]FieldOption, error)) {
	t.Helper()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != path {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "update-options" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload []map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(payload) != 1 || payload[0]["id"] != float64(1) || payload[0]["label"] != "Critical" {
			t.Fatalf("unexpected payload: %#v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1,"label":"Critical"}]}`))
	})

	options, err := updateOptions(context.Background(), client)
	if err != nil {
		t.Fatalf("UpdateOptions error: %v", err)
	}
	if len(options) != 1 || options[0].ID != 1 {
		t.Fatalf("unexpected options: %#v", options)
	}
}

func runFieldDeleteOptionsTest(t *testing.T, path string, deleteOptions func(context.Context, *Client) ([]FieldOption, error)) {
	t.Helper()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != path {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete-options" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload []map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(payload) != 2 || payload[0]["id"] != float64(1) || payload[1]["id"] != float64(2) {
			t.Fatalf("unexpected payload: %#v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}]}`))
	})

	options, err := deleteOptions(context.Background(), client)
	if err != nil {
		t.Fatalf("DeleteOptions error: %v", err)
	}
	if len(options) != 2 {
		t.Fatalf("unexpected options: %#v", options)
	}
}
