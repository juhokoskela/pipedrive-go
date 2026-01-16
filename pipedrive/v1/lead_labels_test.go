package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestLeadLabelsService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leadLabels" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":"f08b42a0-4e75-11ea-9643-03698ef1cfd6","name":"Hot","color":"red","add_time":"2020-02-13T15:31:44.000Z"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	labels, err := client.LeadLabels.List(
		context.Background(),
		WithLeadLabelsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(labels) != 1 || labels[0].ID != "f08b42a0-4e75-11ea-9643-03698ef1cfd6" || labels[0].Name != "Hot" {
		t.Fatalf("unexpected labels: %#v", labels)
	}
	if labels[0].AddTime == nil || labels[0].AddTime.IsZero() {
		t.Fatalf("expected add_time to be parsed")
	}
}

func TestLeadLabelsService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leadLabels" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Hot" {
			t.Fatalf("unexpected name: %#v", payload["name"])
		}
		if payload["color"] != "red" {
			t.Fatalf("unexpected color: %#v", payload["color"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"f08b42a0-4e75-11ea-9643-03698ef1cfd6","name":"Hot","color":"red","add_time":"2020-02-13T15:31:44.000Z","update_time":"2020-02-13T15:31:44.000Z"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	label, err := client.LeadLabels.Create(
		context.Background(),
		WithLeadLabelName("Hot"),
		WithLeadLabelColor(LeadLabelColorRed),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if label.ID != "f08b42a0-4e75-11ea-9643-03698ef1cfd6" || label.Color != LeadLabelColorRed {
		t.Fatalf("unexpected label: %#v", label)
	}
}

func TestLeadLabelsService_Update(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leadLabels/f08b42a0-4e75-11ea-9643-03698ef1cfd6" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["color"] != "blue" {
			t.Fatalf("unexpected color: %#v", payload["color"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"f08b42a0-4e75-11ea-9643-03698ef1cfd6","name":"Hot","color":"blue"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	label, err := client.LeadLabels.Update(
		context.Background(),
		LeadLabelID("f08b42a0-4e75-11ea-9643-03698ef1cfd6"),
		WithLeadLabelColor(LeadLabelColorBlue),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if label.ID != "f08b42a0-4e75-11ea-9643-03698ef1cfd6" || label.Color != LeadLabelColorBlue {
		t.Fatalf("unexpected label: %#v", label)
	}
}

func TestLeadLabelsService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leadLabels/f08b42a0-4e75-11ea-9643-03698ef1cfd6" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"f08b42a0-4e75-11ea-9643-03698ef1cfd6"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.LeadLabels.Delete(context.Background(), LeadLabelID("f08b42a0-4e75-11ea-9643-03698ef1cfd6"))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if result.ID != "f08b42a0-4e75-11ea-9643-03698ef1cfd6" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
