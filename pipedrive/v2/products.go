package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type BillingFrequency string

const (
	BillingFrequencyOneTime      BillingFrequency = "one-time"
	BillingFrequencyAnnually     BillingFrequency = "annually"
	BillingFrequencySemiAnnually BillingFrequency = "semi-annually"
	BillingFrequencyQuarterly    BillingFrequency = "quarterly"
	BillingFrequencyMonthly      BillingFrequency = "monthly"
	BillingFrequencyWeekly       BillingFrequency = "weekly"
)

type ProductSortField string

const (
	ProductSortByAddTime    ProductSortField = "add_time"
	ProductSortByID         ProductSortField = "id"
	ProductSortByName       ProductSortField = "name"
	ProductSortByUpdateTime ProductSortField = "update_time"
)

type ProductSearchField string

const (
	ProductSearchFieldCode         ProductSearchField = "code"
	ProductSearchFieldCustomFields ProductSearchField = "custom_fields"
	ProductSearchFieldName         ProductSearchField = "name"
)

type ProductPrice struct {
	ProductID          *ProductID          `json:"product_id,omitempty"`
	ProductVariationID *ProductVariationID `json:"product_variation_id,omitempty"`
	Currency           string              `json:"currency,omitempty"`
	Price              float64             `json:"price,omitempty"`
	Cost               *float64            `json:"cost,omitempty"`
	DirectCost         *float64            `json:"direct_cost,omitempty"`
	Notes              string              `json:"notes,omitempty"`
}

type Product struct {
	ID                     ProductID              `json:"id"`
	Name                   string                 `json:"name,omitempty"`
	Code                   string                 `json:"code,omitempty"`
	Unit                   string                 `json:"unit,omitempty"`
	Tax                    float64                `json:"tax,omitempty"`
	Category               float64                `json:"category,omitempty"`
	OwnerID                *UserID                `json:"owner_id,omitempty"`
	IsDeleted              bool                   `json:"is_deleted,omitempty"`
	IsLinkable             bool                   `json:"is_linkable,omitempty"`
	VisibleTo              int                    `json:"visible_to,omitempty"`
	CustomFields           map[string]interface{} `json:"custom_fields,omitempty"`
	BillingFrequency       BillingFrequency       `json:"billing_frequency,omitempty"`
	BillingFrequencyCycles *int                   `json:"billing_frequency_cycles,omitempty"`
	Prices                 []ProductPrice         `json:"prices,omitempty"`
}

type ProductSearchItem struct {
	ResultScore float64                `json:"result_score,omitempty"`
	Item        map[string]interface{} `json:"item,omitempty"`
}

type ProductSearchResults struct {
	Items []ProductSearchItem `json:"items,omitempty"`
}

type ProductDeleteResult struct {
	ID ProductID `json:"id"`
}

type ProductsService struct {
	client *Client
}

type GetProductOption interface {
	applyGetProduct(*getProductOptions)
}

type ListProductsOption interface {
	applyListProducts(*listProductsOptions)
}

type CreateProductOption interface {
	applyCreateProduct(*createProductOptions)
}

type UpdateProductOption interface {
	applyUpdateProduct(*updateProductOptions)
}

type DeleteProductOption interface {
	applyDeleteProduct(*deleteProductOptions)
}

type SearchProductsOption interface {
	applySearchProducts(*searchProductsOptions)
}

type DuplicateProductOption interface {
	applyDuplicateProduct(*duplicateProductOptions)
}

type ProductRequestOption interface {
	GetProductOption
	ListProductsOption
	CreateProductOption
	UpdateProductOption
	DeleteProductOption
	SearchProductsOption
	DuplicateProductOption
}

type ProductOption interface {
	CreateProductOption
	UpdateProductOption
}

type getProductOptions struct {
	requestOptions []pipedrive.RequestOption
}

type listProductsOptions struct {
	params         genv2.GetProductsParams
	requestOptions []pipedrive.RequestOption
}

type createProductOptions struct {
	payload        productPayload
	requestOptions []pipedrive.RequestOption
}

type updateProductOptions struct {
	payload        productPayload
	requestOptions []pipedrive.RequestOption
}

type deleteProductOptions struct {
	requestOptions []pipedrive.RequestOption
}

type searchProductsOptions struct {
	params         genv2.SearchProductsParams
	requestOptions []pipedrive.RequestOption
}

