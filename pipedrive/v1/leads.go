package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type LeadValue struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Lead struct {
	ID                LeadID          `json:"id"`
	Title             string          `json:"title,omitempty"`
	OwnerID           *UserID         `json:"owner_id,omitempty"`
	CreatorID         *UserID         `json:"creator_id,omitempty"`
	LabelIDs          []LeadLabelID   `json:"label_ids,omitempty"`
	PersonID          *PersonID       `json:"person_id,omitempty"`
	OrganizationID    *OrganizationID `json:"organization_id,omitempty"`
	SourceName        string          `json:"source_name,omitempty"`
	Origin            string          `json:"origin,omitempty"`
	OriginID          *string         `json:"origin_id,omitempty"`
	Channel           *int            `json:"channel,omitempty"`
	ChannelID         *string         `json:"channel_id,omitempty"`
	IsArchived        bool            `json:"is_archived,omitempty"`
	WasSeen           bool            `json:"was_seen,omitempty"`
	Value             *LeadValue      `json:"value,omitempty"`
	ExpectedCloseDate *string         `json:"expected_close_date,omitempty"`
	NextActivityID    *ActivityID     `json:"next_activity_id,omitempty"`
	AddTime           *DateTime       `json:"add_time,omitempty"`
	UpdateTime        *DateTime       `json:"update_time,omitempty"`
	VisibleTo         string          `json:"visible_to,omitempty"`
	CCEmail           string          `json:"cc_email,omitempty"`
}

type LeadPagination struct {
	Start                 int  `json:"start,omitempty"`
	Limit                 int  `json:"limit,omitempty"`
	MoreItemsInCollection bool `json:"more_items_in_collection,omitempty"`
}

type LeadDeleteResult struct {
	ID LeadID `json:"id"`
}

type LeadsService struct {
	client *Client
}

type ListLeadsOption interface {
	applyListLeads(*listLeadsOptions)
}

type ListArchivedLeadsOption interface {
	applyListArchivedLeads(*listArchivedLeadsOptions)
}

type GetLeadOption interface {
	applyGetLead(*getLeadOptions)
}

type CreateLeadOption interface {
	applyCreateLead(*createLeadOptions)
}

type UpdateLeadOption interface {
	applyUpdateLead(*updateLeadOptions)
}

type DeleteLeadOption interface {
	applyDeleteLead(*deleteLeadOptions)
}

type ListLeadUsersOption interface {
	applyListLeadUsers(*listLeadUsersOptions)
}

type LeadsRequestOption interface {
	ListLeadsOption
	ListArchivedLeadsOption
	GetLeadOption
	CreateLeadOption
	UpdateLeadOption
	DeleteLeadOption
	ListLeadUsersOption
}

type LeadOption interface {
	CreateLeadOption
	UpdateLeadOption
}

type listLeadsOptions struct {
	params         genv1.GetLeadsParams
	requestOptions []pipedrive.RequestOption
}

type listArchivedLeadsOptions struct {
	params         genv1.GetArchivedLeadsParams
	requestOptions []pipedrive.RequestOption
}

type getLeadOptions struct {
	requestOptions []pipedrive.RequestOption
}

type createLeadOptions struct {
	payload        leadPayload
	requestOptions []pipedrive.RequestOption
}

type updateLeadOptions struct {
	payload        leadPayload
	requestOptions []pipedrive.RequestOption
}

type deleteLeadOptions struct {
	requestOptions []pipedrive.RequestOption
}

type listLeadUsersOptions struct {
	requestOptions []pipedrive.RequestOption
}

type leadPayload struct {
	title             *string
	ownerID           *UserID
	labelIDs          []LeadLabelID
	personID          *PersonID
	orgID             *OrganizationID
	value             *LeadValue
	expectedCloseDate *string
	visibleTo         *string
	wasSeen           *bool
	originID          *string
	channel           *int
	channelID         *string
	isArchived        *bool
}

type leadsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o leadsRequestOptions) applyListLeads(cfg *listLeadsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadsRequestOptions) applyListArchivedLeads(cfg *listArchivedLeadsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadsRequestOptions) applyGetLead(cfg *getLeadOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadsRequestOptions) applyCreateLead(cfg *createLeadOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadsRequestOptions) applyUpdateLead(cfg *updateLeadOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadsRequestOptions) applyDeleteLead(cfg *deleteLeadOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadsRequestOptions) applyListLeadUsers(cfg *listLeadUsersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listLeadsOptionFunc func(*listLeadsOptions)

func (f listLeadsOptionFunc) applyListLeads(cfg *listLeadsOptions) {
	f(cfg)
}

type listArchivedLeadsOptionFunc func(*listArchivedLeadsOptions)

func (f listArchivedLeadsOptionFunc) applyListArchivedLeads(cfg *listArchivedLeadsOptions) {
	f(cfg)
}

type leadFieldOption func(*leadPayload)

func (f leadFieldOption) applyCreateLead(cfg *createLeadOptions) {
	f(&cfg.payload)
}

func (f leadFieldOption) applyUpdateLead(cfg *updateLeadOptions) {
	f(&cfg.payload)
}

func WithLeadsRequestOptions(opts ...pipedrive.RequestOption) LeadsRequestOption {
	return leadsRequestOptions{requestOptions: opts}
}

func WithLeadsLimit(limit int) ListLeadsOption {
	return listLeadsOptionFunc(func(cfg *listLeadsOptions) {
		cfg.params.Limit = &limit
	})
}

func WithLeadsStart(start int) ListLeadsOption {
	return listLeadsOptionFunc(func(cfg *listLeadsOptions) {
		cfg.params.Start = &start
	})
}

func WithLeadsOwnerID(id UserID) ListLeadsOption {
	return listLeadsOptionFunc(func(cfg *listLeadsOptions) {
		value := int(id)
		cfg.params.OwnerId = &value
	})
}

func WithLeadsPersonID(id PersonID) ListLeadsOption {
	return listLeadsOptionFunc(func(cfg *listLeadsOptions) {
		value := int(id)
		cfg.params.PersonId = &value
	})
}

func WithLeadsOrganizationID(id OrganizationID) ListLeadsOption {
	return listLeadsOptionFunc(func(cfg *listLeadsOptions) {
		value := int(id)
		cfg.params.OrganizationId = &value
	})
}

func WithLeadsFilterID(id int) ListLeadsOption {
	return listLeadsOptionFunc(func(cfg *listLeadsOptions) {
		cfg.params.FilterId = &id
	})
}

func WithLeadsSort(sort string) ListLeadsOption {
	return listLeadsOptionFunc(func(cfg *listLeadsOptions) {
		value := genv1.GetLeadsParamsSort(sort)
		cfg.params.Sort = &value
	})
}

func WithArchivedLeadsLimit(limit int) ListArchivedLeadsOption {
	return listArchivedLeadsOptionFunc(func(cfg *listArchivedLeadsOptions) {
		cfg.params.Limit = &limit
	})
}

func WithArchivedLeadsStart(start int) ListArchivedLeadsOption {
	return listArchivedLeadsOptionFunc(func(cfg *listArchivedLeadsOptions) {
		cfg.params.Start = &start
	})
}

func WithArchivedLeadsOwnerID(id UserID) ListArchivedLeadsOption {
	return listArchivedLeadsOptionFunc(func(cfg *listArchivedLeadsOptions) {
		value := int(id)
		cfg.params.OwnerId = &value
	})
}

func WithArchivedLeadsPersonID(id PersonID) ListArchivedLeadsOption {
	return listArchivedLeadsOptionFunc(func(cfg *listArchivedLeadsOptions) {
		value := int(id)
		cfg.params.PersonId = &value
	})
}

func WithArchivedLeadsOrganizationID(id OrganizationID) ListArchivedLeadsOption {
	return listArchivedLeadsOptionFunc(func(cfg *listArchivedLeadsOptions) {
		value := int(id)
		cfg.params.OrganizationId = &value
	})
}

func WithArchivedLeadsFilterID(id int) ListArchivedLeadsOption {
	return listArchivedLeadsOptionFunc(func(cfg *listArchivedLeadsOptions) {
		cfg.params.FilterId = &id
	})
}

func WithArchivedLeadsSort(sort string) ListArchivedLeadsOption {
	return listArchivedLeadsOptionFunc(func(cfg *listArchivedLeadsOptions) {
		value := genv1.GetArchivedLeadsParamsSort(sort)
		cfg.params.Sort = &value
	})
}

func WithLeadTitle(title string) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.title = &title
	})
}

