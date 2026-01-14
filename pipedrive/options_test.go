package pipedrive

import (
	"context"
	"net/http"
	"testing"
)

func TestApplyRequestOptions_WithNoRetry(t *testing.T) {
	t.Parallel()

	ctx, _ := ApplyRequestOptions(context.Background(), WithNoRetry())
	if !isNoRetry(ctx) {
		t.Fatalf("expected retry to be disabled in returned context")
	}
}

func TestApplyRequestOptions_WithHeader(t *testing.T) {
	t.Parallel()

	ctx, editors := ApplyRequestOptions(context.Background(), WithHeader("X-Test", "1"))
	if ctx == nil {
		t.Fatalf("expected non-nil context")
	}
	if len(editors) != 1 {
		t.Fatalf("expected 1 editor, got %d", len(editors))
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.test", nil)
	if err := editors[0](ctx, req); err != nil {
		t.Fatalf("editor error: %v", err)
	}
	if got := req.Header.Get("X-Test"); got != "1" {
		t.Fatalf("expected X-Test=1, got %q", got)
	}
}

func TestApplyRequestOptions_WithRequestEditor(t *testing.T) {
	t.Parallel()

	var called bool
	custom := func(_ context.Context, req *http.Request) error {
		called = true
		req.Header.Set("X-Custom", "ok")
		return nil
	}

	ctx, editors := ApplyRequestOptions(context.Background(), WithRequestEditor(custom))
	if len(editors) != 1 {
		t.Fatalf("expected 1 editor, got %d", len(editors))
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.test", nil)
	if err := editors[0](ctx, req); err != nil {
		t.Fatalf("editor error: %v", err)
	}
	if !called {
		t.Fatalf("expected custom editor to be called")
	}
	if got := req.Header.Get("X-Custom"); got != "ok" {
		t.Fatalf("expected X-Custom=ok, got %q", got)
	}
}

func TestApplyRequestOptions_WithRetryPolicy(t *testing.T) {
	t.Parallel()

	p := RetryPolicy{MaxAttempts: 2}
	ctx, _ := ApplyRequestOptions(context.Background(), WithRetryPolicy(p))

	got, ok := retryPolicyFromContext(ctx)
	if !ok {
		t.Fatalf("expected policy override to be set in returned context")
	}
	if got.MaxAttempts != 2 {
		t.Fatalf("expected MaxAttempts=2, got %d", got.MaxAttempts)
	}
}
