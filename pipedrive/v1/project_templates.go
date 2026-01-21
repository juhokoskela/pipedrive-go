package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ProjectTemplate struct {
	ID    ProjectTemplateID `json:"id,omitempty"`
	Name  string            `json:"name,omitempty"`
	Title string            `json:"title,omitempty"`
}

type ProjectTemplatesService struct {
	client *Client
}

type ProjectTemplatesOption interface {
	applyProjectTemplates(*projectTemplatesOptions)
}

type projectTemplatesOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type projectTemplatesOptionFunc func(*projectTemplatesOptions)

func (f projectTemplatesOptionFunc) applyProjectTemplates(cfg *projectTemplatesOptions) {
	f(cfg)
}

func WithProjectTemplatesQuery(values url.Values) ProjectTemplatesOption {
	return projectTemplatesOptionFunc(func(cfg *projectTemplatesOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithProjectTemplatesRequestOptions(opts ...pipedrive.RequestOption) ProjectTemplatesOption {
	return projectTemplatesOptionFunc(func(cfg *projectTemplatesOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newProjectTemplatesOptions(opts []ProjectTemplatesOption) projectTemplatesOptions {
	var cfg projectTemplatesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyProjectTemplates(&cfg)
	}
	return cfg
}

func (s *ProjectTemplatesService) List(ctx context.Context, opts ...ProjectTemplatesOption) ([]ProjectTemplate, error) {
	cfg := newProjectTemplatesOptions(opts)

	var payload struct {
		Data []ProjectTemplate `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/projectTemplates", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *ProjectTemplatesService) Get(ctx context.Context, id ProjectTemplateID, opts ...ProjectTemplatesOption) (*ProjectTemplate, error) {
	cfg := newProjectTemplatesOptions(opts)
	path := fmt.Sprintf("/projectTemplates/%d", id)

	var payload struct {
		Data *ProjectTemplate `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing project template data in response")
	}
	return payload.Data, nil
}
