# Developer Guide

This document outlines development standards, branching, and the release process for GW2CLI.

## Engineering Standards

### Code Quality
- **Concise comments**: Avoid "This function does X" style comments. Only explain non-obvious "why" logic.
- **Commit Style**: Use terse, imperative-style commit messages (e.g., `feat: add cache`, `fix: handle api timeout`).
- **Idiomatic Go**: Follow standard Go formatting (`go fmt`) and modular patterns.

### Architecture
- `cmd/gw2cli/`: CLI entry point and flag parsing.
- `internal/inventory/`: Core business logic and data aggregation.
- `internal/ui/`: Output formatting and terminal interactions.
- `pkg/gw2api/`: Low-level Guild Wars 2 API client.

## Development Workflow

### Branching Strategy
- **main**: Stable branch. Never commit directly to main.
- **feature/<name>**: All new features or bug fixes must be developed on a feature branch.
- **Pull Requests**: Merge feature branches into main via PRs after verification.

### Local Development
```bash
# Build binary
make build

# Run tests
make test

# Run directly
go run cmd/gw2cli/main.go [flags]
```

## Release Process

Automated builds and releases via GitHub Actions.

### 1. Development Builds (Continuous)
Every push to the `main` branch:
- Triggers build workflow.
- Updates the `latest` tag.
- Overwrites the "Latest Development Build" release.

### 2. Versioned Releases (Official)
Official releases are triggered by Git tags:

1. Ensure `main` is up-to-date.
2. Bump `Version` in `cmd/gw2cli/main.go`.
3. Tag and push:

```bash
git checkout main
git pull origin main
git tag 1.x.x
git push origin 1.x.x
```
