package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type UserSettings struct {
	AutofillDealExpectedCloseDate   bool   `json:"autofill_deal_expected_close_date,omitempty"`
	BetaApp                         bool   `json:"beta_app,omitempty"`
	CalltoLinkSyntax                string `json:"callto_link_syntax,omitempty"`
	FileUploadDestination           string `json:"file_upload_destination,omitempty"`
	ListLimit                       int    `json:"list_limit,omitempty"`
	MarketplaceTeam                 bool   `json:"marketplace_team,omitempty"`
	PersonDuplicateCondition        string `json:"person_duplicate_condition,omitempty"`
	PreventSalesphoneCalltoOverride bool   `json:"prevent_salesphone_callto_override,omitempty"`
}

type UserSettingsService struct {
	client *Client
}

type GetUserSettingsOption interface {
	applyGetUserSettings(*getUserSettingsOptions)
}

type UserSettingsRequestOption interface {
	GetUserSettingsOption
}

type getUserSettingsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type userSettingsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o userSettingsRequestOptions) applyGetUserSettings(cfg *getUserSettingsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type getUserSettingsOptionFunc func(*getUserSettingsOptions)

func (f getUserSettingsOptionFunc) applyGetUserSettings(cfg *getUserSettingsOptions) {
	f(cfg)
}

func WithUserSettingsRequestOptions(opts ...pipedrive.RequestOption) UserSettingsRequestOption {
	return userSettingsRequestOptions{requestOptions: opts}
}

func newGetUserSettingsOptions(opts []GetUserSettingsOption) getUserSettingsOptions {
	var cfg getUserSettingsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetUserSettings(&cfg)
	}
	return cfg
}

func (s *UserSettingsService) Get(ctx context.Context, opts ...GetUserSettingsOption) (*UserSettings, error) {
	cfg := newGetUserSettingsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetUserSettings(ctx, toRequestEditors(editors)...)
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
		Data UserSettings `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	settings := payload.Data
	return &settings, nil
}
