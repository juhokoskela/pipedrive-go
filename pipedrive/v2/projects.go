package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ProjectStatus string

const (
	ProjectStatusOpen      ProjectStatus = "open"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusCanceled  ProjectStatus = "canceled"
	ProjectStatusDeleted   ProjectStatus = "deleted"
)

type ProjectSearchField string

const (
	ProjectSearchFieldCustomFields ProjectSearchField = "custom_fields"
	ProjectSearchFieldNotes        ProjectSearchField = "notes"
	ProjectSearchFieldTitle        ProjectSearchField = "title"
	ProjectSearchFieldDescription  ProjectSearchField = "description"
)

type Project struct {
	ID               ProjectID              `json:"id"`
	Title            string                 `json:"title,omitempty"`
	Description      string                 `json:"description,omitempty"`
	Status           ProjectStatus          `json:"status,omitempty"`
	BoardID          ProjectBoardID         `json:"board_id,omitempty"`
	PhaseID          ProjectPhaseID         `json:"phase_id,omitempty"`
	OwnerID          UserID                 `json:"owner_id,omitempty"`
	StartDate        string                 `json:"start_date,omitempty"`
	EndDate          string                 `json:"end_date,omitempty"`
	DealIDs          []DealID               `json:"deal_ids,omitempty"`
	PersonIDs        []PersonID             `json:"person_ids,omitempty"`
	OrganizationIDs  []OrganizationID       `json:"org_ids,omitempty"`
	LabelIDs         []int                  `json:"label_ids,omitempty"`
	HealthStatus     *int                   `json:"health_status,omitempty"`
	AddTime          *time.Time             `json:"add_time,omitempty"`
	UpdateTime       *time.Time             `json:"update_time,omitempty"`
	StatusChangeTime *time.Time             `json:"status_change_time,omitempty"`
	ArchiveTime      *time.Time             `json:"archive_time,omitempty"`
	CustomFields     map[string]interface{} `json:"custom_fields,omitempty"`
}

type ProjectSearchResult struct {
	ResultScore float64           `json:"result_score,omitempty"`
	Item        ProjectSearchItem `json:"item"`
}

type ProjectSearchItem struct {
	ID           ProjectID                  `json:"id"`
	Type         string                     `json:"type,omitempty"`
	Title        string                     `json:"title,omitempty"`
	Status       *ProjectStatus             `json:"status,omitempty"`
	Owner        *ProjectSearchOwner        `json:"owner,omitempty"`
	BoardID      *ProjectBoardID            `json:"board_id,omitempty"`
	Phase        *ProjectSearchPhase        `json:"phase,omitempty"`
	Person       *ProjectSearchPerson       `json:"person,omitempty"`
	Organization *ProjectSearchOrganization `json:"organization,omitempty"`
	Deal         *ProjectSearchDeal         `json:"deal,omitempty"`
	DealCount    int                        `json:"deal_count,omitempty"`
	Description  *string                    `json:"description,omitempty"`
	EndDate      *string                    `json:"end_date,omitempty"`
	CustomFields []string                   `json:"custom_fields,omitempty"`
	Notes        []string                   `json:"notes,omitempty"`
}

type ProjectSearchOwner struct {
	ID *UserID `json:"id,omitempty"`
}

type ProjectSearchPhase struct {
	ID   ProjectPhaseID `json:"id"`
	Name string         `json:"name,omitempty"`
}

type ProjectSearchPerson struct {
	ID   PersonID `json:"id"`
	Name *string  `json:"name,omitempty"`
}

type ProjectSearchOrganization struct {
	ID      OrganizationID `json:"id"`
	Name    *string        `json:"name,omitempty"`
	Address *string        `json:"address,omitempty"`
}

type ProjectSearchDeal struct {
	ID    DealID  `json:"id"`
	Title *string `json:"title,omitempty"`
}

type ProjectChangelogEntry struct {
	ChangeSource          *string                `json:"change_source,omitempty"`
	ChangeSourceUserAgent *string                `json:"change_source_user_agent,omitempty"`
	Time                  *time.Time             `json:"time,omitempty"`
	NewValues             map[string]interface{} `json:"new_values,omitempty"`
	OldValues             map[string]interface{} `json:"old_values,omitempty"`
	ActorUserID           UserID                 `json:"actor_user_id,omitempty"`
}

