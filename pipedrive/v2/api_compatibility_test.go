package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestFieldOption_UnmarshalJSONIdentifierKinds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        string
		wantID       int
		wantStringID string
		wantJSONID   string
		wantErr      bool
	}{
		{
			name:       "integer identifier",
			input:      `{"id":42,"label":"High"}`,
			wantID:     42,
			wantJSONID: "42",
			wantErr:    false,
		},
		{
			name:         "string identifier",
			input:        `{"id":"system-high","label":"High"}`,
			wantStringID: "system-high",
			wantJSONID:   `"system-high"`,
			wantErr:      false,
		},
		{
			name:    "unsupported identifier",
			input:   `{"id":[42],"label":"High"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var option FieldOption
			err := json.Unmarshal([]byte(tt.input), &option)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if option.ID != tt.wantID || option.StringID != tt.wantStringID {
				t.Fatalf("unexpected identifier: ID=%d StringID=%q", option.ID, option.StringID)
			}
			encoded, err := json.Marshal(option)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			var wire map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &wire); err != nil {
				t.Fatalf("decode marshaled option: %v", err)
			}
			if got := string(wire["id"]); got != tt.wantJSONID {
				t.Fatalf("marshaled id = %s, want %s", got, tt.wantJSONID)
			}
		})
	}
}

func TestFieldCreateResponsesAcceptStringOptionIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		path   string
		create func(*Client) (*Field, error)
	}{
		{
			name: "deal field",
			path: "/dealFields",
			create: func(client *Client) (*Field, error) {
				return client.DealFields.Create(context.Background())
			},
		},
		{
			name: "organization field",
			path: "/organizationFields",
			create: func(client *Client) (*Field, error) {
				return client.OrganizationFields.Create(context.Background())
			},
		},
		{
			name: "person field",
			path: "/personFields",
			create: func(client *Client) (*Field, error) {
				return client.PersonFields.Create(context.Background())
			},
		},
		{
			name: "product field",
			path: "/productFields",
			create: func(client *Client) (*Field, error) {
				return client.ProductFields.Create(context.Background())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != tt.path {
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":{"field_code":"cf_1","field_name":"Priority","field_type":"enum","options":[{"id":"system-high","label":"High"}]}}`))
			})

			field, err := tt.create(client)
			if err != nil {
				t.Fatalf("Create error: %v", err)
			}
			if len(field.Options) != 1 || field.Options[0].StringID != "system-high" {
				t.Fatalf("unexpected options: %#v", field.Options)
			}
		})
	}
}

