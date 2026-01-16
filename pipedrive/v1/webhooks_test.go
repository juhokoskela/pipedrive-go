package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestWebhooksService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/webhooks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"company_id":1,"owner_id":2,"user_id":3,"event_action":"create","event_object":"deal","subscription_url":"http://example.org","version":"2.0","is_active":1,"add_time":"2019-10-25T08:25:27.000Z","type":"general","name":"Example webhook"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	webhooks, err := client.Webhooks.List(
		context.Background(),
		WithWebhooksRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(webhooks) != 1 || webhooks[0].Name != "Example webhook" {
		t.Fatalf("unexpected webhooks: %#v", webhooks)
	}
	if !webhooks[0].IsActive {
		t.Fatalf("expected webhook to be active")
	}
}

func TestWebhooksService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/webhooks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["subscription_url"] != "http://example.org" {
			t.Fatalf("unexpected subscription_url: %#v", payload["subscription_url"])
		}
		if payload["event_action"] != "create" {
			t.Fatalf("unexpected event_action: %#v", payload["event_action"])
		}
		if payload["event_object"] != "deal" {
			t.Fatalf("unexpected event_object: %#v", payload["event_object"])
		}
		if payload["name"] != "Webhook" {
			t.Fatalf("unexpected name: %#v", payload["name"])
		}
		if payload["version"] != "2.0" {
			t.Fatalf("unexpected version: %#v", payload["version"])
		}
		if payload["user_id"] != float64(7) {
			t.Fatalf("unexpected user_id: %#v", payload["user_id"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":9,"name":"Webhook","event_action":"create","event_object":"deal","subscription_url":"http://example.org","is_active":1}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	webhook, err := client.Webhooks.Create(
		context.Background(),
		WithWebhookSubscriptionURL("http://example.org"),
		WithWebhookEventAction(WebhookEventActionCreate),
		WithWebhookEventObject(WebhookEventObjectDeal),
		WithWebhookName("Webhook"),
		WithWebhookUserID(7),
		WithWebhookVersion(WebhookVersion2),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if webhook.ID != 9 || webhook.Name != "Webhook" {
		t.Fatalf("unexpected webhook: %#v", webhook)
	}
}

func TestWebhooksService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/webhooks/9" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"status":"ok"}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ok, err := client.Webhooks.Delete(context.Background(), 9)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}
