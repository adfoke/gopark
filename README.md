# GoPark

## Overview
GoPark is a sample Go service built with Gin that exposes user management APIs backed by a SQLite database. The project demonstrates layered architecture with clear separation across routing, handlers, middleware, data access, and configuration packages.

## Project Layout
- `cmd/main.go` bootstraps configuration, database connectivity, migrations, Gin router, and server.
- `config/` provides Viper-powered configuration loading with a default `config.yaml`.
- `internal/handlers`, `internal/routes`, and `internal/middleware` implement HTTP behavior and cross-cutting concerns.
- `internal/db` contains database connection helpers, migrations, and CRUD logic; SQL migrations reside in `internal/migrations`.
- `internal/models` defines domain entities and validation; `internal/docs` hosts Swagger integration stubs.

## Getting Started
Install Go 1.21 or newer, then fetch dependencies:
```sh
go mod tidy
```

Adjust `config/config.yaml` or set environment variables (e.g., `GOPARK_PORT=9090`) to override defaults. The development database lives at `./gopark.db`; the startup process automatically creates the file and runs migrations.

## Running the Service
To launch the API locally:
```sh
go run ./cmd/main.go
```

This starts the server on the configured port (`8080` by default). A health probe is available at `GET /health`. Versioned user endpoints live under `/api/v1/users`, and legacy routes remain at `/user` for backward compatibility.

## Testing
Execute all unit tests with:
```sh
go test ./...
```

Handlers rely on table-driven tests and testify assertions. Before opening a pull request, ensure tests pass and add new cases covering both successful and error paths.

## Tooling & Documentation
- `logrus` provides structured logging; use contextual fields for request tracing.
- `swag` comments in `internal/docs` allow Swagger generation (`swag init`) when the tool is installed.
- The repository ships with an `AGENTS.md` contributor guide summarizing workflows, style rules, and PR expectations.

## Contributing
Create focused branches, write imperative commit messages (e.g., “Add pagination validation”), and describe behavior changes in pull requests. Include the tests you ran and, for API updates, attach sample requests or responses. Please avoid committing secrets or production databases; rely on environment variables for sensitive configuration.
