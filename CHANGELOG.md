# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog, and this project adheres to
Semantic Versioning.

## [Unreleased]

## [1.0.6] - 2026-06-01

### Changed
- Bump the Go toolchain and preferred CI version from 1.26.2 to 1.26.3
- Bump `github/codeql-action` from 4.35.2 to 4.36.0
- Bump `github.com/oapi-codegen/runtime` from 1.4.0 to 1.4.1

## [1.0.5] - 2026-04-14

### Added
- Add a Go Report Card badge to the README.

### Changed
- Bump the Go toolchain and preferred CI version from 1.26.1 to 1.26.2.
- Bump `github.com/oapi-codegen/runtime` to v1.4.0.
- Update GitHub Actions: `actions/setup-go` to v6.4.0 and `github/codeql-action` to v4.35.1.

## [1.0.4] - 2026-03-22

### Added
- Support for `PersonID`, `OrganizationID` and `IncludeFields` parameters in search options for leads
- Support for `OrganizationID` and `IncludeFields` parameters in search options for persons

## [1.0.3] - 2026-03-22

### Added
- Cap successful response bodies at 64 MiB by default, with `Config.MaxResponseSize`, `WithResponseSizeLimit`, and `WithNoResponseSizeLimit` overrides.
- Add `v1.Files.DownloadTo` for streaming large file downloads without the default response cap.
- Add internal replayable multipart body support so upload requests can be retried safely.
- Add `make security` plus `govulncheck` and `gosec` targets.

### Changed
- Update CI to test against Go 1.25.0 and 1.26.1, run lint on the preferred toolchain, and run a dedicated security job.
- Pin GitHub Actions in CI, CodeQL, and release workflows.
- Harden generator CLI output permissions and document auth header precedence and response size controls.
- Bump `github.com/oapi-codegen/runtime` to v1.3.0 and `golang.org/x/oauth2` to v0.36.0.

### Fixed
- Make the retry transport honor per-request retry policy overrides and return a clear error for nil requests.
- Avoid dropping a byte when a response body crosses the configured response size limit.
- Redact webhook HTTP auth passwords in formatted output and JSON serialization.
- Make v1 call log recording uploads and v2 product image uploads replayable for retries.
- Remove stale v1 legacy endpoints and services that should no longer be exposed alongside v2 coverage.

## [1.0.2] - 2026-02-07

### Added
- v2 products: expose `CategoryName` for string categories returned by the API.

### Fixed
- v2 products: tolerate string `category` values in responses to avoid unmarshal errors.

## [1.0.1] - 2026-01-30

### Changed
- Bump github.com/google/uuid to v1.6.0.
- Bump golang.org/x/oauth2 to v0.34.0.
- Update GitHub Actions: actions/checkout to v6, actions/setup-go to v6.

## [1.0.0] - 2026-01-22

### Added
- v2-first SDK with full API v2 coverage and typed service surfaces.
- v1-legacy surface generated from the derived spec for endpoints not in v2.
- Pagination helpers, retries, typed errors, and raw API escape hatch.
- OAuth helper, integration smoke tests, examples, and endpoint tables.
- Release workflow and documentation.
