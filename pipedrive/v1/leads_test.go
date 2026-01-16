package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestLeadsService_List(t *testing.T) {
	t.Parallel()

	labelID := "f08b42a0-4e75-11ea-9643-03698ef1cfd6"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("start"); got != "1" {
			t.Fatalf("unexpected start: %q", got)
		}
		if got := q.Get("owner_id"); got != "7" {
			t.Fatalf("unexpected owner_id: %q", got)
		}
		if got := q.Get("person_id"); got != "9" {
			t.Fatalf("unexpected person_id: %q", got)
		}
		if got := q.Get("organization_id"); got != "11" {
			t.Fatalf("unexpected organization_id: %q", got)
		}
		if got := q.Get("filter_id"); got != "3" {
			t.Fatalf("unexpected filter_id: %q", got)
		}
		if got := q.Get("sort"); got != "add_time DESC" {
			t.Fatalf("unexpected sort: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":"adf21080-0e10-11eb-879b-05d71fb426ec","title":"Lead","label_ids":["` + labelID + `"],"person_id":9,"owner_id":7,"add_time":"2020-10-14T11:30:36.551Z","update_time":"2020-10-14T11:30:36.551Z"}],"additional_data":{"start":1,"limit":2,"more_items_in_collection":false}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	leads, page, err := client.Leads.List(
		context.Background(),
		WithLeadsLimit(2),
		WithLeadsStart(1),
		WithLeadsOwnerID(UserID(7)),
		WithLeadsPersonID(PersonID(9)),
		WithLeadsOrganizationID(OrganizationID(11)),
		WithLeadsFilterID(3),
		WithLeadsSort("add_time DESC"),
		WithLeadsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(leads) != 1 || leads[0].Title != "Lead" {
		t.Fatalf("unexpected leads: %#v", leads)
	}
	if leads[0].AddTime == nil || leads[0].AddTime.IsZero() {
		t.Fatalf("expected add_time to be parsed")
	}
	if page == nil || page.MoreItemsInCollection {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestLeadsService_ListArchived(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads/archived" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":"adf21080-0e10-11eb-879b-05d71fb426ec","title":"Archived","is_archived":true,"add_time":"2020-10-14T11:30:36.551Z"}],"additional_data":{"start":0,"limit":1,"more_items_in_collection":false}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	leads, page, err := client.Leads.ListArchived(context.Background(), WithArchivedLeadsLimit(1))
	if err != nil {
		t.Fatalf("ListArchived error: %v", err)
	}
	if len(leads) != 1 || !leads[0].IsArchived {
		t.Fatalf("unexpected leads: %#v", leads)
	}
	if page == nil || page.MoreItemsInCollection {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestLeadsService_Get(t *testing.T) {
	t.Parallel()

	id := "adf21080-0e10-11eb-879b-05d71fb426ec"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads/"+id {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"` + id + `","title":"Lead","add_time":"2020-10-14T11:30:36.551Z"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	lead, err := client.Leads.Get(context.Background(), LeadID(id))
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if lead.ID != LeadID(id) {
		t.Fatalf("unexpected lead: %#v", lead)
	}
}

func TestLeadsService_Create(t *testing.T) {
	t.Parallel()

	labelID := "f08b42a0-4e75-11ea-9643-03698ef1cfd6"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["title"] != "New Lead" {
			t.Fatalf("unexpected title: %#v", payload["title"])
		}
		if payload["person_id"] != float64(9) {
			t.Fatalf("unexpected person_id: %#v", payload["person_id"])
		}
		if payload["expected_close_date"] != "2024-01-10" {
			t.Fatalf("unexpected expected_close_date: %#v", payload["expected_close_date"])
		}
		if payload["origin_id"] != "origin-1" {
			t.Fatalf("unexpected origin_id: %#v", payload["origin_id"])
		}
		labels := payload["label_ids"].([]interface{})
		if len(labels) != 1 || labels[0] != labelID {
			t.Fatalf("unexpected label_ids: %#v", payload["label_ids"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"adf21080-0e10-11eb-879b-05d71fb426ec","title":"New Lead","label_ids":["` + labelID + `"],"add_time":"2020-10-14T11:30:36.551Z"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	lead, err := client.Leads.Create(
		context.Background(),
		WithLeadTitle("New Lead"),
		WithLeadPersonID(PersonID(9)),
		WithLeadLabelIDs(LeadLabelID(labelID)),
		WithLeadExpectedCloseDate("2024-01-10"),
		WithLeadOriginID("origin-1"),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if lead.Title != "New Lead" {
		t.Fatalf("unexpected lead: %#v", lead)
	}
}

func TestLeadsService_Update(t *testing.T) {
	t.Parallel()

	id := "adf21080-0e10-11eb-879b-05d71fb426ec"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads/"+id {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["title"] != "Updated" {
			t.Fatalf("unexpected title: %#v", payload["title"])
		}
		if payload["is_archived"] != true {
			t.Fatalf("unexpected is_archived: %#v", payload["is_archived"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"` + id + `","title":"Updated","is_archived":true}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	lead, err := client.Leads.Update(
		context.Background(),
		LeadID(id),
		WithLeadTitle("Updated"),
		WithLeadArchived(true),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if lead.Title != "Updated" || !lead.IsArchived {
		t.Fatalf("unexpected lead: %#v", lead)
	}
}

func TestLeadsService_Delete(t *testing.T) {
	t.Parallel()

	id := "adf21080-0e10-11eb-879b-05d71fb426ec"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads/"+id {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"` + id + `"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Leads.Delete(context.Background(), LeadID(id))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if result.ID != LeadID(id) {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestLeadsService_ListPermittedUsers(t *testing.T) {
	t.Parallel()

	id := "adf21080-0e10-11eb-879b-05d71fb426ec"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads/"+id+"/permittedUsers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[101,202]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	users, err := client.Leads.ListPermittedUsers(context.Background(), LeadID(id))
	if err != nil {
		t.Fatalf("ListPermittedUsers error: %v", err)
	}
	if len(users) != 2 || users[1] != 202 {
		t.Fatalf("unexpected users: %#v", users)
	}
}
