package pipedrive

import (
	"context"
	"io"
	"net/http"
)

const defaultMaxResponseSize int64 = 64 << 20

type responseSizeLimitOption struct {
	set       bool
	limit     int64
	unlimited bool
}

type responseSizeLimitKey struct{}

func withResponseSizeLimit(ctx context.Context, opt responseSizeLimitOption) context.Context {
	return context.WithValue(ctx, responseSizeLimitKey{}, opt)
}

func responseSizeLimitFromContext(ctx context.Context) (responseSizeLimitOption, bool) {
	opt, ok := ctx.Value(responseSizeLimitKey{}).(responseSizeLimitOption)
	return opt, ok
}

func normalizeResponseSizeLimit(limit int64) int64 {
	switch {
	case limit < 0:
		return -1
	case limit == 0:
		return defaultMaxResponseSize
	default:
		return limit
	}
}

func effectiveResponseSizeLimit(ctx context.Context, defaultLimit int64) int64 {
	limit := normalizeResponseSizeLimit(defaultLimit)

	override, ok := responseSizeLimitFromContext(ctx)
	if !ok || !override.set {
		return limit
	}
	if override.unlimited {
		return -1
	}
	return normalizeResponseSizeLimit(override.limit)
}

type responseLimitTransport struct {
	next         http.RoundTripper
	defaultLimit int64
}

func newResponseLimitTransport(next http.RoundTripper, defaultLimit int64) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &responseLimitTransport{
		next:         next,
		defaultLimit: defaultLimit,
	}
}

func (t *responseLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.next.RoundTrip(req)
	if err != nil || resp == nil || resp.Body == nil {
		return resp, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, nil
	}

	limit := effectiveResponseSizeLimit(req.Context(), t.defaultLimit)
	if limit < 0 {
		return resp, nil
	}

	resp.Body = &responseLimitReadCloser{
		body:      resp.Body,
		limit:     limit,
		remaining: limit,
	}
	return resp, nil
}

type responseLimitReadCloser struct {
	body      io.ReadCloser
	limit     int64
	remaining int64
}

func (r *responseLimitReadCloser) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return r.body.Read(p)
	}
	if r.remaining <= 0 {
		var probe [1]byte
		n, err := r.body.Read(probe[:])
		switch {
		case n == 0 && err == io.EOF:
			return 0, io.EOF
		case n > 0 || err == nil:
			return 0, &ResponseTooLargeError{Limit: r.limit}
		default:
			return 0, err
		}
	}

	readBuf := p
	if int64(len(readBuf)) > r.remaining {
		readBuf = p[:int(r.remaining)+1]
	}

	n, err := r.body.Read(readBuf)
	if int64(n) <= r.remaining {
		r.remaining -= int64(n)
		return n, err
	}

	n = int(r.remaining)
	r.remaining = 0
	if n < 0 {
		n = 0
	}
	return n, &ResponseTooLargeError{Limit: r.limit}
}

func (r *responseLimitReadCloser) Close() error {
	return r.body.Close()
}
