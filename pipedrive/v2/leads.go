package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type LeadSearchField string

const (
	LeadSearchFieldCustomFields LeadSearchField = "custom_fields"
	LeadSearchFieldNotes        LeadSearchField = "notes"
	LeadSearchFieldTitle        LeadSearchField = "title"
)

type LeadSearchItem struct {
	ResultScore float64                `json:"result_score,omitempty"`
	Item        map[string]interface{} `json:"item,omitempty"`
}

type LeadSearchResults struct {
	Items []LeadSearchItem `json:"items,omitempty"`
}

type LeadConversionJob struct {
	ConversionID ConversionID `json:"conversion_id"`
}

type LeadConversionStatus struct {
	LeadID       *LeadID          `json:"lead_id,omitempty"`
	DealID       *DealID          `json:"deal_id,omitempty"`
	ConversionID ConversionID     `json:"conversion_id,omitempty"`
	Status       ConversionStatus `json:"status,omitempty"`
}

type LeadsService struct {
	client *Client
}

type SearchLeadsOption interface {
	applySearchLeads(*searchLeadsOptions)
}

type ConvertLeadOption interface {
	applyConvertLead(*convertLeadOptions)
}

type GetLeadConversionStatusOption interface {
	applyGetLeadConversionStatus(*getLeadConversionStatusOptions)
}

type LeadRequestOption interface {
	SearchLeadsOption
	ConvertLeadOption
	GetLeadConversionStatusOption
}

type searchLeadsOptions struct {
	params         genv2.SearchLeadsParams
	requestOptions []pipedrive.RequestOption
}

type convertLeadOptions struct {
	payload        leadConversionPayload
	requestOptions []pipedrive.RequestOption
}

type getLeadConversionStatusOptions struct {
	requestOptions []pipedrive.RequestOption
}

type leadConversionPayload struct {
	stageID    *StageID
	pipelineID *PipelineID
}

type leadRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o leadRequestOptions) applySearchLeads(cfg *searchLeadsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadRequestOptions) applyConvertLead(cfg *convertLeadOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadRequestOptions) applyGetLeadConversionStatus(cfg *getLeadConversionStatusOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type searchLeadsOptionFunc func(*searchLeadsOptions)

func (f searchLeadsOptionFunc) applySearchLeads(cfg *searchLeadsOptions) {
	f(cfg)
}

type leadConversionOptionFunc func(*leadConversionPayload)

func (f leadConversionOptionFunc) applyConvertLead(cfg *convertLeadOptions) {
	f(&cfg.payload)
}

func WithLeadRequestOptions(opts ...pipedrive.RequestOption) LeadRequestOption {
	return leadRequestOptions{requestOptions: opts}
}

func WithLeadSearchFields(fields ...LeadSearchField) SearchLeadsOption {
	return searchLeadsOptionFunc(func(cfg *searchLeadsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.SearchLeadsParamsFields(csv)
		cfg.params.Fields = &value
	})
}

func WithLeadSearchExactMatch(enabled bool) SearchLeadsOption {
	return searchLeadsOptionFunc(func(cfg *searchLeadsOptions) {
		cfg.params.ExactMatch = &enabled
	})
}

func WithLeadSearchPageSize(limit int) SearchLeadsOption {
	return searchLeadsOptionFunc(func(cfg *searchLeadsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithLeadSearchCursor(cursor string) SearchLeadsOption {
	return searchLeadsOptionFunc(func(cfg *searchLeadsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithLeadConversionStageID(id StageID) ConvertLeadOption {
	return leadConversionOptionFunc(func(payload *leadConversionPayload) {
		payload.stageID = &id
	})
}

func WithLeadConversionPipelineID(id PipelineID) ConvertLeadOption {
	return leadConversionOptionFunc(func(payload *leadConversionPayload) {
		payload.pipelineID = &id
	})
}

func newSearchLeadsOptions(opts []SearchLeadsOption) searchLeadsOptions {
	var cfg searchLeadsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applySearchLeads(&cfg)
	}
	return cfg
}

func newConvertLeadOptions(opts []ConvertLeadOption) convertLeadOptions {
	var cfg convertLeadOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyConvertLead(&cfg)
	}
	return cfg
}

func newGetLeadConversionStatusOptions(opts []GetLeadConversionStatusOption) getLeadConversionStatusOptions {
	var cfg getLeadConversionStatusOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetLeadConversionStatus(&cfg)
	}
	return cfg
}

func (s *LeadsService) Search(ctx context.Context, term string, opts ...SearchLeadsOption) (*LeadSearchResults, *string, error) {
	cfg := newSearchLeadsOptions(opts)
	cfg.params.Term = term
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.SearchLeadsWithResponse(ctx, &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           *LeadSearchResults `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, nil, fmt.Errorf("missing lead search data in response")
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *LeadsService) ConvertToDeal(ctx context.Context, id LeadID, opts ...ConvertLeadOption) (*LeadConversionJob, error) {
	cfg := newConvertLeadOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	leadUUID, err := parseUUID(string(id), "lead id")
	if err != nil {
		return nil, err
	}

	payload := cfg.payload.toMap()
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.ConvertLeadToDealWithBodyWithResponse(ctx, leadUUID, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payloadResp struct {
		Data *LeadConversionJob `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payloadResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payloadResp.Data == nil {
		return nil, fmt.Errorf("missing lead conversion data in response")
	}
	return payloadResp.Data, nil
}

func (s *LeadsService) ConversionStatus(ctx context.Context, id LeadID, conversionID ConversionID, opts ...GetLeadConversionStatusOption) (*LeadConversionStatus, error) {
	cfg := newGetLeadConversionStatusOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	leadUUID, err := parseUUID(string(id), "lead id")
	if err != nil {
		return nil, err
	}
	conversionUUID, err := parseUUID(string(conversionID), "conversion id")
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.GetLeadConversionStatusWithResponse(ctx, leadUUID, conversionUUID, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payloadResp struct {
		Data *LeadConversionStatus `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payloadResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payloadResp.Data == nil {
		return nil, fmt.Errorf("missing lead conversion status data in response")
	}
	return payloadResp.Data, nil
}

func (p leadConversionPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.stageID != nil {
		body["stage_id"] = int(*p.stageID)
	}
	if p.pipelineID != nil {
		body["pipeline_id"] = int(*p.pipelineID)
	}
	return body
}

func parseUUID(value string, label string) (openapi_types.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return openapi_types.UUID{}, fmt.Errorf("parse %s: %w", label, err)
	}
	return openapi_types.UUID(parsed), nil
}
