package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ProductFieldsService struct {
	client *Client
}

type DeleteProductFieldsOption interface {
	applyDeleteProductFields(*deleteProductFieldsOptions)
}

type ProductFieldsRequestOption interface {
	DeleteProductFieldsOption
}

type deleteProductFieldsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type productFieldsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o productFieldsRequestOptions) applyDeleteProductFields(cfg *deleteProductFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type deleteProductFieldsOptionFunc func(*deleteProductFieldsOptions)

func (f deleteProductFieldsOptionFunc) applyDeleteProductFields(cfg *deleteProductFieldsOptions) {
	f(cfg)
}

func WithProductFieldsRequestOptions(opts ...pipedrive.RequestOption) ProductFieldsRequestOption {
	return productFieldsRequestOptions{requestOptions: opts}
}

func newDeleteProductFieldsOptions(opts []DeleteProductFieldsOption) deleteProductFieldsOptions {
	var cfg deleteProductFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteProductFields(&cfg)
	}
	return cfg
}

func (s *ProductFieldsService) Delete(ctx context.Context, ids []FieldID, opts ...DeleteProductFieldsOption) (*FieldDeleteResult, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("field IDs are required")
	}
	cfg := newDeleteProductFieldsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	params := genv1.DeleteProductFieldsParams{Ids: joinIDs(ids)}
	resp, err := s.client.gen.DeleteProductFields(ctx, &params, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing delete product fields data in response")
	}
	return &FieldDeleteResult{IDs: payload.Data.IDs}, nil
}
