# Repository Guidelines

## Project Structure & Module Organization
- `cmd/main.go` starts the HTTP server and wires dependencies.
- `internal/server`, `internal/routes`, `internal/handlers`, and `internal/middleware` hold the Gin application layers.
- `internal/models` and `internal/db` define persistence types and database helpers; migrations live under `internal/migrations`.
- Shared configuration and defaults reside in `config/`, while `gopark.db` is the local SQLite development database.

## Build, Test, and Development Commands
- `go run ./cmd/main.go` launches the API locally; uses `config.yaml` unless overridden by env vars.
- `go build ./...` ensures the entire workspace compiles and catches dependency issues.
- `go test ./...` runs unit tests; include `-run` to target specific packages when iterating.

## Coding Style & Naming Conventions
- Follow Go 1.20+ defaults: `gofmt` formatting, tabs for indentation, and camelCase for identifiers.
- Group files by responsibility within `internal/` and keep handler files focused on one resource.
- Prefer structured logging via `logrus` and surface errors with context from `internal/handlers/error.go`.

## Testing Guidelines
- Use Go’s testing package with testify assertions; place `_test.go` files beside implementations.
- Name test functions `Test<Component_Action>` and cover both success and failure paths.
- Run `go test ./...` before opening a pull request; add table-driven cases when scenarios vary.

## Commit & Pull Request Guidelines
- Craft commits in imperative mood (e.g., “Add parking slot handler”); keep them scoped and reversible.
- Reference related issues in commit messages or PR descriptions using `Fixes #ID` when applicable.
- PRs should describe the intent, summarize changes, list testing performed, and include screenshots or JSON samples for API changes.

## Configuration & Security Notes
- Secrets should be loaded through environment variables consumed by `config/viper`; never commit credentials or production DB files.
- Review `config/config.yaml` defaults before deploying and ensure TLS and domain-specific settings are overridden per environment.