func TestDocumentedResponseFieldsAreRepresented(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		typ  reflect.Type
		tags []string
	}{
		{
			name: "deal",
			typ:  reflect.TypeOf(Deal{}),
			tags: []string{
				"id", "title", "owner_id", "person_id", "org_id", "pipeline_id", "stage_id",
				"value", "currency", "add_time", "update_time", "stage_change_time", "is_deleted",
				"is_archived", "status", "probability", "lost_reason", "visible_to", "close_time",
				"won_time", "lost_time", "expected_close_date", "label_ids", "origin", "origin_id",
				"channel", "channel_id", "source_lead_id", "arr", "mrr", "acv", "custom_fields",
				"activities_count", "done_activities_count", "email_messages_count", "files_count",
				"first_won_time", "followers_count", "last_activity_id", "last_incoming_mail_time",
				"last_outgoing_mail_time", "next_activity_id", "notes_count", "participants_count",
				"products_count", "smart_bcc_email", "undone_activities_count", "labels",
			},
		},
		{
			name: "person",
			typ:  reflect.TypeOf(Person{}),
			tags: []string{
				"id", "name", "first_name", "last_name", "owner_id", "org_id", "add_time", "update_time",
				"emails", "phones", "is_deleted", "visible_to", "label_ids", "picture_id", "postal_address",
				"notes", "im", "birthday", "job_title", "custom_fields", "activities_count",
				"closed_deals_count", "doi_status", "done_activities_count", "email_messages_count",
				"files_count", "followers_count", "last_activity_id", "last_incoming_mail_time",
				"last_outgoing_mail_time", "lost_deals_count", "marketing_status", "next_activity_id",
				"notes_count", "open_deals_count", "participant_closed_deals_count",
				"participant_open_deals_count", "related_closed_deals_count", "related_lost_deals_count",
				"related_open_deals_count", "related_won_deals_count", "undone_activities_count",
				"won_deals_count", "labels",
			},
		},
		{
			name: "organization",
			typ:  reflect.TypeOf(Organization{}),
			tags: []string{
				"id", "name", "owner_id", "add_time", "update_time", "is_deleted", "visible_to", "address",
				"label_ids", "website", "linkedin", "industry", "annual_revenue", "employee_count", "custom_fields",
				"activities_count", "closed_deals_count", "done_activities_count", "email_messages_count",
				"files_count", "followers_count", "last_activity_id", "lost_deals_count", "next_activity_id",
				"notes_count", "open_deals_count", "people_count", "related_closed_deals_count",
				"related_lost_deals_count", "related_open_deals_count", "related_won_deals_count",
				"undone_activities_count", "won_deals_count", "labels",
			},
		},
		{
			name: "product",
			typ:  reflect.TypeOf(Product{}),
			tags: []string{
				"id", "name", "code", "unit", "tax", "is_deleted", "is_linkable", "visible_to",
				"owner_id", "add_time", "update_time", "description", "category", "custom_fields",
				"billing_frequency", "billing_frequency_cycles", "prices",
			},
		},
		{
			name: "activity",
			typ:  reflect.TypeOf(Activity{}),
			tags: []string{
				"id", "subject", "type", "owner_id", "creator_user_id", "is_deleted", "add_time", "update_time",
				"deal_id", "lead_id", "person_id", "org_id", "project_id", "due_date", "due_time", "duration",
				"busy", "done", "marked_as_done_time", "location", "participants", "attendees",
				"conference_meeting_client", "conference_meeting_url", "conference_meeting_id",
				"public_description", "priority", "note",
			},
		},
		{
			name: "field",
			typ:  reflect.TypeOf(Field{}),
			tags: []string{
				"field_name", "field_code", "description", "field_type", "options", "subfields",
				"is_custom_field", "is_optional_response_field", "ui_visibility", "important_fields", "required_fields",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := make(map[string]bool, tt.typ.NumField())
			for i := 0; i < tt.typ.NumField(); i++ {
				tag := tt.typ.Field(i).Tag.Get("json")
				if comma := indexByte(tag, ','); comma >= 0 {
					tag = tag[:comma]
				}
				got[tag] = true
			}
			for _, tag := range tt.tags {
				if !got[tag] {
					t.Errorf("missing json field %q", tag)
				}
			}
		})
	}
}

func TestCollectionOptionsSerializeExplicitEmptyArrays(t *testing.T) {
	t.Parallel()

	dealCfg := updateDealOptions{}
	WithDealLabelIDs().applyUpdateDeal(&dealCfg)
	assertEmptyArrayValue(t, dealCfg.payload.toMap(), "label_ids")

	personCfg := updatePersonOptions{}
	WithPersonEmails().applyUpdatePerson(&personCfg)
	WithPersonPhones().applyUpdatePerson(&personCfg)
	WithPersonLabelIDs().applyUpdatePerson(&personCfg)
	assertEmptyArrayValue(t, personCfg.payload.toMap(), "emails")
	assertEmptyArrayValue(t, personCfg.payload.toMap(), "phones")
	assertEmptyArrayValue(t, personCfg.payload.toMap(), "label_ids")

	organizationCfg := updateOrganizationOptions{}
	WithOrganizationLabelIDs().applyUpdateOrganization(&organizationCfg)
	assertEmptyArrayValue(t, organizationCfg.payload.toMap(), "label_ids")

	productCfg := updateProductOptions{}
	WithProductPrices().applyUpdateProduct(&productCfg)
	assertEmptyArrayValue(t, productCfg.payload.toMap(), "prices")

	variationCfg := updateProductVariationOptions{}
	WithProductVariationPrices().applyUpdateProductVariation(&variationCfg)
	assertEmptyArrayValue(t, variationCfg.payload.toMap(), "prices")

	activityCfg := updateActivityOptions{}
	WithActivityParticipants().applyUpdateActivity(&activityCfg)
	WithActivityAttendees().applyUpdateActivity(&activityCfg)
	assertEmptyArrayValue(t, activityCfg.payload.toMap(), "participants")
	assertEmptyArrayValue(t, activityCfg.payload.toMap(), "attendees")
}

