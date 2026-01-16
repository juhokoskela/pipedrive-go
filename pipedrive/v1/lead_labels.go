package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type LeadLabelColor string

const (
	LeadLabelColorBlue     LeadLabelColor = "blue"
	LeadLabelColorBrown    LeadLabelColor = "brown"
	LeadLabelColorDarkGray LeadLabelColor = "dark-gray"
	LeadLabelColorGray     LeadLabelColor = "gray"
	LeadLabelColorGreen    LeadLabelColor = "green"
	LeadLabelColorOrange   LeadLabelColor = "orange"
	LeadLabelColorPink     LeadLabelColor = "pink"
	LeadLabelColorPurple   LeadLabelColor = "purple"
	LeadLabelColorRed      LeadLabelColor = "red"
	LeadLabelColorYellow   LeadLabelColor = "yellow"
)

type LeadLabel struct {
	ID         LeadLabelID    `json:"id"`
	Name       string         `json:"name,omitempty"`
	Color      LeadLabelColor `json:"color,omitempty"`
	AddTime    *DateTime      `json:"add_time,omitempty"`
	UpdateTime *DateTime      `json:"update_time,omitempty"`
}

type LeadLabelDeleteResult struct {
	ID LeadLabelID `json:"id"`
}

type LeadLabelsService struct {
	client *Client
}

type ListLeadLabelsOption interface {
	applyListLeadLabels(*listLeadLabelsOptions)
}

type CreateLeadLabelOption interface {
	applyCreateLeadLabel(*createLeadLabelOptions)
}

type UpdateLeadLabelOption interface {
	applyUpdateLeadLabel(*updateLeadLabelOptions)
}

type DeleteLeadLabelOption interface {
	applyDeleteLeadLabel(*deleteLeadLabelOptions)
}

type LeadLabelsRequestOption interface {
	ListLeadLabelsOption
	CreateLeadLabelOption
	UpdateLeadLabelOption
	DeleteLeadLabelOption
}

type LeadLabelOption interface {
	CreateLeadLabelOption
	UpdateLeadLabelOption
}

type listLeadLabelsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type createLeadLabelOptions struct {
	payload        leadLabelPayload
	requestOptions []pipedrive.RequestOption
}

type updateLeadLabelOptions struct {
	payload        leadLabelPayload
	requestOptions []pipedrive.RequestOption
}

type deleteLeadLabelOptions struct {
	requestOptions []pipedrive.RequestOption
}

type leadLabelPayload struct {
	name  *string
	color *LeadLabelColor
}

type leadLabelsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o leadLabelsRequestOptions) applyListLeadLabels(cfg *listLeadLabelsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadLabelsRequestOptions) applyCreateLeadLabel(cfg *createLeadLabelOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadLabelsRequestOptions) applyUpdateLeadLabel(cfg *updateLeadLabelOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o leadLabelsRequestOptions) applyDeleteLeadLabel(cfg *deleteLeadLabelOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listLeadLabelsOptionFunc func(*listLeadLabelsOptions)

func (f listLeadLabelsOptionFunc) applyListLeadLabels(cfg *listLeadLabelsOptions) {
	f(cfg)
}

type leadLabelFieldOption func(*leadLabelPayload)

func (f leadLabelFieldOption) applyCreateLeadLabel(cfg *createLeadLabelOptions) {
	f(&cfg.payload)
}

func (f leadLabelFieldOption) applyUpdateLeadLabel(cfg *updateLeadLabelOptions) {
	f(&cfg.payload)
}

func WithLeadLabelsRequestOptions(opts ...pipedrive.RequestOption) LeadLabelsRequestOption {
	return leadLabelsRequestOptions{requestOptions: opts}
}

func WithLeadLabelName(name string) LeadLabelOption {
	return leadLabelFieldOption(func(cfg *leadLabelPayload) {
		cfg.name = &name
	})
}

func WithLeadLabelColor(color LeadLabelColor) LeadLabelOption {
	return leadLabelFieldOption(func(cfg *leadLabelPayload) {
		cfg.color = &color
	})
}

func newListLeadLabelsOptions(opts []ListLeadLabelsOption) listLeadLabelsOptions {
	var cfg listLeadLabelsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListLeadLabels(&cfg)
	}
	return cfg
}

func newCreateLeadLabelOptions(opts []CreateLeadLabelOption) createLeadLabelOptions {
	var cfg createLeadLabelOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateLeadLabel(&cfg)
	}
	return cfg
}

func newUpdateLeadLabelOptions(opts []UpdateLeadLabelOption) updateLeadLabelOptions {
	var cfg updateLeadLabelOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateLeadLabel(&cfg)
	}
	return cfg
}

func newDeleteLeadLabelOptions(opts []DeleteLeadLabelOption) deleteLeadLabelOptions {
	var cfg deleteLeadLabelOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteLeadLabel(&cfg)
	}
	return cfg
}

func (s *LeadLabelsService) List(ctx context.Context, opts ...ListLeadLabelsOption) ([]LeadLabel, error) {
	cfg := newListLeadLabelsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetLeadLabels(ctx, toRequestEditors(editors)...)
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
		Data []LeadLabel `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *LeadLabelsService) Create(ctx context.Context, opts ...CreateLeadLabelOption) (*LeadLabel, error) {
	cfg := newCreateLeadLabelOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.name == nil || cfg.payload.color == nil {
		return nil, fmt.Errorf("name and color are required")
	}

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddLeadLabelWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *LeadLabel `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing lead label data in response")
	}
	return payload.Data, nil
}

func (s *LeadLabelsService) Update(ctx context.Context, id LeadLabelID, opts ...UpdateLeadLabelOption) (*LeadLabel, error) {
	cfg := newUpdateLeadLabelOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	bodyMap := cfg.payload.toMap()
	if len(bodyMap) == 0 {
		return nil, fmt.Errorf("at least one field is required to update")
	}

	labelUUID, err := parseUUID(string(id), "lead label id")
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateLeadLabelWithBody(ctx, labelUUID, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *LeadLabel `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing lead label data in response")
	}
	return payload.Data, nil
}

func (s *LeadLabelsService) Delete(ctx context.Context, id LeadLabelID, opts ...DeleteLeadLabelOption) (*LeadLabelDeleteResult, error) {
	cfg := newDeleteLeadLabelOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	labelUUID, err := parseUUID(string(id), "lead label id")
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.DeleteLeadLabel(ctx, labelUUID, toRequestEditors(editors)...)
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
		Data *LeadLabelDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing lead label delete data in response")
	}
	return payload.Data, nil
}

func (p leadLabelPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if p.color != nil {
		body["color"] = string(*p.color)
	}
	return body
}
