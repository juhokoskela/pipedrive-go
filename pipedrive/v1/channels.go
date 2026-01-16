package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ChannelProviderType string

const (
	ChannelProviderTypeFacebook ChannelProviderType = "facebook"
	ChannelProviderTypeWhatsapp ChannelProviderType = "whatsapp"
	ChannelProviderTypeOther    ChannelProviderType = "other"
)

type ChannelMessageStatus string

const (
	ChannelMessageStatusSent      ChannelMessageStatus = "sent"
	ChannelMessageStatusDelivered ChannelMessageStatus = "delivered"
	ChannelMessageStatusRead      ChannelMessageStatus = "read"
	ChannelMessageStatusFailed    ChannelMessageStatus = "failed"
)

type Channel struct {
	ID                ChannelID           `json:"id"`
	Name              string              `json:"name,omitempty"`
	AvatarURL         string              `json:"avatar_url,omitempty"`
	ProviderChannelID string              `json:"provider_channel_id,omitempty"`
	MarketplaceClient string              `json:"marketplace_client_id,omitempty"`
	CompanyID         *int                `json:"pd_company_id,omitempty"`
	UserID            *UserID             `json:"pd_user_id,omitempty"`
	CreatedAt         *DateTime           `json:"created_at,omitempty"`
	ProviderType      ChannelProviderType `json:"provider_type,omitempty"`
	TemplateSupport   bool                `json:"template_support,omitempty"`
}

type ChannelMessageAttachment struct {
	ID          string   `json:"id"`
	Type        string   `json:"type,omitempty"`
	Name        string   `json:"name,omitempty"`
	Size        *float64 `json:"size,omitempty"`
	URL         string   `json:"url,omitempty"`
	PreviewURL  string   `json:"preview_url,omitempty"`
	LinkExpires bool     `json:"link_expires,omitempty"`
}

type ChannelMessage struct {
	ID               ChannelMessageID           `json:"id"`
	ChannelID        ChannelID                  `json:"channel_id,omitempty"`
	SenderID         string                     `json:"sender_id,omitempty"`
	ConversationID   ConversationID             `json:"conversation_id,omitempty"`
	Message          string                     `json:"message,omitempty"`
	Status           ChannelMessageStatus       `json:"status,omitempty"`
	CreatedAt        *DateTime                  `json:"created_at,omitempty"`
	ReplyBy          *DateTime                  `json:"reply_by,omitempty"`
	ConversationLink string                     `json:"conversation_link,omitempty"`
	Attachments      []ChannelMessageAttachment `json:"attachments,omitempty"`
}

type ChannelsService struct {
	client *Client
}

type CreateChannelOption interface {
	applyCreateChannel(*createChannelOptions)
}

type DeleteChannelOption interface {
	applyDeleteChannel(*deleteChannelOptions)
}

type ReceiveMessageOption interface {
	applyReceiveMessage(*receiveMessageOptions)
}

type DeleteConversationOption interface {
	applyDeleteConversation(*deleteConversationOptions)
}

type ChannelsRequestOption interface {
	CreateChannelOption
	DeleteChannelOption
	ReceiveMessageOption
	DeleteConversationOption
}

type ChannelOption interface {
	CreateChannelOption
}

type createChannelOptions struct {
	payload        channelPayload
	requestOptions []pipedrive.RequestOption
}

type deleteChannelOptions struct {
	requestOptions []pipedrive.RequestOption
}

type receiveMessageOptions struct {
	payload        channelMessagePayload
	requestOptions []pipedrive.RequestOption
}

type deleteConversationOptions struct {
	requestOptions []pipedrive.RequestOption
}

type channelPayload struct {
	name            *string
	providerID      *string
	avatarURL       *string
	templateSupport *bool
	providerType    *ChannelProviderType
}

type channelMessagePayload struct {
	id               *ChannelMessageID
	channelID        *ChannelID
	senderID         *string
	conversationID   *ConversationID
	message          *string
	status           *ChannelMessageStatus
	createdAt        *time.Time
	replyBy          *time.Time
	conversationLink *string
	attachments      []ChannelMessageAttachment
}

type channelsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o channelsRequestOptions) applyCreateChannel(cfg *createChannelOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o channelsRequestOptions) applyDeleteChannel(cfg *deleteChannelOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o channelsRequestOptions) applyReceiveMessage(cfg *receiveMessageOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o channelsRequestOptions) applyDeleteConversation(cfg *deleteConversationOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type channelFieldOption func(*channelPayload)

func (f channelFieldOption) applyCreateChannel(cfg *createChannelOptions) {
	f(&cfg.payload)
}

type receiveMessageFieldOption func(*channelMessagePayload)

func (f receiveMessageFieldOption) applyReceiveMessage(cfg *receiveMessageOptions) {
	f(&cfg.payload)
}

func WithChannelsRequestOptions(opts ...pipedrive.RequestOption) ChannelsRequestOption {
	return channelsRequestOptions{requestOptions: opts}
}

func WithChannelName(name string) ChannelOption {
	return channelFieldOption(func(cfg *channelPayload) {
		cfg.name = &name
	})
}

func WithChannelProviderID(id string) ChannelOption {
	return channelFieldOption(func(cfg *channelPayload) {
		cfg.providerID = &id
	})
}

func WithChannelAvatarURL(url string) ChannelOption {
	return channelFieldOption(func(cfg *channelPayload) {
		cfg.avatarURL = &url
	})
}

func WithChannelTemplateSupport(enabled bool) ChannelOption {
	return channelFieldOption(func(cfg *channelPayload) {
		cfg.templateSupport = &enabled
	})
}

func WithChannelProviderType(providerType ChannelProviderType) ChannelOption {
	return channelFieldOption(func(cfg *channelPayload) {
		cfg.providerType = &providerType
	})
}

func WithChannelMessageID(id string) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		value := ChannelMessageID(id)
		cfg.id = &value
	})
}

