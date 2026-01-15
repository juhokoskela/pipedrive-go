package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type DealFieldsService struct {
	client *Client
}

type GetDealFieldOption interface {
	applyGetDealField(*getDealFieldOptions)
}

type ListDealFieldsOption interface {
	applyListDealFields(*listDealFieldsOptions)
}

type CreateDealFieldOption interface {
	applyCreateDealField(*createDealFieldOptions)
}

type UpdateDealFieldOption interface {
	applyUpdateDealField(*updateDealFieldOptions)
}

type DeleteDealFieldOption interface {
	applyDeleteDealField(*deleteDealFieldOptions)
}

type AddDealFieldOptionsOption interface {
	applyAddDealFieldOptions(*addDealFieldOptionItemsOptions)
}

type UpdateDealFieldOptionsOption interface {
	applyUpdateDealFieldOptions(*updateDealFieldOptionItemsOptions)
}

type DeleteDealFieldOptionsOption interface {
	applyDeleteDealFieldOptions(*deleteDealFieldOptionItemsOptions)
}

type DealFieldRequestOption interface {
	GetDealFieldOption
	ListDealFieldsOption
	CreateDealFieldOption
	UpdateDealFieldOption
	DeleteDealFieldOption
	AddDealFieldOptionsOption
	UpdateDealFieldOptionsOption
	DeleteDealFieldOptionsOption
}

type DealFieldOption interface {
	CreateDealFieldOption
	UpdateDealFieldOption
}

type getDealFieldOptions struct {
	params         genv2.GetDealFieldParams
	requestOptions []pipedrive.RequestOption
}

type listDealFieldsOptions struct {
	params         genv2.GetDealFieldsParams
	requestOptions []pipedrive.RequestOption
}

type createDealFieldOptions struct {
	payload        fieldPayload
	requestOptions []pipedrive.RequestOption
}

type updateDealFieldOptions struct {
	payload        fieldPayload
	requestOptions []pipedrive.RequestOption
}

type deleteDealFieldOptions struct {
	requestOptions []pipedrive.RequestOption
}

type addDealFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type updateDealFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deleteDealFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type dealFieldRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o dealFieldRequestOptions) applyGetDealField(cfg *getDealFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealFieldRequestOptions) applyListDealFields(cfg *listDealFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealFieldRequestOptions) applyCreateDealField(cfg *createDealFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealFieldRequestOptions) applyUpdateDealField(cfg *updateDealFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealFieldRequestOptions) applyDeleteDealField(cfg *deleteDealFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealFieldRequestOptions) applyAddDealFieldOptions(cfg *addDealFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealFieldRequestOptions) applyUpdateDealFieldOptions(cfg *updateDealFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealFieldRequestOptions) applyDeleteDealFieldOptions(cfg *deleteDealFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getDealFieldOptionFunc func(*getDealFieldOptions)

func (f getDealFieldOptionFunc) applyGetDealField(cfg *getDealFieldOptions) {
	f(cfg)
}

type listDealFieldsOptionFunc func(*listDealFieldsOptions)

func (f listDealFieldsOptionFunc) applyListDealFields(cfg *listDealFieldsOptions) {
	f(cfg)
}

type dealFieldOptionFunc func(*fieldPayload)

func (f dealFieldOptionFunc) applyCreateDealField(cfg *createDealFieldOptions) {
	f(&cfg.payload)
}

func (f dealFieldOptionFunc) applyUpdateDealField(cfg *updateDealFieldOptions) {
	f(&cfg.payload)
}

func WithDealFieldRequestOptions(opts ...pipedrive.RequestOption) DealFieldRequestOption {
	return dealFieldRequestOptions{requestOptions: opts}
}

func WithDealFieldIncludeFields(fields ...FieldIncludeField) GetDealFieldOption {
	return getDealFieldOptionFunc(func(cfg *getDealFieldOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetDealFieldParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithDealFieldsIncludeFields(fields ...FieldIncludeField) ListDealFieldsOption {
	return listDealFieldsOptionFunc(func(cfg *listDealFieldsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetDealFieldsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithDealFieldsPageSize(limit int) ListDealFieldsOption {
	return listDealFieldsOptionFunc(func(cfg *listDealFieldsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithDealFieldsCursor(cursor string) ListDealFieldsOption {
	return listDealFieldsOptionFunc(func(cfg *listDealFieldsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithDealFieldName(name string) DealFieldOption {
	return dealFieldOptionFunc(func(payload *fieldPayload) {
		payload.name = &name
	})
}

func WithDealFieldType(fieldType FieldType) DealFieldOption {
	return dealFieldOptionFunc(func(payload *fieldPayload) {
		payload.fieldType = &fieldType
	})
}

func WithDealFieldDescription(description string) DealFieldOption {
	return dealFieldOptionFunc(func(payload *fieldPayload) {
		payload.description = &description
	})
}

func WithDealFieldOptions(labels ...string) DealFieldOption {
	return dealFieldOptionFunc(func(payload *fieldPayload) {
		payload.addOptions(labels...)
	})
}

func WithDealFieldUIVisibility(value map[string]interface{}) DealFieldOption {
	return dealFieldOptionFunc(func(payload *fieldPayload) {
		payload.uiVisibility = value
	})
}

func WithDealFieldImportantFields(value map[string]interface{}) DealFieldOption {
	return dealFieldOptionFunc(func(payload *fieldPayload) {
		payload.importantFields = value
	})
}

func WithDealFieldRequiredFields(value map[string]interface{}) DealFieldOption {
	return dealFieldOptionFunc(func(payload *fieldPayload) {
		payload.requiredFields = value
	})
}

func newGetDealFieldOptions(opts []GetDealFieldOption) getDealFieldOptions {
	var cfg getDealFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetDealField(&cfg)
	}
	return cfg
}

func newListDealFieldsOptions(opts []ListDealFieldsOption) listDealFieldsOptions {
	var cfg listDealFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListDealFields(&cfg)
	}
	return cfg
}

func newCreateDealFieldOptions(opts []CreateDealFieldOption) createDealFieldOptions {
	var cfg createDealFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateDealField(&cfg)
	}
	return cfg
}

func newUpdateDealFieldOptions(opts []UpdateDealFieldOption) updateDealFieldOptions {
	var cfg updateDealFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateDealField(&cfg)
	}
	return cfg
}

func newDeleteDealFieldOptions(opts []DeleteDealFieldOption) deleteDealFieldOptions {
	var cfg deleteDealFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteDealField(&cfg)
	}
	return cfg
}

func newAddDealFieldOptions(opts []AddDealFieldOptionsOption) addDealFieldOptionItemsOptions {
	var cfg addDealFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddDealFieldOptions(&cfg)
	}
	return cfg
}

func newUpdateDealFieldOptionsOptions(opts []UpdateDealFieldOptionsOption) updateDealFieldOptionItemsOptions {
	var cfg updateDealFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateDealFieldOptions(&cfg)
	}
	return cfg
}

func newDeleteDealFieldOptionsOptions(opts []DeleteDealFieldOptionsOption) deleteDealFieldOptionItemsOptions {
	var cfg deleteDealFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteDealFieldOptions(&cfg)
	}
	return cfg
}

func (s *DealFieldsService) Get(ctx context.Context, fieldCode string, opts ...GetDealFieldOption) (*Field, error) {
	cfg := newGetDealFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetDealFieldWithResponse(ctx, fieldCode, &cfg.params, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing deal field data in response")
	}
	return payload.Data, nil
}

func (s *DealFieldsService) List(ctx context.Context, opts ...ListDealFieldsOption) ([]Field, *string, error) {
	cfg := newListDealFieldsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *DealFieldsService) ListPager(opts ...ListDealFieldsOption) *pipedrive.CursorPager[Field] {
	cfg := newListDealFieldsOptions(opts)
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

func (s *DealFieldsService) ForEach(ctx context.Context, fn func(Field) error, opts ...ListDealFieldsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *DealFieldsService) Create(ctx context.Context, opts ...CreateDealFieldOption) (*Field, error) {
	cfg := newCreateDealFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddDealFieldWithBodyWithResponse(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing deal field data in response")
	}
	return payload.Data, nil
}

func (s *DealFieldsService) Update(ctx context.Context, fieldCode string, opts ...UpdateDealFieldOption) (*Field, error) {
	cfg := newUpdateDealFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateDealFieldWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing deal field data in response")
	}
	return payload.Data, nil
}

func (s *DealFieldsService) Delete(ctx context.Context, fieldCode string, opts ...DeleteDealFieldOption) (*Field, error) {
	cfg := newDeleteDealFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteDealFieldWithResponse(ctx, fieldCode, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing deal field data in response")
	}
	return payload.Data, nil
}

func (s *DealFieldsService) AddOptions(ctx context.Context, fieldCode string, labels []string, opts ...AddDealFieldOptionsOption) ([]FieldOption, error) {
	cfg := newAddDealFieldOptions(opts)
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

	resp, err := s.client.gen.AddDealFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *DealFieldsService) UpdateOptions(ctx context.Context, fieldCode string, updates []FieldOptionUpdate, opts ...UpdateDealFieldOptionsOption) ([]FieldOption, error) {
	cfg := newUpdateDealFieldOptionsOptions(opts)
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

	resp, err := s.client.gen.UpdateDealFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *DealFieldsService) DeleteOptions(ctx context.Context, fieldCode string, ids []int, opts ...DeleteDealFieldOptionsOption) ([]FieldOption, error) {
	cfg := newDeleteDealFieldOptionsOptions(opts)
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

	resp, err := s.client.gen.DeleteDealFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *DealFieldsService) list(ctx context.Context, params genv2.GetDealFieldsParams, requestOptions []pipedrive.RequestOption) ([]Field, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetDealFieldsWithResponse(ctx, &params, toRequestEditors(editors)...)
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
