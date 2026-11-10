package v1

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestUpdatedSinceOptions(t *testing.T) {
	since := time.Date(2026, 11, 10, 14, 30, 0, 0, time.FixedZone("EET", 2*60*60))
	for _, service := range []string{"notes", "leads"} {
		t.Run(service, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/"+service {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				got, err := time.Parse(time.RFC3339, r.URL.Query().Get("updated_since"))
				if err != nil || !got.Equal(since) {
					t.Fatalf("updated_since = %q, error = %v", r.URL.Query().Get("updated_since"), err)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"success":true,"data":[]}`))
			})
			var err error
			if service == "notes" {
				_, _, err = client.Notes.List(context.Background(), WithNotesUpdatedSince(since.Add(-time.Hour)), WithNotesUpdatedSince(since))
			} else {
				_, _, err = client.Leads.List(context.Background(), WithLeadsUpdatedSince(since.Add(-time.Hour)), WithLeadsUpdatedSince(since))
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
