package v2

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

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
