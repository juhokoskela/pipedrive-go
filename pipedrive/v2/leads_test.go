package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestLeadsService_Search(t *testing.T) {
	t.Parallel()

	leadID := "123e4567-e89b-12d3-a456-426614174000"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("term"); got != "alpha" {
			t.Fatalf("unexpected term: %q", got)
		}
		if got := q.Get("fields"); got != "title,notes" {
			t.Fatalf("unexpected fields: %q", got)
		}
		if got := q.Get("exact_match"); got != "true" {
			t.Fatalf("unexpected exact_match: %q", got)
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
		_, _ = w.Write([]byte(`{"data":{"items":[{"result_score":0.8,"item":{"id":"` + leadID + `"}}]},"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	results, next, err := client.Leads.Search(
		context.Background(),
		"alpha",
		WithLeadSearchFields(LeadSearchFieldTitle, LeadSearchFieldNotes),
		WithLeadSearchExactMatch(true),
		WithLeadSearchPageSize(2),
		WithLeadSearchCursor("c1"),
		WithLeadRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(results.Items) != 1 {
		t.Fatalf("unexpected results: %#v", results)
	}
}

func TestLeadsService_ConvertToDeal(t *testing.T) {
	t.Parallel()

	leadID := "123e4567-e89b-12d3-a456-426614174000"
	conversionID := "223e4567-e89b-12d3-a456-426614174111"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads/"+leadID+"/convert/deal" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "2" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["stage_id"] != float64(10) {
			t.Fatalf("unexpected stage_id: %v", payload["stage_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"conversion_id":"` + conversionID + `"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	job, err := client.Leads.ConvertToDeal(
		context.Background(),
		LeadID(leadID),
		WithLeadConversionStageID(StageID(10)),
		WithLeadRequestOptions(pipedrive.WithHeader("X-Test", "2")),
	)
	if err != nil {
		t.Fatalf("ConvertToDeal error: %v", err)
	}
	if job.ConversionID != ConversionID(conversionID) {
		t.Fatalf("unexpected conversion job: %#v", job)
	}
}

func TestLeadsService_ConversionStatus(t *testing.T) {
	t.Parallel()

	leadID := "123e4567-e89b-12d3-a456-426614174000"
	conversionID := "223e4567-e89b-12d3-a456-426614174111"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/leads/"+leadID+"/convert/status/"+conversionID {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"conversion_id":"` + conversionID + `","status":"completed","deal_id":5}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	status, err := client.Leads.ConversionStatus(
		context.Background(),
		LeadID(leadID),
		ConversionID(conversionID),
	)
	if err != nil {
		t.Fatalf("ConversionStatus error: %v", err)
	}
	if status.Status != ConversionStatusCompleted || status.DealID == nil || *status.DealID != 5 {
		t.Fatalf("unexpected status: %#v", status)
	}
}
