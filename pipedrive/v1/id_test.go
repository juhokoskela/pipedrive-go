package v1

import (
	"context"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestValidateID(t *testing.T) {
	t.Parallel()

	for _, id := range []NoteID{0, -1, math.MinInt64} {
		if err := validateID(id, "note id"); err == nil {
			t.Fatalf("expected error for id %d", int64(id))
		}
	}
	if err := validateID(NoteID(1), "note id"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := validateID(NoteID(math.MaxInt64), "note id")
	if math.MaxInt == math.MaxInt64 && err != nil {
		t.Fatalf("unexpected error on 64-bit platform: %v", err)
	}
	if math.MaxInt != math.MaxInt64 && err == nil {
		t.Fatal("expected overflow error where int is 32 bits")
	}
}

func TestNonPositiveIDsRejectedClientSide(t *testing.T) {
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
		{"Notes.Get", func() error { _, err := client.Notes.Get(ctx, 0); return err }},
		{"Notes.Delete", func() error { _, err := client.Notes.Delete(ctx, -1); return err }},
		{"Notes.ListComments", func() error { _, _, err := client.Notes.ListComments(ctx, 0); return err }},
		{"Files.Download", func() error { _, err := client.Files.Download(ctx, 0); return err }},
		{"Files.DownloadTo", func() error { return client.Files.DownloadTo(ctx, 0, io.Discard) }},
		{"Filters.Get", func() error { _, err := client.Filters.Get(ctx, 0); return err }},
		{"Users.Get", func() error { _, err := client.Users.Get(ctx, -5); return err }},
		{"Webhooks.Delete", func() error { _, err := client.Webhooks.Delete(ctx, 0); return err }},
		{"ActivityTypes.Delete", func() error { _, err := client.ActivityTypes.Delete(ctx, 0); return err }},
		{"Pipelines.GetConversionStatistics", func() error {
			_, err := client.Pipelines.GetConversionStatistics(ctx, 0)
			return err
		}},
		{"OrganizationRelationships.Get", func() error {
			_, err := client.OrganizationRelationships.Get(ctx, 0)
			return err
		}},
		{"Deals.Changelog", func() error { _, _, err := client.Deals.Changelog(ctx, 0); return err }},
		{"Deals.Merge", func() error { _, err := client.Deals.Merge(ctx, 0, 1); return err }},
		{"Deals.DeleteParticipant", func() error { _, err := client.Deals.DeleteParticipant(ctx, 1, 0); return err }},
		{"Organizations.Merge", func() error { _, err := client.Organizations.Merge(ctx, 0, 1); return err }},
		{"Persons.Merge", func() error { _, err := client.Persons.Merge(ctx, 0, 1); return err }},
		{"Persons.DeletePicture", func() error { _, err := client.Persons.DeletePicture(ctx, -2); return err }},
		{"Projects.Get", func() error { _, err := client.Projects.Get(ctx, 0); return err }},
		{"Projects.UpdatePlanActivity", func() error {
			_, err := client.Projects.UpdatePlanActivity(ctx, 1, 0, nil)
			return err
		}},
		{"Roles.Get", func() error { _, err := client.Roles.Get(ctx, 0); return err }},
		{"Teams.Get", func() error { _, err := client.Teams.Get(ctx, 0); return err }},
		{"Teams.AddUsers", func() error { _, err := client.Teams.AddUsers(ctx, 0, []UserID{1}); return err }},
		{"Mailbox.GetThread", func() error { _, err := client.Mailbox.GetThread(ctx, 0); return err }},
		{"Tasks.Get", func() error { _, err := client.Tasks.Get(ctx, 0); return err }},
		{"Stages.ListDeals", func() error { _, _, err := client.Stages.ListDeals(ctx, 0); return err }},
		{"Users.Update", func() error { _, err := client.Users.Update(ctx, 0, nil); return err }},
		{"Products.ListDeals", func() error { _, _, err := client.Products.ListDeals(ctx, 0); return err }},
		{"ProjectTemplates.Get", func() error { _, err := client.ProjectTemplates.Get(ctx, 0); return err }},
		{"Files.Get", func() error { _, err := client.Files.Get(ctx, 0); return err }},
		{"Pipelines.ListDeals", func() error { _, _, err := client.Pipelines.ListDeals(ctx, 0); return err }},
	}
	for _, c := range calls {
		err := c.call()
		if err == nil {
			t.Errorf("%s: expected client-side validation error", c.name)
			continue
		}
		if !strings.Contains(err.Error(), "invalid") {
			t.Errorf("%s: expected id validation error, got %v", c.name, err)
		}
	}
}
