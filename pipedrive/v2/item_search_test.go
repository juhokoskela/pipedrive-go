package v2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestItemSearchService_Search(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/itemSearch" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("term"); got != "deal" {
			t.Fatalf("unexpected term: %q", got)
		}
		if got := q.Get("item_types"); got != "deal,person" {
			t.Fatalf("unexpected item_types: %q", got)
		}
		if got := q.Get("fields"); got != "name,email" {
			t.Fatalf("unexpected fields: %q", got)
		}
		if got := q.Get("search_for_related_items"); got != "true" {
			t.Fatalf("unexpected search_for_related_items: %q", got)
		}
		if got := q.Get("exact_match"); got != "true" {
			t.Fatalf("unexpected exact_match: %q", got)
		}
		if got := q.Get("include_fields"); got != "person.picture" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"items":[{"result_score":0.9,"item":{"id":1}}],"related_items":[{"result_score":0.5,"item":{"id":2}}]},"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	results, next, err := client.ItemSearch.Search(
		context.Background(),
		"deal",
		WithItemSearchTypes(ItemSearchTypeDeal, ItemSearchTypePerson),
		WithItemSearchFields(ItemSearchFieldName, ItemSearchFieldEmail),
		WithItemSearchRelatedItems(true),
		WithItemSearchExactMatch(true),
		WithItemSearchIncludeFields(ItemSearchIncludeFieldPersonPicture),
		WithItemSearchPageSize(2),
		WithItemSearchCursor("c1"),
		WithItemSearchRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(results.Items) != 1 || len(results.RelatedItems) != 1 {
		t.Fatalf("unexpected results: %#v", results)
	}
}

func TestItemSearchService_SearchByField(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/itemSearch/field" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("term"); got != "alpha" {
			t.Fatalf("unexpected term: %q", got)
		}
		if got := q.Get("entity_type"); got != "deal" {
			t.Fatalf("unexpected entity_type: %q", got)
		}
		if got := q.Get("field"); got != "custom_field" {
			t.Fatalf("unexpected field: %q", got)
		}
		if got := q.Get("match"); got != "exact" {
			t.Fatalf("unexpected match: %q", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c2" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"result_score":1,"item":{"id":3}}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	items, next, err := client.ItemSearch.SearchByField(
		context.Background(),
		"alpha",
		ItemSearchEntityTypeDeal,
		"custom_field",
		WithItemSearchMatch(ItemSearchMatchExact),
		WithItemSearchByFieldPageSize(1),
		WithItemSearchByFieldCursor("c2"),
	)
	if err != nil {
		t.Fatalf("SearchByField error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(items) != 1 || items[0].ResultScore != 1 {
		t.Fatalf("unexpected items: %#v", items)
	}
}
