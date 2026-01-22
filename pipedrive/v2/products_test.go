package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestProductsService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("owner_id"); got != "5" {
			t.Fatalf("unexpected owner_id: %q", got)
		}
		if got := q.Get("ids"); got != "1,2" {
			t.Fatalf("unexpected ids: %q", got)
		}
		if got := q.Get("filter_id"); got != "3" {
			t.Fatalf("unexpected filter_id: %q", got)
		}
		if got := q.Get("sort_by"); got != "name" {
			t.Fatalf("unexpected sort_by: %q", got)
		}
		if got := q.Get("sort_direction"); got != "asc" {
			t.Fatalf("unexpected sort_direction: %q", got)
		}
		if got := q.Get("custom_fields"); got != "cf_1" {
			t.Fatalf("unexpected custom_fields: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1,"name":"Product"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	products, next, err := client.Products.List(
		context.Background(),
		WithProductsOwnerID(UserID(5)),
		WithProductsIDs(ProductID(1), ProductID(2)),
		WithProductsFilterID(3),
		WithProductsSortBy(ProductSortByName),
		WithProductsSortDirection(SortAsc),
		WithProductsCustomFields("cf_1"),
		WithProductsPageSize(2),
		WithProductsCursor("c1"),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(products) != 1 || products[0].ID != 1 {
		t.Fatalf("unexpected products: %#v", products)
	}
}

func TestProductsService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Widget" {
			t.Fatalf("unexpected name: %v", payload["name"])
		}
		if payload["code"] != "W-1" {
			t.Fatalf("unexpected code: %v", payload["code"])
		}
		if payload["unit"] != "pcs" {
			t.Fatalf("unexpected unit: %v", payload["unit"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":9,"name":"Widget"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	product, err := client.Products.Create(
		context.Background(),
		WithProductName("Widget"),
		WithProductCode("W-1"),
		WithProductUnit("pcs"),
		WithProductRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if product.ID != 9 || product.Name != "Widget" {
		t.Fatalf("unexpected product: %#v", product)
	}
}

func TestProductsService_Search(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("term"); got != "wid" {
			t.Fatalf("unexpected term: %q", got)
		}
		if got := q.Get("fields"); got != "name,code" {
			t.Fatalf("unexpected fields: %q", got)
		}
		if got := q.Get("exact_match"); got != "true" {
			t.Fatalf("unexpected exact_match: %q", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c2" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"items":[{"result_score":0.9,"item":{"id":1}}]},"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	results, next, err := client.Products.Search(
		context.Background(),
		"wid",
		WithProductSearchFields(ProductSearchFieldName, ProductSearchFieldCode),
		WithProductSearchExactMatch(true),
		WithProductSearchPageSize(1),
		WithProductSearchCursor("c2"),
	)
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(results.Items) != 1 {
		t.Fatalf("unexpected results: %#v", results)
	}
}

func TestProductsService_Duplicate(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/duplicate" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":6,"name":"Copy"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	product, err := client.Products.Duplicate(context.Background(), ProductID(5))
	if err != nil {
		t.Fatalf("Duplicate error: %v", err)
	}
	if product.ID != 6 || product.Name != "Copy" {
		t.Fatalf("unexpected product: %#v", product)
	}
}

func TestProductsService_ListVariations(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/variations" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":10,"name":"Variant","product_id":5}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	variations, next, err := client.Products.ListVariations(
		context.Background(),
		ProductID(5),
		WithProductVariationsPageSize(2),
		WithProductVariationsCursor("c1"),
	)
	if err != nil {
		t.Fatalf("ListVariations error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(variations) != 1 || variations[0].ID != 10 {
		t.Fatalf("unexpected variations: %#v", variations)
	}
}

func TestProductsService_CreateVariation(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/variations" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Gold" {
			t.Fatalf("unexpected name: %v", payload["name"])
		}
		prices, ok := payload["prices"].([]interface{})
		if !ok || len(prices) != 1 {
			t.Fatalf("unexpected prices: %#v", payload["prices"])
		}
		price, ok := prices[0].(map[string]interface{})
		if !ok {
			t.Fatalf("unexpected price payload: %#v", prices[0])
		}
		if price["currency"] != "USD" {
			t.Fatalf("unexpected price currency: %#v", price["currency"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":12,"name":"Gold","product_id":5}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	variation, err := client.Products.CreateVariation(
		context.Background(),
		ProductID(5),
		WithProductVariationName("Gold"),
		WithProductVariationPrices(ProductPrice{Currency: "USD", Price: 12.5}),
	)
	if err != nil {
		t.Fatalf("CreateVariation error: %v", err)
	}
	if variation.ID != 12 || variation.Name != "Gold" {
		t.Fatalf("unexpected variation: %#v", variation)
	}
}

func TestProductsService_UpdateVariation(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/variations/12" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "Platinum" {
			t.Fatalf("unexpected name: %v", payload["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":12,"name":"Platinum","product_id":5}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	variation, err := client.Products.UpdateVariation(
		context.Background(),
		ProductID(5),
		ProductVariationID(12),
		WithProductVariationName("Platinum"),
	)
	if err != nil {
		t.Fatalf("UpdateVariation error: %v", err)
	}
	if variation.ID != 12 || variation.Name != "Platinum" {
		t.Fatalf("unexpected variation: %#v", variation)
	}
}

func TestProductsService_DeleteVariation(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/variations/12" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":12}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Products.DeleteVariation(context.Background(), ProductID(5), ProductVariationID(12))
	if err != nil {
		t.Fatalf("DeleteVariation error: %v", err)
	}
	if result.ID != 12 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestProductsService_ListFollowers(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/followers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"user_id":7}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	followers, next, err := client.Products.ListFollowers(
		context.Background(),
		ProductID(5),
		WithProductFollowersPageSize(2),
		WithProductFollowersCursor("c1"),
	)
	if err != nil {
		t.Fatalf("ListFollowers error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(followers) != 1 || followers[0].UserID != 7 {
		t.Fatalf("unexpected followers: %#v", followers)
	}
}

func TestProductsService_AddFollower(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/followers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["user_id"] != float64(7) {
			t.Fatalf("unexpected user_id: %v", payload["user_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user_id":7}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	follower, err := client.Products.AddFollower(context.Background(), ProductID(5), UserID(7))
	if err != nil {
		t.Fatalf("AddFollower error: %v", err)
	}
	if follower.UserID != 7 {
		t.Fatalf("unexpected follower: %#v", follower)
	}
}

func TestProductsService_DeleteFollower(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/followers/7" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user_id":7}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Products.DeleteFollower(context.Background(), ProductID(5), UserID(7))
	if err != nil {
		t.Fatalf("DeleteFollower error: %v", err)
	}
	if result.UserID != 7 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestProductsService_FollowersChangelog(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/followers/changelog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c1" {
			t.Fatalf("unexpected cursor: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"action":"added","actor_user_id":1,"follower_user_id":2,"time":"2024-01-01T10:00:00Z"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	changelog, next, err := client.Products.FollowersChangelog(
		context.Background(),
		ProductID(5),
		WithProductFollowersChangelogPageSize(1),
		WithProductFollowersChangelogCursor("c1"),
	)
	if err != nil {
		t.Fatalf("FollowersChangelog error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(changelog) != 1 || changelog[0].FollowerUserID != 2 {
		t.Fatalf("unexpected changelog: %#v", changelog)
	}
}

func TestProductsService_GetImage(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/images" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":9,"product_id":5,"company_id":"33","public_url":"https://cdn.example.com/p.png","add_time":"2024-01-01T10:00:00Z","mime_type":"image/png","name":"p.png"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	image, err := client.Products.GetImage(context.Background(), ProductID(5))
	if err != nil {
		t.Fatalf("GetImage error: %v", err)
	}
	if image.ID != 9 {
		t.Fatalf("unexpected image: %#v", image)
	}
	if image.CompanyID == nil || *image.CompanyID != "33" {
		t.Fatalf("unexpected company ID: %#v", image.CompanyID)
	}
	if image.PublicURL == nil || *image.PublicURL == "" {
		t.Fatalf("unexpected public URL: %#v", image.PublicURL)
	}
	parsed, err := time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")
	if err != nil {
		t.Fatalf("time parse error: %v", err)
	}
	if image.AddTime == nil || !image.AddTime.Equal(parsed) {
		t.Fatalf("unexpected add time: %#v", image.AddTime)
	}
}

func TestProductsService_UploadImage(t *testing.T) {
	t.Parallel()

	content := []byte("image-bytes")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/images" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/form-data;") {
			t.Fatalf("unexpected content-type: %q", ct)
		}

		reader, err := r.MultipartReader()
		if err != nil {
			t.Fatalf("multipart reader: %v", err)
		}
		part, err := reader.NextPart()
		if err != nil {
			t.Fatalf("multipart part: %v", err)
		}
		if part.FormName() != "data" {
			t.Fatalf("unexpected form name: %q", part.FormName())
		}
		if part.FileName() != "image.png" {
			t.Fatalf("unexpected file name: %q", part.FileName())
		}
		body, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("read part: %v", err)
		}
		if string(body) != string(content) {
			t.Fatalf("unexpected body: %q", string(body))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":11,"product_id":5,"company_id":"33","add_time":"2024-01-01T10:00:00Z"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	image, err := client.Products.UploadImage(
		context.Background(),
		ProductID(5),
		WithProductImageFile("image.png", bytes.NewReader(content)),
	)
	if err != nil {
		t.Fatalf("UploadImage error: %v", err)
	}
	if image.ID != 11 {
		t.Fatalf("unexpected image: %#v", image)
	}
	if image.CompanyID == nil || *image.CompanyID != "33" {
		t.Fatalf("unexpected company ID: %#v", image.CompanyID)
	}
}

func TestProductsService_UpdateImage(t *testing.T) {
	t.Parallel()

	content := []byte("updated-bytes")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/images" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/form-data;") {
			t.Fatalf("unexpected content-type: %q", ct)
		}

		reader, err := r.MultipartReader()
		if err != nil {
			t.Fatalf("multipart reader: %v", err)
		}
		part, err := reader.NextPart()
		if err != nil {
			t.Fatalf("multipart part: %v", err)
		}
		if part.FormName() != "data" {
			t.Fatalf("unexpected form name: %q", part.FormName())
		}
		if part.FileName() != "image.png" {
			t.Fatalf("unexpected file name: %q", part.FileName())
		}
		body, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("read part: %v", err)
		}
		if string(body) != string(content) {
			t.Fatalf("unexpected body: %q", string(body))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":11,"product_id":5,"company_id":"33","add_time":"2024-01-01T10:00:00Z"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	image, err := client.Products.UpdateImage(
		context.Background(),
		ProductID(5),
		WithProductImageFile("image.png", bytes.NewReader(content)),
	)
	if err != nil {
		t.Fatalf("UpdateImage error: %v", err)
	}
	if image.ID != 11 {
		t.Fatalf("unexpected image: %#v", image)
	}
	if image.CompanyID == nil || *image.CompanyID != "33" {
		t.Fatalf("unexpected company ID: %#v", image.CompanyID)
	}
}

func TestProductsService_DeleteImage(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/products/5/images" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":11}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	result, err := client.Products.DeleteImage(context.Background(), ProductID(5))
	if err != nil {
		t.Fatalf("DeleteImage error: %v", err)
	}
	if result.ID != 11 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
