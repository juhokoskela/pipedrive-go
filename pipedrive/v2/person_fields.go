package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type PersonFieldsService struct {
	client *Client
}

type GetPersonFieldOption interface {
	applyGetPersonField(*getPersonFieldOptions)
}

type ListPersonFieldsOption interface {
	applyListPersonFields(*listPersonFieldsOptions)
}

type CreatePersonFieldOption interface {
	applyCreatePersonField(*createPersonFieldOptions)
}

type UpdatePersonFieldOption interface {
	applyUpdatePersonField(*updatePersonFieldOptions)
}

type DeletePersonFieldOption interface {
	applyDeletePersonField(*deletePersonFieldOptions)
}

type AddPersonFieldOptionsOption interface {
	applyAddPersonFieldOptions(*addPersonFieldOptionItemsOptions)
}

type UpdatePersonFieldOptionsOption interface {
	applyUpdatePersonFieldOptions(*updatePersonFieldOptionItemsOptions)
}

type DeletePersonFieldOptionsOption interface {
	applyDeletePersonFieldOptions(*deletePersonFieldOptionItemsOptions)
}

type PersonFieldRequestOption interface {
	GetPersonFieldOption
	ListPersonFieldsOption
	CreatePersonFieldOption
	UpdatePersonFieldOption
	DeletePersonFieldOption
	AddPersonFieldOptionsOption
	UpdatePersonFieldOptionsOption
	DeletePersonFieldOptionsOption
}

type PersonFieldOption interface {
	CreatePersonFieldOption
	UpdatePersonFieldOption
}

type getPersonFieldOptions struct {
	params         genv2.GetPersonFieldParams
	requestOptions []pipedrive.RequestOption
}

type listPersonFieldsOptions struct {
	params         genv2.GetPersonFieldsParams
	requestOptions []pipedrive.RequestOption
}

type createPersonFieldOptions struct {
	payload        fieldPayload
	requestOptions []pipedrive.RequestOption
}

type updatePersonFieldOptions struct {
	payload        fieldPayload
	requestOptions []pipedrive.RequestOption
}

type deletePersonFieldOptions struct {
	requestOptions []pipedrive.RequestOption
}

type addPersonFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type updatePersonFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deletePersonFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type personFieldRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o personFieldRequestOptions) applyGetPersonField(cfg *getPersonFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personFieldRequestOptions) applyListPersonFields(cfg *listPersonFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personFieldRequestOptions) applyCreatePersonField(cfg *createPersonFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personFieldRequestOptions) applyUpdatePersonField(cfg *updatePersonFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personFieldRequestOptions) applyDeletePersonField(cfg *deletePersonFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personFieldRequestOptions) applyAddPersonFieldOptions(cfg *addPersonFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personFieldRequestOptions) applyUpdatePersonFieldOptions(cfg *updatePersonFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personFieldRequestOptions) applyDeletePersonFieldOptions(cfg *deletePersonFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getPersonFieldOptionFunc func(*getPersonFieldOptions)

func (f getPersonFieldOptionFunc) applyGetPersonField(cfg *getPersonFieldOptions) {
	f(cfg)
}

type listPersonFieldsOptionFunc func(*listPersonFieldsOptions)

func (f listPersonFieldsOptionFunc) applyListPersonFields(cfg *listPersonFieldsOptions) {
	f(cfg)
}

type personFieldOptionFunc func(*fieldPayload)

func (f personFieldOptionFunc) applyCreatePersonField(cfg *createPersonFieldOptions) {
	f(&cfg.payload)
}

func (f personFieldOptionFunc) applyUpdatePersonField(cfg *updatePersonFieldOptions) {
	f(&cfg.payload)
}

func WithPersonFieldRequestOptions(opts ...pipedrive.RequestOption) PersonFieldRequestOption {
	return personFieldRequestOptions{requestOptions: opts}
}

func WithPersonFieldIncludeFields(fields ...FieldIncludeField) GetPersonFieldOption {
	return getPersonFieldOptionFunc(func(cfg *getPersonFieldOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetPersonFieldParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithPersonFieldsIncludeFields(fields ...FieldIncludeField) ListPersonFieldsOption {
	return listPersonFieldsOptionFunc(func(cfg *listPersonFieldsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetPersonFieldsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithPersonFieldsPageSize(limit int) ListPersonFieldsOption {
	return listPersonFieldsOptionFunc(func(cfg *listPersonFieldsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithPersonFieldsCursor(cursor string) ListPersonFieldsOption {
	return listPersonFieldsOptionFunc(func(cfg *listPersonFieldsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithPersonFieldName(name string) PersonFieldOption {
	return personFieldOptionFunc(func(payload *fieldPayload) {
		payload.name = &name
	})
}

func WithPersonFieldType(fieldType FieldType) PersonFieldOption {
	return personFieldOptionFunc(func(payload *fieldPayload) {
		payload.fieldType = &fieldType
	})
}

func WithPersonFieldDescription(description string) PersonFieldOption {
	return personFieldOptionFunc(func(payload *fieldPayload) {
		payload.description = &description
	})
}

func WithPersonFieldOptions(labels ...string) PersonFieldOption {
	return personFieldOptionFunc(func(payload *fieldPayload) {
		payload.addOptions(labels...)
	})
}

func WithPersonFieldUIVisibility(value map[string]interface{}) PersonFieldOption {
	return personFieldOptionFunc(func(payload *fieldPayload) {
		payload.uiVisibility = value
	})
}

func WithPersonFieldImportantFields(value map[string]interface{}) PersonFieldOption {
	return personFieldOptionFunc(func(payload *fieldPayload) {
		payload.importantFields = value
	})
}

func WithPersonFieldRequiredFields(value map[string]interface{}) PersonFieldOption {
	return personFieldOptionFunc(func(payload *fieldPayload) {
		payload.requiredFields = value
	})
}

func newGetPersonFieldOptions(opts []GetPersonFieldOption) getPersonFieldOptions {
	var cfg getPersonFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetPersonField(&cfg)
	}
	return cfg
}

func newListPersonFieldsOptions(opts []ListPersonFieldsOption) listPersonFieldsOptions {
	var cfg listPersonFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListPersonFields(&cfg)
	}
	return cfg
}

func newCreatePersonFieldOptions(opts []CreatePersonFieldOption) createPersonFieldOptions {
	var cfg createPersonFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreatePersonField(&cfg)
	}
	return cfg
}

func newUpdatePersonFieldOptions(opts []UpdatePersonFieldOption) updatePersonFieldOptions {
	var cfg updatePersonFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdatePersonField(&cfg)
	}
	return cfg
}

func newDeletePersonFieldOptions(opts []DeletePersonFieldOption) deletePersonFieldOptions {
	var cfg deletePersonFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeletePersonField(&cfg)
	}
	return cfg
}

func newAddPersonFieldOptions(opts []AddPersonFieldOptionsOption) addPersonFieldOptionItemsOptions {
	var cfg addPersonFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddPersonFieldOptions(&cfg)
	}
	return cfg
}

func newUpdatePersonFieldOptionsOptions(opts []UpdatePersonFieldOptionsOption) updatePersonFieldOptionItemsOptions {
	var cfg updatePersonFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdatePersonFieldOptions(&cfg)
	}
	return cfg
}

func newDeletePersonFieldOptionsOptions(opts []DeletePersonFieldOptionsOption) deletePersonFieldOptionItemsOptions {
	var cfg deletePersonFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeletePersonFieldOptions(&cfg)
	}
	return cfg
}

func (s *PersonFieldsService) Get(ctx context.Context, fieldCode string, opts ...GetPersonFieldOption) (*Field, error) {
	cfg := newGetPersonFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetPersonFieldWithResponse(ctx, fieldCode, &cfg.params, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing person field data in response")
	}
	return payload.Data, nil
}

func (s *PersonFieldsService) List(ctx context.Context, opts ...ListPersonFieldsOption) ([]Field, *string, error) {
	cfg := newListPersonFieldsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *PersonFieldsService) ListPager(opts ...ListPersonFieldsOption) *pipedrive.CursorPager[Field] {
	cfg := newListPersonFieldsOptions(opts)
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

func (s *PersonFieldsService) ForEach(ctx context.Context, fn func(Field) error, opts ...ListPersonFieldsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *PersonFieldsService) Create(ctx context.Context, opts ...CreatePersonFieldOption) (*Field, error) {
	cfg := newCreatePersonFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddPersonFieldWithBodyWithResponse(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing person field data in response")
	}
	return payload.Data, nil
}

func (s *PersonFieldsService) Update(ctx context.Context, fieldCode string, opts ...UpdatePersonFieldOption) (*Field, error) {
	cfg := newUpdatePersonFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdatePersonFieldWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing person field data in response")
	}
	return payload.Data, nil
}

func (s *PersonFieldsService) Delete(ctx context.Context, fieldCode string, opts ...DeletePersonFieldOption) (*Field, error) {
	cfg := newDeletePersonFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeletePersonFieldWithResponse(ctx, fieldCode, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing person field data in response")
	}
	return payload.Data, nil
}

func (s *PersonFieldsService) AddOptions(ctx context.Context, fieldCode string, labels []string, opts ...AddPersonFieldOptionsOption) ([]FieldOption, error) {
	cfg := newAddPersonFieldOptions(opts)
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

	resp, err := s.client.gen.AddPersonFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *PersonFieldsService) UpdateOptions(ctx context.Context, fieldCode string, updates []FieldOptionUpdate, opts ...UpdatePersonFieldOptionsOption) ([]FieldOption, error) {
	cfg := newUpdatePersonFieldOptionsOptions(opts)
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

	resp, err := s.client.gen.UpdatePersonFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *PersonFieldsService) DeleteOptions(ctx context.Context, fieldCode string, ids []int, opts ...DeletePersonFieldOptionsOption) ([]FieldOption, error) {
	cfg := newDeletePersonFieldOptionsOptions(opts)
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

	resp, err := s.client.gen.DeletePersonFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *PersonFieldsService) list(ctx context.Context, params genv2.GetPersonFieldsParams, requestOptions []pipedrive.RequestOption) ([]Field, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetPersonFieldsWithResponse(ctx, &params, toRequestEditors(editors)...)
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
