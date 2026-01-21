package v1

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestFilesService_List(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/files" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"name":"doc.pdf"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	query := url.Values{}
	query.Set("limit", "1")
	files, page, err := client.Files.List(
		context.Background(),
		WithFilesQuery(query),
		WithFilesRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(files) != 1 || files[0].ID != 1 || files[0].Name != "doc.pdf" {
		t.Fatalf("unexpected files: %#v", files)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestFilesService_Get(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/files/9" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":9,"name":"contract.pdf"}}`))
	})

	file, err := client.Files.Get(context.Background(), FileID(9))
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if file.ID != 9 || file.Name != "contract.pdf" {
		t.Fatalf("unexpected file: %#v", file)
	}
}

func TestFilesService_Add(t *testing.T) {
	t.Parallel()

	contentType := "multipart/form-data; boundary=test"
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/files" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != contentType {
			t.Fatalf("unexpected content type: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":12,"name":"upload.txt"}}`))
	})

	file, err := client.Files.Add(context.Background(), strings.NewReader("file"), contentType)
	if err != nil {
		t.Fatalf("Add error: %v", err)
	}
	if file.ID != 12 || file.Name != "upload.txt" {
		t.Fatalf("unexpected file: %#v", file)
	}
}

func TestFilesService_AddRemoteFile(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/files/remote" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if got := r.Form.Get("file_type"); got != "gdoc" {
			t.Fatalf("unexpected file_type: %q", got)
		}
		if got := r.Form.Get("item_type"); got != "deal" {
			t.Fatalf("unexpected item_type: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":33,"name":"remote"}}`))
	})

	form := url.Values{}
	form.Set("file_type", "gdoc")
	form.Set("item_type", "deal")
	form.Set("item_id", "1")
	file, err := client.Files.AddRemoteFile(context.Background(), form)
	if err != nil {
		t.Fatalf("AddRemoteFile error: %v", err)
	}
	if file.ID != 33 {
		t.Fatalf("unexpected file: %#v", file)
	}
}

func TestFilesService_LinkRemoteFile(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/files/remoteLink" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if got := r.Form.Get("item_type"); got != "person" {
			t.Fatalf("unexpected item_type: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":44}}`))
	})

	form := url.Values{}
	form.Set("item_type", "person")
	form.Set("item_id", "2")
	form.Set("remote_id", "abc")
	file, err := client.Files.LinkRemoteFile(context.Background(), form)
	if err != nil {
		t.Fatalf("LinkRemoteFile error: %v", err)
	}
	if file.ID != 44 {
		t.Fatalf("unexpected file: %#v", file)
	}
}

func TestFilesService_Update(t *testing.T) {
	t.Parallel()

	contentType := "multipart/form-data; boundary=test"
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/files/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != contentType {
			t.Fatalf("unexpected content type: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5,"name":"updated"}}`))
	})

	file, err := client.Files.Update(context.Background(), FileID(5), strings.NewReader("file"), contentType)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if file.ID != 5 || file.Name != "updated" {
		t.Fatalf("unexpected file: %#v", file)
	}
}

func TestFilesService_Delete(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/files/7" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	ok, err := client.Files.Delete(context.Background(), FileID(7))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok")
	}
}

func TestFilesService_Download(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/files/8/download" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte("payload"))
	})

	data, err := client.Files.Download(context.Background(), FileID(8))
	if err != nil {
		t.Fatalf("Download error: %v", err)
	}
	if string(data) != "payload" {
		t.Fatalf("unexpected data: %q", string(data))
	}
}
