package v2

import (
	"context"
	"encoding/json"
	"fmt"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ActivityFieldsService struct {
	client *Client
}

type GetActivityFieldOption interface {
	applyGetActivityField(*getActivityFieldOptions)
}

type ListActivityFieldsOption interface {
	applyListActivityFields(*listActivityFieldsOptions)
}

type ActivityFieldRequestOption interface {
	GetActivityFieldOption
	ListActivityFieldsOption
}

type getActivityFieldOptions struct {
	params         genv2.GetActivityFieldParams
	requestOptions []pipedrive.RequestOption
}

type listActivityFieldsOptions struct {
	params         genv2.GetActivityFieldsParams
	requestOptions []pipedrive.RequestOption
}

type activityFieldRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o activityFieldRequestOptions) applyGetActivityField(cfg *getActivityFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o activityFieldRequestOptions) applyListActivityFields(cfg *listActivityFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getActivityFieldOptionFunc func(*getActivityFieldOptions)

func (f getActivityFieldOptionFunc) applyGetActivityField(cfg *getActivityFieldOptions) {
	f(cfg)
}

type listActivityFieldsOptionFunc func(*listActivityFieldsOptions)

func (f listActivityFieldsOptionFunc) applyListActivityFields(cfg *listActivityFieldsOptions) {
	f(cfg)
}

func WithActivityFieldRequestOptions(opts ...pipedrive.RequestOption) ActivityFieldRequestOption {
	return activityFieldRequestOptions{requestOptions: opts}
}

func WithActivityFieldIncludeFields(fields ...FieldIncludeField) GetActivityFieldOption {
	return getActivityFieldOptionFunc(func(cfg *getActivityFieldOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetActivityFieldParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithActivityFieldsIncludeFields(fields ...FieldIncludeField) ListActivityFieldsOption {
	return listActivityFieldsOptionFunc(func(cfg *listActivityFieldsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetActivityFieldsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithActivityFieldsPageSize(limit int) ListActivityFieldsOption {
	return listActivityFieldsOptionFunc(func(cfg *listActivityFieldsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithActivityFieldsCursor(cursor string) ListActivityFieldsOption {
	return listActivityFieldsOptionFunc(func(cfg *listActivityFieldsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func newGetActivityFieldOptions(opts []GetActivityFieldOption) getActivityFieldOptions {
	var cfg getActivityFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetActivityField(&cfg)
	}
	return cfg
}

func newListActivityFieldsOptions(opts []ListActivityFieldsOption) listActivityFieldsOptions {
	var cfg listActivityFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListActivityFields(&cfg)
	}
	return cfg
}

func (s *ActivityFieldsService) Get(ctx context.Context, fieldCode string, opts ...GetActivityFieldOption) (*Field, error) {
	cfg := newGetActivityFieldOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetActivityFieldWithResponse(ctx, fieldCode, &cfg.params, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing activity field data in response")
	}
	return payload.Data, nil
}

func (s *ActivityFieldsService) List(ctx context.Context, opts ...ListActivityFieldsOption) ([]Field, *string, error) {
	cfg := newListActivityFieldsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *ActivityFieldsService) ListPager(opts ...ListActivityFieldsOption) *pipedrive.CursorPager[Field] {
	cfg := newListActivityFieldsOptions(opts)
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

func (s *ActivityFieldsService) ForEach(ctx context.Context, fn func(Field) error, opts ...ListActivityFieldsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *ActivityFieldsService) list(ctx context.Context, params genv2.GetActivityFieldsParams, requestOptions []pipedrive.RequestOption) ([]Field, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetActivityFieldsWithResponse(ctx, &params, toRequestEditors(editors)...)
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
