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

type NoteDeal struct {
	Title string `json:"title,omitempty"`
}

type NoteOrganization struct {
	Name string `json:"name,omitempty"`
}

type NotePerson struct {
	Name string `json:"name,omitempty"`
}

type NoteProject struct {
	Title string `json:"title,omitempty"`
}

type NoteUser struct {
	Email   string  `json:"email,omitempty"`
	IconURL *string `json:"icon_url,omitempty"`
	IsYou   bool    `json:"is_you,omitempty"`
	Name    string  `json:"name,omitempty"`
}

type Note struct {
	ID                   NoteID            `json:"id,omitempty"`
	Active               bool              `json:"active_flag,omitempty"`
	AddTime              *DateTime         `json:"add_time,omitempty"`
	Content              string            `json:"content,omitempty"`
	Deal                 *NoteDeal         `json:"deal,omitempty"`
	LeadID               *LeadID           `json:"lead_id,omitempty"`
	DealID               *DealID           `json:"deal_id,omitempty"`
	LastUpdateUserID     *UserID           `json:"last_update_user_id,omitempty"`
	OrgID                *OrganizationID   `json:"org_id,omitempty"`
	Organization         *NoteOrganization `json:"organization,omitempty"`
	Person               *NotePerson       `json:"person,omitempty"`
	PersonID             *PersonID         `json:"person_id,omitempty"`
	ProjectID            *ProjectID        `json:"project_id,omitempty"`
	Project              *NoteProject      `json:"project,omitempty"`
	PinnedToLead         bool              `json:"pinned_to_lead_flag,omitempty"`
	PinnedToDeal         bool              `json:"pinned_to_deal_flag,omitempty"`
	PinnedToOrganization bool              `json:"pinned_to_organization_flag,omitempty"`
	PinnedToPerson       bool              `json:"pinned_to_person_flag,omitempty"`
	PinnedToProject      bool              `json:"pinned_to_project_flag,omitempty"`
	UpdateTime           *DateTime         `json:"update_time,omitempty"`
	User                 *NoteUser         `json:"user,omitempty"`
	UserID               *UserID           `json:"user_id,omitempty"`
}

type NotesPagination struct {
	NextStart             int  `json:"next_start,omitempty"`
	Start                 int  `json:"start,omitempty"`
	Limit                 int  `json:"limit,omitempty"`
	MoreItemsInCollection bool `json:"more_items_in_collection,omitempty"`
}

type NotesAdditionalData struct {
	Pagination *NotesPagination `json:"pagination,omitempty"`
}

type NoteComment struct {
	ID         CommentID `json:"uuid,omitempty"`
	Active     bool      `json:"active_flag,omitempty"`
	AddTime    *DateTime `json:"add_time,omitempty"`
	UpdateTime *DateTime `json:"update_time,omitempty"`
	Content    string    `json:"content,omitempty"`
	ObjectID   string    `json:"object_id,omitempty"`
	ObjectType string    `json:"object_type,omitempty"`
	UserID     *UserID   `json:"user_id,omitempty"`
	UpdaterID  *UserID   `json:"updater_id,omitempty"`
	CompanyID  *int      `json:"company_id,omitempty"`
}

type NoteCommentsAdditionalData struct {
	Pagination *NotesPagination `json:"pagination,omitempty"`
}

type NotesService struct {
	client *Client
}

type ListNotesOption interface {
	applyListNotes(*listNotesOptions)
}

type GetNoteOption interface {
	applyGetNote(*getNoteOptions)
}

type CreateNoteOption interface {
	applyCreateNote(*createNoteOptions)
}

type UpdateNoteOption interface {
	applyUpdateNote(*updateNoteOptions)
}

type DeleteNoteOption interface {
	applyDeleteNote(*deleteNoteOptions)
}

type ListNoteCommentsOption interface {
	applyListNoteComments(*listNoteCommentsOptions)
}

type CreateNoteCommentOption interface {
	applyCreateNoteComment(*createNoteCommentOptions)
}

type GetNoteCommentOption interface {
	applyGetNoteComment(*getNoteCommentOptions)
}

type UpdateNoteCommentOption interface {
	applyUpdateNoteComment(*updateNoteCommentOptions)
}

type DeleteNoteCommentOption interface {
	applyDeleteNoteComment(*deleteNoteCommentOptions)
}

