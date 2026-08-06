package pipedrive

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type AuthProvider interface {
	Apply(req *http.Request) error
}

type MultiAuth []AuthProvider

func (m MultiAuth) Apply(req *http.Request) error {
	for _, p := range m {
		if p == nil {
			continue
		}
		if err := p.Apply(req); err != nil {
			return err
		}
	}
	return nil
}

type APITokenAuth string

func (a APITokenAuth) Apply(req *http.Request) error {
	if a == "" || req == nil {
		return nil
	}
	if req.Header.Get("x-api-token") != "" {
		return nil
	}
	req.Header.Set("x-api-token", string(a))
	return nil
}

type OAuth2Auth struct {
	TokenSource oauth2.TokenSource
}

func (a OAuth2Auth) Apply(req *http.Request) error {
	if a.TokenSource == nil || req == nil {
		return nil
	}
	if req.Header.Get("Authorization") != "" {
		return nil
	}
	token, err := a.TokenSource.Token()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	return nil
}

// authOrigin is the single origin credentials may be attached to. A nil
// *authOrigin applies no restriction; a zero authOrigin matches nothing,
// so credentials fail closed.
type authOrigin struct {
	scheme string
	host   string
}

func authOriginFromBaseURL(baseURL string) *authOrigin {
	if baseURL == "" {
		return nil
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return &authOrigin{}
	}
	return authOriginFromURL(u)
}

func authOriginFromURL(u *url.URL) *authOrigin {
	if u == nil || u.Scheme == "" || u.Host == "" {
		return &authOrigin{}
	}
	scheme := strings.ToLower(u.Scheme)
	return &authOrigin{scheme: scheme, host: canonicalAuthHost(scheme, u.Host)}
}

func (o *authOrigin) matches(u *url.URL) bool {
	if o == nil {
		return true
	}
	if u == nil || o.scheme == "" {
		return false
	}
	scheme := strings.ToLower(u.Scheme)
	return scheme == o.scheme && strings.EqualFold(canonicalAuthHost(scheme, u.Host), o.host)
}

// canonicalAuthHost drops the scheme's default port so
// "api.example.com:443" and "api.example.com" compare equal over https.
func canonicalAuthHost(scheme, host string) string {
	switch scheme {
	case "http":
		return strings.TrimSuffix(host, ":80")
	case "https":
		return strings.TrimSuffix(host, ":443")
	}
	return host
}

func authMiddleware(auth AuthProvider, origin *authOrigin) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if auth == nil {
				return next.RoundTrip(req)
			}
			suppressed := req.Header.Get(suppressAuthHeader) != ""
			if !origin.matches(req.URL) {
				if !suppressed {
					return next.RoundTrip(req)
				}
				cloned := req.Clone(req.Context())
				cloned.Header.Del(suppressAuthHeader)
				return next.RoundTrip(cloned)
			}
			// Credentials go on a clone, never on the caller's request.
			// This is load-bearing for redirect safety: http.Client
			// copies the headers of the request it was handed onto each
			// redirect hop, so a credential written into that request
			// would be replayed to the redirect target regardless of the
			// origin check above. Mutating req in place reintroduces the
			// cross-origin credential leak.
			cloned := req.Clone(req.Context())
			cloned.Header.Del(suppressAuthHeader)
			if suppressed {
				return next.RoundTrip(cloned)
			}
			if err := auth.Apply(cloned); err != nil {
				return nil, err
			}
			return next.RoundTrip(cloned)
		})
	}
}

// credentialHeaders carry caller credentials and must not survive a redirect
// to another origin. Go strips only Authorization, WWW-Authenticate and
// Cookie, and it treats any subdomain of the same registered domain as safe;
// both are looser than the origin this client is pinned to.
var credentialHeaders = []string{"x-api-token", "Authorization"}

// suppressAuthHeader is an internal sentinel the redirect guard sets on a hop
// that leaves the credential scope. The auth middleware honors it and always
// strips it, so it never reaches the wire. It bridges the one decision the
// middleware cannot make alone: with no pinned origin the middleware matches
// every request, but the guard sees the redirect history and knows the chain
// has crossed origins.
const suppressAuthHeader = "Pipedrive-Suppress-Auth"

// redirectCredentialGuard removes credential headers when a redirect leaves
// the credential scope, then delegates to next. When suppressAuth is true it
// also marks the hop so the auth middleware does not re-attach Auth-provider
// credentials at the transport.
//
// authMiddleware cannot cover the editor case: headers attached through
// request editors (WithHeader) are already present on the request before it
// reaches the transport, so they are copied onto redirect hops by http.Client
// rather than applied per-hop by the middleware.
func redirectCredentialGuard(origin *authOrigin, suppressAuth bool, next func(*http.Request, []*http.Request) error) func(*http.Request, []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if req != nil && leavesCredentialScope(origin, req, via) {
			for _, h := range credentialHeaders {
				req.Header.Del(h)
			}
			if suppressAuth {
				req.Header.Set(suppressAuthHeader, "1")
			}
		}
		if next != nil {
			return next(req, via)
		}
		// Reproduce the net/http default policy, which this replaces.
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
}

func leavesCredentialScope(origin *authOrigin, req *http.Request, via []*http.Request) bool {
	if origin != nil {
		return !origin.matches(req.URL)
	}
	// No pinned origin: fall back to comparing against the origin of the
	// chain's initial request. Comparing the first hop (rather than the
	// previous one) keeps a redirect back to the initial origin in scope.
	if len(via) == 0 {
		return false
	}
	return !authOriginFromURL(via[0].URL).matches(req.URL)
}
