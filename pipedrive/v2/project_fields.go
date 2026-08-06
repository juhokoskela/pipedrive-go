package v2

import (
	"context"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ProjectFieldsService struct{ client *Client }

type ListProjectFieldsOption interface {
	applyListProjectFields(*listProjectFieldsOptions)
}
type CreateProjectFieldOption interface {
	applyCreateProjectField(*createProjectFieldOptions)
}
type UpdateProjectFieldOption interface {
	applyUpdateProjectField(*updateProjectFieldOptions)
}

type ProjectFieldOption interface {
	CreateProjectFieldOption
	UpdateProjectFieldOption
}

type ProjectFieldRequestOption interface {
	ListProjectFieldsOption
	CreateProjectFieldOption
	UpdateProjectFieldOption
	projectFieldRequestOptions() []pipedrive.RequestOption
}

type listProjectFieldsOptions struct {
	params         genv2.GetProjectFieldsParams
	requestOptions []pipedrive.RequestOption
}
type createProjectFieldOptions struct {
	payload        projectCustomFieldPayload
	requestOptions []pipedrive.RequestOption
}
type updateProjectFieldOptions struct {
	payload        projectCustomFieldPayload
	requestOptions []pipedrive.RequestOption
}
type projectCustomFieldPayload struct {
	name            *string
	fieldType       *FieldType
	options         []fieldOptionInput
	uiVisibility    map[string]interface{}
	importantFields map[string]interface{}
	requiredFields  map[string]interface{}
}

type listProjectFieldsOptionFunc func(*listProjectFieldsOptions)

func (f listProjectFieldsOptionFunc) applyListProjectFields(cfg *listProjectFieldsOptions) { f(cfg) }

type projectCustomFieldOption func(*projectCustomFieldPayload)

func (f projectCustomFieldOption) applyCreateProjectField(cfg *createProjectFieldOptions) {
	f(&cfg.payload)
}
func (f projectCustomFieldOption) applyUpdateProjectField(cfg *updateProjectFieldOptions) {
	f(&cfg.payload)
}

type projectCustomFieldCreateOption func(*projectCustomFieldPayload)

func (f projectCustomFieldCreateOption) applyCreateProjectField(cfg *createProjectFieldOptions) {
	f(&cfg.payload)
}

type projectFieldRequestOption struct{ options []pipedrive.RequestOption }

func (o projectFieldRequestOption) projectFieldRequestOptions() []pipedrive.RequestOption {
	return o.options
}
func (o projectFieldRequestOption) applyListProjectFields(cfg *listProjectFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o projectFieldRequestOption) applyCreateProjectField(cfg *createProjectFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o projectFieldRequestOption) applyUpdateProjectField(cfg *updateProjectFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}

func WithProjectFieldRequestOptions(opts ...pipedrive.RequestOption) ProjectFieldRequestOption {
	return projectFieldRequestOption{options: opts}
}
func WithProjectFieldsPageSize(limit int) ListProjectFieldsOption {
	return listProjectFieldsOptionFunc(func(cfg *listProjectFieldsOptions) {
		if limit > 0 {
			cfg.params.Limit = &limit
		}
	})
}
func WithProjectFieldsCursor(cursor string) ListProjectFieldsOption {
	return listProjectFieldsOptionFunc(func(cfg *listProjectFieldsOptions) {
		if cursor != "" {
			cfg.params.Cursor = &cursor
		}
	})
}
func WithProjectFieldName(name string) ProjectFieldOption {
	return projectCustomFieldOption(func(p *projectCustomFieldPayload) { p.name = &name })
}
func WithProjectFieldType(fieldType FieldType) CreateProjectFieldOption {
	return projectCustomFieldCreateOption(func(p *projectCustomFieldPayload) { p.fieldType = &fieldType })
}
func WithProjectFieldOptions(labels ...string) CreateProjectFieldOption {
	return projectCustomFieldCreateOption(func(p *projectCustomFieldPayload) {
		for _, label := range labels {
			if label != "" {
				p.options = append(p.options, fieldOptionInput{Label: label})
			}
		}
	})
}
func WithProjectFieldUIVisibility(value map[string]interface{}) ProjectFieldOption {
	return projectCustomFieldOption(func(p *projectCustomFieldPayload) { p.uiVisibility = value })
}
func WithProjectFieldImportantFields(value map[string]interface{}) ProjectFieldOption {
	return projectCustomFieldOption(func(p *projectCustomFieldPayload) { p.importantFields = value })
}
func WithProjectFieldRequiredFields(value map[string]interface{}) ProjectFieldOption {
	return projectCustomFieldOption(func(p *projectCustomFieldPayload) { p.requiredFields = value })
}

func newListProjectFieldsOptions(opts []ListProjectFieldsOption) listProjectFieldsOptions {
	var cfg listProjectFieldsOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyListProjectFields(&cfg)
		}
	}
	return cfg
}
func newCreateProjectFieldOptions(opts []CreateProjectFieldOption) createProjectFieldOptions {
	var cfg createProjectFieldOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyCreateProjectField(&cfg)
		}
	}
	return cfg
}
func newUpdateProjectFieldOptions(opts []UpdateProjectFieldOption) updateProjectFieldOptions {
	var cfg updateProjectFieldOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyUpdateProjectField(&cfg)
		}
	}
	return cfg
}
func projectFieldRequestOptionValues(opts []ProjectFieldRequestOption) []pipedrive.RequestOption {
	var out []pipedrive.RequestOption
	for _, opt := range opts {
		if opt != nil {
			out = append(out, opt.projectFieldRequestOptions()...)
		}
	}
	return out
}
func (p projectCustomFieldPayload) body() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["field_name"] = *p.name
	}
	if p.fieldType != nil {
		body["field_type"] = string(*p.fieldType)
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

func (s *ProjectFieldsService) List(ctx context.Context, opts ...ListProjectFieldsOption) ([]Field, *string, error) {
	cfg := newListProjectFieldsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}
func (s *ProjectFieldsService) ListPager(opts ...ListProjectFieldsOption) *pipedrive.CursorPager[Field] {
	cfg := newListProjectFieldsOptions(opts)
	start := cfg.params.Cursor
	cfg.params.Cursor = nil
	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Field, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if start != nil {
			params.Cursor = start
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}
func (s *ProjectFieldsService) ForEach(ctx context.Context, fn func(Field) error, opts ...ListProjectFieldsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *ProjectFieldsService) Get(ctx context.Context, fieldCode string, opts ...ProjectFieldRequestOption) (*Field, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectFieldRequestOptionValues(opts)...)
	resp, err := s.client.gen.GetProjectField(ctx, fieldCode, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	responseBody, err := readFieldResponseBody(resp)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Field](resp, responseBody, "project field")
}
func (s *ProjectFieldsService) Create(ctx context.Context, opts ...CreateProjectFieldOption) (*Field, error) {
	cfg := newCreateProjectFieldOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.AddProjectFieldWithBody(ctx, "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	responseBody, err := readFieldResponseBody(resp)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Field](resp, responseBody, "project field")
}
func (s *ProjectFieldsService) Update(ctx context.Context, fieldCode string, opts ...UpdateProjectFieldOption) (*Field, error) {
	cfg := newUpdateProjectFieldOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.UpdateProjectFieldWithBody(ctx, fieldCode, "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	responseBody, err := readFieldResponseBody(resp)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Field](resp, responseBody, "project field")
}
func (s *ProjectFieldsService) Delete(ctx context.Context, fieldCode string, opts ...ProjectFieldRequestOption) (*Field, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectFieldRequestOptionValues(opts)...)
	resp, err := s.client.gen.DeleteProjectField(ctx, fieldCode, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	responseBody, err := readFieldResponseBody(resp)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Field](resp, responseBody, "project field delete")
}
func (s *ProjectFieldsService) AddOptions(ctx context.Context, fieldCode string, labels []string, opts ...ProjectFieldRequestOption) ([]FieldOption, error) {
	bodyItems := make([]map[string]interface{}, 0, len(labels))
	for _, label := range labels {
		if label != "" {
			bodyItems = append(bodyItems, map[string]interface{}{"label": label})
		}
	}
	body, err := encodeV2ArrayBody(bodyItems)
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectFieldRequestOptionValues(opts)...)
	resp, err := s.client.gen.AddProjectFieldOptionsWithBody(ctx, fieldCode, "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	responseBody, err := readFieldResponseBody(resp)
	if err != nil {
		return nil, err
	}
	return decodeV2ListNoCursor[FieldOption](resp, responseBody)
}
func (s *ProjectFieldsService) UpdateOptions(ctx context.Context, fieldCode string, updates []FieldOptionUpdate, opts ...ProjectFieldRequestOption) ([]FieldOption, error) {
	bodyItems := make([]map[string]interface{}, 0, len(updates))
	for _, update := range updates {
		bodyItems = append(bodyItems, map[string]interface{}{"id": update.ID, "label": update.Label})
	}
	body, err := encodeV2ArrayBody(bodyItems)
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectFieldRequestOptionValues(opts)...)
	resp, err := s.client.gen.UpdateProjectFieldOptionsWithBody(ctx, fieldCode, "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	responseBody, err := readFieldResponseBody(resp)
	if err != nil {
		return nil, err
	}
	return decodeV2ListNoCursor[FieldOption](resp, responseBody)
}
func (s *ProjectFieldsService) DeleteOptions(ctx context.Context, fieldCode string, ids []int, opts ...ProjectFieldRequestOption) ([]FieldOption, error) {
	bodyItems := make([]map[string]interface{}, 0, len(ids))
	for _, id := range ids {
		bodyItems = append(bodyItems, map[string]interface{}{"id": id})
	}
	body, err := encodeV2ArrayBody(bodyItems)
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectFieldRequestOptionValues(opts)...)
	resp, err := s.client.gen.DeleteProjectFieldOptionsWithBody(ctx, fieldCode, "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	responseBody, err := readFieldResponseBody(resp)
	if err != nil {
		return nil, err
	}
	return decodeV2ListNoCursor[FieldOption](resp, responseBody)
}
func (s *ProjectFieldsService) list(ctx context.Context, params genv2.GetProjectFieldsParams, requestOptions []pipedrive.RequestOption) ([]Field, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)
	resp, err := s.client.gen.GetProjectFields(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	responseBody, err := readFieldResponseBody(resp)
	if err != nil {
		return nil, nil, err
	}
	return decodeV2List[Field](resp, responseBody)
}
