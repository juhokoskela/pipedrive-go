package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

type CurrencyID int64

type Currency struct {
	ID            CurrencyID `json:"id"`
	Code          string     `json:"code,omitempty"`
	Name          string     `json:"name,omitempty"`
	Symbol        string     `json:"symbol,omitempty"`
	Active        bool       `json:"active_flag,omitempty"`
	IsCustom      bool       `json:"is_custom_flag,omitempty"`
	DecimalPoints int        `json:"decimal_points,omitempty"`
}

type ListCurrenciesRequest struct {
	Term string
}

type CurrenciesService struct {
	client *Client
}

func (s *CurrenciesService) List(ctx context.Context, req ListCurrenciesRequest, opts ...pipedrive.RequestOption) ([]Currency, error) {
	ctx, editors := pipedrive.ApplyRequestOptions(ctx, opts...)

	var params *genv1.GetCurrenciesParams
	if req.Term != "" {
		term := req.Term
		params = &genv1.GetCurrenciesParams{Term: &term}
	}

	resp, err := s.client.gen.GetCurrenciesWithResponse(ctx, params, toRequestEditors(editors)...)
	if err != nil {
		return nil, err
	}
	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode > 299 {
		return nil, errorFromResponse(resp.HTTPResponse, resp.Body)
	}

	var payload struct {
		Data []Currency `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}

func errorFromResponse(httpResp *http.Response, body []byte) error {
	if httpResp.StatusCode == http.StatusTooManyRequests {
		return pipedrive.RateLimitErrorFromResponse(httpResp, body, time.Now())
	}
	return pipedrive.APIErrorFromResponse(httpResp, body)
}

func toRequestEditors(editors []pipedrive.RequestEditorFunc) []genv1.RequestEditorFn {
	out := make([]genv1.RequestEditorFn, 0, len(editors))
	for _, editor := range editors {
		if editor == nil {
			continue
		}
		out = append(out, genv1.RequestEditorFn(editor))
	}
	return out
}

