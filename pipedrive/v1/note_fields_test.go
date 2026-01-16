package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestNoteFieldsService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/noteFields" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":3,"name":"Note Field","key":"note_field","field_type":"text"}],"additional_data":{"start":0,"limit":1,"more_items_in_collection":false}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	fields, pagination, err := client.NoteFields.List(
		context.Background(),
		WithNoteFieldsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(fields) != 1 || fields[0].Name != "Note Field" || fields[0].Key != "note_field" {
		t.Fatalf("unexpected fields: %#v", fields)
	}
	if pagination == nil || pagination.Limit != 1 || pagination.MoreItemsInCollection {
		t.Fatalf("unexpected pagination: %#v", pagination)
	}
}