type duplicateProductOptions struct {
	requestOptions []pipedrive.RequestOption
}

type productPayload struct {
	name                   *string
	code                   *string
	description            *string
	unit                   *string
	tax                    *float64
	category               *float64
	ownerID                *UserID
	isLinkable             *bool
	visibleTo              *int
	prices                 []ProductPrice
	billingFrequency       *BillingFrequency
	billingFrequencyCycles *int
}

type productRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o productRequestOptions) applyGetProduct(cfg *getProductOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyListProducts(cfg *listProductsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyCreateProduct(cfg *createProductOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyUpdateProduct(cfg *updateProductOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyDeleteProduct(cfg *deleteProductOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applySearchProducts(cfg *searchProductsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyDuplicateProduct(cfg *duplicateProductOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listProductsOptionFunc func(*listProductsOptions)

func (f listProductsOptionFunc) applyListProducts(cfg *listProductsOptions) {
	f(cfg)
}

type searchProductsOptionFunc func(*searchProductsOptions)

func (f searchProductsOptionFunc) applySearchProducts(cfg *searchProductsOptions) {
	f(cfg)
}

type productFieldOption func(*productPayload)

func (f productFieldOption) applyCreateProduct(cfg *createProductOptions) {
	f(&cfg.payload)
}

func (f productFieldOption) applyUpdateProduct(cfg *updateProductOptions) {
	f(&cfg.payload)
}

func WithProductRequestOptions(opts ...pipedrive.RequestOption) ProductRequestOption {
	return productRequestOptions{requestOptions: opts}
}

func WithProductsOwnerID(id UserID) ListProductsOption {
	return listProductsOptionFunc(func(cfg *listProductsOptions) {
		value := int(id)
		cfg.params.OwnerId = &value
	})
}

func WithProductsIDs(ids ...ProductID) ListProductsOption {
	return listProductsOptionFunc(func(cfg *listProductsOptions) {
		csv := joinIDs(ids)
		if csv == "" {
			return
		}
		cfg.params.Ids = &csv
	})
}

func WithProductsFilterID(id int) ListProductsOption {
	return listProductsOptionFunc(func(cfg *listProductsOptions) {
		cfg.params.FilterId = &id
	})
}

func WithProductsPageSize(limit int) ListProductsOption {
	return listProductsOptionFunc(func(cfg *listProductsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithProductsCursor(cursor string) ListProductsOption {
	return listProductsOptionFunc(func(cfg *listProductsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithProductsSortBy(field ProductSortField) ListProductsOption {
	return listProductsOptionFunc(func(cfg *listProductsOptions) {
		value := genv2.GetProductsParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithProductsSortDirection(direction SortDirection) ListProductsOption {
	return listProductsOptionFunc(func(cfg *listProductsOptions) {
		value := genv2.GetProductsParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithProductsCustomFields(fields ...string) ListProductsOption {
	return listProductsOptionFunc(func(cfg *listProductsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		cfg.params.CustomFields = &csv
	})
}

func WithProductName(name string) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.name = &name
	})
}

func WithProductCode(code string) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.code = &code
	})
}

func WithProductDescription(description string) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.description = &description
	})
}

func WithProductUnit(unit string) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.unit = &unit
	})
}

func WithProductTax(tax float64) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.tax = &tax
	})
}

func WithProductCategory(category float64) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.category = &category
	})
}

func WithProductOwnerID(id UserID) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.ownerID = &id
	})
}

func WithProductLinkable(linkable bool) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.isLinkable = &linkable
	})
}

func WithProductVisibleTo(visibleTo int) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.visibleTo = &visibleTo
	})
}

func WithProductPrices(prices ...ProductPrice) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		if len(prices) == 0 {
			return
		}
		payload.prices = append(payload.prices, prices...)
	})
}

func WithProductBillingFrequency(frequency BillingFrequency) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.billingFrequency = &frequency
	})
}

func WithProductBillingFrequencyCycles(cycles int) ProductOption {
	return productFieldOption(func(payload *productPayload) {
		payload.billingFrequencyCycles = &cycles
	})
}

func WithProductSearchFields(fields ...ProductSearchField) SearchProductsOption {
	return searchProductsOptionFunc(func(cfg *searchProductsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.SearchProductsParamsFields(csv)
		cfg.params.Fields = &value
	})
}

