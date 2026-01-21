package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Team struct {
	ID          TeamID `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type TeamsService struct {
	client *Client
}

type TeamsOption interface {
	applyTeams(*teamsOptions)
}

type teamsOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type teamsOptionFunc func(*teamsOptions)

func (f teamsOptionFunc) applyTeams(cfg *teamsOptions) {
	f(cfg)
}

func WithTeamsQuery(values url.Values) TeamsOption {
	return teamsOptionFunc(func(cfg *teamsOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithTeamsRequestOptions(opts ...pipedrive.RequestOption) TeamsOption {
	return teamsOptionFunc(func(cfg *teamsOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newTeamsOptions(opts []TeamsOption) teamsOptions {
	var cfg teamsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyTeams(&cfg)
	}
	return cfg
}

func (s *TeamsService) List(ctx context.Context, opts ...TeamsOption) ([]Team, error) {
	cfg := newTeamsOptions(opts)

	var payload struct {
		Data []Team `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/legacyTeams", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *TeamsService) Get(ctx context.Context, id TeamID, opts ...TeamsOption) (*Team, error) {
	cfg := newTeamsOptions(opts)
	path := fmt.Sprintf("/legacyTeams/%d", id)

	var payload struct {
		Data *Team `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing team data in response")
	}
	return payload.Data, nil
}

func (s *TeamsService) Create(ctx context.Context, payload map[string]any, opts ...TeamsOption) (*Team, error) {
	cfg := newTeamsOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("team payload is required")
	}

	var resp struct {
		Data *Team `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, "/legacyTeams", cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing team data in response")
	}
	return resp.Data, nil
}

func (s *TeamsService) Update(ctx context.Context, id TeamID, payload map[string]any, opts ...TeamsOption) (*Team, error) {
	cfg := newTeamsOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("team payload is required")
	}
	path := fmt.Sprintf("/legacyTeams/%d", id)

	var resp struct {
		Data *Team `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing team update data in response")
	}
	return resp.Data, nil
}

func (s *TeamsService) ListUsers(ctx context.Context, id TeamID, opts ...TeamsOption) ([]UserID, error) {
	cfg := newTeamsOptions(opts)
	path := fmt.Sprintf("/legacyTeams/%d/users", id)

	var payload struct {
		Data []UserID `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *TeamsService) AddUsers(ctx context.Context, id TeamID, userIDs []UserID, opts ...TeamsOption) ([]UserID, error) {
	cfg := newTeamsOptions(opts)
	if len(userIDs) == 0 {
		return nil, fmt.Errorf("at least one user id is required")
	}
	path := fmt.Sprintf("/legacyTeams/%d/users", id)

	body := map[string]any{"users": userIDs}
	var payload struct {
		Data []UserID `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, path, cfg.query, body, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *TeamsService) DeleteUsers(ctx context.Context, id TeamID, userIDs []UserID, opts ...TeamsOption) ([]UserID, error) {
	cfg := newTeamsOptions(opts)
	if len(userIDs) == 0 {
		return nil, fmt.Errorf("at least one user id is required")
	}
	path := fmt.Sprintf("/legacyTeams/%d/users", id)

	body := map[string]any{"users": userIDs}
	var payload struct {
		Data []UserID `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, path, cfg.query, body, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}
