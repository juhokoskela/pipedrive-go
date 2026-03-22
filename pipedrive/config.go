package pipedrive

import "net/http"

type Config struct {
	BaseURL string

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

	middleware := make([]Middleware, 0, len(cfg.Middleware)+2)
	middleware = append(middleware, cfg.Middleware...)
	if cfg.UserAgent != "" {
		middleware = append(middleware, userAgentMiddleware(cfg.UserAgent))
	}
	if cfg.Auth != nil {
		middleware = append(middleware, authMiddleware(cfg.Auth))
	}

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
