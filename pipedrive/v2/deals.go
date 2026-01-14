package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type DealID int64

type Deal struct {
	ID           DealID                 `json:"id"`
	Title        string                 `json:"title,omitempty"`
	Currency     string                 `json:"currency,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

type ListDealsRequest struct {
	Limit  int
	Cursor string
}

type DealsService struct {
	client *Client
}

func (s *DealsService) Get(ctx context.Context, id DealID, opts ...pipedrive.RequestOption) (*Deal, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, opts...)
	resp, err := s.client.gen.GetDealWithResponse(ctx, int(id), nil, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data *Deal `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("missing deal data in response")
	}
	return payload.Data, nil
}

func (s *DealsService) List(ctx context.Context, req ListDealsRequest, opts ...pipedrive.RequestOption) ([]Deal, *string, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, opts...)

	params := &genv2.GetDealsParams{}
	if req.Limit > 0 {
		limit := req.Limit
		params.Limit = &limit
	}
	if req.Cursor != "" {
		cursor := req.Cursor
		params.Cursor = &cursor
	}

	resp, err := s.client.gen.GetDealsWithResponse(ctx, params, toRequestEditors(editors)...)
	if err != nil {
		return nil, nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data           []Deal `json:"data"`
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

func (s *DealsService) ListPager(req ListDealsRequest, opts ...pipedrive.RequestOption) *pipedrive.CursorPager[Deal] {
	startCursor := req.Cursor
	req.Cursor = ""

	return pipedrive.NewCursorPager(func(ctx context.Context, cursor *string) ([]Deal, *string, error) {
		effective := cursor
		if effective == nil && startCursor != "" {
			effective = &startCursor
		}

		pageReq := req
		if effective != nil {
			pageReq.Cursor = *effective
		}
		return s.List(ctx, pageReq, opts...)
	})
}

func errorFromResponse(httpResp *http.Response, body []byte) error {
	if httpResp.StatusCode == http.StatusTooManyRequests {
		return pipedrive.RateLimitErrorFromResponse(httpResp, body, time.Now())
	}
	return pipedrive.APIErrorFromResponse(httpResp, body)
}

func toRequestEditors(editors []pipedrive.RequestEditorFunc) []genv2.RequestEditorFn {
	out := make([]genv2.RequestEditorFn, 0, len(editors))
	for _, editor := range editors {
		if editor == nil {
			continue
		}
		out = append(out, genv2.RequestEditorFn(editor))
	}
	return out
}
