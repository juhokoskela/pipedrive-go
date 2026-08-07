package v2

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
		if err := validatePathParam(value, "field code"); err == nil {
			t.Fatalf("expected error for %q", value)
		}
	}
	if err := validatePathParam("deal_field_key", "field code"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFieldCode_DotSegmentsRejectedClientSide(t *testing.T) {
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
		{"ActivityFields.Get", func() error { _, err := client.ActivityFields.Get(ctx, ".."); return err }},
		{"DealFields.Get", func() error { _, err := client.DealFields.Get(ctx, ".."); return err }},
		{"DealFields.Update", func() error { _, err := client.DealFields.Update(ctx, ".."); return err }},
		{"DealFields.Delete", func() error { _, err := client.DealFields.Delete(ctx, ".."); return err }},
		{"DealFields.AddOptions", func() error { _, err := client.DealFields.AddOptions(ctx, "..", nil); return err }},
		{"DealFields.UpdateOptions", func() error { _, err := client.DealFields.UpdateOptions(ctx, "..", nil); return err }},
		{"DealFields.DeleteOptions", func() error { _, err := client.DealFields.DeleteOptions(ctx, "..", nil); return err }},
		{"OrganizationFields.Get", func() error { _, err := client.OrganizationFields.Get(ctx, ".."); return err }},
		{"OrganizationFields.Update", func() error { _, err := client.OrganizationFields.Update(ctx, ".."); return err }},
		{"OrganizationFields.Delete", func() error { _, err := client.OrganizationFields.Delete(ctx, ".."); return err }},
		{"OrganizationFields.AddOptions", func() error { _, err := client.OrganizationFields.AddOptions(ctx, "..", nil); return err }},
		{"OrganizationFields.UpdateOptions", func() error { _, err := client.OrganizationFields.UpdateOptions(ctx, "..", nil); return err }},
		{"OrganizationFields.DeleteOptions", func() error { _, err := client.OrganizationFields.DeleteOptions(ctx, "..", nil); return err }},
		{"PersonFields.Get", func() error { _, err := client.PersonFields.Get(ctx, ".."); return err }},
		{"PersonFields.Update", func() error { _, err := client.PersonFields.Update(ctx, ".."); return err }},
		{"PersonFields.Delete", func() error { _, err := client.PersonFields.Delete(ctx, ".."); return err }},
		{"PersonFields.AddOptions", func() error { _, err := client.PersonFields.AddOptions(ctx, "..", nil); return err }},
		{"PersonFields.UpdateOptions", func() error { _, err := client.PersonFields.UpdateOptions(ctx, "..", nil); return err }},
		{"PersonFields.DeleteOptions", func() error { _, err := client.PersonFields.DeleteOptions(ctx, "..", nil); return err }},
		{"ProductFields.Get", func() error { _, err := client.ProductFields.Get(ctx, ".."); return err }},
		{"ProductFields.Update", func() error { _, err := client.ProductFields.Update(ctx, ".."); return err }},
		{"ProductFields.Delete", func() error { _, err := client.ProductFields.Delete(ctx, ".."); return err }},
		{"ProductFields.AddOptions", func() error { _, err := client.ProductFields.AddOptions(ctx, "..", nil); return err }},
		{"ProductFields.UpdateOptions", func() error { _, err := client.ProductFields.UpdateOptions(ctx, "..", nil); return err }},
		{"ProductFields.DeleteOptions", func() error { _, err := client.ProductFields.DeleteOptions(ctx, "..", nil); return err }},
		{"ProjectFields.Get", func() error { _, err := client.ProjectFields.Get(ctx, ".."); return err }},
		{"ProjectFields.Update", func() error { _, err := client.ProjectFields.Update(ctx, ".."); return err }},
		{"ProjectFields.Delete", func() error { _, err := client.ProjectFields.Delete(ctx, ".."); return err }},
		{"ProjectFields.AddOptions", func() error { _, err := client.ProjectFields.AddOptions(ctx, "..", nil); return err }},
		{"ProjectFields.UpdateOptions", func() error { _, err := client.ProjectFields.UpdateOptions(ctx, "..", nil); return err }},
		{"ProjectFields.DeleteOptions", func() error { _, err := client.ProjectFields.DeleteOptions(ctx, "..", nil); return err }},
	}
	for _, c := range calls {
		err := c.call()
		if err == nil {
			t.Errorf("%s: expected client-side validation error", c.name)
			continue
		}
		if !strings.Contains(err.Error(), "invalid") {
			t.Errorf("%s: expected path param validation error, got %v", c.name, err)
		}
	}
}
