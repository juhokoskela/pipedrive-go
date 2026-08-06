package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestProjectPhasesService_AllOperations(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /phases":
			if r.URL.Query().Get("board_id") != "2" {
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":3,"name":"Build","board_id":2,"order_nr":1}]}`))
		case "POST /phases":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["name"] != "Build" || body["board_id"] != float64(2) {
				t.Fatalf("unexpected body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"id":3,"name":"Build","board_id":2}}`))
		case "GET /phases/3":
			_, _ = w.Write([]byte(`{"data":{"id":3,"name":"Build","board_id":2}}`))
		case "PATCH /phases/3":
			_, _ = w.Write([]byte(`{"data":{"id":3,"name":"Test","board_id":2}}`))
		case "DELETE /phases/3":
			_, _ = w.Write([]byte(`{"data":{"id":3}}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	items, err := client.ProjectPhases.List(ctx, 2)
	if err != nil || len(items) != 1 || items[0].ID != 3 {
		t.Fatalf("List = %#v, %v", items, err)
	}
	created, err := client.ProjectPhases.Create(ctx, WithProjectPhaseName("Build"), WithProjectPhaseBoardID(2), WithProjectPhaseOrder(1))
	if err != nil || created.ID != 3 {
		t.Fatalf("Create = %#v, %v", created, err)
	}
	got, err := client.ProjectPhases.Get(ctx, 3)
	if err != nil || got.ID != 3 {
		t.Fatalf("Get = %#v, %v", got, err)
	}
	updated, err := client.ProjectPhases.Update(ctx, 3, WithProjectPhaseName("Test"))
	if err != nil || updated.Name != "Test" {
		t.Fatalf("Update = %#v, %v", updated, err)
	}
	deleted, err := client.ProjectPhases.Delete(ctx, 3)
	if err != nil || deleted.ID != 3 {
		t.Fatalf("Delete = %#v, %v", deleted, err)
	}
}
