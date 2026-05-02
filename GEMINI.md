# Project Instructions

## Branching & Workflow
- **No Direct Pushes to Main:** Never push code directly to the `main` branch.
- **Feature Branches:** Use `feature/<feature-name>` for all new features.
- **Fix Branches:** Use `fix/<bug-name>` or `hotfix/<name>` for bug fixes and hotfixes.
- **Releases:** Trigger releases by pushing a new version tag (e.g., `git tag vX.Y.Z` followed by `git push origin vX.Y.Z`).

## Development Standards
- **Testing:** Every new feature or fix must include automated tests.
- **Local Verification:** Always build and test changes locally (`go build`, `go test ./...`) before pushing.
- **Code Style:** Follow idiomatic Go patterns and use `go-pretty` for table rendering.

## Activity Log
- [x] **Phase 1: Account Summary & Fractal Level:** Added `account` command, fixed playtime labels, and bumped version to 2.8.0.
- [x] **Phase 2: Unauthenticated Support:** Audited subcommands, implemented graceful no-auth error handling, updated README, and bumped version to 2.8.1.

## Next Steps
- Implement `strikes` and `convergences` tracking in `achievements` subcommand.
- Add support for character equipment stats/attributes.
