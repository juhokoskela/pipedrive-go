package v2

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestProductImageAddTime(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, raw, want string
		invalid         bool
	}{
		{name: "live UTC datetime", raw: `"2026-11-10 15:09:15"`, want: "2026-11-10T15:09:15Z"},
		{name: "spec date", raw: `"2026-11-10"`, want: "2026-11-10T00:00:00Z"},
		{name: "RFC3339", raw: `"2026-11-10T17:09:15.123+02:00"`, want: "2026-11-10T15:09:15.123Z"},
		{name: "null", raw: `null`},
		{name: "invalid", raw: `"not a date"`, invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{"data":{"id":9,"product_id":5,"name":"test.png","add_time":%s}}`, tc.raw)
			})
			image, err := client.Products.GetImage(context.Background(), 5)
			if tc.invalid {
				if err == nil {
					t.Fatal("expected invalid time error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if image.ID != 9 || image.ProductID == nil || *image.ProductID != 5 || image.Name == nil || *image.Name != "test.png" {
				t.Fatalf("lost image fields: %#v", image)
			}
			if tc.want == "" {
				if image.AddTime != nil {
					t.Fatalf("expected nil add_time, got %v", image.AddTime)
				}
				return
			}
			want, err := time.Parse(time.RFC3339Nano, tc.want)
			if err != nil {
				t.Fatal(err)
			}
			if image.AddTime == nil || !image.AddTime.Equal(want) {
				t.Fatalf("add_time = %v, want %v", image.AddTime, want)
			}
		})
	}
}
