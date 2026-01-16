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

type FilterType string

const (
	FilterTypeActivity FilterType = "activity"
	FilterTypeDeals    FilterType = "deals"
	FilterTypeLeads    FilterType = "leads"
	FilterTypeOrg      FilterType = "org"
	FilterTypePeople   FilterType = "people"
	FilterTypeProducts FilterType = "products"
	FilterTypeProjects FilterType = "projects"
)

type FilterConditions map[string]interface{}

type Filter struct {
	ID           FilterID         `json:"id,omitempty"`
	Name         string           `json:"name,omitempty"`
	Type         FilterType       `json:"type,omitempty"`
	UserID       *UserID          `json:"user_id,omitempty"`
	VisibleTo    int              `json:"visible_to,omitempty"`
	Active       bool             `json:"active_flag,omitempty"`
	AddTime      *DateTime        `json:"add_time,omitempty"`
	UpdateTime   *DateTime        `json:"update_time,omitempty"`
	CustomViewID int              `json:"custom_view_id,omitempty"`
	Temporary    bool             `json:"temporary_flag,omitempty"`
	Conditions   FilterConditions `json:"conditions,omitempty"`
}

type FilterDeleteResult struct {
	ID FilterID `json:"id"`
}

type FiltersDeleteResult struct {
	IDs []FilterID `json:"id"`
}

type FiltersService struct {
	client *Client
}

type ListFiltersOption interface {
	applyListFilters(*listFiltersOptions)
}

type GetFilterOption interface {
	applyGetFilter(*getFilterOptions)
}

type CreateFilterOption interface {
	applyCreateFilter(*createFilterOptions)
}

type UpdateFilterOption interface {
	applyUpdateFilter(*updateFilterOptions)
}

type DeleteFilterOption interface {
	applyDeleteFilter(*deleteFilterOptions)
}

type DeleteFiltersOption interface {
	applyDeleteFilters(*deleteFiltersOptions)
}

type ListFilterHelpersOption interface {
	applyListFilterHelpers(*listFilterHelpersOptions)
}

type FiltersRequestOption interface {
	ListFiltersOption
	GetFilterOption
	CreateFilterOption
	UpdateFilterOption
	DeleteFilterOption
	DeleteFiltersOption
	ListFilterHelpersOption
}

type FilterOption interface {
	CreateFilterOption
	UpdateFilterOption
}

type listFiltersOptions struct {
	params         genv1.GetFiltersParams
	requestOptions []pipedrive.RequestOption
}

type getFilterOptions struct {
	requestOptions []pipedrive.RequestOption
}

type createFilterOptions struct {
	payload        filterPayload
	requestOptions []pipedrive.RequestOption
}

type updateFilterOptions struct {
	payload        filterPayload
	requestOptions []pipedrive.RequestOption
}

type deleteFilterOptions struct {
	requestOptions []pipedrive.RequestOption
}

type deleteFiltersOptions struct {
	requestOptions []pipedrive.RequestOption
}

type listFilterHelpersOptions struct {
	requestOptions []pipedrive.RequestOption
}

type filterPayload struct {
	name       *string
	filterType *FilterType
	conditions *FilterConditions
}

type filtersRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o filtersRequestOptions) applyListFilters(cfg *listFiltersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o filtersRequestOptions) applyGetFilter(cfg *getFilterOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o filtersRequestOptions) applyCreateFilter(cfg *createFilterOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o filtersRequestOptions) applyUpdateFilter(cfg *updateFilterOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o filtersRequestOptions) applyDeleteFilter(cfg *deleteFilterOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o filtersRequestOptions) applyDeleteFilters(cfg *deleteFiltersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o filtersRequestOptions) applyListFilterHelpers(cfg *listFilterHelpersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listFiltersOptionFunc func(*listFiltersOptions)

func (f listFiltersOptionFunc) applyListFilters(cfg *listFiltersOptions) {
	f(cfg)
}

type filterFieldOption func(*filterPayload)

func (f filterFieldOption) applyCreateFilter(cfg *createFilterOptions) {
	f(&cfg.payload)
}

func (f filterFieldOption) applyUpdateFilter(cfg *updateFilterOptions) {
	f(&cfg.payload)
}

type getFilterOptionFunc func(*getFilterOptions)

func (f getFilterOptionFunc) applyGetFilter(cfg *getFilterOptions) {
	f(cfg)
}

type deleteFilterOptionFunc func(*deleteFilterOptions)

func (f deleteFilterOptionFunc) applyDeleteFilter(cfg *deleteFilterOptions) {
	f(cfg)
}

type deleteFiltersOptionFunc func(*deleteFiltersOptions)

func (f deleteFiltersOptionFunc) applyDeleteFilters(cfg *deleteFiltersOptions) {
	f(cfg)
}

type listFilterHelpersOptionFunc func(*listFilterHelpersOptions)

func (f listFilterHelpersOptionFunc) applyListFilterHelpers(cfg *listFilterHelpersOptions) {
	f(cfg)
}

func WithFiltersRequestOptions(opts ...pipedrive.RequestOption) FiltersRequestOption {
	return filtersRequestOptions{requestOptions: opts}
}

func WithFiltersType(filterType FilterType) ListFiltersOption {
	return listFiltersOptionFunc(func(cfg *listFiltersOptions) {
		value := genv1.GetFiltersParamsType(filterType)
		cfg.params.Type = &value
	})
}

func WithFilterName(name string) FilterOption {
	return filterFieldOption(func(cfg *filterPayload) {
		cfg.name = &name
	})
}

func WithFilterType(filterType FilterType) CreateFilterOption {
	return filterFieldOption(func(cfg *filterPayload) {
		cfg.filterType = &filterType
	})
}

func WithFilterConditions(conditions FilterConditions) FilterOption {
	return filterFieldOption(func(cfg *filterPayload) {
		if conditions == nil {
			cfg.conditions = nil
			return
		}
		clone := make(FilterConditions, len(conditions))
		for key, value := range conditions {
			clone[key] = value
		}
		cfg.conditions = &clone
	})
}

func newListFiltersOptions(opts []ListFiltersOption) listFiltersOptions {
	var cfg listFiltersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListFilters(&cfg)
	}
	return cfg
}

