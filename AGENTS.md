# AGENTS.md

This file provides instructions and context for LLM-based contributors working on this repository.

## Project summary
- Build a Go SDK for Pipedrive with **API v2 full support** and **API v1 legacy-only support** (only endpoints not present in v2).
- Public API is **handwritten and stable**; **generated OpenAPI clients are internal-only**.

## Architecture rules (must follow)
- Do **not** expose generated packages as the public SDK surface.
  - Generated code lives in: `internal/gen/v1`, `internal/gen/v2`
  - Public, stable wrappers live in: `pipedrive/`, `pipedrive/v1`, `pipedrive/v2`
- Services-first public API:
  - `type Client struct { Deals *DealsService; Persons *PersonsService; ... }`
  - All calls accept `context.Context`.
  - Prefer functional options over large request structs.
- Errors must be typed and useful:
  - Preserve status, request id, and response body/payload when possible.
  - Provide `APIError` and `RateLimitError` that support `errors.Is/As`.
- Retries/rate limiting are first-class:
  - Default retry for 429/502/503/504 with backoff + jitter.
  - Respect `Retry-After`.
  - Allow per-request overrides (`WithNoRetry`, `WithRetryPolicy`).
- Transport must be customizable:
  - Support custom `http.Client`/`RoundTripper` and middleware chaining.
- Pagination helpers must feel like Go:
  - Prefer explicit pager or callback iterator patterns.
- Models:
  - Keep exported models curated (typed IDs, sane time handling).
  - Avoid exposing generated “pointer swamp” directly.
- Escape hatch:
  - Provide `Raw().Do(...)` so users aren’t blocked by missing wrappers.

## Development process (required)
- Work **test-driven**:
  - Write/adjust a failing test first for behavior changes.
  - Use `httptest` for request/response semantics.
  - Keep real API integration tests behind env vars and skipped by default.
- Use **feature-focused branches**:
  - Branch name format: `feat/<topic>`, `fix/<topic>`, `chore/<topic>`.
- Use **focused commits**:
  - Max **6 files per commit**. If more changes are needed, split into multiple commits.
  - Keep commits feature-scoped; avoid drive-by refactors.
  - Commit message format (Conventional Commits style):
    - `feat(scope): <summary>`
    - `fix(scope): <summary>`
    - `chore(scope): <summary>`
    - Examples: `fix(middleware): handle nil transport`, `chore(docs): add retry docs`, `feat(ci): add regen drift check`
- When a branch is complete, include **concise PR notes** in the final response:
  - Summary (1–6 bullets)
  - Testing performed
  - Breaking changes / migration notes (if any)

## Useful commands
- `go test ./...`
- `gofmt ./...`
- `make update-specs`
- `make derive-v1-legacy`
- `make generate`

