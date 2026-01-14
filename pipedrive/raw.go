package pipedrive

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RawClient struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewRawClient(baseURL string, httpClient *http.Client) (*RawClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("base url must include scheme and host")
	}
	u.Path = strings.TrimSuffix(u.Path, "/")

	return &RawClient{
		baseURL:    u,
		httpClient: httpClient,
	}, nil
}

func (c *RawClient) Do(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	if c == nil {
		return errors.New("nil RawClient")
	}

	u := *c.baseURL
	u.Path = strings.TrimSuffix(u.Path, "/") + "/" + strings.TrimPrefix(path, "/")
	if query != nil {
		q := u.Query()
		for k, vs := range query {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		switch b := body.(type) {
		case io.Reader:
			reqBody = b
		default:
			buf, err := json.Marshal(b)
			if err != nil {
				return fmt.Errorf("encode json body: %w", err)
			}
			reqBody = bytes.NewReader(buf)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		if _, ok := body.(io.Reader); !ok {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode == http.StatusTooManyRequests {
			return RateLimitErrorFromResponse(resp, respBody, time.Now())
		}
		return APIErrorFromResponse(resp, respBody)
	}

	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode response json: %w", err)
	}
	return nil
}
