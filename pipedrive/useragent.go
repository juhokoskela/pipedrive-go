package pipedrive

import "net/http"

func userAgentMiddleware(userAgent string) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if userAgent != "" && req.Header.Get("User-Agent") == "" {
				cloned := req.Clone(req.Context())
				cloned.Header.Set("User-Agent", userAgent)
				req = cloned
			}
			return next.RoundTrip(req)
		})
	}
}

