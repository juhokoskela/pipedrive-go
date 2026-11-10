package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/v2/pipedrive"
)

func TestDateTime_MarshalPreservesInstant(t *testing.T) {
	t.Parallel()
	original := time.Date(2026, 1, 2, 1, 30, 0, 0, time.FixedZone("UTC+2", 2*60*60))
	data, err := json.Marshal(DateTime{Time: original})
	if err != nil {
		t.Fatal(err)
	}
	var decoded DateTime
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if !decoded.Equal(original) {
		t.Fatalf("round trip = %s, want instant %s", decoded.Time, original)
	}
}

func TestNotes_CreateNormalizesAddTimeToUTC(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			AddTime string `json:"add_time"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body.AddTime != "2026-01-01 23:30:00" {
			t.Errorf("add_time = %q", body.AddTime)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":1}}`))
	}))
	t.Cleanup(srv.Close)
	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Notes.Create(context.Background(), WithNoteContent("test"), WithNoteDealID(1), WithNoteAddTime(time.Date(2026, 1, 2, 1, 30, 0, 0, time.FixedZone("UTC+2", 7200))))
	if err != nil {
		t.Fatal(err)
	}
}
