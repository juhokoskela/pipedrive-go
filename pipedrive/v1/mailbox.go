package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type MailThread struct {
	ID      MailThreadID `json:"id,omitempty"`
	Subject string       `json:"subject,omitempty"`
	Snippet string       `json:"snippet,omitempty"`
	Read    NumberBool   `json:"read_flag,omitempty"`
}

type MailMessage struct {
	ID      MailMessageID `json:"id,omitempty"`
	Subject string        `json:"subject,omitempty"`
	Snippet string        `json:"snippet,omitempty"`
}

type MailboxService struct {
	client *Client
}

type MailboxOption interface {
	applyMailbox(*mailboxOptions)
}

type mailboxOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type mailboxOptionFunc func(*mailboxOptions)

func (f mailboxOptionFunc) applyMailbox(cfg *mailboxOptions) {
	f(cfg)
}

func WithMailboxQuery(values url.Values) MailboxOption {
	return mailboxOptionFunc(func(cfg *mailboxOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithMailboxRequestOptions(opts ...pipedrive.RequestOption) MailboxOption {
	return mailboxOptionFunc(func(cfg *mailboxOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newMailboxOptions(opts []MailboxOption) mailboxOptions {
	var cfg mailboxOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyMailbox(&cfg)
	}
	return cfg
}

func (s *MailboxService) ListThreads(ctx context.Context, opts ...MailboxOption) ([]MailThread, *Pagination, error) {
	cfg := newMailboxOptions(opts)

	var payload struct {
		Data           []MailThread `json:"data"`
		AdditionalData *struct {
			Pagination *Pagination `json:"pagination"`
		} `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/mailbox/mailThreads", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	var page *Pagination
	if payload.AdditionalData != nil {
		page = payload.AdditionalData.Pagination
	}
	return payload.Data, page, nil
}

func (s *MailboxService) GetThread(ctx context.Context, id MailThreadID, opts ...MailboxOption) (*MailThread, error) {
	cfg := newMailboxOptions(opts)
	path := fmt.Sprintf("/mailbox/mailThreads/%d", id)

	var payload struct {
		Data *MailThread `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing mail thread data in response")
	}
	return payload.Data, nil
}

func (s *MailboxService) DeleteThread(ctx context.Context, id MailThreadID, opts ...MailboxOption) (bool, error) {
	cfg := newMailboxOptions(opts)
	path := fmt.Sprintf("/mailbox/mailThreads/%d", id)

	var payload struct {
		Success *bool `json:"success"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return false, err
	}
	if payload.Success == nil {
		return false, fmt.Errorf("missing mail thread delete success in response")
	}
	return *payload.Success, nil
}

func (s *MailboxService) UpdateThread(ctx context.Context, id MailThreadID, form url.Values, opts ...MailboxOption) (*MailThread, error) {
	cfg := newMailboxOptions(opts)
	if len(form) == 0 {
		return nil, fmt.Errorf("form values are required")
	}
	path := fmt.Sprintf("/mailbox/mailThreads/%d", id)

	body := strings.NewReader(form.Encode())
	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	reqOpts = append(reqOpts, pipedrive.WithHeader("Content-Type", "application/x-www-form-urlencoded"))

	var payload struct {
		Data *MailThread `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, body, &payload, reqOpts...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing mail thread update data in response")
	}
	return payload.Data, nil
}

func (s *MailboxService) ListThreadMessages(ctx context.Context, id MailThreadID, opts ...MailboxOption) ([]MailMessage, *Pagination, error) {
	cfg := newMailboxOptions(opts)
	path := fmt.Sprintf("/mailbox/mailThreads/%d/mailMessages", id)

	var payload struct {
		Data           []MailMessage `json:"data"`
		AdditionalData *struct {
			Pagination *Pagination `json:"pagination"`
		} `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	var page *Pagination
	if payload.AdditionalData != nil {
		page = payload.AdditionalData.Pagination
	}
	return payload.Data, page, nil
}

func (s *MailboxService) GetMessage(ctx context.Context, id MailMessageID, opts ...MailboxOption) (*MailMessage, error) {
	cfg := newMailboxOptions(opts)
	path := fmt.Sprintf("/mailbox/mailMessages/%d", id)

	var payload struct {
		Data *MailMessage `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing mail message data in response")
	}
	return payload.Data, nil
}
