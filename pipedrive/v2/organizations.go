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

type OrganizationSearchField string

const (
	OrganizationSearchFieldAddress      OrganizationSearchField = "address"
	OrganizationSearchFieldCustomFields OrganizationSearchField = "custom_fields"
	OrganizationSearchFieldName         OrganizationSearchField = "name"
	OrganizationSearchFieldNotes        OrganizationSearchField = "notes"
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

type OrganizationSearchItem struct {
	ResultScore float64                `json:"result_score,omitempty"`
	Item        map[string]interface{} `json:"item,omitempty"`
}

type OrganizationSearchResults struct {
	Items []OrganizationSearchItem `json:"items,omitempty"`
}

type OrganizationDeleteResult struct {
	ID OrganizationID `json:"id"`
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

type CreateOrganizationOption interface {
	applyCreateOrganization(*createOrganizationOptions)
}

type UpdateOrganizationOption interface {
	applyUpdateOrganization(*updateOrganizationOptions)
}

type DeleteOrganizationOption interface {
	applyDeleteOrganization(*deleteOrganizationOptions)
}

type SearchOrganizationsOption interface {
	applySearchOrganizations(*searchOrganizationsOptions)
}

type GetOrganizationFollowersOption interface {
	applyGetOrganizationFollowers(*getOrganizationFollowersOptions)
}

type AddOrganizationFollowerOption interface {
	applyAddOrganizationFollower(*addOrganizationFollowerOptions)
}

type DeleteOrganizationFollowerOption interface {
	applyDeleteOrganizationFollower(*deleteOrganizationFollowerOptions)
}

type GetOrganizationFollowersChangelogOption interface {
	applyGetOrganizationFollowersChangelog(*getOrganizationFollowersChangelogOptions)
}

type OrganizationRequestOption interface {
	GetOrganizationOption
	ListOrganizationsOption
	CreateOrganizationOption
	UpdateOrganizationOption
	DeleteOrganizationOption
	SearchOrganizationsOption
	GetOrganizationFollowersOption
	AddOrganizationFollowerOption
	DeleteOrganizationFollowerOption
	GetOrganizationFollowersChangelogOption
}

type OrganizationOption interface {
	CreateOrganizationOption
	UpdateOrganizationOption
}

type getOrganizationOptions struct {
	params         genv2.GetOrganizationParams
	requestOptions []pipedrive.RequestOption
}

type listOrganizationsOptions struct {
	params         genv2.GetOrganizationsParams
	requestOptions []pipedrive.RequestOption
}

type createOrganizationOptions struct {
	payload        organizationPayload
	requestOptions []pipedrive.RequestOption
}

type updateOrganizationOptions struct {
	payload        organizationPayload
	requestOptions []pipedrive.RequestOption
}

type deleteOrganizationOptions struct {
	requestOptions []pipedrive.RequestOption
}

type searchOrganizationsOptions struct {
	params         genv2.SearchOrganizationParams
	requestOptions []pipedrive.RequestOption
}

type getOrganizationFollowersOptions struct {
	params         genv2.GetOrganizationFollowersParams
	requestOptions []pipedrive.RequestOption
}

type addOrganizationFollowerOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deleteOrganizationFollowerOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getOrganizationFollowersChangelogOptions struct {
	params         genv2.GetOrganizationFollowersChangelogParams
	requestOptions []pipedrive.RequestOption
}

type organizationPayload struct {
	name         *string
	ownerID      *UserID
	address      *OrganizationAddress
	labelIDs     []int
	visibleTo    *int
	customFields map[string]interface{}
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

func (o organizationRequestOptions) applyCreateOrganization(cfg *createOrganizationOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRequestOptions) applyUpdateOrganization(cfg *updateOrganizationOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRequestOptions) applyDeleteOrganization(cfg *deleteOrganizationOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRequestOptions) applySearchOrganizations(cfg *searchOrganizationsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRequestOptions) applyGetOrganizationFollowers(cfg *getOrganizationFollowersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRequestOptions) applyAddOrganizationFollower(cfg *addOrganizationFollowerOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRequestOptions) applyDeleteOrganizationFollower(cfg *deleteOrganizationFollowerOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRequestOptions) applyGetOrganizationFollowersChangelog(cfg *getOrganizationFollowersChangelogOptions) {
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

type searchOrganizationsOptionFunc func(*searchOrganizationsOptions)

func (f searchOrganizationsOptionFunc) applySearchOrganizations(cfg *searchOrganizationsOptions) {
	f(cfg)
}

type getOrganizationFollowersOptionFunc func(*getOrganizationFollowersOptions)

func (f getOrganizationFollowersOptionFunc) applyGetOrganizationFollowers(cfg *getOrganizationFollowersOptions) {
	f(cfg)
}

type getOrganizationFollowersChangelogOptionFunc func(*getOrganizationFollowersChangelogOptions)

func (f getOrganizationFollowersChangelogOptionFunc) applyGetOrganizationFollowersChangelog(cfg *getOrganizationFollowersChangelogOptions) {
	f(cfg)
}

type organizationFieldOption func(*organizationPayload)

func (f organizationFieldOption) applyCreateOrganization(cfg *createOrganizationOptions) {
	f(&cfg.payload)
}

func (f organizationFieldOption) applyUpdateOrganization(cfg *updateOrganizationOptions) {
	f(&cfg.payload)
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

func WithOrganizationName(name string) OrganizationOption {
	return organizationFieldOption(func(payload *organizationPayload) {
		payload.name = &name
	})
}

func WithOrganizationOwnerID(id UserID) OrganizationOption {
	return organizationFieldOption(func(payload *organizationPayload) {
		payload.ownerID = &id
	})
}

func WithOrganizationAddress(address OrganizationAddress) OrganizationOption {
	return organizationFieldOption(func(payload *organizationPayload) {
		payload.address = &address
	})
}

func WithOrganizationLabelIDs(ids ...int) OrganizationOption {
	return organizationFieldOption(func(payload *organizationPayload) {
		if len(ids) == 0 {
			return
		}
		payload.labelIDs = append(payload.labelIDs, ids...)
	})
}

func WithOrganizationVisibleTo(visibleTo int) OrganizationOption {
	return organizationFieldOption(func(payload *organizationPayload) {
		payload.visibleTo = &visibleTo
	})
}

func WithOrganizationCustomFieldsMap(fields map[string]interface{}) OrganizationOption {
	return organizationFieldOption(func(payload *organizationPayload) {
		payload.customFields = fields
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

func WithOrganizationSearchFields(fields ...OrganizationSearchField) SearchOrganizationsOption {
	return searchOrganizationsOptionFunc(func(cfg *searchOrganizationsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.SearchOrganizationParamsFields(csv)
		cfg.params.Fields = &value
	})
}

func WithOrganizationSearchExactMatch(enabled bool) SearchOrganizationsOption {
	return searchOrganizationsOptionFunc(func(cfg *searchOrganizationsOptions) {
		cfg.params.ExactMatch = &enabled
	})
}

func WithOrganizationSearchPageSize(limit int) SearchOrganizationsOption {
	return searchOrganizationsOptionFunc(func(cfg *searchOrganizationsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithOrganizationSearchCursor(cursor string) SearchOrganizationsOption {
	return searchOrganizationsOptionFunc(func(cfg *searchOrganizationsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithOrganizationFollowersPageSize(limit int) GetOrganizationFollowersOption {
	return getOrganizationFollowersOptionFunc(func(cfg *getOrganizationFollowersOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithOrganizationFollowersCursor(cursor string) GetOrganizationFollowersOption {
	return getOrganizationFollowersOptionFunc(func(cfg *getOrganizationFollowersOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithOrganizationFollowersChangelogPageSize(limit int) GetOrganizationFollowersChangelogOption {
	return getOrganizationFollowersChangelogOptionFunc(func(cfg *getOrganizationFollowersChangelogOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithOrganizationFollowersChangelogCursor(cursor string) GetOrganizationFollowersChangelogOption {
	return getOrganizationFollowersChangelogOptionFunc(func(cfg *getOrganizationFollowersChangelogOptions) {
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

func newCreateOrganizationOptions(opts []CreateOrganizationOption) createOrganizationOptions {
	var cfg createOrganizationOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateOrganization(&cfg)
	}
	return cfg
}

func newUpdateOrganizationOptions(opts []UpdateOrganizationOption) updateOrganizationOptions {
	var cfg updateOrganizationOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateOrganization(&cfg)
	}
	return cfg
}

func newDeleteOrganizationOptions(opts []DeleteOrganizationOption) deleteOrganizationOptions {
	var cfg deleteOrganizationOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteOrganization(&cfg)
	}
	return cfg
}

func newSearchOrganizationsOptions(opts []SearchOrganizationsOption) searchOrganizationsOptions {
	var cfg searchOrganizationsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applySearchOrganizations(&cfg)
	}
	return cfg
}

func newGetOrganizationFollowersOptions(opts []GetOrganizationFollowersOption) getOrganizationFollowersOptions {
	var cfg getOrganizationFollowersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetOrganizationFollowers(&cfg)
	}
	return cfg
}

func newAddOrganizationFollowerOptions(opts []AddOrganizationFollowerOption) addOrganizationFollowerOptions {
	var cfg addOrganizationFollowerOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddOrganizationFollower(&cfg)
	}
	return cfg
}

func newDeleteOrganizationFollowerOptions(opts []DeleteOrganizationFollowerOption) deleteOrganizationFollowerOptions {
	var cfg deleteOrganizationFollowerOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteOrganizationFollower(&cfg)
	}
	return cfg
}

func newGetOrganizationFollowersChangelogOptions(opts []GetOrganizationFollowersChangelogOption) getOrganizationFollowersChangelogOptions {
	var cfg getOrganizationFollowersChangelogOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetOrganizationFollowersChangelog(&cfg)
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

func (s *OrganizationsService) Create(ctx context.Context, opts ...CreateOrganizationOption) (*Organization, error) {
	cfg := newCreateOrganizationOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddOrganizationWithBodyWithResponse(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *OrganizationsService) Update(ctx context.Context, id OrganizationID, opts ...UpdateOrganizationOption) (*Organization, error) {
	cfg := newUpdateOrganizationOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateOrganizationWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *OrganizationsService) Delete(ctx context.Context, id OrganizationID, opts ...DeleteOrganizationOption) (*OrganizationDeleteResult, error) {
	cfg := newDeleteOrganizationOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteOrganizationWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *OrganizationDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization delete data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationsService) Search(ctx context.Context, term string, opts ...SearchOrganizationsOption) (*OrganizationSearchResults, *string, error) {
	cfg := newSearchOrganizationsOptions(opts)
	cfg.params.Term = term
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.SearchOrganizationWithResponse(ctx, &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           *OrganizationSearchResults `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, nil, fmt.Errorf("missing organization search data in response")
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *OrganizationsService) ListFollowers(ctx context.Context, id OrganizationID, opts ...GetOrganizationFollowersOption) ([]Follower, *string, error) {
	cfg := newGetOrganizationFollowersOptions(opts)
	return s.listFollowers(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *OrganizationsService) ListFollowersPager(id OrganizationID, opts ...GetOrganizationFollowersOption) *pipedrive.CursorPager[Follower] {
	cfg := newGetOrganizationFollowersOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Follower, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listFollowers(ctx, id, params, cfg.requestOptions)
	})
}

func (s *OrganizationsService) ForEachFollowers(ctx context.Context, id OrganizationID, fn func(Follower) error, opts ...GetOrganizationFollowersOption) error {
	return s.ListFollowersPager(id, opts...).ForEach(ctx, fn)
}

func (s *OrganizationsService) AddFollower(ctx context.Context, id OrganizationID, userID UserID, opts ...AddOrganizationFollowerOption) (*Follower, error) {
	cfg := newAddOrganizationFollowerOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(map[string]interface{}{
		"user_id": int(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddOrganizationFollowerWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Follower `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing follower data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationsService) DeleteFollower(ctx context.Context, id OrganizationID, followerID UserID, opts ...DeleteOrganizationFollowerOption) (*FollowerDeleteResult, error) {
	cfg := newDeleteOrganizationFollowerOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteOrganizationFollowerWithResponse(ctx, int(id), int(followerID), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *FollowerDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing delete follower data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationsService) FollowersChangelog(ctx context.Context, id OrganizationID, opts ...GetOrganizationFollowersChangelogOption) ([]FollowerChangelog, *string, error) {
	cfg := newGetOrganizationFollowersChangelogOptions(opts)
	return s.followersChangelog(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *OrganizationsService) FollowersChangelogPager(id OrganizationID, opts ...GetOrganizationFollowersChangelogOption) *pipedrive.CursorPager[FollowerChangelog] {
	cfg := newGetOrganizationFollowersChangelogOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]FollowerChangelog, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.followersChangelog(ctx, id, params, cfg.requestOptions)
	})
}

func (s *OrganizationsService) ForEachFollowersChangelog(ctx context.Context, id OrganizationID, fn func(FollowerChangelog) error, opts ...GetOrganizationFollowersChangelogOption) error {
	return s.FollowersChangelogPager(id, opts...).ForEach(ctx, fn)
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

func (s *OrganizationsService) listFollowers(ctx context.Context, id OrganizationID, params genv2.GetOrganizationFollowersParams, requestOptions []pipedrive.RequestOption) ([]Follower, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetOrganizationFollowersWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Follower `json:"data"`
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

func (s *OrganizationsService) followersChangelog(ctx context.Context, id OrganizationID, params genv2.GetOrganizationFollowersChangelogParams, requestOptions []pipedrive.RequestOption) ([]FollowerChangelog, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetOrganizationFollowersChangelogWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []FollowerChangelog `json:"data"`
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

func (p organizationPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if p.ownerID != nil {
		body["owner_id"] = int(*p.ownerID)
	}
	if p.address != nil {
		body["address"] = p.address
	}
	if len(p.labelIDs) > 0 {
		body["label_ids"] = p.labelIDs
	}
	if p.visibleTo != nil {
		body["visible_to"] = *p.visibleTo
	}
	if p.customFields != nil {
		body["custom_fields"] = p.customFields
	}
	return body
}
