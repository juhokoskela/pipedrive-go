package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type CallLogOutcome string

const (
	CallLogOutcomeConnected     CallLogOutcome = "connected"
	CallLogOutcomeNoAnswer      CallLogOutcome = "no_answer"
	CallLogOutcomeLeftMessage   CallLogOutcome = "left_message"
	CallLogOutcomeLeftVoicemail CallLogOutcome = "left_voicemail"
	CallLogOutcomeWrongNumber   CallLogOutcome = "wrong_number"
	CallLogOutcomeBusy          CallLogOutcome = "busy"
)

type CallLog struct {
	ID             CallLogID       `json:"id"`
	ActivityID     *ActivityID     `json:"activity_id,omitempty"`
	Subject        string          `json:"subject,omitempty"`
	Duration       string          `json:"duration,omitempty"`
	Outcome        CallLogOutcome  `json:"outcome,omitempty"`
	FromPhone      string          `json:"from_phone_number,omitempty"`
	ToPhone        string          `json:"to_phone_number,omitempty"`
	StartTime      *DateTime       `json:"start_time,omitempty"`
	EndTime        *DateTime       `json:"end_time,omitempty"`
	PersonID       *PersonID       `json:"person_id,omitempty"`
	OrganizationID *OrganizationID `json:"org_id,omitempty"`
	DealID         *DealID         `json:"deal_id,omitempty"`
	LeadID         *LeadID         `json:"lead_id,omitempty"`
	UserID         *UserID         `json:"user_id,omitempty"`
	CompanyID      *int            `json:"company_id,omitempty"`
	HasRecording   bool            `json:"has_recording,omitempty"`
	Note           string          `json:"note,omitempty"`
}

type CallLogsPagination struct {
	Start                 int  `json:"start,omitempty"`
	Limit                 int  `json:"limit,omitempty"`
	MoreItemsInCollection bool `json:"more_items_in_collection,omitempty"`
	NextStart             *int `json:"next_start,omitempty"`
}

type CallLogsService struct {
	client *Client
}

type ListCallLogsOption interface {
	applyListCallLogs(*listCallLogsOptions)
}

type CreateCallLogOption interface {
	applyCreateCallLog(*createCallLogOptions)
}

type GetCallLogOption interface {
	applyGetCallLog(*getCallLogOptions)
}

type DeleteCallLogOption interface {
	applyDeleteCallLog(*deleteCallLogOptions)
}

type AddCallLogRecordingOption interface {
	applyAddCallLogRecording(*addCallLogRecordingOptions)
}

type CallLogsRequestOption interface {
	ListCallLogsOption
	CreateCallLogOption
	GetCallLogOption
	DeleteCallLogOption
	AddCallLogRecordingOption
}

type CallLogOption interface {
	CreateCallLogOption
}

type listCallLogsOptions struct {
	params         genv1.GetUserCallLogsParams
	requestOptions []pipedrive.RequestOption
}

type createCallLogOptions struct {
	payload        callLogPayload
	requestOptions []pipedrive.RequestOption
}

type getCallLogOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deleteCallLogOptions struct {
	requestOptions []pipedrive.RequestOption
}

type addCallLogRecordingOptions struct {
	requestOptions []pipedrive.RequestOption
}

type callLogPayload struct {
	userID     *UserID
	activityID *ActivityID
	subject    *string
	duration   *string
	outcome    *CallLogOutcome
	fromPhone  *string
	toPhone    *string
	startTime  *time.Time
	endTime    *time.Time
	personID   *PersonID
	orgID      *OrganizationID
	dealID     *DealID
	leadID     *LeadID
	note       *string
}

type callLogsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o callLogsRequestOptions) applyListCallLogs(cfg *listCallLogsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o callLogsRequestOptions) applyCreateCallLog(cfg *createCallLogOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o callLogsRequestOptions) applyGetCallLog(cfg *getCallLogOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o callLogsRequestOptions) applyDeleteCallLog(cfg *deleteCallLogOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o callLogsRequestOptions) applyAddCallLogRecording(cfg *addCallLogRecordingOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listCallLogsOptionFunc func(*listCallLogsOptions)

func (f listCallLogsOptionFunc) applyListCallLogs(cfg *listCallLogsOptions) {
	f(cfg)
}

type callLogFieldOption func(*callLogPayload)

func (f callLogFieldOption) applyCreateCallLog(cfg *createCallLogOptions) {
	f(&cfg.payload)
}

func WithCallLogsRequestOptions(opts ...pipedrive.RequestOption) CallLogsRequestOption {
	return callLogsRequestOptions{requestOptions: opts}
}

func WithCallLogsStart(start int) ListCallLogsOption {
	return listCallLogsOptionFunc(func(cfg *listCallLogsOptions) {
		cfg.params.Start = &start
	})
}

func WithCallLogsLimit(limit int) ListCallLogsOption {
	return listCallLogsOptionFunc(func(cfg *listCallLogsOptions) {
		cfg.params.Limit = &limit
	})
}

func WithCallLogUserID(id UserID) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.userID = &id
	})
}

func WithCallLogActivityID(id ActivityID) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.activityID = &id
	})
}

func WithCallLogSubject(subject string) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.subject = &subject
	})
}

func WithCallLogDuration(duration string) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.duration = &duration
	})
}

func WithCallLogOutcome(outcome CallLogOutcome) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.outcome = &outcome
	})
}

func WithCallLogFromPhoneNumber(number string) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.fromPhone = &number
	})
}