type NotesRequestOption interface {
	ListNotesOption
	GetNoteOption
	CreateNoteOption
	UpdateNoteOption
	DeleteNoteOption
	ListNoteCommentsOption
	CreateNoteCommentOption
	GetNoteCommentOption
	UpdateNoteCommentOption
	DeleteNoteCommentOption
}

type NoteOption interface {
	CreateNoteOption
	UpdateNoteOption
}

type NoteCommentOption interface {
	CreateNoteCommentOption
	UpdateNoteCommentOption
}

type listNotesOptions struct {
	params         genv1.GetNotesParams
	leadID         *LeadID
	requestOptions []pipedrive.RequestOption
}

type getNoteOptions struct {
	requestOptions []pipedrive.RequestOption
}

type createNoteOptions struct {
	payload        notePayload
	requestOptions []pipedrive.RequestOption
}

type updateNoteOptions struct {
	payload        notePayload
	requestOptions []pipedrive.RequestOption
}

type deleteNoteOptions struct {
	requestOptions []pipedrive.RequestOption
}

type listNoteCommentsOptions struct {
	params         genv1.GetNoteCommentsParams
	requestOptions []pipedrive.RequestOption
}

type createNoteCommentOptions struct {
	payload        noteCommentPayload
	requestOptions []pipedrive.RequestOption
}

type getNoteCommentOptions struct {
	requestOptions []pipedrive.RequestOption
}

type updateNoteCommentOptions struct {
	payload        noteCommentPayload
	requestOptions []pipedrive.RequestOption
}

type deleteNoteCommentOptions struct {
	requestOptions []pipedrive.RequestOption
}

type notePayload struct {
	content              *string
	leadID               *LeadID
	dealID               *DealID
	personID             *PersonID
	orgID                *OrganizationID
	projectID            *ProjectID
	userID               *UserID
	addTime              *time.Time
	pinnedToLead         *bool
	pinnedToDeal         *bool
	pinnedToOrganization *bool
	pinnedToPerson       *bool
	pinnedToProject      *bool
}

type noteCommentPayload struct {
	content *string
}

type notesRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o notesRequestOptions) applyListNotes(cfg *listNotesOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o notesRequestOptions) applyGetNote(cfg *getNoteOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o notesRequestOptions) applyCreateNote(cfg *createNoteOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o notesRequestOptions) applyUpdateNote(cfg *updateNoteOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o notesRequestOptions) applyDeleteNote(cfg *deleteNoteOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o notesRequestOptions) applyListNoteComments(cfg *listNoteCommentsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o notesRequestOptions) applyCreateNoteComment(cfg *createNoteCommentOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o notesRequestOptions) applyGetNoteComment(cfg *getNoteCommentOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o notesRequestOptions) applyUpdateNoteComment(cfg *updateNoteCommentOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o notesRequestOptions) applyDeleteNoteComment(cfg *deleteNoteCommentOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listNotesOptionFunc func(*listNotesOptions)

func (f listNotesOptionFunc) applyListNotes(cfg *listNotesOptions) {
	f(cfg)
}

type getNoteOptionFunc func(*getNoteOptions)

func (f getNoteOptionFunc) applyGetNote(cfg *getNoteOptions) {
	f(cfg)
}

type noteFieldOption func(*notePayload)

func (f noteFieldOption) applyCreateNote(cfg *createNoteOptions) {
	f(&cfg.payload)
}

func (f noteFieldOption) applyUpdateNote(cfg *updateNoteOptions) {
	f(&cfg.payload)
}

type deleteNoteOptionFunc func(*deleteNoteOptions)

func (f deleteNoteOptionFunc) applyDeleteNote(cfg *deleteNoteOptions) {
	f(cfg)
}

type listNoteCommentsOptionFunc func(*listNoteCommentsOptions)

func (f listNoteCommentsOptionFunc) applyListNoteComments(cfg *listNoteCommentsOptions) {
	f(cfg)
}

type noteCommentFieldOption func(*noteCommentPayload)

func (f noteCommentFieldOption) applyCreateNoteComment(cfg *createNoteCommentOptions) {
	f(&cfg.payload)
}

func (f noteCommentFieldOption) applyUpdateNoteComment(cfg *updateNoteCommentOptions) {
	f(&cfg.payload)
}

type getNoteCommentOptionFunc func(*getNoteCommentOptions)

func (f getNoteCommentOptionFunc) applyGetNoteComment(cfg *getNoteCommentOptions) {
	f(cfg)
}

type deleteNoteCommentOptionFunc func(*deleteNoteCommentOptions)

func (f deleteNoteCommentOptionFunc) applyDeleteNoteComment(cfg *deleteNoteCommentOptions) {
	f(cfg)
}

func WithNotesRequestOptions(opts ...pipedrive.RequestOption) NotesRequestOption {
	return notesRequestOptions{requestOptions: opts}
}

func WithNotesUserID(id UserID) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := int(id)
		cfg.params.UserId = &value
	})
}

