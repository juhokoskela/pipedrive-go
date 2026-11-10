package v2

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
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
		{"WithActivitiesDealID", func() error { _, _, err := c.Activities.List(ctx, WithActivitiesDealID(DealID(large))); return err }},
		{"WithActivitiesPersonID", func() error { _, _, err := c.Activities.List(ctx, WithActivitiesPersonID(PersonID(large))); return err }},
		{"WithActivitiesOrgID", func() error {
			_, _, err := c.Activities.List(ctx, WithActivitiesOrgID(OrganizationID(large)))
			return err
		}},
		{"WithDealsOwnerID", func() error { _, _, err := c.Deals.List(ctx, WithDealsOwnerID(UserID(large))); return err }},
		{"WithDealsPersonID", func() error { _, _, err := c.Deals.List(ctx, WithDealsPersonID(PersonID(large))); return err }},
		{"WithDealsOrganizationID", func() error {
			_, _, err := c.Deals.List(ctx, WithDealsOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithDealsPipelineID", func() error { _, _, err := c.Deals.List(ctx, WithDealsPipelineID(PipelineID(large))); return err }},
		{"WithDealsStageID", func() error { _, _, err := c.Deals.List(ctx, WithDealsStageID(StageID(large))); return err }},
		{"WithArchivedDealsOwnerID", func() error {
			_, _, err := c.Deals.ListArchived(ctx, WithArchivedDealsOwnerID(UserID(large)))
			return err
		}},
		{"WithArchivedDealsPersonID", func() error {
			_, _, err := c.Deals.ListArchived(ctx, WithArchivedDealsPersonID(PersonID(large)))
			return err
		}},
		{"WithArchivedDealsOrganizationID", func() error {
			_, _, err := c.Deals.ListArchived(ctx, WithArchivedDealsOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithArchivedDealsPipelineID", func() error {
			_, _, err := c.Deals.ListArchived(ctx, WithArchivedDealsPipelineID(PipelineID(large)))
			return err
		}},
		{"WithArchivedDealsStageID", func() error {
			_, _, err := c.Deals.ListArchived(ctx, WithArchivedDealsStageID(StageID(large)))
			return err
		}},
		{"WithDealSearchPersonID", func() error {
			_, _, err := c.Deals.Search(ctx, "term", WithDealSearchPersonID(PersonID(large)))
			return err
		}},
		{"WithDealSearchOrganizationID", func() error {
			_, _, err := c.Deals.Search(ctx, "term", WithDealSearchOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithLeadSearchPersonID", func() error {
			_, _, err := c.Leads.Search(ctx, "term", WithLeadSearchPersonID(PersonID(large)))
			return err
		}},
		{"WithLeadSearchOrganizationID", func() error {
			_, _, err := c.Leads.Search(ctx, "term", WithLeadSearchOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithPersonsOrgID", func() error { _, _, err := c.Persons.List(ctx, WithPersonsOrgID(OrganizationID(large))); return err }},
		{"WithPersonsDealID", func() error { _, _, err := c.Persons.List(ctx, WithPersonsDealID(DealID(large))); return err }},
		{"WithPersonSearchOrganizationID", func() error {
			_, _, err := c.Persons.Search(ctx, "term", WithPersonSearchOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithProductsOwnerID", func() error { _, _, err := c.Products.List(ctx, WithProductsOwnerID(UserID(large))); return err }},
		{"WithProjectsDealID", func() error { _, _, err := c.Projects.List(ctx, WithProjectsDealID(DealID(large))); return err }},
		{"WithProjectsPersonID", func() error { _, _, err := c.Projects.List(ctx, WithProjectsPersonID(PersonID(large))); return err }},
		{"WithProjectsOrganizationID", func() error {
			_, _, err := c.Projects.List(ctx, WithProjectsOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithProjectSearchPersonID", func() error {
			_, _, err := c.Projects.Search(ctx, "term", WithProjectSearchPersonID(PersonID(large)))
			return err
		}},
		{"WithProjectSearchOrganizationID", func() error {
			_, _, err := c.Projects.Search(ctx, "term", WithProjectSearchOrganizationID(OrganizationID(large)))
			return err
		}},
		{"WithStagesPipelineID", func() error { _, _, err := c.Stages.List(ctx, WithStagesPipelineID(PipelineID(large))); return err }},
		{"WithTasksAssigneeID", func() error { _, _, err := c.Tasks.List(ctx, WithTasksAssigneeID(UserID(large))); return err }},
		{"WithTasksProjectID", func() error { _, _, err := c.Tasks.List(ctx, WithTasksProjectID(ProjectID(large))); return err }},
		{"WithProjectsPhaseID/List", func() error {
			_, _, err := c.Projects.List(ctx, WithProjectsPhaseID(ProjectPhaseID(large)))
			return err
		}},
		{"WithProjectsPhaseID/ListArchived", func() error {
			_, _, err := c.Projects.ListArchived(ctx, WithProjectsPhaseID(ProjectPhaseID(large)))
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
