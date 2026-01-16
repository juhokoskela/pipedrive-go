package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type PersonFieldsService struct {
	client *Client
}

type DeletePersonFieldsOption interface {
	applyDeletePersonFields(*deletePersonFieldsOptions)
}

type PersonFieldsRequestOption interface {
	DeletePersonFieldsOption
}

type deletePersonFieldsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type personFieldsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o personFieldsRequestOptions) applyDeletePersonFields(cfg *deletePersonFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type deletePersonFieldsOptionFunc func(*deletePersonFieldsOptions)

func (f deletePersonFieldsOptionFunc) applyDeletePersonFields(cfg *deletePersonFieldsOptions) {
	f(cfg)
}

func WithPersonFieldsRequestOptions(opts ...pipedrive.RequestOption) PersonFieldsRequestOption {
	return personFieldsRequestOptions{requestOptions: opts}
}

func newDeletePersonFieldsOptions(opts []DeletePersonFieldsOption) deletePersonFieldsOptions {
	var cfg deletePersonFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeletePersonFields(&cfg)
	}
	return cfg
}

func (s *PersonFieldsService) Delete(ctx context.Context, ids []FieldID, opts ...DeletePersonFieldsOption) (*FieldDeleteResult, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("field IDs are required")
	}
	cfg := newDeletePersonFieldsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	params := genv1.DeletePersonFieldsParams{Ids: joinIDs(ids)}
	resp, err := s.client.gen.DeletePersonFields(ctx, &params, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing delete person fields data in response")
	}
	return &FieldDeleteResult{IDs: payload.Data.IDs}, nil
}
