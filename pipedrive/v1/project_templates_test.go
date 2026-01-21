package v1

import (
	"context"
	"net/http"
	"testing"
)

func TestProjectTemplatesService_List(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projectTemplates" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"Template"}]}`))
	})

	templates, err := client.ProjectTemplates.List(context.Background())
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(templates) != 1 || templates[0].ID != 1 {
		t.Fatalf("unexpected templates: %#v", templates)
	}
}

func TestProjectTemplatesService_Get(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projectTemplates/3" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":3,"name":"Template"}}`))
	})

	template, err := client.ProjectTemplates.Get(context.Background(), ProjectTemplateID(3))
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if template.ID != 3 {
		t.Fatalf("unexpected template: %#v", template)
	}
}
