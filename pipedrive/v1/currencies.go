package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

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

	resp, err := s.client.gen.GetCurrencies(ctx, params, toRequestEditors(editors)...)
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
		Data []Currency `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return payload.Data, nil
}
