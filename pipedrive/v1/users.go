package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type UserAccessApp string

const (
	UserAccessAppGlobal          UserAccessApp = "global"
	UserAccessAppSales           UserAccessApp = "sales"
	UserAccessAppCampaigns       UserAccessApp = "campaigns"
	UserAccessAppProjects        UserAccessApp = "projects"
	UserAccessAppAccountSettings UserAccessApp = "account_settings"
	UserAccessAppPartnership     UserAccessApp = "partnership"
)

type UserAccess struct {
	App             UserAccessApp   `json:"app,omitempty"`
	Admin           bool            `json:"admin,omitempty"`
	PermissionSetID PermissionSetID `json:"permission_set_id,omitempty"`
}

type User struct {
	ID                UserID       `json:"id,omitempty"`
	Name              string       `json:"name,omitempty"`
	DefaultCurrency   string       `json:"default_currency,omitempty"`
	Locale            string       `json:"locale,omitempty"`
	Lang              int          `json:"lang,omitempty"`
	Email             string       `json:"email,omitempty"`
	Phone             *string      `json:"phone,omitempty"`
	Activated         bool         `json:"activated,omitempty"`
	LastLogin         *DateTime    `json:"last_login,omitempty"`
	Created           *DateTime    `json:"created,omitempty"`
	Modified          *DateTime    `json:"modified,omitempty"`
	HasCreatedCompany bool         `json:"has_created_company,omitempty"`
	Access            []UserAccess `json:"access,omitempty"`
	Active            bool         `json:"active_flag,omitempty"`
	TimezoneName      string       `json:"timezone_name,omitempty"`
	TimezoneOffset    string       `json:"timezone_offset,omitempty"`
	RoleID            int          `json:"role_id,omitempty"`
	IconURL           *string      `json:"icon_url,omitempty"`
	IsYou             bool         `json:"is_you,omitempty"`
	IsDeleted         bool         `json:"is_deleted,omitempty"`
}

type UserLanguage struct {
	LanguageCode string `json:"language_code,omitempty"`
	CountryCode  string `json:"country_code,omitempty"`
}

type CurrentUser struct {
	User
	CompanyID       int           `json:"company_id,omitempty"`
	CompanyName     string        `json:"company_name,omitempty"`
	CompanyDomain   string        `json:"company_domain,omitempty"`
	CompanyCountry  string        `json:"company_country,omitempty"`
	CompanyIndustry string        `json:"company_industry,omitempty"`
	Language        *UserLanguage `json:"language,omitempty"`
}

type UserPermissions struct {
	CanAddCustomFields          bool `json:"can_add_custom_fields,omitempty"`
	CanAddProducts              bool `json:"can_add_products,omitempty"`
	CanAddProspectsAsLeads      bool `json:"can_add_prospects_as_leads,omitempty"`
	CanBulkEditItems            bool `json:"can_bulk_edit_items,omitempty"`
	CanChangeVisibilityOfItems  bool `json:"can_change_visibility_of_items,omitempty"`
	CanConvertDealsToLeads      bool `json:"can_convert_deals_to_leads,omitempty"`
	CanCreateOwnWorkflow        bool `json:"can_create_own_workflow,omitempty"`
	CanDeleteActivities         bool `json:"can_delete_activities,omitempty"`
	CanDeleteCustomFields       bool `json:"can_delete_custom_fields,omitempty"`
	CanDeleteDeals              bool `json:"can_delete_deals,omitempty"`
	CanEditCustomFields         bool `json:"can_edit_custom_fields,omitempty"`
	CanEditDealsClosedDate      bool `json:"can_edit_deals_closed_date,omitempty"`
	CanEditProducts             bool `json:"can_edit_products,omitempty"`
	CanEditSharedFilters        bool `json:"can_edit_shared_filters,omitempty"`
	CanExportDataFromLists      bool `json:"can_export_data_from_lists,omitempty"`
	CanFollowOtherUsers         bool `json:"can_follow_other_users,omitempty"`
	CanMergeDeals               bool `json:"can_merge_deals,omitempty"`
	CanMergeOrganizations       bool `json:"can_merge_organizations,omitempty"`
	CanMergePeople              bool `json:"can_merge_people,omitempty"`
	CanModifyLabels             bool `json:"can_modify_labels,omitempty"`
	CanSeeCompanyWideStatistics bool `json:"can_see_company_wide_statistics,omitempty"`
	CanSeeDealsListSummary      bool `json:"can_see_deals_list_summary,omitempty"`
	CanSeeHiddenItemsNames      bool `json:"can_see_hidden_items_names,omitempty"`
	CanSeeOtherUsers            bool `json:"can_see_other_users,omitempty"`
	CanSeeOtherUsersStatistics  bool `json:"can_see_other_users_statistics,omitempty"`
	CanSeeSecurityDashboard     bool `json:"can_see_security_dashboard,omitempty"`
	CanShareFilters             bool `json:"can_share_filters,omitempty"`
	CanShareInsights            bool `json:"can_share_insights,omitempty"`
	CanUseAPI                   bool `json:"can_use_api,omitempty"`
	CanUseEmailTracking         bool `json:"can_use_email_tracking,omitempty"`
	CanUseImport                bool `json:"can_use_import,omitempty"`
}

