package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/juhokoskela/pipedrive-go/internal/multipartbody"
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
	ProjectID   *ProjectID      `json:"project_id,omitempty"`
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

type UploadFileOption interface {
	applyUploadFile(*uploadFileOptions)
}

type filesOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type uploadFileOptions struct {
	payload        uploadFilePayload
	requestOptions []pipedrive.RequestOption
}

type uploadFilePayload struct {
	dealID     *DealID
	personID   *PersonID
	orgID      *OrganizationID
	productID  *ProductID
	activityID *ActivityID
	leadID     *LeadID
	projectID  *ProjectID
}

type filesOptionFunc func(*filesOptions)

func (f filesOptionFunc) applyFiles(cfg *filesOptions) {
	f(cfg)
}

type uploadFileOptionFunc func(*uploadFileOptions)

func (f uploadFileOptionFunc) applyUploadFile(cfg *uploadFileOptions) {
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

func WithUploadFileRequestOptions(opts ...pipedrive.RequestOption) UploadFileOption {
	return uploadFileOptionFunc(func(cfg *uploadFileOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func WithFileDealID(id DealID) UploadFileOption {
	return uploadFileOptionFunc(func(cfg *uploadFileOptions) {
		cfg.payload.dealID = &id
	})
}

func WithFilePersonID(id PersonID) UploadFileOption {
	return uploadFileOptionFunc(func(cfg *uploadFileOptions) {
		cfg.payload.personID = &id
	})
}

func WithFileOrganizationID(id OrganizationID) UploadFileOption {
	return uploadFileOptionFunc(func(cfg *uploadFileOptions) {
		cfg.payload.orgID = &id
	})
}

func WithFileProductID(id ProductID) UploadFileOption {
	return uploadFileOptionFunc(func(cfg *uploadFileOptions) {
		cfg.payload.productID = &id
	})
}

func WithFileActivityID(id ActivityID) UploadFileOption {
	return uploadFileOptionFunc(func(cfg *uploadFileOptions) {
		cfg.payload.activityID = &id
	})
}

func WithFileLeadID(id LeadID) UploadFileOption {
	return uploadFileOptionFunc(func(cfg *uploadFileOptions) {
		cfg.payload.leadID = &id
	})
}

func WithFileProjectID(id ProjectID) UploadFileOption {
	return uploadFileOptionFunc(func(cfg *uploadFileOptions) {
		cfg.payload.projectID = &id
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

func newUploadFileOptions(opts []UploadFileOption) uploadFileOptions {
	var cfg uploadFileOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyUploadFile(&cfg)
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
	if err := validateID(id, "file id"); err != nil {
		return nil, err
	}
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

// Upload adds a file by encoding content as a multipart/form-data body.
// Unlike Add, the encoded body is replayable, so uploads participate in
// retries even when content itself is not seekable.
func (s *FilesService) Upload(ctx context.Context, fileName string, content io.Reader, opts ...UploadFileOption) (*File, error) {
	if fileName == "" || content == nil {
		return nil, fmt.Errorf("file name and content are required")
	}
	cfg := newUploadFileOptions(opts)
	fields, err := cfg.payload.toMultipartFields()
	if err != nil {
		return nil, err
	}

	contentType, body, err := multipartbody.NewFileWithFields("file", fileName, content, fields)
	if err != nil {
		return nil, err
	}

	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	reqOpts = append(reqOpts, pipedrive.WithHeader("Content-Type", contentType))

	var payload struct {
		Data *File `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, "/files", nil, body, &payload, reqOpts...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing file data in response")
	}
	return payload.Data, nil
}

func (p uploadFilePayload) toMultipartFields() (url.Values, error) {
	fields := make(url.Values)
	if p.dealID != nil {
		if err := validateID(*p.dealID, "deal id"); err != nil {
			return nil, err
		}
		fields.Set("deal_id", strconv.FormatInt(int64(*p.dealID), 10))
	}
	if p.personID != nil {
		if err := validateID(*p.personID, "person id"); err != nil {
			return nil, err
		}
		fields.Set("person_id", strconv.FormatInt(int64(*p.personID), 10))
	}
	if p.orgID != nil {
		if err := validateID(*p.orgID, "organization id"); err != nil {
			return nil, err
		}
		fields.Set("org_id", strconv.FormatInt(int64(*p.orgID), 10))
	}
	if p.productID != nil {
		if err := validateID(*p.productID, "product id"); err != nil {
			return nil, err
		}
		fields.Set("product_id", strconv.FormatInt(int64(*p.productID), 10))
	}
	if p.activityID != nil {
		if err := validateID(*p.activityID, "activity id"); err != nil {
			return nil, err
		}
		fields.Set("activity_id", strconv.FormatInt(int64(*p.activityID), 10))
	}
	if p.leadID != nil {
		if _, err := parseUUID(string(*p.leadID), "lead id"); err != nil {
			return nil, err
		}
		fields.Set("lead_id", string(*p.leadID))
	}
	if p.projectID != nil {
		if err := validateID(*p.projectID, "project id"); err != nil {
			return nil, err
		}
		fields.Set("project_id", strconv.FormatInt(int64(*p.projectID), 10))
	}
	return fields, nil
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
	if err := validateID(id, "file id"); err != nil {
		return nil, err
	}
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
	if err := validateID(id, "file id"); err != nil {
		return false, err
	}
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
	if err := validateID(id, "file id"); err != nil {
		return nil, err
	}
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

func (s *FilesService) DownloadTo(ctx context.Context, id FileID, dst io.Writer, opts ...FilesOption) error {
	if err := validateID(id, "file id"); err != nil {
		return err
	}
	if dst == nil {
		return fmt.Errorf("download destination is required")
	}

	cfg := newFilesOptions(opts)
	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	reqOpts = append(reqOpts, pipedrive.WithNoResponseSizeLimit())
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, reqOpts...)

	resp, err := s.client.gen.DownloadFile(ctx, int(id), toRequestEditors(editors)...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			return fmt.Errorf("read response: %w", err)
		}
		return errorFromResponse(resp, body)
	}
	if _, err := io.Copy(dst, resp.Body); err != nil {
		return fmt.Errorf("write download: %w", err)
	}
	return nil
}
