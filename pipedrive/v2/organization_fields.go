package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type OrganizationFieldsService struct {
	client *Client
}

type GetOrganizationFieldOption interface {
	applyGetOrganizationField(*getOrganizationFieldOptions)
}

type ListOrganizationFieldsOption interface {
	applyListOrganizationFields(*listOrganizationFieldsOptions)
}

type CreateOrganizationFieldOption interface {
	applyCreateOrganizationField(*createOrganizationFieldOptions)
}

type UpdateOrganizationFieldOption interface {
	applyUpdateOrganizationField(*updateOrganizationFieldOptions)
}

type DeleteOrganizationFieldOption interface {
	applyDeleteOrganizationField(*deleteOrganizationFieldOptions)
}

type AddOrganizationFieldOptionsOption interface {
	applyAddOrganizationFieldOptions(*addOrganizationFieldOptionItemsOptions)
}

type UpdateOrganizationFieldOptionsOption interface {
	applyUpdateOrganizationFieldOptions(*updateOrganizationFieldOptionItemsOptions)
}

type DeleteOrganizationFieldOptionsOption interface {
	applyDeleteOrganizationFieldOptions(*deleteOrganizationFieldOptionItemsOptions)
}

type OrganizationFieldRequestOption interface {
	GetOrganizationFieldOption
	ListOrganizationFieldsOption
	CreateOrganizationFieldOption
	UpdateOrganizationFieldOption
	DeleteOrganizationFieldOption
	AddOrganizationFieldOptionsOption
	UpdateOrganizationFieldOptionsOption
	DeleteOrganizationFieldOptionsOption
}

type OrganizationFieldOption interface {
	CreateOrganizationFieldOption
	UpdateOrganizationFieldOption
}

type getOrganizationFieldOptions struct {
	params         genv2.GetOrganizationFieldParams
	requestOptions []pipedrive.RequestOption
}

type listOrganizationFieldsOptions struct {
	params         genv2.GetOrganizationFieldsParams
	requestOptions []pipedrive.RequestOption
}

type createOrganizationFieldOptions struct {
	payload        fieldPayload
	requestOptions []pipedrive.RequestOption
}

type updateOrganizationFieldOptions struct {
	payload        fieldPayload
	requestOptions []pipedrive.RequestOption
}

type deleteOrganizationFieldOptions struct {
	requestOptions []pipedrive.RequestOption
}

type addOrganizationFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type updateOrganizationFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deleteOrganizationFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type organizationFieldRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o organizationFieldRequestOptions) applyGetOrganizationField(cfg *getOrganizationFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationFieldRequestOptions) applyListOrganizationFields(cfg *listOrganizationFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationFieldRequestOptions) applyCreateOrganizationField(cfg *createOrganizationFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationFieldRequestOptions) applyUpdateOrganizationField(cfg *updateOrganizationFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationFieldRequestOptions) applyDeleteOrganizationField(cfg *deleteOrganizationFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationFieldRequestOptions) applyAddOrganizationFieldOptions(cfg *addOrganizationFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationFieldRequestOptions) applyUpdateOrganizationFieldOptions(cfg *updateOrganizationFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationFieldRequestOptions) applyDeleteOrganizationFieldOptions(cfg *deleteOrganizationFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getOrganizationFieldOptionFunc func(*getOrganizationFieldOptions)

func (f getOrganizationFieldOptionFunc) applyGetOrganizationField(cfg *getOrganizationFieldOptions) {
	f(cfg)
}

type listOrganizationFieldsOptionFunc func(*listOrganizationFieldsOptions)

func (f listOrganizationFieldsOptionFunc) applyListOrganizationFields(cfg *listOrganizationFieldsOptions) {
	f(cfg)
}

type organizationFieldOptionFunc func(*fieldPayload)

func (f organizationFieldOptionFunc) applyCreateOrganizationField(cfg *createOrganizationFieldOptions) {
	f(&cfg.payload)
}

func (f organizationFieldOptionFunc) applyUpdateOrganizationField(cfg *updateOrganizationFieldOptions) {
	f(&cfg.payload)
}

func WithOrganizationFieldRequestOptions(opts ...pipedrive.RequestOption) OrganizationFieldRequestOption {
	return organizationFieldRequestOptions{requestOptions: opts}
}

func WithOrganizationFieldIncludeFields(fields ...FieldIncludeField) GetOrganizationFieldOption {
	return getOrganizationFieldOptionFunc(func(cfg *getOrganizationFieldOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetOrganizationFieldParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithOrganizationFieldsIncludeFields(fields ...FieldIncludeField) ListOrganizationFieldsOption {
	return listOrganizationFieldsOptionFunc(func(cfg *listOrganizationFieldsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetOrganizationFieldsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithOrganizationFieldsPageSize(limit int) ListOrganizationFieldsOption {
	return listOrganizationFieldsOptionFunc(func(cfg *listOrganizationFieldsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithOrganizationFieldsCursor(cursor string) ListOrganizationFieldsOption {
	return listOrganizationFieldsOptionFunc(func(cfg *listOrganizationFieldsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithOrganizationFieldName(name string) OrganizationFieldOption {
	return organizationFieldOptionFunc(func(payload *fieldPayload) {
		payload.name = &name
	})
}

func WithOrganizationFieldType(fieldType FieldType) OrganizationFieldOption {
	return organizationFieldOptionFunc(func(payload *fieldPayload) {
		payload.fieldType = &fieldType
	})
}

func WithOrganizationFieldDescription(description string) OrganizationFieldOption {
	return organizationFieldOptionFunc(func(payload *fieldPayload) {
		payload.description = &description
	})
}

func WithOrganizationFieldOptions(labels ...string) OrganizationFieldOption {
	return organizationFieldOptionFunc(func(payload *fieldPayload) {
		payload.addOptions(labels...)
	})
}

func WithOrganizationFieldUIVisibility(value map[string]interface{}) OrganizationFieldOption {
	return organizationFieldOptionFunc(func(payload *fieldPayload) {
		payload.uiVisibility = value
	})
}

func WithOrganizationFieldImportantFields(value map[string]interface{}) OrganizationFieldOption {
	return organizationFieldOptionFunc(func(payload *fieldPayload) {
		payload.importantFields = value
	})
}

func WithOrganizationFieldRequiredFields(value map[string]interface{}) OrganizationFieldOption {
	return organizationFieldOptionFunc(func(payload *fieldPayload) {
		payload.requiredFields = value
	})
}

func newGetOrganizationFieldOptions(opts []GetOrganizationFieldOption) getOrganizationFieldOptions {
	var cfg getOrganizationFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetOrganizationField(&cfg)
	}
	return cfg
}

func newListOrganizationFieldsOptions(opts []ListOrganizationFieldsOption) listOrganizationFieldsOptions {
	var cfg listOrganizationFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListOrganizationFields(&cfg)
	}
	return cfg
}

func newCreateOrganizationFieldOptions(opts []CreateOrganizationFieldOption) createOrganizationFieldOptions {
	var cfg createOrganizationFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateOrganizationField(&cfg)
	}
	return cfg
}

func newUpdateOrganizationFieldOptions(opts []UpdateOrganizationFieldOption) updateOrganizationFieldOptions {
	var cfg updateOrganizationFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateOrganizationField(&cfg)
	}
	return cfg
}

func newDeleteOrganizationFieldOptions(opts []DeleteOrganizationFieldOption) deleteOrganizationFieldOptions {
	var cfg deleteOrganizationFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteOrganizationField(&cfg)
	}
	return cfg
}

func newAddOrganizationFieldOptions(opts []AddOrganizationFieldOptionsOption) addOrganizationFieldOptionItemsOptions {
	var cfg addOrganizationFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddOrganizationFieldOptions(&cfg)
	}
	return cfg
}

func newUpdateOrganizationFieldOptionsOptions(opts []UpdateOrganizationFieldOptionsOption) updateOrganizationFieldOptionItemsOptions {
	var cfg updateOrganizationFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateOrganizationFieldOptions(&cfg)
	}
	return cfg
}

func newDeleteOrganizationFieldOptionsOptions(opts []DeleteOrganizationFieldOptionsOption) deleteOrganizationFieldOptionItemsOptions {
	var cfg deleteOrganizationFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteOrganizationFieldOptions(&cfg)
	}
	return cfg
}

func (s *OrganizationFieldsService) Get(ctx context.Context, fieldCode string, opts ...GetOrganizationFieldOption) (*Field, error) {
	cfg := newGetOrganizationFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetOrganizationFieldWithResponse(ctx, fieldCode, &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Field `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization field data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationFieldsService) List(ctx context.Context, opts ...ListOrganizationFieldsOption) ([]Field, *string, error) {
	cfg := newListOrganizationFieldsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *OrganizationFieldsService) ListPager(opts ...ListOrganizationFieldsOption) *pipedrive.CursorPager[Field] {
	cfg := newListOrganizationFieldsOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Field, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *OrganizationFieldsService) ForEach(ctx context.Context, fn func(Field) error, opts ...ListOrganizationFieldsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *OrganizationFieldsService) Create(ctx context.Context, opts ...CreateOrganizationFieldOption) (*Field, error) {
	cfg := newCreateOrganizationFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddOrganizationFieldWithBodyWithResponse(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Field `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization field data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationFieldsService) Update(ctx context.Context, fieldCode string, opts ...UpdateOrganizationFieldOption) (*Field, error) {
	cfg := newUpdateOrganizationFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateOrganizationFieldWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Field `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization field data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationFieldsService) Delete(ctx context.Context, fieldCode string, opts ...DeleteOrganizationFieldOption) (*Field, error) {
	cfg := newDeleteOrganizationFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteOrganizationFieldWithResponse(ctx, fieldCode, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Field `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization field data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationFieldsService) AddOptions(ctx context.Context, fieldCode string, labels []string, opts ...AddOrganizationFieldOptionsOption) ([]FieldOption, error) {
	cfg := newAddOrganizationFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	payload := make([]fieldOptionInput, 0, len(labels))
	for _, label := range labels {
		if label == "" {
			continue
		}
		payload = append(payload, fieldOptionInput{Label: label})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddOrganizationFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payloadResp struct {
		Data []FieldOption `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payloadResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payloadResp.Data, nil
}

func (s *OrganizationFieldsService) UpdateOptions(ctx context.Context, fieldCode string, updates []FieldOptionUpdate, opts ...UpdateOrganizationFieldOptionsOption) ([]FieldOption, error) {
	cfg := newUpdateOrganizationFieldOptionsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	payload := make([]struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	}, 0, len(updates))
	for _, update := range updates {
		payload = append(payload, struct {
			ID    int    `json:"id"`
			Label string `json:"label"`
		}{
			ID:    update.ID,
			Label: update.Label,
		})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateOrganizationFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payloadResp struct {
		Data []FieldOption `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payloadResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payloadResp.Data, nil
}

func (s *OrganizationFieldsService) DeleteOptions(ctx context.Context, fieldCode string, ids []int, opts ...DeleteOrganizationFieldOptionsOption) ([]FieldOption, error) {
	cfg := newDeleteOrganizationFieldOptionsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	payload := make([]struct {
		ID int `json:"id"`
	}, 0, len(ids))
	for _, id := range ids {
		payload = append(payload, struct {
			ID int `json:"id"`
		}{
			ID: id,
		})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.DeleteOrganizationFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payloadResp struct {
		Data []FieldOption `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payloadResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payloadResp.Data, nil
}

func (s *OrganizationFieldsService) list(ctx context.Context, params genv2.GetOrganizationFieldsParams, requestOptions []pipedrive.RequestOption) ([]Field, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetOrganizationFieldsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Field `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}