func WithNotesLeadID(id LeadID) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		cfg.leadID = &id
	})
}

func WithNotesDealID(id DealID) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := int(id)
		cfg.params.DealId = &value
	})
}

func WithNotesPersonID(id PersonID) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := int(id)
		cfg.params.PersonId = &value
	})
}

func WithNotesOrganizationID(id OrganizationID) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := int(id)
		cfg.params.OrgId = &value
	})
}

func WithNotesProjectID(id ProjectID) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := int(id)
		cfg.params.ProjectId = &value
	})
}

func WithNotesStart(start int) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		cfg.params.Start = &start
	})
}

func WithNotesLimit(limit int) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		cfg.params.Limit = &limit
	})
}

func WithNotesSort(sort string) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		cfg.params.Sort = &sort
	})
}

func WithNotesStartDate(start time.Time) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		cfg.params.StartDate = &openapi_types.Date{Time: start}
	})
}

func WithNotesEndDate(end time.Time) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		cfg.params.EndDate = &openapi_types.Date{Time: end}
	})
}

func WithNotesPinnedToLead(flag bool) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := genv1.GetNotesParamsPinnedToLeadFlag(boolToNumber(flag))
		cfg.params.PinnedToLeadFlag = &value
	})
}

func WithNotesPinnedToDeal(flag bool) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := genv1.GetNotesParamsPinnedToDealFlag(boolToNumber(flag))
		cfg.params.PinnedToDealFlag = &value
	})
}

func WithNotesPinnedToOrganization(flag bool) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := genv1.GetNotesParamsPinnedToOrganizationFlag(boolToNumber(flag))
		cfg.params.PinnedToOrganizationFlag = &value
	})
}

func WithNotesPinnedToPerson(flag bool) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := genv1.GetNotesParamsPinnedToPersonFlag(boolToNumber(flag))
		cfg.params.PinnedToPersonFlag = &value
	})
}

func WithNotesPinnedToProject(flag bool) ListNotesOption {
	return listNotesOptionFunc(func(cfg *listNotesOptions) {
		value := genv1.GetNotesParamsPinnedToProjectFlag(boolToNumber(flag))
		cfg.params.PinnedToProjectFlag = &value
	})
}

func WithNoteContent(content string) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.content = &content
	})
}

func WithNoteLeadID(id LeadID) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.leadID = &id
	})
}

func WithNoteDealID(id DealID) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.dealID = &id
	})
}

func WithNotePersonID(id PersonID) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.personID = &id
	})
}

func WithNoteOrganizationID(id OrganizationID) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.orgID = &id
	})
}

func WithNoteProjectID(id ProjectID) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.projectID = &id
	})
}

func WithNoteUserID(id UserID) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.userID = &id
	})
}

func WithNoteAddTime(addTime time.Time) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.addTime = &addTime
	})
}

func WithNotePinnedToLead(flag bool) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.pinnedToLead = &flag
	})
}

func WithNotePinnedToDeal(flag bool) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.pinnedToDeal = &flag
	})
}

func WithNotePinnedToOrganization(flag bool) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.pinnedToOrganization = &flag
	})
}

func WithNotePinnedToPerson(flag bool) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.pinnedToPerson = &flag
	})
}

func WithNotePinnedToProject(flag bool) NoteOption {
	return noteFieldOption(func(cfg *notePayload) {
		cfg.pinnedToProject = &flag
	})
}

func WithNoteCommentsStart(start int) ListNoteCommentsOption {
	return listNoteCommentsOptionFunc(func(cfg *listNoteCommentsOptions) {
		cfg.params.Start = &start
	})
}

func WithNoteCommentsLimit(limit int) ListNoteCommentsOption {
	return listNoteCommentsOptionFunc(func(cfg *listNoteCommentsOptions) {
		cfg.params.Limit = &limit
	})
}

func WithNoteCommentContent(content string) NoteCommentOption {
	return noteCommentFieldOption(func(cfg *noteCommentPayload) {
		cfg.content = &content
	})
}