type ProjectDeleteResult struct {
	ID ProjectID `json:"id"`
}

type ProjectsService struct {
	client *Client
}

type ListProjectsOption interface{ applyListProjects(*listProjectsOptions) }
type ListArchivedProjectsOption interface {
	applyListArchivedProjects(*listArchivedProjectsOptions)
}

type ProjectListOption interface {
	ListProjectsOption
	ListArchivedProjectsOption
}

type SearchProjectsOption interface{ applySearchProjects(*searchProjectsOptions) }
type CreateProjectOption interface{ applyCreateProject(*createProjectOptions) }
type UpdateProjectOption interface{ applyUpdateProject(*updateProjectOptions) }
type ListProjectChangelogOption interface {
	applyListProjectChangelog(*listProjectChangelogOptions)
}

type ProjectRequestOption interface {
	ListProjectsOption
	ListArchivedProjectsOption
	SearchProjectsOption
	CreateProjectOption
	UpdateProjectOption
	ListProjectChangelogOption
	projectRequestOptions() []pipedrive.RequestOption
}

type ProjectOption interface {
	CreateProjectOption
	UpdateProjectOption
}

type listProjectsOptions struct {
	params         genv2.GetProjectsParams
	requestOptions []pipedrive.RequestOption
}

type listArchivedProjectsOptions struct {
	params         genv2.GetArchivedProjectsParams
	requestOptions []pipedrive.RequestOption
}

type searchProjectsOptions struct {
	params         genv2.SearchProjectsParams
	requestOptions []pipedrive.RequestOption
}

type listProjectChangelogOptions struct {
	params         genv2.GetProjectChangelogParams
	requestOptions []pipedrive.RequestOption
}

type createProjectOptions struct {
	payload        projectPayload
	requestOptions []pipedrive.RequestOption
}

type updateProjectOptions struct {
	payload        projectPayload
	requestOptions []pipedrive.RequestOption
}

type projectPayload struct {
	title           *string
	description     *string
	status          *ProjectStatus
	boardID         *ProjectBoardID
	phaseID         *ProjectPhaseID
	ownerID         *UserID
	startDate       *string
	endDate         *string
	healthStatus    *int
	templateID      *ProjectTemplateID
	dealIDs         []DealID
	dealIDsSet      bool
	personIDs       []PersonID
	personIDsSet    bool
	organizationIDs []OrganizationID
	orgIDsSet       bool
	labelIDs        []int
	labelIDsSet     bool
	customFields    map[string]interface{}
}

type projectRequestOption struct{ options []pipedrive.RequestOption }

func (o projectRequestOption) projectRequestOptions() []pipedrive.RequestOption { return o.options }
func (o projectRequestOption) applyListProjects(cfg *listProjectsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o projectRequestOption) applyListArchivedProjects(cfg *listArchivedProjectsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o projectRequestOption) applySearchProjects(cfg *searchProjectsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o projectRequestOption) applyCreateProject(cfg *createProjectOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o projectRequestOption) applyUpdateProject(cfg *updateProjectOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}
func (o projectRequestOption) applyListProjectChangelog(cfg *listProjectChangelogOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.options...)
}

type listProjectsOptionFunc func(*listProjectsOptions)

func (f listProjectsOptionFunc) applyListProjects(cfg *listProjectsOptions) { f(cfg) }

type listArchivedProjectsOptionFunc func(*listArchivedProjectsOptions)

func (f listArchivedProjectsOptionFunc) applyListArchivedProjects(cfg *listArchivedProjectsOptions) {
	f(cfg)
}

type searchProjectsOptionFunc func(*searchProjectsOptions)

func (f searchProjectsOptionFunc) applySearchProjects(cfg *searchProjectsOptions) { f(cfg) }

type listProjectChangelogOptionFunc func(*listProjectChangelogOptions)

func (f listProjectChangelogOptionFunc) applyListProjectChangelog(cfg *listProjectChangelogOptions) {
	f(cfg)
}

type projectFieldOption func(*projectPayload)

func (f projectFieldOption) applyCreateProject(cfg *createProjectOptions) { f(&cfg.payload) }
func (f projectFieldOption) applyUpdateProject(cfg *updateProjectOptions) { f(&cfg.payload) }

func WithProjectRequestOptions(opts ...pipedrive.RequestOption) ProjectRequestOption {
	return projectRequestOption{options: opts}
}

