package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Role struct {
	ID          RoleID `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type RoleAssignment struct {
	UserID UserID `json:"user_id,omitempty"`
	RoleID RoleID `json:"role_id,omitempty"`
}

type RolesService struct {
	client *Client
}

type RolesOption interface {
	applyRoles(*rolesOptions)
}

type rolesOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type rolesOptionFunc func(*rolesOptions)

func (f rolesOptionFunc) applyRoles(cfg *rolesOptions) {
	f(cfg)
}

func WithRolesQuery(values url.Values) RolesOption {
	return rolesOptionFunc(func(cfg *rolesOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithRolesRequestOptions(opts ...pipedrive.RequestOption) RolesOption {
	return rolesOptionFunc(func(cfg *rolesOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newRolesOptions(opts []RolesOption) rolesOptions {
	var cfg rolesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyRoles(&cfg)
	}
	return cfg
}

func (s *RolesService) List(ctx context.Context, opts ...RolesOption) ([]Role, error) {
	cfg := newRolesOptions(opts)

	var payload struct {
		Data []Role `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/roles", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *RolesService) Get(ctx context.Context, id RoleID, opts ...RolesOption) (*Role, error) {
	cfg := newRolesOptions(opts)
	path := fmt.Sprintf("/roles/%d", id)

	var payload struct {
		Data *Role `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing role data in response")
	}
	return payload.Data, nil
}

func (s *RolesService) Create(ctx context.Context, payload map[string]any, opts ...RolesOption) (*Role, error) {
	cfg := newRolesOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("role payload is required")
	}

	var resp struct {
		Data *Role `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, "/roles", cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing role data in response")
	}
	return resp.Data, nil
}

func (s *RolesService) Update(ctx context.Context, id RoleID, payload map[string]any, opts ...RolesOption) (*Role, error) {
	cfg := newRolesOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("role payload is required")
	}
	path := fmt.Sprintf("/roles/%d", id)

	var resp struct {
		Data *Role `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing role update data in response")
	}
	return resp.Data, nil
}

func (s *RolesService) Delete(ctx context.Context, id RoleID, opts ...RolesOption) (bool, error) {
	cfg := newRolesOptions(opts)
	path := fmt.Sprintf("/roles/%d", id)

	var resp struct {
		Success *bool `json:"success"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, path, cfg.query, nil, &resp, cfg.requestOptions...); err != nil {
		return false, err
	}
	if resp.Success == nil {
		return false, fmt.Errorf("missing role delete success in response")
	}
	return *resp.Success, nil
}

func (s *RolesService) ListAssignments(ctx context.Context, id RoleID, opts ...RolesOption) ([]RoleAssignment, *Pagination, error) {
	cfg := newRolesOptions(opts)
	path := fmt.Sprintf("/roles/%d/assignments", id)

	var payload struct {
		Data           []RoleAssignment `json:"data"`
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

func (s *RolesService) AddAssignment(ctx context.Context, id RoleID, userID UserID, opts ...RolesOption) (*RoleAssignment, error) {
	cfg := newRolesOptions(opts)
	path := fmt.Sprintf("/roles/%d/assignments", id)

	body := map[string]any{"user_id": int(userID)}
	var payload struct {
		Data *RoleAssignment `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, path, cfg.query, body, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing role assignment data in response")
	}
	return payload.Data, nil
}

func (s *RolesService) DeleteAssignment(ctx context.Context, id RoleID, userID UserID, opts ...RolesOption) (bool, error) {
	cfg := newRolesOptions(opts)
	path := fmt.Sprintf("/roles/%d/assignments", id)

	body := map[string]any{"user_id": int(userID)}
	var payload struct {
		Success *bool `json:"success"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, path, cfg.query, body, &payload, cfg.requestOptions...); err != nil {
		return false, err
	}
	if payload.Success == nil {
		return false, fmt.Errorf("missing role assignment delete success in response")
	}
	return *payload.Success, nil
}

func (s *RolesService) ListPipelines(ctx context.Context, id RoleID, opts ...RolesOption) ([]map[string]any, error) {
	cfg := newRolesOptions(opts)
	path := fmt.Sprintf("/roles/%d/pipelines", id)

	var payload struct {
		Data []map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *RolesService) UpdatePipelines(ctx context.Context, id RoleID, payload map[string]any, opts ...RolesOption) (map[string]any, error) {
	cfg := newRolesOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("role pipelines payload is required")
	}
	path := fmt.Sprintf("/roles/%d/pipelines", id)

	var resp struct {
		Data map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing role pipelines data in response")
	}
	return resp.Data, nil
}

func (s *RolesService) ListSettings(ctx context.Context, id RoleID, opts ...RolesOption) ([]map[string]any, error) {
	cfg := newRolesOptions(opts)
	path := fmt.Sprintf("/roles/%d/settings", id)

	var payload struct {
		Data []map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *RolesService) UpsertSetting(ctx context.Context, id RoleID, payload map[string]any, opts ...RolesOption) (map[string]any, error) {
	cfg := newRolesOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("role setting payload is required")
	}
	path := fmt.Sprintf("/roles/%d/settings", id)

	var resp struct {
		Data map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, path, cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing role setting data in response")
	}
	return resp.Data, nil
}
