package v2

import (
	"bytes"
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

type ActivityLocation struct {
	Value           string `json:"value,omitempty"`
	Country         string `json:"country,omitempty"`
	AdminAreaLevel1 string `json:"admin_area_level_1,omitempty"`
	AdminAreaLevel2 string `json:"admin_area_level_2,omitempty"`
	Locality        string `json:"locality,omitempty"`
	Sublocality     string `json:"sublocality,omitempty"`
	Route           string `json:"route,omitempty"`
	StreetNumber    string `json:"street_number,omitempty"`
	Subpremise      string `json:"subpremise,omitempty"`
	PostalCode      string `json:"postal_code,omitempty"`
}

type ActivityParticipant struct {
	PersonID *PersonID `json:"person_id,omitempty"`
	Primary  bool      `json:"primary,omitempty"`
}

type Activity struct {
	ID                ActivityID            `json:"id"`
	Subject           string                `json:"subject,omitempty"`
	Type              string                `json:"type,omitempty"`
	Done              bool                  `json:"done,omitempty"`
	Busy              bool                  `json:"busy,omitempty"`
	DueDate           string                `json:"due_date,omitempty"`
	DueTime           string                `json:"due_time,omitempty"`
	Duration          string                `json:"duration,omitempty"`
	OwnerID           *UserID               `json:"owner_id,omitempty"`
	PersonID          *PersonID             `json:"person_id,omitempty"`
	OrgID             *OrganizationID       `json:"org_id,omitempty"`
	DealID            *DealID               `json:"deal_id,omitempty"`
	LeadID            *LeadID               `json:"lead_id,omitempty"`
	ProjectID         *ProjectID            `json:"project_id,omitempty"`
	Location          *ActivityLocation     `json:"location,omitempty"`
	Participants      []ActivityParticipant `json:"participants,omitempty"`
	Attendees         []ActivityAttendee    `json:"attendees,omitempty"`
	PublicDescription string                `json:"public_description,omitempty"`
	Priority          *int                  `json:"priority,omitempty"`
	Note              string                `json:"note,omitempty"`
	AddTime           *time.Time            `json:"add_time,omitempty"`
	UpdateTime        *time.Time            `json:"update_time,omitempty"`
}

type ActivityDeleteResult struct {
	ID ActivityID `json:"id"`
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

type CreateActivityOption interface {
	applyCreateActivity(*createActivityOptions)
}

type UpdateActivityOption interface {
	applyUpdateActivity(*updateActivityOptions)
}

type DeleteActivityOption interface {
	applyDeleteActivity(*deleteActivityOptions)
}

type ActivityRequestOption interface {
	GetActivityOption
	ListActivitiesOption
	CreateActivityOption
	UpdateActivityOption
	DeleteActivityOption
}

type ActivityOption interface {
	CreateActivityOption
	UpdateActivityOption
}

type getActivityOptions struct {
	params         genv2.GetActivityParams
	requestOptions []pipedrive.RequestOption
}

type listActivitiesOptions struct {
	params         genv2.GetActivitiesParams
	requestOptions []pipedrive.RequestOption
}

type createActivityOptions struct {
	payload        activityPayload
	requestOptions []pipedrive.RequestOption
}

type updateActivityOptions struct {
	payload        activityPayload
	requestOptions []pipedrive.RequestOption
}

type deleteActivityOptions struct {
	requestOptions []pipedrive.RequestOption
}

type activityPayload struct {
	subject           *string
	activityType      *string
	ownerID           *UserID
	dealID            *DealID
	leadID            *LeadID
	personID          *PersonID
	orgID             *OrganizationID
	projectID         *ProjectID
	dueDate           *string
	dueTime           *string
	duration          *string
	busy              *bool
	done              *bool
	location          *ActivityLocation
	participants      []ActivityParticipant
	attendees         []ActivityAttendee
	publicDescription *string
	priority          *int
	note              *string
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

func (o activityRequestOptions) applyCreateActivity(cfg *createActivityOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o activityRequestOptions) applyUpdateActivity(cfg *updateActivityOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o activityRequestOptions) applyDeleteActivity(cfg *deleteActivityOptions) {
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

type activityFieldOption func(*activityPayload)

func (f activityFieldOption) applyCreateActivity(cfg *createActivityOptions) {
	f(&cfg.payload)
}

func (f activityFieldOption) applyUpdateActivity(cfg *updateActivityOptions) {
	f(&cfg.payload)
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

func WithActivitySubject(subject string) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.subject = &subject
	})
}

func WithActivityType(activityType string) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.activityType = &activityType
	})
}

