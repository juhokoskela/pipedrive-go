# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog, and this project adheres to
Semantic Versioning.

## [Unreleased]

### Added
- Nothing yet.

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