func WithChannelMessageChannelID(id string) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		value := ChannelID(id)
		cfg.channelID = &value
	})
}

func WithChannelMessageSenderID(id string) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		cfg.senderID = &id
	})
}

func WithChannelMessageConversationID(id string) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		value := ConversationID(id)
		cfg.conversationID = &value
	})
}

func WithChannelMessageBody(body string) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		cfg.message = &body
	})
}

func WithChannelMessageStatus(status ChannelMessageStatus) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		cfg.status = &status
	})
}

func WithChannelMessageCreatedAt(t time.Time) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		cfg.createdAt = &t
	})
}

func WithChannelMessageReplyBy(t time.Time) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		cfg.replyBy = &t
	})
}

func WithChannelMessageConversationLink(link string) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		cfg.conversationLink = &link
	})
}

func WithChannelMessageAttachments(attachments ...ChannelMessageAttachment) ReceiveMessageOption {
	return receiveMessageFieldOption(func(cfg *channelMessagePayload) {
		cfg.attachments = append([]ChannelMessageAttachment{}, attachments...)
	})
}

func newCreateChannelOptions(opts []CreateChannelOption) createChannelOptions {
	var cfg createChannelOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyCreateChannel(&cfg)
	}
	return cfg
}

func newDeleteChannelOptions(opts []DeleteChannelOption) deleteChannelOptions {
	var cfg deleteChannelOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteChannel(&cfg)
	}
	return cfg
}

func newReceiveMessageOptions(opts []ReceiveMessageOption) receiveMessageOptions {
	var cfg receiveMessageOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyReceiveMessage(&cfg)
	}
	return cfg
}

func newDeleteConversationOptions(opts []DeleteConversationOption) deleteConversationOptions {
	var cfg deleteConversationOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteConversation(&cfg)
	}
	return cfg
}

func (s *ChannelsService) Create(ctx context.Context, opts ...CreateChannelOption) (*Channel, error) {
	cfg := newCreateChannelOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.name == nil || cfg.payload.providerID == nil {
		return nil, fmt.Errorf("name and provider channel ID are required")
	}

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.AddChannelWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *Channel `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing channel data in response")
	}
	return payload.Data, nil
}

func (s *ChannelsService) Delete(ctx context.Context, id ChannelID, opts ...DeleteChannelOption) (bool, error) {
	cfg := newDeleteChannelOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteChannel(ctx, string(id), toRequestEditors(editors)...)
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
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}
	return payload.Success, nil
}

func (s *ChannelsService) ReceiveMessage(ctx context.Context, opts ...ReceiveMessageOption) (*ChannelMessage, error) {
	cfg := newReceiveMessageOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	if cfg.payload.id == nil || cfg.payload.channelID == nil || cfg.payload.senderID == nil || cfg.payload.conversationID == nil || cfg.payload.message == nil || cfg.payload.status == nil || cfg.payload.createdAt == nil {
		return nil, fmt.Errorf("id, channel ID, sender ID, conversation ID, message, status, and created_at are required")
	}

	body, err := json.Marshal(cfg.payload.toMap())
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	resp, err := s.client.gen.ReceiveMessageWithBody(ctx, "application/json", bytes.NewReader(body), toRequestEditors(editors)...)
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
		Data *ChannelMessage `json:"data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing message data in response")
	}
	return payload.Data, nil
}

func (s *ChannelsService) DeleteConversation(ctx context.Context, channelID ChannelID, conversationID ConversationID, opts ...DeleteConversationOption) (bool, error) {
	cfg := newDeleteConversationOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DeleteConversation(ctx, string(channelID), string(conversationID), toRequestEditors(editors)...)
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
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}
	return payload.Success, nil
}

func (p channelPayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.name != nil {
		body["name"] = *p.name
	}
	if p.providerID != nil {
		body["provider_channel_id"] = *p.providerID
	}
	if p.avatarURL != nil {
		body["avatar_url"] = *p.avatarURL
	}
	if p.templateSupport != nil {
		body["template_support"] = *p.templateSupport
	}
	if p.providerType != nil {
		body["provider_type"] = string(*p.providerType)
	}
	return body
}

func (p channelMessagePayload) toMap() map[string]interface{} {
	body := map[string]interface{}{}
	if p.id != nil {
		body["id"] = string(*p.id)
	}
	if p.channelID != nil {
		body["channel_id"] = string(*p.channelID)
	}
	if p.senderID != nil {
		body["sender_id"] = *p.senderID
	}
	if p.conversationID != nil {
		body["conversation_id"] = string(*p.conversationID)
	}
	if p.message != nil {
		body["message"] = *p.message
	}
	if p.status != nil {
		body["status"] = string(*p.status)
	}
	if p.createdAt != nil {
		body["created_at"] = formatV1Time(*p.createdAt)
	}
	if p.replyBy != nil {
		body["reply_by"] = formatV1Time(*p.replyBy)
	}
	if p.conversationLink != nil {
		body["conversation_link"] = *p.conversationLink
	}
	if len(p.attachments) > 0 {
		body["attachments"] = p.attachments
	}
	return body
}
