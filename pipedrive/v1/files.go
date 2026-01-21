package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type File struct {
	ID          FileID          `json:"id,omitempty"`
	Name        string          `json:"name,omitempty"`
	FileName    string          `json:"file_name,omitempty"`
	FileType    string          `json:"file_type,omitempty"`
	FileSize    int             `json:"file_size,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	UserID      *UserID         `json:"user_id,omitempty"`
	DealID      *DealID         `json:"deal_id,omitempty"`
	PersonID    *PersonID       `json:"person_id,omitempty"`
	OrgID       *OrganizationID `json:"org_id,omitempty"`
	ProductID   *ProductID      `json:"product_id,omitempty"`
	ActivityID  *ActivityID     `json:"activity_id,omitempty"`
	LeadID      *LeadID         `json:"lead_id,omitempty"`
	Active      bool            `json:"active_flag,omitempty"`
	Inline      bool            `json:"inline_flag,omitempty"`
	AddTime     *DateTime       `json:"add_time,omitempty"`
	UpdateTime  *DateTime       `json:"update_time,omitempty"`
}

type FilesService struct {
	client *Client
}

type FilesOption interface {
	applyFiles(*filesOptions)
}

type filesOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type filesOptionFunc func(*filesOptions)

func (f filesOptionFunc) applyFiles(cfg *filesOptions) {
	f(cfg)
}

func WithFilesQuery(values url.Values) FilesOption {
	return filesOptionFunc(func(cfg *filesOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithFilesRequestOptions(opts ...pipedrive.RequestOption) FilesOption {
	return filesOptionFunc(func(cfg *filesOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newFilesOptions(opts []FilesOption) filesOptions {
	var cfg filesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyFiles(&cfg)
	}
	return cfg
}

func (s *FilesService) List(ctx context.Context, opts ...FilesOption) ([]File, *Pagination, error) {
	cfg := newFilesOptions(opts)

	var payload struct {
		Data           []File `json:"data"`
		AdditionalData *struct {
			Pagination *Pagination `json:"pagination"`
		} `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/files", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	var page *Pagination
	if payload.AdditionalData != nil {
		page = payload.AdditionalData.Pagination
	}
	return payload.Data, page, nil
}

func (s *FilesService) Get(ctx context.Context, id FileID, opts ...FilesOption) (*File, error) {
	cfg := newFilesOptions(opts)
	path := fmt.Sprintf("/files/%d", id)

	var payload struct {
		Data *File `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing file data in response")
	}
	return payload.Data, nil
}

func (s *FilesService) Add(ctx context.Context, body io.Reader, contentType string, opts ...FilesOption) (*File, error) {
	cfg := newFilesOptions(opts)
	if body == nil {
		return nil, fmt.Errorf("file body is required")
	}
	if contentType == "" {
		return nil, fmt.Errorf("content type is required")
	}

	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	reqOpts = append(reqOpts, pipedrive.WithHeader("Content-Type", contentType))

	var payload struct {
		Data *File `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, "/files", cfg.query, body, &payload, reqOpts...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing file data in response")
	}
	return payload.Data, nil
}

func (s *FilesService) AddRemoteFile(ctx context.Context, form url.Values, opts ...FilesOption) (*File, error) {
	cfg := newFilesOptions(opts)
	if len(form) == 0 {
		return nil, fmt.Errorf("form values are required")
	}

	body := strings.NewReader(form.Encode())
	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	reqOpts = append(reqOpts, pipedrive.WithHeader("Content-Type", "application/x-www-form-urlencoded"))

	var payload struct {
		Data *File `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, "/files/remote", cfg.query, body, &payload, reqOpts...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing remote file data in response")
	}
	return payload.Data, nil
}

func (s *FilesService) LinkRemoteFile(ctx context.Context, form url.Values, opts ...FilesOption) (*File, error) {
	cfg := newFilesOptions(opts)
	if len(form) == 0 {
		return nil, fmt.Errorf("form values are required")
	}

	body := strings.NewReader(form.Encode())
	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	reqOpts = append(reqOpts, pipedrive.WithHeader("Content-Type", "application/x-www-form-urlencoded"))

	var payload struct {
		Data *File `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, "/files/remoteLink", cfg.query, body, &payload, reqOpts...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing linked file data in response")
	}
	return payload.Data, nil
}

func (s *FilesService) Update(ctx context.Context, id FileID, body io.Reader, contentType string, opts ...FilesOption) (*File, error) {
	cfg := newFilesOptions(opts)
	if body == nil {
		return nil, fmt.Errorf("file body is required")
	}
	if contentType == "" {
		return nil, fmt.Errorf("content type is required")
	}

	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	reqOpts = append(reqOpts, pipedrive.WithHeader("Content-Type", contentType))

	path := fmt.Sprintf("/files/%d", id)
	var payload struct {
		Data *File `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, body, &payload, reqOpts...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing file update data in response")
	}
	return payload.Data, nil
}

func (s *FilesService) Delete(ctx context.Context, id FileID, opts ...FilesOption) (bool, error) {
	cfg := newFilesOptions(opts)
	path := fmt.Sprintf("/files/%d", id)

	var payload struct {
		Success *bool `json:"success"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return false, err
	}
	if payload.Success == nil {
		return false, fmt.Errorf("missing file delete success in response")
	}
	return *payload.Success, nil
}

func (s *FilesService) Download(ctx context.Context, id FileID, opts ...FilesOption) ([]byte, error) {
	cfg := newFilesOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.DownloadFile(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errorFromResponse(resp, body)
	}
	return body, nil
}
