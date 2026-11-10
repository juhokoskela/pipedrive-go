package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ProjectGroup struct {
	ID   ProjectGroupID `json:"id,omitempty"`
	Name string         `json:"name,omitempty"`
}

type ProjectTask struct {
	ID    ProjectPlanTaskID `json:"id,omitempty"`
	Title string            `json:"title,omitempty"`
}

type ProjectsService struct {
	client *Client
}

type ProjectsOption interface {
	applyProjects(*projectsOptions)
}

type projectsOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type projectsOptionFunc func(*projectsOptions)

func (f projectsOptionFunc) applyProjects(cfg *projectsOptions) {
	f(cfg)
}

func WithProjectsQuery(values url.Values) ProjectsOption {
	return projectsOptionFunc(func(cfg *projectsOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithProjectsRequestOptions(opts ...pipedrive.RequestOption) ProjectsOption {
	return projectsOptionFunc(func(cfg *projectsOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newProjectsOptions(opts []ProjectsOption) projectsOptions {
	var cfg projectsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyProjects(&cfg)
	}
	return cfg
}

func (s *ProjectsService) ListActivities(ctx context.Context, id ProjectID, opts ...ProjectsOption) ([]Activity, *Pagination, error) {
	if err := validateID(id, "project id"); err != nil {
		return nil, nil, err
	}
	cfg := newProjectsOptions(opts)
	path := fmt.Sprintf("/projects/%d/activities", id)

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

func (s *ProjectsService) ListGroups(ctx context.Context, id ProjectID, opts ...ProjectsOption) ([]ProjectGroup, error) {
	if err := validateID(id, "project id"); err != nil {
		return nil, err
	}
	cfg := newProjectsOptions(opts)
	path := fmt.Sprintf("/projects/%d/groups", id)

	var payload struct {
		Data []ProjectGroup `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *ProjectsService) GetPlan(ctx context.Context, id ProjectID, opts ...ProjectsOption) ([]map[string]any, error) {
	if err := validateID(id, "project id"); err != nil {
		return nil, err
	}
	cfg := newProjectsOptions(opts)
	path := fmt.Sprintf("/projects/%d/plan", id)

	var payload struct {
		Data []map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing project plan data in response")
	}
	return payload.Data, nil
}

func (s *ProjectsService) UpdatePlanActivity(ctx context.Context, id ProjectID, activityID ProjectPlanActivityID, payload map[string]any, opts ...ProjectsOption) (map[string]any, error) {
	if err := validateID(id, "project id"); err != nil {
		return nil, err
	}
	if err := validateID(activityID, "project plan activity id"); err != nil {
		return nil, err
	}
	cfg := newProjectsOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("project plan activity payload is required")
	}
	path := fmt.Sprintf("/projects/%d/plan/activities/%d", id, activityID)

	var resp struct {
		Data map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing project plan activity data in response")
	}
	return resp.Data, nil
}

func (s *ProjectsService) UpdatePlanTask(ctx context.Context, id ProjectID, taskID ProjectPlanTaskID, payload map[string]any, opts ...ProjectsOption) (map[string]any, error) {
	if err := validateID(id, "project id"); err != nil {
		return nil, err
	}
	if err := validateID(taskID, "project plan task id"); err != nil {
		return nil, err
	}
	cfg := newProjectsOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("project plan task payload is required")
	}
	path := fmt.Sprintf("/projects/%d/plan/tasks/%d", id, taskID)

	var resp struct {
		Data map[string]any `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing project plan task data in response")
	}
	return resp.Data, nil
}

func (s *ProjectsService) ListTasks(ctx context.Context, id ProjectID, opts ...ProjectsOption) ([]ProjectTask, error) {
	if err := validateID(id, "project id"); err != nil {
		return nil, err
	}
	cfg := newProjectsOptions(opts)
	path := fmt.Sprintf("/projects/%d/tasks", id)

	var payload struct {
		Data []ProjectTask `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}
