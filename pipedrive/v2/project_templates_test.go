package v2

import (
	"context"
	"net/http"
	"testing"
)

func TestProjectTemplatesService_AllOperations(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/projectTemplates":
			if r.URL.Query().Get("limit") != "10" {
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":6,"title":"Launch","projects_board_id":2}],"additional_data":{"next_cursor":null}}`))
		case "/projectTemplates/6":
			_, _ = w.Write([]byte(`{"data":{"id":6,"title":"Launch","projects_board_id":2}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})
	ctx := context.Background()
	items, _, err := client.ProjectTemplates.List(ctx, WithProjectTemplatesPageSize(10))
	if err != nil || len(items) != 1 || items[0].ID != 6 {
		t.Fatalf("List = %#v, %v", items, err)
	}
	got, err := client.ProjectTemplates.Get(ctx, 6)
	if err != nil || got.ID != 6 {
		t.Fatalf("Get = %#v, %v", got, err)
	}
}
