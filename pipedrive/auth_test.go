package pipedrive

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"golang.org/x/oauth2"
)

func TestNewHTTPClient_AppliesUserAgentAndAPIToken(t *testing.T) {
	t.Parallel()

	base := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("User-Agent"); got != "pipedrive-go/test" {
			t.Fatalf("unexpected user-agent: %q", got)
		}
		if got := req.Header.Get("x-api-token"); got != "token123" {
			t.Fatalf("unexpected x-api-token: %q", got)
		}
		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	httpClient := NewHTTPClient(Config{
		HTTPClient: &http.Client{Transport: base},
		UserAgent:  "pipedrive-go/test",
		Auth:       APITokenAuth("token123"),
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	resp, err := httpClient.Transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
}

func TestNewHTTPClient_AuthNotSentCrossOrigin(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var attackerHeaders []http.Header
	attacker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attackerHeaders = append(attackerHeaders, r.Header.Clone())
		mu.Unlock()
	}))
	t.Cleanup(attacker.Close)

	var apiHeaders []http.Header
	var api *httptest.Server
	api = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		apiHeaders = append(apiHeaders, r.Header.Clone())
		mu.Unlock()
		switch r.URL.Path {
		case "/cross":
			http.Redirect(w, r, attacker.URL+"/steal", http.StatusFound)
		case "/same":
			http.Redirect(w, r, api.URL+"/landed", http.StatusFound)
		}
	}))
	t.Cleanup(api.Close)

	client := NewHTTPClient(Config{
		BaseURL: api.URL,
		Auth: MultiAuth{
			APITokenAuth("secret-token"),
			OAuth2Auth{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "secret-bearer"})},
		},
	})

	resp, err := client.Get(api.URL + "/cross")
	if err != nil {
		t.Fatalf("cross-origin request error: %v", err)
	}
	_ = resp.Body.Close()

	resp, err = client.Get(api.URL + "/same")
	if err != nil {
		t.Fatalf("same-origin request error: %v", err)
	}
	_ = resp.Body.Close()

	mu.Lock()
	defer mu.Unlock()

	if len(attackerHeaders) != 1 {
		t.Fatalf("expected 1 attacker request, got %d", len(attackerHeaders))
	}
	for _, h := range []string{"x-api-token", "Authorization"} {
		if got := attackerHeaders[0].Get(h); got != "" {
			t.Fatalf("credential header %s leaked to cross-origin redirect target: %q", h, got)
		}
	}

	if len(apiHeaders) != 3 {
		t.Fatalf("expected 3 api requests, got %d", len(apiHeaders))
	}
	for i, h := range apiHeaders {
		if got := h.Get("x-api-token"); got != "secret-token" {
			t.Fatalf("api request %d missing x-api-token, got %q", i, got)
		}
		if got := h.Get("Authorization"); got != "Bearer secret-bearer" {
			t.Fatalf("api request %d missing authorization, got %q", i, got)
		}
	}
}

func TestNewHTTPClient_EditorHeadersStrippedCrossOrigin(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var attackerHeaders []http.Header
	attacker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attackerHeaders = append(attackerHeaders, r.Header.Clone())
		mu.Unlock()
	}))
	t.Cleanup(attacker.Close)

	var apiHeaders []http.Header
	var api *httptest.Server
	api = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		apiHeaders = append(apiHeaders, r.Header.Clone())
		mu.Unlock()
		switch r.URL.Path {
		case "/cross":
			http.Redirect(w, r, attacker.URL+"/steal", http.StatusFound)
		case "/same":
			http.Redirect(w, r, api.URL+"/landed", http.StatusFound)
		}
	}))
	t.Cleanup(api.Close)

	client := NewHTTPClient(Config{BaseURL: api.URL})

	// The escape hatch: a credential header attached by a request editor
	// rather than by the Auth provider.
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, api.URL+"/cross", nil)
	req.Header.Set("x-api-token", "editor-token")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("cross-origin request error: %v", err)
	}
	_ = resp.Body.Close()

	req, _ = http.NewRequestWithContext(context.Background(), http.MethodGet, api.URL+"/same", nil)
	req.Header.Set("x-api-token", "editor-token")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("same-origin request error: %v", err)
	}
	_ = resp.Body.Close()

	mu.Lock()
	defer mu.Unlock()

	if len(attackerHeaders) != 1 {
		t.Fatalf("expected 1 attacker request, got %d", len(attackerHeaders))
	}
	if got := attackerHeaders[0].Get("x-api-token"); got != "" {
		t.Fatalf("editor-set token leaked to cross-origin redirect target: %q", got)
	}

	// Same-origin redirects must keep working.
	if len(apiHeaders) != 3 {
		t.Fatalf("expected 3 api requests, got %d", len(apiHeaders))
	}
	for i, h := range apiHeaders {
		if got := h.Get("x-api-token"); got != "editor-token" {
			t.Fatalf("api request %d lost its token, got %q", i, got)
		}
	}
}