func WithProjectFilterID(id int) ProjectListOption {
	return projectListFilterOption{id: id}
}

type projectListFilterOption struct{ id int }

func (o projectListFilterOption) applyListProjects(cfg *listProjectsOptions) {
	cfg.params.FilterId = &o.id
}
func (o projectListFilterOption) applyListArchivedProjects(cfg *listArchivedProjectsOptions) {
	cfg.params.FilterId = &o.id
}

func WithProjectsStatus(statuses ...ProjectStatus) ProjectListOption {
	return projectListStatusOption{status: joinCSV(statuses)}
}

type projectListStatusOption struct{ status string }

func (o projectListStatusOption) applyListProjects(cfg *listProjectsOptions) {
	if o.status != "" {
		cfg.params.Status = &o.status
	}
}
func (o projectListStatusOption) applyListArchivedProjects(cfg *listArchivedProjectsOptions) {
	if o.status != "" {
		cfg.params.Status = &o.status
	}
}

func WithProjectsPhaseID(id ProjectPhaseID) ProjectListOption {
	return projectListPhaseOption{id: int(id)}
}

type projectListPhaseOption struct{ id int }

func (o projectListPhaseOption) applyListProjects(cfg *listProjectsOptions) {
	cfg.params.PhaseId = &o.id
}
func (o projectListPhaseOption) applyListArchivedProjects(cfg *listArchivedProjectsOptions) {
	cfg.params.PhaseId = &o.id
}

func WithProjectsDealID(id DealID) ListProjectsOption {
	return listProjectsOptionFunc(func(cfg *listProjectsOptions) { value := int(id); cfg.params.DealId = &value })
}

func WithProjectsPersonID(id PersonID) ListProjectsOption {
	return listProjectsOptionFunc(func(cfg *listProjectsOptions) { value := int(id); cfg.params.PersonId = &value })
}

func WithProjectsOrganizationID(id OrganizationID) ListProjectsOption {
	return listProjectsOptionFunc(func(cfg *listProjectsOptions) { value := int(id); cfg.params.OrgId = &value })
}

func WithProjectsPageSize(limit int) ProjectListOption {
	return projectListPageSizeOption{limit: limit}
}

type projectListPageSizeOption struct{ limit int }

func (o projectListPageSizeOption) applyListProjects(cfg *listProjectsOptions) {
	if o.limit > 0 {
		cfg.params.Limit = &o.limit
	}
}
func (o projectListPageSizeOption) applyListArchivedProjects(cfg *listArchivedProjectsOptions) {
	if o.limit > 0 {
		cfg.params.Limit = &o.limit
	}
}

func WithProjectsCursor(cursor string) ProjectListOption {
	return projectListCursorOption{cursor: cursor}
}

type projectListCursorOption struct{ cursor string }

func (o projectListCursorOption) applyListProjects(cfg *listProjectsOptions) {
	if o.cursor != "" {
		cfg.params.Cursor = &o.cursor
	}
}
func (o projectListCursorOption) applyListArchivedProjects(cfg *listArchivedProjectsOptions) {
	if o.cursor != "" {
		cfg.params.Cursor = &o.cursor
	}
}

func WithProjectSearchFields(fields ...ProjectSearchField) SearchProjectsOption {
	return searchProjectsOptionFunc(func(cfg *searchProjectsOptions) {
		value := joinCSV(fields)
		if value != "" {
			typed := genv2.SearchProjectsParamsFields(value)
			cfg.params.Fields = &typed
		}
	})
}

func WithProjectSearchExactMatch(exact bool) SearchProjectsOption {
	return searchProjectsOptionFunc(func(cfg *searchProjectsOptions) { cfg.params.ExactMatch = &exact })
}

func WithProjectSearchPersonID(id PersonID) SearchProjectsOption {
	return searchProjectsOptionFunc(func(cfg *searchProjectsOptions) { value := int(id); cfg.params.PersonId = &value })
}

func WithProjectSearchOrganizationID(id OrganizationID) SearchProjectsOption {
	return searchProjectsOptionFunc(func(cfg *searchProjectsOptions) { value := int(id); cfg.params.OrganizationId = &value })
}

func WithProjectSearchPageSize(limit int) SearchProjectsOption {
	return searchProjectsOptionFunc(func(cfg *searchProjectsOptions) {
		if limit > 0 {
			cfg.params.Limit = &limit
		}
	})
}

