package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type PipelineStageConversion struct {
	FromStageID    StageID `json:"from_stage_id,omitempty"`
	ToStageID      StageID `json:"to_stage_id,omitempty"`
	ConversionRate int     `json:"conversion_rate,omitempty"`
}

type PipelineConversionStatistics struct {
	StageConversions []PipelineStageConversion `json:"stage_conversions,omitempty"`
	WonConversion    int                       `json:"won_conversion,omitempty"`
	LostConversion   int                       `json:"lost_conversion,omitempty"`
}

type PipelineMovementCount struct {
	Count int `json:"count,omitempty"`
}

type PipelineDealSummary struct {
	Count           int                `json:"count,omitempty"`
	DealIDs         []DealID           `json:"deals_ids,omitempty"`
	Values          map[string]float64 `json:"values,omitempty"`
	FormattedValues map[string]string  `json:"formatted_values,omitempty"`
}

type PipelineStageAge struct {
	StageID StageID `json:"stage_id,omitempty"`
	Value   int     `json:"value,omitempty"`
}

type PipelineAverageAge struct {
	AcrossAllStages int                `json:"across_all_stages,omitempty"`
	ByStages        []PipelineStageAge `json:"by_stages,omitempty"`
}

type PipelineMovementStatistics struct {
	MovementsBetweenStages PipelineMovementCount `json:"movements_between_stages,omitempty"`
	NewDeals               PipelineDealSummary   `json:"new_deals,omitempty"`
	DealsLeftOpen          PipelineDealSummary   `json:"deals_left_open,omitempty"`
	WonDeals               PipelineDealSummary   `json:"won_deals,omitempty"`
	LostDeals              PipelineDealSummary   `json:"lost_deals,omitempty"`
	AverageAgeInDays       PipelineAverageAge    `json:"average_age_in_days,omitempty"`
}

type PipelineDealsAdditionalData struct {
	Pagination
	Summary      map[string]any `json:"summary,omitempty"`
	DealsSummary map[string]any `json:"deals_summary,omitempty"`
}

type PipelinesService struct {
	client *Client
}

type GetPipelineConversionStatisticsOption interface {
	applyGetPipelineConversionStatistics(*getPipelineConversionStatisticsOptions)
}

type GetPipelineMovementStatisticsOption interface {
	applyGetPipelineMovementStatistics(*getPipelineMovementStatisticsOptions)
}

type PipelineDealsOption interface {
	applyPipelineDeals(*pipelineDealsOptions)
}

type PipelinesRequestOption interface {
	GetPipelineConversionStatisticsOption
	GetPipelineMovementStatisticsOption
}

type getPipelineConversionStatisticsOptions struct {
	params         genv1.GetPipelineConversionStatisticsParams
	requestOptions []pipedrive.RequestOption
}

type getPipelineMovementStatisticsOptions struct {
	params         genv1.GetPipelineMovementStatisticsParams
	requestOptions []pipedrive.RequestOption
}

type pipelineDealsOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type pipelinesRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o pipelinesRequestOptions) applyGetPipelineConversionStatistics(cfg *getPipelineConversionStatisticsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o pipelinesRequestOptions) applyGetPipelineMovementStatistics(cfg *getPipelineMovementStatisticsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getPipelineConversionStatisticsOptionFunc func(*getPipelineConversionStatisticsOptions)

func (f getPipelineConversionStatisticsOptionFunc) applyGetPipelineConversionStatistics(cfg *getPipelineConversionStatisticsOptions) {
	f(cfg)
}

type getPipelineMovementStatisticsOptionFunc func(*getPipelineMovementStatisticsOptions)

func (f getPipelineMovementStatisticsOptionFunc) applyGetPipelineMovementStatistics(cfg *getPipelineMovementStatisticsOptions) {
	f(cfg)
}

type pipelineDealsOptionFunc func(*pipelineDealsOptions)

func (f pipelineDealsOptionFunc) applyPipelineDeals(cfg *pipelineDealsOptions) {
	f(cfg)
}

func WithPipelinesRequestOptions(opts ...pipedrive.RequestOption) PipelinesRequestOption {
	return pipelinesRequestOptions{requestOptions: opts}
}

func WithPipelineConversionStartDate(start time.Time) GetPipelineConversionStatisticsOption {
	return getPipelineConversionStatisticsOptionFunc(func(cfg *getPipelineConversionStatisticsOptions) {
		cfg.params.StartDate = openapi_types.Date{Time: start}
	})
}

func WithPipelineConversionEndDate(end time.Time) GetPipelineConversionStatisticsOption {
	return getPipelineConversionStatisticsOptionFunc(func(cfg *getPipelineConversionStatisticsOptions) {
		cfg.params.EndDate = openapi_types.Date{Time: end}
	})
}

func WithPipelineConversionUserID(userID UserID) GetPipelineConversionStatisticsOption {
	return getPipelineConversionStatisticsOptionFunc(func(cfg *getPipelineConversionStatisticsOptions) {
		value := int(userID)
		cfg.params.UserId = &value
	})
}

func WithPipelineMovementStartDate(start time.Time) GetPipelineMovementStatisticsOption {
	return getPipelineMovementStatisticsOptionFunc(func(cfg *getPipelineMovementStatisticsOptions) {
		cfg.params.StartDate = openapi_types.Date{Time: start}
	})
}

func WithPipelineMovementEndDate(end time.Time) GetPipelineMovementStatisticsOption {
	return getPipelineMovementStatisticsOptionFunc(func(cfg *getPipelineMovementStatisticsOptions) {
		cfg.params.EndDate = openapi_types.Date{Time: end}
	})
}

func WithPipelineMovementUserID(userID UserID) GetPipelineMovementStatisticsOption {
	return getPipelineMovementStatisticsOptionFunc(func(cfg *getPipelineMovementStatisticsOptions) {
		value := int(userID)
		cfg.params.UserId = &value
	})
}

func WithPipelineDealsQuery(values url.Values) PipelineDealsOption {
	return pipelineDealsOptionFunc(func(cfg *pipelineDealsOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithPipelineDealsRequestOptions(opts ...pipedrive.RequestOption) PipelineDealsOption {
	return pipelineDealsOptionFunc(func(cfg *pipelineDealsOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newGetPipelineConversionStatisticsOptions(opts []GetPipelineConversionStatisticsOption) getPipelineConversionStatisticsOptions {
	var cfg getPipelineConversionStatisticsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetPipelineConversionStatistics(&cfg)
	}
	return cfg
}

func newGetPipelineMovementStatisticsOptions(opts []GetPipelineMovementStatisticsOption) getPipelineMovementStatisticsOptions {
	var cfg getPipelineMovementStatisticsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetPipelineMovementStatistics(&cfg)
	}
	return cfg
}

func newPipelineDealsOptions(opts []PipelineDealsOption) pipelineDealsOptions {
	var cfg pipelineDealsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyPipelineDeals(&cfg)
	}
	return cfg
}

func (s *PipelinesService) GetConversionStatistics(ctx context.Context, id PipelineID, opts ...GetPipelineConversionStatisticsOption) (*PipelineConversionStatistics, error) {
	cfg := newGetPipelineConversionStatisticsOptions(opts)
	if cfg.params.StartDate.IsZero() || cfg.params.EndDate.IsZero() {
		return nil, fmt.Errorf("start and end dates are required")
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetPipelineConversionStatistics(ctx, int(id), &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errorFromResponse(resp, respBody)
	}

	var payload struct {
		Data *PipelineConversionStatistics `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing pipeline conversion data in response")
	}
	return payload.Data, nil
}

func (s *PipelinesService) GetMovementStatistics(ctx context.Context, id PipelineID, opts ...GetPipelineMovementStatisticsOption) (*PipelineMovementStatistics, error) {
	cfg := newGetPipelineMovementStatisticsOptions(opts)
	if cfg.params.StartDate.IsZero() || cfg.params.EndDate.IsZero() {
		return nil, fmt.Errorf("start and end dates are required")
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetPipelineMovementStatistics(ctx, int(id), &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errorFromResponse(resp, respBody)
	}

	var payload struct {
		Data *PipelineMovementStatistics `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing pipeline movement data in response")
	}
	return payload.Data, nil
}

func (s *PipelinesService) ListDeals(ctx context.Context, id PipelineID, opts ...PipelineDealsOption) ([]Deal, *PipelineDealsAdditionalData, error) {
	cfg := newPipelineDealsOptions(opts)
	path := fmt.Sprintf("/pipelines/%d/deals", id)

	var payload struct {
		Data           []Deal                       `json:"data"`
		AdditionalData *PipelineDealsAdditionalData `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	return payload.Data, payload.AdditionalData, nil
}
