package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestProjectFieldsService_AllOperations(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Test"); got != "project-fields" {
			t.Fatalf("unexpected request header: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /projectFields":
			if query := r.URL.Query(); query.Get("limit") != "20" || query.Get("cursor") != "fields-start" {
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"field_code":"cf","field_name":"Risk","field_type":"enum","options":[{"id":"system-high","label":"High"}]}],"additional_data":{"next_cursor":null}}`))
		case "POST /projectFields":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["field_name"] != "Risk" || body["field_type"] != "enum" ||
				len(body["options"].([]any)) != 2 || body["ui_visibility"].(map[string]any)["enabled"] != true ||
				body["important_fields"].(map[string]any)["pipeline_id"] != float64(1) ||
				body["required_fields"].(map[string]any)["stage_id"] != float64(2) {
				t.Fatalf("unexpected body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"field_code":"cf","field_name":"Risk","field_type":"enum"}}`))
		case "GET /projectFields/cf":
			_, _ = w.Write([]byte(`{"data":{"field_code":"cf","field_name":"Risk","field_type":"enum"}}`))
		case "PATCH /projectFields/cf":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["field_name"] != "Updated" || body["ui_visibility"].(map[string]any)["enabled"] != false ||
				body["important_fields"].(map[string]any)["pipeline_id"] != float64(3) ||
				body["required_fields"].(map[string]any)["stage_id"] != float64(4) {
				t.Fatalf("unexpected update body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"data":{"field_code":"cf","field_name":"Updated","field_type":"enum"}}`))
		case "DELETE /projectFields/cf":
			_, _ = w.Write([]byte(`{"data":{"field_code":"cf","field_name":"Risk","field_type":"enum"}}`))
		case "POST /projectFields/cf/options":
			var body []map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body) != 1 || body[0]["label"] != "High" {
				t.Fatalf("unexpected add options body: %#v, %v", body, err)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1,"label":"High"}]}`))
		case "PATCH /projectFields/cf/options":
			var body []map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body) != 1 ||
				body[0]["id"] != float64(1) || body[0]["label"] != "Critical" {
				t.Fatalf("unexpected update options body: %#v, %v", body, err)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1,"label":"Critical"}]}`))
		case "DELETE /projectFields/cf/options":
			var body []map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body) != 1 || body[0]["id"] != float64(1) {
				t.Fatalf("unexpected delete options body: %#v, %v", body, err)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":1,"label":"Critical"}]}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	requestOpt := WithProjectFieldRequestOptions(pipedrive.WithHeader("X-Test", "project-fields"))
	items, _, err := client.ProjectFields.List(ctx,
		WithProjectFieldsPageSize(20),
		WithProjectFieldsCursor("fields-start"),
		requestOpt,
	)
	if err != nil || len(items) != 1 || items[0].FieldCode != "cf" || items[0].Options[0].StringID != "system-high" {
		t.Fatalf("List = %#v, %v", items, err)
	}
	created, err := client.ProjectFields.Create(ctx,
		WithProjectFieldName("Risk"),
		WithProjectFieldType(FieldTypeEnum),
		WithProjectFieldOptions("High", "Low"),
		WithProjectFieldUIVisibility(map[string]interface{}{"enabled": true}),
		WithProjectFieldImportantFields(map[string]interface{}{"pipeline_id": 1}),
		WithProjectFieldRequiredFields(map[string]interface{}{"stage_id": 2}),
		requestOpt,
	)
	if err != nil || created.FieldCode != "cf" {
		t.Fatalf("Create = %#v, %v", created, err)
	}
	got, err := client.ProjectFields.Get(ctx, "cf", requestOpt)
	if err != nil || got.FieldCode != "cf" {
		t.Fatalf("Get = %#v, %v", got, err)
	}
	updated, err := client.ProjectFields.Update(ctx, "cf",
		WithProjectFieldName("Updated"),
		WithProjectFieldUIVisibility(map[string]interface{}{"enabled": false}),
		WithProjectFieldImportantFields(map[string]interface{}{"pipeline_id": 3}),
		WithProjectFieldRequiredFields(map[string]interface{}{"stage_id": 4}),
		requestOpt,
	)
	if err != nil || updated.FieldName != "Updated" {
		t.Fatalf("Update = %#v, %v", updated, err)
	}
	deleted, err := client.ProjectFields.Delete(ctx, "cf", requestOpt)
	if err != nil || deleted.FieldCode != "cf" {
		t.Fatalf("Delete = %#v, %v", deleted, err)
	}
	added, err := client.ProjectFields.AddOptions(ctx, "cf", []string{"", "High"}, requestOpt)
	if err != nil || len(added) != 1 || added[0].Label != "High" {
		t.Fatalf("AddOptions = %#v, %v", added, err)
	}
	changed, err := client.ProjectFields.UpdateOptions(ctx, "cf", []FieldOptionUpdate{{ID: 1, Label: "Critical"}}, requestOpt)
	if err != nil || len(changed) != 1 || changed[0].Label != "Critical" {
		t.Fatalf("UpdateOptions = %#v, %v", changed, err)
	}
	removed, err := client.ProjectFields.DeleteOptions(ctx, "cf", []int{1}, requestOpt)
	if err != nil || len(removed) != 1 || removed[0].ID != 1 {
		t.Fatalf("DeleteOptions = %#v, %v", removed, err)
	}
}

func TestProjectFieldsService_ListPager(t *testing.T) {
	t.Parallel()

	runFieldListPagerTest(t, "/projectFields", func(client *Client) *pipedrive.CursorPager[Field] {
		return client.ProjectFields.ListPager(WithProjectFieldsPageSize(2), WithProjectFieldsCursor("start"))
	})
}

func TestProjectFieldsService_ForEach(t *testing.T) {
	t.Parallel()

	runFieldForEachTest(t, "/projectFields", func(ctx context.Context, client *Client, fn func(Field) error) error {
		return client.ProjectFields.ForEach(ctx, fn)
	})
}