func WithProjectSearchCursor(cursor string) SearchProjectsOption {
	return searchProjectsOptionFunc(func(cfg *searchProjectsOptions) {
		if cursor != "" {
			cfg.params.Cursor = &cursor
		}
	})
}

func WithProjectChangelogPageSize(limit int) ListProjectChangelogOption {
	return listProjectChangelogOptionFunc(func(cfg *listProjectChangelogOptions) {
		if limit > 0 {
			cfg.params.Limit = &limit
		}
	})
}

func WithProjectChangelogCursor(cursor string) ListProjectChangelogOption {
	return listProjectChangelogOptionFunc(func(cfg *listProjectChangelogOptions) {
		if cursor != "" {
			cfg.params.Cursor = &cursor
		}
	})
}

func WithProjectTitle(title string) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.title = &title })
}
func WithProjectDescription(description string) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.description = &description })
}
func WithProjectStatus(status ProjectStatus) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.status = &status })
}
func WithProjectBoardID(id ProjectBoardID) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.boardID = &id })
}
func WithProjectPhaseID(id ProjectPhaseID) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.phaseID = &id })
}
func WithProjectOwnerID(id UserID) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.ownerID = &id })
}
func WithProjectStartDate(date string) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.startDate = &date })
}
func WithProjectEndDate(date string) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.endDate = &date })
}
func WithProjectHealthStatus(status int) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.healthStatus = &status })
}
func WithProjectTemplateID(id ProjectTemplateID) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.templateID = &id })
}
func WithProjectDealIDs(ids ...DealID) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.dealIDsSet = true; p.dealIDs = append([]DealID{}, ids...) })
}
func WithProjectPersonIDs(ids ...PersonID) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.personIDsSet = true; p.personIDs = append([]PersonID{}, ids...) })
}
func WithProjectOrganizationIDs(ids ...OrganizationID) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.orgIDsSet = true; p.organizationIDs = append([]OrganizationID{}, ids...) })
}
func WithProjectLabelIDs(ids ...int) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.labelIDsSet = true; p.labelIDs = append([]int{}, ids...) })
}
func WithProjectCustomFields(fields map[string]interface{}) ProjectOption {
	return projectFieldOption(func(p *projectPayload) { p.customFields = fields })
}

func newListProjectsOptions(opts []ListProjectsOption) listProjectsOptions {
	var cfg listProjectsOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyListProjects(&cfg)
		}
	}
	return cfg
}

func newListArchivedProjectsOptions(opts []ListArchivedProjectsOption) listArchivedProjectsOptions {
	var cfg listArchivedProjectsOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyListArchivedProjects(&cfg)
		}
	}
	return cfg
}

func newSearchProjectsOptions(term string, opts []SearchProjectsOption) searchProjectsOptions {
	cfg := searchProjectsOptions{params: genv2.SearchProjectsParams{Term: term}}
	for _, opt := range opts {
		if opt != nil {
			opt.applySearchProjects(&cfg)
		}
	}
	return cfg
}

func newCreateProjectOptions(opts []CreateProjectOption) createProjectOptions {
	var cfg createProjectOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyCreateProject(&cfg)
		}
	}
	return cfg
}

func newUpdateProjectOptions(opts []UpdateProjectOption) updateProjectOptions {
	var cfg updateProjectOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyUpdateProject(&cfg)
		}
	}
	return cfg
}

func newProjectChangelogOptions(opts []ListProjectChangelogOption) listProjectChangelogOptions {
	var cfg listProjectChangelogOptions
	for _, opt := range opts {
		if opt != nil {
			opt.applyListProjectChangelog(&cfg)
		}
	}
	return cfg
}

func projectRequestOptionValues(opts []ProjectRequestOption) []pipedrive.RequestOption {
	var out []pipedrive.RequestOption
	for _, opt := range opts {
		if opt != nil {
			out = append(out, opt.projectRequestOptions()...)
		}
	}
	return out
}

