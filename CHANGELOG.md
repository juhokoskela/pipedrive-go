# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog, and this project adheres to
Semantic Versioning.

## [Unreleased]

## [1.11.0] - 2026-08-18

### Added

- Add v2 organization create and update options for website, LinkedIn URL,
  industry, annual revenue and employee count, including explicit clear options
  for nullable values.
- Add v2 person create and update options for postal address, notes, instant
  messaging accounts, birthday and job title. `PersonAddress` is an alias of
  `OrganizationAddress`, preserving compatibility with existing address values.
- Add v1 support for the `nova` and `partnership` permission-set app values and
  the `nova` user-access app value.

### Changed

- Refresh the official v1 and v2 OpenAPI specifications, re-derive the v1
  legacy surface and regenerate the internal clients.

## [1.10.1] - 2026-08-14

### Changed

- Bump the preferred Go toolchain, CI, security scanning and OpenAPI drift
  checks from Go 1.26.5 to Go 1.26.6. Go 1.25 remains supported.
- Bump `github/codeql-action` from v4.37.5 to v4.37.6, keeping the `init`
  and `analyze` steps on the same pinned revision.

## [1.10.0] - 2026-08-07

### Security

- Credentials from `Config.Auth` are now attached only to requests whose origin
  matches `Config.BaseURL`. Previously the transport middleware re-applied them
  on every redirect hop, so a cross-origin redirect delivered `x-api-token` and
  OAuth bearer tokens to the redirect target. An unparseable `BaseURL` now fails
  closed and no credentials are sent at all. With an empty `Config.BaseURL`
  credentials still apply to every first-party request, but a redirect leaving
  the initial request's origin suppresses them on that hop.
- Credential headers set directly on a request (for example through
  `WithHeader("x-api-token", …)`) are stripped when a redirect leaves the pinned
  origin. `NewHTTPClient` installs a `CheckRedirect` for this; a `CheckRedirect`
  already present on a supplied `HTTPClient` is still called, and the standard
  10-redirect limit is preserved. Credential headers the preserved callback
  re-adds to an off-origin hop are removed again after it returns.
- The v1 OAuth client is built without the `Auth` provider, so an API token is
  no longer sent to the OAuth host.
- Non-2xx response bodies are truncated at 1 MiB. They previously bypassed the
  response size cap entirely and were buffered in full into `APIError.Body`.
- A server-supplied `Retry-After` is capped and no longer overflows when
  converted to a duration.

### Added

- The public v2 SDK now covers all 158 operations in the current Pipedrive
  specification, including Projects, Tasks, Project Boards, Project Phases,
  Project Templates and Project Fields.
- Documented v2 response properties and request options are now represented for
  Deals, Persons, Organizations, Products, Activities and Fields.
- `Config.OAuthBaseURL` directs v1 token exchange and refresh at a non-production
  host. Empty keeps the previous default of `https://oauth.pipedrive.com`.
- `RetryPolicy.MaxRetryAfter` caps how long a server-provided `Retry-After` may
  delay a retry. Zero uses the new 1 minute default; negative disables the cap
  and restores the previous unbounded behaviour.

### Fixed

- Field-option IDs returned by Pipedrive now decode from either integer or
  string representations without exposing generated-client types.
- Deal-ID cursor pagers reject missing identifiers locally instead of issuing
  invalid API requests.
- v1 `RefreshTokens` posted to `/oauth/token/` with a trailing slash. A redirect
  from the real endpoint downgrades the POST to GET and drops the form body,
  silently breaking token refresh.

### Changed

These are behavioural changes that existing code may depend on:

- **Collection options can explicitly clear values.** Zero-argument collection
  options now serialize an empty array; explicit JSON null updates are also
  preserved instead of being omitted.
- Person, Organization and Product field-description options remain callable
  for source compatibility but are deprecated no-ops because API v2 does not
  accept those descriptions.
- **Empty custom-field maps no longer clear fields.** `WithOrganizationCustomFieldsMap`
  and `WithProjectCustomFields` sent `"custom_fields": {}` when given an empty
  map, which clears every custom field on the record; `WithDealCustomFieldsMap`
  ignored it. All three now ignore an empty map, matching the deal behaviour. If
  you relied on an empty map to clear fields, pass the field keys explicitly
  with the values you want cleared.
- **Query options are last-wins.** Merging a raw `WithXQuery(url.Values)` with a
  typed option for the same key previously emitted the key twice and left
  precedence to the server. The later option now replaces the earlier one.
  Multiple values passed in a single call are still preserved.
- **Identifiers are validated client-side.** Every identifier passed as a
  method argument — primary IDs, secondary IDs such as `mergeWithID`, follower
  and assignment user IDs, and each entry of bulk ID slices — rejects zero,
  negative values and values that would not survive conversion to `int` on
  32-bit platforms, returning an error instead of being sent. Path parameters
  that are opaque strings reject `""`, `"."` and `".."`, which URL resolution
  would otherwise collapse into a different endpoint. Identifiers supplied
  through filter options (for example `WithDealsOwnerID`) are passed through
  unvalidated, as the server treats them as search criteria rather than
  resource addresses.
- **Custom field keys containing a comma are rejected.** Keys are joined into one
  comma-separated parameter, so an embedded comma was indistinguishable from two
  keys. The error surfaces from `List`, `Get` and the pagers.
- v1 `Persons.AddPicture` now requires a non-empty content type, matching
  `FilesService`.
- The official v1 and v2 specifications and internal generated clients were
  refreshed, the v1 legacy-only surface was re-derived, and
  `github.com/oapi-codegen/runtime` was updated to v1.6.0.

## [1.0.9] - 2026-07-16

### Changed

- Bump GitHub Actions: `actions/setup-go` from v6.5.0 to v7.0.0 and `softprops/action-gh-release` from v3.0.1 to v3.0.2
- Bump `github.com/oapi-codegen/runtime` from v1.4.2 to v1.5.0

## [1.0.8] - 2026-07-13

### Changed
- Bump the Go toolchain and preferred CI version from 1.26.4 to 1.26.5
- Bump GitHub Actions: `actions/checkout`to v7, `actions/setup-go` to v6.5.0, `softprops/action-gh-release` to v3.0.1, and `actions/codeql/*`to v4.37.0
- Bump `github.com/oapi-codegen/runtime` to v1.4.2

## [1.0.7] - 2026-06-05

### Changed

- Bump the Go toolchain and preferred CI version from 1.26.3 to 1.26.4
- Bump CodeQL & Checkout action versions

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
