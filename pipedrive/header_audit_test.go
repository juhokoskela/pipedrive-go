package pipedrive

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithHeader_LastValueWinsRegardlessOfCase(t *testing.T) {
	t.Parallel()
	opts := []RequestOption{WithHeader("X-Test", "first"), WithHeader("x-test", "last")}
	var cfg requestOptions
	for _, opt := range opts {
		opt(&cfg)
	}
	if len(cfg.headers) != 1 || cfg.headers["X-Test"] != "last" {
		t.Fatalf("headers = %v, want one X-Test header with value last", cfg.headers)
	}
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
	if err := raw.Do(context.Background(), http.MethodGet, "/", nil, nil, nil, opts...); err != nil {
		t.Fatal(err)
	}
}