func TestNewHTTPClient_RedirectLimitPreserved(t *testing.T) {
	t.Parallel()

	var hops int
	var mu sync.Mutex
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		hops++
		mu.Unlock()
		http.Redirect(w, r, "/loop", http.StatusFound)
	}))
	t.Cleanup(srv.Close)

	client := NewHTTPClient(Config{BaseURL: srv.URL})
	resp, err := client.Get(srv.URL + "/loop")
	if err == nil {
		_ = resp.Body.Close()
		t.Fatal("expected an error from an endless redirect loop")
	}

	mu.Lock()
	defer mu.Unlock()
	if hops > 11 {
		t.Fatalf("redirect cap not enforced, server saw %d hops", hops)
	}
}

func TestNewHTTPClient_UserCheckRedirectPreserved(t *testing.T) {
	t.Parallel()

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/start" {
			http.Redirect(w, r, srv.URL+"/next", http.StatusFound)
		}
	}))
	t.Cleanup(srv.Close)

	sentinel := errors.New("user check redirect ran")
	client := NewHTTPClient(Config{
		BaseURL: srv.URL,
		HTTPClient: &http.Client{
			CheckRedirect: func(*http.Request, []*http.Request) error { return sentinel },
		},
	})

	_, err := client.Get(srv.URL + "/start")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected the user-provided CheckRedirect to run, got %v", err)
	}
}

func TestNewHTTPClient_InvalidBaseURLFailsClosed(t *testing.T) {
	t.Parallel()

	base := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("x-api-token"); got != "" {
			t.Fatalf("credentials applied despite invalid base url: %q", got)
		}
		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})

	httpClient := NewHTTPClient(Config{
		BaseURL:    "not a url",
		HTTPClient: &http.Client{Transport: base},
		Auth:       APITokenAuth("token123"),
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	resp, err := httpClient.Transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
}

func TestAPITokenAuth_Apply(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	if err := APITokenAuth("token123").Apply(req); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if got := req.Header.Get("x-api-token"); got != "token123" {
		t.Fatalf("unexpected token header: %q", got)
	}

	if err := APITokenAuth("other").Apply(req); err != nil {
		t.Fatalf("Apply existing header error: %v", err)
	}
	if got := req.Header.Get("x-api-token"); got != "token123" {
		t.Fatalf("expected existing token to be preserved, got %q", got)
	}

	if err := APITokenAuth("").Apply(req); err != nil {
		t.Fatalf("empty auth should not error: %v", err)
	}
	if err := APITokenAuth("token").Apply(nil); err != nil {
		t.Fatalf("nil request should not error: %v", err)
	}
}

func TestOAuth2Auth_Apply(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	auth := OAuth2Auth{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "access123"})}
	if err := auth.Apply(req); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer access123" {
		t.Fatalf("unexpected authorization header: %q", got)
	}

	req.Header.Set("Authorization", "Bearer existing")
	if err := auth.Apply(req); err != nil {
		t.Fatalf("Apply existing header error: %v", err)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer existing" {
		t.Fatalf("expected existing authorization to be preserved, got %q", got)
	}

	if err := (OAuth2Auth{}).Apply(req); err != nil {
		t.Fatalf("empty auth should not error: %v", err)
	}
	if err := auth.Apply(nil); err != nil {
		t.Fatalf("nil request should not error: %v", err)
	}
}

func TestMultiAuth_ApplySkipsNilAndStopsOnError(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	wantErr := errors.New("auth failed")
	auth := MultiAuth{
		nil,
		APITokenAuth("token123"),
		errAuth{err: wantErr},
		APITokenAuth("after-error"),
	}
	err = auth.Apply(req)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected auth error, got %v", err)
	}
	if got := req.Header.Get("x-api-token"); got != "token123" {
		t.Fatalf("unexpected token header: %q", got)
	}
}

type errAuth struct {
	err error
}

