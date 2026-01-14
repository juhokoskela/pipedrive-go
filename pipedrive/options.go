package pipedrive

import (
	"context"
	"net/http"
)

type RequestEditorFunc func(ctx context.Context, req *http.Request) error

type RequestOption func(*requestOptions)

type requestOptions struct {
	headers map[string]string
	editors []RequestEditorFunc

	noRetry bool
}

func WithHeader(key, value string) RequestOption {
	return func(o *requestOptions) {
		if o.headers == nil {
			o.headers = make(map[string]string)
		}
		o.headers[key] = value
	}
}

func WithRequestEditor(fn RequestEditorFunc) RequestOption {
	return func(o *requestOptions) {
		if fn == nil {
			return
		}
		o.editors = append(o.editors, fn)
	}
}

func WithNoRetry() RequestOption {
	return func(o *requestOptions) {
		o.noRetry = true
	}
}

func ApplyRequestOptions(ctx context.Context, opts ...RequestOption) (context.Context, []RequestEditorFunc) {
	if ctx == nil {
		ctx = context.Background()
	}

	var o requestOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&o)
	}

	if o.noRetry {
		ctx = withNoRetry(ctx)
	}

	editors := make([]RequestEditorFunc, 0, len(o.editors)+1)
	if len(o.headers) > 0 {
		headers := make(map[string]string, len(o.headers))
		for k, v := range o.headers {
			headers[k] = v
		}
		editors = append(editors, func(_ context.Context, req *http.Request) error {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
			return nil
		})
	}
	editors = append(editors, o.editors...)

	return ctx, editors
}

