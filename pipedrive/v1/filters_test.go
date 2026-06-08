package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestFiltersService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/filters" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("type"); got != "deals" {
			t.Fatalf("unexpected type: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Pipeline Deals","type":"deals"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	filters, err := client.Filters.List(
		context.Background(),
		WithFiltersType(FilterTypeDeals),
		WithFiltersRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(filters) != 1 || filters[0].Name != "Pipeline Deals" {
		t.Fatalf("unexpected filters: %#v", filters)
	}
}

func TestFiltersService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/filters/42" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":42,"name":"My Filter","type":"deals"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	filter, err := client.Filters.Get(
		context.Background(),
		FilterID(42),
		WithFiltersRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if filter.ID != 42 || filter.Name != "My Filter" {
		t.Fatalf("unexpected filter: %#v", filter)
	}
}

func TestFiltersService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/filters" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "create" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["name"] != "New Filter" {
			t.Fatalf("unexpected name: %#v", payload["name"])
		}
		if payload["type"] != "deals" {
			t.Fatalf("unexpected type: %#v", payload["type"])
		}
		if payload["conditions"] == nil {
			t.Fatalf("missing conditions")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":10,"name":"New Filter","type":"deals","conditions":{"glue":"and","conditions":[]}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	conditions := FilterConditions{
		"glue":       "and",
		"conditions": []interface{}{},
	}
	filter, err := client.Filters.Create(
		context.Background(),
		WithFilterName("New Filter"),
		WithFilterType(FilterTypeDeals),
		WithFilterConditions(conditions),
		WithFiltersRequestOptions(pipedrive.WithHeader("X-Test", "create")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if filter.ID != 10 || filter.Name != "New Filter" {
		t.Fatalf("unexpected filter: %#v", filter)
	}
}

func TestFiltersService_Update(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/filters/10" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "update" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["name"] != "Updated Filter" {
			t.Fatalf("unexpected name: %#v", payload["name"])
		}
		if _, ok := payload["type"]; ok {
			t.Fatalf("unexpected type in update payload")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":10,"name":"Updated Filter","type":"deals"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	filter, err := client.Filters.Update(
		context.Background(),
		FilterID(10),
		WithFilterName("Updated Filter"),
		WithFilterConditions(FilterConditions{"glue": "and"}),
		WithFiltersRequestOptions(pipedrive.WithHeader("X-Test", "update")),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if filter.Name != "Updated Filter" {
		t.Fatalf("unexpected filter: %#v", filter)
	}
}

func TestFiltersService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/filters/10" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":10}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Filters.Delete(
		context.Background(),
		FilterID(10),
		WithFiltersRequestOptions(pipedrive.WithHeader("X-Test", "delete")),
	)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if result.ID != 10 {
		t.Fatalf("unexpected delete result: %#v", result)
	}
}

func TestFiltersService_DeleteBulk(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/filters" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ids"); got != "1,2,3" {
			t.Fatalf("unexpected ids: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "delete-bulk" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":[1,2,3]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Filters.DeleteBulk(
		context.Background(),
		[]FilterID{1, 2, 3},
		WithFiltersRequestOptions(pipedrive.WithHeader("X-Test", "delete-bulk")),
	)
	if err != nil {
		t.Fatalf("DeleteBulk error: %v", err)
	}
	if result == nil || len(result.IDs) != 3 || result.IDs[2] != 3 {
		t.Fatalf("unexpected delete result: %#v", result)
	}
}

func TestFiltersService_ListHelpers(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/filters/helpers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "helpers" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"example":1}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	helpers, err := client.Filters.ListHelpers(
		context.Background(),
		WithFiltersRequestOptions(pipedrive.WithHeader("X-Test", "helpers")),
	)
	if err != nil {
		t.Fatalf("ListHelpers error: %v", err)
	}
	if helpers["success"] != true {
		t.Fatalf("unexpected helpers: %#v", helpers)
	}
}

func TestFiltersService_ErrorResponses(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"success":false,"error":"server error"}`))
	})

	conditions := FilterConditions{"glue": "and"}
	if _, err := client.Filters.List(context.Background()); err == nil {
		t.Fatalf("expected List error")
	}
	if _, err := client.Filters.Get(context.Background(), FilterID(1)); err == nil {
		t.Fatalf("expected Get error")
	}
	if _, err := client.Filters.Create(
		context.Background(),
		WithFilterName("Filter"),
		WithFilterType(FilterTypeDeals),
		WithFilterConditions(conditions),
	); err == nil {
		t.Fatalf("expected Create error")
	}
	if _, err := client.Filters.Update(
		context.Background(),
		FilterID(1),
		WithFilterName("Updated"),
	); err == nil {
		t.Fatalf("expected Update error")
	}
	if _, err := client.Filters.Delete(context.Background(), FilterID(1)); err == nil {
		t.Fatalf("expected Delete error")
	}
	if _, err := client.Filters.DeleteBulk(context.Background(), []FilterID{1, 2}); err == nil {
		t.Fatalf("expected DeleteBulk error")
	}
	if _, err := client.Filters.ListHelpers(context.Background()); err == nil {
		t.Fatalf("expected ListHelpers error")
	}
}

func TestFiltersService_MissingDataErrors(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	conditions := FilterConditions{"glue": "and"}
	if _, err := client.Filters.Get(context.Background(), FilterID(1)); err == nil {
		t.Fatalf("expected Get missing data error")
	}
	if _, err := client.Filters.Create(
		context.Background(),
		WithFilterName("Filter"),
		WithFilterType(FilterTypeDeals),
		WithFilterConditions(conditions),
	); err == nil {
		t.Fatalf("expected Create missing data error")
	}
	if _, err := client.Filters.Update(
		context.Background(),
		FilterID(1),
		WithFilterName("Updated"),
	); err == nil {
		t.Fatalf("expected Update missing data error")
	}
	if _, err := client.Filters.Delete(context.Background(), FilterID(1)); err == nil {
		t.Fatalf("expected Delete missing data error")
	}
	if _, err := client.Filters.DeleteBulk(context.Background(), []FilterID{1, 2}); err == nil {
		t.Fatalf("expected DeleteBulk missing data error")
	}
}