func newGetFilterOptions(opts []GetFilterOption) getFilterOptions {
	var cfg getFilterOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetFilter(&cfg)
	}
	return cfg
}

func newCreateFilterOptions(opts []CreateFilterOption) createFilterOptions {
	var cfg createFilterOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateFilter(&cfg)
	}
	return cfg
}

func newUpdateFilterOptions(opts []UpdateFilterOption) updateFilterOptions {
	var cfg updateFilterOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateFilter(&cfg)
	}
	return cfg
}

func newDeleteFilterOptions(opts []DeleteFilterOption) deleteFilterOptions {
	var cfg deleteFilterOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteFilter(&cfg)
	}
	return cfg
}

func newDeleteFiltersOptions(opts []DeleteFiltersOption) deleteFiltersOptions {
	var cfg deleteFiltersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteFilters(&cfg)
	}
	return cfg
}

func newListFilterHelpersOptions(opts []ListFilterHelpersOption) listFilterHelpersOptions {
	var cfg listFilterHelpersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListFilterHelpers(&cfg)
	}
	return cfg
}

func (s *FiltersService) List(ctx context.Context, opts ...ListFiltersOption) ([]Filter, error) {
	cfg := newListFiltersOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetFilters(ctx, &cfg.params, toRequestEditors(editors)...)
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
		Data []Filter `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *FiltersService) Get(ctx context.Context, id FilterID, opts ...GetFilterOption) (*Filter, error) {
	cfg := newGetFilterOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetFilter(ctx, int(id), toRequestEditors(editors)...)
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
		Data *Filter `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing filter data in response")
	}
	return payload.Data, nil
}

func (s *FiltersService) Create(ctx context.Context, opts ...CreateFilterOption) (*Filter, error) {
	cfg := newCreateFilterOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.name == nil {
		return nil, fmt.Errorf("name is required")
	}
	if cfg.payload.filterType == nil {
		return nil, fmt.Errorf("filter type is required")
	}
	if cfg.payload.conditions == nil {
		return nil, fmt.Errorf("conditions are required")
	}

	body, err := json.Marshal(cfg.payload.toMap(true))
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddFilterWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *Filter `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing filter data in response")
	}
	return payload.Data, nil
}

func (s *FiltersService) Update(ctx context.Context, id FilterID, opts ...UpdateFilterOption) (*Filter, error) {
	cfg := newUpdateFilterOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.name == nil && cfg.payload.conditions == nil {
		return nil, fmt.Errorf("name or conditions are required")
	}

	body, err := json.Marshal(cfg.payload.toMap(false))
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateFilterWithBody(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *Filter `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing filter data in response")
	}
	return payload.Data, nil
}

func (s *FiltersService) Delete(ctx context.Context, id FilterID, opts ...DeleteFilterOption) (*FilterDeleteResult, error) {
	cfg := newDeleteFilterOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteFilter(ctx, int(id), toRequestEditors(editors)...)
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
		Data *struct {
			ID FilterID `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing delete filter data in response")
	}
	return &FilterDeleteResult{ID: payload.Data.ID}, nil
}

func (s *FiltersService) DeleteBulk(ctx context.Context, ids []FilterID, opts ...DeleteFiltersOption) (*FiltersDeleteResult, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("filter IDs are required")
	}
	cfg := newDeleteFiltersOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	params := genv1.DeleteFiltersParams{Ids: joinIDs(ids)}
	resp, err := s.client.gen.DeleteFilters(ctx, &params, toRequestEditors(editors)...)
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
		Data *struct {
			IDs []FilterID `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing delete filters data in response")
	}
	return &FiltersDeleteResult{IDs: payload.Data.IDs}, nil
}

func (s *FiltersService) ListHelpers(ctx context.Context, opts ...ListFilterHelpersOption) (map[string]interface{}, error) {
	cfg := newListFilterHelpersOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetFilterHelpers(ctx, toRequestEditors(editors)...)
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

	var payload map[string]interface{}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload, nil
}

func (p filterPayload) toMap(includeType bool) map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if includeType && p.filterType != nil {
		body["type"] = string(*p.filterType)
	}
	if p.conditions != nil {
		body["conditions"] = map[string]interface{}(*p.conditions)
	}
	return body
}
