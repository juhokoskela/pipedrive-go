# Releasing

This module follows semantic versioning. The Pipedrive API version (v1/v2)
is part of the package path, not the module version.

## Release process

1. Update `CHANGELOG.md` with a new version section.
2. Run tests and checks:

```sh
make fmt
make lint
go test ./...
```

3. Tag and push:

```sh
git tag -a vX.Y.Z -m "vX.Y.Z"
git push origin vX.Y.Z
```

Pushing the tag triggers the release workflow in `.github/workflows/release.yml`.
It verifies tests and publishes GitHub release notes from `CHANGELOG.md`.

## Release notes

Release notes are sourced from the `CHANGELOG.md` section that matches the tag.
`scripts/release-notes.sh` extracts the notes and fails if the section is missing.
