package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Product struct {
	ID         ProductID `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Code       string    `json:"code,omitempty"`
	OwnerID    *UserID   `json:"owner_id,omitempty"`
	Active     bool      `json:"active_flag,omitempty"`
	AddTime    *DateTime `json:"add_time,omitempty"`
	UpdateTime *DateTime `json:"update_time,omitempty"`
}

type ProductsService struct {
	client *Client
}

type ProductsOption interface {
	applyProducts(*productsOptions)
}

type productsOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type productsOptionFunc func(*productsOptions)

func (f productsOptionFunc) applyProducts(cfg *productsOptions) {
	f(cfg)
}

func WithProductsQuery(values url.Values) ProductsOption {
	return productsOptionFunc(func(cfg *productsOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithProductsRequestOptions(opts ...pipedrive.RequestOption) ProductsOption {
	return productsOptionFunc(func(cfg *productsOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newProductsOptions(opts []ProductsOption) productsOptions {
	var cfg productsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyProducts(&cfg)
	}
	return cfg
}

func (s *ProductsService) ListDeals(ctx context.Context, id ProductID, opts ...ProductsOption) ([]Deal, *Pagination, error) {
	cfg := newProductsOptions(opts)
	path := fmt.Sprintf("/products/%d/deals", id)

	var payload struct {
		Data           []Deal `json:"data"`
		AdditionalData *struct {
			Pagination *Pagination `json:"pagination"`
		} `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	var page *Pagination
	if payload.AdditionalData != nil {
		page = payload.AdditionalData.Pagination
	}
	return payload.Data, page, nil
}

func (s *ProductsService) ListFiles(ctx context.Context, id ProductID, opts ...ProductsOption) ([]File, *Pagination, error) {
	cfg := newProductsOptions(opts)
	path := fmt.Sprintf("/products/%d/files", id)

	var payload struct {
		Data           []File `json:"data"`
		AdditionalData *struct {
			Pagination *Pagination `json:"pagination"`
		} `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	var page *Pagination
	if payload.AdditionalData != nil {
		page = payload.AdditionalData.Pagination
	}
	return payload.Data, page, nil
}

func (s *ProductsService) ListUsers(ctx context.Context, id ProductID, opts ...ProductsOption) ([]User, error) {
	cfg := newProductsOptions(opts)
	path := fmt.Sprintf("/products/%d/permittedUsers", id)

	var payload struct {
		Data []User `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}
