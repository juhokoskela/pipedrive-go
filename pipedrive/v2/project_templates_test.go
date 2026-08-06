package v2

import (
	"context"
	"net/http"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestProjectTemplatesService_AllOperations(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Test"); got != "project-templates" {
			t.Fatalf("unexpected request header: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/projectTemplates":
			if r.URL.Query().Get("limit") != "10" || r.URL.Query().Get("cursor") != "templates-start" {
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":6,"title":"Launch","projects_board_id":2}],"additional_data":{"next_cursor":null}}`))
		case "/projectTemplates/6":
			_, _ = w.Write([]byte(`{"data":{"id":6,"title":"Launch","projects_board_id":2}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})
	ctx := context.Background()
	requestOpt := WithProjectTemplateRequestOptions(pipedrive.WithHeader("X-Test", "project-templates"))
	items, _, err := client.ProjectTemplates.List(ctx,
		WithProjectTemplatesPageSize(10),
		WithProjectTemplatesCursor("templates-start"),
		requestOpt,
	)
	if err != nil || len(items) != 1 || items[0].ID != 6 {
		t.Fatalf("List = %#v, %v", items, err)
	}
	got, err := client.ProjectTemplates.Get(ctx, 6, requestOpt)
	if err != nil || got.ID != 6 {
		t.Fatalf("Get = %#v, %v", got, err)
	}
}

func TestProjectTemplatesService_PagerAndIterator(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/projectTemplates" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("cursor") {
		case "start":
			_, _ = w.Write([]byte(`{"data":[{"id":1}],"additional_data":{"next_cursor":"next"}}`))
		case "next":
			_, _ = w.Write([]byte(`{"data":[{"id":2}],"additional_data":{"next_cursor":null}}`))
		default:
			_, _ = w.Write([]byte(`{"data":[{"id":3}],"additional_data":{"next_cursor":null}}`))
		}
	})

	ctx := context.Background()
	pager := client.ProjectTemplates.ListPager(WithProjectTemplatesPageSize(2), WithProjectTemplatesCursor("start"))
	var ids []ProjectTemplateID
	for pager.Next(ctx) {
		for _, template := range pager.Items() {
			ids = append(ids, template.ID)
		}
	}
	if err := pager.Err(); err != nil || len(ids) != 2 || ids[1] != 2 {
		t.Fatalf("pager = %v, %v", ids, err)
	}

	var iterated []ProjectTemplateID
	if err := client.ProjectTemplates.ForEach(ctx, func(template ProjectTemplate) error {
		iterated = append(iterated, template.ID)
		return nil
	}); err != nil {
		t.Fatalf("ForEach error: %v", err)
	}
	if len(iterated) != 1 || iterated[0] != 3 {
		t.Fatalf("unexpected iterated templates: %v", iterated)
	}
}
