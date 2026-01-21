package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func TestOrganizationsService_ListCollection(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/collection" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Org"}],"additional_data":{"next_cursor":"c1"}}`))
	})

	query := url.Values{}
	query.Set("limit", "1")
	orgs, page, err := client.Organizations.ListCollection(context.Background(), WithOrganizationsQuery(query))
	if err != nil {
		t.Fatalf("ListCollection error: %v", err)
	}
	if len(orgs) != 1 || orgs[0].ID != 1 || orgs[0].Name != "Org" {
		t.Fatalf("unexpected orgs: %#v", orgs)
	}
	if page == nil || page.NextCursor == nil || *page.NextCursor != "c1" {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestOrganizationsService_Delete(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ids"); got != "2,3" {
			t.Fatalf("unexpected ids: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	ok, err := client.Organizations.Delete(context.Background(), []OrganizationID{2, 3})
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}

func TestOrganizationsService_Merge(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/4/merge" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["merge_with_id"] != float64(5) {
			t.Fatalf("unexpected merge_with_id: %#v", payload["merge_with_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":4,"name":"Merged"}}`))
	})

	org, err := client.Organizations.Merge(context.Background(), OrganizationID(4), OrganizationID(5))
	if err != nil {
		t.Fatalf("Merge error: %v", err)
	}
	if org.ID != 4 || org.Name != "Merged" {
		t.Fatalf("unexpected org: %#v", org)
	}
}

func TestOrganizationsService_ListActivities(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/7/activities" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"subject":"Call"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	activities, page, err := client.Organizations.ListActivities(context.Background(), OrganizationID(7))
	if err != nil {
		t.Fatalf("ListActivities error: %v", err)
	}
	if len(activities) != 1 || activities[0].ID != 1 {
		t.Fatalf("unexpected activities: %#v", activities)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestOrganizationsService_Changelog(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/9/changelog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"field_key":"name"}],"additional_data":{"next_cursor":"c2"}}`))
	})

	changes, page, err := client.Organizations.Changelog(context.Background(), OrganizationID(9))
	if err != nil {
		t.Fatalf("Changelog error: %v", err)
	}
	if len(changes) != 1 || changes[0]["field_key"] != "name" {
		t.Fatalf("unexpected changes: %#v", changes)
	}
	if page == nil || page.NextCursor == nil || *page.NextCursor != "c2" {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestOrganizationsService_ListDeals(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/11/deals" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":2,"title":"Deal"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	deals, page, err := client.Organizations.ListDeals(context.Background(), OrganizationID(11))
	if err != nil {
		t.Fatalf("ListDeals error: %v", err)
	}
	if len(deals) != 1 || deals[0].ID != 2 {
		t.Fatalf("unexpected deals: %#v", deals)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestOrganizationsService_ListFiles(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/13/files" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":5,"name":"file.txt"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	files, page, err := client.Organizations.ListFiles(context.Background(), OrganizationID(13))
	if err != nil {
		t.Fatalf("ListFiles error: %v", err)
	}
	if len(files) != 1 || files[0].ID != 5 {
		t.Fatalf("unexpected files: %#v", files)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestOrganizationsService_ListMailMessages(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/15/mailMessages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":9,"subject":"Hello"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	messages, page, err := client.Organizations.ListMailMessages(context.Background(), OrganizationID(15))
	if err != nil {
		t.Fatalf("ListMailMessages error: %v", err)
	}
	if len(messages) != 1 || messages[0].ID != 9 {
		t.Fatalf("unexpected messages: %#v", messages)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestOrganizationsService_ListPersons(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/17/persons" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":3,"name":"Person"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	people, page, err := client.Organizations.ListPersons(context.Background(), OrganizationID(17))
	if err != nil {
		t.Fatalf("ListPersons error: %v", err)
	}
	if len(people) != 1 || people[0].ID != 3 {
		t.Fatalf("unexpected people: %#v", people)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestOrganizationsService_ListUpdates(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/19/flow" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"type":"note"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	updates, page, err := client.Organizations.ListUpdates(context.Background(), OrganizationID(19))
	if err != nil {
		t.Fatalf("ListUpdates error: %v", err)
	}
	if len(updates) != 1 || updates[0]["type"] != "note" {
		t.Fatalf("unexpected updates: %#v", updates)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestOrganizationsService_ListUsers(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizations/21/permittedUsers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":9,"name":"User"}]}`))
	})

	users, err := client.Organizations.ListUsers(context.Background(), OrganizationID(21))
	if err != nil {
		t.Fatalf("ListUsers error: %v", err)
	}
	if len(users) != 1 || users[0].ID != 9 {
		t.Fatalf("unexpected users: %#v", users)
	}
}
