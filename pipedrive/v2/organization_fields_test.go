package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestOrganizationFieldsService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizationFields" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("include_fields"); got != "ui_visibility" {
			t.Fatalf("unexpected include_fields: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("cursor"); got != "c2" {
			t.Fatalf("unexpected cursor: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"field_code":"cf_1","field_name":"Priority"}],"additional_data":{"next_cursor":null}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	fields, next, err := client.OrganizationFields.List(
		context.Background(),
		WithOrganizationFieldsIncludeFields(FieldIncludeField("ui_visibility")),
		WithOrganizationFieldsPageSize(2),
		WithOrganizationFieldsCursor("c2"),
		WithOrganizationFieldRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if next != nil {
		t.Fatalf("expected nil cursor, got %q", *next)
	}
	if len(fields) != 1 || fields[0].FieldCode != "cf_1" {
		t.Fatalf("unexpected fields: %#v", fields)
	}
}

func TestOrganizationFieldsService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizationFields" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "2" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["field_name"] != "Priority" {
			t.Fatalf("unexpected field_name: %v", payload["field_name"])
		}
		if payload["field_type"] != "enum" {
			t.Fatalf("unexpected field_type: %v", payload["field_type"])
		}
		if payload["description"] != "Organization priority" {
			t.Fatalf("unexpected description: %v", payload["description"])
		}
		options, ok := payload["options"].([]interface{})
		if !ok || len(options) != 1 {
			t.Fatalf("unexpected options: %#v", payload["options"])
		}
		option, ok := options[0].(map[string]interface{})
		if !ok || option["label"] != "High" {
			t.Fatalf("unexpected option: %#v", options[0])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"field_code":"cf_1","field_name":"Priority"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	field, err := client.OrganizationFields.Create(
		context.Background(),
		WithOrganizationFieldName("Priority"),
		WithOrganizationFieldType(FieldTypeEnum),
		WithOrganizationFieldDescription("Organization priority"),
		WithOrganizationFieldOptions("High"),
		WithOrganizationFieldRequestOptions(pipedrive.WithHeader("X-Test", "2")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if field.FieldCode != "cf_1" || field.FieldName != "Priority" {
		t.Fatalf("unexpected field: %#v", field)
	}
}

func TestOrganizationFieldsService_AddOptions(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/organizationFields/cf_1/options" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var payload []map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(payload) != 2 || payload[0]["label"] != "Critical" || payload[1]["label"] != "High" {
			t.Fatalf("unexpected payload: %#v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":1,"label":"Critical"}]}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	options, err := client.OrganizationFields.AddOptions(
		context.Background(),
		"cf_1",
		[]string{"Critical", "High"},
	)
	if err != nil {
		t.Fatalf("AddOptions error: %v", err)
	}
	if len(options) != 1 || options[0].ID != 1 {
		t.Fatalf("unexpected options: %#v", options)
	}
}

func TestOrganizationFieldsService_Get(t *testing.T) {
	t.Parallel()

	runFieldGetTest(t, "/organizationFields/cf_1", func(ctx context.Context, client *Client) (*Field, error) {
		return client.OrganizationFields.Get(
			ctx,
			"cf_1",
			WithOrganizationFieldIncludeFields(FieldIncludeField("ui_visibility"), FieldIncludeField("important_fields")),
			WithOrganizationFieldRequestOptions(pipedrive.WithHeader("X-Test", "get")),
		)
	})
}

func TestOrganizationFieldsService_ListPager(t *testing.T) {
	t.Parallel()

	runFieldListPagerTest(t, "/organizationFields", func(client *Client) *pipedrive.CursorPager[Field] {
		return client.OrganizationFields.ListPager(WithOrganizationFieldsPageSize(2), WithOrganizationFieldsCursor("start"))
	})
}

func TestOrganizationFieldsService_ForEach(t *testing.T) {
	t.Parallel()

	runFieldForEachTest(t, "/organizationFields", func(ctx context.Context, client *Client, fn func(Field) error) error {
		return client.OrganizationFields.ForEach(ctx, fn)
	})
}

func TestOrganizationFieldsService_Update(t *testing.T) {
	t.Parallel()

	runFieldUpdateTest(t, "/organizationFields/cf_1", func(ctx context.Context, client *Client) (*Field, error) {
		return client.OrganizationFields.Update(
			ctx,
			"cf_1",
			WithOrganizationFieldName("Priority updated"),
			WithOrganizationFieldType(FieldTypeVarchar),
			WithOrganizationFieldDescription("Updated description"),
			WithOrganizationFieldUIVisibility(map[string]interface{}{"add": true}),
			WithOrganizationFieldImportantFields(map[string]interface{}{"pipeline_id": 1}),
			WithOrganizationFieldRequiredFields(map[string]interface{}{"stage_id": 2}),
			WithOrganizationFieldRequestOptions(pipedrive.WithHeader("X-Test", "update")),
		)
	})
}

func TestOrganizationFieldsService_Delete(t *testing.T) {
	t.Parallel()

	runFieldDeleteTest(t, "/organizationFields/cf_1", func(ctx context.Context, client *Client) (*Field, error) {
		return client.OrganizationFields.Delete(
			ctx,
			"cf_1",
			WithOrganizationFieldRequestOptions(pipedrive.WithHeader("X-Test", "delete")),
		)
	})
}

func TestOrganizationFieldsService_AddOptionsWithRequestOptions(t *testing.T) {
	t.Parallel()

	runFieldAddOptionsRequestOptionsTest(t, "/organizationFields/cf_1/options", func(ctx context.Context, client *Client) ([]FieldOption, error) {
		return client.OrganizationFields.AddOptions(
			ctx,
			"cf_1",
			[]string{"Critical"},
			WithOrganizationFieldRequestOptions(pipedrive.WithHeader("X-Test", "add-options")),
		)
	})
}

func TestOrganizationFieldsService_UpdateOptions(t *testing.T) {
	t.Parallel()

	runFieldUpdateOptionsTest(t, "/organizationFields/cf_1/options", func(ctx context.Context, client *Client) ([]FieldOption, error) {
		return client.OrganizationFields.UpdateOptions(
			ctx,
			"cf_1",
			[]FieldOptionUpdate{{ID: 1, Label: "Critical"}},
			WithOrganizationFieldRequestOptions(pipedrive.WithHeader("X-Test", "update-options")),
		)
	})
}

func TestOrganizationFieldsService_DeleteOptions(t *testing.T) {
	t.Parallel()

	runFieldDeleteOptionsTest(t, "/organizationFields/cf_1/options", func(ctx context.Context, client *Client) ([]FieldOption, error) {
		return client.OrganizationFields.DeleteOptions(
			ctx,
			"cf_1",
			[]int{1, 2},
			WithOrganizationFieldRequestOptions(pipedrive.WithHeader("X-Test", "delete-options")),
		)
	})
}
