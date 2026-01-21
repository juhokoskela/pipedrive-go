package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestPipelinesService_GetConversionStatistics(t *testing.T) {
	t.Parallel()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/pipelines/9/conversion_statistics" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("start_date"); got != "2024-01-01" {
			t.Fatalf("unexpected start_date: %q", got)
		}
		if got := q.Get("end_date"); got != "2024-01-31" {
			t.Fatalf("unexpected end_date: %q", got)
		}
		if got := q.Get("user_id"); got != "7" {
			t.Fatalf("unexpected user_id: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"stage_conversions":[{"from_stage_id":1,"to_stage_id":2,"conversion_rate":50}],"won_conversion":30,"lost_conversion":20}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	stats, err := client.Pipelines.GetConversionStatistics(
		context.Background(),
		PipelineID(9),
		WithPipelineConversionStartDate(start),
		WithPipelineConversionEndDate(end),
		WithPipelineConversionUserID(UserID(7)),
		WithPipelinesRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("GetConversionStatistics error: %v", err)
	}
	if stats.WonConversion != 30 || len(stats.StageConversions) != 1 {
		t.Fatalf("unexpected stats: %#v", stats)
	}
}

func TestPipelinesService_GetMovementStatistics(t *testing.T) {
	t.Parallel()

	start := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/pipelines/9/movement_statistics" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("start_date"); got != "2024-02-01" {
			t.Fatalf("unexpected start_date: %q", got)
		}
		if got := q.Get("end_date"); got != "2024-02-28" {
			t.Fatalf("unexpected end_date: %q", got)
		}
		if got := q.Get("user_id"); got != "7" {
			t.Fatalf("unexpected user_id: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"movements_between_stages":{"count":1},"new_deals":{"count":1,"deals_ids":[1,2],"values":{"USD":10},"formatted_values":{"USD":"US$10"}},"deals_left_open":{"count":0},"won_deals":{"count":0},"lost_deals":{"count":1,"deals_ids":[3]},"average_age_in_days":{"across_all_stages":2,"by_stages":[{"stage_id":10,"value":15}]}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	stats, err := client.Pipelines.GetMovementStatistics(
		context.Background(),
		PipelineID(9),
		WithPipelineMovementStartDate(start),
		WithPipelineMovementEndDate(end),
		WithPipelineMovementUserID(UserID(7)),
	)
	if err != nil {
		t.Fatalf("GetMovementStatistics error: %v", err)
	}
	if stats.MovementsBetweenStages.Count != 1 || stats.AverageAgeInDays.AcrossAllStages != 2 {
		t.Fatalf("unexpected stats: %#v", stats)
	}
}

func TestPipelinesService_GetConversionStatisticsRequiresDates(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Pipelines.GetConversionStatistics(context.Background(), PipelineID(1))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPipelinesService_ListDeals(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/pipelines/9/deals" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("stage_id"); got != "2" {
			t.Fatalf("unexpected stage_id: %q", got)
		}
		if got := q.Get("get_summary"); got != "1" {
			t.Fatalf("unexpected get_summary: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"title":"Deal"}],"additional_data":{"start":0,"limit":1,"more_items_in_collection":false,"deals_summary":{"count":1}}}`))
	})

	query := url.Values{}
	query.Set("stage_id", "2")
	query.Set("get_summary", "1")
	deals, additional, err := client.Pipelines.ListDeals(context.Background(), PipelineID(9), WithPipelineDealsQuery(query))
	if err != nil {
		t.Fatalf("ListDeals error: %v", err)
	}
	if len(deals) != 1 || deals[0].ID != 1 {
		t.Fatalf("unexpected deals: %#v", deals)
	}
	if additional == nil || additional.Limit != 1 || additional.DealsSummary == nil {
		t.Fatalf("unexpected additional data: %#v", additional)
	}
}
