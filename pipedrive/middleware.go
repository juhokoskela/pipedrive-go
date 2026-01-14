package pipedrive

import "net/http"

type Middleware func(next http.RoundTripper) http.RoundTripper

func chainMiddleware(base http.RoundTripper, middleware []Middleware) http.RoundTripper {
	rt := base
	for i := len(middleware) - 1; i >= 0; i-- {
		rt = middleware[i](rt)
	}
	return rt
}

