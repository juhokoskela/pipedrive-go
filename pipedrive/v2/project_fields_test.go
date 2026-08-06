package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestProjectFieldsService_AllOperations(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /projectFields":
			_, _ = w.Write([]byte(`{"data":[{"field_code":"cf","field_name":"Risk","field_type":"enum","options":[{"id":"system-high","label":"High"}]}],"additional_data":{"next_cursor":null}}`))
		case "POST /projectFields":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["field_name"] != "Risk" || body["field_type"] != "enum" {
				t.Fatalf("unexpected body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"field_code":"cf","field_name":"Risk","field_type":"enum"}}`))
		case "GET /projectFields/cf":
			_, _ = w.Write([]byte(`{"data":{"field_code":"cf","field_name":"Risk","field_type":"enum"}}`))
		case "PATCH /projectFields/cf":
			_, _ = w.Write([]byte(`{"data":{"field_code":"cf","field_name":"Updated","field_type":"enum"}}`))
		case "DELETE /projectFields/cf":
			_, _ = w.Write([]byte(`{"data":{"field_code":"cf","field_name":"Risk","field_type":"enum"}}`))
		case "POST /projectFields/cf/options":
			_, _ = w.Write([]byte(`{"data":[{"id":1,"label":"High"}]}`))
		case "PATCH /projectFields/cf/options":
			_, _ = w.Write([]byte(`{"data":[{"id":1,"label":"Critical"}]}`))
		case "DELETE /projectFields/cf/options":
			_, _ = w.Write([]byte(`{"data":[{"id":1,"label":"Critical"}]}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	items, _, err := client.ProjectFields.List(ctx)
	if err != nil || len(items) != 1 || items[0].FieldCode != "cf" || items[0].Options[0].StringID != "system-high" {
		t.Fatalf("List = %#v, %v", items, err)
	}
	created, err := client.ProjectFields.Create(ctx, WithProjectFieldName("Risk"), WithProjectFieldType(FieldTypeEnum), WithProjectFieldOptions("High"))
	if err != nil || created.FieldCode != "cf" {
		t.Fatalf("Create = %#v, %v", created, err)
	}
	got, err := client.ProjectFields.Get(ctx, "cf")
	if err != nil || got.FieldCode != "cf" {
		t.Fatalf("Get = %#v, %v", got, err)
	}
	updated, err := client.ProjectFields.Update(ctx, "cf", WithProjectFieldName("Updated"))
	if err != nil || updated.FieldName != "Updated" {
		t.Fatalf("Update = %#v, %v", updated, err)
	}
	deleted, err := client.ProjectFields.Delete(ctx, "cf")
	if err != nil || deleted.FieldCode != "cf" {
		t.Fatalf("Delete = %#v, %v", deleted, err)
	}
	added, err := client.ProjectFields.AddOptions(ctx, "cf", []string{"High"})
	if err != nil || len(added) != 1 || added[0].Label != "High" {
		t.Fatalf("AddOptions = %#v, %v", added, err)
	}
	changed, err := client.ProjectFields.UpdateOptions(ctx, "cf", []FieldOptionUpdate{{ID: 1, Label: "Critical"}})
	if err != nil || len(changed) != 1 || changed[0].Label != "Critical" {
		t.Fatalf("UpdateOptions = %#v, %v", changed, err)
	}
	removed, err := client.ProjectFields.DeleteOptions(ctx, "cf", []int{1})
	if err != nil || len(removed) != 1 || removed[0].ID != 1 {
		t.Fatalf("DeleteOptions = %#v, %v", removed, err)
	}
}
