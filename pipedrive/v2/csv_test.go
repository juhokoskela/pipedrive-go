package v2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestValidateCSVValues(t *testing.T) {
	t.Parallel()

	if err := validateCSVValues([]string{"a", "b,c"}, "custom field key"); err == nil {
		t.Fatal("expected error for value containing a comma")
	}
	if err := validateCSVValues([]string{"a", "b"}, "custom field key"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCustomFieldsWithCommaRejected(t *testing.T) {
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
		{"Deals.List", func() error { _, _, err := client.Deals.List(ctx, WithDealsCustomFields("a,b")); return err }},
		{"Deals.Get", func() error { _, err := client.Deals.Get(ctx, 1, WithDealCustomFields("a,b")); return err }},
		{"Deals.ListArchived", func() error {
			_, _, err := client.Deals.ListArchived(ctx, WithArchivedDealsCustomFields("a,b"))
			return err
		}},
		{"Organizations.List", func() error {
			_, _, err := client.Organizations.List(ctx, WithOrganizationsCustomFields("a,b"))
			return err
		}},
		{"Organizations.Get", func() error {
			_, err := client.Organizations.Get(ctx, 1, WithOrganizationCustomFields("a,b"))
			return err
		}},
		{"Persons.List", func() error { _, _, err := client.Persons.List(ctx, WithPersonsCustomFields("a,b")); return err }},
		{"Persons.Get", func() error { _, err := client.Persons.Get(ctx, 1, WithPersonCustomFields("a,b")); return err }},
		{"Products.List", func() error { _, _, err := client.Products.List(ctx, WithProductsCustomFields("a,b")); return err }},
	}
	for _, c := range calls {
		err := c.call()
		if err == nil {
			t.Errorf("%s: expected comma rejection", c.name)
			continue
		}
		if !strings.Contains(err.Error(), "comma") {
			t.Errorf("%s: expected comma validation error, got %v", c.name, err)
		}
	}

	// Pagers surface the same error on first fetch.
	pager := client.Deals.ListPager(WithDealsCustomFields("a,b"))
	if pager.Next(ctx) {
		t.Error("expected pager fetch to fail")
	}
	if pager.Err() == nil {
		t.Error("expected pager to surface comma rejection")
	}
}

func TestEmptyCustomFieldsMapDoesNotClearFields(t *testing.T) {
	t.Parallel()

	var dealCfg createDealOptions
	WithDealCustomFieldsMap(map[string]interface{}{}).applyCreateDeal(&dealCfg)
	if _, ok := dealCfg.payload.toMap()["custom_fields"]; ok {
		t.Error("deal: empty custom fields map must not emit custom_fields")
	}

	var orgCfg createOrganizationOptions
	WithOrganizationCustomFieldsMap(map[string]interface{}{}).applyCreateOrganization(&orgCfg)
	if _, ok := orgCfg.payload.toMap()["custom_fields"]; ok {
		t.Error("organization: empty custom fields map must not emit custom_fields")
	}

	var projectCfg createProjectOptions
	WithProjectCustomFields(map[string]interface{}{}).applyCreateProject(&projectCfg)
	if _, ok := projectCfg.payload.body()["custom_fields"]; ok {
		t.Error("project: empty custom fields map must not emit custom_fields")
	}

	// A non-empty map still round-trips.
	var populated createDealOptions
	WithDealCustomFieldsMap(map[string]interface{}{"key": "value"}).applyCreateDeal(&populated)
	if _, ok := populated.payload.toMap()["custom_fields"]; !ok {
		t.Error("deal: non-empty custom fields map must emit custom_fields")
	}
}

func TestPagerVariantsRejectEmptyDealIDs(t *testing.T) {
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

	products := client.Deals.ListProductsAcrossDealsPager(nil)
	if products.Next(ctx) || products.Err() == nil {
		t.Error("expected ListProductsAcrossDealsPager to reject empty deal IDs")
	}
	installments := client.Deals.ListInstallmentsPager(nil)
	if installments.Next(ctx) || installments.Err() == nil {
		t.Error("expected ListInstallmentsPager to reject empty deal IDs")
	}
}
