# pipedrive-go

Go SDK for the Pipedrive API with a handwritten, stable public surface.

- API v2: full support
- API v1: legacy-only (endpoints not present in v2)

Generated OpenAPI clients are internal-only; use the `pipedrive/v1` and
`pipedrive/v2` packages.

## API versions

This SDK is v2-first. Use `pipedrive/v2` for new work and only reach for
`pipedrive/v1` when an endpoint is not available in v2. The v1 surface is
derived from the v1 OpenAPI spec with all v2-covered operations removed, so
endpoints migrate out of v1 automatically as v2 grows.

Endpoint tables:
- v2: `docs/endpoints-v2.md`
- v1 legacy: `docs/endpoints-v1-legacy.md`

## Install

```sh
go get github.com/juhokoskela/pipedrive-go@latest
```

Go 1.25+ is required.

## Quickstart (API token)

v2:

```go
package main

import (
	"context"
	"log"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
	v2 "github.com/juhokoskela/pipedrive-go/pipedrive/v2"
)

func main() {
	client, err := v2.NewClient(pipedrive.Config{
		Auth: pipedrive.APITokenAuth("YOUR_API_TOKEN"),
	})
	if err != nil {
		log.Fatal(err)
	}

	pipelines, _, err := client.Pipelines.List(
		context.Background(),
		v2.WithPipelinesPageSize(50),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("pipelines=%d", len(pipelines))
}
```

v1 legacy:

```go
client, err := v1.NewClient(pipedrive.Config{
	Auth: pipedrive.APITokenAuth("YOUR_API_TOKEN"),
})
if err != nil {
	log.Fatal(err)
}

currencies, err := client.Currencies.List(context.Background(), v1.ListCurrenciesRequest{})
if err != nil {
	log.Fatal(err)
}
log.Printf("currencies=%d", len(currencies))
```

## Pagination

Cursor pagination helpers are exposed as `ListPager` and `ForEach`:

```go
pager := client.Deals.ListPager(v2.WithDealsPageSize(100))
err := pager.ForEach(context.Background(), func(d v2.Deal) error {
	log.Printf("deal id=%d title=%s", d.ID, d.Title)
	return nil
})
if err != nil {
	log.Fatal(err)
}
```

## OAuth2

Use the v1 OAuth helper to build the authorize URL and exchange tokens, then
pass the access token into a v2 client.

```go
import (
	"context"
	"encoding/base64"
	"log"

	"golang.org/x/oauth2"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
	v1 "github.com/juhokoskela/pipedrive-go/pipedrive/v1"
	v2 "github.com/juhokoskela/pipedrive-go/pipedrive/v2"
)

ctx := context.Background()
oauthClient, _ := v1.NewClient(pipedrive.Config{})

authURL, err := oauthClient.OAuth.Authorize(
	ctx,
	v1.WithOAuthClientID(clientID),
	v1.WithOAuthRedirectURI(redirectURI),
	v1.WithOAuthState("state123"),
)
if err != nil {
	log.Fatal(err)
}
log.Printf("open: %s", authURL)

authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret))
tokens, err := oauthClient.OAuth.GetTokens(
	ctx,
	v1.WithOAuthAuthorization(authHeader),
	v1.WithOAuthCode(code),
	v1.WithOAuthRedirectURI(redirectURI),
)
if err != nil {
	log.Fatal(err)
}

ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tokens.AccessToken})
apiClient, _ := v2.NewClient(pipedrive.Config{
	BaseURL: tokens.APIDomain + "/api/v2",
	Auth:    pipedrive.OAuth2Auth{TokenSource: ts},
})
```

## Retries and per-request options

Retries are enabled by default for 429/502/503/504 with backoff and jitter.
Override globally with `RetryPolicy`, or per call with request options:

```go
policy := pipedrive.DefaultRetryPolicy()
policy.MaxAttempts = 2

client, _ := v2.NewClient(pipedrive.Config{
	Auth:        pipedrive.APITokenAuth("YOUR_API_TOKEN"),
	RetryPolicy: &policy,
})

deal, err := client.Deals.Get(
	context.Background(),
	v2.DealID(123),
	v2.WithDealRequestOptions(pipedrive.WithNoRetry()),
)
```

## Custom HTTP and middleware

```go
client, _ := v2.NewClient(pipedrive.Config{
	Auth:       pipedrive.APITokenAuth("YOUR_API_TOKEN"),
	HTTPClient: customClient,
	Middleware: []pipedrive.Middleware{loggingMiddleware},
	UserAgent:  "my-app/1.0",
})
```

## Raw API escape hatch

```go
import "net/http"

var out struct {
	Data []v2.Pipeline `json:"data"`
}
err := client.Raw.Do(context.Background(), http.MethodGet, "/pipelines", nil, nil, &out)
```

## Integration checks

Opt-in integration tests (skipped by default):

```sh
PIPEDRIVE_API_TOKEN=... go test ./pipedrive -run Integration -v
PIPEDRIVE_API_TOKEN=... PIPEDRIVE_INTEGRATION_WRITE=1 go test ./pipedrive -run IntegrationV2OrganizationCreateDelete -v
```

Smoke CLI:

```sh
PIPEDRIVE_API_TOKEN=... go run ./cmd/smoke
```

Optional env overrides:
- `PIPEDRIVE_BASE_URL_V1`
- `PIPEDRIVE_BASE_URL_V2`
- `PIPEDRIVE_SMOKE_TIMEOUT`

## Examples

See the runnable examples in `examples/`:
- `examples/token` for API token usage
- `examples/oauth` for OAuth2 usage

## Versioning

The module follows semantic versioning and does not encode Pipedrive API
versioning in the module path. The v1/v2 distinction is in the package path,
not the module version.
See `RELEASING.md` for release steps.

## Development

```sh
make update-specs
make derive-v1-legacy
make generate
make docs
make fmt
make lint
go test ./...
```
