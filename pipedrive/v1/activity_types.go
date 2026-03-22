package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ActivityType struct {
	ID         ActivityTypeID `json:"id"`
	Name       string         `json:"name,omitempty"`
	IconKey    string         `json:"icon_key,omitempty"`
	Color      string         `json:"color,omitempty"`
	Order      int            `json:"order_nr,omitempty"`
	Key        string         `json:"key_string,omitempty"`
	Active     bool           `json:"active_flag,omitempty"`
	IsCustom   bool           `json:"is_custom_flag,omitempty"`
	AddTime    *DateTime      `json:"add_time,omitempty"`
	UpdateTime *DateTime      `json:"update_time,omitempty"`
}

type ActivityTypesService struct {
	client *Client
}

type ListActivityTypesOption interface {
	applyListActivityTypes(*listActivityTypesOptions)
}

type CreateActivityTypeOption interface {
	applyCreateActivityType(*createActivityTypeOptions)
}

type UpdateActivityTypeOption interface {
	applyUpdateActivityType(*updateActivityTypeOptions)
}

type DeleteActivityTypeOption interface {
	applyDeleteActivityType(*deleteActivityTypeOptions)
}

type ActivityTypesRequestOption interface {
	ListActivityTypesOption
	CreateActivityTypeOption
	UpdateActivityTypeOption
	DeleteActivityTypeOption
}

type ActivityTypeOption interface {
	CreateActivityTypeOption
	UpdateActivityTypeOption
}

type listActivityTypesOptions struct {
	requestOptions []pipedrive.RequestOption
}

type createActivityTypeOptions struct {
	payload        activityTypePayload
	requestOptions []pipedrive.RequestOption
}

type updateActivityTypeOptions struct {
	payload        activityTypePayload
	requestOptions []pipedrive.RequestOption
}

type deleteActivityTypeOptions struct {
	requestOptions []pipedrive.RequestOption
}

type activityTypePayload struct {
	name    *string
	iconKey *string
	color   *string
	order   *int
}

type activityTypesRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o activityTypesRequestOptions) applyListActivityTypes(cfg *listActivityTypesOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o activityTypesRequestOptions) applyCreateActivityType(cfg *createActivityTypeOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o activityTypesRequestOptions) applyUpdateActivityType(cfg *updateActivityTypeOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o activityTypesRequestOptions) applyDeleteActivityType(cfg *deleteActivityTypeOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listActivityTypesOptionFunc func(*listActivityTypesOptions)

func (f listActivityTypesOptionFunc) applyListActivityTypes(cfg *listActivityTypesOptions) {
	f(cfg)
}

type activityTypeFieldOption func(*activityTypePayload)

func (f activityTypeFieldOption) applyCreateActivityType(cfg *createActivityTypeOptions) {
	f(&cfg.payload)
}

func (f activityTypeFieldOption) applyUpdateActivityType(cfg *updateActivityTypeOptions) {
	f(&cfg.payload)
}

func WithActivityTypesRequestOptions(opts ...pipedrive.RequestOption) ActivityTypesRequestOption {
	return activityTypesRequestOptions{requestOptions: opts}
}

func WithActivityTypeName(name string) ActivityTypeOption {
	return activityTypeFieldOption(func(cfg *activityTypePayload) {
		cfg.name = &name
	})
}

func WithActivityTypeIconKey(iconKey string) ActivityTypeOption {
	return activityTypeFieldOption(func(cfg *activityTypePayload) {
		cfg.iconKey = &iconKey
	})
}

func WithActivityTypeColor(color string) ActivityTypeOption {
	return activityTypeFieldOption(func(cfg *activityTypePayload) {
		cfg.color = &color
	})
}

func WithActivityTypeOrder(order int) UpdateActivityTypeOption {
	return activityTypeFieldOption(func(cfg *activityTypePayload) {
		cfg.order = &order
	})
}

func newListActivityTypesOptions(opts []ListActivityTypesOption) listActivityTypesOptions {
	var cfg listActivityTypesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListActivityTypes(&cfg)
	}
	return cfg
}

func newCreateActivityTypeOptions(opts []CreateActivityTypeOption) createActivityTypeOptions {
	var cfg createActivityTypeOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateActivityType(&cfg)
	}
	return cfg
}

func newUpdateActivityTypeOptions(opts []UpdateActivityTypeOption) updateActivityTypeOptions {
	var cfg updateActivityTypeOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateActivityType(&cfg)
	}
	return cfg
}

func newDeleteActivityTypeOptions(opts []DeleteActivityTypeOption) deleteActivityTypeOptions {
	var cfg deleteActivityTypeOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteActivityType(&cfg)
	}
	return cfg
}

func (s *ActivityTypesService) List(ctx context.Context, opts ...ListActivityTypesOption) ([]ActivityType, error) {
	cfg := newListActivityTypesOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetActivityTypes(ctx, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errorFromResponse(resp, body)
	}

	var payload struct {
		Data []ActivityType `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *ActivityTypesService) Create(ctx context.Context, opts ...CreateActivityTypeOption) (*ActivityType, error) {
	cfg := newCreateActivityTypeOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.name == nil || cfg.payload.iconKey == nil {
		return nil, fmt.Errorf("name and icon key are required")
	}

	body, err := json.Marshal(cfg.payload.toMap(false))
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddActivityTypeWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *ActivityType `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing activity type data in response")
	}
	return payload.Data, nil
}

func (s *ActivityTypesService) Update(ctx context.Context, id ActivityTypeID, opts ...UpdateActivityTypeOption) (*ActivityType, error) {
	cfg := newUpdateActivityTypeOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	bodyMap := cfg.payload.toMap(true)
	if len(bodyMap) == 0 {
		return nil, fmt.Errorf("at least one field is required to update")
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateActivityTypeWithBody(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *ActivityType `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing activity type data in response")
	}
	return payload.Data, nil
}

func (s *ActivityTypesService) Delete(ctx context.Context, id ActivityTypeID, opts ...DeleteActivityTypeOption) (*ActivityType, error) {
	cfg := newDeleteActivityTypeOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteActivityType(ctx, int(id), toRequestEditors(editors)...)
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
		Data *ActivityType `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing activity type data in response")
	}
	return payload.Data, nil
}

func (p activityTypePayload) toMap(includeOrder bool) map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if p.iconKey != nil {
		body["icon_key"] = *p.iconKey
	}
	if p.color != nil {
		body["color"] = *p.color
	}
	if includeOrder && p.order != nil {
		body["order_nr"] = *p.order
	}
	return body
}
