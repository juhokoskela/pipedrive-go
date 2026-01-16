package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type LeadFieldsService struct {
	client *Client
}

type ListLeadFieldsOption interface {
	applyListLeadFields(*listLeadFieldsOptions)
}

type LeadFieldsRequestOption interface {
	ListLeadFieldsOption
}

type listLeadFieldsOptions struct {
	params         genv1.GetLeadFieldsParams
	requestOptions []pipedrive.RequestOption
}

type leadFieldsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o leadFieldsRequestOptions) applyListLeadFields(cfg *listLeadFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listLeadFieldsOptionFunc func(*listLeadFieldsOptions)

func (f listLeadFieldsOptionFunc) applyListLeadFields(cfg *listLeadFieldsOptions) {
	f(cfg)
}

func WithLeadFieldsRequestOptions(opts ...pipedrive.RequestOption) LeadFieldsRequestOption {
	return leadFieldsRequestOptions{requestOptions: opts}
}

func WithLeadFieldsStart(start int) ListLeadFieldsOption {
	return listLeadFieldsOptionFunc(func(cfg *listLeadFieldsOptions) {
		cfg.params.Start = &start
	})
}

func WithLeadFieldsLimit(limit int) ListLeadFieldsOption {
	return listLeadFieldsOptionFunc(func(cfg *listLeadFieldsOptions) {
		cfg.params.Limit = &limit
	})
}

func newListLeadFieldsOptions(opts []ListLeadFieldsOption) listLeadFieldsOptions {
	var cfg listLeadFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListLeadFields(&cfg)
	}
	return cfg
}

func (s *LeadFieldsService) List(ctx context.Context, opts ...ListLeadFieldsOption) ([]Field, *FieldPagination, error) {
	cfg := newListLeadFieldsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetLeadFields(ctx, &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp, respBody)
	}

	var payload struct {
		Data           []Field          `json:"data"`
		AdditionalData *FieldPagination `json:"additional_data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, payload.AdditionalData, nil
}
