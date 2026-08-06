package v2

import (
	"context"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ProjectTemplate struct {
	ID          ProjectTemplateID `json:"id"`
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	BoardID     ProjectBoardID    `json:"projects_board_id,omitempty"`
	OwnerID     UserID            `json:"owner_id,omitempty"`
	AddTime     *time.Time        `json:"add_time,omitempty"`
	UpdateTime  *time.Time        `json:"update_time,omitempty"`
}

type ProjectTemplatesService struct{ client *Client }

type ListProjectTemplatesOption interface {
	applyListProjectTemplates(*listProjectTemplatesOptions)
}
type ProjectTemplateRequestOption interface {
	ListProjectTemplatesOption
	projectTemplateRequestOptions() []pipedrive.RequestOption
}

type listProjectTemplatesOptions struct {
	params         genv2.GetProjectTemplatesParams
	requestOptions []pipedrive.RequestOption
}

type listProjectTemplatesOptionFunc func(*listProjectTemplatesOptions)

func (f listProjectTemplatesOptionFunc) applyListProjectTemplates(cfg *listProjectTemplatesOptions) {
	f(cfg)
}

type projectTemplateRequestOption struct{ options []pipedrive.RequestOption }

func (o projectTemplateRequestOption) projectTemplateRequestOptions() []pipedrive.RequestOption {
	return o.options
}
func (o projectTemplateRequestOption) applyListProjectTemplates(cfg *listProjectTemplatesOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}

func WithProjectTemplateRequestOptions(opts ...pipedrive.RequestOption) ProjectTemplateRequestOption {
	return projectTemplateRequestOption{options: opts}
}
func WithProjectTemplatesPageSize(limit int) ListProjectTemplatesOption {
	return listProjectTemplatesOptionFunc(func(cfg *listProjectTemplatesOptions) {
		if limit > 0 {
			cfg.params.Limit = &limit
		}
	})
}
func WithProjectTemplatesCursor(cursor string) ListProjectTemplatesOption {
	return listProjectTemplatesOptionFunc(func(cfg *listProjectTemplatesOptions) {
		if cursor != "" {
			cfg.params.Cursor = &cursor
		}
	})
}

func newListProjectTemplatesOptions(opts []ListProjectTemplatesOption) listProjectTemplatesOptions {
	var cfg listProjectTemplatesOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyListProjectTemplates(&cfg)
		}
	}
	return cfg
}
func projectTemplateRequestOptionValues(opts []ProjectTemplateRequestOption) []pipedrive.RequestOption {
	var out []pipedrive.RequestOption
	for _, opt := range opts {
		if opt != nil {
			out = append(out, opt.projectTemplateRequestOptions()...)
		}
	}
	return out
}

func (s *ProjectTemplatesService) List(ctx context.Context, opts ...ListProjectTemplatesOption) ([]ProjectTemplate, *string, error) {
	cfg := newListProjectTemplatesOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *ProjectTemplatesService) ListPager(opts ...ListProjectTemplatesOption) *pipedrive.CursorPager[ProjectTemplate] {
	cfg := newListProjectTemplatesOptions(opts)
	start := cfg.params.Cursor
	cfg.params.Cursor = nil
	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]ProjectTemplate, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if start != nil {
			params.Cursor = start
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *ProjectTemplatesService) ForEach(ctx context.Context, fn func(ProjectTemplate) error, opts ...ListProjectTemplatesOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *ProjectTemplatesService) Get(ctx context.Context, id ProjectTemplateID, opts ...ProjectTemplateRequestOption) (*ProjectTemplate, error) {
	if err := validateID(id, "project template id"); err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectTemplateRequestOptionValues(opts)...)
	resp, err := s.client.gen.GetProjectTemplateWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectTemplate](resp.HTTPResponse, resp.Body, "project template")
}

func (s *ProjectTemplatesService) list(ctx context.Context, params genv2.GetProjectTemplatesParams, requestOptions []pipedrive.RequestOption) ([]ProjectTemplate, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)
	resp, err := s.client.gen.GetProjectTemplatesWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	return decodeV2List[ProjectTemplate](resp.HTTPResponse, resp.Body)
}
