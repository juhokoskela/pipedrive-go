package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type DealFieldsService struct {
	client *Client
}

type DeleteDealFieldsOption interface {
	applyDeleteDealFields(*deleteDealFieldsOptions)
}

type DealFieldsRequestOption interface {
	DeleteDealFieldsOption
}

type deleteDealFieldsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type dealFieldsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o dealFieldsRequestOptions) applyDeleteDealFields(cfg *deleteDealFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type deleteDealFieldsOptionFunc func(*deleteDealFieldsOptions)

func (f deleteDealFieldsOptionFunc) applyDeleteDealFields(cfg *deleteDealFieldsOptions) {
	f(cfg)
}

func WithDealFieldsRequestOptions(opts ...pipedrive.RequestOption) DealFieldsRequestOption {
	return dealFieldsRequestOptions{requestOptions: opts}
}

func newDeleteDealFieldsOptions(opts []DeleteDealFieldsOption) deleteDealFieldsOptions {
	var cfg deleteDealFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteDealFields(&cfg)
	}
	return cfg
}

func (s *DealFieldsService) Delete(ctx context.Context, ids []FieldID, opts ...DeleteDealFieldsOption) (*FieldDeleteResult, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("field IDs are required")
	}
	cfg := newDeleteDealFieldsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	params := genv1.DeleteDealFieldsParams{Ids: joinIDs(ids)}
	resp, err := s.client.gen.DeleteDealFields(ctx, &params, toRequestEditors(editors)...)
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
			IDs []FieldID `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing delete deal fields data in response")
	}
	return &FieldDeleteResult{IDs: payload.Data.IDs}, nil
}