type UsersService struct {
	client *Client
}

type ListUsersOption interface {
	applyListUsers(*listUsersOptions)
}

type GetUserOption interface {
	applyGetUser(*getUserOptions)
}

type GetCurrentUserOption interface {
	applyGetCurrentUser(*getCurrentUserOptions)
}

type GetUserPermissionsOption interface {
	applyGetUserPermissions(*getUserPermissionsOptions)
}

type UsersRequestOption interface {
	ListUsersOption
	GetUserOption
	GetCurrentUserOption
	GetUserPermissionsOption
}

type listUsersOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getUserOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getCurrentUserOptions struct {
	requestOptions []pipedrive.RequestOption
}

type getUserPermissionsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type usersRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o usersRequestOptions) applyListUsers(cfg *listUsersOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o usersRequestOptions) applyGetUser(cfg *getUserOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o usersRequestOptions) applyGetCurrentUser(cfg *getCurrentUserOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o usersRequestOptions) applyGetUserPermissions(cfg *getUserPermissionsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listUsersOptionFunc func(*listUsersOptions)

func (f listUsersOptionFunc) applyListUsers(cfg *listUsersOptions) {
	f(cfg)
}

type getUserOptionFunc func(*getUserOptions)

func (f getUserOptionFunc) applyGetUser(cfg *getUserOptions) {
	f(cfg)
}

type getCurrentUserOptionFunc func(*getCurrentUserOptions)

func (f getCurrentUserOptionFunc) applyGetCurrentUser(cfg *getCurrentUserOptions) {
	f(cfg)
}

type getUserPermissionsOptionFunc func(*getUserPermissionsOptions)

func (f getUserPermissionsOptionFunc) applyGetUserPermissions(cfg *getUserPermissionsOptions) {
	f(cfg)
}

func WithUsersRequestOptions(opts ...pipedrive.RequestOption) UsersRequestOption {
	return usersRequestOptions{requestOptions: opts}
}

func newListUsersOptions(opts []ListUsersOption) listUsersOptions {
	var cfg listUsersOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListUsers(&cfg)
	}
	return cfg
}

func newGetUserOptions(opts []GetUserOption) getUserOptions {
	var cfg getUserOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetUser(&cfg)
	}
	return cfg
}

func newGetCurrentUserOptions(opts []GetCurrentUserOption) getCurrentUserOptions {
	var cfg getCurrentUserOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetCurrentUser(&cfg)
	}
	return cfg
}

func newGetUserPermissionsOptions(opts []GetUserPermissionsOption) getUserPermissionsOptions {
	var cfg getUserPermissionsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetUserPermissions(&cfg)
	}
	return cfg
}

func (s *UsersService) List(ctx context.Context, opts ...ListUsersOption) ([]User, error) {
	cfg := newListUsersOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetUsers(ctx, toRequestEditors(editors)...)
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
		Data []User `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *UsersService) Get(ctx context.Context, id UserID, opts ...GetUserOption) (*User, error) {
	cfg := newGetUserOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetUser(ctx, int(id), toRequestEditors(editors)...)
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
		Data *User `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing user data in response")
	}
	return payload.Data, nil
}

func (s *UsersService) GetCurrent(ctx context.Context, opts ...GetCurrentUserOption) (*CurrentUser, error) {
	cfg := newGetCurrentUserOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetCurrentUser(ctx, toRequestEditors(editors)...)
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
		Data *CurrentUser `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing current user data in response")
	}
	return payload.Data, nil
}

func (s *UsersService) GetPermissions(ctx context.Context, id UserID, opts ...GetUserPermissionsOption) (*UserPermissions, error) {
	cfg := newGetUserPermissionsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetUserPermissions(ctx, int(id), toRequestEditors(editors)...)
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
		Data *UserPermissions `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing user permissions data in response")
	}
	return payload.Data, nil
}
