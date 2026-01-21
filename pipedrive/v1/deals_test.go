package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestDealsService_ListCollection(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/collection" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"title":"Deal A"}],"additional_data":{"next_cursor":"c1"}}`))
	})

	query := url.Values{}
	query.Set("limit", "1")
	deals, page, err := client.Deals.ListCollection(
		context.Background(),
		WithDealsQuery(query),
		WithDealsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("ListCollection error: %v", err)
	}
	if len(deals) != 1 || deals[0].ID != 1 || deals[0].Title != "Deal A" {
		t.Fatalf("unexpected deals: %#v", deals)
	}
	if page == nil || page.NextCursor == nil || *page.NextCursor != "c1" {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestDealsService_Summary(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/summary" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("status"); got != "open" {
			t.Fatalf("unexpected status: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"total_count":2,"total_currency_converted_value":123.5}}`))
	})

	query := url.Values{}
	query.Set("status", "open")
	summary, err := client.Deals.Summary(context.Background(), WithDealsQuery(query))
	if err != nil {
		t.Fatalf("Summary error: %v", err)
	}
	if summary.TotalCount != 2 || summary.TotalCurrencyConvertedValue != 123.5 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}

func TestDealsService_ArchivedSummary(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/summary/archived" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"total_count":1}}`))
	})

	summary, err := client.Deals.ArchivedSummary(context.Background())
	if err != nil {
		t.Fatalf("ArchivedSummary error: %v", err)
	}
	if summary.TotalCount != 1 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}

func TestDealsService_Timeline(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/timeline" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"periods":[{"count":1}]}}`))
	})

	timeline, err := client.Deals.Timeline(context.Background())
	if err != nil {
		t.Fatalf("Timeline error: %v", err)
	}
	if _, ok := timeline["periods"]; !ok {
		t.Fatalf("expected timeline periods: %#v", timeline)
	}
}

func TestDealsService_ArchivedTimeline(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/timeline/archived" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"periods":[]}}`))
	})

	_, err := client.Deals.ArchivedTimeline(context.Background())
	if err != nil {
		t.Fatalf("ArchivedTimeline error: %v", err)
	}
}

