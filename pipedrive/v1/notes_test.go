package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestNotesService_List(t *testing.T) {
	t.Parallel()

	leadID := "adf21080-0e10-11eb-879b-05d71fb426ec"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("user_id"); got != "10" {
			t.Fatalf("unexpected user_id: %q", got)
		}
		if got := q.Get("lead_id"); got != leadID {
			t.Fatalf("unexpected lead_id: %q", got)
		}
		if got := q.Get("deal_id"); got != "2" {
			t.Fatalf("unexpected deal_id: %q", got)
		}
		if got := q.Get("person_id"); got != "3" {
			t.Fatalf("unexpected person_id: %q", got)
		}
		if got := q.Get("org_id"); got != "4" {
			t.Fatalf("unexpected org_id: %q", got)
		}
		if got := q.Get("project_id"); got != "5" {
			t.Fatalf("unexpected project_id: %q", got)
		}
		if got := q.Get("start"); got != "1" {
			t.Fatalf("unexpected start: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}
		if got := q.Get("sort"); got != "add_time DESC" {
			t.Fatalf("unexpected sort: %q", got)
		}
		if got := q.Get("start_date"); got != "2024-01-02" {
			t.Fatalf("unexpected start_date: %q", got)
		}
		if got := q.Get("end_date"); got != "2024-01-10" {
			t.Fatalf("unexpected end_date: %q", got)
		}
		if got := q.Get("pinned_to_deal_flag"); got != "1" {
			t.Fatalf("unexpected pinned_to_deal_flag: %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "1" {
			t.Fatalf("unexpected header X-Test: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":1,"active_flag":true,"add_time":"2019-12-09 13:59:21","content":"Hello","lead_id":"` + leadID + `","deal_id":2,"person_id":3,"org_id":4,"project_id":5,"pinned_to_deal_flag":true,"update_time":"2019-12-09 14:26:11","user":{"email":"user@email.com","is_you":true,"name":"User Name"}}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false,"next_start":1}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	start := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	notes, additional, err := client.Notes.List(
		context.Background(),
		WithNotesUserID(UserID(10)),
		WithNotesLeadID(LeadID(leadID)),
		WithNotesDealID(DealID(2)),
		WithNotesPersonID(PersonID(3)),
		WithNotesOrganizationID(OrganizationID(4)),
		WithNotesProjectID(ProjectID(5)),
		WithNotesStart(1),
		WithNotesLimit(2),
		WithNotesSort("add_time DESC"),
		WithNotesStartDate(start),
		WithNotesEndDate(end),
		WithNotesPinnedToDeal(true),
		WithNotesRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(notes) != 1 || notes[0].ID != 1 || notes[0].Content != "Hello" {
		t.Fatalf("unexpected notes: %#v", notes)
	}
	if notes[0].AddTime == nil || notes[0].AddTime.IsZero() {
		t.Fatalf("expected add_time to be parsed")
	}
	if !notes[0].PinnedToDeal {
		t.Fatalf("expected pinned_to_deal_flag to be true")
	}
	if additional == nil || additional.Pagination == nil || additional.Pagination.NextStart != 1 {
		t.Fatalf("unexpected additional data: %#v", additional)
	}
}

func TestNotesService_Create(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["content"] != "Hello" {
			t.Fatalf("unexpected content: %#v", payload["content"])
		}
		if payload["deal_id"] != float64(2) {
			t.Fatalf("unexpected deal_id: %#v", payload["deal_id"])
		}
		if payload["user_id"] != float64(7) {
			t.Fatalf("unexpected user_id: %#v", payload["user_id"])
		}
		if payload["add_time"] != "2024-01-01 10:00:00" {
			t.Fatalf("unexpected add_time: %#v", payload["add_time"])
		}
		if payload["pinned_to_deal_flag"] != float64(1) {
			t.Fatalf("unexpected pinned_to_deal_flag: %#v", payload["pinned_to_deal_flag"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":2,"content":"Hello","deal_id":2}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	note, err := client.Notes.Create(
		context.Background(),
		WithNoteContent("Hello"),
		WithNoteDealID(DealID(2)),
		WithNoteUserID(UserID(7)),
		WithNoteAddTime(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)),
		WithNotePinnedToDeal(true),
	)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if note.ID != 2 || note.Content != "Hello" {
		t.Fatalf("unexpected note: %#v", note)
	}
}

func TestNotesService_CreateRequiresContent(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Notes.Create(context.Background(), WithNoteDealID(DealID(2)))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestNotesService_CreateRequiresTarget(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Notes.Create(context.Background(), WithNoteContent("Hello"))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestNotesService_Get(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5,"content":"Hello","deal_id":2}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	note, err := client.Notes.Get(context.Background(), NoteID(5))
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if note.ID != 5 || note.Content != "Hello" {
		t.Fatalf("unexpected note: %#v", note)
	}
}

func TestNotesService_Update(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["content"] != "Updated" {
			t.Fatalf("unexpected content: %#v", payload["content"])
		}
		if payload["pinned_to_deal_flag"] != float64(0) {
			t.Fatalf("unexpected pinned_to_deal_flag: %#v", payload["pinned_to_deal_flag"])
		}
		if payload["deal_id"] != float64(2) {
			t.Fatalf("unexpected deal_id: %#v", payload["deal_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":5,"content":"Updated","deal_id":2}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	note, err := client.Notes.Update(
		context.Background(),
		NoteID(5),
		WithNoteContent("Updated"),
		WithNotePinnedToDeal(false),
		WithNoteDealID(DealID(2)),
	)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if note.ID != 5 || note.Content != "Updated" {
		t.Fatalf("unexpected note: %#v", note)
	}
}

func TestNotesService_UpdateRequiresFields(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Notes.Update(context.Background(), NoteID(5))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestNotesService_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes/5" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":true}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ok, err := client.Notes.Delete(context.Background(), NoteID(5))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if !ok {
		t.Fatalf("unexpected delete result: %v", ok)
	}
}

func TestNotesService_ListComments(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes/5/comments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("start"); got != "1" {
			t.Fatalf("unexpected start: %q", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("unexpected limit: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"uuid":"46c3b0e1-db35-59ca-1828-4817378dff71","content":"This is a comment","add_time":"2021-06-22T07:18:16.750Z"}],"additional_data":{"pagination":{"start":0,"limit":1,"more_items_in_collection":false,"next_start":1}}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	comments, additional, err := client.Notes.ListComments(
		context.Background(),
		NoteID(5),
		WithNoteCommentsStart(1),
		WithNoteCommentsLimit(2),
	)
	if err != nil {
		t.Fatalf("ListComments error: %v", err)
	}
	if len(comments) != 1 || comments[0].Content != "This is a comment" {
		t.Fatalf("unexpected comments: %#v", comments)
	}
	if comments[0].AddTime == nil || comments[0].AddTime.IsZero() {
		t.Fatalf("expected add_time to be parsed")
	}
	if additional == nil || additional.Pagination == nil || additional.Pagination.Limit != 1 {
		t.Fatalf("unexpected additional data: %#v", additional)
	}
}

func TestNotesService_CreateComment(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes/5/comments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["content"] != "Comment" {
			t.Fatalf("unexpected content: %#v", payload["content"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"uuid":"46c3b0e1-db35-59ca-1828-4817378dff71","content":"Comment"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	comment, err := client.Notes.CreateComment(
		context.Background(),
		NoteID(5),
		WithNoteCommentContent("Comment"),
	)
	if err != nil {
		t.Fatalf("CreateComment error: %v", err)
	}
	if comment.ID == "" || comment.Content != "Comment" {
		t.Fatalf("unexpected comment: %#v", comment)
	}
}

func TestNotesService_CreateCommentRequiresContent(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Notes.CreateComment(context.Background(), NoteID(5))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestNotesService_GetComment(t *testing.T) {
	t.Parallel()

	commentID := "46c3b0e1-db35-59ca-1828-4817378dff71"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes/5/comments/"+commentID {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"uuid":"` + commentID + `","content":"Comment"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	comment, err := client.Notes.GetComment(context.Background(), NoteID(5), CommentID(commentID))
	if err != nil {
		t.Fatalf("GetComment error: %v", err)
	}
	if comment.ID != CommentID(commentID) || comment.Content != "Comment" {
		t.Fatalf("unexpected comment: %#v", comment)
	}
}

func TestNotesService_UpdateComment(t *testing.T) {
	t.Parallel()

	commentID := "46c3b0e1-db35-59ca-1828-4817378dff71"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes/5/comments/"+commentID {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["content"] != "Updated" {
			t.Fatalf("unexpected content: %#v", payload["content"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"uuid":"` + commentID + `","content":"Updated"}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	comment, err := client.Notes.UpdateComment(
		context.Background(),
		NoteID(5),
		CommentID(commentID),
		WithNoteCommentContent("Updated"),
	)
	if err != nil {
		t.Fatalf("UpdateComment error: %v", err)
	}
	if comment.Content != "Updated" {
		t.Fatalf("unexpected comment: %#v", comment)
	}
}

func TestNotesService_UpdateCommentRequiresContent(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Notes.UpdateComment(context.Background(), NoteID(5), CommentID("46c3b0e1-db35-59ca-1828-4817378dff71"))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestNotesService_DeleteComment(t *testing.T) {
	t.Parallel()

	commentID := "46c3b0e1-db35-59ca-1828-4817378dff71"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/notes/5/comments/"+commentID {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":true}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ok, err := client.Notes.DeleteComment(context.Background(), NoteID(5), CommentID(commentID))
	if err != nil {
		t.Fatalf("DeleteComment error: %v", err)
	}
	if !ok {
		t.Fatalf("unexpected delete result: %v", ok)
	}
}