func (p projectPayload) body() map[string]interface{} {
	body := map[string]interface{}{}
	if p.title != nil {
		body["title"] = *p.title
	}
	if p.description != nil {
		body["description"] = *p.description
	}
	if p.status != nil {
		body["status"] = string(*p.status)
	}
	if p.boardID != nil {
		body["board_id"] = int(*p.boardID)
	}
	if p.phaseID != nil {
		body["phase_id"] = int(*p.phaseID)
	}
	if p.ownerID != nil {
		body["owner_id"] = int(*p.ownerID)
	}
	if p.startDate != nil {
		body["start_date"] = *p.startDate
	}
	if p.endDate != nil {
		body["end_date"] = *p.endDate
	}
	if p.healthStatus != nil {
		body["health_status"] = *p.healthStatus
	}
	if p.templateID != nil {
		body["template_id"] = int(*p.templateID)
	}
	if p.dealIDsSet {
		body["deal_ids"] = p.dealIDs
	}
	if p.personIDsSet {
		body["person_ids"] = p.personIDs
	}
	if p.orgIDsSet {
		body["org_ids"] = p.organizationIDs
	}
	if p.labelIDsSet {
		body["label_ids"] = p.labelIDs
	}
	if p.customFields != nil {
		body["custom_fields"] = p.customFields
	}
	return body
}

func (s *ProjectsService) List(ctx context.Context, opts ...ListProjectsOption) ([]Project, *string, error) {
	cfg := newListProjectsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *ProjectsService) ListPager(opts ...ListProjectsOption) *pipedrive.CursorPager[Project] {
	cfg := newListProjectsOptions(opts)
	start := cfg.params.Cursor
	cfg.params.Cursor = nil
	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Project, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if start != nil {
			params.Cursor = start
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *ProjectsService) ForEach(ctx context.Context, fn func(Project) error, opts ...ListProjectsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *ProjectsService) ListArchived(ctx context.Context, opts ...ListArchivedProjectsOption) ([]Project, *string, error) {
	cfg := newListArchivedProjectsOptions(opts)
	return s.listArchived(ctx, cfg.params, cfg.requestOptions)
}

func (s *ProjectsService) ListArchivedPager(opts ...ListArchivedProjectsOption) *pipedrive.CursorPager[Project] {
	cfg := newListArchivedProjectsOptions(opts)
	start := cfg.params.Cursor
	cfg.params.Cursor = nil
	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Project, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if start != nil {
			params.Cursor = start
		}
		return s.listArchived(ctx, params, cfg.requestOptions)
	})
}

func (s *ProjectsService) ForEachArchived(ctx context.Context, fn func(Project) error, opts ...ListArchivedProjectsOption) error {
	return s.ListArchivedPager(opts...).ForEach(ctx, fn)
}

func (s *ProjectsService) Search(ctx context.Context, term string, opts ...SearchProjectsOption) ([]ProjectSearchResult, *string, error) {
	cfg := newSearchProjectsOptions(term, opts)
	return s.search(ctx, cfg.params, cfg.requestOptions)
}

func (s *ProjectsService) SearchPager(term string, opts ...SearchProjectsOption) *pipedrive.CursorPager[ProjectSearchResult] {
	cfg := newSearchProjectsOptions(term, opts)
	start := cfg.params.Cursor
	cfg.params.Cursor = nil
	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]ProjectSearchResult, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if start != nil {
			params.Cursor = start
		}
		return s.search(ctx, params, cfg.requestOptions)
	})
}

func (s *ProjectsService) ForEachSearch(ctx context.Context, term string, fn func(ProjectSearchResult) error, opts ...SearchProjectsOption) error {
	return s.SearchPager(term, opts...).ForEach(ctx, fn)
}

func (s *ProjectsService) Get(ctx context.Context, id ProjectID, opts ...ProjectRequestOption) (*Project, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectRequestOptionValues(opts)...)
	resp, err := s.client.gen.GetProjectWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Project](resp.HTTPResponse, resp.Body, "project")
}

func (s *ProjectsService) Create(ctx context.Context, opts ...CreateProjectOption) (*Project, error) {
	cfg := newCreateProjectOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.AddProjectWithBodyWithResponse(ctx, "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Project](resp.HTTPResponse, resp.Body, "project")
}

func (s *ProjectsService) Update(ctx context.Context, id ProjectID, opts ...UpdateProjectOption) (*Project, error) {
	cfg := newUpdateProjectOptions(opts)
	body, err := encodeV2Body(cfg.payload.body())
	if err != nil {
		return nil, err
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)
	resp, err := s.client.gen.UpdateProjectWithBodyWithResponse(ctx, int(id), "application/json", body, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Project](resp.HTTPResponse, resp.Body, "project")
}

func (s *ProjectsService) Delete(ctx context.Context, id ProjectID, opts ...ProjectRequestOption) (*ProjectDeleteResult, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectRequestOptionValues(opts)...)
	resp, err := s.client.gen.DeleteProjectWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[ProjectDeleteResult](resp.HTTPResponse, resp.Body, "project delete")
}

