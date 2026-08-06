package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestProjectBoardsService_AllOperations(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /boards":
			_, _ = w.Write([]byte(`{"data":[{"id":2,"name":"Delivery","order_nr":1}]}`))
		case "POST /boards":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["name"] != "Delivery" || body["order_nr"] != float64(1) {
				t.Fatalf("unexpected body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"id":2,"name":"Delivery","order_nr":1}}`))
		case "GET /boards/2":
			_, _ = w.Write([]byte(`{"data":{"id":2,"name":"Delivery"}}`))
		case "PATCH /boards/2":
			_, _ = w.Write([]byte(`{"data":{"id":2,"name":"Updated"}}`))
		case "DELETE /boards/2":
			_, _ = w.Write([]byte(`{"data":{"id":2}}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	items, err := client.ProjectBoards.List(ctx)
	if err != nil || len(items) != 1 || items[0].ID != 2 {
		t.Fatalf("List = %#v, %v", items, err)
	}
	created, err := client.ProjectBoards.Create(ctx, WithProjectBoardName("Delivery"), WithProjectBoardOrder(1))
	if err != nil || created.ID != 2 {
		t.Fatalf("Create = %#v, %v", created, err)
	}
	got, err := client.ProjectBoards.Get(ctx, 2)
	if err != nil || got.ID != 2 {
		t.Fatalf("Get = %#v, %v", got, err)
	}
	updated, err := client.ProjectBoards.Update(ctx, 2, WithProjectBoardName("Updated"))
	if err != nil || updated.Name != "Updated" {
		t.Fatalf("Update = %#v, %v", updated, err)
	}
	deleted, err := client.ProjectBoards.Delete(ctx, 2)
	if err != nil || deleted.ID != 2 {
		t.Fatalf("Delete = %#v, %v", deleted, err)
	}
}
