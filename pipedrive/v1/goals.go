package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type GoalInterval string

const (
	GoalIntervalWeekly    GoalInterval = "weekly"
	GoalIntervalMonthly   GoalInterval = "monthly"
	GoalIntervalQuarterly GoalInterval = "quarterly"
	GoalIntervalYearly    GoalInterval = "yearly"
)

type GoalTypeName string

const (
	GoalTypeNameDealsWon            GoalTypeName = "deals_won"
	GoalTypeNameDealsProgressed     GoalTypeName = "deals_progressed"
	GoalTypeNameActivitiesCompleted GoalTypeName = "activities_completed"
	GoalTypeNameActivitiesAdded     GoalTypeName = "activities_added"
	GoalTypeNameDealsStarted        GoalTypeName = "deals_started"
	GoalTypeNameRevenueForecast     GoalTypeName = "revenue_forecast"
)

type GoalAssigneeType string

const (
	GoalAssigneeTypePerson  GoalAssigneeType = "person"
	GoalAssigneeTypeCompany GoalAssigneeType = "company"
	GoalAssigneeTypeTeam    GoalAssigneeType = "team"
)

type GoalTrackingMetric string

const (
	GoalTrackingMetricQuantity GoalTrackingMetric = "quantity"
	GoalTrackingMetricSum      GoalTrackingMetric = "sum"
)

type GoalTypeParams struct {
	PipelineIDs     []PipelineID     `json:"pipeline_id,omitempty"`
	StageID         *StageID         `json:"stage_id,omitempty"`
	ActivityTypeIDs []ActivityTypeID `json:"activity_type_id,omitempty"`
}

type GoalType struct {
	Name   GoalTypeName    `json:"name,omitempty"`
	Params *GoalTypeParams `json:"params,omitempty"`
}

type GoalAssignee struct {
	ID   int              `json:"id,omitempty"`
	Type GoalAssigneeType `json:"type,omitempty"`
}

type GoalDuration struct {
	Start string  `json:"start,omitempty"`
	End   *string `json:"end,omitempty"`
}

type GoalExpectedOutcome struct {
	Target         float64            `json:"target,omitempty"`
	TrackingMetric GoalTrackingMetric `json:"tracking_metric,omitempty"`
	CurrencyID     *CurrencyID        `json:"currency_id,omitempty"`
}

type Goal struct {
	ID              GoalID               `json:"id,omitempty"`
	OwnerID         *UserID              `json:"owner_id,omitempty"`
	Title           string               `json:"title,omitempty"`
	Type            *GoalType            `json:"type,omitempty"`
	Assignee        *GoalAssignee        `json:"assignee,omitempty"`
	Interval        GoalInterval         `json:"interval,omitempty"`
	Duration        *GoalDuration        `json:"duration,omitempty"`
	ExpectedOutcome *GoalExpectedOutcome `json:"expected_outcome,omitempty"`
	IsActive        bool                 `json:"is_active,omitempty"`
	ReportIDs       []string             `json:"report_ids,omitempty"`
}

type GoalResult struct {
	Progress int   `json:"progress,omitempty"`
	Goal     *Goal `json:"goal,omitempty"`
}

type GoalsService struct {
	client *Client
}

type ListGoalsOption interface {
	applyListGoals(*listGoalsOptions)
}

type CreateGoalOption interface {
	applyCreateGoal(*createGoalOptions)
}

type UpdateGoalOption interface {
	applyUpdateGoal(*updateGoalOptions)
}

type DeleteGoalOption interface {
	applyDeleteGoal(*deleteGoalOptions)
}

type GetGoalResultOption interface {
	applyGetGoalResult(*getGoalResultOptions)
}

type GoalsRequestOption interface {
	ListGoalsOption
	CreateGoalOption
	UpdateGoalOption
	DeleteGoalOption
	GetGoalResultOption
}

type GoalOption interface {
	CreateGoalOption
	UpdateGoalOption
}

type listGoalsOptions struct {
	params         genv1.GetGoalsParams
	requestOptions []pipedrive.RequestOption
}

type createGoalOptions struct {
	payload        goalPayload
	requestOptions []pipedrive.RequestOption
}

type updateGoalOptions struct {
	payload        goalPayload
	requestOptions []pipedrive.RequestOption
}

type deleteGoalOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getGoalResultOptions struct {
	params         genv1.GetGoalResultParams
	requestOptions []pipedrive.RequestOption
}

type goalPayload struct {
	title           *string
	assignee        *goalAssigneePayload
	goalType        *goalTypePayload
	expectedOutcome *goalExpectedOutcomePayload
	duration        *goalDurationPayload
	interval        *GoalInterval
}

