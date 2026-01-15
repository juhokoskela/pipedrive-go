package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type OrganizationIncludeField string

const (
	OrganizationIncludeFieldActivitiesCount         OrganizationIncludeField = "activities_count"
	OrganizationIncludeFieldClosedDealsCount        OrganizationIncludeField = "closed_deals_count"
	OrganizationIncludeFieldDoneActivitiesCount     OrganizationIncludeField = "done_activities_count"
	OrganizationIncludeFieldEmailMessagesCount      OrganizationIncludeField = "email_messages_count"
	OrganizationIncludeFieldFilesCount              OrganizationIncludeField = "files_count"
	OrganizationIncludeFieldFollowersCount          OrganizationIncludeField = "followers_count"
	OrganizationIncludeFieldLastActivityID          OrganizationIncludeField = "last_activity_id"
	OrganizationIncludeFieldLostDealsCount          OrganizationIncludeField = "lost_deals_count"
	OrganizationIncludeFieldNextActivityID          OrganizationIncludeField = "next_activity_id"
	OrganizationIncludeFieldNotesCount              OrganizationIncludeField = "notes_count"
	OrganizationIncludeFieldOpenDealsCount          OrganizationIncludeField = "open_deals_count"
	OrganizationIncludeFieldPeopleCount             OrganizationIncludeField = "people_count"
	OrganizationIncludeFieldRelatedClosedDealsCount OrganizationIncludeField = "related_closed_deals_count"
	OrganizationIncludeFieldRelatedLostDealsCount   OrganizationIncludeField = "related_lost_deals_count"
	OrganizationIncludeFieldRelatedOpenDealsCount   OrganizationIncludeField = "related_open_deals_count"
	OrganizationIncludeFieldRelatedWonDealsCount    OrganizationIncludeField = "related_won_deals_count"
	OrganizationIncludeFieldUndoneActivitiesCount   OrganizationIncludeField = "undone_activities_count"
	OrganizationIncludeFieldWonDealsCount           OrganizationIncludeField = "won_deals_count"
)

type OrganizationSortField string

const (
	OrganizationSortByAddTime    OrganizationSortField = "add_time"
	OrganizationSortByID         OrganizationSortField = "id"
	OrganizationSortByUpdateTime OrganizationSortField = "update_time"
)

type OrganizationAddress struct {
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

type Organization struct {
	ID           OrganizationID         `json:"id"`
	Name         string                 `json:"name,omitempty"`
	OwnerID      *UserID                `json:"owner_id,omitempty"`
	AddTime      *time.Time             `json:"add_time,omitempty"`
	UpdateTime   *time.Time             `json:"update_time,omitempty"`
	Address      *OrganizationAddress   `json:"address,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

type OrganizationsService struct {
	client *Client
}

type GetOrganizationOption interface {
	applyGetOrganization(*getOrganizationOptions)
}

type ListOrganizationsOption interface {
	applyListOrganizations(*listOrganizationsOptions)
}

type OrganizationRequestOption interface {
	GetOrganizationOption
	ListOrganizationsOption
}

type getOrganizationOptions struct {
	params         genv2.GetOrganizationParams
	requestOptions []pipedrive.RequestOption
}

type listOrganizationsOptions struct {
	params         genv2.GetOrganizationsParams
	requestOptions []pipedrive.RequestOption
}

type organizationRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o organizationRequestOptions) applyGetOrganization(cfg *getOrganizationOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRequestOptions) applyListOrganizations(cfg *listOrganizationsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getOrganizationOptionFunc func(*getOrganizationOptions)

func (f getOrganizationOptionFunc) applyGetOrganization(cfg *getOrganizationOptions) {
	f(cfg)
}

type listOrganizationsOptionFunc func(*listOrganizationsOptions)

func (f listOrganizationsOptionFunc) applyListOrganizations(cfg *listOrganizationsOptions) {
	f(cfg)
}

func WithOrganizationRequestOptions(opts ...pipedrive.RequestOption) OrganizationRequestOption {
	return organizationRequestOptions{requestOptions: opts}
}

func WithOrganizationIncludeFields(fields ...OrganizationIncludeField) GetOrganizationOption {
	return getOrganizationOptionFunc(func(cfg *getOrganizationOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetOrganizationParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithOrganizationCustomFields(fields ...string) GetOrganizationOption {
	return getOrganizationOptionFunc(func(cfg *getOrganizationOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		cfg.params.CustomFields = &csv
	})
}

func WithOrganizationsFilterID(id int) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		cfg.params.FilterId = &id
	})
}

func WithOrganizationsOwnerID(id int) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		cfg.params.OwnerId = &id
	})
}

func WithOrganizationsUpdatedSince(t time.Time) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		value := formatTime(t)
		cfg.params.UpdatedSince = &value
	})
}

func WithOrganizationsUpdatedUntil(t time.Time) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		value := formatTime(t)
		cfg.params.UpdatedUntil = &value
	})
}

func WithOrganizationsSortBy(field OrganizationSortField) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		value := genv2.GetOrganizationsParamsSortBy(field)
		cfg.params.SortBy = &value
	})
}

func WithOrganizationsSortDirection(direction SortDirection) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		value := genv2.GetOrganizationsParamsSortDirection(direction)
		cfg.params.SortDirection = &value
	})
}

func WithOrganizationsIncludeFields(fields ...OrganizationIncludeField) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.GetOrganizationsParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithOrganizationsCustomFields(fields ...string) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		cfg.params.CustomFields = &csv
	})
}

func WithOrganizationsIDs(ids ...OrganizationID) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		csv := joinIDs(ids)
		if csv == "" {
			return
		}
		cfg.params.Ids = &csv
	})
}

func WithOrganizationsPageSize(limit int) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithOrganizationsCursor(cursor string) ListOrganizationsOption {
	return listOrganizationsOptionFunc(func(cfg *listOrganizationsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func newGetOrganizationOptions(opts []GetOrganizationOption) getOrganizationOptions {
	var cfg getOrganizationOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetOrganization(&cfg)
	}
	return cfg
}

func newListOrganizationsOptions(opts []ListOrganizationsOption) listOrganizationsOptions {
	var cfg listOrganizationsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListOrganizations(&cfg)
	}
	return cfg
}

func (s *OrganizationsService) Get(ctx context.Context, id OrganizationID, opts ...GetOrganizationOption) (*Organization, error) {
	cfg := newGetOrganizationOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetOrganizationWithResponse(ctx, int(id), &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Organization `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationsService) List(ctx context.Context, opts ...ListOrganizationsOption) ([]Organization, *string, error) {
	cfg := newListOrganizationsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *OrganizationsService) ListPager(opts ...ListOrganizationsOption) *pipedrive.CursorPager[Organization] {
	cfg := newListOrganizationsOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Organization, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.list(ctx, params, cfg.requestOptions)
	})
}

func (s *OrganizationsService) ForEach(ctx context.Context, fn func(Organization) error, opts ...ListOrganizationsOption) error {
	return s.ListPager(opts...).ForEach(ctx, fn)
}

func (s *OrganizationsService) list(ctx context.Context, params genv2.GetOrganizationsParams, requestOptions []pipedrive.RequestOption) ([]Organization, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetOrganizationsWithResponse(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Organization `json:"data"`
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
