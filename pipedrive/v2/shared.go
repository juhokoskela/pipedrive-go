package v2

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	genv2 "github.com/juhokoskela/pipedrive-go/internal/gen/v2"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

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

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func joinCSV[T ~string](values []T) string {
	if len(values) == 0 {
		return ""
	}
	out := make([]string, 0, len(values))
	for _, v := range values {
		if v == "" {
			continue
		}
		out = append(out, string(v))
	}
	return strings.Join(out, ",")
}

// validatePathParam rejects identifier values that URL resolution would
// collapse into a different endpoint: "" yields a trailing-slash path,
// while "." and ".." survive percent-escaping and resolve as dot segments
// when the request URL is built.
func validatePathParam(value, label string) error {
	switch value {
	case "", ".", "..":
		return fmt.Errorf("invalid %s %q", label, value)
	}
	return nil
}

// validateID rejects identifiers that cannot be sent as a path parameter.
// Generated clients take path IDs as int, so on platforms where int is 32
// bits an int64 ID above MaxInt would silently wrap and address a
// different resource. Non-positive IDs never identify a resource.
// Once this returns nil, converting id to int is lossless.
func validateID[T ~int64](id T, label string) error {
	if id <= 0 {
		return fmt.Errorf("invalid %s %d", label, int64(id))
	}
	if int64(id) > math.MaxInt {
		return fmt.Errorf("%s %d overflows int on this platform", label, int64(id))
	}
	return nil
}

// validateCSVValues rejects values containing a comma. These values are
// joined into a single comma-separated query parameter, so an embedded
// comma is indistinguishable from two separate values server-side.
func validateCSVValues(values []string, label string) error {
	for _, v := range values {
		if strings.Contains(v, ",") {
			return fmt.Errorf("invalid %s %q: must not contain a comma", label, v)
		}
	}
	return nil
}

func joinIDs[T ~int64](ids []T) string {
	if len(ids) == 0 {
		return ""
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, strconv.FormatInt(int64(id), 10))
	}
	return strings.Join(out, ",")
}