type goalAssigneePayload struct {
	id           *int
	assigneeType *GoalAssigneeType
}

type goalTypePayload struct {
	name   *GoalTypeName
	params *goalTypeParamsPayload
}

type goalTypeParamsPayload struct {
	pipelineIDs     *[]PipelineID
	stageID         *StageID
	activityTypeIDs *[]ActivityTypeID
}

type goalExpectedOutcomePayload struct {
	target         *float64
	trackingMetric *GoalTrackingMetric
	currencyID     *CurrencyID
}

type goalDurationPayload struct {
	start *string
	end   *string
}

type goalsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o goalsRequestOptions) applyListGoals(cfg *listGoalsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o goalsRequestOptions) applyCreateGoal(cfg *createGoalOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o goalsRequestOptions) applyUpdateGoal(cfg *updateGoalOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o goalsRequestOptions) applyDeleteGoal(cfg *deleteGoalOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o goalsRequestOptions) applyGetGoalResult(cfg *getGoalResultOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listGoalsOptionFunc func(*listGoalsOptions)

func (f listGoalsOptionFunc) applyListGoals(cfg *listGoalsOptions) {
	f(cfg)
}

type goalFieldOption func(*goalPayload)

func (f goalFieldOption) applyCreateGoal(cfg *createGoalOptions) {
	f(&cfg.payload)
}

func (f goalFieldOption) applyUpdateGoal(cfg *updateGoalOptions) {
	f(&cfg.payload)
}

type deleteGoalOptionFunc func(*deleteGoalOptions)

func (f deleteGoalOptionFunc) applyDeleteGoal(cfg *deleteGoalOptions) {
	f(cfg)
}

type getGoalResultOptionFunc func(*getGoalResultOptions)

func (f getGoalResultOptionFunc) applyGetGoalResult(cfg *getGoalResultOptions) {
	f(cfg)
}

func WithGoalsRequestOptions(opts ...pipedrive.RequestOption) GoalsRequestOption {
	return goalsRequestOptions{requestOptions: opts}
}

func WithGoalsTypeName(name GoalTypeName) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		value := genv1.GetGoalsParamsTypeName(name)
		cfg.params.TypeName = &value
	})
}

func WithGoalsTitle(title string) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		cfg.params.Title = &title
	})
}

func WithGoalsActive(active bool) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		cfg.params.IsActive = &active
	})
}

func WithGoalsAssigneeID(id int) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		cfg.params.AssigneeId = &id
	})
}

func WithGoalsAssigneeType(assigneeType GoalAssigneeType) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		value := genv1.GetGoalsParamsAssigneeType(assigneeType)
		cfg.params.AssigneeType = &value
	})
}

func WithGoalsExpectedOutcomeTarget(target float64) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		value := float32(target)
		cfg.params.ExpectedOutcomeTarget = &value
	})
}

func WithGoalsExpectedOutcomeTrackingMetric(metric GoalTrackingMetric) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		value := genv1.GetGoalsParamsExpectedOutcomeTrackingMetric(metric)
		cfg.params.ExpectedOutcomeTrackingMetric = &value
	})
}

func WithGoalsExpectedOutcomeCurrencyID(id CurrencyID) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		value := int(id)
		cfg.params.ExpectedOutcomeCurrencyId = &value
	})
}

func WithGoalsTypePipelineIDs(ids ...PipelineID) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		if len(ids) == 0 {
			return
		}
		values := make([]int, 0, len(ids))
		for _, id := range ids {
			values = append(values, int(id))
		}
		cfg.params.TypeParamsPipelineId = &values
	})
}

func WithGoalsTypeStageID(id StageID) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		value := int(id)
		cfg.params.TypeParamsStageId = &value
	})
}

func WithGoalsTypeActivityTypeIDs(ids ...ActivityTypeID) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		if len(ids) == 0 {
			return
		}
		values := make([]int, 0, len(ids))
		for _, id := range ids {
			values = append(values, int(id))
		}
		cfg.params.TypeParamsActivityTypeId = &values
	})
}

func WithGoalsPeriodStart(start time.Time) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		cfg.params.PeriodStart = &openapi_types.Date{Time: start}
	})
}

func WithGoalsPeriodEnd(end time.Time) ListGoalsOption {
	return listGoalsOptionFunc(func(cfg *listGoalsOptions) {
		cfg.params.PeriodEnd = &openapi_types.Date{Time: end}
	})
}

func WithGoalTitle(title string) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		cfg.title = &title
	})
}

func WithGoalAssigneeID(id int) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		assignee := cfg.ensureAssignee()
		assignee.id = &id
	})
}

func WithGoalAssigneeType(assigneeType GoalAssigneeType) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		assignee := cfg.ensureAssignee()
		assignee.assigneeType = &assigneeType
	})
}