func WithLeadOwnerID(id UserID) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.ownerID = &id
	})
}

func WithLeadLabelIDs(ids ...LeadLabelID) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.labelIDs = append([]LeadLabelID{}, ids...)
	})
}

func WithLeadPersonID(id PersonID) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.personID = &id
	})
}

func WithLeadOrganizationID(id OrganizationID) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.orgID = &id
	})
}

func WithLeadValue(amount float64, currency string) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.value = &LeadValue{Amount: amount, Currency: currency}
	})
}

func WithLeadExpectedCloseDate(date string) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.expectedCloseDate = &date
	})
}

func WithLeadVisibleTo(value string) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.visibleTo = &value
	})
}

func WithLeadWasSeen(seen bool) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.wasSeen = &seen
	})
}

func WithLeadOriginID(originID string) CreateLeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.originID = &originID
	})
}

func WithLeadChannel(id int) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.channel = &id
	})
}

func WithLeadChannelID(id string) LeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.channelID = &id
	})
}

func WithLeadArchived(archived bool) UpdateLeadOption {
	return leadFieldOption(func(cfg *leadPayload) {
		cfg.isArchived = &archived
	})
}

func newListLeadsOptions(opts []ListLeadsOption) listLeadsOptions {
	var cfg listLeadsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListLeads(&cfg)
	}
	return cfg
}

func newListArchivedLeadsOptions(opts []ListArchivedLeadsOption) listArchivedLeadsOptions {
	var cfg listArchivedLeadsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListArchivedLeads(&cfg)
	}
	return cfg
}

func newGetLeadOptions(opts []GetLeadOption) getLeadOptions {
	var cfg getLeadOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetLead(&cfg)
	}
	return cfg
}

func newCreateLeadOptions(opts []CreateLeadOption) createLeadOptions {
	var cfg createLeadOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateLead(&cfg)
	}
	return cfg
}

func newUpdateLeadOptions(opts []UpdateLeadOption) updateLeadOptions {
	var cfg updateLeadOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateLead(&cfg)
	}
	return cfg
}

func newDeleteLeadOptions(opts []DeleteLeadOption) deleteLeadOptions {
	var cfg deleteLeadOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteLead(&cfg)
	}
	return cfg
}

func newListLeadUsersOptions(opts []ListLeadUsersOption) listLeadUsersOptions {
	var cfg listLeadUsersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListLeadUsers(&cfg)
	}
	return cfg
}

func (s *LeadsService) List(ctx context.Context, opts ...ListLeadsOption) ([]Lead, *LeadPagination, error) {
	cfg := newListLeadsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *LeadsService) ListArchived(ctx context.Context, opts ...ListArchivedLeadsOption) ([]Lead, *LeadPagination, error) {
	cfg := newListArchivedLeadsOptions(opts)
	return s.listArchived(ctx, cfg.params, cfg.requestOptions)
}

