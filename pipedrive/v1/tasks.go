package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Task struct {
	ID               TaskID     `json:"id,omitempty"`
	Title            string     `json:"title,omitempty"`
	Description      string     `json:"description,omitempty"`
	DueDate          *string    `json:"due_date,omitempty"`
	Done             NumberBool `json:"done,omitempty"`
	AssigneeID       *UserID    `json:"assignee_id,omitempty"`
	CreatorID        *UserID    `json:"creator_id,omitempty"`
	ParentTaskID     *TaskID    `json:"parent_task_id,omitempty"`
	ProjectID        *ProjectID `json:"project_id,omitempty"`
	AddTime          *DateTime  `json:"add_time,omitempty"`
	UpdateTime       *DateTime  `json:"update_time,omitempty"`
	MarkedAsDoneTime *DateTime  `json:"marked_as_done_time,omitempty"`
}

type TasksService struct {
	client *Client
}

type TasksOption interface {
	applyTasks(*tasksOptions)
}

type tasksOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type tasksOptionFunc func(*tasksOptions)

func (f tasksOptionFunc) applyTasks(cfg *tasksOptions) {
	f(cfg)
}

func WithTasksQuery(values url.Values) TasksOption {
	return tasksOptionFunc(func(cfg *tasksOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithTasksRequestOptions(opts ...pipedrive.RequestOption) TasksOption {
	return tasksOptionFunc(func(cfg *tasksOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newTasksOptions(opts []TasksOption) tasksOptions {
	var cfg tasksOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyTasks(&cfg)
	}
	return cfg
}

func (s *TasksService) List(ctx context.Context, opts ...TasksOption) ([]Task, *CollectionPagination, error) {
	cfg := newTasksOptions(opts)

	var payload struct {
		Data           []Task                `json:"data"`
		AdditionalData *CollectionPagination `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/tasks", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *TasksService) Get(ctx context.Context, id TaskID, opts ...TasksOption) (*Task, error) {
	cfg := newTasksOptions(opts)
	path := fmt.Sprintf("/tasks/%d", id)

	var payload struct {
		Data *Task `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing task data in response")
	}
	return payload.Data, nil
}

func (s *TasksService) Create(ctx context.Context, payload map[string]any, opts ...TasksOption) (*Task, error) {
	cfg := newTasksOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("task payload is required")
	}

	var resp struct {
		Data *Task `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, "/tasks", cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing task data in response")
	}
	return resp.Data, nil
}

func (s *TasksService) Update(ctx context.Context, id TaskID, payload map[string]any, opts ...TasksOption) (*Task, error) {
	cfg := newTasksOptions(opts)
	if len(payload) == 0 {
		return nil, fmt.Errorf("task payload is required")
	}
	path := fmt.Sprintf("/tasks/%d", id)

	var resp struct {
		Data *Task `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, payload, &resp, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("missing task update data in response")
	}
	return resp.Data, nil
}

func (s *TasksService) Delete(ctx context.Context, id TaskID, opts ...TasksOption) (bool, error) {
	cfg := newTasksOptions(opts)
	path := fmt.Sprintf("/tasks/%d", id)

	var resp struct {
		Success *bool `json:"success"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, path, cfg.query, nil, &resp, cfg.requestOptions...); err != nil {
		return false, err
	}
	if resp.Success == nil {
		return false, fmt.Errorf("missing task delete success in response")
	}
	return *resp.Success, nil
}