func WithCallLogToPhoneNumber(number string) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.toPhone = &number
	})
}

func WithCallLogStartTime(t time.Time) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.startTime = &t
	})
}

func WithCallLogEndTime(t time.Time) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.endTime = &t
	})
}

func WithCallLogPersonID(id PersonID) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.personID = &id
	})
}

func WithCallLogOrganizationID(id OrganizationID) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.orgID = &id
	})
}

func WithCallLogDealID(id DealID) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.dealID = &id
	})
}

func WithCallLogLeadID(id LeadID) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.leadID = &id
	})
}

func WithCallLogNote(note string) CallLogOption {
	return callLogFieldOption(func(cfg *callLogPayload) {
		cfg.note = &note
	})
}

func newListCallLogsOptions(opts []ListCallLogsOption) listCallLogsOptions {
	var cfg listCallLogsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListCallLogs(&cfg)
	}
	return cfg
}

func newCreateCallLogOptions(opts []CreateCallLogOption) createCallLogOptions {
	var cfg createCallLogOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateCallLog(&cfg)
	}
	return cfg
}

func newGetCallLogOptions(opts []GetCallLogOption) getCallLogOptions {
	var cfg getCallLogOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetCallLog(&cfg)
	}
	return cfg
}

func newDeleteCallLogOptions(opts []DeleteCallLogOption) deleteCallLogOptions {
	var cfg deleteCallLogOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteCallLog(&cfg)
	}
	return cfg
}

func newAddCallLogRecordingOptions(opts []AddCallLogRecordingOption) addCallLogRecordingOptions {
	var cfg addCallLogRecordingOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAddCallLogRecording(&cfg)
	}
	return cfg
}

func (s *CallLogsService) List(ctx context.Context, opts ...ListCallLogsOption) ([]CallLog, *CallLogsPagination, error) {
	cfg := newListCallLogsOptions(opts)
	return s.list(ctx, cfg.params, cfg.requestOptions)
}

func (s *CallLogsService) Create(ctx context.Context, opts ...CreateCallLogOption) (*CallLog, error) {
	cfg := newCreateCallLogOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.toPhone == nil || cfg.payload.outcome == nil || cfg.payload.startTime == nil || cfg.payload.endTime == nil {
		return nil, fmt.Errorf("to phone number, outcome, start time, and end time are required")
	}

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddCallLogWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *CallLog `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing call log data in response")
	}
	return payload.Data, nil
}

func (s *CallLogsService) Get(ctx context.Context, id CallLogID, opts ...GetCallLogOption) (*CallLog, error) {
	cfg := newGetCallLogOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetCallLog(ctx, string(id), toRequestEditors(editors)...)
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
		Data *CallLog `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing call log data in response")
	}
	return payload.Data, nil
}

func (s *CallLogsService) Delete(ctx context.Context, id CallLogID, opts ...DeleteCallLogOption) (bool, error) {
	cfg := newDeleteCallLogOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteCallLog(ctx, string(id), toRequestEditors(editors)...)
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
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}
	return payload.Success, nil
}

func (s *CallLogsService) AddRecording(ctx context.Context, id CallLogID, fileName string, content io.Reader, opts ...AddCallLogRecordingOption) (bool, error) {
	if fileName == "" || content == nil {
		return false, fmt.Errorf("file name and content are required")
	}
	cfg := newAddCallLogRecordingOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return false, fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, content); err != nil {
		return false, fmt.Errorf("write form file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return false, fmt.Errorf("close multipart writer: %w", err)
	}

	resp, err := s.client.gen.AddCallLogAudioFileWithBody(ctx, string(id), writer.FormDataContentType(), &buf, toRequestEditors(editors)...)
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
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}
	return payload.Success, nil
}

func (s *CallLogsService) list(ctx context.Context, params genv1.GetUserCallLogsParams, requestOptions []pipedrive.RequestOption) ([]CallLog, *CallLogsPagination, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetUserCallLogs(ctx, &params, toRequestEditors(editors)...)
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
		Data           []CallLog `json:"data"`
		AdditionalData *struct {
			Pagination *CallLogsPagination `json:"pagination"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var page *CallLogsPagination
	if payload.AdditionalData != nil {
		page = payload.AdditionalData.Pagination
	}
	return payload.Data, page, nil
}

func (p callLogPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.userID != nil {
		body["user_id"] = int(*p.userID)
	}
	if p.activityID != nil {
		body["activity_id"] = int(*p.activityID)
	}
	if p.subject != nil {
		body["subject"] = *p.subject
	}
	if p.duration != nil {
		body["duration"] = *p.duration
	}
	if p.outcome != nil {
		body["outcome"] = string(*p.outcome)
	}
	if p.fromPhone != nil {
		body["from_phone_number"] = *p.fromPhone
	}
	if p.toPhone != nil {
		body["to_phone_number"] = *p.toPhone
	}
	if p.startTime != nil {
		body["start_time"] = formatV1Time(*p.startTime)
	}
	if p.endTime != nil {
		body["end_time"] = formatV1Time(*p.endTime)
	}
	if p.personID != nil {
		body["person_id"] = int(*p.personID)
	}
	if p.orgID != nil {
		body["org_id"] = int(*p.orgID)
	}
	if p.dealID != nil {
		body["deal_id"] = int(*p.dealID)
	}
	if p.leadID != nil {
		body["lead_id"] = string(*p.leadID)
	}
	if p.note != nil {
		body["note"] = *p.note
	}
	return body
}
