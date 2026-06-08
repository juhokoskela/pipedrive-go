package v1

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestCallLogsService_List(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/callLogs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("start"); got != "0" {
			t.Fatalf("unexpected start: %q", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":"log-1","outcome":"connected","to_phone_number":"+123","start_time":"2022-12-12 01:01:01","end_time":"2022-12-12 01:02:01","user_id":7,"company_id":9}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false,"next_start":1}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	logs, page, err := client.CallLogs.List(
		context.Background(),
		WithCallLogsStart(0),
		WithCallLogsLimit(1),
		WithCallLogsRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(logs) != 1 || logs[0].ID != "log-1" {
		t.Fatalf("unexpected logs: %#v", logs)
	}
	if logs[0].StartTime == nil || logs[0].StartTime.IsZero() {
		t.Fatalf("expected start_time to be parsed")
	}
	if page == nil || page.MoreItemsInCollection {
		t.Fatalf("unexpected pagination: %#v", page)
	}
	if page.NextStart == nil || *page.NextStart != 1 {
		t.Fatalf("unexpected next_start: %#v", page)
	}
}

func TestCallLogsService_Create(t *testing.T) {
	t.Parallel()

	start := time.Date(2022, 12, 12, 1, 1, 1, 0, time.UTC)
	end := time.Date(2022, 12, 12, 1, 2, 1, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/callLogs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "create" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !strings.Contains(string(body), "\"user_id\":7") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"activity_id\":8") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"subject\":\"Discovery call\"") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"duration\":\"00:01:00\"") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"from_phone_number\":\"+321\"") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"to_phone_number\":\"+123\"") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"outcome\":\"connected\"") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"start_time\":\"2022-12-12 01:01:01\"") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"person_id\":9") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"org_id\":10") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"deal_id\":11") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"lead_id\":\"adf21080-0e10-11eb-879b-05d71fb426ec\"") {
			t.Fatalf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "\"note\":\"Call notes\"") {
			t.Fatalf("unexpected body: %s", body)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"log-1","outcome":"connected","to_phone_number":"+123","start_time":"2022-12-12 01:01:01","end_time":"2022-12-12 01:02:01","has_recording":false}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	log, err := client.CallLogs.Create(
		context.Background(),
		WithCallLogUserID(UserID(7)),
		WithCallLogActivityID(ActivityID(8)),
		WithCallLogSubject("Discovery call"),
		WithCallLogDuration("00:01:00"),
		WithCallLogFromPhoneNumber("+321"),
		WithCallLogToPhoneNumber("+123"),
		WithCallLogOutcome(CallLogOutcomeConnected),
		WithCallLogStartTime(start),
		WithCallLogEndTime(end),
		WithCallLogPersonID(PersonID(9)),
		WithCallLogOrganizationID(OrganizationID(10)),
		WithCallLogDealID(DealID(11)),
		WithCallLogLeadID(LeadID("adf21080-0e10-11eb-879b-05d71fb426ec")),
		WithCallLogNote("Call notes"),
		WithCallLogsRequestOptions(pipedrive.WithHeader("X-Test", "create")),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if log.ID != "log-1" || log.Outcome != CallLogOutcomeConnected {
		t.Fatalf("unexpected log: %#v", log)
	}
	if log.StartTime == nil || log.StartTime.IsZero() {
		t.Fatalf("expected start_time to be parsed")
	}
}

func TestCallLogsService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/callLogs/log-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "get" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":"log-1","outcome":"busy","to_phone_number":"+123","start_time":"2022-12-12 01:01:01","end_time":"2022-12-12 01:02:01"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	log, err := client.CallLogs.Get(
		context.Background(),
		CallLogID("log-1"),
		WithCallLogsRequestOptions(pipedrive.WithHeader("X-Test", "get")),
	)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if log.ID != "log-1" || log.Outcome != CallLogOutcomeBusy {
		t.Fatalf("unexpected log: %#v", log)
	}
}

func TestCallLogsService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/callLogs/log-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "delete" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ok, err := client.CallLogs.Delete(
		context.Background(),
		CallLogID("log-1"),
		WithCallLogsRequestOptions(pipedrive.WithHeader("X-Test", "delete")),
	)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete success")
	}
}

func TestCallLogsService_AddRecording(t *testing.T) {
	t.Parallel()

	audio := []byte("test-audio")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/callLogs/log-1/recordings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "recording" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		contentType := r.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			t.Fatalf("parse content type: %v", err)
		}
		if mediaType != "multipart/form-data" {
			t.Fatalf("unexpected media type: %s", mediaType)
		}

		reader := multipart.NewReader(r.Body, params["boundary"])
		part, err := reader.NextPart()
		if err != nil {
			t.Fatalf("read part: %v", err)
		}
		if part.FormName() != "file" {
			t.Fatalf("unexpected form name: %s", part.FormName())
		}
		if part.FileName() != "recording.mp3" {
			t.Fatalf("unexpected filename: %s", part.FileName())
		}
		payload, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("read part: %v", err)
		}
		if !bytes.Equal(payload, audio) {
			t.Fatalf("unexpected payload: %q", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ok, err := client.CallLogs.AddRecording(
		context.Background(),
		CallLogID("log-1"),
		"recording.mp3",
		bytes.NewReader(audio),
		WithCallLogsRequestOptions(pipedrive.WithHeader("X-Test", "recording")),
	)
	if err != nil {
		t.Fatalf("AddRecording error: %v", err)
	}
	if !ok {
		t.Fatalf("expected recording success")
	}
}

func TestCallLogsService_AddRecording_SourceReadError(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{
		BaseURL: "https://example.test",
		HTTPClient: &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			defer req.Body.Close()
			_, err := io.ReadAll(req.Body)
			return nil, err
		})},
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.CallLogs.AddRecording(
		context.Background(),
		CallLogID("log-1"),
		"recording.mp3",
		&errReader{err: errors.New("boom")},
	)
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected source read error, got %v", err)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type errReader struct {
	err error
}

func (r *errReader) Read(p []byte) (int, error) {
	return 0, r.err
}
