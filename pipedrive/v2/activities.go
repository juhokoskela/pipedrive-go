package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ActivityIncludeField string

const (
	ActivityIncludeFieldAttendees ActivityIncludeField = "attendees"
)

type ActivitySortField string

const (
	ActivitySortByAddTime    ActivitySortField = "add_time"
	ActivitySortByDueDate    ActivitySortField = "due_date"
	ActivitySortByID         ActivitySortField = "id"
	ActivitySortByUpdateTime ActivitySortField = "update_time"
)

type ActivityAttendee struct {
	Email       string    `json:"email,omitempty"`
	Name        string    `json:"name,omitempty"`
	Status      string    `json:"status,omitempty"`
	PersonID    *PersonID `json:"person_id,omitempty"`
	UserID      *UserID   `json:"user_id,omitempty"`
	IsOrganizer bool      `json:"is_organizer,omitempty"`
}

type Activity struct {
	ID         ActivityID         `json:"id"`
	Subject    string             `json:"subject,omitempty"`
	Type       string             `json:"type,omitempty"`
	Done       bool               `json:"done,omitempty"`
	DueDate    string             `json:"due_date,omitempty"`
	DueTime    string             `json:"due_time,omitempty"`
	OwnerID    *UserID            `json:"owner_id,omitempty"`
	PersonID   *PersonID          `json:"person_id,omitempty"`
	OrgID      *OrganizationID    `json:"org_id,omitempty"`
	DealID     *DealID            `json:"deal_id,omitempty"`
	LeadID     *LeadID            `json:"lead_id,omitempty"`
	AddTime    *time.Time         `json:"add_time,omitempty"`
	UpdateTime *time.Time         `json:"update_time,omitempty"`
	Attendees  []ActivityAttendee `json:"attendees,omitempty"`
}

type ActivitiesService struct {
	client *Client
}

type GetActivityOption interface {
	applyGetActivity(*getActivityOptions)
}

type ListActivitiesOption interface {
	applyListActivities(*listActivitiesOptions)
}

type ActivityRequestOption interface {
	GetActivityOption
	ListActivitiesOption
}

type getActivityOptions struct {
	params         genv2.GetActivityParams
	requestOptions []pipedrive.RequestOption
}

type listActivitiesOptions struct {
	params         genv2.GetActivitiesParams
	requestOptions []pipedrive.RequestOption
}

type activityRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o activityRequestOptions) applyGetActivity(cfg *getActivityOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o activityRequestOptions) applyListActivities(cfg *listActivitiesOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getActivityOptionFunc func(*getActivityOptions)

func (f getActivityOptionFunc) applyGetActivity(cfg *getActivityOptions) {
	f(cfg)
}

type listActivitiesOptionFunc func(*listActivitiesOptions)

func (f listActivitiesOptionFunc) applyListActivities(cfg *listActivitiesOptions) {
	f(cfg)
}

func WithActivityRequestOptions(opts ...pipedrive.RequestOption) ActivityRequestOption {
	return activityRequestOptions{requestOptions: opts}
}

func WithActivityIncludeFields(fields ...ActivityIncludeField) GetActivityOption {
	return getActivityOptionFunc(func(cfg *getActivityOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetActivityParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithActivitiesFilterID(id int) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		cfg.params.FilterId = &id
	})
}

func WithActivitiesOwnerID(id int) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		cfg.params.OwnerId = &id
	})
}

func WithActivitiesDealID(id DealID) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		value := int(id)
		cfg.params.DealId = &value
	})
}

func WithActivitiesLeadID(id LeadID) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		value := string(id)
		cfg.params.LeadId = &value
	})
}

func WithActivitiesPersonID(id PersonID) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		value := int(id)
		cfg.params.PersonId = &value
	})
}

func WithActivitiesOrgID(id OrganizationID) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		value := int(id)
		cfg.params.OrgId = &value
	})
}

func WithActivitiesDone(done bool) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		cfg.params.Done = &done
	})
}

func WithActivitiesUpdatedSince(t time.Time) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		value := formatTime(t)
		cfg.params.UpdatedSince = &value
	})
}

func WithActivitiesUpdatedUntil(t time.Time) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		value := formatTime(t)
		cfg.params.UpdatedUntil = &value
	})
}

func WithActivitiesSortBy(field ActivitySortField) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		value := genv2.GetActivitiesParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithActivitiesSortDirection(direction SortDirection) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		value := genv2.GetActivitiesParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithActivitiesIncludeFields(fields ...ActivityIncludeField) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetActivitiesParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithActivitiesIDs(ids ...ActivityID) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		csv := joinIDs(ids)
		if csv == "" {
			return
		}
		cfg.params.Ids = &csv
	})
}

func WithActivitiesPageSize(limit int) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithActivitiesCursor(cursor string) ListActivitiesOption {
	return listActivitiesOptionFunc(func(cfg *listActivitiesOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func newGetActivityOptions(opts []GetActivityOption) getActivityOptions {
	var cfg getActivityOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetActivity(&cfg)
	}
	return cfg
}

func newListActivitiesOptions(opts []ListActivitiesOption) listActivitiesOptions {
	var cfg listActivitiesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListActivities(&cfg)
	}
	return cfg
}

func (s *ActivitiesService) Get(ctx context.Context, id ActivityID, opts ...GetActivityOption) (*Activity, error) {
	cfg := newGetActivityOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetActivityWithResponse(ctx, int(id), &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Activity `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing activity data in response")
	}
	return payload.Data, nil
}

func (s *ActivitiesService) List(ctx context.Context, opts ...ListActivitiesOption) ([]Activity, *string, error) {
	cfg := newListActivitiesOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *ActivitiesService) ListPager(opts ...ListActivitiesOption) *pipedrive.CursorPager[Activity] {
	cfg := newListActivitiesOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Activity, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *ActivitiesService) ForEach(ctx context.Context, fn func(Activity) error, opts ...ListActivitiesOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *ActivitiesService) list(ctx context.Context, params genv2.GetActivitiesParams, requestOptions []pipedrive.RequestOption) ([]Activity, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetActivitiesWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Activity `json:"data"`
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
	return payload.Data, next, nil
}
