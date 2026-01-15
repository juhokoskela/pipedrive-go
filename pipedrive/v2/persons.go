package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type PersonIncludeField string

const (
	PersonIncludeFieldActivitiesCount             PersonIncludeField = "activities_count"
	PersonIncludeFieldClosedDealsCount            PersonIncludeField = "closed_deals_count"
	PersonIncludeFieldDoiStatus                   PersonIncludeField = "doi_status"
	PersonIncludeFieldDoneActivitiesCount         PersonIncludeField = "done_activities_count"
	PersonIncludeFieldEmailMessagesCount          PersonIncludeField = "email_messages_count"
	PersonIncludeFieldFilesCount                  PersonIncludeField = "files_count"
	PersonIncludeFieldFollowersCount              PersonIncludeField = "followers_count"
	PersonIncludeFieldLastActivityID              PersonIncludeField = "last_activity_id"
	PersonIncludeFieldLastIncomingMailTime        PersonIncludeField = "last_incoming_mail_time"
	PersonIncludeFieldLastOutgoingMailTime        PersonIncludeField = "last_outgoing_mail_time"
	PersonIncludeFieldLostDealsCount              PersonIncludeField = "lost_deals_count"
	PersonIncludeFieldMarketingStatus             PersonIncludeField = "marketing_status"
	PersonIncludeFieldNextActivityID              PersonIncludeField = "next_activity_id"
	PersonIncludeFieldNotesCount                  PersonIncludeField = "notes_count"
	PersonIncludeFieldOpenDealsCount              PersonIncludeField = "open_deals_count"
	PersonIncludeFieldParticipantClosedDealsCount PersonIncludeField = "participant_closed_deals_count"
	PersonIncludeFieldParticipantOpenDealsCount   PersonIncludeField = "participant_open_deals_count"
	PersonIncludeFieldRelatedClosedDealsCount     PersonIncludeField = "related_closed_deals_count"
	PersonIncludeFieldRelatedLostDealsCount       PersonIncludeField = "related_lost_deals_count"
	PersonIncludeFieldRelatedOpenDealsCount       PersonIncludeField = "related_open_deals_count"
	PersonIncludeFieldRelatedWonDealsCount        PersonIncludeField = "related_won_deals_count"
	PersonIncludeFieldUndoneActivitiesCount       PersonIncludeField = "undone_activities_count"
	PersonIncludeFieldWonDealsCount               PersonIncludeField = "won_deals_count"
)

type PersonSortField string

const (
	PersonSortByAddTime    PersonSortField = "add_time"
	PersonSortByID         PersonSortField = "id"
	PersonSortByUpdateTime PersonSortField = "update_time"
)

type LabeledValue struct {
	Value   string `json:"value,omitempty"`
	Primary bool   `json:"primary,omitempty"`
	Label   string `json:"label,omitempty"`
}

type Person struct {
	ID           PersonID               `json:"id"`
	Name         string                 `json:"name,omitempty"`
	FirstName    string                 `json:"first_name,omitempty"`
	LastName     string                 `json:"last_name,omitempty"`
	OwnerID      *UserID                `json:"owner_id,omitempty"`
	OrgID        *OrganizationID        `json:"org_id,omitempty"`
	AddTime      *time.Time             `json:"add_time,omitempty"`
	UpdateTime   *time.Time             `json:"update_time,omitempty"`
	Emails       []LabeledValue         `json:"emails,omitempty"`
	Phones       []LabeledValue         `json:"phones,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

type PersonsService struct {
	client *Client
}

type GetPersonOption interface {
	applyGetPerson(*getPersonOptions)
}

type ListPersonsOption interface {
	applyListPersons(*listPersonsOptions)
}

type PersonRequestOption interface {
	GetPersonOption
	ListPersonsOption
}

type getPersonOptions struct {
	params         genv2.GetPersonParams
	requestOptions []pipedrive.RequestOption
}

type listPersonsOptions struct {
	params         genv2.GetPersonsParams
	requestOptions []pipedrive.RequestOption
}

type personRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o personRequestOptions) applyGetPerson(cfg *getPersonOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personRequestOptions) applyListPersons(cfg *listPersonsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getPersonOptionFunc func(*getPersonOptions)

func (f getPersonOptionFunc) applyGetPerson(cfg *getPersonOptions) {
	f(cfg)
}

type listPersonsOptionFunc func(*listPersonsOptions)

func (f listPersonsOptionFunc) applyListPersons(cfg *listPersonsOptions) {
	f(cfg)
}

func WithPersonRequestOptions(opts ...pipedrive.RequestOption) PersonRequestOption {
	return personRequestOptions{requestOptions: opts}
}

func WithPersonIncludeFields(fields ...PersonIncludeField) GetPersonOption {
	return getPersonOptionFunc(func(cfg *getPersonOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetPersonParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithPersonCustomFields(fields ...string) GetPersonOption {
	return getPersonOptionFunc(func(cfg *getPersonOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		cfg.params.CustomFields = &csv
	})
}

func WithPersonsFilterID(id int) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		cfg.params.FilterId = &id
	})
}

func WithPersonsOwnerID(id int) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		cfg.params.OwnerId = &id
	})
}

func WithPersonsOrgID(id OrganizationID) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		value := int(id)
		cfg.params.OrgId = &value
	})
}

func WithPersonsDealID(id DealID) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		value := int(id)
		cfg.params.DealId = &value
	})
}

func WithPersonsUpdatedSince(t time.Time) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		value := formatTime(t)
		cfg.params.UpdatedSince = &value
	})
}

func WithPersonsUpdatedUntil(t time.Time) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		value := formatTime(t)
		cfg.params.UpdatedUntil = &value
	})
}

func WithPersonsSortBy(field PersonSortField) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		value := genv2.GetPersonsParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithPersonsSortDirection(direction SortDirection) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		value := genv2.GetPersonsParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithPersonsIncludeFields(fields ...PersonIncludeField) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetPersonsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithPersonsCustomFields(fields ...string) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		cfg.params.CustomFields = &csv
	})
}

func WithPersonsIDs(ids ...PersonID) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		csv := joinIDs(ids)
		if csv == "" {
			return
		}
		cfg.params.Ids = &csv
	})
}

func WithPersonsPageSize(limit int) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithPersonsCursor(cursor string) ListPersonsOption {
	return listPersonsOptionFunc(func(cfg *listPersonsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func newGetPersonOptions(opts []GetPersonOption) getPersonOptions {
	var cfg getPersonOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetPerson(&cfg)
	}
	return cfg
}

func newListPersonsOptions(opts []ListPersonsOption) listPersonsOptions {
	var cfg listPersonsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListPersons(&cfg)
	}
	return cfg
}

func (s *PersonsService) Get(ctx context.Context, id PersonID, opts ...GetPersonOption) (*Person, error) {
	cfg := newGetPersonOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetPersonWithResponse(ctx, int(id), &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Person `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing person data in response")
	}
	return payload.Data, nil
}

func (s *PersonsService) List(ctx context.Context, opts ...ListPersonsOption) ([]Person, *string, error) {
	cfg := newListPersonsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *PersonsService) ListPager(opts ...ListPersonsOption) *pipedrive.CursorPager[Person] {
	cfg := newListPersonsOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Person, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *PersonsService) ForEach(ctx context.Context, fn func(Person) error, opts ...ListPersonsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *PersonsService) list(ctx context.Context, params genv2.GetPersonsParams, requestOptions []pipedrive.RequestOption) ([]Person, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetPersonsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Person `json:"data"`
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
