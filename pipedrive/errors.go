package pipedrive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type APIError struct {
	Status    int
	Code      string
	Message   string
	Body      []byte
	RequestID string
}

func (e *APIError) Error() string {
	if e == nil {
		return "pipedrive: <nil>"
	}
	if e.Code != "" && e.Message != "" {
		return fmt.Sprintf("pipedrive: http %d %s: %s", e.Status, e.Code, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("pipedrive: http %d: %s", e.Status, e.Message)
	}
	if e.Code != "" {
		return fmt.Sprintf("pipedrive: http %d %s", e.Status, e.Code)
	}
	return fmt.Sprintf("pipedrive: http %d", e.Status)
}

type RateLimitError struct {
	*APIError

	RetryAfter time.Duration
	Limit      int
	Remaining  int
	Reset      time.Time
}

func (e *RateLimitError) Unwrap() error { return e.APIError }

func (e *RateLimitError) Error() string {
	if e == nil || e.APIError == nil {
		return "pipedrive: rate limit"
	}
	if e.RetryAfter > 0 {
		return fmt.Sprintf("%s (retry after %s)", e.APIError.Error(), e.RetryAfter)
	}
	return e.APIError.Error()
}

type pipedriveErrorPayload struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Error     string `json:"error"`
	ErrorInfo string `json:"error_info"`
}

func apiErrorFromResponse(resp *http.Response, body []byte) *APIError {
	err := &APIError{
		Status:    resp.StatusCode,
		Body:      body,
		RequestID: resp.Header.Get("X-Request-Id"),
	}

	var payload pipedriveErrorPayload
	if json.Unmarshal(body, &payload) == nil {
		if payload.Code != "" {
			err.Code = payload.Code
		}
		switch {
		case payload.Message != "":
			err.Message = payload.Message
		case payload.Error != "":
			err.Message = payload.Error
		}
	}

	return err
}

func rateLimitErrorFromResponse(resp *http.Response, body []byte, now time.Time) *RateLimitError {
	apiErr := apiErrorFromResponse(resp, body)
	rl := &RateLimitError{
		APIError: apiErr,
	}

	rl.RetryAfter = parseRetryAfter(resp.Header.Get("Retry-After"), now)
	rl.Limit = parseIntHeader(resp.Header.Get("X-RateLimit-Limit"))
	rl.Remaining = parseIntHeader(resp.Header.Get("X-RateLimit-Remaining"))
	rl.Reset = parseResetHeader(resp.Header.Get("X-RateLimit-Reset"))

	return rl
}

func parseRetryAfter(value string, now time.Time) time.Duration {
	if value == "" {
		return 0
	}
	if secs, err := strconv.Atoi(value); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(value); err == nil {
		if t.After(now) {
			return t.Sub(now)
		}
		return 0
	}
	return 0
}

func parseIntHeader(value string) int {
	if value == "" {
		return 0
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}

func parseResetHeader(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	if unix, err := strconv.ParseInt(value, 10, 64); err == nil {
		if unix > 1_000_000_000_000 {
			return time.UnixMilli(unix).UTC()
		}
		return time.Unix(unix, 0).UTC()
	}

	if t, err := http.ParseTime(value); err == nil {
		return t.UTC()
	}

	return time.Time{}
}

