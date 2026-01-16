package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestLeadFieldsService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leadFields" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("start"); got != "50" {
			t.Fatalf("unexpected start: %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "25" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"My Field","key":"my_field","field_type":"varchar","active_flag":true}],"additional_data":{"start":50,"limit":25,"more_items_in_collection":false}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	fields, pagination, err := client.LeadFields.List(
		context.Background(),
		WithLeadFieldsStart(50),
		WithLeadFieldsLimit(25),
		WithLeadFieldsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(fields) != 1 {
		t.Fatalf("unexpected fields: %#v", fields)
	}
	if fields[0].Name != "My Field" || fields[0].Key != "my_field" || fields[0].FieldType != "varchar" {
		t.Fatalf("unexpected field: %#v", fields[0])
	}
	if pagination == nil || pagination.Start != 50 || pagination.Limit != 25 || pagination.MoreItemsInCollection {
		t.Fatalf("unexpected pagination: %#v", pagination)
	}
}
