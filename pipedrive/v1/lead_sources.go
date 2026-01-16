package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type LeadSource struct {
	Name string `json:"name,omitempty"`
}

type LeadSourcesService struct {
	client *Client
}

type ListLeadSourcesOption interface {
	applyListLeadSources(*listLeadSourcesOptions)
}

type LeadSourcesRequestOption interface {
	ListLeadSourcesOption
}

type listLeadSourcesOptions struct {
	requestOptions []pipedrive.RequestOption
}

type leadSourcesRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o leadSourcesRequestOptions) applyListLeadSources(cfg *listLeadSourcesOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listLeadSourcesOptionFunc func(*listLeadSourcesOptions)

func (f listLeadSourcesOptionFunc) applyListLeadSources(cfg *listLeadSourcesOptions) {
	f(cfg)
}

func WithLeadSourcesRequestOptions(opts ...pipedrive.RequestOption) LeadSourcesRequestOption {
	return leadSourcesRequestOptions{requestOptions: opts}
}

func newListLeadSourcesOptions(opts []ListLeadSourcesOption) listLeadSourcesOptions {
	var cfg listLeadSourcesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListLeadSources(&cfg)
	}
	return cfg
}

func (s *LeadSourcesService) List(ctx context.Context, opts ...ListLeadSourcesOption) ([]LeadSource, error) {
	cfg := newListLeadSourcesOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetLeadSources(ctx, toRequestEditors(editors)...)
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
		Data []LeadSource `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}
