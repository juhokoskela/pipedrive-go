package pipedrive

import "net/http"

type Middleware func(next http.RoundTripper) http.RoundTripper

func chainMiddleware(base http.RoundTripper, middleware []Middleware) http.RoundTripper {
	rt := base
	for i := len(middleware) - 1; i >= 0; i-- {
		rt = &middlewareTransport{RoundTripper: middleware[i](rt), next: rt}
	}
	return rt
}

type middlewareTransport struct {
	http.RoundTripper
	next http.RoundTripper
}

func (t *middlewareTransport) CloseIdleConnections() {
	if closer, ok := t.RoundTripper.(interface{ CloseIdleConnections() }); ok {
		closer.CloseIdleConnections()
		return
	}
	closeIdleConnections(t.next)
}

func closeIdleConnections(transport http.RoundTripper) {
	if closer, ok := transport.(interface{ CloseIdleConnections() }); ok {
		closer.CloseIdleConnections()
	}
}
