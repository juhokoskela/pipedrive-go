package v1

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/juhokoskela/pipedrive-go/v2/pipedrive"
)

func TestScalarQueryIDsRejectOverflowBeforeRequest(t *testing.T) {
	if strconv.IntSize != 32 {
		t.Skip("requires 32-bit int")
	}
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("unexpected HTTP request for overflowing ID")
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()
	c, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	const large = int64(math.MaxInt32) + 1
	tests := []struct {
		name string
		call func() error
	}{
		{"WithGoalsExpectedOutcomeCurrencyID", func() error {
			_, err := c.Goals.List(ctx, WithGoalsExpectedOutcomeCurrencyID(CurrencyID(large)))
			return err
		}},
		{"WithGoalsTypeStageID", func() error { _, err := c.Goals.List(ctx, WithGoalsTypeStageID(StageID(large))); return err }},
		{"WithLeadsOwnerID", func() error { _, _, err := c.Leads.List(ctx, WithLeadsOwnerID(UserID(large))); return err }},
		{"WithLeadsPersonID", func() error { _, _, err := c.Leads.List(ctx, WithLeadsPersonID(PersonID(large))); return err }},
		{"WithLeadsOrganizationID", func() error {
			_, _, err := c.Leads.List(ctx, WithLeadsOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithArchivedLeadsOwnerID", func() error {
			_, _, err := c.Leads.ListArchived(ctx, WithArchivedLeadsOwnerID(UserID(large)))
			return err
		}},
		{"WithArchivedLeadsPersonID", func() error {
			_, _, err := c.Leads.ListArchived(ctx, WithArchivedLeadsPersonID(PersonID(large)))
			return err
		}},
		{"WithArchivedLeadsOrganizationID", func() error {
			_, _, err := c.Leads.ListArchived(ctx, WithArchivedLeadsOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithNotesUserID", func() error { _, _, err := c.Notes.List(ctx, WithNotesUserID(UserID(large))); return err }},
		{"WithNotesDealID", func() error { _, _, err := c.Notes.List(ctx, WithNotesDealID(DealID(large))); return err }},
		{"WithNotesPersonID", func() error { _, _, err := c.Notes.List(ctx, WithNotesPersonID(PersonID(large))); return err }},
		{"WithNotesOrganizationID", func() error {
			_, _, err := c.Notes.List(ctx, WithNotesOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithNotesProjectID", func() error { _, _, err := c.Notes.List(ctx, WithNotesProjectID(ProjectID(large))); return err }},
		{"WithNotesTaskID", func() error { _, _, err := c.Notes.List(ctx, WithNotesTaskID(TaskID(large))); return err }},
		{"WithPipelineConversionUserID", func() error {
			_, err := c.Pipelines.GetConversionStatistics(ctx, 1, WithPipelineConversionUserID(UserID(large)))
			return err
		}},
		{"WithPipelineMovementUserID", func() error {
			_, err := c.Pipelines.GetMovementStatistics(ctx, 1, WithPipelineMovementUserID(UserID(large)))
			return err
		}},
		{"WithGoalsTypePipelineIDs", func() error { _, err := c.Goals.List(ctx, WithGoalsTypePipelineIDs(PipelineID(large))); return err }},
		{"WithGoalsTypeActivityTypeIDs", func() error {
			_, err := c.Goals.List(ctx, WithGoalsTypeActivityTypeIDs(ActivityTypeID(large)))
			return err
		}},
		{"relationship list", func() error {
			_, _, err := c.OrganizationRelationships.List(ctx, WithOrganizationRelationshipsOrgID(OrganizationID(large)))
			return err
		}},
		{"relationship get", func() error {
			_, err := c.OrganizationRelationships.Get(ctx, 1, WithOrganizationRelationshipOrgID(OrganizationID(large)))
			return err
		}},
		{"meeting user", func() error {
			_, err := c.Meetings.CreateUserProviderLink(ctx, WithUserProviderLinkID("1e3943c9-6395-462b-b432-1f252c017f3d"), WithUserProviderLinkUserID(UserID(large)), WithUserProviderLinkCompanyID(1), WithUserProviderLinkMarketplaceClientID("test"))
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); err == nil || !strings.Contains(err.Error(), "overflows int") {
				t.Fatalf("expected overflow error, got %v", err)
			}
		})
	}
}
