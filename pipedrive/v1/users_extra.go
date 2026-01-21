package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func (s *UsersService) Create(ctx context.Context, payload map[string]any, opts ...pipedrive.RequestOption) (*User, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("user payload is required")
	}

	var resp struct {
		Data *User `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, "/users", nil, payload, &resp, opts...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing user data in response")
	}
	return resp.Data, nil
}

func (s *UsersService) Update(ctx context.Context, id UserID, payload map[string]any, opts ...pipedrive.RequestOption) (*User, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("user payload is required")
	}
	path := fmt.Sprintf("/users/%d", id)

	var resp struct {
		Data *User `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, nil, payload, &resp, opts...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing user update data in response")
	}
	return resp.Data, nil
}

func (s *UsersService) FindByName(ctx context.Context, query url.Values, opts ...pipedrive.RequestOption) ([]User, error) {
	var payload struct {
		Data []User `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/users/find", query, nil, &payload, opts...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *UsersService) ListRoleAssignments(ctx context.Context, id UserID, query url.Values, opts ...pipedrive.RequestOption) ([]map[string]any, error) {
	path := fmt.Sprintf("/users/%d/roleAssignments", id)

	var payload struct {
		Data []map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, query, nil, &payload, opts...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *UsersService) ListRoleSettings(ctx context.Context, id UserID, opts ...pipedrive.RequestOption) ([]map[string]any, error) {
	path := fmt.Sprintf("/users/%d/roleSettings", id)

	var payload struct {
		Data []map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, nil, nil, &payload, opts...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *UsersService) ListTeams(ctx context.Context, id UserID, query url.Values, opts ...pipedrive.RequestOption) ([]Team, error) {
	path := fmt.Sprintf("/legacyTeams/user/%d", id)

	var payload struct {
		Data []Team `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, query, nil, &payload, opts...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}
