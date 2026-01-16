package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type UserConnections struct {
	GoogleID *string `json:"google,omitempty"`
}

type UserConnectionsService struct {
	client *Client
}

type GetUserConnectionsOption interface {
	applyGetUserConnections(*getUserConnectionsOptions)
}

type UserConnectionsRequestOption interface {
	GetUserConnectionsOption
}

type getUserConnectionsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type userConnectionsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o userConnectionsRequestOptions) applyGetUserConnections(cfg *getUserConnectionsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getUserConnectionsOptionFunc func(*getUserConnectionsOptions)

func (f getUserConnectionsOptionFunc) applyGetUserConnections(cfg *getUserConnectionsOptions) {
	f(cfg)
}

func WithUserConnectionsRequestOptions(opts ...pipedrive.RequestOption) UserConnectionsRequestOption {
	return userConnectionsRequestOptions{requestOptions: opts}
}

func newGetUserConnectionsOptions(opts []GetUserConnectionsOption) getUserConnectionsOptions {
	var cfg getUserConnectionsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetUserConnections(&cfg)
	}
	return cfg
}

func (s *UserConnectionsService) Get(ctx context.Context, opts ...GetUserConnectionsOption) (*UserConnections, error) {
	cfg := newGetUserConnectionsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetUserConnections(ctx, toRequestEditors(editors)...)
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
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	connections := &UserConnections{}
	if payload.Data != nil {
		if raw, ok := payload.Data["google"]; ok {
			value, err := parseStringOrFalse(raw)
			if err != nil {
				return nil, fmt.Errorf("decode response: %w", err)
			}
			connections.GoogleID = value
		}
	}

	return connections, nil
}

func parseStringOrFalse(raw json.RawMessage) (*string, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) || bytes.Equal(raw, []byte("false")) {
		return nil, nil
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	return &value, nil
}
