package v1

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func TestOAuthService_Authorize(t *testing.T) {
	t.Parallel()

	client, err := NewClient(pipedrive.Config{BaseURL: "http://example.com"})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	client.OAuth.baseURL = "https://auth.example.com"

	urlStr, err := client.OAuth.Authorize(
		context.Background(),
		WithOAuthClientID("client-1"),
		WithOAuthRedirectURI("https://app.example.com/cb"),
		WithOAuthState("abc"),
	)
	if err != nil {
		t.Fatalf("Authorize error: %v", err)
	}
	parsed, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	if parsed.Path != "/oauth/authorize" {
		t.Fatalf("unexpected path: %s", parsed.Path)
	}
	q := parsed.Query()
	if q.Get("client_id") != "client-1" {
		t.Fatalf("unexpected client_id: %q", q.Get("client_id"))
	}
	if q.Get("redirect_uri") != "https://app.example.com/cb" {
		t.Fatalf("unexpected redirect_uri: %q", q.Get("redirect_uri"))
	}
	if q.Get("state") != "abc" {
		t.Fatalf("unexpected state: %q", q.Get("state"))
	}
}

func TestOAuthService_GetTokens(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/oauth/token" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Basic abc" {
			t.Fatalf("unexpected authorization: %q", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("unexpected content type: %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != "code=auth&grant_type=authorization_code&redirect_uri=https%3A%2F%2Fapp.example.com%2Fcb" {
			t.Fatalf("unexpected body: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"token","refresh_token":"refresh","expires_in":3600,"api_domain":"https://company.pipedrive.com"}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	oauthRaw, err := pipedrive.NewRawClient(srv.URL, srv.Client())
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}
	client.OAuth.raw = oauthRaw
	client.OAuth.baseURL = srv.URL

	tokens, err := client.OAuth.GetTokens(
		context.Background(),
		WithOAuthAuthorization("Basic abc"),
		WithOAuthCode("auth"),
		WithOAuthRedirectURI("https://app.example.com/cb"),
	)
	if err != nil {
		t.Fatalf("GetTokens error: %v", err)
	}
	if tokens.AccessToken != "token" || tokens.RefreshToken != "refresh" {
		t.Fatalf("unexpected tokens: %#v", tokens)
	}
}

func TestOAuthService_RefreshTokens(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/oauth/token/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Basic abc" {
			t.Fatalf("unexpected authorization: %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != "grant_type=refresh_token&refresh_token=refresh" {
			t.Fatalf("unexpected body: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"token2","refresh_token":"refresh2"}`))
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(pipedrive.Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	oauthRaw, err := pipedrive.NewRawClient(srv.URL, srv.Client())
	if err != nil {
		t.Fatalf("NewRawClient error: %v", err)
	}
	client.OAuth.raw = oauthRaw
	client.OAuth.baseURL = srv.URL

	tokens, err := client.OAuth.RefreshTokens(
		context.Background(),
		WithOAuthAuthorization("Basic abc"),
		WithOAuthRefreshToken("refresh"),
	)
	if err != nil {
		t.Fatalf("RefreshTokens error: %v", err)
	}
	if tokens.AccessToken != "token2" || tokens.RefreshToken != "refresh2" {
		t.Fatalf("unexpected tokens: %#v", tokens)
	}
}
