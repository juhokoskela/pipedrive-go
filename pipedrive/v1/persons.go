package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Person struct {
	ID         PersonID        `json:"id,omitempty"`
	Name       string          `json:"name,omitempty"`
	FirstName  string          `json:"first_name,omitempty"`
	LastName   string          `json:"last_name,omitempty"`
	OwnerID    *UserID         `json:"owner_id,omitempty"`
	OrgID      *OrganizationID `json:"org_id,omitempty"`
	AddTime    *DateTime       `json:"add_time,omitempty"`
	UpdateTime *DateTime       `json:"update_time,omitempty"`
}

type PersonsService struct {
	client *Client
}

type PersonsOption interface {
	applyPersons(*personsOptions)
}

type personsOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type personsOptionFunc func(*personsOptions)

func (f personsOptionFunc) applyPersons(cfg *personsOptions) {
	f(cfg)
}

func WithPersonsQuery(values url.Values) PersonsOption {
	return personsOptionFunc(func(cfg *personsOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithPersonsRequestOptions(opts ...pipedrive.RequestOption) PersonsOption {
	return personsOptionFunc(func(cfg *personsOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newPersonsOptions(opts []PersonsOption) personsOptions {
	var cfg personsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyPersons(&cfg)
	}
	return cfg
}

func (s *PersonsService) ListCollection(ctx context.Context, opts ...PersonsOption) ([]Person, *CollectionPagination, error) {
	cfg := newPersonsOptions(opts)

	var payload struct {
		Data           []Person              `json:"data"`
		AdditionalData *CollectionPagination `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/persons/collection", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *PersonsService) Delete(ctx context.Context, ids []PersonID, opts ...PersonsOption) (bool, error) {
	cfg := newPersonsOptions(opts)
	query := mergeQueryValues(url.Values{}, cfg.query)
	if len(ids) == 0 {
		return false, fmt.Errorf("at least one person id is required")
	}
	query.Set("ids", joinIDs(ids))

	var payload struct {
		Success *bool `json:"success"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, "/persons", query, nil, &payload, cfg.requestOptions...); err != nil {
		return false, err
	}
	if payload.Success == nil {
		return false, fmt.Errorf("missing persons delete success in response")
	}
	return *payload.Success, nil
}

func (s *PersonsService) Merge(ctx context.Context, id PersonID, mergeWithID PersonID, opts ...PersonsOption) (*Person, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/merge", id)

	body := map[string]any{"merge_with_id": int(mergeWithID)}
	var payload struct {
		Data *Person `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, body, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing merged person data in response")
	}
	return payload.Data, nil
}

func (s *PersonsService) ListActivities(ctx context.Context, id PersonID, opts ...PersonsOption) ([]Activity, *Pagination, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/activities", id)

	var payload struct {
		Data           []Activity `json:"data"`
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

func (s *PersonsService) Changelog(ctx context.Context, id PersonID, opts ...PersonsOption) ([]map[string]any, *CollectionPagination, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/changelog", id)

	var payload struct {
		Data           []map[string]any      `json:"data"`
		AdditionalData *CollectionPagination `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *PersonsService) ListDeals(ctx context.Context, id PersonID, opts ...PersonsOption) ([]Deal, *Pagination, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/deals", id)

	var payload struct {
		Data           []Deal `json:"data"`
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

func (s *PersonsService) ListFiles(ctx context.Context, id PersonID, opts ...PersonsOption) ([]File, *Pagination, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/files", id)

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

func (s *PersonsService) ListMailMessages(ctx context.Context, id PersonID, opts ...PersonsOption) ([]MailMessage, *Pagination, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/mailMessages", id)

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

func (s *PersonsService) ListProducts(ctx context.Context, id PersonID, opts ...PersonsOption) ([]Product, *Pagination, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/products", id)

	var payload struct {
		Data           []Product `json:"data"`
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

func (s *PersonsService) ListUpdates(ctx context.Context, id PersonID, opts ...PersonsOption) ([]map[string]any, *Pagination, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/flow", id)

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

func (s *PersonsService) ListUsers(ctx context.Context, id PersonID, opts ...PersonsOption) ([]User, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/permittedUsers", id)

	var payload struct {
		Data []User `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *PersonsService) AddPicture(ctx context.Context, id PersonID, body io.Reader, contentType string, opts ...PersonsOption) (map[string]any, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/picture", id)

	if body == nil {
		return nil, fmt.Errorf("picture body is required")
	}
	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	if contentType != "" {
		reqOpts = append(reqOpts, pipedrive.WithHeader("Content-Type", contentType))
	}

	var payload struct {
		Data map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, path, cfg.query, body, &payload, reqOpts...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing person picture data in response")
	}
	return payload.Data, nil
}

func (s *PersonsService) DeletePicture(ctx context.Context, id PersonID, opts ...PersonsOption) (PersonID, error) {
	cfg := newPersonsOptions(opts)
	path := fmt.Sprintf("/persons/%d/picture", id)

	var payload struct {
		Data *struct {
			ID PersonID `json:"id"`
		} `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return 0, err
	}
	if payload.Data == nil {
		return 0, fmt.Errorf("missing person picture delete data in response")
	}
	return payload.Data.ID, nil
}
