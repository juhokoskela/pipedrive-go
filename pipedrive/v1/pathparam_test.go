package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestValidatePathParam(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"", ".", ".."} {
		if err := validatePathParam(value, "id"); err == nil {
			t.Fatalf("expected error for %q", value)
		}
	}
	if err := validatePathParam("abc-123", "id"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStringIDs_DotSegmentsRejectedClientSide(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("no request should reach the server, got %s %s", r.Method, r.URL)
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	ctx := context.Background()

	calls := []struct {
		name string
		call func() error
	}{
		{"CallLogs.Get", func() error { _, err := client.CallLogs.Get(ctx, ".."); return err }},
		{"CallLogs.Delete", func() error { _, err := client.CallLogs.Delete(ctx, ".."); return err }},
		{"CallLogs.AddRecording", func() error { _, err := client.CallLogs.AddRecording(ctx, "..", "x", nil); return err }},
		{"Channels.Delete", func() error { _, err := client.Channels.Delete(ctx, ".."); return err }},
		{"Channels.DeleteConversation#channelID", func() error { _, err := client.Channels.DeleteConversation(ctx, "..", "x"); return err }},
		{"Channels.DeleteConversation#conversationID", func() error { _, err := client.Channels.DeleteConversation(ctx, "x", ".."); return err }},
		{"Goals.Update", func() error { _, err := client.Goals.Update(ctx, ".."); return err }},
		{"Goals.Delete", func() error { _, err := client.Goals.Delete(ctx, ".."); return err }},
		{"Goals.GetResult", func() error { _, err := client.Goals.GetResult(ctx, ".."); return err }},
		{"PermissionSets.Get", func() error { _, err := client.PermissionSets.Get(ctx, ".."); return err }},
		{"PermissionSets.ListAssignments", func() error { _, err := client.PermissionSets.ListAssignments(ctx, ".."); return err }},
		{"Leads.ListPermittedUsers", func() error { _, err := client.Leads.ListPermittedUsers(ctx, "not-a-uuid"); return err }},
	}
	for _, c := range calls {
		err := c.call()
		if err == nil {
			t.Errorf("%s: expected client-side validation error", c.name)
			continue
		}
		if c.name != "Leads.ListPermittedUsers" && !strings.Contains(err.Error(), "invalid") {
			t.Errorf("%s: expected path param validation error, got %v", c.name, err)
		}
	}
}
