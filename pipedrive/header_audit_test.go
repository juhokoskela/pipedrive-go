package pipedrive

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithHeader_LastValueWinsRegardlessOfCase(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Test"); got != "last" {
			t.Errorf("X-Test = %q, want last", got)
		}
	}))
	t.Cleanup(srv.Close)
	raw, err := NewRawClient(srv.URL, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	for range 64 {
		if err := raw.Do(context.Background(), http.MethodGet, "/", nil, nil, nil,
			WithHeader("X-Test", "first"), WithHeader("x-test", "last")); err != nil {
			t.Fatal(err)
		}
	}
}
