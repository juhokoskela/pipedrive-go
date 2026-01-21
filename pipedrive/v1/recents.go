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

type RecentsItemType string

const (
	RecentsItemActivity     RecentsItemType = "activity"
	RecentsItemActivityType RecentsItemType = "activityType"
	RecentsItemDeal         RecentsItemType = "deal"
	RecentsItemFile         RecentsItemType = "file"
	RecentsItemFilter       RecentsItemType = "filter"
	RecentsItemNote         RecentsItemType = "note"
	RecentsItemPerson       RecentsItemType = "person"
	RecentsItemOrganization RecentsItemType = "organization"
	RecentsItemPipeline     RecentsItemType = "pipeline"
	RecentsItemProduct      RecentsItemType = "product"
	RecentsItemStage        RecentsItemType = "stage"
	RecentsItemUser         RecentsItemType = "user"
)

type Recent struct {
	Item RecentsItemType `json:"item,omitempty"`
	ID   int64           `json:"id,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

type RecentsPagination struct {
	Start                 int  `json:"start,omitempty"`
	Limit                 int  `json:"limit,omitempty"`
	MoreItemsInCollection bool `json:"more_items_in_collection,omitempty"`
}

type RecentsAdditionalData struct {
	LastTimestampOnPage string             `json:"last_timestamp_on_page,omitempty"`
	SinceTimestamp      string             `json:"since_timestamp,omitempty"`
	Pagination          *RecentsPagination `json:"pagination,omitempty"`
}

type RecentsService struct {
	client *Client
}

type ListRecentsOption interface {
	applyListRecents(*listRecentsOptions)
}

type RecentsRequestOption interface {
	ListRecentsOption
}

type listRecentsOptions struct {
	params         genv1.GetRecentsParams
	requestOptions []pipedrive.RequestOption
}

type recentsRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o recentsRequestOptions) applyListRecents(cfg *listRecentsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type listRecentsOptionFunc func(*listRecentsOptions)

func (f listRecentsOptionFunc) applyListRecents(cfg *listRecentsOptions) {
	f(cfg)
}

func WithRecentsRequestOptions(opts ...pipedrive.RequestOption) RecentsRequestOption {
	return recentsRequestOptions{requestOptions: opts}
}

func WithRecentsSince(since time.Time) ListRecentsOption {
	return listRecentsOptionFunc(func(cfg *listRecentsOptions) {
		cfg.params.SinceTimestamp = formatV1Time(since)
	})
}

func WithRecentsItems(items ...RecentsItemType) ListRecentsOption {
	return listRecentsOptionFunc(func(cfg *listRecentsOptions) {
		if len(items) == 0 {
			return
		}
		values := make([]string, 0, len(items))
		for _, item := range items {
			if item == "" {
				continue
			}
			values = append(values, string(item))
		}
		if len(values) == 0 {
			return
		}
		csv := strings.Join(values, ",")
		value := genv1.GetRecentsParamsItems(csv)
		cfg.params.Items = &value
	})
}

func WithRecentsStart(start int) ListRecentsOption {
	return listRecentsOptionFunc(func(cfg *listRecentsOptions) {
		cfg.params.Start = &start
	})
}

func WithRecentsLimit(limit int) ListRecentsOption {
	return listRecentsOptionFunc(func(cfg *listRecentsOptions) {
		cfg.params.Limit = &limit
	})
}

func newListRecentsOptions(opts []ListRecentsOption) listRecentsOptions {
	var cfg listRecentsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyListRecents(&cfg)
	}
	return cfg
}

func (s *RecentsService) List(ctx context.Context, opts ...ListRecentsOption) ([]Recent, *RecentsAdditionalData, error) {
	cfg := newListRecentsOptions(opts)
	if cfg.params.SinceTimestamp == "" {
		return nil, nil, fmt.Errorf("since timestamp is required")
	}
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.GetRecents(ctx, &cfg.params, toRequestEditors(editors)...)
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
		Data           []Recent               `json:"data"`
		AdditionalData *RecentsAdditionalData `json:"additional_data"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, payload.AdditionalData, nil
}
