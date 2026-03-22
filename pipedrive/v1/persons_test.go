package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestPersonsService_Merge(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/4/merge" {
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

	person, err := client.Persons.Merge(context.Background(), PersonID(4), PersonID(5))
	if err != nil {
		t.Fatalf("Merge error: %v", err)
	}
	if person.ID != 4 || person.Name != "Merged" {
		t.Fatalf("unexpected person: %#v", person)
	}
}

func TestPersonsService_Changelog(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/9/changelog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"field_key":"name"}],"additional_data":{"next_cursor":"c2"}}`))
	})

	changes, page, err := client.Persons.Changelog(context.Background(), PersonID(9))
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

func TestPersonsService_ListFiles(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/13/files" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":5,"name":"file.txt"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	files, page, err := client.Persons.ListFiles(context.Background(), PersonID(13))
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

func TestPersonsService_ListMailMessages(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/15/mailMessages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":9,"subject":"Hello"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	messages, page, err := client.Persons.ListMailMessages(context.Background(), PersonID(15))
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

func TestPersonsService_ListProducts(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/17/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":3,"name":"Product"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	products, page, err := client.Persons.ListProducts(context.Background(), PersonID(17))
	if err != nil {
		t.Fatalf("ListProducts error: %v", err)
	}
	if len(products) != 1 || products[0].ID != 3 {
		t.Fatalf("unexpected products: %#v", products)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestPersonsService_ListUpdates(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/19/flow" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"type":"note"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	updates, page, err := client.Persons.ListUpdates(context.Background(), PersonID(19))
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

func TestPersonsService_ListUsers(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/21/permittedUsers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":9,"name":"User"}]}`))
	})

	users, err := client.Persons.ListUsers(context.Background(), PersonID(21))
	if err != nil {
		t.Fatalf("ListUsers error: %v", err)
	}
	if len(users) != 1 || users[0].ID != 9 {
		t.Fatalf("unexpected users: %#v", users)
	}
}

func TestPersonsService_AddPicture(t *testing.T) {
	t.Parallel()

	contentType := "multipart/form-data; boundary=test"
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/25/picture" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != contentType {
			t.Fatalf("unexpected content type: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":123}}`))
	})

	data, err := client.Persons.AddPicture(context.Background(), PersonID(25), strings.NewReader("file"), contentType)
	if err != nil {
		t.Fatalf("AddPicture error: %v", err)
	}
	if data["id"] != float64(123) {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestPersonsService_DeletePicture(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/persons/27/picture" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":27}}`))
	})

	id, err := client.Persons.DeletePicture(context.Background(), PersonID(27))
	if err != nil {
		t.Fatalf("DeletePicture error: %v", err)
	}
	if id != 27 {
		t.Fatalf("unexpected id: %d", id)
	}
}
