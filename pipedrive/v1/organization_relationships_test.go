package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestOrganizationRelationshipsService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizationRelationships" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("org_id"); got != "10" {
			t.Fatalf("unexpected org_id: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"type":"parent","related_organization_name":"Telia","calculated_type":"daughter","calculated_related_org_id":1480,"rel_owner_org_id":{"value":1481,"name":"Pipedrive Inc."},"rel_linked_org_id":{"value":1480,"name":"Telia"},"add_time":"2020-09-22 08:58:28","update_time":"2020-09-22 08:58:28","active_flag":"true"}],"additional_data":{"pagination":{"start":0,"limit":100,"more_items_in_collection":false}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	relationships, additional, err := client.OrganizationRelationships.List(
		context.Background(),
		WithOrganizationRelationshipsOrgID(OrganizationID(10)),
		WithOrganizationRelationshipsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(relationships) != 1 || relationships[0].ID != 1 || relationships[0].CalculatedType != "daughter" {
		t.Fatalf("unexpected relationships: %#v", relationships)
	}
	if relationships[0].AddTime == nil || relationships[0].AddTime.IsZero() {
		t.Fatalf("expected add_time to be parsed")
	}
	if additional == nil || additional.Pagination == nil {
		t.Fatalf("unexpected additional data: %#v", additional)
	}
}

func TestOrganizationRelationshipsService_ListRequiresOrgID(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.OrganizationRelationships.List(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOrganizationRelationshipsService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizationRelationships" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["type"] != "parent" {
			t.Fatalf("unexpected type: %#v", payload["type"])
		}
		if payload["rel_owner_org_id"] != float64(1) {
			t.Fatalf("unexpected rel_owner_org_id: %#v", payload["rel_owner_org_id"])
		}
		if payload["rel_linked_org_id"] != float64(2) {
			t.Fatalf("unexpected rel_linked_org_id: %#v", payload["rel_linked_org_id"])
		}
		if payload["org_id"] != float64(10) {
			t.Fatalf("unexpected org_id: %#v", payload["org_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":2,"type":"parent","rel_owner_org_id":{"value":1},"rel_linked_org_id":{"value":2}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	relationship, err := client.OrganizationRelationships.Create(
		context.Background(),
		WithOrganizationRelationshipType(OrganizationRelationshipTypeParent),
		WithOrganizationRelationshipOwnerID(OrganizationID(1)),
		WithOrganizationRelationshipLinkedID(OrganizationID(2)),
		WithOrganizationRelationshipOrgID(OrganizationID(10)),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if relationship.ID != 2 {
		t.Fatalf("unexpected relationship: %#v", relationship)
	}
}

func TestOrganizationRelationshipsService_CreateRequiresFields(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.OrganizationRelationships.Create(context.Background(), WithOrganizationRelationshipType(OrganizationRelationshipTypeParent))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOrganizationRelationshipsService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizationRelationships/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("org_id"); got != "10" {
			t.Fatalf("unexpected org_id: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5,"type":"parent","rel_owner_org_id":{"value":1},"rel_linked_org_id":{"value":2}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	relationship, err := client.OrganizationRelationships.Get(
		context.Background(),
		OrganizationRelationshipID(5),
		WithOrganizationRelationshipOrgID(OrganizationID(10)),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if relationship.ID != 5 {
		t.Fatalf("unexpected relationship: %#v", relationship)
	}
}

func TestOrganizationRelationshipsService_Update(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizationRelationships/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["type"] != "related" {
			t.Fatalf("unexpected type: %#v", payload["type"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5,"type":"related","rel_owner_org_id":{"value":1},"rel_linked_org_id":{"value":2}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	relationship, err := client.OrganizationRelationships.Update(
		context.Background(),
		OrganizationRelationshipID(5),
		WithOrganizationRelationshipType(OrganizationRelationshipTypeRelated),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if relationship.ID != 5 || relationship.Type != OrganizationRelationshipTypeRelated {
		t.Fatalf("unexpected relationship: %#v", relationship)
	}
}

func TestOrganizationRelationshipsService_UpdateRequiresFields(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.OrganizationRelationships.Update(context.Background(), OrganizationRelationshipID(5))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOrganizationRelationshipsService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizationRelationships/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.OrganizationRelationships.Delete(context.Background(), OrganizationRelationshipID(5))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if result.ID != 5 {
		t.Fatalf("unexpected delete result: %#v", result)
	}
}
