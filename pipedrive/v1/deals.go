package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type Deal struct {
	ID         DealID          `json:"id,omitempty"`
	Title      string          `json:"title,omitempty"`
	Status     string          `json:"status,omitempty"`
	Value      float64         `json:"value,omitempty"`
	Currency   string          `json:"currency,omitempty"`
	OwnerID    *UserID         `json:"user_id,omitempty"`
	PersonID   *PersonID       `json:"person_id,omitempty"`
	OrgID      *OrganizationID `json:"org_id,omitempty"`
	StageID    *StageID        `json:"stage_id,omitempty"`
	PipelineID *PipelineID     `json:"pipeline_id,omitempty"`
	AddTime    *DateTime       `json:"add_time,omitempty"`
	UpdateTime *DateTime       `json:"update_time,omitempty"`
}

type DealsSummaryValues struct {
	Count                   int     `json:"count,omitempty"`
	Value                   float64 `json:"value,omitempty"`
	ValueConverted          float64 `json:"value_converted,omitempty"`
	ValueConvertedFormatted string  `json:"value_converted_formatted,omitempty"`
	ValueFormatted          string  `json:"value_formatted,omitempty"`
}

type DealsSummaryWeightedValues struct {
	Count          int     `json:"count,omitempty"`
	Value          float64 `json:"value,omitempty"`
	ValueFormatted string  `json:"value_formatted,omitempty"`
}

type DealsSummary struct {
	TotalCount                                   int                         `json:"total_count,omitempty"`
	TotalCurrencyConvertedValue                  float64                     `json:"total_currency_converted_value,omitempty"`
	TotalCurrencyConvertedValueFormatted         string                      `json:"total_currency_converted_value_formatted,omitempty"`
	TotalWeightedCurrencyConvertedValue          float64                     `json:"total_weighted_currency_converted_value,omitempty"`
	TotalWeightedCurrencyConvertedValueFormatted string                      `json:"total_weighted_currency_converted_value_formatted,omitempty"`
	ValuesTotal                                  *DealsSummaryValues         `json:"values_total,omitempty"`
	WeightedValuesTotal                          *DealsSummaryWeightedValues `json:"weighted_values_total,omitempty"`
}

type DealsTimeline map[string]any

type DealsService struct {
	client *Client
}

type DealsOption interface {
	applyDeals(*dealsOptions)
}

type dealsOptions struct {
	query          url.Values
	requestOptions []pipedrive.RequestOption
}

type dealsOptionFunc func(*dealsOptions)

func (f dealsOptionFunc) applyDeals(cfg *dealsOptions) {
	f(cfg)
}

func WithDealsQuery(values url.Values) DealsOption {
	return dealsOptionFunc(func(cfg *dealsOptions) {
		cfg.query = mergeQueryValues(cfg.query, values)
	})
}

func WithDealsRequestOptions(opts ...pipedrive.RequestOption) DealsOption {
	return dealsOptionFunc(func(cfg *dealsOptions) {
		cfg.requestOptions = append(cfg.requestOptions, opts...)
	})
}

func newDealsOptions(opts []DealsOption) dealsOptions {
	var cfg dealsOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyDeals(&cfg)
	}
	return cfg
}

