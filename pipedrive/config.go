package pipedrive

import "net/http"

type Config struct {
	HTTPClient *http.Client

	Middleware []Middleware

	RetryPolicy *RetryPolicy

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

	transport = chainMiddleware(transport, middleware)

	if cfg.RetryPolicy != nil {
		transport = newRetryTransport(transport, *cfg.RetryPolicy, retryTransportOptions{})
	}

	clone.Transport = transport
	return clone
}
