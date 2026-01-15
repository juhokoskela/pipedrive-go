package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ProductFieldsService struct {
	client *Client
}

type GetProductFieldOption interface {
	applyGetProductField(*getProductFieldOptions)
}

type ListProductFieldsOption interface {
	applyListProductFields(*listProductFieldsOptions)
}

type CreateProductFieldOption interface {
	applyCreateProductField(*createProductFieldOptions)
}

type UpdateProductFieldOption interface {
	applyUpdateProductField(*updateProductFieldOptions)
}

type DeleteProductFieldOption interface {
	applyDeleteProductField(*deleteProductFieldOptions)
}

type AddProductFieldOptionsOption interface {
	applyAddProductFieldOptions(*addProductFieldOptionItemsOptions)
}

type UpdateProductFieldOptionsOption interface {
	applyUpdateProductFieldOptions(*updateProductFieldOptionItemsOptions)
}

type DeleteProductFieldOptionsOption interface {
	applyDeleteProductFieldOptions(*deleteProductFieldOptionItemsOptions)
}

type ProductFieldRequestOption interface {
	GetProductFieldOption
	ListProductFieldsOption
	CreateProductFieldOption
	UpdateProductFieldOption
	DeleteProductFieldOption
	AddProductFieldOptionsOption
	UpdateProductFieldOptionsOption
	DeleteProductFieldOptionsOption
}

type ProductFieldOption interface {
	CreateProductFieldOption
	UpdateProductFieldOption
}

type getProductFieldOptions struct {
	params         genv2.GetProductFieldParams
	requestOptions []pipedrive.RequestOption
}

type listProductFieldsOptions struct {
	params         genv2.GetProductFieldsParams
	requestOptions []pipedrive.RequestOption
}

type createProductFieldOptions struct {
	payload        fieldPayload
	requestOptions []pipedrive.RequestOption
}

type updateProductFieldOptions struct {
	payload        fieldPayload
	requestOptions []pipedrive.RequestOption
}

type deleteProductFieldOptions struct {
	requestOptions []pipedrive.RequestOption
}

type addProductFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type updateProductFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deleteProductFieldOptionItemsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type productFieldRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o productFieldRequestOptions) applyGetProductField(cfg *getProductFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productFieldRequestOptions) applyListProductFields(cfg *listProductFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productFieldRequestOptions) applyCreateProductField(cfg *createProductFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productFieldRequestOptions) applyUpdateProductField(cfg *updateProductFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productFieldRequestOptions) applyDeleteProductField(cfg *deleteProductFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productFieldRequestOptions) applyAddProductFieldOptions(cfg *addProductFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productFieldRequestOptions) applyUpdateProductFieldOptions(cfg *updateProductFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productFieldRequestOptions) applyDeleteProductFieldOptions(cfg *deleteProductFieldOptionItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getProductFieldOptionFunc func(*getProductFieldOptions)

func (f getProductFieldOptionFunc) applyGetProductField(cfg *getProductFieldOptions) {
	f(cfg)
}

type listProductFieldsOptionFunc func(*listProductFieldsOptions)

func (f listProductFieldsOptionFunc) applyListProductFields(cfg *listProductFieldsOptions) {
	f(cfg)
}

type productFieldOptionFunc func(*fieldPayload)

func (f productFieldOptionFunc) applyCreateProductField(cfg *createProductFieldOptions) {
	f(&cfg.payload)
}

func (f productFieldOptionFunc) applyUpdateProductField(cfg *updateProductFieldOptions) {
	f(&cfg.payload)
}

func WithProductFieldRequestOptions(opts ...pipedrive.RequestOption) ProductFieldRequestOption {
	return productFieldRequestOptions{requestOptions: opts}
}

func WithProductFieldIncludeFields(fields ...FieldIncludeField) GetProductFieldOption {
	return getProductFieldOptionFunc(func(cfg *getProductFieldOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetProductFieldParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithProductFieldsIncludeFields(fields ...FieldIncludeField) ListProductFieldsOption {
	return listProductFieldsOptionFunc(func(cfg *listProductFieldsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetProductFieldsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithProductFieldsPageSize(limit int) ListProductFieldsOption {
	return listProductFieldsOptionFunc(func(cfg *listProductFieldsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithProductFieldsCursor(cursor string) ListProductFieldsOption {
	return listProductFieldsOptionFunc(func(cfg *listProductFieldsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithProductFieldName(name string) ProductFieldOption {
	return productFieldOptionFunc(func(payload *fieldPayload) {
		payload.name = &name
	})
}

func WithProductFieldType(fieldType FieldType) ProductFieldOption {
	return productFieldOptionFunc(func(payload *fieldPayload) {
		payload.fieldType = &fieldType
	})
}

func WithProductFieldDescription(description string) ProductFieldOption {
	return productFieldOptionFunc(func(payload *fieldPayload) {
		payload.description = &description
	})
}

func WithProductFieldOptions(labels ...string) ProductFieldOption {
	return productFieldOptionFunc(func(payload *fieldPayload) {
		payload.addOptions(labels...)
	})
}

func WithProductFieldUIVisibility(value map[string]interface{}) ProductFieldOption {
	return productFieldOptionFunc(func(payload *fieldPayload) {
		payload.uiVisibility = value
	})
}

func WithProductFieldImportantFields(value map[string]interface{}) ProductFieldOption {
	return productFieldOptionFunc(func(payload *fieldPayload) {
		payload.importantFields = value
	})
}

func WithProductFieldRequiredFields(value map[string]interface{}) ProductFieldOption {
	return productFieldOptionFunc(func(payload *fieldPayload) {
		payload.requiredFields = value
	})
}

func newGetProductFieldOptions(opts []GetProductFieldOption) getProductFieldOptions {
	var cfg getProductFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetProductField(&cfg)
	}
	return cfg
}

func newListProductFieldsOptions(opts []ListProductFieldsOption) listProductFieldsOptions {
	var cfg listProductFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListProductFields(&cfg)
	}
	return cfg
}

func newCreateProductFieldOptions(opts []CreateProductFieldOption) createProductFieldOptions {
	var cfg createProductFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateProductField(&cfg)
	}
	return cfg
}

func newUpdateProductFieldOptions(opts []UpdateProductFieldOption) updateProductFieldOptions {
	var cfg updateProductFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateProductField(&cfg)
	}
	return cfg
}

func newDeleteProductFieldOptions(opts []DeleteProductFieldOption) deleteProductFieldOptions {
	var cfg deleteProductFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteProductField(&cfg)
	}
	return cfg
}

func newAddProductFieldOptions(opts []AddProductFieldOptionsOption) addProductFieldOptionItemsOptions {
	var cfg addProductFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddProductFieldOptions(&cfg)
	}
	return cfg
}

func newUpdateProductFieldOptionsOptions(opts []UpdateProductFieldOptionsOption) updateProductFieldOptionItemsOptions {
	var cfg updateProductFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateProductFieldOptions(&cfg)
	}
	return cfg
}

func newDeleteProductFieldOptionsOptions(opts []DeleteProductFieldOptionsOption) deleteProductFieldOptionItemsOptions {
	var cfg deleteProductFieldOptionItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteProductFieldOptions(&cfg)
	}
	return cfg
}

func (s *ProductFieldsService) Get(ctx context.Context, fieldCode string, opts ...GetProductFieldOption) (*Field, error) {
	cfg := newGetProductFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetProductFieldWithResponse(ctx, fieldCode, &cfg.params, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing product field data in response")
	}
	return payload.Data, nil
}

func (s *ProductFieldsService) List(ctx context.Context, opts ...ListProductFieldsOption) ([]Field, *string, error) {
	cfg := newListProductFieldsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *ProductFieldsService) ListPager(opts ...ListProductFieldsOption) *pipedrive.CursorPager[Field] {
	cfg := newListProductFieldsOptions(opts)
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

func (s *ProductFieldsService) ForEach(ctx context.Context, fn func(Field) error, opts ...ListProductFieldsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *ProductFieldsService) Create(ctx context.Context, opts ...CreateProductFieldOption) (*Field, error) {
	cfg := newCreateProductFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddProductFieldWithBodyWithResponse(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing product field data in response")
	}
	return payload.Data, nil
}

func (s *ProductFieldsService) Update(ctx context.Context, fieldCode string, opts ...UpdateProductFieldOption) (*Field, error) {
	cfg := newUpdateProductFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateProductFieldWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing product field data in response")
	}
	return payload.Data, nil
}

func (s *ProductFieldsService) Delete(ctx context.Context, fieldCode string, opts ...DeleteProductFieldOption) (*Field, error) {
	cfg := newDeleteProductFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteProductFieldWithResponse(ctx, fieldCode, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing product field data in response")
	}
	return payload.Data, nil
}

func (s *ProductFieldsService) AddOptions(ctx context.Context, fieldCode string, labels []string, opts ...AddProductFieldOptionsOption) ([]FieldOption, error) {
	cfg := newAddProductFieldOptions(opts)
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

	resp, err := s.client.gen.AddProductFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *ProductFieldsService) UpdateOptions(ctx context.Context, fieldCode string, updates []FieldOptionUpdate, opts ...UpdateProductFieldOptionsOption) ([]FieldOption, error) {
	cfg := newUpdateProductFieldOptionsOptions(opts)
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

	resp, err := s.client.gen.UpdateProductFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *ProductFieldsService) DeleteOptions(ctx context.Context, fieldCode string, ids []int, opts ...DeleteProductFieldOptionsOption) ([]FieldOption, error) {
	cfg := newDeleteProductFieldOptionsOptions(opts)
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

	resp, err := s.client.gen.DeleteProductFieldOptionsWithBodyWithResponse(ctx, fieldCode, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *ProductFieldsService) list(ctx context.Context, params genv2.GetProductFieldsParams, requestOptions []pipedrive.RequestOption) ([]Field, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetProductFieldsWithResponse(ctx, &params, toRequestEditors(editors)...)
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