func TestNullableOptionsSerializeExplicitNull(t *testing.T) {
	t.Parallel()

	dealCfg := updateDealOptions{}
	ClearDealProbability().applyUpdateDeal(&dealCfg)
	ClearDealLostReason().applyUpdateDeal(&dealCfg)
	ClearDealCloseTime().applyUpdateDeal(&dealCfg)
	dealBody := dealCfg.payload.toMap()
	assertNullValue(t, dealBody, "probability")
	assertNullValue(t, dealBody, "lost_reason")
	assertNullValue(t, dealBody, "close_time")

	productCfg := updateProductOptions{}
	ClearProductBillingFrequencyCycles().applyUpdateProduct(&productCfg)
	assertNullValue(t, productCfg.payload.toMap(), "billing_frequency_cycles")
}

func TestLegacyEmptyDealOptionsRemainOmitted(t *testing.T) {
	t.Parallel()

	cfg := updateDealOptions{}
	WithDealLostReason("").applyUpdateDeal(&cfg)
	WithDealCloseTime("").applyUpdateDeal(&cfg)

	body := cfg.payload.toMap()
	assertMissingValue(t, body, "lost_reason")
	assertMissingValue(t, body, "close_time")
}

func TestPersonCustomFieldsOption(t *testing.T) {
	t.Parallel()

	cfg := updatePersonOptions{}
	fields := map[string]interface{}{"custom-key": "custom-value"}
	WithPersonCustomFieldsMap(fields).applyUpdatePerson(&cfg)

	got, ok := cfg.payload.toMap()["custom_fields"].(map[string]interface{})
	if !ok || !reflect.DeepEqual(got, fields) {
		t.Fatalf("unexpected custom_fields: %#v", got)
	}
}

func TestUnsupportedFieldDescriptionsAreOmitted(t *testing.T) {
	t.Parallel()

	personCfg := createPersonFieldOptions{}
	WithPersonFieldDescription("unsupported").applyCreatePersonField(&personCfg)
	assertMissingValue(t, personCfg.payload.toMap(), "description")

	organizationCfg := createOrganizationFieldOptions{}
	WithOrganizationFieldDescription("unsupported").applyCreateOrganizationField(&organizationCfg)
	assertMissingValue(t, organizationCfg.payload.toMap(), "description")

	productCfg := createProductFieldOptions{}
	WithProductFieldDescription("unsupported").applyCreateProductField(&productCfg)
	assertMissingValue(t, productCfg.payload.toMap(), "description")

	dealCfg := createDealFieldOptions{}
	WithDealFieldDescription("supported").applyCreateDealField(&dealCfg)
	if got := dealCfg.payload.toMap()["description"]; got != "supported" {
		t.Fatalf("unexpected deal field description: %#v", got)
	}
}