func (s *DealsService) ListCollection(ctx context.Context, opts ...DealsOption) ([]Deal, *CollectionPagination, error) {
	cfg := newDealsOptions(opts)

	var payload struct {
		Data           []Deal                `json:"data"`
		AdditionalData *CollectionPagination `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/deals/collection", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *DealsService) Summary(ctx context.Context, opts ...DealsOption) (*DealsSummary, error) {
	cfg := newDealsOptions(opts)

	var payload struct {
		Data *DealsSummary `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/deals/summary", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deals summary data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) ArchivedSummary(ctx context.Context, opts ...DealsOption) (*DealsSummary, error) {
	cfg := newDealsOptions(opts)

	var payload struct {
		Data *DealsSummary `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/deals/summary/archived", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing archived deals summary data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) Timeline(ctx context.Context, opts ...DealsOption) (DealsTimeline, error) {
	cfg := newDealsOptions(opts)

	var payload struct {
		Data DealsTimeline `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/deals/timeline", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *DealsService) ArchivedTimeline(ctx context.Context, opts ...DealsOption) (DealsTimeline, error) {
	cfg := newDealsOptions(opts)

	var payload struct {
		Data DealsTimeline `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, "/deals/timeline/archived", cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *DealsService) ListActivities(ctx context.Context, id DealID, opts ...DealsOption) ([]Activity, *Pagination, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/activities", id)

	var payload struct {
		Data           []Activity `json:"data"`
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

func (s *DealsService) Changelog(ctx context.Context, id DealID, opts ...DealsOption) ([]map[string]any, *CollectionPagination, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/changelog", id)

	var payload struct {
		Data           []map[string]any      `json:"data"`
		AdditionalData *CollectionPagination `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *DealsService) ListFiles(ctx context.Context, id DealID, opts ...DealsOption) ([]File, *Pagination, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/files", id)

	var payload struct {
		Data           []File `json:"data"`
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

func (s *DealsService) ListMailMessages(ctx context.Context, id DealID, opts ...DealsOption) ([]MailMessage, *Pagination, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/mailMessages", id)

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

func (s *DealsService) ListParticipants(ctx context.Context, id DealID, opts ...DealsOption) ([]Person, *Pagination, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/participants", id)

	var payload struct {
		Data           []Person `json:"data"`
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

func (s *DealsService) AddParticipant(ctx context.Context, id DealID, personID PersonID, opts ...DealsOption) (*Person, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/participants", id)

	body := map[string]any{"person_id": int(personID)}
	var payload struct {
		Data *Person `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, path, cfg.query, body, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal participant data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) DeleteParticipant(ctx context.Context, id DealID, participantID DealParticipantID, opts ...DealsOption) (bool, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/participants/%d", id, participantID)

	var payload struct {
		Success *bool `json:"success"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return false, err
	}
	if payload.Success == nil {
		return false, fmt.Errorf("missing deal participant delete success in response")
	}
	return *payload.Success, nil
}

func (s *DealsService) ParticipantsChangelog(ctx context.Context, id DealID, opts ...DealsOption) ([]map[string]any, *CollectionPagination, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/participantsChangelog", id)

	var payload struct {
		Data           []map[string]any      `json:"data"`
		AdditionalData *CollectionPagination `json:"additional_data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, nil, err
	}
	return payload.Data, payload.AdditionalData, nil
}

func (s *DealsService) ListPersons(ctx context.Context, id DealID, opts ...DealsOption) ([]Person, *Pagination, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/persons", id)

	var payload struct {
		Data           []Person `json:"data"`
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

func (s *DealsService) ListUpdates(ctx context.Context, id DealID, opts ...DealsOption) ([]map[string]any, *Pagination, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/flow", id)

	var payload struct {
		Data           []map[string]any `json:"data"`
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

func (s *DealsService) ListUsers(ctx context.Context, id DealID, opts ...DealsOption) ([]User, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/permittedUsers", id)

	var payload struct {
		Data []User `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodGet, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (s *DealsService) Merge(ctx context.Context, id DealID, mergeWithID DealID, opts ...DealsOption) (*Deal, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/merge", id)

	body := map[string]any{"merge_with_id": int(mergeWithID)}
	var payload struct {
		Data *Deal `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPut, path, cfg.query, body, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing merged deal data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) Duplicate(ctx context.Context, id DealID, opts ...DealsOption) (*Deal, error) {
	cfg := newDealsOptions(opts)
	path := fmt.Sprintf("/deals/%d/duplicate", id)

	var payload struct {
		Data *Deal `json:"data"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodPost, path, cfg.query, nil, &payload, cfg.requestOptions...); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing duplicated deal data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) Delete(ctx context.Context, ids []DealID, opts ...DealsOption) (bool, error) {
	cfg := newDealsOptions(opts)
	query := mergeQueryValues(url.Values{}, cfg.query)
	if len(ids) == 0 {
		return false, fmt.Errorf("at least one deal id is required")
	}
	query.Set("ids", joinIDs(ids))

	var payload struct {
		Success *bool `json:"success"`
	}
	if err := s.client.Raw.Do(ctx, http.MethodDelete, "/deals", query, nil, &payload, cfg.requestOptions...); err != nil {
		return false, err
	}
	if payload.Success == nil {
		return false, fmt.Errorf("missing deal delete success in response")
	}
	return *payload.Success, nil
}
