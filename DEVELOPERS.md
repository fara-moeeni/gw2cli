# Developer Guide

This document outlines the development standards, branching strategy, and release process for GW2CLI.

## Engineering Standards

### Zero-AI Traceability
All contributions must remain zero-AI traceable to maintain a professional, humanized engineering style.
- **No verbose comments**: Avoid "This function does X" style comments. Only explain non-obvious "why" logic.
- **Humanized Commits**: Use terse, imperative-style commit messages (e.g., `feat: add cache`, `fix: handle api timeout`).
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
# Build the binary
make build

# Run tests
make test

# Run directly
go run cmd/gw2cli/main.go [flags]
```

## Release Process

The project uses GitHub Actions for automated builds and releases.

### 1. Development Builds (Continuous)
Every push to the `main` branch automatically:
- Triggers the build workflow.
- Updates the `latest` tag in the repository.
- Overwrites the "Latest Development Build" release with new binaries.

### 2. Versioned Releases (Official)
Official releases are triggered by Git tags. To force a new versioned build:

1. Ensure `main` is up-to-date.
2. Bump the `Version` constant in `cmd/gw2cli/main.go`.
3. Tag the commit and push:

```bash
git checkout main
git pull origin main
git tag 1.x.x
git push origin 1.x.x
```

The CI/CD pipeline will detect the tag and create a new GitHub Release named after the tag version with all cross-compiled binaries attached.
