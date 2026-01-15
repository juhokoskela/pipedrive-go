package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type PipelineID int64

type PipelineSortField string

const (
	PipelineSortByID         PipelineSortField = "id"
	PipelineSortByUpdateTime PipelineSortField = "update_time"
	PipelineSortByAddTime    PipelineSortField = "add_time"
)

type Pipeline struct {
	ID                     PipelineID `json:"id"`
	Name                   string     `json:"name,omitempty"`
	Order                  int        `json:"order_nr,omitempty"`
	IsDeleted              bool       `json:"is_deleted,omitempty"`
	DealProbabilityEnabled bool       `json:"is_deal_probability_enabled,omitempty"`
	AddTime                *time.Time `json:"add_time,omitempty"`
	UpdateTime             *time.Time `json:"update_time,omitempty"`
}

type PipelinesService struct {
	client *Client
}

type GetPipelineOption interface {
	applyGetPipeline(*getPipelineOptions)
}

type ListPipelinesOption interface {
	applyListPipelines(*listPipelinesOptions)
}

type CreatePipelineOption interface {
	applyCreatePipeline(*createPipelineOptions)
}

type UpdatePipelineOption interface {
	applyUpdatePipeline(*updatePipelineOptions)
}

type DeletePipelineOption interface {
	applyDeletePipeline(*deletePipelineOptions)
}

type PipelineRequestOption interface {
	GetPipelineOption
	ListPipelinesOption
	CreatePipelineOption
	UpdatePipelineOption
	DeletePipelineOption
}

type PipelineOption interface {
	CreatePipelineOption
	UpdatePipelineOption
}

type getPipelineOptions struct {
	requestOptions []pipedrive.RequestOption
}

type listPipelinesOptions struct {
	params         genv2.GetPipelinesParams
	requestOptions []pipedrive.RequestOption
}

type createPipelineOptions struct {
	payload        pipelinePayload
	requestOptions []pipedrive.RequestOption
}

type updatePipelineOptions struct {
	payload        pipelinePayload
	requestOptions []pipedrive.RequestOption
}

type deletePipelineOptions struct {
	requestOptions []pipedrive.RequestOption
}

type pipelinePayload struct {
	name                   *string
	dealProbabilityEnabled *bool
}

type pipelineRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o pipelineRequestOptions) applyGetPipeline(cfg *getPipelineOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o pipelineRequestOptions) applyListPipelines(cfg *listPipelinesOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o pipelineRequestOptions) applyCreatePipeline(cfg *createPipelineOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o pipelineRequestOptions) applyUpdatePipeline(cfg *updatePipelineOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o pipelineRequestOptions) applyDeletePipeline(cfg *deletePipelineOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listPipelinesOptionFunc func(*listPipelinesOptions)

func (f listPipelinesOptionFunc) applyListPipelines(cfg *listPipelinesOptions) {
	f(cfg)
}

type createPipelineOptionFunc func(*createPipelineOptions)

func (f createPipelineOptionFunc) applyCreatePipeline(cfg *createPipelineOptions) {
	f(cfg)
}

type pipelineFieldOption func(*pipelinePayload)

func (f pipelineFieldOption) applyCreatePipeline(cfg *createPipelineOptions) {
	f(&cfg.payload)
}

func (f pipelineFieldOption) applyUpdatePipeline(cfg *updatePipelineOptions) {
	f(&cfg.payload)
}

func WithPipelineRequestOptions(opts ...pipedrive.RequestOption) PipelineRequestOption {
	return pipelineRequestOptions{requestOptions: opts}
}

func WithPipelinesSortBy(field PipelineSortField) ListPipelinesOption {
	return listPipelinesOptionFunc(func(cfg *listPipelinesOptions) {
		value := genv2.GetPipelinesParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithPipelinesSortDirection(direction SortDirection) ListPipelinesOption {
	return listPipelinesOptionFunc(func(cfg *listPipelinesOptions) {
		value := genv2.GetPipelinesParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithPipelinesPageSize(limit int) ListPipelinesOption {
	return listPipelinesOptionFunc(func(cfg *listPipelinesOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithPipelinesCursor(cursor string) ListPipelinesOption {
	return listPipelinesOptionFunc(func(cfg *listPipelinesOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithPipelineName(name string) PipelineOption {
	return pipelineFieldOption(func(payload *pipelinePayload) {
		payload.name = &name
	})
}

func WithPipelineDealProbabilityEnabled(enabled bool) PipelineOption {
	return pipelineFieldOption(func(payload *pipelinePayload) {
		payload.dealProbabilityEnabled = &enabled
	})
}

func newGetPipelineOptions(opts []GetPipelineOption) getPipelineOptions {
	var cfg getPipelineOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetPipeline(&cfg)
	}
	return cfg
}

func newListPipelinesOptions(opts []ListPipelinesOption) listPipelinesOptions {
	var cfg listPipelinesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListPipelines(&cfg)
	}
	return cfg
}

func newCreatePipelineOptions(opts []CreatePipelineOption) createPipelineOptions {
	var cfg createPipelineOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreatePipeline(&cfg)
	}
	return cfg
}

func newUpdatePipelineOptions(opts []UpdatePipelineOption) updatePipelineOptions {
	var cfg updatePipelineOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdatePipeline(&cfg)
	}
	return cfg
}

func newDeletePipelineOptions(opts []DeletePipelineOption) deletePipelineOptions {
	var cfg deletePipelineOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeletePipeline(&cfg)
	}
	return cfg
}

func (s *PipelinesService) List(ctx context.Context, opts ...ListPipelinesOption) ([]Pipeline, *string, error) {
	cfg := newListPipelinesOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *PipelinesService) ListPager(opts ...ListPipelinesOption) *pipedrive.CursorPager[Pipeline] {
	cfg := newListPipelinesOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Pipeline, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *PipelinesService) ForEach(ctx context.Context, fn func(Pipeline) error, opts ...ListPipelinesOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *PipelinesService) Get(ctx context.Context, id PipelineID, opts ...GetPipelineOption) (*Pipeline, error) {
	cfg := newGetPipelineOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetPipelineWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Pipeline `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing pipeline data in response")
	}
	return payload.Data, nil
}

func (s *PipelinesService) Create(ctx context.Context, opts ...CreatePipelineOption) (*Pipeline, error) {
	cfg := newCreatePipelineOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body := genv2.AddPipelineJSONBody{}
	if cfg.payload.name != nil {
		body.Name = *cfg.payload.name
	}
	if cfg.payload.dealProbabilityEnabled != nil {
		body.IsDealProbabilityEnabled = cfg.payload.dealProbabilityEnabled
	}

	resp, err := s.client.gen.AddPipelineWithResponse(ctx, body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Pipeline `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing pipeline data in response")
	}
	return payload.Data, nil
}

func (s *PipelinesService) Update(ctx context.Context, id PipelineID, opts ...UpdatePipelineOption) (*Pipeline, error) {
	cfg := newUpdatePipelineOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body := genv2.UpdatePipelineJSONBody{}
	if cfg.payload.name != nil {
		body.Name = cfg.payload.name
	}
	if cfg.payload.dealProbabilityEnabled != nil {
		body.IsDealProbabilityEnabled = cfg.payload.dealProbabilityEnabled
	}

	resp, err := s.client.gen.UpdatePipelineWithResponse(ctx, int(id), body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Pipeline `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing pipeline data in response")
	}
	return payload.Data, nil
}

type PipelineDeleteResult struct {
	ID PipelineID `json:"id"`
}

func (s *PipelinesService) Delete(ctx context.Context, id PipelineID, opts ...DeletePipelineOption) (*PipelineDeleteResult, error) {
	cfg := newDeletePipelineOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeletePipelineWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *PipelineDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing pipeline data in response")
	}
	return payload.Data, nil
}

func (s *PipelinesService) list(ctx context.Context, params genv2.GetPipelinesParams, requestOptions []pipedrive.RequestOption) ([]Pipeline, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetPipelinesWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Pipeline `json:"data"`
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
