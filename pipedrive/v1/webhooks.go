package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type WebhookID int64

type WebhookType string

const (
	WebhookTypeGeneral     WebhookType = "general"
	WebhookTypeApplication WebhookType = "application"
	WebhookTypeAutomation  WebhookType = "automation"
)

type WebhookEventAction string

const (
	WebhookEventActionCreate WebhookEventAction = "create"
	WebhookEventActionChange WebhookEventAction = "change"
	WebhookEventActionDelete WebhookEventAction = "delete"
	WebhookEventActionAll    WebhookEventAction = "*"
)

type WebhookEventObject string

const (
	WebhookEventObjectActivity     WebhookEventObject = "activity"
	WebhookEventObjectDeal         WebhookEventObject = "deal"
	WebhookEventObjectLead         WebhookEventObject = "lead"
	WebhookEventObjectNote         WebhookEventObject = "note"
	WebhookEventObjectOrganization WebhookEventObject = "organization"
	WebhookEventObjectPerson       WebhookEventObject = "person"
	WebhookEventObjectPipeline     WebhookEventObject = "pipeline"
	WebhookEventObjectProduct      WebhookEventObject = "product"
	WebhookEventObjectStage        WebhookEventObject = "stage"
	WebhookEventObjectUser         WebhookEventObject = "user"
	WebhookEventObjectAll          WebhookEventObject = "*"
)

type WebhookVersion string

const (
	WebhookVersion1 WebhookVersion = "1.0"
	WebhookVersion2 WebhookVersion = "2.0"
)

type Webhook struct {
	ID               WebhookID   `json:"id,omitempty"`
	CompanyID        int         `json:"company_id,omitempty"`
	OwnerID          *UserID     `json:"owner_id,omitempty"`
	UserID           *UserID     `json:"user_id,omitempty"`
	EventAction      string      `json:"event_action,omitempty"`
	EventObject      string      `json:"event_object,omitempty"`
	SubscriptionURL  string      `json:"subscription_url,omitempty"`
	Version          string      `json:"version,omitempty"`
	IsActive         NumberBool  `json:"is_active,omitempty"`
	AddTime          *DateTime   `json:"add_time,omitempty"`
	RemoveTime       *DateTime   `json:"remove_time,omitempty"`
	Type             WebhookType `json:"type,omitempty"`
	HTTPAuthUser     *string     `json:"http_auth_user,omitempty"`
	HTTPAuthPassword *string     `json:"http_auth_password,omitempty"`
	RemoveReason     *string     `json:"remove_reason,omitempty"`
	LastDeliveryTime *DateTime   `json:"last_delivery_time,omitempty"`
	LastHTTPStatus   *int        `json:"last_http_status,omitempty"`
	AdminID          *UserID     `json:"admin_id,omitempty"`
	Name             string      `json:"name,omitempty"`
}

type WebhooksService struct {
	client *Client
}

type ListWebhooksOption interface {
	applyListWebhooks(*listWebhooksOptions)
}

type CreateWebhookOption interface {
	applyCreateWebhook(*createWebhookOptions)
}

type DeleteWebhookOption interface {
	applyDeleteWebhook(*deleteWebhookOptions)
}

type WebhooksRequestOption interface {
	ListWebhooksOption
	CreateWebhookOption
	DeleteWebhookOption
}

type WebhookOption interface {
	CreateWebhookOption
}

type listWebhooksOptions struct {
	requestOptions []pipedrive.RequestOption
}

type createWebhookOptions struct {
	payload        webhookPayload
	requestOptions []pipedrive.RequestOption
}

type deleteWebhookOptions struct {
	requestOptions []pipedrive.RequestOption
}

type webhookPayload struct {
	subscriptionURL  *string
	eventAction      *WebhookEventAction
	eventObject      *WebhookEventObject
	name             *string
	userID           *UserID
	httpAuthUser     *string
	httpAuthPassword *string
	version          *WebhookVersion
}

type webhooksRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o webhooksRequestOptions) applyListWebhooks(cfg *listWebhooksOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o webhooksRequestOptions) applyCreateWebhook(cfg *createWebhookOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o webhooksRequestOptions) applyDeleteWebhook(cfg *deleteWebhookOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listWebhooksOptionFunc func(*listWebhooksOptions)

func (f listWebhooksOptionFunc) applyListWebhooks(cfg *listWebhooksOptions) {
	f(cfg)
}

type webhookFieldOption func(*webhookPayload)

func (f webhookFieldOption) applyCreateWebhook(cfg *createWebhookOptions) {
	f(&cfg.payload)
}

type deleteWebhookOptionFunc func(*deleteWebhookOptions)

func (f deleteWebhookOptionFunc) applyDeleteWebhook(cfg *deleteWebhookOptions) {
	f(cfg)
}

func WithWebhooksRequestOptions(opts ...pipedrive.RequestOption) WebhooksRequestOption {
	return webhooksRequestOptions{requestOptions: opts}
}

func WithWebhookSubscriptionURL(url string) WebhookOption {
	return webhookFieldOption(func(cfg *webhookPayload) {
		cfg.subscriptionURL = &url
	})
}

func WithWebhookEventAction(action WebhookEventAction) WebhookOption {
	return webhookFieldOption(func(cfg *webhookPayload) {
		cfg.eventAction = &action
	})
}

func WithWebhookEventObject(object WebhookEventObject) WebhookOption {
	return webhookFieldOption(func(cfg *webhookPayload) {
		cfg.eventObject = &object
	})
}

func WithWebhookName(name string) WebhookOption {
	return webhookFieldOption(func(cfg *webhookPayload) {
		cfg.name = &name
	})
}

func WithWebhookUserID(id UserID) CreateWebhookOption {
	return webhookFieldOption(func(cfg *webhookPayload) {
		cfg.userID = &id
	})
}

func WithWebhookHTTPAuthUser(user string) CreateWebhookOption {
	return webhookFieldOption(func(cfg *webhookPayload) {
		cfg.httpAuthUser = &user
	})
}

func WithWebhookHTTPAuthPassword(password string) CreateWebhookOption {
	return webhookFieldOption(func(cfg *webhookPayload) {
		cfg.httpAuthPassword = &password
	})
}

func WithWebhookVersion(version WebhookVersion) CreateWebhookOption {
	return webhookFieldOption(func(cfg *webhookPayload) {
		cfg.version = &version
	})
}

func newListWebhooksOptions(opts []ListWebhooksOption) listWebhooksOptions {
	var cfg listWebhooksOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListWebhooks(&cfg)
	}
	return cfg
}

func newCreateWebhookOptions(opts []CreateWebhookOption) createWebhookOptions {
	var cfg createWebhookOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateWebhook(&cfg)
	}
	return cfg
}

func newDeleteWebhookOptions(opts []DeleteWebhookOption) deleteWebhookOptions {
	var cfg deleteWebhookOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteWebhook(&cfg)
	}
	return cfg
}

func (s *WebhooksService) List(ctx context.Context, opts ...ListWebhooksOption) ([]Webhook, error) {
	cfg := newListWebhooksOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetWebhooks(ctx, toRequestEditors(editors)...)
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
		Data []Webhook `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func (s *WebhooksService) Create(ctx context.Context, opts ...CreateWebhookOption) (*Webhook, error) {
	cfg := newCreateWebhookOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.subscriptionURL == nil {
		return nil, fmt.Errorf("subscription URL is required")
	}
	if cfg.payload.eventAction == nil {
		return nil, fmt.Errorf("event action is required")
	}
	if cfg.payload.eventObject == nil {
		return nil, fmt.Errorf("event object is required")
	}
	if cfg.payload.name == nil {
		return nil, fmt.Errorf("name is required")
	}

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddWebhookWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *Webhook `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing webhook data in response")
	}
	return payload.Data, nil
}

func (s *WebhooksService) Delete(ctx context.Context, id WebhookID, opts ...DeleteWebhookOption) (bool, error) {
	cfg := newDeleteWebhookOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteWebhook(ctx, int(id), toRequestEditors(editors)...)
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
		Success *bool `json:"success"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}
	if payload.Success == nil {
		return false, fmt.Errorf("missing webhook delete success in response")
	}
	return *payload.Success, nil
}

func (p webhookPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.subscriptionURL != nil {
		body["subscription_url"] = *p.subscriptionURL
	}
	if p.eventAction != nil {
		body["event_action"] = string(*p.eventAction)
	}
	if p.eventObject != nil {
		body["event_object"] = string(*p.eventObject)
	}
	if p.name != nil {
		body["name"] = *p.name
	}
	if p.userID != nil {
		body["user_id"] = int(*p.userID)
	}
	if p.httpAuthUser != nil {
		body["http_auth_user"] = *p.httpAuthUser
	}
	if p.httpAuthPassword != nil {
		body["http_auth_password"] = *p.httpAuthPassword
	}
	if p.version != nil {
		body["version"] = string(*p.version)
	}
	return body
}
