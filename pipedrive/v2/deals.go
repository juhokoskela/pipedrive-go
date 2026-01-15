package v2

import (
	"context"
	"encoding/json"
	"fmt"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Deal struct {
	ID           DealID                 `json:"id"`
	Title        string                 `json:"title,omitempty"`
	Currency     string                 `json:"currency,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

type DealsService struct {
	client *Client
}

type GetDealOption interface {
	applyGetDeal(*getDealOptions)
}

type ListDealsOption interface {
	applyListDeals(*listDealsOptions)
}

type DealRequestOption interface {
	GetDealOption
	ListDealsOption
}

type getDealOptions struct {
	requestOptions []pipedrive.RequestOption
}

type listDealsOptions struct {
	params         genv2.GetDealsParams
	requestOptions []pipedrive.RequestOption
}

type dealRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o dealRequestOptions) applyGetDeal(cfg *getDealOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o dealRequestOptions) applyListDeals(cfg *listDealsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listDealsOptionFunc func(*listDealsOptions)

func (f listDealsOptionFunc) applyListDeals(cfg *listDealsOptions) {
	f(cfg)
}

func WithDealRequestOptions(opts ...pipedrive.RequestOption) DealRequestOption {
	return dealRequestOptions{requestOptions: opts}
}

func WithDealsPageSize(limit int) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithDealsCursor(cursor string) ListDealsOption {
	return listDealsOptionFunc(func(cfg *listDealsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func newGetDealOptions(opts []GetDealOption) getDealOptions {
	var cfg getDealOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetDeal(&cfg)
	}
	return cfg
}

func newListDealsOptions(opts []ListDealsOption) listDealsOptions {
	var cfg listDealsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListDeals(&cfg)
	}
	return cfg
}

func (s *DealsService) Get(ctx context.Context, id DealID, opts ...GetDealOption) (*Deal, error) {
	cfg := newGetDealOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetDealWithResponse(ctx, int(id), nil, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Deal `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) List(ctx context.Context, opts ...ListDealsOption) ([]Deal, *string, error) {
	cfg := newListDealsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *DealsService) ListPager(opts ...ListDealsOption) *pipedrive.CursorPager[Deal] {
	cfg := newListDealsOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Deal, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *DealsService) ForEach(ctx context.Context, fn func(Deal) error, opts ...ListDealsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *DealsService) list(ctx context.Context, params genv2.GetDealsParams, requestOptions []pipedrive.RequestOption) ([]Deal, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetDealsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Deal `json:"data"`
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
