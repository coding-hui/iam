# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WeCoding IAM is an Identity and Access Management system written in Go (backend) with a React/TypeScript frontend. It provides RBAC access control (Casbin-based), multiple authentication methods (JWT, Basic, SecretKey), audit logging, and MFA support (TOTP).

## Common Commands

### Backend (Go)

```bash
make build          # Build all binaries for host platform
make lint           # Run golangci-lint
make test           # Run unit tests with race detection and coverage
make cover          # Run tests + generate HTML coverage report
make format         # Run gofmt, goimports, golines, swag fmt
make tidy           # go mod tidy
make gen            # Run code generation (error codes, Swagger docs)
make swag           # Generate Swagger docs
make serve-swagger  # Serve Swagger UI
make all            # Full pipeline: tidy gen add-copyright format lint build
make clean          # Remove build output (_output/)
```

### Frontend (web/)

```bash
cd web
pnpm install
pnpm start          # Start dev server (UMI_ENV=dev)
pnpm build          # Production build
pnpm lint           # ESLint + Prettier + TypeScript check
pnpm test           # Jest unit tests
```

### Running a Single Test

```bash
go test ./internal/authz/... -v -run TestRoleCreate
```

## Architecture

### Backend: Ory-style Registry Pattern

All services are accessed through a central `driver.Registry` (`internal/driver/registry_default.go`), which uses `sync.Once` for lazy initialization. The registry wires up:

- **L1 (Infrastructure)**: Logger (logrus), OpenTelemetry tracer, Config
- **L1 (Low-level)**: Identity, Session, Self-service flows
- **L2**: Authorization (Casbin RBAC), Role, Policy
- **L3**: Token management, Lockout (brute-force protection), Audit, Courier (webhooks)

### Entry Point Flow

`cmd/apiserver/main.go` → `apiserver.Run()` → creates `driver.Registry` → initializes services → creates Gin router → starts HTTP server with graceful shutdown.

### HTTP Layer

All REST API routes are in `internal/api/router.go` under `/api/v1`. Global middleware includes: RequestID, Logging, Recovery, CORS. Key route groups: `/identities`, `/sessions`, `/login`, `/roles`, `/policies`, `/authz/check`, `/tokens`, `/audit/events`.

### Persistence

GORM-based SQL implementations in `internal/persistence/sql/` support MySQL and SQLite (switchable via config). The cache layer (`internal/cache/`) uses Redis with in-memory fallback.

### Authorization Engine

Casbin-based RBAC in `internal/authz/`. Roles are granted to users; policies combine resources, operations, and authorization effects.

## Key Patterns

- **Registry**: All services accessed via `driver.Registry` — never instantiate services directly
- **Config**: All config structs in `internal/config/` — server, database, redis, jwt, mfa, audit, notify
- **Error Codes**: Generated in `pkg/code/` — do not edit manually
- **Graceful Shutdown**: Use `pkg/shutdown/` manager for shutdown callbacks
- **CLI Apps**: Generic CLI framework in `pkg/app/` (cobra-based with viper)

## Configuration

Server config: `configs/iam-apiserver.yaml`. Database and Redis must be configured before running.

## Development Notes

- Go 1.24, Node >= 20.0.0, pnpm 8.5.1
- Backend entry: `cmd/apiserver/main.go`
- Frontend entry: `web/` (Umi.js 4 + React 19 + Ant Design)
- Build output: `_output/platforms/`
- API docs: `api/swagger/` (Swagger 2.0)