func WithGoalTypeName(name GoalTypeName) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		goalType := cfg.ensureGoalType()
		goalType.name = &name
	})
}

func WithGoalTypePipelineIDs(ids ...PipelineID) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		params := cfg.ensureGoalTypeParams()
		clone := append([]PipelineID{}, ids...)
		params.pipelineIDs = &clone
	})
}

func WithGoalTypeStageID(id StageID) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		params := cfg.ensureGoalTypeParams()
		params.stageID = &id
	})
}

func WithGoalTypeActivityTypeIDs(ids ...ActivityTypeID) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		params := cfg.ensureGoalTypeParams()
		clone := append([]ActivityTypeID{}, ids...)
		params.activityTypeIDs = &clone
	})
}

func WithGoalExpectedOutcomeTarget(target float64) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		outcome := cfg.ensureExpectedOutcome()
		outcome.target = &target
	})
}

func WithGoalExpectedOutcomeTrackingMetric(metric GoalTrackingMetric) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		outcome := cfg.ensureExpectedOutcome()
		outcome.trackingMetric = &metric
	})
}

func WithGoalExpectedOutcomeCurrencyID(id CurrencyID) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		outcome := cfg.ensureExpectedOutcome()
		outcome.currencyID = &id
	})
}

func WithGoalDurationStart(start string) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		duration := cfg.ensureDuration()
		duration.start = &start
	})
}

func WithGoalDurationEnd(end string) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		duration := cfg.ensureDuration()
		duration.end = &end
	})
}

func WithGoalInterval(interval GoalInterval) GoalOption {
	return goalFieldOption(func(cfg *goalPayload) {
		cfg.interval = &interval
	})
}

func WithGoalResultStartDate(start time.Time) GetGoalResultOption {
	return getGoalResultOptionFunc(func(cfg *getGoalResultOptions) {
		cfg.params.PeriodStart = openapi_types.Date{Time: start}
	})
}

func WithGoalResultEndDate(end time.Time) GetGoalResultOption {
	return getGoalResultOptionFunc(func(cfg *getGoalResultOptions) {
		cfg.params.PeriodEnd = openapi_types.Date{Time: end}
	})
}

func newListGoalsOptions(opts []ListGoalsOption) listGoalsOptions {
	var cfg listGoalsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListGoals(&cfg)
	}
	return cfg
}

func newCreateGoalOptions(opts []CreateGoalOption) createGoalOptions {
	var cfg createGoalOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateGoal(&cfg)
	}
	return cfg
}

func newUpdateGoalOptions(opts []UpdateGoalOption) updateGoalOptions {
	var cfg updateGoalOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateGoal(&cfg)
	}
	return cfg
}

func newDeleteGoalOptions(opts []DeleteGoalOption) deleteGoalOptions {
	var cfg deleteGoalOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteGoal(&cfg)
	}
	return cfg
}

func newGetGoalResultOptions(opts []GetGoalResultOption) getGoalResultOptions {
	var cfg getGoalResultOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetGoalResult(&cfg)
	}
	return cfg
}

func (p *goalPayload) ensureAssignee() *goalAssigneePayload {
	if p.assignee == nil {
		p.assignee = &goalAssigneePayload{}
	}
	return p.assignee
}

func (p *goalPayload) ensureGoalType() *goalTypePayload {
	if p.goalType == nil {
		p.goalType = &goalTypePayload{}
	}
	return p.goalType
}

func (p *goalPayload) ensureGoalTypeParams() *goalTypeParamsPayload {
	goalType := p.ensureGoalType()
	if goalType.params == nil {
		goalType.params = &goalTypeParamsPayload{}
	}
	return goalType.params
}

func (p *goalPayload) ensureExpectedOutcome() *goalExpectedOutcomePayload {
	if p.expectedOutcome == nil {
		p.expectedOutcome = &goalExpectedOutcomePayload{}
	}
	return p.expectedOutcome
}

func (p *goalPayload) ensureDuration() *goalDurationPayload {
	if p.duration == nil {
		p.duration = &goalDurationPayload{}
	}
	return p.duration
}

func (p goalAssigneePayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.id != nil {
		body["id"] = *p.id
	}
	if p.assigneeType != nil {
		body["type"] = string(*p.assigneeType)
	}
	return body
}

func (p goalTypeParamsPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.pipelineIDs != nil {
		body["pipeline_id"] = *p.pipelineIDs
	}
	if p.stageID != nil {
		body["stage_id"] = *p.stageID
	}
	if p.activityTypeIDs != nil {
		body["activity_type_id"] = *p.activityTypeIDs
	}
	return body
}

