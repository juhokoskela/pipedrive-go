package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Activity struct {
	ID                ActivityID      `json:"id"`
	Subject           string          `json:"subject,omitempty"`
	Type              string          `json:"type,omitempty"`
	Done              bool            `json:"done,omitempty"`
	Busy              bool            `json:"busy_flag,omitempty"`
	Active            bool            `json:"active_flag,omitempty"`
	DueDate           string          `json:"due_date,omitempty"`
	DueTime           string          `json:"due_time,omitempty"`
	Duration          string          `json:"duration,omitempty"`
	UserID            *UserID         `json:"user_id,omitempty"`
	UpdateUserID      *UserID         `json:"update_user_id,omitempty"`
	DealID            *DealID         `json:"deal_id,omitempty"`
	PersonID          *PersonID       `json:"person_id,omitempty"`
	OrganizationID    *OrganizationID `json:"org_id,omitempty"`
	LeadID            *LeadID         `json:"lead_id,omitempty"`
	ProjectID         *ProjectID      `json:"project_id,omitempty"`
	PublicDescription string          `json:"public_description,omitempty"`
	SourceTimezone    string          `json:"source_timezone,omitempty"`
	AddTime           *DateTime       `json:"add_time,omitempty"`
	UpdateTime        *DateTime       `json:"update_time,omitempty"`
	MarkedAsDoneTime  *DateTime       `json:"marked_as_done_time,omitempty"`
}

type ActivitiesDeleteResult struct {
	IDs []ActivityID `json:"id"`
}

type ActivitiesService struct {
	client *Client
}

type ListActivitiesCollectionOption interface {
	applyListActivitiesCollection(*listActivitiesCollectionOptions)
}

type DeleteActivitiesOption interface {
	applyDeleteActivities(*deleteActivitiesOptions)
}

type ActivitiesRequestOption interface {
	ListActivitiesCollectionOption
	DeleteActivitiesOption
}

type listActivitiesCollectionOptions struct {
	params         genv1.GetActivitiesCollectionParams
	requestOptions []pipedrive.RequestOption
}

type deleteActivitiesOptions struct {
	requestOptions []pipedrive.RequestOption
}

type activitiesRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o activitiesRequestOptions) applyListActivitiesCollection(cfg *listActivitiesCollectionOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o activitiesRequestOptions) applyDeleteActivities(cfg *deleteActivitiesOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listActivitiesCollectionOptionFunc func(*listActivitiesCollectionOptions)

func (f listActivitiesCollectionOptionFunc) applyListActivitiesCollection(cfg *listActivitiesCollectionOptions) {
	f(cfg)
}

func WithActivitiesRequestOptions(opts ...pipedrive.RequestOption) ActivitiesRequestOption {
	return activitiesRequestOptions{requestOptions: opts}
}

func WithActivitiesCollectionCursor(cursor string) ListActivitiesCollectionOption {
	return listActivitiesCollectionOptionFunc(func(cfg *listActivitiesCollectionOptions) {
		cfg.params.Cursor = &cursor
	})
}

func WithActivitiesCollectionLimit(limit int) ListActivitiesCollectionOption {
	return listActivitiesCollectionOptionFunc(func(cfg *listActivitiesCollectionOptions) {
		cfg.params.Limit = &limit
	})
}

func WithActivitiesCollectionSince(since time.Time) ListActivitiesCollectionOption {
	return listActivitiesCollectionOptionFunc(func(cfg *listActivitiesCollectionOptions) {
		value := formatV1Time(since)
		cfg.params.Since = &value
	})
}

func WithActivitiesCollectionUntil(until time.Time) ListActivitiesCollectionOption {
	return listActivitiesCollectionOptionFunc(func(cfg *listActivitiesCollectionOptions) {
		value := formatV1Time(until)
		cfg.params.Until = &value
	})
}

func WithActivitiesCollectionUserID(userID UserID) ListActivitiesCollectionOption {
	return listActivitiesCollectionOptionFunc(func(cfg *listActivitiesCollectionOptions) {
		value := int(userID)
		cfg.params.UserId = &value
	})
}

func WithActivitiesCollectionDone(done bool) ListActivitiesCollectionOption {
	return listActivitiesCollectionOptionFunc(func(cfg *listActivitiesCollectionOptions) {
		cfg.params.Done = &done
	})
}

func WithActivitiesCollectionTypes(types ...string) ListActivitiesCollectionOption {
	return listActivitiesCollectionOptionFunc(func(cfg *listActivitiesCollectionOptions) {
		if len(types) == 0 {
			return
		}
		csv := strings.Join(types, ",")
		cfg.params.Type = &csv
	})
}

func newListActivitiesCollectionOptions(opts []ListActivitiesCollectionOption) listActivitiesCollectionOptions {
	var cfg listActivitiesCollectionOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListActivitiesCollection(&cfg)
	}
	return cfg
}

func newDeleteActivitiesOptions(opts []DeleteActivitiesOption) deleteActivitiesOptions {
	var cfg deleteActivitiesOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeleteActivities(&cfg)
	}
	return cfg
}

func (s *ActivitiesService) ListCollection(ctx context.Context, opts ...ListActivitiesCollectionOption) ([]Activity, *string, error) {
	cfg := newListActivitiesCollectionOptions(opts)
	return s.listCollection(ctx, cfg.params, cfg.requestOptions)
}

func (s *ActivitiesService) ListCollectionPager(opts ...ListActivitiesCollectionOption) *pipedrive.CursorPager[Activity] {
	cfg := newListActivitiesCollectionOptions(opts)
	startCursor := cfg.params.Cursor
	cfg.params.Cursor = nil

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Activity, *string, error) {
		params := cfg.params
		if cursor != nil {
			params.Cursor = cursor
		} else if startCursor != nil {
			params.Cursor = startCursor
		}
		return s.listCollection(ctx, params, cfg.requestOptions)
	})
}

func (s *ActivitiesService) ForEachCollection(ctx context.Context, fn func(Activity) error, opts ...ListActivitiesCollectionOption) error {
	return s.ListCollectionPager(opts...).ForEach(ctx, fn)
}

func (s *ActivitiesService) Delete(ctx context.Context, ids []ActivityID, opts ...DeleteActivitiesOption) (*ActivitiesDeleteResult, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("activity IDs are required")
	}
	cfg := newDeleteActivitiesOptions(opts)
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	params := genv1.DeleteActivitiesParams{Ids: joinIDs(ids)}
	resp, err := s.client.gen.DeleteActivities(ctx, &params, toRequestEditors(editors)...)
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

	var payload struct {
		Data *struct {
			IDs []ActivityID `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing delete activities data in response")
	}
	return &ActivitiesDeleteResult{IDs: payload.Data.IDs}, nil
}

func (s *ActivitiesService) listCollection(ctx context.Context, params genv1.GetActivitiesCollectionParams, requestOptions []pipedrive.RequestOption) ([]Activity, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, requestOptions...)

	resp, err := s.client.gen.GetActivitiesCollection(ctx, &params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp, body)
	}

	var payload struct {
		Data           []Activity `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}