func TestNewEntityQueryOptions(t *testing.T) {
	t.Parallel()

	dealGet := getDealOptions{}
	WithDealIncludeLabels(true).applyGetDeal(&dealGet)
	WithDealIncludeOptionLabels(true).applyGetDeal(&dealGet)
	if dealGet.params.IncludeLabels == nil || !*dealGet.params.IncludeLabels ||
		dealGet.params.IncludeOptionLabels == nil || !*dealGet.params.IncludeOptionLabels {
		t.Fatalf("unexpected deal get params: %#v", dealGet.params)
	}

	dealsList := listDealsOptions{}
	WithDealsIncludeLabels(true).applyListDeals(&dealsList)
	WithDealsIncludeOptionLabels(true).applyListDeals(&dealsList)
	if dealsList.params.IncludeLabels == nil || !*dealsList.params.IncludeLabels ||
		dealsList.params.IncludeOptionLabels == nil || !*dealsList.params.IncludeOptionLabels {
		t.Fatalf("unexpected deals list params: %#v", dealsList.params)
	}

	personGet := getPersonOptions{}
	WithPersonIncludeLabels(true).applyGetPerson(&personGet)
	WithPersonIncludeOptionLabels(true).applyGetPerson(&personGet)
	if personGet.params.IncludeLabels == nil || !*personGet.params.IncludeLabels ||
		personGet.params.IncludeOptionLabels == nil || !*personGet.params.IncludeOptionLabels {
		t.Fatalf("unexpected person get params: %#v", personGet.params)
	}

	personsList := listPersonsOptions{}
	WithPersonsIncludeLabels(true).applyListPersons(&personsList)
	WithPersonsIncludeOptionLabels(true).applyListPersons(&personsList)
	if personsList.params.IncludeLabels == nil || !*personsList.params.IncludeLabels ||
		personsList.params.IncludeOptionLabels == nil || !*personsList.params.IncludeOptionLabels {
		t.Fatalf("unexpected persons list params: %#v", personsList.params)
	}

	organizationGet := getOrganizationOptions{}
	WithOrganizationIncludeLabels(true).applyGetOrganization(&organizationGet)
	WithOrganizationIncludeOptionLabels(true).applyGetOrganization(&organizationGet)
	if organizationGet.params.IncludeLabels == nil || !*organizationGet.params.IncludeLabels ||
		organizationGet.params.IncludeOptionLabels == nil || !*organizationGet.params.IncludeOptionLabels {
		t.Fatalf("unexpected organization get params: %#v", organizationGet.params)
	}

	organizationsList := listOrganizationsOptions{}
	WithOrganizationsIncludeLabels(true).applyListOrganizations(&organizationsList)
	WithOrganizationsIncludeOptionLabels(true).applyListOrganizations(&organizationsList)
	if organizationsList.params.IncludeLabels == nil || !*organizationsList.params.IncludeLabels ||
		organizationsList.params.IncludeOptionLabels == nil || !*organizationsList.params.IncludeOptionLabels {
		t.Fatalf("unexpected organizations list params: %#v", organizationsList.params)
	}

	updatedSince := time.Date(2026, 8, 1, 12, 30, 0, 0, time.FixedZone("test", 2*60*60))
	productsList := listProductsOptions{}
	WithProductsUpdatedSince(updatedSince).applyListProducts(&productsList)
	if productsList.params.UpdatedSince == nil || *productsList.params.UpdatedSince != updatedSince.Format(time.RFC3339) {
		t.Fatalf("unexpected products updated_since: %#v", productsList.params.UpdatedSince)
	}
}

func assertEmptyArrayValue(t *testing.T, body map[string]interface{}, key string) {
	t.Helper()

	value, ok := body[key]
	if !ok {
		t.Fatalf("missing %q", key)
	}
	valueOf := reflect.ValueOf(value)
	if valueOf.Kind() != reflect.Slice || valueOf.Len() != 0 {
		t.Fatalf("%q is not an empty array: %#v", key, value)
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &wire); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got := string(wire[key]); got != "[]" {
		t.Fatalf("%q JSON is %s, want []", key, got)
	}
}

func assertNullValue(t *testing.T, body map[string]interface{}, key string) {
	t.Helper()

	value, ok := body[key]
	if !ok {
		t.Fatalf("missing %q", key)
	}
	if value != nil {
		t.Fatalf("%q is not null: %#v", key, value)
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &wire); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got := string(wire[key]); got != "null" {
		t.Fatalf("%q JSON is %s, want null", key, got)
	}
}

func assertMissingValue(t *testing.T, body map[string]interface{}, key string) {
	t.Helper()

	if value, ok := body[key]; ok {
		t.Fatalf("unexpected %q: %#v", key, value)
	}
}

func indexByte(value string, target byte) int {
	for i := 0; i < len(value); i++ {
		if value[i] == target {
			return i
		}
	}
	return -1
}
