package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestUserConnectionsService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
		want    *string
	}{
		{
			name:    "connected",
			payload: `{"success":true,"data":{"google":"google-123"}}`,
			want:    ptr("google-123"),
		},
		{
			name:    "not-connected",
			payload: `{"success":true,"data":{"google":false}}`,
			want:    nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatalf("unexpected method: %s", r.Method)
				}
				if r.URL.Path != "/userConnections" {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				if got := r.Header.Get("X-Test"); got != "1" {
					t.Fatalf("unexpected header X-Test: %q", got)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tc.payload))
			}))
			t.Cleanup(srv.Close)

			client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
			if err != nil {
				t.Fatalf("NewClient error: %v", err)
			}

			connections, err := client.UserConnections.Get(
				context.Background(),
				WithUserConnectionsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
			)
			if err != nil {
				t.Fatalf("Get error: %v", err)
			}
			if tc.want == nil {
				if connections.GoogleID != nil {
					t.Fatalf("expected nil GoogleID, got %v", *connections.GoogleID)
				}
				return
			}
			if connections.GoogleID == nil || *connections.GoogleID != *tc.want {
				t.Fatalf("unexpected GoogleID: %#v", connections.GoogleID)
			}
		})
	}
}

func ptr(value string) *string {
	return &value
}
