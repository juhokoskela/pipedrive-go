package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Organization struct {
	ID         OrganizationID `json:"id,omitempty"`
	Name       string         `json:"name,omitempty"`
	OwnerID    *UserID        `json:"owner_id,omitempty"`
	AddTime    *DateTime      `json:"add_time,omitempty"`
	UpdateTime *DateTime      `json:"update_time,omitempty"`
}

type OrganizationsService struct {
	client *Client
}

type OrganizationsOption interface {
	applyOrganizations(*organizationsOptions)
}

type organizationsOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type organizationsOptionFunc func(*organizationsOptions)

func (f organizationsOptionFunc) applyOrganizations(cfg *organizationsOptions) {
	f(cfg)
}

func WithOrganizationsQuery(values url.Values) OrganizationsOption {
	return organizationsOptionFunc(func(cfg *organizationsOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithOrganizationsRequestOptions(opts ...pipedrive.RequestOption) OrganizationsOption {
	return organizationsOptionFunc(func(cfg *organizationsOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newOrganizationsOptions(opts []OrganizationsOption) organizationsOptions {
	var cfg organizationsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyOrganizations(&cfg)
	}
	return cfg
}

func (s *OrganizationsService) Merge(ctx context.Context, id OrganizationID, mergeWithID OrganizationID, opts ...OrganizationsOption) (*Organization, error) {
	cfg := newOrganizationsOptions(opts)
	path := fmt.Sprintf("/organizations/%d/merge", id)

	body := map[string]any{"merge_with_id": int(mergeWithID)}
	var payload struct {
		Data *Organization `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, body, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing merged organization data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationsService) Changelog(ctx context.Context, id OrganizationID, opts ...OrganizationsOption) ([]map[string]any, *CollectionPagination, error) {
	cfg := newOrganizationsOptions(opts)
	path := fmt.Sprintf("/organizations/%d/changelog", id)

	var payload struct {
		Data           []map[string]any      `json:"data"`
		AdditionalData *CollectionPagination `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *OrganizationsService) ListFiles(ctx context.Context, id OrganizationID, opts ...OrganizationsOption) ([]File, *Pagination, error) {
	cfg := newOrganizationsOptions(opts)
	path := fmt.Sprintf("/organizations/%d/files", id)

	var payload struct {
		Data           []File `json:"data"`
		AdditionalData *struct {
			Pagination *Pagination `json:"pagination"`
		} `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	var page *Pagination
	if payload.AdditionalData != nil {
		page = payload.AdditionalData.Pagination
	}
	return payload.Data, page, nil
}

func (s *OrganizationsService) ListMailMessages(ctx context.Context, id OrganizationID, opts ...OrganizationsOption) ([]MailMessage, *Pagination, error) {
	cfg := newOrganizationsOptions(opts)
	path := fmt.Sprintf("/organizations/%d/mailMessages", id)

	var payload struct {
		Data           []MailMessage `json:"data"`
		AdditionalData *struct {
			Pagination *Pagination `json:"pagination"`
		} `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	var page *Pagination
	if payload.AdditionalData != nil {
		page = payload.AdditionalData.Pagination
	}
	return payload.Data, page, nil
}

func (s *OrganizationsService) ListUpdates(ctx context.Context, id OrganizationID, opts ...OrganizationsOption) ([]map[string]any, *Pagination, error) {
	cfg := newOrganizationsOptions(opts)
	path := fmt.Sprintf("/organizations/%d/flow", id)

	var payload struct {
		Data           []map[string]any `json:"data"`
		AdditionalData *struct {
			Pagination *Pagination `json:"pagination"`
		} `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	var page *Pagination
	if payload.AdditionalData != nil {
		page = payload.AdditionalData.Pagination
	}
	return payload.Data, page, nil
}

func (s *OrganizationsService) ListUsers(ctx context.Context, id OrganizationID, opts ...OrganizationsOption) ([]User, error) {
	cfg := newOrganizationsOptions(opts)
	path := fmt.Sprintf("/organizations/%d/permittedUsers", id)

	var payload struct {
		Data []User `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}