func (a errAuth) Apply(*http.Request) error {
	return a.err
}

func TestCanonicalAuthHost(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		scheme, host, want string
	}{
		{"http", "example.com:80", "example.com"},
		{"https", "example.com:443", "example.com"},
		{"http", "example.com:8080", "example.com:8080"},
		{"https", "example.com:8443", "example.com:8443"},
		{"https", "127.0.0.1:4443", "127.0.0.1:4443"},
		{"ws", "example.com:443", "example.com:443"},
	} {
		if got := canonicalAuthHost(tc.scheme, tc.host); got != tc.want {
			t.Errorf("canonicalAuthHost(%q, %q) = %q, want %q", tc.scheme, tc.host, got, tc.want)
		}
	}
}

func TestAuthOriginFromBaseURL_ParseErrorFailsClosed(t *testing.T) {
	t.Parallel()

	// Unparseable, scheme-less and host-less base URLs all fail closed:
	// credentials must never be attached on a guess about the origin.
	for _, baseURL := range []string{
		"https://\x7f.invalid",
		"api.example.com",
		"https://",
	} {
		origin := authOriginFromBaseURL(baseURL)
		if origin == nil {
			t.Fatalf("expected a fail-closed origin for base URL %q", baseURL)
		}
		u, err := url.Parse("https://example.com")
		if err != nil {
			t.Fatalf("parse comparison URL: %v", err)
		}
		if origin.matches(u) {
			t.Fatalf("fail-closed origin for base URL %q must not match any request URL", baseURL)
		}
	}
}

func TestLeavesCredentialScope(t *testing.T) {
	t.Parallel()

	first, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.example.com/deals", nil)
	same, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.example.com/persons", nil)
	cross, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://evil.example.net/steal", nil)

	// A nil origin with no previous hops leaves nothing to compare against.
	if leavesCredentialScope(nil, first, nil) {
		t.Error("nil origin with empty via must not leave scope")
	}
	// A nil origin falls back to the previous hop's origin.
	if leavesCredentialScope(nil, same, []*http.Request{first}) {
		t.Error("same-origin redirect must not leave scope")
	}
	if !leavesCredentialScope(nil, cross, []*http.Request{first}) {
		t.Error("cross-origin redirect must leave scope")
	}
	// A pinned origin decides by itself, ignoring the previous hop.
	pinned := authOriginFromBaseURL("https://api.example.com")
	if !leavesCredentialScope(pinned, cross, []*http.Request{cross}) {
		t.Error("pinned origin must flag a cross-origin request")
	}
}

func TestNewHTTPClient_NilOriginRedirectSuppressesCredentials(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var attackerHeaders []http.Header
	attacker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attackerHeaders = append(attackerHeaders, r.Header.Clone())
		mu.Unlock()
	}))
	t.Cleanup(attacker.Close)

	var firstHeaders []http.Header
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		firstHeaders = append(firstHeaders, r.Header.Clone())
		mu.Unlock()
		http.Redirect(w, r, attacker.URL+"/steal", http.StatusFound)
	}))
	t.Cleanup(api.Close)

	// No BaseURL: Auth still applies to every first-party request, but a
	// redirect that leaves the initial origin must not carry credentials.
	client := NewHTTPClient(Config{
		Auth: MultiAuth{
			APITokenAuth("secret-token"),
			OAuth2Auth{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "secret-bearer"})},
		},
	})

	resp, err := client.Get(api.URL + "/start")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	_ = resp.Body.Close()

	mu.Lock()
	defer mu.Unlock()

	if len(firstHeaders) != 1 {
		t.Fatalf("expected 1 api request, got %d", len(firstHeaders))
	}
	if got := firstHeaders[0].Get("x-api-token"); got != "secret-token" {
		t.Fatalf("first request missing x-api-token, got %q", got)
	}

	if len(attackerHeaders) != 1 {
		t.Fatalf("expected 1 attacker request, got %d", len(attackerHeaders))
	}
	for _, h := range []string{"x-api-token", "Authorization"} {
		if got := attackerHeaders[0].Get(h); got != "" {
			t.Fatalf("header %s reached the cross-origin redirect target: %q", h, got)
		}
	}
}

