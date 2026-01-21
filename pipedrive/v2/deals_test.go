package v2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestDealsService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals/1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":1,"title":"Test deal"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	deal, err := client.Deals.Get(context.Background(), 1, WithDealRequestOptions(pipedrive.WithHeader("X-Test", "1")))
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if deal.ID != 1 || deal.Title != "Test deal" {
		t.Fatalf("unexpected deal: %#v", deal)
	}
}

func TestDealsService_ListPager(t *testing.T) {
	t.Parallel()

	var listCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		switch r.URL.Path {
		case "/deals":
			listCalls++
			cursor := r.URL.Query().Get("cursor")
			w.Header().Set("Content-Type", "application/json")
			if listCalls == 1 {
				if cursor != "" {
					t.Fatalf("expected no cursor on first page, got %q", cursor)
				}
				if got := r.URL.Query().Get("limit"); got != "2" {
					t.Fatalf("expected limit=2, got %q", got)
				}
				_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":"c2"}}`))
				return
			}
			if listCalls == 2 {
				if cursor != "c2" {
					t.Fatalf("expected cursor c2 on second page, got %q", cursor)
				}
				if got := r.URL.Query().Get("limit"); got != "2" {
					t.Fatalf("expected limit=2, got %q", got)
				}
				_, _ = w.Write([]byte(`{"data":[{"id":3}],"additional_data":{"next_cursor":null}}`))
				return
			}
			t.Fatalf("unexpected listCalls=%d", listCalls)
		default:
			if strings.HasPrefix(r.URL.Path, "/deals/") {
				t.Fatalf("unexpected deal path: %s", r.URL.Path)
			}
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	pager := client.Deals.ListPager(WithDealsPageSize(2))

	var ids []DealID
	for pager.Next(context.Background()) {
		for _, d := range pager.Items() {
			ids = append(ids, d.ID)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pager error: %v", err)
	}
	if want := []DealID{1, 2, 3}; len(ids) != len(want) || ids[0] != want[0] || ids[1] != want[1] || ids[2] != want[2] {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestDealsService_ForEach(t *testing.T) {
	t.Parallel()

	var listCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/deals" {
			http.NotFound(w, r)
			return
		}
		listCalls++
		cursor := r.URL.Query().Get("cursor")
		w.Header().Set("Content-Type", "application/json")
		if listCalls == 1 {
			if cursor != "" {
				t.Fatalf("expected no cursor on first page, got %q", cursor)
			}
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("expected limit=2, got %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1},{"id":2}],"additional_data":{"next_cursor":"c2"}}`))
			return
		}
		if listCalls == 2 {
			if cursor != "c2" {
				t.Fatalf("expected cursor c2 on second page, got %q", cursor)
			}
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("expected limit=2, got %q", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":3}],"additional_data":{"next_cursor":null}}`))
			return
		}
		t.Fatalf("unexpected listCalls=%d", listCalls)
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	var ids []DealID
	err = client.Deals.ForEach(context.Background(), func(d Deal) error {
		ids = append(ids, d.ID)
		return nil
	}, WithDealsPageSize(2))
	if err != nil {
		t.Fatalf("ForEach error: %v", err)
	}
	if want := []DealID{1, 2, 3}; len(ids) != len(want) || ids[0] != want[0] || ids[1] != want[1] || ids[2] != want[2] {
		t.Fatalf("unexpected ids: %v", ids)
	}
}
