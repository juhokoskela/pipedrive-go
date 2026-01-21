package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type PermissionSetApp string

const (
	PermissionSetAppSales           PermissionSetApp = "sales"
	PermissionSetAppProjects        PermissionSetApp = "projects"
	PermissionSetAppCampaigns       PermissionSetApp = "campaigns"
	PermissionSetAppGlobal          PermissionSetApp = "global"
	PermissionSetAppAccountSettings PermissionSetApp = "account_settings"
)

type PermissionSetType string

const (
	PermissionSetTypeAdmin   PermissionSetType = "admin"
	PermissionSetTypeManager PermissionSetType = "manager"
	PermissionSetTypeRegular PermissionSetType = "regular"
	PermissionSetTypeCustom  PermissionSetType = "custom"
)

type PermissionSet struct {
	ID              PermissionSetID   `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	Description     string            `json:"description,omitempty"`
	App             PermissionSetApp  `json:"app,omitempty"`
	Type            PermissionSetType `json:"type,omitempty"`
	AssignmentCount int               `json:"assignment_count,omitempty"`
	Contents        []string          `json:"contents,omitempty"`
}

type PermissionSetAssignment struct {
	UserID          UserID          `json:"user_id,omitempty"`
	PermissionSetID PermissionSetID `json:"permission_set_id,omitempty"`
	Name            string          `json:"name,omitempty"`
}

type PermissionSetsService struct {
	client *Client
}

type ListPermissionSetsOption interface {
	applyListPermissionSets(*listPermissionSetsOptions)
}

type GetPermissionSetOption interface {
	applyGetPermissionSet(*getPermissionSetOptions)
}

type ListPermissionSetAssignmentsOption interface {
	applyListPermissionSetAssignments(*listPermissionSetAssignmentsOptions)
}

type PermissionSetsRequestOption interface {
	ListPermissionSetsOption
	GetPermissionSetOption
	ListPermissionSetAssignmentsOption
}

type listPermissionSetsOptions struct {
	params         genv1.GetPermissionSetsParams
	requestOptions []pipedrive.RequestOption
}

type getPermissionSetOptions struct {
	requestOptions []pipedrive.RequestOption
}

type listPermissionSetAssignmentsOptions struct {
	params         genv1.GetPermissionSetAssignmentsParams
	requestOptions []pipedrive.RequestOption
}

type permissionSetsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o permissionSetsRequestOptions) applyListPermissionSets(cfg *listPermissionSetsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o permissionSetsRequestOptions) applyGetPermissionSet(cfg *getPermissionSetOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o permissionSetsRequestOptions) applyListPermissionSetAssignments(cfg *listPermissionSetAssignmentsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listPermissionSetsOptionFunc func(*listPermissionSetsOptions)

func (f listPermissionSetsOptionFunc) applyListPermissionSets(cfg *listPermissionSetsOptions) {
	f(cfg)
}

type getPermissionSetOptionFunc func(*getPermissionSetOptions)

func (f getPermissionSetOptionFunc) applyGetPermissionSet(cfg *getPermissionSetOptions) {
	f(cfg)
}

type listPermissionSetAssignmentsOptionFunc func(*listPermissionSetAssignmentsOptions)

func (f listPermissionSetAssignmentsOptionFunc) applyListPermissionSetAssignments(cfg *listPermissionSetAssignmentsOptions) {
	f(cfg)
}

func WithPermissionSetsRequestOptions(opts ...pipedrive.RequestOption) PermissionSetsRequestOption {
	return permissionSetsRequestOptions{requestOptions: opts}
}

func WithPermissionSetsApp(app PermissionSetApp) ListPermissionSetsOption {
	return listPermissionSetsOptionFunc(func(cfg *listPermissionSetsOptions) {
		value := genv1.GetPermissionSetsParamsApp(app)
		cfg.params.App = &value
	})
}

func WithPermissionSetAssignmentsStart(start int) ListPermissionSetAssignmentsOption {
	return listPermissionSetAssignmentsOptionFunc(func(cfg *listPermissionSetAssignmentsOptions) {
		cfg.params.Start = &start
	})
}

func WithPermissionSetAssignmentsLimit(limit int) ListPermissionSetAssignmentsOption {
	return listPermissionSetAssignmentsOptionFunc(func(cfg *listPermissionSetAssignmentsOptions) {
		cfg.params.Limit = &limit
	})
}

func newListPermissionSetsOptions(opts []ListPermissionSetsOption) listPermissionSetsOptions {
	var cfg listPermissionSetsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListPermissionSets(&cfg)
	}
	return cfg
}

func newGetPermissionSetOptions(opts []GetPermissionSetOption) getPermissionSetOptions {
	var cfg getPermissionSetOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetPermissionSet(&cfg)
	}
	return cfg
}

func newListPermissionSetAssignmentsOptions(opts []ListPermissionSetAssignmentsOption) listPermissionSetAssignmentsOptions {
	var cfg listPermissionSetAssignmentsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListPermissionSetAssignments(&cfg)
	}
	return cfg
}

func (s *PermissionSetsService) List(ctx context.Context, opts ...ListPermissionSetsOption) ([]PermissionSet, error) {
	cfg := newListPermissionSetsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetPermissionSets(ctx, &cfg.params, toRequestEditors(editors)...)
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
		Data []PermissionSet `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *PermissionSetsService) Get(ctx context.Context, id PermissionSetID, opts ...GetPermissionSetOption) (*PermissionSet, error) {
	cfg := newGetPermissionSetOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetPermissionSet(ctx, string(id), toRequestEditors(editors)...)
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

	set, err := decodePermissionSet(respBody)
	if err != nil {
		return nil, err
	}
	if set == nil {
		return nil, fmt.Errorf("missing permission set data in response")
	}
	return set, nil
}

func (s *PermissionSetsService) ListAssignments(ctx context.Context, id PermissionSetID, opts ...ListPermissionSetAssignmentsOption) ([]PermissionSetAssignment, error) {
	cfg := newListPermissionSetAssignmentsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetPermissionSetAssignments(ctx, string(id), &cfg.params, toRequestEditors(editors)...)
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
		Data []PermissionSetAssignment `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func decodePermissionSet(respBody []byte) (*PermissionSet, error) {
	var wrapped struct {
		Data *PermissionSet `json:"data"`
	}
	if err := json.Unmarshal(respBody, &wrapped); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if wrapped.Data != nil {
		return wrapped.Data, nil
	}

	var set PermissionSet
	if err := json.Unmarshal(respBody, &set); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if set.ID == "" {
		return nil, fmt.Errorf("missing permission set data in response")
	}
	return &set, nil
}
