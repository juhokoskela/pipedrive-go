package v1

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestMailboxService_ListThreads(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/mailbox/mailThreads" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("unexpected limit: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"subject":"Hello"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	query := url.Values{}
	query.Set("limit", "1")
	threads, page, err := client.Mailbox.ListThreads(context.Background(), WithMailboxQuery(query))
	if err != nil {
		t.Fatalf("ListThreads error: %v", err)
	}
	if len(threads) != 1 || threads[0].ID != 1 {
		t.Fatalf("unexpected threads: %#v", threads)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestMailboxService_GetThread(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/mailbox/mailThreads/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5,"subject":"Thread"}}`))
	})

	thread, err := client.Mailbox.GetThread(context.Background(), MailThreadID(5))
	if err != nil {
		t.Fatalf("GetThread error: %v", err)
	}
	if thread.ID != 5 {
		t.Fatalf("unexpected thread: %#v", thread)
	}
}

func TestMailboxService_DeleteThread(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/mailbox/mailThreads/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	ok, err := client.Mailbox.DeleteThread(context.Background(), MailThreadID(5))
	if err != nil {
		t.Fatalf("DeleteThread error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok")
	}
}

func TestMailboxService_UpdateThread(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/mailbox/mailThreads/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("unexpected content type: %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != "read_flag=1" {
			t.Fatalf("unexpected body: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5,"subject":"Updated"}}`))
	})

	form := url.Values{}
	form.Set("read_flag", "1")
	thread, err := client.Mailbox.UpdateThread(context.Background(), MailThreadID(5), form)
	if err != nil {
		t.Fatalf("UpdateThread error: %v", err)
	}
	if thread.ID != 5 || thread.Subject != "Updated" {
		t.Fatalf("unexpected thread: %#v", thread)
	}
}

func TestMailboxService_ListThreadMessages(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/mailbox/mailThreads/7/mailMessages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":2,"subject":"Msg"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false}}}`))
	})

	msgs, page, err := client.Mailbox.ListThreadMessages(context.Background(), MailThreadID(7))
	if err != nil {
		t.Fatalf("ListThreadMessages error: %v", err)
	}
	if len(msgs) != 1 || msgs[0].ID != 2 {
		t.Fatalf("unexpected messages: %#v", msgs)
	}
	if page == nil || page.Limit != 1 {
		t.Fatalf("unexpected page: %#v", page)
	}
}

func TestMailboxService_GetMessage(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/mailbox/mailMessages/3" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":3,"subject":"Hello"}}`))
	})

	msg, err := client.Mailbox.GetMessage(context.Background(), MailMessageID(3))
	if err != nil {
		t.Fatalf("GetMessage error: %v", err)
	}
	if msg.ID != 3 {
		t.Fatalf("unexpected message: %#v", msg)
	}
}