func newListNotesOptions(opts []ListNotesOption) listNotesOptions {
	var cfg listNotesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListNotes(&cfg)
	}
	return cfg
}

func newGetNoteOptions(opts []GetNoteOption) getNoteOptions {
	var cfg getNoteOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetNote(&cfg)
	}
	return cfg
}

func newCreateNoteOptions(opts []CreateNoteOption) createNoteOptions {
	var cfg createNoteOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateNote(&cfg)
	}
	return cfg
}

func newUpdateNoteOptions(opts []UpdateNoteOption) updateNoteOptions {
	var cfg updateNoteOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateNote(&cfg)
	}
	return cfg
}

func newDeleteNoteOptions(opts []DeleteNoteOption) deleteNoteOptions {
	var cfg deleteNoteOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteNote(&cfg)
	}
	return cfg
}

func newListNoteCommentsOptions(opts []ListNoteCommentsOption) listNoteCommentsOptions {
	var cfg listNoteCommentsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListNoteComments(&cfg)
	}
	return cfg
}

func newCreateNoteCommentOptions(opts []CreateNoteCommentOption) createNoteCommentOptions {
	var cfg createNoteCommentOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateNoteComment(&cfg)
	}
	return cfg
}

func newGetNoteCommentOptions(opts []GetNoteCommentOption) getNoteCommentOptions {
	var cfg getNoteCommentOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetNoteComment(&cfg)
	}
	return cfg
}

func newUpdateNoteCommentOptions(opts []UpdateNoteCommentOption) updateNoteCommentOptions {
	var cfg updateNoteCommentOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateNoteComment(&cfg)
	}
	return cfg
}

func newDeleteNoteCommentOptions(opts []DeleteNoteCommentOption) deleteNoteCommentOptions {
	var cfg deleteNoteCommentOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteNoteComment(&cfg)
	}
	return cfg
}

func (s *NotesService) List(ctx context.Context, opts ...ListNotesOption) ([]Note, *NotesAdditionalData, error) {
	cfg := newListNotesOptions(opts)
	if cfg.leadID != nil {
		leadUUID, err := parseUUID(string(*cfg.leadID), "lead id")
		if err != nil {
			return nil, nil, err
		}
		cfg.params.LeadId = &leadUUID
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetNotes(ctx, &cfg.params, toRequestEditors(editors)...)
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
		Data           []Note               `json:"data"`
		AdditionalData *NotesAdditionalData `json:"additional_data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *NotesService) Get(ctx context.Context, id NoteID, opts ...GetNoteOption) (*Note, error) {
	cfg := newGetNoteOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetNote(ctx, int(id), toRequestEditors(editors)...)
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
		Data *Note `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing note data in response")
	}
	return payload.Data, nil
}

func (s *NotesService) Create(ctx context.Context, opts ...CreateNoteOption) (*Note, error) {
	cfg := newCreateNoteOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.content == nil {
		return nil, fmt.Errorf("content is required")
	}
	if !cfg.payload.hasTarget() {
		return nil, fmt.Errorf("note target is required")
	}
	bodyMap, err := cfg.payload.toMap()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddNoteWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *Note `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing note data in response")
	}
	return payload.Data, nil
}

func (s *NotesService) Update(ctx context.Context, id NoteID, opts ...UpdateNoteOption) (*Note, error) {
	cfg := newUpdateNoteOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	bodyMap, err := cfg.payload.toMap()
	if err != nil {
		return nil, err
	}
	if len(bodyMap) == 0 {
		return nil, fmt.Errorf("no note fields provided")
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateNoteWithBody(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *Note `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing note data in response")
	}
	return payload.Data, nil
}

