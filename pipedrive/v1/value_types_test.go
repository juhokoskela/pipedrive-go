package v1

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestNumberBool_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "null", input: "null", want: false},
		{name: "empty", input: "", want: false},
		{name: "false", input: "false", want: false},
		{name: "zero", input: "0", want: false},
		{name: "true", input: "true", want: true},
		{name: "one", input: "1", want: true},
		{name: "numeric", input: "2", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got NumberBool
			if err := got.UnmarshalJSON([]byte(tt.input)); err != nil {
				t.Fatalf("UnmarshalJSON error: %v", err)
			}
			if bool(got) != tt.want {
				t.Fatalf("NumberBool = %v, want %v", bool(got), tt.want)
			}
		})
	}
}

func TestNumberBool_UnmarshalJSONErrors(t *testing.T) {
	t.Parallel()

	var nilBool *NumberBool
	if err := nilBool.UnmarshalJSON([]byte("true")); err == nil {
		t.Fatalf("expected nil receiver error")
	}

	var value NumberBool
	if err := value.UnmarshalJSON([]byte(`"bad"`)); err == nil || !strings.Contains(err.Error(), "invalid value") {
		t.Fatalf("expected invalid value error, got %v", err)
	}
}

func TestDateTime_JSON(t *testing.T) {
	t.Parallel()

	var parsed DateTime
	if err := parsed.UnmarshalJSON([]byte(`"2024-01-02T03:04:05Z"`)); err != nil {
		t.Fatalf("UnmarshalJSON RFC3339 error: %v", err)
	}
	if parsed.IsZero() {
		t.Fatalf("expected parsed time")
	}

	var v1Parsed DateTime
	if err := v1Parsed.UnmarshalJSON([]byte(`"2024-01-02 03:04"`)); err != nil {
		t.Fatalf("UnmarshalJSON v1 minute layout error: %v", err)
	}

	var nullParsed DateTime
	if err := nullParsed.UnmarshalJSON([]byte("null")); err != nil {
		t.Fatalf("UnmarshalJSON null error: %v", err)
	}
	if !nullParsed.IsZero() {
		t.Fatalf("expected null to leave zero value")
	}

	var emptyParsed DateTime
	if err := emptyParsed.UnmarshalJSON([]byte(`""`)); err != nil {
		t.Fatalf("UnmarshalJSON empty string error: %v", err)
	}
	if !emptyParsed.IsZero() {
		t.Fatalf("expected empty string to leave zero value")
	}

	encoded, err := DateTime{Time: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)}.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	if string(encoded) != `"2024-01-02 03:04:05"` {
		t.Fatalf("unexpected encoded time: %s", encoded)
	}

	zero, err := (DateTime{}).MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON zero error: %v", err)
	}
	if string(zero) != "null" {
		t.Fatalf("unexpected zero encoding: %s", zero)
	}
}

func TestDateTime_UnmarshalJSONErrors(t *testing.T) {
	t.Parallel()

	var nilTime *DateTime
	if err := nilTime.UnmarshalJSON([]byte(`"2024-01-02"`)); err == nil {
		t.Fatalf("expected nil receiver error")
	}

	var value DateTime
	if err := value.UnmarshalJSON([]byte("not-json")); err == nil {
		t.Fatalf("expected decode error")
	}
	if err := value.UnmarshalJSON([]byte(`"not-a-time"`)); err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestErrorFromResponse(t *testing.T) {
	t.Parallel()

	apiErr := errorFromResponse(&http.Response{StatusCode: http.StatusBadRequest, Header: http.Header{}}, []byte(`{"error":"bad"}`))
	var typedAPI *pipedrive.APIError
	if !errors.As(apiErr, &typedAPI) {
		t.Fatalf("expected APIError, got %T", apiErr)
	}

	rateErr := errorFromResponse(&http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{}}, []byte(`{"error":"slow"}`))
	var typedRate *pipedrive.RateLimitError
	if !errors.As(rateErr, &typedRate) {
		t.Fatalf("expected RateLimitError, got %T", rateErr)
	}
}

func TestRequestEditorConversionSkipsNil(t *testing.T) {
	t.Parallel()

	editors := toRequestEditors([]pipedrive.RequestEditorFunc{
		nil,
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("X-Edited", "1")
			return nil
		},
	})
	if len(editors) != 1 {
		t.Fatalf("expected one editor, got %d", len(editors))
	}

	req, err := http.NewRequest(http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}
	if err := editors[0](context.Background(), req); err != nil {
		t.Fatalf("editor error: %v", err)
	}
	if got := req.Header.Get("X-Edited"); got != "1" {
		t.Fatalf("unexpected edited header: %q", got)
	}
}

func TestJoinIDs(t *testing.T) {
	t.Parallel()

	if got := joinIDs([]FieldID{}); got != "" {
		t.Fatalf("empty joinIDs = %q", got)
	}
	if got := joinIDs([]FieldID{1, 2}); got != "1,2" {
		t.Fatalf("joinIDs = %q", got)
	}
}
