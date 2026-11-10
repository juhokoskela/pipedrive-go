package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juhokoskela/pipedrive-go/v2/pipedrive"
)

func TestCreateActivityPreservesLargeRelatedIDs(t *testing.T) {
	t.Parallel()
	const large = int64(1<<32 + 7)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		for _, key := range []string{"owner_id", "deal_id", "org_id", "project_id"} {
			var id int64
			if err := json.Unmarshal(body[key], &id); err != nil {
				t.Error(err)
			}
			if id != large {
				t.Errorf("%s = %d, want %d", key, id, large)
			}
		}
		_, _ = w.Write([]byte(`{"data":{"id":1}}`))
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Activities.Create(context.Background(), WithActivityOwnerID(UserID(large)), WithActivityDealID(DealID(large)), WithActivityOrgID(OrganizationID(large)), WithActivityProjectID(ProjectID(large)))
	if err != nil {
		t.Fatal(err)
	}
}
