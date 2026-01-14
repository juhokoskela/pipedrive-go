# Architecture

This repository implements a Go SDK for Pipedrive with a **stable handwritten public API** and **unexported generated OpenAPI clients**.

## Goals
- **Go-idiomatic** SDK surface: services, `context.Context`, functional options, typed errors.
- **API coverage policy**
  - **API v2**: supported in full.
  - **API v1**: supported **only** for endpoints not yet present in v2 (“v1-legacy”).
- **Stability across regeneration**: users should not be exposed to generator noise.
- **First-class reliability**: retries and rate limiting should be painless and configurable.
- **Transport customization**: users can plug in their own `http.Client`/`RoundTripper` and middleware.

## Non-goals
- Perfect one-to-one type parity with Pipedrive OpenAPI schemas.
- Exposing generated clients/types as the primary SDK surface.

## Repository layout

### Specs
- `openapi/upstream/`
  - `v1.yaml` — upstream snapshot of API v1 spec
  - `v2.yaml` — upstream snapshot of API v2 spec
- `openapi/derived/`
  - `v1-legacy.yaml` — derived from v1 by removing operations present in v2 (by `operationId`)

### Generated code (never exported)
- `internal/gen/v2/` — raw OpenAPI-generated client for API v2
- `internal/gen/v1/` — raw OpenAPI-generated client for **v1-legacy** spec only

### Public SDK (stable façade)
- `pipedrive/` — shared runtime and common types:
  - `Config`, auth, middleware, retry, pagination, error types, raw escape hatch
- `pipedrive/v2/` — handwritten API v2 façade (services + curated models)
- `pipedrive/v1/` — handwritten API v1-legacy façade (services + curated models)

## Public API shape

### Service-based client (not a “god client”)
Each versioned package exposes a small client with services:

```go
type Client struct {
    Deals   *DealsService
    Persons *PersonsService
    // ...
}
```

Each method is Go-idiomatic:

```go
func (s *DealsService) Get(ctx context.Context, id DealID, opts ...GetDealOption) (*Deal, error)
```

### Functional options
Options keep signatures stable and extensible:
- request options: `WithHeader(k, v)`, `WithNoRetry()`, `WithRetryPolicy(p)`
- list options: `WithPageSize(n)`, `WithCursor(c)`, `WithSort(...)`, `WithFields(...)`

Options are passed down to the internal request pipeline without leaking generated types.

## Request pipeline

All requests flow through a single shared runtime:

1. Build request (path params, query params, JSON body)
2. Apply auth + headers + user agent
3. Apply middleware chain (RoundTripper-based)
4. Execute request
5. Decode response
6. Convert errors into typed Go errors
7. Apply retry policy (when enabled)

### Transport + middleware
Users can supply a custom `http.Client` or transport and provide middleware hooks.

```go
type Middleware func(next http.RoundTripper) http.RoundTripper
```

This enables logging, tracing (OpenTelemetry), metrics, or request mutation without SDK forks.

## Authentication

The SDK supports:
- API token (`x-api-token` header)
- OAuth2 access token (via an auth provider callback)

Precedence matches the Node SDK: **API token wins** if both are configured.

Auth is modeled as an interface so users can rotate tokens:

```go
type AuthProvider interface {
    Apply(req *http.Request) error
}
```

## Errors

The SDK returns errors that are useful in production:
- Preserve HTTP status code
- Preserve request ID (when present)
- Preserve the raw body and/or structured Pipedrive error payload
- Support `errors.Is/As`

Recommended core types:
- `APIError` — general non-2xx responses
- `RateLimitError` — 429 responses with parsed rate limit headers

## Retries and rate limiting

Retries are first-class and configurable:
- default retries on: **429, 502, 503, 504**
- exponential backoff + jitter
- respect `Retry-After` when present
- allow per-request override: `WithNoRetry()` / `WithRetryPolicy(...)`

The retry implementation must be observable (hooks) and must never silently hide permanent failures.

## Pagination

Pagination helpers should hide API quirks and feel like Go:

### Explicit pager
```go
pager := client.Deals.ListPager(req)
for pager.Next(ctx) {
    deals := pager.Items()
    // ...
}
if err := pager.Err(); err != nil { ... }
```

### Callback iterator
```go
err := client.Deals.ForEach(ctx, req, func(d Deal) error {
    return nil
})
```

Avoid channels by default to keep cancellation semantics simple.

## Models

Generated models are typically not user-friendly (pointer swamp, odd names, optional wrapper types).

Strategy:
- Expose curated “domain” models for common resources: `Deal`, `Person`, `Organization`, …
- Use typed IDs: `type DealID int64`
- Prefer `time.Time` only where formats are reliable; otherwise keep raw strings
- Preserve custom fields via a `map[string]any` or specialized type
- Map between façade models and generated models internally

## Escape hatches

Even a good SDK cannot cover every use case immediately. Provide a supported escape hatch:

- `Raw().Do(ctx, method, path, query, body, &out)` per API version

This avoids forcing users to import internal generated packages.

## v1-legacy policy enforcement

The v1 façade is intentionally limited. We enforce this mechanically:
- Derive `openapi/derived/v1-legacy.yaml` by removing operations whose `operationId` exists in v2.
- Generate `internal/gen/v1` from that derived spec.

When an endpoint appears in v2, it disappears from v1-legacy on the next regeneration.

## Testing strategy

### Unit tests (default)
- Retry/backoff behavior and `Retry-After` handling
- Error decoding into `APIError`/`RateLimitError`
- Middleware ordering and transport customization
- Pagination helpers (pager/iterator) with `httptest`

### Contract/integration tests (opt-in)
- Small suite hitting real endpoints behind env vars (API token/OAuth)
- Skipped by default so CI is safe and deterministic

### Regeneration drift check
CI should fail if:
- `make generate` changes `internal/gen/*` without being committed
- wrappers stop compiling against regenerated internals

