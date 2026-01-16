package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestMeetingsService_CreateUserProviderLink(t *testing.T) {
	t.Parallel()

	const (
		userProviderID = "1e3943c9-6395-462b-b432-1f252c017f3d"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/meetings/userProviderLinks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["user_provider_id"] != userProviderID {
			t.Fatalf("unexpected user_provider_id: %#v", payload["user_provider_id"])
		}
		if payload["user_id"] != float64(123) {
			t.Fatalf("unexpected user_id: %#v", payload["user_id"])
		}
		if payload["company_id"] != float64(456) {
			t.Fatalf("unexpected company_id: %#v", payload["company_id"])
		}
		if payload["marketplace_client_id"] != "marketplace-client" {
			t.Fatalf("unexpected marketplace_client_id: %#v", payload["marketplace_client_id"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"message":"The user was added successfully"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Meetings.CreateUserProviderLink(
		context.Background(),
		WithUserProviderLinkID(UserProviderLinkID(userProviderID)),
		WithUserProviderLinkUserID(123),
		WithUserProviderLinkCompanyID(456),
		WithUserProviderLinkMarketplaceClientID("marketplace-client"),
		WithMeetingsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("CreateUserProviderLink error: %v", err)
	}
	if result.Message != "The user was added successfully" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestMeetingsService_DeleteUserProviderLink(t *testing.T) {
	t.Parallel()

	const (
		userProviderID = "1e3943c9-6395-462b-b432-1f252c017f3d"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/meetings/userProviderLinks/"+userProviderID {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"message":"The user data was successfully removed"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Meetings.DeleteUserProviderLink(
		context.Background(),
		UserProviderLinkID(userProviderID),
		WithMeetingsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("DeleteUserProviderLink error: %v", err)
	}
	if result.Message != "The user data was successfully removed" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