func TestDealsService_ListActivities(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/7/activities" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":3,"subject":"Call"}],"additional_data":{"pagination":{"start":0,"limit":2,"more_items_in_collection":false}}}`))
	})

	query := url.Values{}
	query.Set("limit", "2")
	activities, page, err := client.Deals.ListActivities(context.Background(), DealID(7), WithDealsQuery(query))
	if err != nil {
		t.Fatalf("ListActivities error: %v", err)
	}
	if len(activities) != 1 || activities[0].ID != 3 {
		t.Fatalf("unexpected activities: %#v", activities)
	}
	if page == nil || page.Limit != 2 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestDealsService_Changelog(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/9/changelog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"field_key":"title"}],"additional_data":{"next_cursor":"c2"}}`))
	})

	changes, page, err := client.Deals.Changelog(context.Background(), DealID(9))
	if err != nil {
		t.Fatalf("Changelog error: %v", err)
	}
	if len(changes) != 1 || changes[0]["field_key"] != "title" {
		t.Fatalf("unexpected changes: %#v", changes)
	}
	if page == nil || page.NextCursor == nil || *page.NextCursor != "c2" {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestDealsService_ListFiles(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/5/files" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":11,"name":"file.txt"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	files, page, err := client.Deals.ListFiles(context.Background(), DealID(5))
	if err != nil {
		t.Fatalf("ListFiles error: %v", err)
	}
	if len(files) != 1 || files[0].ID != 11 {
		t.Fatalf("unexpected files: %#v", files)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestDealsService_ListMailMessages(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/6/mailMessages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":21,"subject":"Hello"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	messages, page, err := client.Deals.ListMailMessages(context.Background(), DealID(6))
	if err != nil {
		t.Fatalf("ListMailMessages error: %v", err)
	}
	if len(messages) != 1 || messages[0].ID != 21 {
		t.Fatalf("unexpected messages: %#v", messages)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestDealsService_ListParticipants(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/8/participants" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":13,"name":"Person"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	participants, page, err := client.Deals.ListParticipants(context.Background(), DealID(8))
	if err != nil {
		t.Fatalf("ListParticipants error: %v", err)
	}
	if len(participants) != 1 || participants[0].ID != 13 {
		t.Fatalf("unexpected participants: %#v", participants)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestDealsService_AddParticipant(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/10/participants" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["person_id"] != float64(44) {
			t.Fatalf("unexpected person_id: %#v", payload["person_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":44,"name":"Participant"}}`))
	})

	person, err := client.Deals.AddParticipant(context.Background(), DealID(10), PersonID(44))
	if err != nil {
		t.Fatalf("AddParticipant error: %v", err)
	}
	if person.ID != 44 {
		t.Fatalf("unexpected person: %#v", person)
	}
}

func TestDealsService_DeleteParticipant(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/10/participants/55" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	ok, err := client.Deals.DeleteParticipant(context.Background(), DealID(10), DealParticipantID(55))
	if err != nil {
		t.Fatalf("DeleteParticipant error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}

func TestDealsService_ParticipantsChangelog(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/12/participantsChangelog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"action":"added"}],"additional_data":{"next_cursor":"n1"}}`))
	})

	changes, page, err := client.Deals.ParticipantsChangelog(context.Background(), DealID(12))
	if err != nil {
		t.Fatalf("ParticipantsChangelog error: %v", err)
	}
	if len(changes) != 1 || changes[0]["action"] != "added" {
		t.Fatalf("unexpected changes: %#v", changes)
	}
	if page == nil || page.NextCursor == nil || *page.NextCursor != "n1" {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestDealsService_ListPersons(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/15/persons" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":31,"name":"Person"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	people, page, err := client.Deals.ListPersons(context.Background(), DealID(15))
	if err != nil {
		t.Fatalf("ListPersons error: %v", err)
	}
	if len(people) != 1 || people[0].ID != 31 {
		t.Fatalf("unexpected people: %#v", people)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestDealsService_ListUpdates(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/17/flow" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"type":"note"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	updates, page, err := client.Deals.ListUpdates(context.Background(), DealID(17))
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

func TestDealsService_ListUsers(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/20/permittedUsers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":9,"name":"User"}]}`))
	})

	users, err := client.Deals.ListUsers(context.Background(), DealID(20))
	if err != nil {
		t.Fatalf("ListUsers error: %v", err)
	}
	if len(users) != 1 || users[0].ID != 9 {
		t.Fatalf("unexpected users: %#v", users)
	}
}

func TestDealsService_Merge(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/22/merge" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["merge_with_id"] != float64(23) {
			t.Fatalf("unexpected merge_with_id: %#v", payload["merge_with_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":22,"title":"Merged"}}`))
	})

	deal, err := client.Deals.Merge(context.Background(), DealID(22), DealID(23))
	if err != nil {
		t.Fatalf("Merge error: %v", err)
	}
	if deal.ID != 22 || deal.Title != "Merged" {
		t.Fatalf("unexpected deal: %#v", deal)
	}
}

func TestDealsService_Duplicate(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/24/duplicate" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":99,"title":"Duplicate"}}`))
	})

	deal, err := client.Deals.Duplicate(context.Background(), DealID(24))
	if err != nil {
		t.Fatalf("Duplicate error: %v", err)
	}
	if deal.ID != 99 {
		t.Fatalf("unexpected deal: %#v", deal)
	}
}

func TestDealsService_Delete(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ids"); got != "1,2" {
			t.Fatalf("unexpected ids: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	ok, err := client.Deals.Delete(context.Background(), []DealID{1, 2})
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}