func (s *NotesService) Delete(ctx context.Context, id NoteID, opts ...DeleteNoteOption) (bool, error) {
	cfg := newDeleteNoteOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteNote(ctx, int(id), toRequestEditors(editors)...)
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
		Data bool `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *NotesService) ListComments(ctx context.Context, id NoteID, opts ...ListNoteCommentsOption) ([]NoteComment, *NoteCommentsAdditionalData, error) {
	cfg := newListNoteCommentsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetNoteComments(ctx, int(id), &cfg.params, toRequestEditors(editors)...)
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
		Data           []NoteComment               `json:"data"`
		AdditionalData *NoteCommentsAdditionalData `json:"additional_data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *NotesService) CreateComment(ctx context.Context, id NoteID, opts ...CreateNoteCommentOption) (*NoteComment, error) {
	cfg := newCreateNoteCommentOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.content == nil {
		return nil, fmt.Errorf("content is required")
	}

	bodyMap := map[string]interface{}{
		"content": *cfg.payload.content,
	}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddNoteCommentWithBody(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *NoteComment `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing comment data in response")
	}
	return payload.Data, nil
}

func (s *NotesService) GetComment(ctx context.Context, id NoteID, commentID CommentID, opts ...GetNoteCommentOption) (*NoteComment, error) {
	cfg := newGetNoteCommentOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	commentUUID, err := parseUUID(string(commentID), "comment id")
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.GetComment(ctx, int(id), commentUUID, toRequestEditors(editors)...)
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
		Data *NoteComment `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing comment data in response")
	}
	return payload.Data, nil
}

func (s *NotesService) UpdateComment(ctx context.Context, id NoteID, commentID CommentID, opts ...UpdateNoteCommentOption) (*NoteComment, error) {
	cfg := newUpdateNoteCommentOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.content == nil {
		return nil, fmt.Errorf("content is required")
	}
	commentUUID, err := parseUUID(string(commentID), "comment id")
	if err != nil {
		return nil, err
	}

	bodyMap := map[string]interface{}{
		"content": *cfg.payload.content,
	}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateCommentForNoteWithBody(ctx, int(id), commentUUID, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *NoteComment `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing comment data in response")
	}
	return payload.Data, nil
}

func (s *NotesService) DeleteComment(ctx context.Context, id NoteID, commentID CommentID, opts ...DeleteNoteCommentOption) (bool, error) {
	cfg := newDeleteNoteCommentOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	commentUUID, err := parseUUID(string(commentID), "comment id")
	if err != nil {
		return false, err
	}

	resp, err := s.client.gen.DeleteComment(ctx, int(id), commentUUID, toRequestEditors(editors)...)
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
		Data bool `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (p notePayload) hasTarget() bool {
	return p.leadID != nil || p.dealID != nil || p.personID != nil || p.orgID != nil || p.projectID != nil
}

func (p notePayload) toMap() (map[string]interface{}, error) {
	body := map[string]interface{}{}
	if p.content != nil {
		body["content"] = *p.content
	}
	if p.leadID != nil {
		if _, err := parseUUID(string(*p.leadID), "lead id"); err != nil {
			return nil, err
		}
		body["lead_id"] = string(*p.leadID)
	}
	if p.dealID != nil {
		body["deal_id"] = int64(*p.dealID)
	}
	if p.personID != nil {
		body["person_id"] = int64(*p.personID)
	}
	if p.orgID != nil {
		body["org_id"] = int64(*p.orgID)
	}
	if p.projectID != nil {
		body["project_id"] = int64(*p.projectID)
	}
	if p.userID != nil {
		body["user_id"] = int64(*p.userID)
	}
	if p.addTime != nil {
		body["add_time"] = formatV1Time(*p.addTime)
	}
	if p.pinnedToLead != nil {
		if p.leadID == nil {
			return nil, fmt.Errorf("lead ID is required for pinned to lead flag")
		}
		body["pinned_to_lead_flag"] = boolToNumber(*p.pinnedToLead)
	}
	if p.pinnedToDeal != nil {
		if p.dealID == nil {
			return nil, fmt.Errorf("deal ID is required for pinned to deal flag")
		}
		body["pinned_to_deal_flag"] = boolToNumber(*p.pinnedToDeal)
	}
	if p.pinnedToOrganization != nil {
		if p.orgID == nil {
			return nil, fmt.Errorf("organization ID is required for pinned to organization flag")
		}
		body["pinned_to_organization_flag"] = boolToNumber(*p.pinnedToOrganization)
	}
	if p.pinnedToPerson != nil {
		if p.personID == nil {
			return nil, fmt.Errorf("person ID is required for pinned to person flag")
		}
		body["pinned_to_person_flag"] = boolToNumber(*p.pinnedToPerson)
	}
	if p.pinnedToProject != nil {
		if p.projectID == nil {
			return nil, fmt.Errorf("project ID is required for pinned to project flag")
		}
		body["pinned_to_project_flag"] = boolToNumber(*p.pinnedToProject)
	}
	return body, nil
}

func boolToNumber(value bool) float32 {
	if value {
		return 1
	}
	return 0
}