func (s *ProjectsService) Archive(ctx context.Context, id ProjectID, opts ...ProjectRequestOption) (*Project, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectRequestOptionValues(opts)...)
	resp, err := s.client.gen.ArchiveProjectWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2Data[Project](resp.HTTPResponse, resp.Body, "archived project")
}

func (s *ProjectsService) ListChangelog(ctx context.Context, id ProjectID, opts ...ListProjectChangelogOption) ([]ProjectChangelogEntry, *string, error) {
	cfg := newProjectChangelogOptions(opts)
	return s.listChangelog(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *ProjectsService) ChangelogPager(id ProjectID, opts ...ListProjectChangelogOption) *pipedrive.CursorPager[ProjectChangelogEntry] {
	cfg := newProjectChangelogOptions(opts)
	start := cfg.params.Cursor
	cfg.params.Cursor = nil
	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]ProjectChangelogEntry, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if start != nil {
			params.Cursor = start
		}
		return s.listChangelog(ctx, id, params, cfg.requestOptions)
	})
}

func (s *ProjectsService) ForEachChangelog(ctx context.Context, id ProjectID, fn func(ProjectChangelogEntry) error, opts ...ListProjectChangelogOption) error {
	return s.ChangelogPager(id, opts...).ForEach(ctx, fn)
}

func (s *ProjectsService) ListPermittedUsers(ctx context.Context, id ProjectID, opts ...ProjectRequestOption) ([]UserID, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, projectRequestOptionValues(opts)...)
	resp, err := s.client.gen.GetProjectUsersWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	return decodeV2ListNoCursor[UserID](resp.HTTPResponse, resp.Body)
}

func (s *ProjectsService) list(ctx context.Context, params genv2.GetProjectsParams, requestOptions []pipedrive.RequestOption) ([]Project, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)
	resp, err := s.client.gen.GetProjectsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	return decodeV2List[Project](resp.HTTPResponse, resp.Body)
}

func (s *ProjectsService) listArchived(ctx context.Context, params genv2.GetArchivedProjectsParams, requestOptions []pipedrive.RequestOption) ([]Project, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)
	resp, err := s.client.gen.GetArchivedProjectsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	return decodeV2List[Project](resp.HTTPResponse, resp.Body)
}

func (s *ProjectsService) search(ctx context.Context, params genv2.SearchProjectsParams, requestOptions []pipedrive.RequestOption) ([]ProjectSearchResult, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)
	resp, err := s.client.gen.SearchProjectsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}
	var payload struct {
		Data struct {
			Items []ProjectSearchResult `json:"items"`
		} `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data.Items, next, nil
}

func (s *ProjectsService) listChangelog(ctx context.Context, id ProjectID, params genv2.GetProjectChangelogParams, requestOptions []pipedrive.RequestOption) ([]ProjectChangelogEntry, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)
	resp, err := s.client.gen.GetProjectChangelogWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	return decodeV2List[ProjectChangelogEntry](resp.HTTPResponse, resp.Body)
}

func encodeV2Body(body map[string]interface{}) (*bytes.Reader, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}
	return bytes.NewReader(encoded), nil
}

func encodeV2ArrayBody(body []map[string]interface{}) (*bytes.Reader, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}
	return bytes.NewReader(encoded), nil
}

func decodeV2Data[T any](httpResp *http.Response, body []byte, label string) (*T, error) {
	if httpResp.StatusCode < 200 || httpResp.StatusCode > 299 {
		return nil, errorFromResponse(httpResp, body)
	}
	var payload struct {
		Data *T `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing %s data in response", label)
	}
	return payload.Data, nil
}

func decodeV2List[T any](httpResp *http.Response, body []byte) ([]T, *string, error) {
	if httpResp.StatusCode < 200 || httpResp.StatusCode > 299 {
		return nil, nil, errorFromResponse(httpResp, body)
	}
	var payload struct {
		Data           []T `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func decodeV2ListNoCursor[T any](httpResp *http.Response, body []byte) ([]T, error) {
	items, _, err := decodeV2List[T](httpResp, body)
	return items, err
}
