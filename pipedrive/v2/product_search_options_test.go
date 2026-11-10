package v2

import (
	"context"
	"net/http"
	"testing"
)

func TestProductSearchIncludeFields(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("include_fields"); got != "product.price" {
			t.Fatalf("include_fields = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"items":[]}}`))
	})
	_, _, err := client.Products.Search(context.Background(), "sample", WithProductSearchIncludeFields(ProductSearchIncludeFieldPrice), WithProductSearchIncludeFields())
	if err != nil {
		t.Fatal(err)
	}
}
