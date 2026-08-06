package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FieldType string

const (
	FieldTypeInt              FieldType = "int"
	FieldTypeDouble           FieldType = "double"
	FieldTypeBoolean          FieldType = "boolean"
	FieldTypeVarchar          FieldType = "varchar"
	FieldTypeText             FieldType = "text"
	FieldTypePhone            FieldType = "phone"
	FieldTypeVarcharOptions   FieldType = "varchar_options"
	FieldTypeVarcharAuto      FieldType = "varchar_auto"
	FieldTypeDate             FieldType = "date"
	FieldTypeDateRange        FieldType = "daterange"
	FieldTypeTime             FieldType = "time"
	FieldTypeTimeRange        FieldType = "timerange"
	FieldTypeEnum             FieldType = "enum"
	FieldTypeSet              FieldType = "set"
	FieldTypeAddress          FieldType = "address"
	FieldTypeMonetary         FieldType = "monetary"
	FieldTypeDeal             FieldType = "deal"
	FieldTypeDeals            FieldType = "deals"
	FieldTypeLead             FieldType = "lead"
	FieldTypeOrg              FieldType = "org"
	FieldTypePeople           FieldType = "people"
	FieldTypeProject          FieldType = "project"
	FieldTypeStage            FieldType = "stage"
	FieldTypeUser             FieldType = "user"
	FieldTypeActivity         FieldType = "activity"
	FieldTypeJSON             FieldType = "json"
	FieldTypePicture          FieldType = "picture"
	FieldTypeStatus           FieldType = "status"
	FieldTypeVisibleTo        FieldType = "visible_to"
	FieldTypePriceList        FieldType = "price_list"
	FieldTypeBillingFrequency FieldType = "billing_frequency"
	FieldTypeProjectsBoard    FieldType = "projects_board"
	FieldTypeProjectsPhase    FieldType = "projects_phase"
)

type FieldIncludeField string

type FieldOption struct {
	// ID contains numeric custom-field option identifiers. It remains an int for
	// source compatibility with earlier releases.
	ID int `json:"-"`
	// StringID contains built-in field option identifiers when Pipedrive returns
	// the documented string form. Exactly one of ID and StringID is normally set.
	StringID   string     `json:"-"`
	Label      string     `json:"label,omitempty"`
	Color      *string    `json:"color,omitempty"`
	UpdateTime *time.Time `json:"update_time,omitempty"`
	AddTime    *time.Time `json:"add_time,omitempty"`
}

