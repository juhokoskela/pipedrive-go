package v1

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestChannelsService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/channels" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "My Channel" {
			t.Fatalf("unexpected name: %#v", payload["name"])
		}
		if payload["provider_channel_id"] != "provider-1" {
			t.Fatalf("unexpected provider_channel_id: %#v", payload["provider_channel_id"])
		}
		if payload["provider_type"] != "other" {
			t.Fatalf("unexpected provider_type: %#v", payload["provider_type"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"chan-1","name":"My Channel","provider_channel_id":"provider-1","created_at":"2022-03-01 00:00:00","provider_type":"other","template_support":false}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	channel, err := client.Channels.Create(
		context.Background(),
		WithChannelName("My Channel"),
		WithChannelProviderID("provider-1"),
		WithChannelProviderType(ChannelProviderTypeOther),
		WithChannelsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if channel.ID != "chan-1" || channel.Name != "My Channel" {
		t.Fatalf("unexpected channel: %#v", channel)
	}
	if channel.CreatedAt == nil || channel.CreatedAt.IsZero() {
		t.Fatalf("expected created_at to be parsed")
	}
}

func TestChannelsService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/channels/chan-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ok, err := client.Channels.Delete(context.Background(), ChannelID("chan-1"))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}

func TestChannelsService_ReceiveMessage(t *testing.T) {
	t.Parallel()

	created := time.Date(2022, 3, 1, 7, 58, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/channels/messages/receive" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !strings.Contains(string(body), "\"id\":\"msg-1\"") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"status\":\"sent\"") {
			t.Fatalf("unexpected body: %s", body)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"msg-1","channel_id":"chan-1","sender_id":"user-1","conversation_id":"conv-1","message":"Hello","status":"sent","created_at":"2022-03-01 07:58","attachments":[{"id":"att-1","type":"image/png","name":"Image","size":123,"url":"http://example.com/a.png","preview_url":"http://example.com/a.preview.png"}]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	msg, err := client.Channels.ReceiveMessage(
		context.Background(),
		WithChannelMessageID("msg-1"),
		WithChannelMessageChannelID("chan-1"),
		WithChannelMessageSenderID("user-1"),
		WithChannelMessageConversationID("conv-1"),
		WithChannelMessageBody("Hello"),
		WithChannelMessageStatus(ChannelMessageStatusSent),
		WithChannelMessageCreatedAt(created),
	)
	if err != nil {
		t.Fatalf("ReceiveMessage error: %v", err)
	}
	if msg.ID != "msg-1" || msg.Status != ChannelMessageStatusSent {
		t.Fatalf("unexpected message: %#v", msg)
	}
	if msg.CreatedAt == nil || msg.CreatedAt.IsZero() {
		t.Fatalf("expected created_at to be parsed")
	}
	if len(msg.Attachments) != 1 || msg.Attachments[0].ID != "att-1" {
		t.Fatalf("unexpected attachments: %#v", msg.Attachments)
	}
}

func TestChannelsService_DeleteConversation(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/channels/chan-1/conversations/conv-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ok, err := client.Channels.DeleteConversation(context.Background(), ChannelID("chan-1"), ConversationID("conv-1"))
	if err != nil {
		t.Fatalf("DeleteConversation error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}
