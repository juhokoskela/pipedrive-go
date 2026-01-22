package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"time"

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

type ProductVariation struct {
	ID        ProductVariationID `json:"id"`
	Name      string             `json:"name,omitempty"`
	ProductID *ProductID         `json:"product_id,omitempty"`
	Prices    []ProductPrice     `json:"prices,omitempty"`
}

type ProductImage struct {
	ID        ProductImageID `json:"id"`
	ProductID *ProductID     `json:"product_id,omitempty"`
	CompanyID *string        `json:"company_id,omitempty"`
	PublicURL *string        `json:"public_url,omitempty"`
	AddTime   *time.Time     `json:"add_time,omitempty"`
	MimeType  *string        `json:"mime_type,omitempty"`
	Name      *string        `json:"name,omitempty"`
}

type ProductDeleteResult struct {
	ID ProductID `json:"id"`
}

type ProductVariationDeleteResult struct {
	ID ProductVariationID `json:"id"`
}

type ProductImageDeleteResult struct {
	ID ProductImageID `json:"id"`
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

type ListProductVariationsOption interface {
	applyListProductVariations(*listProductVariationsOptions)
}

type CreateProductVariationOption interface {
	applyCreateProductVariation(*createProductVariationOptions)
}

type UpdateProductVariationOption interface {
	applyUpdateProductVariation(*updateProductVariationOptions)
}

type DeleteProductVariationOption interface {
	applyDeleteProductVariation(*deleteProductVariationOptions)
}

type GetProductImageOption interface {
	applyGetProductImage(*getProductImageOptions)
}

type UploadProductImageOption interface {
	applyUploadProductImage(*uploadProductImageOptions)
}

type UpdateProductImageOption interface {
	applyUpdateProductImage(*updateProductImageOptions)
}

type DeleteProductImageOption interface {
	applyDeleteProductImage(*deleteProductImageOptions)
}

type GetProductFollowersOption interface {
	applyGetProductFollowers(*getProductFollowersOptions)
}

type AddProductFollowerOption interface {
	applyAddProductFollower(*addProductFollowerOptions)
}

type DeleteProductFollowerOption interface {
	applyDeleteProductFollower(*deleteProductFollowerOptions)
}

type GetProductFollowersChangelogOption interface {
	applyGetProductFollowersChangelog(*getProductFollowersChangelogOptions)
}

type ProductRequestOption interface {
	GetProductOption
	ListProductsOption
	CreateProductOption
	UpdateProductOption
	DeleteProductOption
	SearchProductsOption
	DuplicateProductOption
	ListProductVariationsOption
	CreateProductVariationOption
	UpdateProductVariationOption
	DeleteProductVariationOption
	GetProductImageOption
	UploadProductImageOption
	UpdateProductImageOption
	DeleteProductImageOption
	GetProductFollowersOption
	AddProductFollowerOption
	DeleteProductFollowerOption
	GetProductFollowersChangelogOption
}

type ProductOption interface {
	CreateProductOption
	UpdateProductOption
}

type ProductVariationOption interface {
	CreateProductVariationOption
	UpdateProductVariationOption
}

type ProductImageOption interface {
	UploadProductImageOption
	UpdateProductImageOption
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

type listProductVariationsOptions struct {
	params         genv2.GetProductVariationsParams
	requestOptions []pipedrive.RequestOption
}

type createProductVariationOptions struct {
	payload        productVariationPayload
	requestOptions []pipedrive.RequestOption
}

type updateProductVariationOptions struct {
	payload        productVariationPayload
	requestOptions []pipedrive.RequestOption
}

type deleteProductVariationOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getProductImageOptions struct {
	requestOptions []pipedrive.RequestOption
}

type uploadProductImageOptions struct {
	payload        productImagePayload
	requestOptions []pipedrive.RequestOption
}

type updateProductImageOptions struct {
	payload        productImagePayload
	requestOptions []pipedrive.RequestOption
}

type deleteProductImageOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getProductFollowersOptions struct {
	params         genv2.GetProductFollowersParams
	requestOptions []pipedrive.RequestOption
}

type addProductFollowerOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deleteProductFollowerOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getProductFollowersChangelogOptions struct {
	params         genv2.GetProductFollowersChangelogParams
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

type productVariationPayload struct {
	name   *string
	prices []ProductPrice
}

type productImagePayload struct {
	fileName string
	reader   io.Reader
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

func (o productRequestOptions) applyListProductVariations(cfg *listProductVariationsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyCreateProductVariation(cfg *createProductVariationOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyUpdateProductVariation(cfg *updateProductVariationOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyDeleteProductVariation(cfg *deleteProductVariationOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyGetProductImage(cfg *getProductImageOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyUploadProductImage(cfg *uploadProductImageOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyUpdateProductImage(cfg *updateProductImageOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyDeleteProductImage(cfg *deleteProductImageOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyGetProductFollowers(cfg *getProductFollowersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyAddProductFollower(cfg *addProductFollowerOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyDeleteProductFollower(cfg *deleteProductFollowerOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o productRequestOptions) applyGetProductFollowersChangelog(cfg *getProductFollowersChangelogOptions) {
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

type listProductVariationsOptionFunc func(*listProductVariationsOptions)

func (f listProductVariationsOptionFunc) applyListProductVariations(cfg *listProductVariationsOptions) {
	f(cfg)
}

type getProductFollowersOptionFunc func(*getProductFollowersOptions)

func (f getProductFollowersOptionFunc) applyGetProductFollowers(cfg *getProductFollowersOptions) {
	f(cfg)
}

type getProductFollowersChangelogOptionFunc func(*getProductFollowersChangelogOptions)

func (f getProductFollowersChangelogOptionFunc) applyGetProductFollowersChangelog(cfg *getProductFollowersChangelogOptions) {
	f(cfg)
}

type productFieldOption func(*productPayload)

func (f productFieldOption) applyCreateProduct(cfg *createProductOptions) {
	f(&cfg.payload)
}

func (f productFieldOption) applyUpdateProduct(cfg *updateProductOptions) {
	f(&cfg.payload)
}

type productVariationFieldOption func(*productVariationPayload)

func (f productVariationFieldOption) applyCreateProductVariation(cfg *createProductVariationOptions) {
	f(&cfg.payload)
}

func (f productVariationFieldOption) applyUpdateProductVariation(cfg *updateProductVariationOptions) {
	f(&cfg.payload)
}

type productImageFieldOption func(*productImagePayload)

func (f productImageFieldOption) applyUploadProductImage(cfg *uploadProductImageOptions) {
	f(&cfg.payload)
}

func (f productImageFieldOption) applyUpdateProductImage(cfg *updateProductImageOptions) {
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

func WithProductVariationsPageSize(limit int) ListProductVariationsOption {
	return listProductVariationsOptionFunc(func(cfg *listProductVariationsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithProductVariationsCursor(cursor string) ListProductVariationsOption {
	return listProductVariationsOptionFunc(func(cfg *listProductVariationsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithProductVariationName(name string) ProductVariationOption {
	return productVariationFieldOption(func(payload *productVariationPayload) {
		payload.name = &name
	})
}

func WithProductVariationPrices(prices ...ProductPrice) ProductVariationOption {
	return productVariationFieldOption(func(payload *productVariationPayload) {
		if len(prices) == 0 {
			return
		}
		payload.prices = append(payload.prices, prices...)
	})
}

func WithProductImageFile(name string, reader io.Reader) ProductImageOption {
	return productImageFieldOption(func(payload *productImagePayload) {
		payload.fileName = name
		payload.reader = reader
	})
}

func WithProductFollowersPageSize(limit int) GetProductFollowersOption {
	return getProductFollowersOptionFunc(func(cfg *getProductFollowersOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithProductFollowersCursor(cursor string) GetProductFollowersOption {
	return getProductFollowersOptionFunc(func(cfg *getProductFollowersOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithProductFollowersChangelogPageSize(limit int) GetProductFollowersChangelogOption {
	return getProductFollowersChangelogOptionFunc(func(cfg *getProductFollowersChangelogOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithProductFollowersChangelogCursor(cursor string) GetProductFollowersChangelogOption {
	return getProductFollowersChangelogOptionFunc(func(cfg *getProductFollowersChangelogOptions) {
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

func newGetProductImageOptions(opts []GetProductImageOption) getProductImageOptions {
	var cfg getProductImageOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetProductImage(&cfg)
	}
	return cfg
}

func newUploadProductImageOptions(opts []UploadProductImageOption) uploadProductImageOptions {
	var cfg uploadProductImageOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUploadProductImage(&cfg)
	}
	return cfg
}

func newUpdateProductImageOptions(opts []UpdateProductImageOption) updateProductImageOptions {
	var cfg updateProductImageOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateProductImage(&cfg)
	}
	return cfg
}

func newDeleteProductImageOptions(opts []DeleteProductImageOption) deleteProductImageOptions {
	var cfg deleteProductImageOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteProductImage(&cfg)
	}
	return cfg
}

func newListProductVariationsOptions(opts []ListProductVariationsOption) listProductVariationsOptions {
	var cfg listProductVariationsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListProductVariations(&cfg)
	}
	return cfg
}

func newCreateProductVariationOptions(opts []CreateProductVariationOption) createProductVariationOptions {
	var cfg createProductVariationOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateProductVariation(&cfg)
	}
	return cfg
}

func newUpdateProductVariationOptions(opts []UpdateProductVariationOption) updateProductVariationOptions {
	var cfg updateProductVariationOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateProductVariation(&cfg)
	}
	return cfg
}

func newDeleteProductVariationOptions(opts []DeleteProductVariationOption) deleteProductVariationOptions {
	var cfg deleteProductVariationOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteProductVariation(&cfg)
	}
	return cfg
}

func newGetProductFollowersOptions(opts []GetProductFollowersOption) getProductFollowersOptions {
	var cfg getProductFollowersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetProductFollowers(&cfg)
	}
	return cfg
}

func newAddProductFollowerOptions(opts []AddProductFollowerOption) addProductFollowerOptions {
	var cfg addProductFollowerOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddProductFollower(&cfg)
	}
	return cfg
}

func newDeleteProductFollowerOptions(opts []DeleteProductFollowerOption) deleteProductFollowerOptions {
	var cfg deleteProductFollowerOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteProductFollower(&cfg)
	}
	return cfg
}

func newGetProductFollowersChangelogOptions(opts []GetProductFollowersChangelogOption) getProductFollowersChangelogOptions {
	var cfg getProductFollowersChangelogOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetProductFollowersChangelog(&cfg)
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

func (s *ProductsService) ListVariations(ctx context.Context, id ProductID, opts ...ListProductVariationsOption) ([]ProductVariation, *string, error) {
	cfg := newListProductVariationsOptions(opts)
	return s.listVariations(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *ProductsService) ListVariationsPager(id ProductID, opts ...ListProductVariationsOption) *pipedrive.CursorPager[ProductVariation] {
	cfg := newListProductVariationsOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]ProductVariation, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listVariations(ctx, id, params, cfg.requestOptions)
	})
}

func (s *ProductsService) ForEachVariations(ctx context.Context, id ProductID, fn func(ProductVariation) error, opts ...ListProductVariationsOption) error {
	return s.ListVariationsPager(id, opts...).ForEach(ctx, fn)
}

func (s *ProductsService) CreateVariation(ctx context.Context, id ProductID, opts ...CreateProductVariationOption) (*ProductVariation, error) {
	cfg := newCreateProductVariationOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddProductVariationWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *ProductVariation `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product variation data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) UpdateVariation(ctx context.Context, id ProductID, variationID ProductVariationID, opts ...UpdateProductVariationOption) (*ProductVariation, error) {
	cfg := newUpdateProductVariationOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateProductVariationWithBodyWithResponse(ctx, int(id), int(variationID), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *ProductVariation `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product variation data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) DeleteVariation(ctx context.Context, id ProductID, variationID ProductVariationID, opts ...DeleteProductVariationOption) (*ProductVariationDeleteResult, error) {
	cfg := newDeleteProductVariationOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteProductVariationWithResponse(ctx, int(id), int(variationID), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *ProductVariationDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product variation delete data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) GetImage(ctx context.Context, id ProductID, opts ...GetProductImageOption) (*ProductImage, error) {
	cfg := newGetProductImageOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetProductImageWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *ProductImage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product image data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) UploadImage(ctx context.Context, id ProductID, opts ...UploadProductImageOption) (*ProductImage, error) {
	cfg := newUploadProductImageOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	contentType, body, err := cfg.payload.toMultipart()
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.UploadProductImageWithBodyWithResponse(ctx, int(id), contentType, body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *ProductImage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product image data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) UpdateImage(ctx context.Context, id ProductID, opts ...UpdateProductImageOption) (*ProductImage, error) {
	cfg := newUpdateProductImageOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	contentType, body, err := cfg.payload.toMultipart()
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.UpdateProductImageWithBodyWithResponse(ctx, int(id), contentType, body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *ProductImage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product image data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) DeleteImage(ctx context.Context, id ProductID, opts ...DeleteProductImageOption) (*ProductImageDeleteResult, error) {
	cfg := newDeleteProductImageOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteProductImageWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *ProductImageDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing product image delete data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) ListFollowers(ctx context.Context, id ProductID, opts ...GetProductFollowersOption) ([]Follower, *string, error) {
	cfg := newGetProductFollowersOptions(opts)
	return s.listFollowers(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *ProductsService) ListFollowersPager(id ProductID, opts ...GetProductFollowersOption) *pipedrive.CursorPager[Follower] {
	cfg := newGetProductFollowersOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Follower, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listFollowers(ctx, id, params, cfg.requestOptions)
	})
}

func (s *ProductsService) ForEachFollowers(ctx context.Context, id ProductID, fn func(Follower) error, opts ...GetProductFollowersOption) error {
	return s.ListFollowersPager(id, opts...).ForEach(ctx, fn)
}

func (s *ProductsService) AddFollower(ctx context.Context, id ProductID, userID UserID, opts ...AddProductFollowerOption) (*Follower, error) {
	cfg := newAddProductFollowerOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(map[string]interface{}{
		"user_id": int(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddProductFollowerWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Follower `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing follower data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) DeleteFollower(ctx context.Context, id ProductID, followerID UserID, opts ...DeleteProductFollowerOption) (*FollowerDeleteResult, error) {
	cfg := newDeleteProductFollowerOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteProductFollowerWithResponse(ctx, int(id), int(followerID), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *FollowerDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing delete follower data in response")
	}
	return payload.Data, nil
}

func (s *ProductsService) FollowersChangelog(ctx context.Context, id ProductID, opts ...GetProductFollowersChangelogOption) ([]FollowerChangelog, *string, error) {
	cfg := newGetProductFollowersChangelogOptions(opts)
	return s.followersChangelog(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *ProductsService) FollowersChangelogPager(id ProductID, opts ...GetProductFollowersChangelogOption) *pipedrive.CursorPager[FollowerChangelog] {
	cfg := newGetProductFollowersChangelogOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]FollowerChangelog, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.followersChangelog(ctx, id, params, cfg.requestOptions)
	})
}

func (s *ProductsService) ForEachFollowersChangelog(ctx context.Context, id ProductID, fn func(FollowerChangelog) error, opts ...GetProductFollowersChangelogOption) error {
	return s.FollowersChangelogPager(id, opts...).ForEach(ctx, fn)
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

func (s *ProductsService) listVariations(ctx context.Context, id ProductID, params genv2.GetProductVariationsParams, requestOptions []pipedrive.RequestOption) ([]ProductVariation, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetProductVariationsWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []ProductVariation `json:"data"`
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

func (s *ProductsService) listFollowers(ctx context.Context, id ProductID, params genv2.GetProductFollowersParams, requestOptions []pipedrive.RequestOption) ([]Follower, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetProductFollowersWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Follower `json:"data"`
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

func (s *ProductsService) followersChangelog(ctx context.Context, id ProductID, params genv2.GetProductFollowersChangelogParams, requestOptions []pipedrive.RequestOption) ([]FollowerChangelog, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetProductFollowersChangelogWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []FollowerChangelog `json:"data"`
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

func (p productVariationPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if len(p.prices) > 0 {
		body["prices"] = p.prices
	}
	return body
}

func (p productImagePayload) toMultipart() (string, *bytes.Buffer, error) {
	if p.reader == nil || p.fileName == "" {
		return "", nil, fmt.Errorf("product image file is required")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("data", p.fileName)
	if err != nil {
		return "", nil, fmt.Errorf("create multipart file: %w", err)
	}
	if _, err := io.Copy(part, p.reader); err != nil {
		return "", nil, fmt.Errorf("write multipart file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", nil, fmt.Errorf("close multipart writer: %w", err)
	}
	return writer.FormDataContentType(), &buf, nil
}
