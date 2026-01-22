# Contributing

Thanks for helping improve this SDK.

## Development setup

```sh
make update-specs
make derive-v1-legacy
make generate
make docs
make fmt
make lint
go test ./...
```

Integration tests are opt-in:

```sh
PIPEDRIVE_API_TOKEN=... go test ./pipedrive -run Integration -v
```

## Branches

Use feature-focused branches:

- `feat/<topic>`
- `fix/<topic>`
- `chore/<topic>`

## Commits

- Use Conventional Commits: `feat(scope): ...`, `fix(scope): ...`, `chore(scope): ...`
- Keep commits focused (max 6 files per commit)
- Avoid drive-by refactors

## Code style

- Use `gofmt` (`make fmt`)
- Prefer functional options for request configuration
- Ensure all public calls accept `context.Context`
- Keep generated packages internal; public surface lives in `pipedrive/`, `pipedrive/v1`, `pipedrive/v2`

## Tests

- Add/adjust tests for behavior changes using `httptest`
- Keep real API tests behind env vars (skipped by default)
