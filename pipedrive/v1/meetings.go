package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type UserProviderLinkResult struct {
	Message string `json:"message,omitempty"`
}

type MeetingsService struct {
	client *Client
}

type CreateUserProviderLinkOption interface {
	applyCreateUserProviderLink(*createUserProviderLinkOptions)
}

type DeleteUserProviderLinkOption interface {
	applyDeleteUserProviderLink(*deleteUserProviderLinkOptions)
}

type MeetingsRequestOption interface {
	CreateUserProviderLinkOption
	DeleteUserProviderLinkOption
}

type UserProviderLinkOption interface {
	CreateUserProviderLinkOption
}

type createUserProviderLinkOptions struct {
	payload        userProviderLinkPayload
	requestOptions []pipedrive.RequestOption
}

type deleteUserProviderLinkOptions struct {
	requestOptions []pipedrive.RequestOption
}

type userProviderLinkPayload struct {
	userProviderID      *UserProviderLinkID
	userID              *UserID
	companyID           *int
	marketplaceClientID *string
}

type meetingsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o meetingsRequestOptions) applyCreateUserProviderLink(cfg *createUserProviderLinkOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o meetingsRequestOptions) applyDeleteUserProviderLink(cfg *deleteUserProviderLinkOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type userProviderLinkFieldOption func(*userProviderLinkPayload)

func (f userProviderLinkFieldOption) applyCreateUserProviderLink(cfg *createUserProviderLinkOptions) {
	f(&cfg.payload)
}

type deleteUserProviderLinkOptionFunc func(*deleteUserProviderLinkOptions)

func (f deleteUserProviderLinkOptionFunc) applyDeleteUserProviderLink(cfg *deleteUserProviderLinkOptions) {
	f(cfg)
}

func WithMeetingsRequestOptions(opts ...pipedrive.RequestOption) MeetingsRequestOption {
	return meetingsRequestOptions{requestOptions: opts}
}

func WithUserProviderLinkID(id UserProviderLinkID) UserProviderLinkOption {
	return userProviderLinkFieldOption(func(cfg *userProviderLinkPayload) {
		cfg.userProviderID = &id
	})
}

func WithUserProviderLinkUserID(id UserID) UserProviderLinkOption {
	return userProviderLinkFieldOption(func(cfg *userProviderLinkPayload) {
		cfg.userID = &id
	})
}

func WithUserProviderLinkCompanyID(id int) UserProviderLinkOption {
	return userProviderLinkFieldOption(func(cfg *userProviderLinkPayload) {
		cfg.companyID = &id
	})
}

func WithUserProviderLinkMarketplaceClientID(id string) UserProviderLinkOption {
	return userProviderLinkFieldOption(func(cfg *userProviderLinkPayload) {
		cfg.marketplaceClientID = &id
	})
}

func newCreateUserProviderLinkOptions(opts []CreateUserProviderLinkOption) createUserProviderLinkOptions {
	var cfg createUserProviderLinkOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateUserProviderLink(&cfg)
	}
	return cfg
}

func newDeleteUserProviderLinkOptions(opts []DeleteUserProviderLinkOption) deleteUserProviderLinkOptions {
	var cfg deleteUserProviderLinkOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteUserProviderLink(&cfg)
	}
	return cfg
}

func (s *MeetingsService) CreateUserProviderLink(ctx context.Context, opts ...CreateUserProviderLinkOption) (*UserProviderLinkResult, error) {
	cfg := newCreateUserProviderLinkOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.userProviderID == nil {
		return nil, fmt.Errorf("user provider ID is required")
	}
	if cfg.payload.userID == nil {
		return nil, fmt.Errorf("user ID is required")
	}
	if cfg.payload.companyID == nil {
		return nil, fmt.Errorf("company ID is required")
	}
	if cfg.payload.marketplaceClientID == nil {
		return nil, fmt.Errorf("marketplace client ID is required")
	}

	userProviderID, err := parseUUID(string(*cfg.payload.userProviderID), "user provider id")
	if err != nil {
		return nil, err
	}

	body := genv1.SaveUserProviderLinkJSONRequestBody{
		CompanyId:           *cfg.payload.companyID,
		MarketplaceClientId: *cfg.payload.marketplaceClientID,
		UserId:              int(*cfg.payload.userID),
		UserProviderId:      userProviderID,
	}

	resp, err := s.client.gen.SaveUserProviderLink(ctx, body, toRequestEditors(editors)...)
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
		Data *UserProviderLinkResult `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing user provider link data in response")
	}
	return payload.Data, nil
}

func (s *MeetingsService) DeleteUserProviderLink(ctx context.Context, id UserProviderLinkID, opts ...DeleteUserProviderLinkOption) (*UserProviderLinkResult, error) {
	cfg := newDeleteUserProviderLinkOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	userProviderID, err := parseUUID(string(id), "user provider id")
	if err != nil {
		return nil, err
	}

	resp, err := s.client.gen.DeleteUserProviderLink(ctx, userProviderID, toRequestEditors(editors)...)
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
		Data *UserProviderLinkResult `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing user provider link data in response")
	}
	return payload.Data, nil
}
