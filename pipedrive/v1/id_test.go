package v1

import (
	"context"
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
		{"ActivityTypes.Update", func() error { _, err := client.ActivityTypes.Update(ctx, 0); return err }},
		{"ActivityTypes.Delete", func() error { _, err := client.ActivityTypes.Delete(ctx, 0); return err }},
		{"Deals.Changelog", func() error { _, _, err := client.Deals.Changelog(ctx, 0); return err }},
		{"Deals.ListFiles", func() error { _, _, err := client.Deals.ListFiles(ctx, 0); return err }},
		{"Deals.ListMailMessages", func() error { _, _, err := client.Deals.ListMailMessages(ctx, 0); return err }},
		{"Deals.ListParticipants", func() error { _, _, err := client.Deals.ListParticipants(ctx, 0); return err }},
		{"Deals.AddParticipant", func() error { _, err := client.Deals.AddParticipant(ctx, 0, 1); return err }},
		{"Deals.DeleteParticipant#id", func() error { _, err := client.Deals.DeleteParticipant(ctx, 0, 1); return err }},
		{"Deals.DeleteParticipant#participantID", func() error { _, err := client.Deals.DeleteParticipant(ctx, 1, 0); return err }},
		{"Deals.ParticipantsChangelog", func() error { _, _, err := client.Deals.ParticipantsChangelog(ctx, 0); return err }},
		{"Deals.ListUpdates", func() error { _, _, err := client.Deals.ListUpdates(ctx, 0); return err }},
		{"Deals.ListUsers", func() error { _, err := client.Deals.ListUsers(ctx, 0); return err }},
		{"Deals.Merge", func() error { _, err := client.Deals.Merge(ctx, 0, 1); return err }},
		{"Deals.Duplicate", func() error { _, err := client.Deals.Duplicate(ctx, 0); return err }},
		{"Files.Get", func() error { _, err := client.Files.Get(ctx, 0); return err }},
		{"Files.Update", func() error { _, err := client.Files.Update(ctx, 0, nil, ""); return err }},
		{"Files.Delete", func() error { _, err := client.Files.Delete(ctx, 0); return err }},
		{"Files.Download", func() error { _, err := client.Files.Download(ctx, 0); return err }},
		{"Files.DownloadTo", func() error { err := client.Files.DownloadTo(ctx, 0, nil); return err }},
		{"Filters.Get", func() error { _, err := client.Filters.Get(ctx, 0); return err }},
		{"Filters.Update", func() error { _, err := client.Filters.Update(ctx, 0); return err }},
		{"Filters.Delete", func() error { _, err := client.Filters.Delete(ctx, 0); return err }},
		{"Mailbox.GetThread", func() error { _, err := client.Mailbox.GetThread(ctx, 0); return err }},
		{"Mailbox.DeleteThread", func() error { _, err := client.Mailbox.DeleteThread(ctx, 0); return err }},
		{"Mailbox.UpdateThread", func() error { _, err := client.Mailbox.UpdateThread(ctx, 0, nil); return err }},
		{"Mailbox.ListThreadMessages", func() error { _, _, err := client.Mailbox.ListThreadMessages(ctx, 0); return err }},
		{"Mailbox.GetMessage", func() error { _, err := client.Mailbox.GetMessage(ctx, 0); return err }},
		{"Notes.Get", func() error { _, err := client.Notes.Get(ctx, 0); return err }},
		{"Notes.Update", func() error { _, err := client.Notes.Update(ctx, 0); return err }},
		{"Notes.Delete", func() error { _, err := client.Notes.Delete(ctx, 0); return err }},
		{"Notes.ListComments", func() error { _, _, err := client.Notes.ListComments(ctx, 0); return err }},
		{"Notes.CreateComment", func() error { _, err := client.Notes.CreateComment(ctx, 0); return err }},
		{"Notes.GetComment", func() error { _, err := client.Notes.GetComment(ctx, 0, "x"); return err }},
		{"Notes.UpdateComment", func() error { _, err := client.Notes.UpdateComment(ctx, 0, "x"); return err }},
		{"Notes.DeleteComment", func() error { _, err := client.Notes.DeleteComment(ctx, 0, "x"); return err }},
		{"OrganizationRelationships.Get", func() error { _, err := client.OrganizationRelationships.Get(ctx, 0); return err }},
		{"OrganizationRelationships.Update", func() error { _, err := client.OrganizationRelationships.Update(ctx, 0); return err }},
		{"OrganizationRelationships.Delete", func() error { _, err := client.OrganizationRelationships.Delete(ctx, 0); return err }},
		{"Organizations.Merge", func() error { _, err := client.Organizations.Merge(ctx, 0, 1); return err }},
		{"Organizations.Changelog", func() error { _, _, err := client.Organizations.Changelog(ctx, 0); return err }},
		{"Organizations.ListFiles", func() error { _, _, err := client.Organizations.ListFiles(ctx, 0); return err }},
		{"Organizations.ListMailMessages", func() error { _, _, err := client.Organizations.ListMailMessages(ctx, 0); return err }},
		{"Organizations.ListUpdates", func() error { _, _, err := client.Organizations.ListUpdates(ctx, 0); return err }},
		{"Organizations.ListUsers", func() error { _, err := client.Organizations.ListUsers(ctx, 0); return err }},
		{"Persons.Merge", func() error { _, err := client.Persons.Merge(ctx, 0, 1); return err }},
		{"Persons.Changelog", func() error { _, _, err := client.Persons.Changelog(ctx, 0); return err }},
		{"Persons.ListFiles", func() error { _, _, err := client.Persons.ListFiles(ctx, 0); return err }},
		{"Persons.ListMailMessages", func() error { _, _, err := client.Persons.ListMailMessages(ctx, 0); return err }},
		{"Persons.ListProducts", func() error { _, _, err := client.Persons.ListProducts(ctx, 0); return err }},
		{"Persons.ListUpdates", func() error { _, _, err := client.Persons.ListUpdates(ctx, 0); return err }},
		{"Persons.ListUsers", func() error { _, err := client.Persons.ListUsers(ctx, 0); return err }},
		{"Persons.AddPicture", func() error { _, err := client.Persons.AddPicture(ctx, 0, nil, ""); return err }},
		{"Persons.DeletePicture", func() error { _, err := client.Persons.DeletePicture(ctx, 0); return err }},
		{"Pipelines.GetConversionStatistics", func() error { _, err := client.Pipelines.GetConversionStatistics(ctx, 0); return err }},
		{"Pipelines.GetMovementStatistics", func() error { _, err := client.Pipelines.GetMovementStatistics(ctx, 0); return err }},
		{"Pipelines.ListDeals", func() error { _, _, err := client.Pipelines.ListDeals(ctx, 0); return err }},
		{"Products.ListDeals", func() error { _, _, err := client.Products.ListDeals(ctx, 0); return err }},
		{"Products.ListFiles", func() error { _, _, err := client.Products.ListFiles(ctx, 0); return err }},
		{"Products.ListUsers", func() error { _, err := client.Products.ListUsers(ctx, 0); return err }},
		{"ProjectTemplates.Get", func() error { _, err := client.ProjectTemplates.Get(ctx, 0); return err }},
		{"Projects.Get", func() error { _, err := client.Projects.Get(ctx, 0); return err }},
		{"Projects.Update", func() error { _, err := client.Projects.Update(ctx, 0, nil); return err }},
		{"Projects.Delete", func() error { _, err := client.Projects.Delete(ctx, 0); return err }},
		{"Projects.Archive", func() error { _, err := client.Projects.Archive(ctx, 0); return err }},
		{"Projects.GetBoard", func() error { _, err := client.Projects.GetBoard(ctx, 0); return err }},
		{"Projects.GetPhase", func() error { _, err := client.Projects.GetPhase(ctx, 0); return err }},
		{"Projects.ListActivities", func() error { _, _, err := client.Projects.ListActivities(ctx, 0); return err }},
		{"Projects.ListGroups", func() error { _, err := client.Projects.ListGroups(ctx, 0); return err }},
		{"Projects.GetPlan", func() error { _, err := client.Projects.GetPlan(ctx, 0); return err }},
		{"Projects.UpdatePlanActivity#id", func() error { _, err := client.Projects.UpdatePlanActivity(ctx, 0, 1, nil); return err }},
		{"Projects.UpdatePlanActivity#activityID", func() error { _, err := client.Projects.UpdatePlanActivity(ctx, 1, 0, nil); return err }},
		{"Projects.UpdatePlanTask#id", func() error { _, err := client.Projects.UpdatePlanTask(ctx, 0, 1, nil); return err }},
		{"Projects.UpdatePlanTask#taskID", func() error { _, err := client.Projects.UpdatePlanTask(ctx, 1, 0, nil); return err }},
		{"Projects.ListTasks", func() error { _, err := client.Projects.ListTasks(ctx, 0); return err }},
		{"Roles.Get", func() error { _, err := client.Roles.Get(ctx, 0); return err }},
		{"Roles.Update", func() error { _, err := client.Roles.Update(ctx, 0, nil); return err }},
		{"Roles.Delete", func() error { _, err := client.Roles.Delete(ctx, 0); return err }},
		{"Roles.ListAssignments", func() error { _, _, err := client.Roles.ListAssignments(ctx, 0); return err }},
		{"Roles.AddAssignment", func() error { _, err := client.Roles.AddAssignment(ctx, 0, 1); return err }},
		{"Roles.DeleteAssignment", func() error { _, err := client.Roles.DeleteAssignment(ctx, 0, 1); return err }},
		{"Roles.ListPipelines", func() error { _, err := client.Roles.ListPipelines(ctx, 0); return err }},
		{"Roles.UpdatePipelines", func() error { _, err := client.Roles.UpdatePipelines(ctx, 0, nil); return err }},
		{"Roles.ListSettings", func() error { _, err := client.Roles.ListSettings(ctx, 0); return err }},
		{"Roles.UpsertSetting", func() error { _, err := client.Roles.UpsertSetting(ctx, 0, nil); return err }},
		{"Stages.ListDeals", func() error { _, _, err := client.Stages.ListDeals(ctx, 0); return err }},
		{"Tasks.Get", func() error { _, err := client.Tasks.Get(ctx, 0); return err }},
		{"Tasks.Update", func() error { _, err := client.Tasks.Update(ctx, 0, nil); return err }},
		{"Tasks.Delete", func() error { _, err := client.Tasks.Delete(ctx, 0); return err }},
		{"Teams.Get", func() error { _, err := client.Teams.Get(ctx, 0); return err }},
		{"Teams.Update", func() error { _, err := client.Teams.Update(ctx, 0, nil); return err }},
		{"Teams.ListUsers", func() error { _, err := client.Teams.ListUsers(ctx, 0); return err }},
		{"Teams.AddUsers", func() error { _, err := client.Teams.AddUsers(ctx, 0, nil); return err }},
		{"Teams.DeleteUsers", func() error { _, err := client.Teams.DeleteUsers(ctx, 0, nil); return err }},
		{"Users.Get", func() error { _, err := client.Users.Get(ctx, 0); return err }},
		{"Users.GetPermissions", func() error { _, err := client.Users.GetPermissions(ctx, 0); return err }},
		{"Users.Update", func() error { _, err := client.Users.Update(ctx, 0, nil); return err }},
		{"Users.ListRoleAssignments", func() error { _, err := client.Users.ListRoleAssignments(ctx, 0, nil); return err }},
		{"Users.ListRoleSettings", func() error { _, err := client.Users.ListRoleSettings(ctx, 0); return err }},
		{"Users.ListTeams", func() error { _, err := client.Users.ListTeams(ctx, 0, nil); return err }},
		{"Webhooks.Delete", func() error { _, err := client.Webhooks.Delete(ctx, 0); return err }},
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
