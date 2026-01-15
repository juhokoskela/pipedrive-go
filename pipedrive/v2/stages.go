package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type StageID int64

type StageSortField string

const (
	StageSortByID         StageSortField = "id"
	StageSortByUpdateTime StageSortField = "update_time"
	StageSortByAddTime    StageSortField = "add_time"
	StageSortByOrder      StageSortField = "order_nr"
)

type Stage struct {
	ID              StageID    `json:"id"`
	Order           int        `json:"order_nr,omitempty"`
	Name            string     `json:"name,omitempty"`
	IsDeleted       bool       `json:"is_deleted,omitempty"`
	DealProbability int        `json:"deal_probability,omitempty"`
	PipelineID      PipelineID `json:"pipeline_id,omitempty"`
	DealRotEnabled  bool       `json:"is_deal_rot_enabled,omitempty"`
	DaysToRotten    *int       `json:"days_to_rotten,omitempty"`
	AddTime         *time.Time `json:"add_time,omitempty"`
	UpdateTime      *time.Time `json:"update_time,omitempty"`
}

type StagesService struct {
	client *Client
}

type GetStageOption interface {
	applyGetStage(*getStageOptions)
}

type ListStagesOption interface {
	applyListStages(*listStagesOptions)
}

type CreateStageOption interface {
	applyCreateStage(*createStageOptions)
}

type UpdateStageOption interface {
	applyUpdateStage(*updateStageOptions)
}

type DeleteStageOption interface {
	applyDeleteStage(*deleteStageOptions)
}

type StageRequestOption interface {
	GetStageOption
	ListStagesOption
	CreateStageOption
	UpdateStageOption
	DeleteStageOption
}

type StageOption interface {
	CreateStageOption
	UpdateStageOption
}

type getStageOptions struct {
	requestOptions []pipedrive.RequestOption
}

type listStagesOptions struct {
	params         genv2.GetStagesParams
	requestOptions []pipedrive.RequestOption
}

type createStageOptions struct {
	payload        stagePayload
	requestOptions []pipedrive.RequestOption
}

type updateStageOptions struct {
	payload        stagePayload
	requestOptions []pipedrive.RequestOption
}

type deleteStageOptions struct {
	requestOptions []pipedrive.RequestOption
}

type stagePayload struct {
	name            *string
	pipelineID      *PipelineID
	dealProbability *int
	dealRotEnabled  *bool
	daysToRotten    *int
}

type stageRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o stageRequestOptions) applyGetStage(cfg *getStageOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o stageRequestOptions) applyListStages(cfg *listStagesOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o stageRequestOptions) applyCreateStage(cfg *createStageOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o stageRequestOptions) applyUpdateStage(cfg *updateStageOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o stageRequestOptions) applyDeleteStage(cfg *deleteStageOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listStagesOptionFunc func(*listStagesOptions)

func (f listStagesOptionFunc) applyListStages(cfg *listStagesOptions) {
	f(cfg)
}

type stageFieldOption func(*stagePayload)

func (f stageFieldOption) applyCreateStage(cfg *createStageOptions) {
	f(&cfg.payload)
}

func (f stageFieldOption) applyUpdateStage(cfg *updateStageOptions) {
	f(&cfg.payload)
}

func WithStageRequestOptions(opts ...pipedrive.RequestOption) StageRequestOption {
	return stageRequestOptions{requestOptions: opts}
}

func WithStagesPipelineID(id PipelineID) ListStagesOption {
	return listStagesOptionFunc(func(cfg *listStagesOptions) {
		value := int(id)
		cfg.params.PipelineId = &value
	})
}

func WithStagesSortBy(field StageSortField) ListStagesOption {
	return listStagesOptionFunc(func(cfg *listStagesOptions) {
		value := genv2.GetStagesParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithStagesSortDirection(direction SortDirection) ListStagesOption {
	return listStagesOptionFunc(func(cfg *listStagesOptions) {
		value := genv2.GetStagesParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithStagesPageSize(limit int) ListStagesOption {
	return listStagesOptionFunc(func(cfg *listStagesOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithStagesCursor(cursor string) ListStagesOption {
	return listStagesOptionFunc(func(cfg *listStagesOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithStageName(name string) StageOption {
	return stageFieldOption(func(payload *stagePayload) {
		payload.name = &name
	})
}

func WithStagePipelineID(id PipelineID) StageOption {
	return stageFieldOption(func(payload *stagePayload) {
		payload.pipelineID = &id
	})
}

func WithStageDealProbability(probability int) StageOption {
	return stageFieldOption(func(payload *stagePayload) {
		payload.dealProbability = &probability
	})
}

func WithStageDealRotEnabled(enabled bool) StageOption {
	return stageFieldOption(func(payload *stagePayload) {
		payload.dealRotEnabled = &enabled
	})
}

func WithStageDaysToRotten(days int) StageOption {
	return stageFieldOption(func(payload *stagePayload) {
		payload.daysToRotten = &days
	})
}

func newGetStageOptions(opts []GetStageOption) getStageOptions {
	var cfg getStageOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetStage(&cfg)
	}
	return cfg
}

func newListStagesOptions(opts []ListStagesOption) listStagesOptions {
	var cfg listStagesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListStages(&cfg)
	}
	return cfg
}

func newCreateStageOptions(opts []CreateStageOption) createStageOptions {
	var cfg createStageOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateStage(&cfg)
	}
	return cfg
}

func newUpdateStageOptions(opts []UpdateStageOption) updateStageOptions {
	var cfg updateStageOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateStage(&cfg)
	}
	return cfg
}

func newDeleteStageOptions(opts []DeleteStageOption) deleteStageOptions {
	var cfg deleteStageOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteStage(&cfg)
	}
	return cfg
}

func (s *StagesService) List(ctx context.Context, opts ...ListStagesOption) ([]Stage, *string, error) {
	cfg := newListStagesOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *StagesService) ListPager(opts ...ListStagesOption) *pipedrive.CursorPager[Stage] {
	cfg := newListStagesOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Stage, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *StagesService) ForEach(ctx context.Context, fn func(Stage) error, opts ...ListStagesOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *StagesService) Get(ctx context.Context, id StageID, opts ...GetStageOption) (*Stage, error) {
	cfg := newGetStageOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetStageWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Stage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing stage data in response")
	}
	return payload.Data, nil
}

func (s *StagesService) Create(ctx context.Context, opts ...CreateStageOption) (*Stage, error) {
	cfg := newCreateStageOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body := genv2.AddStageJSONRequestBody{}
	if cfg.payload.name != nil {
		body.Name = *cfg.payload.name
	}
	if cfg.payload.pipelineID != nil {
		body.PipelineId = int(*cfg.payload.pipelineID)
	}
	if cfg.payload.dealProbability != nil {
		body.DealProbability = cfg.payload.dealProbability
	}
	if cfg.payload.dealRotEnabled != nil {
		body.IsDealRotEnabled = cfg.payload.dealRotEnabled
	}
	if cfg.payload.daysToRotten != nil {
		body.DaysToRotten = cfg.payload.daysToRotten
	}

	resp, err := s.client.gen.AddStageWithResponse(ctx, body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Stage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing stage data in response")
	}
	return payload.Data, nil
}

func (s *StagesService) Update(ctx context.Context, id StageID, opts ...UpdateStageOption) (*Stage, error) {
	cfg := newUpdateStageOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body := genv2.UpdateStageJSONRequestBody{}
	if cfg.payload.name != nil {
		body.Name = cfg.payload.name
	}
	if cfg.payload.pipelineID != nil {
		id := int(*cfg.payload.pipelineID)
		body.PipelineId = &id
	}
	if cfg.payload.dealProbability != nil {
		body.DealProbability = cfg.payload.dealProbability
	}
	if cfg.payload.dealRotEnabled != nil {
		body.IsDealRotEnabled = cfg.payload.dealRotEnabled
	}
	if cfg.payload.daysToRotten != nil {
		body.DaysToRotten = cfg.payload.daysToRotten
	}

	resp, err := s.client.gen.UpdateStageWithResponse(ctx, int(id), body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Stage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing stage data in response")
	}
	return payload.Data, nil
}

type StageDeleteResult struct {
	ID StageID `json:"id"`
}

func (s *StagesService) Delete(ctx context.Context, id StageID, opts ...DeleteStageOption) (*StageDeleteResult, error) {
	cfg := newDeleteStageOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteStageWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *StageDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing stage data in response")
	}
	return payload.Data, nil
}

func (s *StagesService) list(ctx context.Context, params genv2.GetStagesParams, requestOptions []pipedrive.RequestOption) ([]Stage, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetStagesWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Stage `json:"data"`
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
