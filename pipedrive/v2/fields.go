package v2

import "time"

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
	ID         int        `json:"id,omitempty"`
	Label      string     `json:"label,omitempty"`
	Color      *string    `json:"color,omitempty"`
	UpdateTime *time.Time `json:"update_time,omitempty"`
	AddTime    *time.Time `json:"add_time,omitempty"`
}

type FieldSubfield struct {
	FieldCode string    `json:"field_code,omitempty"`
	FieldName string    `json:"field_name,omitempty"`
	FieldType FieldType `json:"field_type,omitempty"`
}

type Field struct {
	FieldName               string          `json:"field_name,omitempty"`
	FieldCode               string          `json:"field_code,omitempty"`
	Description             string          `json:"description,omitempty"`
	FieldType               FieldType       `json:"field_type,omitempty"`
	Options                 []FieldOption   `json:"options,omitempty"`
	Subfields               []FieldSubfield `json:"subfields,omitempty"`
	IsCustomField           bool            `json:"is_custom_field,omitempty"`
	IsOptionalResponseField bool            `json:"is_optional_response_field,omitempty"`
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
