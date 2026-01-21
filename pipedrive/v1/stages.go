package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type StagesDeleteResult struct {
	IDs []StageID `json:"id"`
}

type StagesService struct {
	client *Client
}

type DeleteStagesOption interface {
	applyDeleteStages(*deleteStagesOptions)
}

type StageDealsOption interface {
	applyStageDeals(*stageDealsOptions)
}

type StagesRequestOption interface {
	DeleteStagesOption
}

type deleteStagesOptions struct {
	requestOptions []pipedrive.RequestOption
}

type stageDealsOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type stagesRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o stagesRequestOptions) applyDeleteStages(cfg *deleteStagesOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type deleteStagesOptionFunc func(*deleteStagesOptions)

func (f deleteStagesOptionFunc) applyDeleteStages(cfg *deleteStagesOptions) {
	f(cfg)
}

type stageDealsOptionFunc func(*stageDealsOptions)

func (f stageDealsOptionFunc) applyStageDeals(cfg *stageDealsOptions) {
	f(cfg)
}

func WithStagesRequestOptions(opts ...pipedrive.RequestOption) StagesRequestOption {
	return stagesRequestOptions{requestOptions: opts}
}

func WithStageDealsQuery(values url.Values) StageDealsOption {
	return stageDealsOptionFunc(func(cfg *stageDealsOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithStageDealsRequestOptions(opts ...pipedrive.RequestOption) StageDealsOption {
	return stageDealsOptionFunc(func(cfg *stageDealsOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newDeleteStagesOptions(opts []DeleteStagesOption) deleteStagesOptions {
	var cfg deleteStagesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteStages(&cfg)
	}
	return cfg
}

func newStageDealsOptions(opts []StageDealsOption) stageDealsOptions {
	var cfg stageDealsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyStageDeals(&cfg)
	}
	return cfg
}

func (s *StagesService) Delete(ctx context.Context, ids []StageID, opts ...DeleteStagesOption) (*StagesDeleteResult, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("stage IDs are required")
	}
	cfg := newDeleteStagesOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	params := genv1.DeleteStagesParams{Ids: joinIDs(ids)}
	resp, err := s.client.gen.DeleteStages(ctx, &params, toRequestEditors(editors)...)
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
		Data *struct {
			IDs []StageID `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing delete stages data in response")
	}
	return &StagesDeleteResult{IDs: payload.Data.IDs}, nil
}

func (s *StagesService) ListDeals(ctx context.Context, id StageID, opts ...StageDealsOption) ([]Deal, *Pagination, error) {
	cfg := newStageDealsOptions(opts)
	path := fmt.Sprintf("/stages/%d/deals", id)

	var payload struct {
		Data           []Deal      `json:"data"`
		AdditionalData *Pagination `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	return payload.Data, payload.AdditionalData, nil
}
