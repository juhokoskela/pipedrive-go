package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type NoteFieldsService struct {
	client *Client
}

type ListNoteFieldsOption interface {
	applyListNoteFields(*listNoteFieldsOptions)
}

type NoteFieldsRequestOption interface {
	ListNoteFieldsOption
}

type listNoteFieldsOptions struct {
	requestOptions []pipedrive.RequestOption
}

type noteFieldsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o noteFieldsRequestOptions) applyListNoteFields(cfg *listNoteFieldsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listNoteFieldsOptionFunc func(*listNoteFieldsOptions)

func (f listNoteFieldsOptionFunc) applyListNoteFields(cfg *listNoteFieldsOptions) {
	f(cfg)
}

func WithNoteFieldsRequestOptions(opts ...pipedrive.RequestOption) NoteFieldsRequestOption {
	return noteFieldsRequestOptions{requestOptions: opts}
}

func newListNoteFieldsOptions(opts []ListNoteFieldsOption) listNoteFieldsOptions {
	var cfg listNoteFieldsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListNoteFields(&cfg)
	}
	return cfg
}

func (s *NoteFieldsService) List(ctx context.Context, opts ...ListNoteFieldsOption) ([]Field, *FieldPagination, error) {
	cfg := newListNoteFieldsOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetNoteFields(ctx, toRequestEditors(editors)...)
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
		Data           []Field          `json:"data"`
		AdditionalData *FieldPagination `json:"additional_data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, payload.AdditionalData, nil
}