func (o *FieldOption) UnmarshalJSON(data []byte) error {
	if o == nil {
		return fmt.Errorf("v2.FieldOption: UnmarshalJSON on nil receiver")
	}

	var wire struct {
		ID         json.RawMessage `json:"id"`
		Label      string          `json:"label,omitempty"`
		Color      *string         `json:"color,omitempty"`
		UpdateTime *time.Time      `json:"update_time,omitempty"`
		AddTime    *time.Time      `json:"add_time,omitempty"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return fmt.Errorf("v2.FieldOption: decode: %w", err)
	}

	decoded := FieldOption{
		Label:      wire.Label,
		Color:      wire.Color,
		UpdateTime: wire.UpdateTime,
		AddTime:    wire.AddTime,
	}
	id := bytes.TrimSpace(wire.ID)
	if len(id) == 0 {
		*o = decoded
		return nil
	}
	if bytes.Equal(id, []byte("null")) {
		return fmt.Errorf("v2.FieldOption: id must be an integer or string")
	}

	if id[0] == '"' {
		if err := json.Unmarshal(id, &decoded.StringID); err != nil {
			return fmt.Errorf("v2.FieldOption: decode string id: %w", err)
		}
		*o = decoded
		return nil
	}
	if err := json.Unmarshal(id, &decoded.ID); err != nil {
		return fmt.Errorf("v2.FieldOption: decode integer id: %w", err)
	}
	*o = decoded
	return nil
}

func (o FieldOption) MarshalJSON() ([]byte, error) {
	var id interface{}
	switch {
	case o.StringID != "":
		id = o.StringID
	case o.ID != 0:
		id = o.ID
	}

	wire := struct {
		ID         interface{} `json:"id,omitempty"`
		Label      string      `json:"label,omitempty"`
		Color      *string     `json:"color,omitempty"`
		UpdateTime *time.Time  `json:"update_time,omitempty"`
		AddTime    *time.Time  `json:"add_time,omitempty"`
	}{
		ID:         id,
		Label:      o.Label,
		Color:      o.Color,
		UpdateTime: o.UpdateTime,
		AddTime:    o.AddTime,
	}
	data, err := json.Marshal(wire)
	if err != nil {
		return nil, fmt.Errorf("v2.FieldOption: encode: %w", err)
	}
	return data, nil
}

type EntityLabel struct {
	ID    int    `json:"id"`
	Label string `json:"label,omitempty"`
}

type FieldSubfield struct {
	FieldCode string    `json:"field_code,omitempty"`
	FieldName string    `json:"field_name,omitempty"`
	FieldType FieldType `json:"field_type,omitempty"`
}

type Field struct {
	FieldName               string                 `json:"field_name,omitempty"`
	FieldCode               string                 `json:"field_code,omitempty"`
	Description             string                 `json:"description,omitempty"`
	FieldType               FieldType              `json:"field_type,omitempty"`
	Options                 []FieldOption          `json:"options,omitempty"`
	Subfields               []FieldSubfield        `json:"subfields,omitempty"`
	IsCustomField           bool                   `json:"is_custom_field,omitempty"`
	IsOptionalResponseField bool                   `json:"is_optional_response_field,omitempty"`
	UIVisibility            map[string]interface{} `json:"ui_visibility,omitempty"`
	ImportantFields         map[string]interface{} `json:"important_fields,omitempty"`
	RequiredFields          map[string]interface{} `json:"required_fields,omitempty"`
}

type optionalSlice[T any] struct {
	value []T
	set   bool
}

func (o *optionalSlice[T]) append(values ...T) {
	if !o.set {
		o.value = make([]T, 0, len(values))
		o.set = true
	}
	o.value = append(o.value, values...)
}

type nullableValue[T any] struct {
	value *T
	set   bool
}

func (o *nullableValue[T]) assign(value T) {
	o.value = &value
	o.set = true
}

func (o *nullableValue[T]) clear() {
	o.value = nil
	o.set = true
}

type FieldOptionUpdate struct {
	ID    int
	Label string
}

type fieldOptionInput struct {
	Label string `json:"label"`
}

type fieldPayload struct {
	name            *string
	fieldType       *FieldType
	description     *string
	options         []fieldOptionInput
	uiVisibility    map[string]interface{}
	importantFields map[string]interface{}
	requiredFields  map[string]interface{}
}

func (p *fieldPayload) addOptions(labels ...string) {
	for _, label := range labels {
		if label == "" {
			continue
		}
		p.options = append(p.options, fieldOptionInput{Label: label})
	}
}

func (p fieldPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["field_name"] = *p.name
	}
	if p.fieldType != nil {
		body["field_type"] = string(*p.fieldType)
	}
	if p.description != nil {
		body["description"] = *p.description
	}
	if len(p.options) > 0 {
		body["options"] = p.options
	}
	if p.uiVisibility != nil {
		body["ui_visibility"] = p.uiVisibility
	}
	if p.importantFields != nil {
		body["important_fields"] = p.importantFields
	}
	if p.requiredFields != nil {
		body["required_fields"] = p.requiredFields
	}
	return body
}

func readFieldResponseBody(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("read field response: missing HTTP response body")
	}

	body, readErr := io.ReadAll(resp.Body)
	closeErr := resp.Body.Close()
	switch {
	case readErr != nil && closeErr != nil:
		return nil, errors.Join(
			fmt.Errorf("read field response body: %w", readErr),
			fmt.Errorf("close field response body: %w", closeErr),
		)
	case readErr != nil:
		return nil, fmt.Errorf("read field response body: %w", readErr)
	default:
		// A successful read preserves the payload needed for status-derived errors.
		return body, nil
	}
}
