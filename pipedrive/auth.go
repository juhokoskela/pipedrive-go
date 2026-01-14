package pipedrive

import (
	"net/http"

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

func authMiddleware(auth AuthProvider) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if auth != nil {
				cloned := req.Clone(req.Context())
				if err := auth.Apply(cloned); err != nil {
					return nil, err
				}
				req = cloned
			}
			return next.RoundTrip(req)
		})
	}
}

