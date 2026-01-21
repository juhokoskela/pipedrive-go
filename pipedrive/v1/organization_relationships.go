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

type OrganizationRelationshipType string

const (
	OrganizationRelationshipTypeParent  OrganizationRelationshipType = "parent"
	OrganizationRelationshipTypeRelated OrganizationRelationshipType = "related"
)

type OrganizationRelationshipOrg struct {
	ID          OrganizationID `json:"value,omitempty"`
	Name        string         `json:"name,omitempty"`
	PeopleCount int            `json:"people_count,omitempty"`
	OwnerID     *UserID        `json:"owner_id,omitempty"`
	Address     string         `json:"address,omitempty"`
	CcEmail     string         `json:"cc_email,omitempty"`
}

type OrganizationRelationship struct {
	ID                      OrganizationRelationshipID   `json:"id,omitempty"`
	Type                    OrganizationRelationshipType `json:"type,omitempty"`
	RelOwnerOrg             *OrganizationRelationshipOrg `json:"rel_owner_org_id,omitempty"`
	RelLinkedOrg            *OrganizationRelationshipOrg `json:"rel_linked_org_id,omitempty"`
	AddTime                 *DateTime                    `json:"add_time,omitempty"`
	UpdateTime              *DateTime                    `json:"update_time,omitempty"`
	ActiveFlag              string                       `json:"active_flag,omitempty"`
	CalculatedType          string                       `json:"calculated_type,omitempty"`
	CalculatedRelatedOrgID  *OrganizationID              `json:"calculated_related_org_id,omitempty"`
	RelatedOrganizationName string                       `json:"related_organization_name,omitempty"`
}

type OrganizationRelationshipsPagination struct {
	Start                 int  `json:"start,omitempty"`
	Limit                 int  `json:"limit,omitempty"`
	MoreItemsInCollection bool `json:"more_items_in_collection,omitempty"`
}

type OrganizationRelationshipsAdditionalData struct {
	Pagination *OrganizationRelationshipsPagination `json:"pagination,omitempty"`
}

type OrganizationRelationshipDeleteResult struct {
	ID OrganizationRelationshipID `json:"id"`
}

type OrganizationRelationshipsService struct {
	client *Client
}

type ListOrganizationRelationshipsOption interface {
	applyListOrganizationRelationships(*listOrganizationRelationshipsOptions)
}

type GetOrganizationRelationshipOption interface {
	applyGetOrganizationRelationship(*getOrganizationRelationshipOptions)
}

type CreateOrganizationRelationshipOption interface {
	applyCreateOrganizationRelationship(*createOrganizationRelationshipOptions)
}

type UpdateOrganizationRelationshipOption interface {
	applyUpdateOrganizationRelationship(*updateOrganizationRelationshipOptions)
}

type DeleteOrganizationRelationshipOption interface {
	applyDeleteOrganizationRelationship(*deleteOrganizationRelationshipOptions)
}

type OrganizationRelationshipsRequestOption interface {
	ListOrganizationRelationshipsOption
	GetOrganizationRelationshipOption
	CreateOrganizationRelationshipOption
	UpdateOrganizationRelationshipOption
	DeleteOrganizationRelationshipOption
}

type OrganizationRelationshipOption interface {
	CreateOrganizationRelationshipOption
	UpdateOrganizationRelationshipOption
}

type listOrganizationRelationshipsOptions struct {
	orgID          *OrganizationID
	requestOptions []pipedrive.RequestOption
}

type getOrganizationRelationshipOptions struct {
	orgID          *OrganizationID
	requestOptions []pipedrive.RequestOption
}

type createOrganizationRelationshipOptions struct {
	payload        organizationRelationshipPayload
	requestOptions []pipedrive.RequestOption
}

type updateOrganizationRelationshipOptions struct {
	payload        organizationRelationshipPayload
	requestOptions []pipedrive.RequestOption
}

type deleteOrganizationRelationshipOptions struct {
	requestOptions []pipedrive.RequestOption
}

type organizationRelationshipPayload struct {
	orgID            *OrganizationID
	relationshipType *OrganizationRelationshipType
	ownerOrgID       *OrganizationID
	linkedOrgID      *OrganizationID
}

type organizationRelationshipsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o organizationRelationshipsRequestOptions) applyListOrganizationRelationships(cfg *listOrganizationRelationshipsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRelationshipsRequestOptions) applyGetOrganizationRelationship(cfg *getOrganizationRelationshipOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRelationshipsRequestOptions) applyCreateOrganizationRelationship(cfg *createOrganizationRelationshipOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRelationshipsRequestOptions) applyUpdateOrganizationRelationship(cfg *updateOrganizationRelationshipOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o organizationRelationshipsRequestOptions) applyDeleteOrganizationRelationship(cfg *deleteOrganizationRelationshipOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listOrganizationRelationshipsOptionFunc func(*listOrganizationRelationshipsOptions)

func (f listOrganizationRelationshipsOptionFunc) applyListOrganizationRelationships(cfg *listOrganizationRelationshipsOptions) {
	f(cfg)
}

type getOrganizationRelationshipOptionFunc func(*getOrganizationRelationshipOptions)

func (f getOrganizationRelationshipOptionFunc) applyGetOrganizationRelationship(cfg *getOrganizationRelationshipOptions) {
	f(cfg)
}

type organizationRelationshipFieldOption func(*organizationRelationshipPayload)

func (f organizationRelationshipFieldOption) applyCreateOrganizationRelationship(cfg *createOrganizationRelationshipOptions) {
	f(&cfg.payload)
}

func (f organizationRelationshipFieldOption) applyUpdateOrganizationRelationship(cfg *updateOrganizationRelationshipOptions) {
	f(&cfg.payload)
}

type deleteOrganizationRelationshipOptionFunc func(*deleteOrganizationRelationshipOptions)

func (f deleteOrganizationRelationshipOptionFunc) applyDeleteOrganizationRelationship(cfg *deleteOrganizationRelationshipOptions) {
	f(cfg)
}

func WithOrganizationRelationshipsRequestOptions(opts ...pipedrive.RequestOption) OrganizationRelationshipsRequestOption {
	return organizationRelationshipsRequestOptions{requestOptions: opts}
}

func WithOrganizationRelationshipsOrgID(id OrganizationID) ListOrganizationRelationshipsOption {
	return listOrganizationRelationshipsOptionFunc(func(cfg *listOrganizationRelationshipsOptions) {
		cfg.orgID = &id
	})
}

type organizationRelationshipOrgIDOption struct {
	id OrganizationID
}

func (o organizationRelationshipOrgIDOption) applyCreateOrganizationRelationship(cfg *createOrganizationRelationshipOptions) {
	cfg.payload.orgID = &o.id
}

func (o organizationRelationshipOrgIDOption) applyUpdateOrganizationRelationship(cfg *updateOrganizationRelationshipOptions) {
	cfg.payload.orgID = &o.id
}

func (o organizationRelationshipOrgIDOption) applyGetOrganizationRelationship(cfg *getOrganizationRelationshipOptions) {
	cfg.orgID = &o.id
}

func WithOrganizationRelationshipOrgID(id OrganizationID) organizationRelationshipOrgIDOption {
	return organizationRelationshipOrgIDOption{id: id}
}

func WithOrganizationRelationshipType(value OrganizationRelationshipType) OrganizationRelationshipOption {
	return organizationRelationshipFieldOption(func(cfg *organizationRelationshipPayload) {
		cfg.relationshipType = &value
	})
}

func WithOrganizationRelationshipOwnerID(id OrganizationID) OrganizationRelationshipOption {
	return organizationRelationshipFieldOption(func(cfg *organizationRelationshipPayload) {
		cfg.ownerOrgID = &id
	})
}

func WithOrganizationRelationshipLinkedID(id OrganizationID) OrganizationRelationshipOption {
	return organizationRelationshipFieldOption(func(cfg *organizationRelationshipPayload) {
		cfg.linkedOrgID = &id
	})
}

func newListOrganizationRelationshipsOptions(opts []ListOrganizationRelationshipsOption) listOrganizationRelationshipsOptions {
	var cfg listOrganizationRelationshipsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListOrganizationRelationships(&cfg)
	}
	return cfg
}

func newGetOrganizationRelationshipOptions(opts []GetOrganizationRelationshipOption) getOrganizationRelationshipOptions {
	var cfg getOrganizationRelationshipOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetOrganizationRelationship(&cfg)
	}
	return cfg
}

func newCreateOrganizationRelationshipOptions(opts []CreateOrganizationRelationshipOption) createOrganizationRelationshipOptions {
	var cfg createOrganizationRelationshipOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateOrganizationRelationship(&cfg)
	}
	return cfg
}