func TestNewHTTPClient_NilOriginRedirectBackToInitialOrigin(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	apiHits := make(map[string][]http.Header)
	var api *httptest.Server
	bouncer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, api.URL+"/landed", http.StatusFound)
	}))
	t.Cleanup(bouncer.Close)
	api = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		apiHits[r.URL.Path] = append(apiHits[r.URL.Path], r.Header.Clone())
		mu.Unlock()
		if r.URL.Path == "/start" {
			http.Redirect(w, r, bouncer.URL+"/bounce", http.StatusFound)
		}
	}))
	t.Cleanup(api.Close)

	client := NewHTTPClient(Config{Auth: APITokenAuth("secret-token")})
	resp, err := client.Get(api.URL + "/start")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	_ = resp.Body.Close()

	mu.Lock()
	defer mu.Unlock()

	landed := apiHits["/landed"]
	if len(landed) != 1 {
		t.Fatalf("expected 1 landed request, got %d", len(landed))
	}
	if got := landed[0].Get("x-api-token"); got != "secret-token" {
		t.Fatalf("redirect back to the initial origin must be re-authenticated, got %q", got)
	}
}

func TestNewHTTPClient_EditorHeadersStrippedWithoutAuth(t *testing.T) {
	t.Parallel()

	var attackerHeaders []http.Header
	attacker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attackerHeaders = append(attackerHeaders, r.Header.Clone())
	}))
	t.Cleanup(attacker.Close)

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, attacker.URL+"/x", http.StatusFound)
	}))
	t.Cleanup(api.Close)

	// Editor-set credentials with no Auth provider: the guard strips them
	// even when no auth middleware is installed.
	client := NewHTTPClient(Config{BaseURL: api.URL})
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, api.URL+"/x", nil)
	req.Header.Set("x-api-token", "editor-token")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	_ = resp.Body.Close()

	if len(attackerHeaders) != 1 {
		t.Fatalf("expected 1 attacker request, got %d", len(attackerHeaders))
	}
	if got := attackerHeaders[0].Get("x-api-token"); got != "" {
		t.Fatalf("editor token reached the cross-origin redirect target: %q", got)
	}
}

func TestNewHTTPClient_UserCheckRedirectCannotRestoreCredentials(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var attackerHeaders []http.Header
	attacker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attackerHeaders = append(attackerHeaders, r.Header.Clone())
		mu.Unlock()
	}))
	t.Cleanup(attacker.Close)

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, attacker.URL+"/steal", http.StatusFound)
	}))
	t.Cleanup(api.Close)

	// A preserved callback that copies headers from the redirect history back
	// onto the hop request — the guard must re-strip after it runs.
	restore := func(req *http.Request, via []*http.Request) error {
		for _, h := range []string{"x-api-token", "Authorization"} {
			if v := via[0].Header.Get(h); v != "" {
				req.Header.Set(h, v)
			}
		}
		return nil
	}

	client := NewHTTPClient(Config{
		BaseURL:    api.URL,
		HTTPClient: &http.Client{CheckRedirect: restore},
	})

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, api.URL+"/cross", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("x-api-token", "editor-token")
	req.Header.Set("Authorization", "Bearer editor-bearer")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(attackerHeaders) != 1 {
		t.Fatalf("expected 1 attacker request, got %d", len(attackerHeaders))
	}
	for _, h := range []string{"x-api-token", "Authorization"} {
		if got := attackerHeaders[0].Get(h); got != "" {
			t.Fatalf("callback-restored %s reached the cross-origin redirect target: %q", h, got)
		}
	}
}

func TestNewHTTPClient_UserCheckRedirectCannotRewriteCredentialsOffOrigin(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var attackerHeaders http.Header
	attacker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attackerHeaders = r.Header.Clone()
		mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(attacker.Close)

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/same-origin" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Redirect(w, r, "/same-origin", http.StatusFound)
	}))
	t.Cleanup(api.Close)

	rewrite := func(req *http.Request, _ []*http.Request) error {
		u, err := url.Parse(attacker.URL + "/steal")
		if err != nil {
			return err
		}
		req.URL = u
		return nil
	}

	client := NewHTTPClient(Config{
		BaseURL:    api.URL,
		HTTPClient: &http.Client{CheckRedirect: rewrite},
	})

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, api.URL+"/start", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("x-api-token", "editor-token")
	req.Header.Set("Authorization", "Bearer editor-bearer")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if attackerHeaders == nil {
		t.Fatal("expected rewritten request to reach attacker server")
	}
	for _, h := range []string{"x-api-token", "Authorization"} {
		if got := attackerHeaders.Get(h); got != "" {
			t.Fatalf("callback-rewritten %s reached the cross-origin target: %q", h, got)
		}
	}
}