func (p goalTypePayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = string(*p.name)
	}
	if p.params != nil {
		params := p.params.toMap()
		if len(params) > 0 {
			body["params"] = params
		}
	}
	return body
}

func (p goalExpectedOutcomePayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.target != nil {
		body["target"] = *p.target
	}
	if p.trackingMetric != nil {
		body["tracking_metric"] = string(*p.trackingMetric)
	}
	if p.currencyID != nil {
		body["currency_id"] = *p.currencyID
	}
	return body
}

func (p goalDurationPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.start != nil {
		body["start"] = *p.start
	}
	if p.end != nil {
		body["end"] = *p.end
	}
	return body
}

func (p goalPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.title != nil {
		body["title"] = *p.title
	}
	if p.assignee != nil {
		assignee := p.assignee.toMap()
		if len(assignee) > 0 {
			body["assignee"] = assignee
		}
	}
	if p.goalType != nil {
		goalType := p.goalType.toMap()
		if len(goalType) > 0 {
			body["type"] = goalType
		}
	}
	if p.expectedOutcome != nil {
		outcome := p.expectedOutcome.toMap()
		if len(outcome) > 0 {
			body["expected_outcome"] = outcome
		}
	}
	if p.duration != nil {
		duration := p.duration.toMap()
		if len(duration) > 0 {
			body["duration"] = duration
		}
	}
	if p.interval != nil {
		body["interval"] = string(*p.interval)
	}
	return body
}

func (s *GoalsService) List(ctx context.Context, opts ...ListGoalsOption) ([]Goal, error) {
	cfg := newListGoalsOptions(opts)
	if (cfg.params.PeriodStart == nil) != (cfg.params.PeriodEnd == nil) {
		return nil, fmt.Errorf("period start and end must be provided together")
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetGoals(ctx, &cfg.params, toRequestEditors(editors)...)
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
		Data struct {
			Goals []Goal `json:"goals"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data.Goals, nil
}

func (s *GoalsService) Create(ctx context.Context, opts ...CreateGoalOption) (*Goal, error) {
	cfg := newCreateGoalOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.assignee == nil || cfg.payload.assignee.id == nil || cfg.payload.assignee.assigneeType == nil {
		return nil, fmt.Errorf("assignee id and type are required")
	}
	if cfg.payload.goalType == nil || cfg.payload.goalType.name == nil {
		return nil, fmt.Errorf("goal type name is required")
	}
	if cfg.payload.expectedOutcome == nil || cfg.payload.expectedOutcome.target == nil || cfg.payload.expectedOutcome.trackingMetric == nil {
		return nil, fmt.Errorf("expected outcome target and tracking metric are required")
	}
	if cfg.payload.duration == nil || cfg.payload.duration.start == nil {
		return nil, fmt.Errorf("duration start is required")
	}
	if cfg.payload.interval == nil {
		return nil, fmt.Errorf("interval is required")
	}

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddGoalWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *struct {
			Goal *Goal `json:"goal"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil || payload.Data.Goal == nil {
		return nil, fmt.Errorf("missing goal data in response")
	}
	return payload.Data.Goal, nil
}

func (s *GoalsService) Update(ctx context.Context, id GoalID, opts ...UpdateGoalOption) (*Goal, error) {
	cfg := newUpdateGoalOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	bodyMap := cfg.payload.toMap()
	if len(bodyMap) == 0 {
		return nil, fmt.Errorf("at least one field is required to update")
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateGoalWithBody(ctx, string(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *struct {
			Goal *Goal `json:"goal"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil || payload.Data.Goal == nil {
		return nil, fmt.Errorf("missing goal data in response")
	}
	return payload.Data.Goal, nil
}

func (s *GoalsService) Delete(ctx context.Context, id GoalID, opts ...DeleteGoalOption) (bool, error) {
	cfg := newDeleteGoalOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteGoal(ctx, string(id), toRequestEditors(editors)...)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false, errorFromResponse(resp, respBody)
	}

	var payload struct {
		Success *bool `json:"success"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}
	if payload.Success == nil {
		return false, fmt.Errorf("missing goal delete success in response")
	}
	return *payload.Success, nil
}

func (s *GoalsService) GetResult(ctx context.Context, id GoalID, opts ...GetGoalResultOption) (*GoalResult, error) {
	cfg := newGetGoalResultOptions(opts)
	if cfg.params.PeriodStart.IsZero() || cfg.params.PeriodEnd.IsZero() {
		return nil, fmt.Errorf("period start and end are required")
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetGoalResult(ctx, string(id), &cfg.params, toRequestEditors(editors)...)
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
		Data *GoalResult `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing goal result data in response")
	}
	return payload.Data, nil
}
