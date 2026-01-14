package pipedrive

import "net/http"

type Config struct {
	HTTPClient *http.Client

	Middleware []Middleware

	RetryPolicy *RetryPolicy
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

	transport = chainMiddleware(transport, cfg.Middleware)

	if cfg.RetryPolicy != nil {
		transport = newRetryTransport(transport, *cfg.RetryPolicy, retryTransportOptions{})
	}

	clone.Transport = transport
	return clone
}