func (s *LeadsService) Get(ctx context.Context, id LeadID, opts ...GetLeadOption) (*Lead, error) {
	cfg := newGetLeadOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	leadUUID, err := parseUUID(string(id), "lead id")
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.GetLead(ctx, leadUUID, toRequestEditors(editors)...)
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
		Data *Lead `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing lead data in response")
	}
	return payload.Data, nil
}

func (s *LeadsService) Create(ctx context.Context, opts ...CreateLeadOption) (*Lead, error) {
	cfg := newCreateLeadOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.title == nil {
		return nil, fmt.Errorf("title is required")
	}
	if cfg.payload.personID == nil && cfg.payload.orgID == nil {
		return nil, fmt.Errorf("person ID or organization ID is required")
	}

	body, err := json.Marshal(cfg.payload.toMap(false))
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddLeadWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *Lead `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing lead data in response")
	}
	return payload.Data, nil
}

func (s *LeadsService) Update(ctx context.Context, id LeadID, opts ...UpdateLeadOption) (*Lead, error) {
	cfg := newUpdateLeadOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	bodyMap := cfg.payload.toMap(true)
	if len(bodyMap) == 0 {
		return nil, fmt.Errorf("at least one field is required to update")
	}

	leadUUID, err := parseUUID(string(id), "lead id")
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateLeadWithBody(ctx, leadUUID, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *Lead `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing lead data in response")
	}
	return payload.Data, nil
}

func (s *LeadsService) Delete(ctx context.Context, id LeadID, opts ...DeleteLeadOption) (*LeadDeleteResult, error) {
	cfg := newDeleteLeadOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	leadUUID, err := parseUUID(string(id), "lead id")
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.DeleteLead(ctx, leadUUID, toRequestEditors(editors)...)
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
		Data *LeadDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing lead delete data in response")
	}
	return payload.Data, nil
}

func (s *LeadsService) ListPermittedUsers(ctx context.Context, id LeadID, opts ...ListLeadUsersOption) ([]UserID, error) {
	cfg := newListLeadUsersOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetLeadUsers(ctx, string(id), toRequestEditors(editors)...)
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
		Data []UserID `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *LeadsService) list(ctx context.Context, params genv1.GetLeadsParams, requestOptions []pipedrive.RequestOption) ([]Lead, *LeadPagination, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetLeads(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp, respBody)
	}

	var payload struct {
		Data           []Lead          `json:"data"`
		AdditionalData *LeadPagination `json:"additional_data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *LeadsService) listArchived(ctx context.Context, params genv1.GetArchivedLeadsParams, requestOptions []pipedrive.RequestOption) ([]Lead, *LeadPagination, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetArchivedLeads(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp, respBody)
	}

	var payload struct {
		Data           []Lead          `json:"data"`
		AdditionalData *LeadPagination `json:"additional_data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, payload.AdditionalData, nil
}

func (p leadPayload) toMap(includeUpdateFields bool) map[string]interface{} {
	body := map[string]interface{}{}
	if p.title != nil {
		body["title"] = *p.title
	}
	if p.ownerID != nil {
		body["owner_id"] = int(*p.ownerID)
	}
	if len(p.labelIDs) > 0 {
		labels := make([]string, 0, len(p.labelIDs))
		for _, id := range p.labelIDs {
			labels = append(labels, string(id))
		}
		body["label_ids"] = labels
	}
	if p.personID != nil {
		body["person_id"] = int(*p.personID)
	}
	if p.orgID != nil {
		body["organization_id"] = int(*p.orgID)
	}
	if p.value != nil {
		body["value"] = map[string]interface{}{
			"amount":   p.value.Amount,
			"currency": p.value.Currency,
		}
	}
	if p.expectedCloseDate != nil {
		body["expected_close_date"] = *p.expectedCloseDate
	}
	if p.visibleTo != nil {
		body["visible_to"] = *p.visibleTo
	}
	if p.wasSeen != nil {
		body["was_seen"] = *p.wasSeen
	}
	if p.originID != nil {
		body["origin_id"] = *p.originID
	}
	if p.channel != nil {
		body["channel"] = *p.channel
	}
	if p.channelID != nil {
		body["channel_id"] = *p.channelID
	}
	if includeUpdateFields && p.isArchived != nil {
		body["is_archived"] = *p.isArchived
	}
	return body
}
