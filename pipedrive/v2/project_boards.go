package v2

import (
	"context"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ProjectBoard struct {
	ID         ProjectBoardID `json:"id"`
	Name       string         `json:"name,omitempty"`
	Order      int            `json:"order_nr,omitempty"`
	AddTime    *time.Time     `json:"add_time,omitempty"`
	UpdateTime *time.Time     `json:"update_time,omitempty"`
}

type ProjectBoardDeleteResult struct {
	ID ProjectBoardID `json:"id"`
}

type ProjectBoardsService struct{ client *Client }

type CreateProjectBoardOption interface {
	applyCreateProjectBoard(*createProjectBoardOptions)
}
type UpdateProjectBoardOption interface {
	applyUpdateProjectBoard(*updateProjectBoardOptions)
}

type ProjectBoardOption interface {
	CreateProjectBoardOption
	UpdateProjectBoardOption
}

type ProjectBoardRequestOption interface {
	CreateProjectBoardOption
	UpdateProjectBoardOption
	projectBoardRequestOptions() []pipedrive.RequestOption
}

type createProjectBoardOptions struct {
	payload        projectBoardPayload
	requestOptions []pipedrive.RequestOption
}

type updateProjectBoardOptions struct {
	payload        projectBoardPayload
	requestOptions []pipedrive.RequestOption
}

type projectBoardPayload struct {
	name  *string
	order *int
}

type projectBoardFieldOption func(*projectBoardPayload)

func (f projectBoardFieldOption) applyCreateProjectBoard(cfg *createProjectBoardOptions) {
	f(&cfg.payload)
}
func (f projectBoardFieldOption) applyUpdateProjectBoard(cfg *updateProjectBoardOptions) {
	f(&cfg.payload)
}

type projectBoardRequestOption struct{ options []pipedrive.RequestOption }

func (o projectBoardRequestOption) projectBoardRequestOptions() []pipedrive.RequestOption {
	return o.options
}
func (o projectBoardRequestOption) applyCreateProjectBoard(cfg *createProjectBoardOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o projectBoardRequestOption) applyUpdateProjectBoard(cfg *updateProjectBoardOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}

func WithProjectBoardRequestOptions(opts ...pipedrive.RequestOption) ProjectBoardRequestOption {
	return projectBoardRequestOption{options: opts}
}

func WithProjectBoardName(name string) ProjectBoardOption {
	return projectBoardFieldOption(func(payload *projectBoardPayload) { payload.name = &name })
}

func WithProjectBoardOrder(order int) ProjectBoardOption {
	return projectBoardFieldOption(func(payload *projectBoardPayload) { payload.order = &order })
}

func newCreateProjectBoardOptions(opts []CreateProjectBoardOption) createProjectBoardOptions {
	var cfg createProjectBoardOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyCreateProjectBoard(&cfg)
		}
	}
	return cfg
}

func newUpdateProjectBoardOptions(opts []UpdateProjectBoardOption) updateProjectBoardOptions {
	var cfg updateProjectBoardOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyUpdateProjectBoard(&cfg)
		}
	}
	return cfg
}

func projectBoardRequestOptionValues(opts []ProjectBoardRequestOption) []pipedrive.RequestOption {
	var out []pipedrive.RequestOption
	for _, opt := range opts {
		if opt != nil {
			out = append(out, opt.projectBoardRequestOptions()...)
		}
	}
	return out
}

func (p projectBoardPayload) body() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if p.order != nil {
		body["order_nr"] = *p.order
	}
	return body
}

func (s *ProjectBoardsService) List(ctx context.Context, opts ...ProjectBoardRequestOption) ([]ProjectBoard, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectBoardRequestOptionValues(opts)...)
	resp, err := s.client.gen.GetProjectsBoardsWithResponse(ctx, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2ListNoCursor[ProjectBoard](resp.HTTPResponse, resp.Body)
}

func (s *ProjectBoardsService) Get(ctx context.Context, id ProjectBoardID, opts ...ProjectBoardRequestOption) (*ProjectBoard, error) {
	if err := validateID(id, "project board id"); err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectBoardRequestOptionValues(opts)...)
	resp, err := s.client.gen.GetProjectsBoardWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectBoard](resp.HTTPResponse, resp.Body, "project board")
}

func (s *ProjectBoardsService) Create(ctx context.Context, opts ...CreateProjectBoardOption) (*ProjectBoard, error) {
	cfg := newCreateProjectBoardOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.AddProjectBoardWithBodyWithResponse(ctx, "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectBoard](resp.HTTPResponse, resp.Body, "project board")
}

func (s *ProjectBoardsService) Update(ctx context.Context, id ProjectBoardID, opts ...UpdateProjectBoardOption) (*ProjectBoard, error) {
	if err := validateID(id, "project board id"); err != nil {
		return nil, err
	}
	cfg := newUpdateProjectBoardOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.UpdateProjectBoardWithBodyWithResponse(ctx, int(id), "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectBoard](resp.HTTPResponse, resp.Body, "project board")
}

func (s *ProjectBoardsService) Delete(ctx context.Context, id ProjectBoardID, opts ...ProjectBoardRequestOption) (*ProjectBoardDeleteResult, error) {
	if err := validateID(id, "project board id"); err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectBoardRequestOptionValues(opts)...)
	resp, err := s.client.gen.DeleteProjectBoardWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectBoardDeleteResult](resp.HTTPResponse, resp.Body, "project board delete")
}
