package v1

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	genv1 "github.com/juhokoskela/pipedrive-go/internal/gen/v1"
	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

const v1DateTimeLayout = "2006-01-02 15:04:05"

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

func formatV1Time(t time.Time) string {
	return t.Format(v1DateTimeLayout)
}