func WithActivityOwnerID(id UserID) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.ownerID = &id
	})
}

func WithActivityDealID(id DealID) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.dealID = &id
	})
}

func WithActivityLeadID(id LeadID) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.leadID = &id
	})
}

func WithActivityPersonID(id PersonID) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.personID = &id
	})
}

func WithActivityOrgID(id OrganizationID) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.orgID = &id
	})
}

func WithActivityProjectID(id ProjectID) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.projectID = &id
	})
}

func WithActivityDueDate(date string) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.dueDate = &date
	})
}

func WithActivityDueTime(time string) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.dueTime = &time
	})
}

func WithActivityDuration(duration string) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.duration = &duration
	})
}

func WithActivityBusy(busy bool) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.busy = &busy
	})
}

func WithActivityDone(done bool) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.done = &done
	})
}

func WithActivityLocation(location ActivityLocation) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.location = &location
	})
}

func WithActivityParticipants(participants ...ActivityParticipant) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		if len(participants) == 0 {
			return
		}
		payload.participants = append(payload.participants, participants...)
	})
}

func WithActivityAttendees(attendees ...ActivityAttendee) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		if len(attendees) == 0 {
			return
		}
		payload.attendees = append(payload.attendees, attendees...)
	})
}

func WithActivityPublicDescription(description string) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.publicDescription = &description
	})
}

func WithActivityPriority(priority int) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.priority = &priority
	})
}

func WithActivityNote(note string) ActivityOption {
	return activityFieldOption(func(payload *activityPayload) {
		payload.note = &note
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

func newCreateActivityOptions(opts []CreateActivityOption) createActivityOptions {
	var cfg createActivityOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateActivity(&cfg)
	}
	return cfg
}

func newUpdateActivityOptions(opts []UpdateActivityOption) updateActivityOptions {
	var cfg updateActivityOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateActivity(&cfg)
	}
	return cfg
}

func newDeleteActivityOptions(opts []DeleteActivityOption) deleteActivityOptions {
	var cfg deleteActivityOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteActivity(&cfg)
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

func (s *ActivitiesService) Create(ctx context.Context, opts ...CreateActivityOption) (*Activity, error) {
	cfg := newCreateActivityOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddActivityWithBodyWithResponse(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *ActivitiesService) Update(ctx context.Context, id ActivityID, opts ...UpdateActivityOption) (*Activity, error) {
	cfg := newUpdateActivityOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateActivityWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *ActivitiesService) Delete(ctx context.Context, id ActivityID, opts ...DeleteActivityOption) (*ActivityDeleteResult, error) {
	cfg := newDeleteActivityOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteActivityWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *ActivityDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing delete activity data in response")
	}
	return payload.Data, nil
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

func (p activityPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.subject != nil {
		body["subject"] = *p.subject
	}
	if p.activityType != nil {
		body["type"] = *p.activityType
	}
	if p.ownerID != nil {
		body["owner_id"] = int(*p.ownerID)
	}
	if p.dealID != nil {
		body["deal_id"] = int(*p.dealID)
	}
	if p.leadID != nil {
		body["lead_id"] = string(*p.leadID)
	}
	if p.personID != nil {
		body["person_id"] = int(*p.personID)
	}
	if p.orgID != nil {
		body["org_id"] = int(*p.orgID)
	}
	if p.projectID != nil {
		body["project_id"] = int(*p.projectID)
	}
	if p.dueDate != nil {
		body["due_date"] = *p.dueDate
	}
	if p.dueTime != nil {
		body["due_time"] = *p.dueTime
	}
	if p.duration != nil {
		body["duration"] = *p.duration
	}
	if p.busy != nil {
		body["busy"] = *p.busy
	}
	if p.done != nil {
		body["done"] = *p.done
	}
	if p.location != nil {
		body["location"] = p.location
	}
	if len(p.participants) > 0 {
		body["participants"] = p.participants
	}
	if len(p.attendees) > 0 {
		body["attendees"] = p.attendees
	}
	if p.publicDescription != nil {
		body["public_description"] = *p.publicDescription
	}
	if p.priority != nil {
		body["priority"] = *p.priority
	}
	if p.note != nil {
		body["note"] = *p.note
	}
	return body
}
