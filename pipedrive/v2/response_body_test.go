package v2

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestResponseBodyCloseError_PreservesAPIError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(context.Context, *Client) error
	}{
		{
			name: "product",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.Products.Get(ctx, 42)
				return err
			},
		},
		{
			name: "field",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.ProjectFields.Get(ctx, "custom_field")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			const payload = `{"code":"invalid_request","message":"bad request"}`
			client, err := NewClient(pipedrive.Config{
				BaseURL: "https://example.test",
				HTTPClient: &http.Client{Transport: responseRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Header: http.Header{
							"Content-Type": []string{"application/json"},
							"X-Request-Id": []string{"request-123"},
						},
						Body:    &closeErrorReadCloser{Reader: strings.NewReader(payload)},
						Request: req,
					}, nil
				})},
			})
			if err != nil {
				t.Fatalf("NewClient error: %v", err)
			}

			err = tt.call(context.Background(), client)
			var apiErr *pipedrive.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T: %v", err, err)
			}
			if apiErr.Status != http.StatusBadRequest || apiErr.RequestID != "request-123" ||
				apiErr.Code != "invalid_request" || apiErr.Message != "bad request" ||
				string(apiErr.Body) != payload {
				t.Fatalf("unexpected APIError: %#v", apiErr)
			}
		})
	}
}

type responseRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f responseRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type closeErrorReadCloser struct {
	io.Reader
}

func (c *closeErrorReadCloser) Close() error {
	return errors.New("close failed")
}

