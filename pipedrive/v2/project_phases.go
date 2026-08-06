package v2

import (
	"context"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ProjectPhase struct {
	ID         ProjectPhaseID `json:"id"`
	Name       string         `json:"name,omitempty"`
	BoardID    ProjectBoardID `json:"board_id,omitempty"`
	Order      int            `json:"order_nr,omitempty"`
	AddTime    *time.Time     `json:"add_time,omitempty"`
	UpdateTime *time.Time     `json:"update_time,omitempty"`
}

type ProjectPhaseDeleteResult struct {
	ID ProjectPhaseID `json:"id"`
}

type ProjectPhasesService struct{ client *Client }

type CreateProjectPhaseOption interface {
	applyCreateProjectPhase(*createProjectPhaseOptions)
}
type UpdateProjectPhaseOption interface {
	applyUpdateProjectPhase(*updateProjectPhaseOptions)
}

type ProjectPhaseOption interface {
	CreateProjectPhaseOption
	UpdateProjectPhaseOption
}

type ProjectPhaseRequestOption interface {
	CreateProjectPhaseOption
	UpdateProjectPhaseOption
	projectPhaseRequestOptions() []pipedrive.RequestOption
}

type createProjectPhaseOptions struct {
	payload        projectPhasePayload
	requestOptions []pipedrive.RequestOption
}
type updateProjectPhaseOptions struct {
	payload        projectPhasePayload
	requestOptions []pipedrive.RequestOption
}
type projectPhasePayload struct {
	name    *string
	boardID *ProjectBoardID
	order   *int
}

type projectPhaseFieldOption func(*projectPhasePayload)

func (f projectPhaseFieldOption) applyCreateProjectPhase(cfg *createProjectPhaseOptions) {
	f(&cfg.payload)
}
func (f projectPhaseFieldOption) applyUpdateProjectPhase(cfg *updateProjectPhaseOptions) {
	f(&cfg.payload)
}

type projectPhaseRequestOption struct{ options []pipedrive.RequestOption }

func (o projectPhaseRequestOption) projectPhaseRequestOptions() []pipedrive.RequestOption {
	return o.options
}
func (o projectPhaseRequestOption) applyCreateProjectPhase(cfg *createProjectPhaseOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o projectPhaseRequestOption) applyUpdateProjectPhase(cfg *updateProjectPhaseOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}

func WithProjectPhaseRequestOptions(opts ...pipedrive.RequestOption) ProjectPhaseRequestOption {
	return projectPhaseRequestOption{options: opts}
}
func WithProjectPhaseName(name string) ProjectPhaseOption {
	return projectPhaseFieldOption(func(p *projectPhasePayload) { p.name = &name })
}
func WithProjectPhaseBoardID(id ProjectBoardID) ProjectPhaseOption {
	return projectPhaseFieldOption(func(p *projectPhasePayload) { p.boardID = &id })
}
func WithProjectPhaseOrder(order int) ProjectPhaseOption {
	return projectPhaseFieldOption(func(p *projectPhasePayload) { p.order = &order })
}

func newCreateProjectPhaseOptions(opts []CreateProjectPhaseOption) createProjectPhaseOptions {
	var cfg createProjectPhaseOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyCreateProjectPhase(&cfg)
		}
	}
	return cfg
}
func newUpdateProjectPhaseOptions(opts []UpdateProjectPhaseOption) updateProjectPhaseOptions {
	var cfg updateProjectPhaseOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyUpdateProjectPhase(&cfg)
		}
	}
	return cfg
}
func projectPhaseRequestOptionValues(opts []ProjectPhaseRequestOption) []pipedrive.RequestOption {
	var out []pipedrive.RequestOption
	for _, opt := range opts {
		if opt != nil {
			out = append(out, opt.projectPhaseRequestOptions()...)
		}
	}
	return out
}
func (p projectPhasePayload) body() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if p.boardID != nil {
		body["board_id"] = int(*p.boardID)
	}
	if p.order != nil {
		body["order_nr"] = *p.order
	}
	return body
}

func (s *ProjectPhasesService) List(ctx context.Context, boardID ProjectBoardID, opts ...ProjectPhaseRequestOption) ([]ProjectPhase, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectPhaseRequestOptionValues(opts)...)
	params := genv2.GetProjectsPhasesParams{BoardId: int(boardID)}
	resp, err := s.client.gen.GetProjectsPhasesWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2ListNoCursor[ProjectPhase](resp.HTTPResponse, resp.Body)
}

func (s *ProjectPhasesService) Get(ctx context.Context, id ProjectPhaseID, opts ...ProjectPhaseRequestOption) (*ProjectPhase, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectPhaseRequestOptionValues(opts)...)
	resp, err := s.client.gen.GetProjectsPhaseWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectPhase](resp.HTTPResponse, resp.Body, "project phase")
}

func (s *ProjectPhasesService) Create(ctx context.Context, opts ...CreateProjectPhaseOption) (*ProjectPhase, error) {
	cfg := newCreateProjectPhaseOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.AddProjectPhaseWithBodyWithResponse(ctx, "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectPhase](resp.HTTPResponse, resp.Body, "project phase")
}

func (s *ProjectPhasesService) Update(ctx context.Context, id ProjectPhaseID, opts ...UpdateProjectPhaseOption) (*ProjectPhase, error) {
	cfg := newUpdateProjectPhaseOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.UpdateProjectPhaseWithBodyWithResponse(ctx, int(id), "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectPhase](resp.HTTPResponse, resp.Body, "project phase")
}

func (s *ProjectPhasesService) Delete(ctx context.Context, id ProjectPhaseID, opts ...ProjectPhaseRequestOption) (*ProjectPhaseDeleteResult, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectPhaseRequestOptionValues(opts)...)
	resp, err := s.client.gen.DeleteProjectPhaseWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectPhaseDeleteResult](resp.HTTPResponse, resp.Body, "project phase delete")
}