func WithProductSearchExactMatch(enabled bool) SearchProductsOption {
	return searchProductsOptionFunc(func(cfg *searchProductsOptions) {
		cfg.params.ExactMatch = &enabled
	})
}

func WithProductSearchPageSize(limit int) SearchProductsOption {
	return searchProductsOptionFunc(func(cfg *searchProductsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithProductSearchCursor(cursor string) SearchProductsOption {
	return searchProductsOptionFunc(func(cfg *searchProductsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func newGetProductOptions(opts []GetProductOption) getProductOptions {
	var cfg getProductOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetProduct(&cfg)
	}
	return cfg
}

func newListProductsOptions(opts []ListProductsOption) listProductsOptions {
	var cfg listProductsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListProducts(&cfg)
	}
	return cfg
}

func newCreateProductOptions(opts []CreateProductOption) createProductOptions {
	var cfg createProductOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateProduct(&cfg)
	}
	return cfg
}

func newUpdateProductOptions(opts []UpdateProductOption) updateProductOptions {
	var cfg updateProductOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateProduct(&cfg)
	}
	return cfg
}

func newDeleteProductOptions(opts []DeleteProductOption) deleteProductOptions {
	var cfg deleteProductOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteProduct(&cfg)
	}
	return cfg
}

func newSearchProductsOptions(opts []SearchProductsOption) searchProductsOptions {
	var cfg searchProductsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applySearchProducts(&cfg)
	}
	return cfg
}

func newDuplicateProductOptions(opts []DuplicateProductOption) duplicateProductOptions {
	var cfg duplicateProductOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDuplicateProduct(&cfg)
	}
	return cfg
}

func (s *ProductsService) Get(ctx context.Context, id ProductID, opts ...GetProductOption) (*Product, error) {
	cfg := newGetProductOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetProductWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Product `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) List(ctx context.Context, opts ...ListProductsOption) ([]Product, *string, error) {
	cfg := newListProductsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *ProductsService) ListPager(opts ...ListProductsOption) *pipedrive.CursorPager[Product] {
	cfg := newListProductsOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Product, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *ProductsService) ForEach(ctx context.Context, fn func(Product) error, opts ...ListProductsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *ProductsService) Create(ctx context.Context, opts ...CreateProductOption) (*Product, error) {
	cfg := newCreateProductOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddProductWithBodyWithResponse(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Product `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) Update(ctx context.Context, id ProductID, opts ...UpdateProductOption) (*Product, error) {
	cfg := newUpdateProductOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateProductWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Product `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) Delete(ctx context.Context, id ProductID, opts ...DeleteProductOption) (*ProductDeleteResult, error) {
	cfg := newDeleteProductOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteProductWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *ProductDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product delete data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) Search(ctx context.Context, term string, opts ...SearchProductsOption) (*ProductSearchResults, *string, error) {
	cfg := newSearchProductsOptions(opts)
	cfg.params.Term = term
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.SearchProductsWithResponse(ctx, &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           *ProductSearchResults `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, nil, fmt.Errorf("missing product search data in response")
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *ProductsService) Duplicate(ctx context.Context, id ProductID, opts ...DuplicateProductOption) (*Product, error) {
	cfg := newDuplicateProductOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DuplicateProductWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Product `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) list(ctx context.Context, params genv2.GetProductsParams, requestOptions []pipedrive.RequestOption) ([]Product, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetProductsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Product `json:"data"`
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

func (p productPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if p.code != nil {
		body["code"] = *p.code
	}
	if p.description != nil {
		body["description"] = *p.description
	}
	if p.unit != nil {
		body["unit"] = *p.unit
	}
	if p.tax != nil {
		body["tax"] = *p.tax
	}
	if p.category != nil {
		body["category"] = *p.category
	}
	if p.ownerID != nil {
		body["owner_id"] = int(*p.ownerID)
	}
	if p.isLinkable != nil {
		body["is_linkable"] = *p.isLinkable
	}
	if p.visibleTo != nil {
		body["visible_to"] = *p.visibleTo
	}
	if len(p.prices) > 0 {
		body["prices"] = p.prices
	}
	if p.billingFrequency != nil {
		body["billing_frequency"] = *p.billingFrequency
	}
	if p.billingFrequencyCycles != nil {
		body["billing_frequency_cycles"] = *p.billingFrequencyCycles
	}
	return body
}
