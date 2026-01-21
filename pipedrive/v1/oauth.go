package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

const DefaultOAuthBaseURL = "https://oauth.pipedrive.com"

type OAuthTokens struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	APIDomain    string `json:"api_domain,omitempty"`
}

type OAuthService struct {
	raw     *pipedrive.RawClient
	baseURL string
}

type AuthorizeOption interface {
	applyAuthorize(*authorizeOptions)
}

type GetTokensOption interface {
	applyGetTokens(*getTokensOptions)
}

type RefreshTokensOption interface {
	applyRefreshTokens(*refreshTokensOptions)
}

type OAuthRequestOption interface {
	GetTokensOption
	RefreshTokensOption
}

type authorizeOptions struct {
	clientID    string
	redirectURI string
	state       *string
}

type getTokensOptions struct {
	authorization  string
	code           string
	redirectURI    string
	requestOptions []pipedrive.RequestOption
}

type refreshTokensOptions struct {
	authorization  string
	refreshToken   string
	requestOptions []pipedrive.RequestOption
}

type authorizeOptionFunc func(*authorizeOptions)

func (f authorizeOptionFunc) applyAuthorize(cfg *authorizeOptions) {
	f(cfg)
}

type oauthRequestOptions struct {
	requestOptions []pipedrive.RequestOption
}

func (o oauthRequestOptions) applyGetTokens(cfg *getTokensOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

func (o oauthRequestOptions) applyRefreshTokens(cfg *refreshTokensOptions) {
	cfg.requestOptions = append(cfg.requestOptions, o.requestOptions...)
}

type oauthAuthorizationOption struct {
	authorization string
}

func (o oauthAuthorizationOption) applyGetTokens(cfg *getTokensOptions) {
	cfg.authorization = o.authorization
}

func (o oauthAuthorizationOption) applyRefreshTokens(cfg *refreshTokensOptions) {
	cfg.authorization = o.authorization
}

type oauthCodeOption struct {
	code string
}

func (o oauthCodeOption) applyGetTokens(cfg *getTokensOptions) {
	cfg.code = o.code
}

type oauthRedirectURIOption struct {
	redirectURI string
}

func (o oauthRedirectURIOption) applyAuthorize(cfg *authorizeOptions) {
	cfg.redirectURI = o.redirectURI
}

func (o oauthRedirectURIOption) applyGetTokens(cfg *getTokensOptions) {
	cfg.redirectURI = o.redirectURI
}

type oauthRefreshTokenOption struct {
	refreshToken string
}

func (o oauthRefreshTokenOption) applyRefreshTokens(cfg *refreshTokensOptions) {
	cfg.refreshToken = o.refreshToken
}

func WithOAuthClientID(id string) AuthorizeOption {
	return authorizeOptionFunc(func(cfg *authorizeOptions) {
		cfg.clientID = id
	})
}

func WithOAuthRedirectURI(uri string) oauthRedirectURIOption {
	return oauthRedirectURIOption{redirectURI: uri}
}

func WithOAuthState(state string) AuthorizeOption {
	return authorizeOptionFunc(func(cfg *authorizeOptions) {
		cfg.state = &state
	})
}

func WithOAuthAuthorization(value string) OAuthRequestOption {
	return oauthAuthorizationOption{authorization: value}
}

func WithOAuthCode(code string) GetTokensOption {
	return oauthCodeOption{code: code}
}

func WithOAuthRefreshToken(token string) RefreshTokensOption {
	return oauthRefreshTokenOption{refreshToken: token}
}

func WithOAuthRequestOptions(opts ...pipedrive.RequestOption) OAuthRequestOption {
	return oauthRequestOptions{requestOptions: opts}
}

func newAuthorizeOptions(opts []AuthorizeOption) authorizeOptions {
	var cfg authorizeOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyAuthorize(&cfg)
	}
	return cfg
}

func newGetTokensOptions(opts []GetTokensOption) getTokensOptions {
	var cfg getTokensOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyGetTokens(&cfg)
	}
	return cfg
}

func newRefreshTokensOptions(opts []RefreshTokensOption) refreshTokensOptions {
	var cfg refreshTokensOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyRefreshTokens(&cfg)
	}
	return cfg
}

func (s *OAuthService) Authorize(ctx context.Context, opts ...AuthorizeOption) (string, error) {
	_ = ctx
	cfg := newAuthorizeOptions(opts)
	if cfg.clientID == "" {
		return "", fmt.Errorf("client id is required")
	}
	if cfg.redirectURI == "" {
		return "", fmt.Errorf("redirect uri is required")
	}

	base, err := url.Parse(s.baseURL)
	if err != nil {
		return "", fmt.Errorf("parse oauth base url: %w", err)
	}
	base.Path = strings.TrimSuffix(base.Path, "/") + "/oauth/authorize"
	query := base.Query()
	query.Set("client_id", cfg.clientID)
	query.Set("redirect_uri", cfg.redirectURI)
	if cfg.state != nil {
		query.Set("state", *cfg.state)
	}
	base.RawQuery = query.Encode()
	return base.String(), nil
}

func (s *OAuthService) GetTokens(ctx context.Context, opts ...GetTokensOption) (*OAuthTokens, error) {
	cfg := newGetTokensOptions(opts)
	if cfg.authorization == "" {
		return nil, fmt.Errorf("authorization header is required")
	}
	if cfg.code == "" {
		return nil, fmt.Errorf("authorization code is required")
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", cfg.code)
	if cfg.redirectURI != "" {
		form.Set("redirect_uri", cfg.redirectURI)
	}

	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	reqOpts = append(reqOpts, pipedrive.WithHeader("Authorization", cfg.authorization))
	reqOpts = append(reqOpts, pipedrive.WithHeader("Content-Type", "application/x-www-form-urlencoded"))

	var payload OAuthTokens
	body := strings.NewReader(form.Encode())
	if err := s.raw.Do(ctx, http.MethodPost, "/oauth/token", nil, body, &payload, reqOpts...); err != nil {
		return nil, err
	}
	return &payload, nil
}

func (s *OAuthService) RefreshTokens(ctx context.Context, opts ...RefreshTokensOption) (*OAuthTokens, error) {
	cfg := newRefreshTokensOptions(opts)
	if cfg.authorization == "" {
		return nil, fmt.Errorf("authorization header is required")
	}
	if cfg.refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", cfg.refreshToken)

	reqOpts := append([]pipedrive.RequestOption{}, cfg.requestOptions...)
	reqOpts = append(reqOpts, pipedrive.WithHeader("Authorization", cfg.authorization))
	reqOpts = append(reqOpts, pipedrive.WithHeader("Content-Type", "application/x-www-form-urlencoded"))

	var payload OAuthTokens
	body := strings.NewReader(form.Encode())
	if err := s.raw.Do(ctx, http.MethodPost, "/oauth/token/", nil, body, &payload, reqOpts...); err != nil {
		return nil, err
	}
	return &payload, nil
}
