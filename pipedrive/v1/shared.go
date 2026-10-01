package v1

import (
	"net/http"
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

func validateIDs[T ~int64](ids []T, label string) error {
	for _, id := range ids {
		if err := validateID(id, label); err != nil {
			return err
		}
	}
	return nil
}

func intIDs[T ~int64](ids []T) []int {
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		out = append(out, int(id))
	}
	return out
}

func formatV1Time(t time.Time) string {
	return t.Format(v1DateTimeLayout)
}
