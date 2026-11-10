package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/v2/pipedrive"
)

func TestListUsersAcceptsIDOnlyResponses(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[ 123 , 456 ]}`))
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tests := []struct {
		name string
		call func() ([]User, error)
	}{
		{"deals", func() ([]User, error) { return c.Deals.ListUsers(ctx, 1) }},
		{"persons", func() ([]User, error) { return c.Persons.ListUsers(ctx, 1) }},
		{"organizations", func() ([]User, error) { return c.Organizations.ListUsers(ctx, 1) }},
		{"products", func() ([]User, error) { return c.Products.ListUsers(ctx, 1) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := tt.call()
			if err != nil {
				t.Fatal(err)
			}
			if len(users) != 2 || users[0].ID != 123 || users[1].ID != 456 {
				t.Fatalf("users = %#v", users)
			}
		})
	}
}

func TestProductDealsAcceptExpandedRelatedIDs(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"id":1,"title":"Example","user_id": {"id":10,"value":10} ,"person_id": {"value":20,"name":"Person"} ,"org_id": {"value":30,"name":"Organization"}},{"id":2,"user_id":11,"person_id":21,"org_id":31},{"id":3,"user_id":null,"person_id":null,"org_id":null}]}`))
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	deals, _, err := c.Products.ListDeals(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(deals) != 3 {
		t.Fatalf("deals = %#v", deals)
	}
	for i, d := range deals[:2] {
		if d.OwnerID == nil || int64(*d.OwnerID) != int64(10+i) || d.PersonID == nil || int64(*d.PersonID) != int64(20+i) || d.OrgID == nil || int64(*d.OrgID) != int64(30+i) {
			t.Fatalf("deal = %#v", d)
		}
	}
	if deals[0].Title != "Example" || deals[2].OwnerID != nil || deals[2].PersonID != nil || deals[2].OrgID != nil {
		t.Fatalf("deals = %#v", deals)
	}
}
