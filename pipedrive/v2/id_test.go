package v2

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

	for _, id := range []DealID{0, -1, math.MinInt64} {
		if err := validateID(id, "deal id"); err == nil {
			t.Fatalf("expected error for id %d", int64(id))
		}
	}
	if err := validateID(DealID(1), "deal id"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Representable on 64-bit platforms, rejected where int is 32 bits.
	err := validateID(DealID(math.MaxInt64), "deal id")
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
		{"Deals.Get", func() error { _, err := client.Deals.Get(ctx, 0); return err }},
		{"Deals.Delete", func() error { _, err := client.Deals.Delete(ctx, -1); return err }},
		{"Deals.DeleteFollower", func() error { _, err := client.Deals.DeleteFollower(ctx, 1, 0); return err }},
		{"Deals.DeleteProduct", func() error { _, err := client.Deals.DeleteProduct(ctx, 1, -2); return err }},
		{"Persons.Get", func() error { _, err := client.Persons.Get(ctx, 0); return err }},
		{"Organizations.Get", func() error { _, err := client.Organizations.Get(ctx, -5); return err }},
		{"Products.Get", func() error { _, err := client.Products.Get(ctx, 0); return err }},
		{"Activities.Get", func() error { _, err := client.Activities.Get(ctx, 0); return err }},
		{"Projects.Get", func() error { _, err := client.Projects.Get(ctx, 0); return err }},
		{"Tasks.Get", func() error { _, err := client.Tasks.Get(ctx, 0); return err }},
		{"Stages.Get", func() error { _, err := client.Stages.Get(ctx, 0); return err }},
		{"Pipelines.Get", func() error { _, err := client.Pipelines.Get(ctx, 0); return err }},
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

func TestPagerIDsValidatedClientSide(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("no request should reach the server, got %s %s", r.Method, r.URL)
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if _, _, err := client.Deals.ListFollowers(context.Background(), 0); err == nil {
		t.Fatal("expected client-side validation error from ListFollowers")
	}
}
