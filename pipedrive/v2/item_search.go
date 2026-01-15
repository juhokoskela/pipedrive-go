package v2

import (
	"context"
	"encoding/json"
	"fmt"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type ItemSearchType string

const (
	ItemSearchTypeDeal           ItemSearchType = "deal"
	ItemSearchTypeFile           ItemSearchType = "file"
	ItemSearchTypeLead           ItemSearchType = "lead"
	ItemSearchTypeMailAttachment ItemSearchType = "mail_attachment"
	ItemSearchTypeOrganization   ItemSearchType = "organization"
	ItemSearchTypePerson         ItemSearchType = "person"
	ItemSearchTypeProduct        ItemSearchType = "product"
	ItemSearchTypeProject        ItemSearchType = "project"
)

type ItemSearchField string

const (
	ItemSearchFieldAddress      ItemSearchField = "address"
	ItemSearchFieldCode         ItemSearchField = "code"
	ItemSearchFieldCustomFields ItemSearchField = "custom_fields"
	ItemSearchFieldDescription  ItemSearchField = "description"
	ItemSearchFieldEmail        ItemSearchField = "email"
	ItemSearchFieldName         ItemSearchField = "name"
	ItemSearchFieldNotes        ItemSearchField = "notes"
	ItemSearchFieldPhone        ItemSearchField = "phone"
	ItemSearchFieldTitle        ItemSearchField = "title"
)

type ItemSearchIncludeField string

const (
	ItemSearchIncludeFieldDealCCEmail   ItemSearchIncludeField = "deal.cc_email"
	ItemSearchIncludeFieldPersonPicture ItemSearchIncludeField = "person.picture"
	ItemSearchIncludeFieldProductPrice  ItemSearchIncludeField = "product.price"
)

type ItemSearchEntityType string

const (
	ItemSearchEntityTypeDeal         ItemSearchEntityType = "deal"
	ItemSearchEntityTypeLead         ItemSearchEntityType = "lead"
	ItemSearchEntityTypeOrganization ItemSearchEntityType = "organization"
	ItemSearchEntityTypePerson       ItemSearchEntityType = "person"
	ItemSearchEntityTypeProduct      ItemSearchEntityType = "product"
	ItemSearchEntityTypeProject      ItemSearchEntityType = "project"
)

type ItemSearchMatch string

const (
	ItemSearchMatchBeginning ItemSearchMatch = "beginning"
	ItemSearchMatchExact     ItemSearchMatch = "exact"
	ItemSearchMatchMiddle    ItemSearchMatch = "middle"
)

type ItemSearchItem struct {
	ResultScore float64                `json:"result_score,omitempty"`
	Item        map[string]interface{} `json:"item,omitempty"`
}

type ItemSearchResults struct {
	Items        []ItemSearchItem `json:"items,omitempty"`
	RelatedItems []ItemSearchItem `json:"related_items,omitempty"`
}

type ItemSearchService struct {
	client *Client
}

type SearchItemsOption interface {
	applySearchItems(*searchItemsOptions)
}

type SearchItemsByFieldOption interface {
	applySearchItemsByField(*searchItemsByFieldOptions)
}

type ItemSearchRequestOption interface {
	SearchItemsOption
	SearchItemsByFieldOption
}

type searchItemsOptions struct {
	params         genv2.SearchItemParams
	requestOptions []pipedrive.RequestOption
}

type searchItemsByFieldOptions struct {
	params         genv2.SearchItemByFieldParams
	requestOptions []pipedrive.RequestOption
}

type itemSearchRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o itemSearchRequestOptions) applySearchItems(cfg *searchItemsOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o itemSearchRequestOptions) applySearchItemsByField(cfg *searchItemsByFieldOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type searchItemsOptionFunc func(*searchItemsOptions)

func (f searchItemsOptionFunc) applySearchItems(cfg *searchItemsOptions) {
	f(cfg)
}

type searchItemsByFieldOptionFunc func(*searchItemsByFieldOptions)

func (f searchItemsByFieldOptionFunc) applySearchItemsByField(cfg *searchItemsByFieldOptions) {
	f(cfg)
}

func WithItemSearchRequestOptions(opts ...pipedrive.RequestOption) ItemSearchRequestOption {
	return itemSearchRequestOptions{requestOptions: opts}
}

func WithItemSearchTypes(types ...ItemSearchType) SearchItemsOption {
	return searchItemsOptionFunc(func(cfg *searchItemsOptions) {
		csv := joinCSV(types)
		if csv == "" {
			return
		}
		value := genv2.SearchItemParamsItemTypes(csv)
		cfg.params.ItemTypes = &value
	})
}

func WithItemSearchFields(fields ...ItemSearchField) SearchItemsOption {
	return searchItemsOptionFunc(func(cfg *searchItemsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.SearchItemParamsFields(csv)
		cfg.params.Fields = &value
	})
}

func WithItemSearchIncludeFields(fields ...ItemSearchIncludeField) SearchItemsOption {
	return searchItemsOptionFunc(func(cfg *searchItemsOptions) {
		csv := joinCSV(fields)
		if csv == "" {
			return
		}
		value := genv2.SearchItemParamsIncludeFields(csv)
		cfg.params.IncludeFields = &value
	})
}

func WithItemSearchRelatedItems(enabled bool) SearchItemsOption {
	return searchItemsOptionFunc(func(cfg *searchItemsOptions) {
		cfg.params.SearchForRelatedItems = &enabled
	})
}

func WithItemSearchExactMatch(enabled bool) SearchItemsOption {
	return searchItemsOptionFunc(func(cfg *searchItemsOptions) {
		cfg.params.ExactMatch = &enabled
	})
}

func WithItemSearchPageSize(limit int) SearchItemsOption {
	return searchItemsOptionFunc(func(cfg *searchItemsOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithItemSearchCursor(cursor string) SearchItemsOption {
	return searchItemsOptionFunc(func(cfg *searchItemsOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func WithItemSearchMatch(match ItemSearchMatch) SearchItemsByFieldOption {
	return searchItemsByFieldOptionFunc(func(cfg *searchItemsByFieldOptions) {
		if match == "" {
			return
		}
		value := genv2.SearchItemByFieldParamsMatch(match)
		cfg.params.Match = &value
	})
}

func WithItemSearchByFieldPageSize(limit int) SearchItemsByFieldOption {
	return searchItemsByFieldOptionFunc(func(cfg *searchItemsByFieldOptions) {
		if limit <= 0 {
			return
		}
		cfg.params.Limit = &limit
	})
}

func WithItemSearchByFieldCursor(cursor string) SearchItemsByFieldOption {
	return searchItemsByFieldOptionFunc(func(cfg *searchItemsByFieldOptions) {
		if cursor == "" {
			return
		}
		cfg.params.Cursor = &cursor
	})
}

func newSearchItemsOptions(opts []SearchItemsOption) searchItemsOptions {
	var cfg searchItemsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applySearchItems(&cfg)
	}
	return cfg
}

func newSearchItemsByFieldOptions(opts []SearchItemsByFieldOption) searchItemsByFieldOptions {
	var cfg searchItemsByFieldOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applySearchItemsByField(&cfg)
	}
	return cfg
}

func (s *ItemSearchService) Search(ctx context.Context, term string, opts ...SearchItemsOption) (*ItemSearchResults, *string, error) {
	cfg := newSearchItemsOptions(opts)
	cfg.params.Term = term
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.SearchItemWithResponse(ctx, &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           *ItemSearchResults `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, nil, fmt.Errorf("missing item search data in response")
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}

func (s *ItemSearchService) SearchByField(ctx context.Context, term string, entityType ItemSearchEntityType, field string, opts ...SearchItemsByFieldOption) ([]ItemSearchItem, *string, error) {
	cfg := newSearchItemsByFieldOptions(opts)
	cfg.params.Term = term
	cfg.params.EntityType = genv2.SearchItemByFieldParamsEntityType(entityType)
	cfg.params.Field = field
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, cfg.requestOptions...)

	resp, err := s.client.gen.SearchItemByFieldWithResponse(ctx, &cfg.params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []ItemSearchItem `json:"data"`
		AdditionalData *struct {
			NextCursor *string `json:"next_cursor"`
		} `json:"additional_data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	var next *string
	if payload.AdditionalData != nil {
		next = payload.AdditionalData.NextCursor
	}
	return payload.Data, next, nil
}
