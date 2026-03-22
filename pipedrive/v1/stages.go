package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type StagesService struct {
	client *Client
}

type StageDealsOption interface {
	applyStageDeals(*stageDealsOptions)
}

type stageDealsOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type stageDealsOptionFunc func(*stageDealsOptions)

func (f stageDealsOptionFunc) applyStageDeals(cfg *stageDealsOptions) {
	f(cfg)
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
