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
		{"DealFields.Get", func() error { _, err := client.DealFields.Get(ctx, ".."); return err }},
		{"DealFields.Delete", func() error { _, err := client.DealFields.Delete(ctx, "."); return err }},
		{"OrganizationFields.Get", func() error { _, err := client.OrganizationFields.Get(ctx, ""); return err }},
		{"PersonFields.Get", func() error { _, err := client.PersonFields.Get(ctx, ".."); return err }},
		{"ProductFields.Get", func() error { _, err := client.ProductFields.Get(ctx, ".."); return err }},
		{"ProjectFields.Get", func() error { _, err := client.ProjectFields.Get(ctx, ".."); return err }},
		{"ActivityFields.Get", func() error { _, err := client.ActivityFields.Get(ctx, ".."); return err }},
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
