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
		{"Activities.Get", func() error { _, err := client.Activities.Get(ctx, 0); return err }},
		{"Activities.Update", func() error { _, err := client.Activities.Update(ctx, 0); return err }},
		{"Activities.Delete", func() error { _, err := client.Activities.Delete(ctx, 0); return err }},
		{"Deals.Get", func() error { _, err := client.Deals.Get(ctx, 0); return err }},
		{"Deals.Update", func() error { _, err := client.Deals.Update(ctx, 0); return err }},
		{"Deals.Delete", func() error { _, err := client.Deals.Delete(ctx, 0); return err }},
		{"Deals.ConvertToLead", func() error { _, err := client.Deals.ConvertToLead(ctx, 0); return err }},
		{"Deals.ConversionStatus", func() error { _, err := client.Deals.ConversionStatus(ctx, 0, "x"); return err }},
		{"Deals.AddFollower#id", func() error { _, err := client.Deals.AddFollower(ctx, 0, 1); return err }},
		{"Deals.AddFollower#userID", func() error { _, err := client.Deals.AddFollower(ctx, 1, 0); return err }},
		{"Deals.ListProductsAcrossDeals", func() error {
			_, _, err := client.Deals.ListProductsAcrossDeals(ctx, []DealID{1, 0})
			return err
		}},
		{"Deals.ListInstallments", func() error { _, _, err := client.Deals.ListInstallments(ctx, []DealID{1, 0}); return err }},
		{"Deals.DeleteFollower#id", func() error { _, err := client.Deals.DeleteFollower(ctx, 0, 1); return err }},
		{"Deals.DeleteFollower#followerID", func() error { _, err := client.Deals.DeleteFollower(ctx, 1, 0); return err }},
		{"Deals.AddProduct", func() error { _, err := client.Deals.AddProduct(ctx, 0); return err }},
		{"Deals.AddProducts", func() error { _, err := client.Deals.AddProducts(ctx, 0, nil); return err }},
		{"Deals.UpdateProduct#id", func() error { _, err := client.Deals.UpdateProduct(ctx, 0, 1); return err }},
		{"Deals.UpdateProduct#attachmentID", func() error { _, err := client.Deals.UpdateProduct(ctx, 1, 0); return err }},
		{"Deals.DeleteProduct#id", func() error { _, err := client.Deals.DeleteProduct(ctx, 0, 1); return err }},
		{"Deals.DeleteProduct#attachmentID", func() error { _, err := client.Deals.DeleteProduct(ctx, 1, 0); return err }},
		{"Deals.DeleteProducts", func() error { _, err := client.Deals.DeleteProducts(ctx, 0); return err }},
		{"Deals.ListAdditionalDiscounts", func() error { _, err := client.Deals.ListAdditionalDiscounts(ctx, 0); return err }},
		{"Deals.AddAdditionalDiscount", func() error { _, err := client.Deals.AddAdditionalDiscount(ctx, 0); return err }},
		{"Deals.UpdateAdditionalDiscount", func() error { _, err := client.Deals.UpdateAdditionalDiscount(ctx, 0, "x"); return err }},
		{"Deals.DeleteAdditionalDiscount", func() error { _, err := client.Deals.DeleteAdditionalDiscount(ctx, 0, "x"); return err }},
		{"Deals.AddInstallment", func() error { _, err := client.Deals.AddInstallment(ctx, 0); return err }},
		{"Deals.UpdateInstallment#id", func() error { _, err := client.Deals.UpdateInstallment(ctx, 0, 1); return err }},
		{"Deals.UpdateInstallment#installmentID", func() error { _, err := client.Deals.UpdateInstallment(ctx, 1, 0); return err }},
		{"Deals.DeleteInstallment#id", func() error { _, err := client.Deals.DeleteInstallment(ctx, 0, 1); return err }},
		{"Deals.DeleteInstallment#installmentID", func() error { _, err := client.Deals.DeleteInstallment(ctx, 1, 0); return err }},
		{"Organizations.Get", func() error { _, err := client.Organizations.Get(ctx, 0); return err }},
		{"Organizations.Update", func() error { _, err := client.Organizations.Update(ctx, 0); return err }},
		{"Organizations.Delete", func() error { _, err := client.Organizations.Delete(ctx, 0); return err }},
		{"Organizations.AddFollower#id", func() error { _, err := client.Organizations.AddFollower(ctx, 0, 1); return err }},
		{"Organizations.AddFollower#userID", func() error { _, err := client.Organizations.AddFollower(ctx, 1, 0); return err }},
		{"Organizations.DeleteFollower#id", func() error { _, err := client.Organizations.DeleteFollower(ctx, 0, 1); return err }},
		{"Organizations.DeleteFollower#followerID", func() error { _, err := client.Organizations.DeleteFollower(ctx, 1, 0); return err }},
		{"Persons.Get", func() error { _, err := client.Persons.Get(ctx, 0); return err }},
		{"Persons.Update", func() error { _, err := client.Persons.Update(ctx, 0); return err }},
		{"Persons.Delete", func() error { _, err := client.Persons.Delete(ctx, 0); return err }},
		{"Persons.AddFollower#id", func() error { _, err := client.Persons.AddFollower(ctx, 0, 1); return err }},
		{"Persons.AddFollower#userID", func() error { _, err := client.Persons.AddFollower(ctx, 1, 0); return err }},
		{"Persons.DeleteFollower#id", func() error { _, err := client.Persons.DeleteFollower(ctx, 0, 1); return err }},
		{"Persons.DeleteFollower#followerID", func() error { _, err := client.Persons.DeleteFollower(ctx, 1, 0); return err }},
		{"Persons.GetPicture", func() error { _, err := client.Persons.GetPicture(ctx, 0); return err }},
		{"Pipelines.Get", func() error { _, err := client.Pipelines.Get(ctx, 0); return err }},
		{"Pipelines.Update", func() error { _, err := client.Pipelines.Update(ctx, 0); return err }},
		{"Pipelines.Delete", func() error { _, err := client.Pipelines.Delete(ctx, 0); return err }},
		{"Products.Get", func() error { _, err := client.Products.Get(ctx, 0); return err }},
		{"Products.Update", func() error { _, err := client.Products.Update(ctx, 0); return err }},
		{"Products.Delete", func() error { _, err := client.Products.Delete(ctx, 0); return err }},
		{"Products.Duplicate", func() error { _, err := client.Products.Duplicate(ctx, 0); return err }},
		{"Products.CreateVariation", func() error { _, err := client.Products.CreateVariation(ctx, 0); return err }},
		{"Products.UpdateVariation#id", func() error { _, err := client.Products.UpdateVariation(ctx, 0, 1); return err }},
		{"Products.UpdateVariation#variationID", func() error { _, err := client.Products.UpdateVariation(ctx, 1, 0); return err }},
		{"Products.DeleteVariation#id", func() error { _, err := client.Products.DeleteVariation(ctx, 0, 1); return err }},
		{"Products.DeleteVariation#variationID", func() error { _, err := client.Products.DeleteVariation(ctx, 1, 0); return err }},
		{"Products.GetImage", func() error { _, err := client.Products.GetImage(ctx, 0); return err }},
		{"Products.UploadImage", func() error { _, err := client.Products.UploadImage(ctx, 0); return err }},
		{"Products.UpdateImage", func() error { _, err := client.Products.UpdateImage(ctx, 0); return err }},
		{"Products.DeleteImage", func() error { _, err := client.Products.DeleteImage(ctx, 0); return err }},
		{"Products.AddFollower#id", func() error { _, err := client.Products.AddFollower(ctx, 0, 1); return err }},
		{"Products.AddFollower#userID", func() error { _, err := client.Products.AddFollower(ctx, 1, 0); return err }},
		{"Products.DeleteFollower#id", func() error { _, err := client.Products.DeleteFollower(ctx, 0, 1); return err }},
		{"Products.DeleteFollower#followerID", func() error { _, err := client.Products.DeleteFollower(ctx, 1, 0); return err }},
		{"ProjectBoards.Get", func() error { _, err := client.ProjectBoards.Get(ctx, 0); return err }},
		{"ProjectBoards.Update", func() error { _, err := client.ProjectBoards.Update(ctx, 0); return err }},
		{"ProjectBoards.Delete", func() error { _, err := client.ProjectBoards.Delete(ctx, 0); return err }},
		{"ProjectPhases.List", func() error { _, err := client.ProjectPhases.List(ctx, 0); return err }},
		{"ProjectPhases.Get", func() error { _, err := client.ProjectPhases.Get(ctx, 0); return err }},
		{"ProjectPhases.Update", func() error { _, err := client.ProjectPhases.Update(ctx, 0); return err }},
		{"ProjectPhases.Delete", func() error { _, err := client.ProjectPhases.Delete(ctx, 0); return err }},
		{"ProjectTemplates.Get", func() error { _, err := client.ProjectTemplates.Get(ctx, 0); return err }},
		{"Projects.Get", func() error { _, err := client.Projects.Get(ctx, 0); return err }},
		{"Projects.Update", func() error { _, err := client.Projects.Update(ctx, 0); return err }},
		{"Projects.Delete", func() error { _, err := client.Projects.Delete(ctx, 0); return err }},
		{"Projects.Archive", func() error { _, err := client.Projects.Archive(ctx, 0); return err }},
		{"Projects.ListPermittedUsers", func() error { _, err := client.Projects.ListPermittedUsers(ctx, 0); return err }},
		{"Stages.Get", func() error { _, err := client.Stages.Get(ctx, 0); return err }},
		{"Stages.Update", func() error { _, err := client.Stages.Update(ctx, 0); return err }},
		{"Stages.Delete", func() error { _, err := client.Stages.Delete(ctx, 0); return err }},
		{"Tasks.Get", func() error { _, err := client.Tasks.Get(ctx, 0); return err }},
		{"Tasks.Update", func() error { _, err := client.Tasks.Update(ctx, 0); return err }},
		{"Tasks.Delete", func() error { _, err := client.Tasks.Delete(ctx, 0); return err }},
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

	// Guards living in unexported list helpers fire through their public
	// wrappers.
	noop := func(FollowerChangelog) error { return nil }
	helperCalls := []struct {
		name string
		call func() error
	}{
		{"Deals.ListProducts", func() error { _, _, err := client.Deals.ListProducts(ctx, 0); return err }},
		{"Deals.ForEachFollowersChangelog", func() error { return client.Deals.ForEachFollowersChangelog(ctx, 0, noop) }},
		{"Organizations.ListFollowers", func() error { _, _, err := client.Organizations.ListFollowers(ctx, 0); return err }},
		{"Organizations.ForEachFollowersChangelog", func() error {
			return client.Organizations.ForEachFollowersChangelog(ctx, 0, noop)
		}},
		{"Persons.ListFollowers", func() error { _, _, err := client.Persons.ListFollowers(ctx, 0); return err }},
		{"Persons.ForEachFollowersChangelog", func() error { return client.Persons.ForEachFollowersChangelog(ctx, 0, noop) }},
		{"Products.ListVariations", func() error { _, _, err := client.Products.ListVariations(ctx, 0); return err }},
		{"Products.ListFollowers", func() error { _, _, err := client.Products.ListFollowers(ctx, 0); return err }},
		{"Products.ForEachFollowersChangelog", func() error { return client.Products.ForEachFollowersChangelog(ctx, 0, noop) }},
		{"Projects.ListChangelog", func() error { _, _, err := client.Projects.ListChangelog(ctx, 0); return err }},
		{"Users.ListFollowers", func() error { _, _, err := client.Users.ListFollowers(ctx, 0); return err }},
	}
	for _, c := range helperCalls {
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

func TestFieldDeleteOptionsRejectsEveryNonPositiveID(t *testing.T) {
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

	methods := []struct {
		name string
		call func([]int) error
	}{
		{"DealFields.DeleteOptions", func(ids []int) error {
			_, err := client.DealFields.DeleteOptions(ctx, "cf_1", ids)
			return err
		}},
		{"OrganizationFields.DeleteOptions", func(ids []int) error {
			_, err := client.OrganizationFields.DeleteOptions(ctx, "cf_1", ids)
			return err
		}},
		{"PersonFields.DeleteOptions", func(ids []int) error {
			_, err := client.PersonFields.DeleteOptions(ctx, "cf_1", ids)
			return err
		}},
		{"ProductFields.DeleteOptions", func(ids []int) error {
			_, err := client.ProductFields.DeleteOptions(ctx, "cf_1", ids)
			return err
		}},
	}
	tests := []struct {
		name string
		ids  []int
	}{
		{"zero after valid ID", []int{1, 0}},
		{"negative after valid ID", []int{1, -1}},
	}

	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					err := method.call(test.ids)
					if err == nil {
						t.Fatal("expected client-side validation error")
					}
					if !strings.Contains(err.Error(), "invalid field option id") {
						t.Fatalf("expected field option id validation error, got %v", err)
					}
				})
			}
		})
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
