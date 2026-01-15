package v2

import (
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
