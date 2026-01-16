package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type BillingAddon struct {
	Code string `json:"code,omitempty"`
}

type BillingService struct {
	client *Client
}

type ListBillingAddonsOption interface {
	applyListBillingAddons(*listBillingAddonsOptions)
}

type BillingRequestOption interface {
	ListBillingAddonsOption
}

type listBillingAddonsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type billingRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o billingRequestOptions) applyListBillingAddons(cfg *listBillingAddonsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listBillingAddonsOptionFunc func(*listBillingAddonsOptions)

func (f listBillingAddonsOptionFunc) applyListBillingAddons(cfg *listBillingAddonsOptions) {
	f(cfg)
}

func WithBillingRequestOptions(opts ...pipedrive.RequestOption) BillingRequestOption {
	return billingRequestOptions{requestOptions: opts}
}

func newListBillingAddonsOptions(opts []ListBillingAddonsOption) listBillingAddonsOptions {
	var cfg listBillingAddonsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListBillingAddons(&cfg)
	}
	return cfg
}

func (s *BillingService) ListAddons(ctx context.Context, opts ...ListBillingAddonsOption) ([]BillingAddon, error) {
	cfg := newListBillingAddonsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetCompanyAddons(ctx, toRequestEditors(editors)...)
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
		Data []BillingAddon `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}
