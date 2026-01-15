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

type PersonMarketingStatus string

const (
	PersonMarketingStatusNoConsent    PersonMarketingStatus = "no_consent"
	PersonMarketingStatusUnsubscribed PersonMarketingStatus = "unsubscribed"
	PersonMarketingStatusSubscribed   PersonMarketingStatus = "subscribed"
	PersonMarketingStatusArchived     PersonMarketingStatus = "archived"
)

type PersonSearchField string

const (
	PersonSearchFieldCustomFields PersonSearchField = "custom_fields"
	PersonSearchFieldEmail        PersonSearchField = "email"
	PersonSearchFieldName         PersonSearchField = "name"
	PersonSearchFieldNotes        PersonSearchField = "notes"
	PersonSearchFieldPhone        PersonSearchField = "phone"
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

type PersonSearchItem struct {
	ResultScore float64                `json:"result_score,omitempty"`
	Item        map[string]interface{} `json:"item,omitempty"`
}

type PersonSearchResults struct {
	Items []PersonSearchItem `json:"items,omitempty"`
}

type PersonPicture struct {
	ID            int               `json:"id,omitempty"`
	ItemType      string            `json:"item_type,omitempty"`
	ItemID        *PersonID         `json:"item_id,omitempty"`
	AddedByUserID *UserID           `json:"added_by_user_id,omitempty"`
	Active        bool              `json:"active_flag,omitempty"`
	FileSize      int               `json:"file_size,omitempty"`
	Pictures      map[string]string `json:"pictures,omitempty"`
}

type PersonDeleteResult struct {
	ID PersonID `json:"id"`
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

type CreatePersonOption interface {
	applyCreatePerson(*createPersonOptions)
}

type UpdatePersonOption interface {
	applyUpdatePerson(*updatePersonOptions)
}

type DeletePersonOption interface {
	applyDeletePerson(*deletePersonOptions)
}

type SearchPersonsOption interface {
	applySearchPersons(*searchPersonsOptions)
}

type GetPersonFollowersOption interface {
	applyGetPersonFollowers(*getPersonFollowersOptions)
}

type AddPersonFollowerOption interface {
	applyAddPersonFollower(*addPersonFollowerOptions)
}

type DeletePersonFollowerOption interface {
	applyDeletePersonFollower(*deletePersonFollowerOptions)
}

type GetPersonFollowersChangelogOption interface {
	applyGetPersonFollowersChangelog(*getPersonFollowersChangelogOptions)
}

type GetPersonPictureOption interface {
	applyGetPersonPicture(*getPersonPictureOptions)
}

type PersonRequestOption interface {
	GetPersonOption
	ListPersonsOption
	CreatePersonOption
	UpdatePersonOption
	DeletePersonOption
	SearchPersonsOption
	GetPersonFollowersOption
	AddPersonFollowerOption
	DeletePersonFollowerOption
	GetPersonFollowersChangelogOption
	GetPersonPictureOption
}

type PersonOption interface {
	CreatePersonOption
	UpdatePersonOption
}

type getPersonOptions struct {
	params         genv2.GetPersonParams
	requestOptions []pipedrive.RequestOption
}

type listPersonsOptions struct {
	params         genv2.GetPersonsParams
	requestOptions []pipedrive.RequestOption
}

type createPersonOptions struct {
	payload        personPayload
	requestOptions []pipedrive.RequestOption
}

type updatePersonOptions struct {
	payload        personPayload
	requestOptions []pipedrive.RequestOption
}

type deletePersonOptions struct {
	requestOptions []pipedrive.RequestOption
}

type searchPersonsOptions struct {
	params         genv2.SearchPersonsParams
	requestOptions []pipedrive.RequestOption
}

type getPersonFollowersOptions struct {
	params         genv2.GetPersonFollowersParams
	requestOptions []pipedrive.RequestOption
}

type addPersonFollowerOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deletePersonFollowerOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getPersonFollowersChangelogOptions struct {
	params         genv2.GetPersonFollowersChangelogParams
	requestOptions []pipedrive.RequestOption
}

type getPersonPictureOptions struct {
	requestOptions []pipedrive.RequestOption
}

type personPayload struct {
	name            *string
	ownerID         *UserID
	orgID           *OrganizationID
	emails          []LabeledValue
	phones          []LabeledValue
	labelIDs        []int
	visibleTo       *int
	marketingStatus *PersonMarketingStatus
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

func (o personRequestOptions) applyCreatePerson(cfg *createPersonOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personRequestOptions) applyUpdatePerson(cfg *updatePersonOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personRequestOptions) applyDeletePerson(cfg *deletePersonOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personRequestOptions) applySearchPersons(cfg *searchPersonsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personRequestOptions) applyGetPersonFollowers(cfg *getPersonFollowersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personRequestOptions) applyAddPersonFollower(cfg *addPersonFollowerOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personRequestOptions) applyDeletePersonFollower(cfg *deletePersonFollowerOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personRequestOptions) applyGetPersonFollowersChangelog(cfg *getPersonFollowersChangelogOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o personRequestOptions) applyGetPersonPicture(cfg *getPersonPictureOptions) {
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

type searchPersonsOptionFunc func(*searchPersonsOptions)

func (f searchPersonsOptionFunc) applySearchPersons(cfg *searchPersonsOptions) {
	f(cfg)
}

type getPersonFollowersOptionFunc func(*getPersonFollowersOptions)

func (f getPersonFollowersOptionFunc) applyGetPersonFollowers(cfg *getPersonFollowersOptions) {
	f(cfg)
}

type getPersonFollowersChangelogOptionFunc func(*getPersonFollowersChangelogOptions)

func (f getPersonFollowersChangelogOptionFunc) applyGetPersonFollowersChangelog(cfg *getPersonFollowersChangelogOptions) {
	f(cfg)
}

type personFieldOption func(*personPayload)

func (f personFieldOption) applyCreatePerson(cfg *createPersonOptions) {
	f(&cfg.payload)
}

func (f personFieldOption) applyUpdatePerson(cfg *updatePersonOptions) {
	f(&cfg.payload)
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

func WithPersonName(name string) PersonOption {
	return personFieldOption(func(payload *personPayload) {
		payload.name = &name
	})
}

func WithPersonOwnerID(id UserID) PersonOption {
	return personFieldOption(func(payload *personPayload) {
		payload.ownerID = &id
	})
}

func WithPersonOrgID(id OrganizationID) PersonOption {
	return personFieldOption(func(payload *personPayload) {
		payload.orgID = &id
	})
}

func WithPersonEmails(emails ...LabeledValue) PersonOption {
	return personFieldOption(func(payload *personPayload) {
		if len(emails) == 0 {
			return
		}
		payload.emails = append(payload.emails, emails...)
	})
}

func WithPersonPhones(phones ...LabeledValue) PersonOption {
	return personFieldOption(func(payload *personPayload) {
		if len(phones) == 0 {
			return
		}
		payload.phones = append(payload.phones, phones...)
	})
}

func WithPersonLabelIDs(ids ...int) PersonOption {
	return personFieldOption(func(payload *personPayload) {
		if len(ids) == 0 {
			return
		}
		payload.labelIDs = append(payload.labelIDs, ids...)
	})
}

func WithPersonVisibleTo(visibleTo int) PersonOption {
	return personFieldOption(func(payload *personPayload) {
		payload.visibleTo = &visibleTo
	})
}

func WithPersonMarketingStatus(status PersonMarketingStatus) PersonOption {
	return personFieldOption(func(payload *personPayload) {
		payload.marketingStatus = &status
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

func WithPersonSearchFields(fields ...PersonSearchField) SearchPersonsOption {
	return searchPersonsOptionFunc(func(cfg *searchPersonsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.SearchPersonsParamsFields(csv)
		cfg.params.Fields = &value
	})
}

func WithPersonSearchExactMatch(enabled bool) SearchPersonsOption {
	return searchPersonsOptionFunc(func(cfg *searchPersonsOptions) {
		cfg.params.ExactMatch = &enabled
	})
}

func WithPersonSearchPageSize(limit int) SearchPersonsOption {
	return searchPersonsOptionFunc(func(cfg *searchPersonsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithPersonSearchCursor(cursor string) SearchPersonsOption {
	return searchPersonsOptionFunc(func(cfg *searchPersonsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithPersonFollowersPageSize(limit int) GetPersonFollowersOption {
	return getPersonFollowersOptionFunc(func(cfg *getPersonFollowersOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithPersonFollowersCursor(cursor string) GetPersonFollowersOption {
	return getPersonFollowersOptionFunc(func(cfg *getPersonFollowersOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithPersonFollowersChangelogPageSize(limit int) GetPersonFollowersChangelogOption {
	return getPersonFollowersChangelogOptionFunc(func(cfg *getPersonFollowersChangelogOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithPersonFollowersChangelogCursor(cursor string) GetPersonFollowersChangelogOption {
	return getPersonFollowersChangelogOptionFunc(func(cfg *getPersonFollowersChangelogOptions) {
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

func newCreatePersonOptions(opts []CreatePersonOption) createPersonOptions {
	var cfg createPersonOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreatePerson(&cfg)
	}
	return cfg
}

func newUpdatePersonOptions(opts []UpdatePersonOption) updatePersonOptions {
	var cfg updatePersonOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdatePerson(&cfg)
	}
	return cfg
}

func newDeletePersonOptions(opts []DeletePersonOption) deletePersonOptions {
	var cfg deletePersonOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeletePerson(&cfg)
	}
	return cfg
}

func newSearchPersonsOptions(opts []SearchPersonsOption) searchPersonsOptions {
	var cfg searchPersonsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applySearchPersons(&cfg)
	}
	return cfg
}

func newGetPersonFollowersOptions(opts []GetPersonFollowersOption) getPersonFollowersOptions {
	var cfg getPersonFollowersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetPersonFollowers(&cfg)
	}
	return cfg
}

func newAddPersonFollowerOptions(opts []AddPersonFollowerOption) addPersonFollowerOptions {
	var cfg addPersonFollowerOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddPersonFollower(&cfg)
	}
	return cfg
}

func newDeletePersonFollowerOptions(opts []DeletePersonFollowerOption) deletePersonFollowerOptions {
	var cfg deletePersonFollowerOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeletePersonFollower(&cfg)
	}
	return cfg
}

func newGetPersonFollowersChangelogOptions(opts []GetPersonFollowersChangelogOption) getPersonFollowersChangelogOptions {
	var cfg getPersonFollowersChangelogOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetPersonFollowersChangelog(&cfg)
	}
	return cfg
}

func newGetPersonPictureOptions(opts []GetPersonPictureOption) getPersonPictureOptions {
	var cfg getPersonPictureOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetPersonPicture(&cfg)
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

func (s *PersonsService) Create(ctx context.Context, opts ...CreatePersonOption) (*Person, error) {
	cfg := newCreatePersonOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddPersonWithBodyWithResponse(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *PersonsService) Update(ctx context.Context, id PersonID, opts ...UpdatePersonOption) (*Person, error) {
	cfg := newUpdatePersonOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdatePersonWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *PersonsService) Delete(ctx context.Context, id PersonID, opts ...DeletePersonOption) (*PersonDeleteResult, error) {
	cfg := newDeletePersonOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeletePersonWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *PersonDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing person delete data in response")
	}
	return payload.Data, nil
}

func (s *PersonsService) Search(ctx context.Context, term string, opts ...SearchPersonsOption) (*PersonSearchResults, *string, error) {
	cfg := newSearchPersonsOptions(opts)
	cfg.params.Term = term
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.SearchPersonsWithResponse(ctx, &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           *PersonSearchResults `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, nil, fmt.Errorf("missing person search data in response")
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *PersonsService) ListFollowers(ctx context.Context, id PersonID, opts ...GetPersonFollowersOption) ([]Follower, *string, error) {
	cfg := newGetPersonFollowersOptions(opts)
	return s.listFollowers(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *PersonsService) ListFollowersPager(id PersonID, opts ...GetPersonFollowersOption) *pipedrive.CursorPager[Follower] {
	cfg := newGetPersonFollowersOptions(opts)
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

func (s *PersonsService) ForEachFollowers(ctx context.Context, id PersonID, fn func(Follower) error, opts ...GetPersonFollowersOption) error {
	return s.ListFollowersPager(id, opts...).ForEach(ctx, fn)
}

func (s *PersonsService) AddFollower(ctx context.Context, id PersonID, userID UserID, opts ...AddPersonFollowerOption) (*Follower, error) {
	cfg := newAddPersonFollowerOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	body, err := json.Marshal(map[string]interface{}{
		"user_id": int(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddPersonFollowerWithBodyWithResponse(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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

func (s *PersonsService) DeleteFollower(ctx context.Context, id PersonID, followerID UserID, opts ...DeletePersonFollowerOption) (*FollowerDeleteResult, error) {
	cfg := newDeletePersonFollowerOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeletePersonFollowerWithResponse(ctx, int(id), int(followerID), toRequestEditors(editors)...)
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

func (s *PersonsService) FollowersChangelog(ctx context.Context, id PersonID, opts ...GetPersonFollowersChangelogOption) ([]FollowerChangelog, *string, error) {
	cfg := newGetPersonFollowersChangelogOptions(opts)
	return s.followersChangelog(ctx, id, cfg.params, cfg.requestOptions)
}

func (s *PersonsService) FollowersChangelogPager(id PersonID, opts ...GetPersonFollowersChangelogOption) *pipedrive.CursorPager[FollowerChangelog] {
	cfg := newGetPersonFollowersChangelogOptions(opts)
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

func (s *PersonsService) ForEachFollowersChangelog(ctx context.Context, id PersonID, fn func(FollowerChangelog) error, opts ...GetPersonFollowersChangelogOption) error {
	return s.FollowersChangelogPager(id, opts...).ForEach(ctx, fn)
}

func (s *PersonsService) GetPicture(ctx context.Context, id PersonID, opts ...GetPersonPictureOption) (*PersonPicture, error) {
	cfg := newGetPersonPictureOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetPersonPictureWithResponse(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *PersonPicture `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing person picture data in response")
	}
	return payload.Data, nil
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

func (s *PersonsService) listFollowers(ctx context.Context, id PersonID, params genv2.GetPersonFollowersParams, requestOptions []pipedrive.RequestOption) ([]Follower, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetPersonFollowersWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
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

func (s *PersonsService) followersChangelog(ctx context.Context, id PersonID, params genv2.GetPersonFollowersChangelogParams, requestOptions []pipedrive.RequestOption) ([]FollowerChangelog, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetPersonFollowersChangelogWithResponse(ctx, int(id), &params, toRequestEditors(editors)...)
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

func (p personPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if p.ownerID != nil {
		body["owner_id"] = int(*p.ownerID)
	}
	if p.orgID != nil {
		body["org_id"] = int(*p.orgID)
	}
	if len(p.emails) > 0 {
		body["emails"] = p.emails
	}
	if len(p.phones) > 0 {
		body["phones"] = p.phones
	}
	if len(p.labelIDs) > 0 {
		body["label_ids"] = p.labelIDs
	}
	if p.visibleTo != nil {
		body["visible_to"] = *p.visibleTo
	}
	if p.marketingStatus != nil {
		body["marketing_status"] = *p.marketingStatus
	}
	return body
}
