package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type OrganizationFieldsService struct {
	client *Client
}

type DeleteOrganizationFieldsOption interface {
	applyDeleteOrganizationFields(*deleteOrganizationFieldsOptions)
}

type OrganizationFieldsRequestOption interface {
	DeleteOrganizationFieldsOption
}

type deleteOrganizationFieldsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type organizationFieldsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o organizationFieldsRequestOptions) applyDeleteOrganizationFields(cfg *deleteOrganizationFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type deleteOrganizationFieldsOptionFunc func(*deleteOrganizationFieldsOptions)

func (f deleteOrganizationFieldsOptionFunc) applyDeleteOrganizationFields(cfg *deleteOrganizationFieldsOptions) {
	f(cfg)
}

func WithOrganizationFieldsRequestOptions(opts ...pipedrive.RequestOption) OrganizationFieldsRequestOption {
	return organizationFieldsRequestOptions{requestOptions: opts}
}

func newDeleteOrganizationFieldsOptions(opts []DeleteOrganizationFieldsOption) deleteOrganizationFieldsOptions {
	var cfg deleteOrganizationFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteOrganizationFields(&cfg)
	}
	return cfg
}

func (s *OrganizationFieldsService) Delete(ctx context.Context, ids []FieldID, opts ...DeleteOrganizationFieldsOption) (*FieldDeleteResult, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("field IDs are required")
	}
	cfg := newDeleteOrganizationFieldsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	params := genv1.DeleteOrganizationFieldsParams{Ids: joinIDs(ids)}
	resp, err := s.client.gen.DeleteOrganizationFields(ctx, &params, toRequestEditors(editors)...)
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
		return nil, fmt.Errorf("missing delete organization fields data in response")
	}
	return &FieldDeleteResult{IDs: payload.Data.IDs}, nil
}