func newUpdateOrganizationRelationshipOptions(opts []UpdateOrganizationRelationshipOption) updateOrganizationRelationshipOptions {
	var cfg updateOrganizationRelationshipOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUpdateOrganizationRelationship(&cfg)
	}
	return cfg
}

func newDeleteOrganizationRelationshipOptions(opts []DeleteOrganizationRelationshipOption) deleteOrganizationRelationshipOptions {
	var cfg deleteOrganizationRelationshipOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteOrganizationRelationship(&cfg)
	}
	return cfg
}

func (s *OrganizationRelationshipsService) List(ctx context.Context, opts ...ListOrganizationRelationshipsOption) ([]OrganizationRelationship, *OrganizationRelationshipsAdditionalData, error) {
	cfg := newListOrganizationRelationshipsOptions(opts)
	if cfg.orgID == nil {
		return nil, nil, fmt.Errorf("organization ID is required")
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	params := genv1.GetOrganizationRelationshipsParams{OrgId: int(*cfg.orgID)}
	resp, err := s.client.gen.GetOrganizationRelationships(ctx, &params, toRequestEditors(editors)...)
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
		Data           []OrganizationRelationship               `json:"data"`
		AdditionalData *OrganizationRelationshipsAdditionalData `json:"additional_data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *OrganizationRelationshipsService) Get(ctx context.Context, id OrganizationRelationshipID, opts ...GetOrganizationRelationshipOption) (*OrganizationRelationship, error) {
	cfg := newGetOrganizationRelationshipOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	var params *genv1.GetOrganizationRelationshipParams
	if cfg.orgID != nil {
		value := int(*cfg.orgID)
		params = &genv1.GetOrganizationRelationshipParams{OrgId: &value}
	}

	resp, err := s.client.gen.GetOrganizationRelationship(ctx, int(id), params, toRequestEditors(editors)...)
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
		Data *OrganizationRelationship `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization relationship data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationRelationshipsService) Create(ctx context.Context, opts ...CreateOrganizationRelationshipOption) (*OrganizationRelationship, error) {
	cfg := newCreateOrganizationRelationshipOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.relationshipType == nil {
		return nil, fmt.Errorf("relationship type is required")
	}
	if cfg.payload.ownerOrgID == nil {
		return nil, fmt.Errorf("owner organization ID is required")
	}
	if cfg.payload.linkedOrgID == nil {
		return nil, fmt.Errorf("linked organization ID is required")
	}

	bodyMap := map[string]interface{}{
		"type":              string(*cfg.payload.relationshipType),
		"rel_owner_org_id":  int64(*cfg.payload.ownerOrgID),
		"rel_linked_org_id": int64(*cfg.payload.linkedOrgID),
	}
	if cfg.payload.orgID != nil {
		bodyMap["org_id"] = int64(*cfg.payload.orgID)
	}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddOrganizationRelationshipWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *OrganizationRelationship `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization relationship data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationRelationshipsService) Update(ctx context.Context, id OrganizationRelationshipID, opts ...UpdateOrganizationRelationshipOption) (*OrganizationRelationship, error) {
	cfg := newUpdateOrganizationRelationshipOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	bodyMap := map[string]interface{}{}
	if cfg.payload.relationshipType != nil {
		bodyMap["type"] = string(*cfg.payload.relationshipType)
	}
	if cfg.payload.ownerOrgID != nil {
		bodyMap["rel_owner_org_id"] = int64(*cfg.payload.ownerOrgID)
	}
	if cfg.payload.linkedOrgID != nil {
		bodyMap["rel_linked_org_id"] = int64(*cfg.payload.linkedOrgID)
	}
	if cfg.payload.orgID != nil {
		bodyMap["org_id"] = int64(*cfg.payload.orgID)
	}
	if len(bodyMap) == 0 {
		return nil, fmt.Errorf("no organization relationship fields provided")
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.UpdateOrganizationRelationshipWithBody(ctx, int(id), "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *OrganizationRelationship `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization relationship data in response")
	}
	return payload.Data, nil
}

func (s *OrganizationRelationshipsService) Delete(ctx context.Context, id OrganizationRelationshipID, opts ...DeleteOrganizationRelationshipOption) (*OrganizationRelationshipDeleteResult, error) {
	cfg := newDeleteOrganizationRelationshipOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteOrganizationRelationship(ctx, int(id), toRequestEditors(editors)...)
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
		Data *OrganizationRelationshipDeleteResult `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing organization relationship delete data in response")
	}
	return payload.Data, nil
}
