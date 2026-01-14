# Milestones

This project builds a Go SDK for the Pipedrive API with:
- Full support for **API v2**
- Support for **API v1 endpoints that are not yet available in v2** (“v1-legacy”)
- A **stable, Go-idiomatic façade** over **unexported generated clients**

## M0 — Repo scaffold (docs + tooling)
- Add `ARCHITECTURE.md`, `MILESTONES.md`, `AGENTS.md`
- Establish baseline tooling: `go.mod`, `Makefile`, `golangci-lint` (or equivalent), CI skeleton
- Decide/pin OpenAPI generator(s) and codegen inputs/outputs

**Acceptance criteria**
- `go test ./...` passes (even if only tiny placeholder tests)
- CI runs lint/test and fails on generated-diff drift (once generation exists)

## M1 — OpenAPI pipeline (v2 + v1-legacy derivation)
- Snapshot upstream specs to `openapi/upstream/v1.yaml` and `openapi/upstream/v2.yaml`
- Implement spec-diff tool to derive `openapi/derived/v1-legacy.yaml`:
  - remove all v1 operations whose `operationId` exists in v2
  - output a coverage report (counts + missing groups)
- Add `make update-specs` and `make derive-v1-legacy`

**Acceptance criteria**
- Spec derivation has tests (no-op on re-run, stable output)
- Coverage report is reproducible and checked in (or generated in CI)

## M2 — Generate internal clients (unexported)
- Generate raw clients into `internal/gen/v2` and `internal/gen/v1` from:
  - `openapi/upstream/v2.yaml`
  - `openapi/derived/v1-legacy.yaml`
- Add `make generate` that is deterministic and pinned to specific generator version(s)

**Acceptance criteria**
- `internal/gen/*` is not imported by external packages (enforced by package structure)
- CI verifies `make generate` produces no diff

## M3 — Shared runtime (public, stable)
- Implement `pipedrive` core:
  - `Config`, `AuthProvider`, `Middleware` (RoundTripper-based)
  - request/response pipeline (headers, user-agent, request-id capture)
  - typed errors: `APIError`, `RateLimitError` supporting `errors.Is/As`
  - retry policy: default retries on 429/502/503/504, backoff+jitter, `Retry-After`
  - escape hatch: `Raw().Do(...)`

**Acceptance criteria**
- Unit tests for retry, error decoding, and middleware ordering
- No exported dependency on generated types

## M4 — v2 façade (services + pagination)
- Create `pipedrive/v2` client with service structs:
  - `DealsService`, `PersonsService`, `OrganizationsService`, …
- Implement pagination helpers:
  - `ListPager` + `ForEach` patterns
- Introduce curated domain models and ID types for top resources

**Acceptance criteria**
- Public API is stable and Go-idiomatic (context everywhere, options, typed IDs)
- Wrapper tests validate request shapes and response decoding with `httptest`

## M5 — v1-legacy façade (gap coverage)
- Create `pipedrive/v1` façade wrapping `internal/gen/v1`
- Add/verify only endpoints missing in v2 (enforced by spec derivation)

**Acceptance criteria**
- No overlap drift: adding an endpoint to v2 automatically removes it from v1-legacy on regen

## M6 — Integration tests (optional but recommended)
- Add contract/integration suite behind env vars:
  - `PIPEDRIVE_API_TOKEN` and/or OAuth credentials
- Prefer “small real calls” that validate auth, retries, pagination semantics

**Acceptance criteria**
- Integration tests are skipped by default and safe to run in CI on demand

## M7 — Documentation + examples
- README: install, auth (token + OAuth2), retries, pagination, raw API
- Examples in `examples/` (token + OAuth2)
- Generated endpoint tables for v2 and v1-legacy (like the Node SDK docs)

**Acceptance criteria**
- Examples compile
- Docs clearly explain v2-first + v1-legacy policy

## M8 — Release process
- Add semantic versioning, changelog strategy, and release automation
- Ensure module versioning is independent from Pipedrive API v1/v2 naming

**Acceptance criteria**
- Reproducible release notes and tags

