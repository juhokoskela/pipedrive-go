package pipedrive

import "net/http"

type Config struct {
	// BaseURL is the API origin requests are sent to. When set, credentials
	// from Auth are only attached to requests targeting this origin (scheme
	// and host), including redirect hops — a cross-origin redirect target
	// never receives them. When empty, Auth applies to every request made
	// through the client, but a redirect that leaves the initial request's
	// origin still suppresses credentials on that hop.
	BaseURL string

	// OAuthBaseURL is the origin the v1 OAuth service uses for authorize
	// URLs and token exchange/refresh. Empty uses the production default
	// (https://oauth.pipedrive.com). Requests to it authenticate with
	// explicit client credentials; Auth is never attached.
	OAuthBaseURL string

	HTTPClient *http.Client

	Middleware []Middleware

	RetryPolicy *RetryPolicy

	// MaxResponseSize caps successful response bodies in bytes.
	// Zero uses the default 64 MiB cap. Negative values disable the cap.
	MaxResponseSize int64

	UserAgent string
	Auth      AuthProvider
}

func NewHTTPClient(cfg Config) *http.Client {
	base := cfg.HTTPClient
	if base == nil {
		base = http.DefaultClient
	}

	clone := new(http.Client)
	*clone = *base

	transport := clone.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	origin := authOriginFromBaseURL(cfg.BaseURL)

	middleware := make([]Middleware, 0, len(cfg.Middleware)+2)
	middleware = append(middleware, cfg.Middleware...)
	if cfg.UserAgent != "" {
		middleware = append(middleware, userAgentMiddleware(cfg.UserAgent))
	}
	if cfg.Auth != nil {
		middleware = append(middleware, authMiddleware(cfg.Auth, origin))
	}

	// Catches credential headers set directly on the request (for example
	// through WithHeader), which the auth middleware never sees, and marks
	// off-scope hops so the middleware does not re-attach Auth credentials.
	clone.CheckRedirect = redirectCredentialGuard(origin, cfg.Auth != nil, base.CheckRedirect)

	transport = newResponseLimitTransport(transport, cfg.MaxResponseSize)
	transport = chainMiddleware(transport, middleware)

	policy := cfg.RetryPolicy
	if policy == nil {
		p := DefaultRetryPolicy()
		policy = &p
	}
	transport = newRetryTransport(transport, *policy, retryTransportOptions{})

	clone.Transport = transport
	return clone
}