func TestV2Services_ResponseHandling(t *testing.T) {
	t.Parallel()
	calls := []struct {
		name string
		list bool
		call func(*Client) error
	}{
		{name: "stage get", call: func(c *Client) error { _, err := c.Stages.Get(t.Context(), 1); return err }},
		{name: "stage create", call: func(c *Client) error {
			_, err := c.Stages.Create(t.Context(), WithStageName("Stage"), WithStagePipelineID(1))
			return err
		}},
		{name: "stage update", call: func(c *Client) error { _, err := c.Stages.Update(t.Context(), 1, WithStageName("Stage")); return err }},
		{name: "stage delete", call: func(c *Client) error { _, err := c.Stages.Delete(t.Context(), 1); return err }},
		{name: "stage list", list: true, call: func(c *Client) error { _, _, err := c.Stages.List(t.Context()); return err }},
		{name: "activity get", call: func(c *Client) error { _, err := c.Activities.Get(t.Context(), 1); return err }},
		{name: "activity create", call: func(c *Client) error {
			_, err := c.Activities.Create(t.Context(), WithActivitySubject("Call"))
			return err
		}},
		{name: "activity update", call: func(c *Client) error {
			_, err := c.Activities.Update(t.Context(), 1, WithActivitySubject("Call"))
			return err
		}},
		{name: "activity delete", call: func(c *Client) error { _, err := c.Activities.Delete(t.Context(), 1); return err }},
		{name: "activity list", list: true, call: func(c *Client) error { _, _, err := c.Activities.List(t.Context()); return err }},
		{name: "deal get", call: func(c *Client) error { _, err := c.Deals.Get(t.Context(), 1); return err }},
		{name: "product get", call: func(c *Client) error { _, err := c.Products.Get(t.Context(), 1); return err }},
		{name: "product variation", list: true, call: func(c *Client) error { _, _, err := c.Products.ListVariations(t.Context(), 1); return err }},
		{name: "person get", call: func(c *Client) error { _, err := c.Persons.Get(t.Context(), 1); return err }},
		{name: "organization get", call: func(c *Client) error { _, err := c.Organizations.Get(t.Context(), 1); return err }},
		{name: "lead search", call: func(c *Client) error { _, _, err := c.Leads.Search(t.Context(), "term"); return err }},
		{name: "user followers", list: true, call: func(c *Client) error { _, _, err := c.Users.ListFollowers(t.Context(), 1); return err }},
		{name: "item search", call: func(c *Client) error { _, _, err := c.ItemSearch.Search(t.Context(), "term"); return err }},
		{name: "pipeline get", call: func(c *Client) error { _, err := c.Pipelines.Get(t.Context(), 1); return err }},
		{name: "project get", call: func(c *Client) error { _, err := c.Projects.Get(t.Context(), 1); return err }},
		{name: "project list", list: true, call: func(c *Client) error { _, _, err := c.Projects.List(t.Context()); return err }},
		{name: "project board", call: func(c *Client) error { _, err := c.ProjectBoards.Get(t.Context(), 1); return err }},
		{name: "project phase", call: func(c *Client) error { _, err := c.ProjectPhases.Get(t.Context(), 1); return err }},
		{name: "project template", call: func(c *Client) error { _, err := c.ProjectTemplates.Get(t.Context(), 1); return err }},
		{name: "task get", call: func(c *Client) error { _, err := c.Tasks.Get(t.Context(), 1); return err }},
		{name: "deal field delete", call: func(c *Client) error { _, err := c.DealFields.Delete(t.Context(), "field"); return err }},
		{name: "organization field delete", call: func(c *Client) error { _, err := c.OrganizationFields.Delete(t.Context(), "field"); return err }},
		{name: "person field delete", call: func(c *Client) error { _, err := c.PersonFields.Delete(t.Context(), "field"); return err }},
		{name: "product field delete", call: func(c *Client) error { _, err := c.ProductFields.Delete(t.Context(), "field"); return err }},
		{name: "activity field get", call: func(c *Client) error { _, err := c.ActivityFields.Get(t.Context(), "field"); return err }},
		{name: "project field get", call: func(c *Client) error { _, err := c.ProjectFields.Get(t.Context(), "field"); return err }},
	}
	readErr := errors.New("read failed")
	closeErr := errors.New("close failed")
	for _, call := range calls {
		t.Run(call.name, func(t *testing.T) {
			for _, tt := range []struct {
				name    string
				status  int
				payload string
				readErr error
				limit   int64
			}{
				{name: "success", status: http.StatusOK},
				{name: "malformed JSON", status: http.StatusOK, payload: `{"data":`},
				{name: "API error", status: http.StatusBadRequest, payload: `{"code":"invalid_request","message":"bad request"}`},
				{name: "rate limit", status: http.StatusTooManyRequests, payload: `{"code":"rate_limit","message":"slow down"}`},
				{name: "read failure", status: http.StatusOK, readErr: readErr},
				{name: "response limit", status: http.StatusOK, limit: 8},
			} {
				t.Run(tt.name, func(t *testing.T) {
					payload := tt.payload
					if payload == "" {
						if call.list {
							payload = `{"data":[{"id":1}]}`
						} else {
							payload = `{"data":{"id":1}}`
						}
					}
					body := &trackedResponseBody{Reader: strings.NewReader(payload), readErr: tt.readErr, closeErr: closeErr}
					client, err := NewClient(pipedrive.Config{
						BaseURL:         "https://example.test",
						MaxResponseSize: tt.limit,
						RetryPolicy:     &pipedrive.RetryPolicy{MaxAttempts: 1},
						HTTPClient: &http.Client{Transport: responseRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
							return &http.Response{StatusCode: tt.status, Header: http.Header{
								"Content-Type": []string{"application/json"},
								"X-Request-Id": []string{"request-123"},
								"Retry-After":  []string{"7"},
							}, Body: body, Request: req}, nil
						})},
					})
					if err != nil {
						t.Fatal(err)
					}
					err = call.call(client)
					if body.closes != 1 {
						t.Errorf("body closed %d times, want 1", body.closes)
					}
					switch {
					case tt.readErr != nil:
						if !errors.Is(err, tt.readErr) || !errors.Is(err, closeErr) {
							t.Errorf("error = %v, want read error", err)
						}
					case tt.limit > 0:
						var limitErr *pipedrive.ResponseTooLargeError
						if !errors.As(err, &limitErr) || limitErr.Limit != tt.limit {
							t.Errorf("error = %v, want size limit error", err)
						}
					case tt.status >= 400:
						var apiErr *pipedrive.APIError
						if !errors.As(err, &apiErr) {
							t.Fatalf("error = %T %v, want APIError", err, err)
						}
						if apiErr.Status != tt.status || apiErr.RequestID != "request-123" || string(apiErr.Body) != payload {
							t.Errorf("APIError = %#v", apiErr)
						}
						if tt.status == http.StatusTooManyRequests {
							var rateErr *pipedrive.RateLimitError
							if !errors.As(err, &rateErr) {
								t.Fatalf("error = %T %v, want RateLimitError", err, err)
							}
							if rateErr.RetryAfter != 7*time.Second {
								t.Errorf("RetryAfter = %s", rateErr.RetryAfter)
							}
						}
					case tt.name == "malformed JSON":
						var syntaxErr *json.SyntaxError
						if !errors.As(err, &syntaxErr) || !strings.HasPrefix(err.Error(), "decode response: ") {
							t.Errorf("error = %v, want contextual JSON syntax error", err)
						}
					default:
						if err != nil {
							t.Fatal(err)
						}
					}
				})
			}
		})
	}
}

type trackedResponseBody struct {
	io.Reader
	readErr  error
	closeErr error
	closes   int
}

func (b *trackedResponseBody) Read(p []byte) (int, error) {
	if b.readErr != nil {
		return 0, b.readErr
	}
	return b.Reader.Read(p)
}

func (b *trackedResponseBody) Close() error {
	b.closes++
	return b.closeErr
}
