# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WeCoding IAM — a Go-based Identity and Access Management system with RBAC via Casbin, JWT/Basic/SecretKey authentication, and a React frontend (git submodule in `web/`). Module path: `github.com/coding-hui/iam`, requires Go 1.24.

## Common Commands

### Build
```bash
make build                                    # Build all binaries (iam-apiserver, iam-authzserver, iamctl)
make build BINS="iam-apiserver"               # Build specific binary
make build.multiarch                         # Multi-platform (linux_amd64, linux_arm64)
make image                                   # Build Docker images for host arch
make push                                    # Build + push Docker images
go run ./cmd/iam-apiserver/main.go -c configs/iam-apiserver.yaml  # Run API server
```

### Test
```bash
make test                                     # Run all unit tests (race + coverage)
make cover                                    # Tests + HTML coverage report in _output/coverage.html
go test -race -short -v ./internal/apiserver/domain/service/...  # Run single package tests
```

### Lint & Format
```bash
make lint                                     # golangci-lint v2.1 (config: .golangci.yaml)
make format                                   # gofmt + goimports + golines (max 180) + swag fmt
```

### Code Generation
```bash
make gen                                      # Generate error codes -> pkg/code/code_generated.go
make gen.proto                                # Generate protobuf Go files from api/proto/*.proto
make swag                                     # Regenerate Swagger docs (api/swagger/)
make ca                                       # Generate CA certificates for all components
```

### Install & Service Management
```bash
make install INSTALL_MODE=local               # Local install (SQLite + Redis, default)
make install INSTALL_MODE=docker              # Docker install
make start/stop/status/restart/logs           # Service management after install
```

### Other
```bash
make tidy                                     # go mod tidy
make tools                                    # Install dev tools (swag, golangci-lint, mockgen, etc.)
make check-updates                            # Check outdated Go dependencies
make add-copyright                            # Add license headers (required by CI)
make verify-copyright                         # Verify license headers (runs in CI)
make clean                                    # Remove _output directory
```

### Frontend (web/ submodule)
```bash
cd web && pnpm install && pnpm start          # Dev server at localhost:8000
cd web && pnpm test                           # Run frontend tests
```

## Architecture

**Three binaries**, each with a cobra-based app framework:

1. **iam-apiserver** (`cmd/iam-apiserver/`) — Main REST API server (Gin) + optional gRPC server. Default port 8080 (HTTP), 9443 (gRPC). Connects to MySQL/SQLite + Redis.
2. **iam-authzserver** (`cmd/iam-authzserver/`) — Authorization server using Casbin. Fetches policies from apiserver via gRPC, caches in Redis.
3. **iamctl** (`internal/iamctl/`) — CLI tool with commands: jwt, user, login, set, new, info, validate, version, authz.

**Default credentials**: Admin user `ADMIN/WECODING` seeded on first run.

### DDD Layering (apiserver)

The apiserver follows Domain-Driven Design with four layers:

- **Interfaces** (`internal/apiserver/interfaces/api/`) — Gin HTTP handlers + DTO assemblers
- **Domain** (`internal/apiserver/domain/`) — Business logic
  - `model/` — Entity structs (GORM models with `iam_` table prefix, self-register via `init()`)
  - `repository/` — Repository interfaces (Factory pattern with typed accessors like `UserRepository()`)
  - `service/` — Domain services implementing business logic; `InitServiceBean()` registers all in IoC
- **Infrastructure** (`internal/apiserver/infrastructure/`) — Technical implementations
  - `datastore/mysqldb/` — MySQL connection + GORM AutoMigrate
  - `datastore/sqlitedb/` — SQLite alternative (preferred when configured)
  - `datastore/sql/` — GORM-based repository implementations
  - `cache/` — Redis + in-memory cache (factory pattern)
- **Event** (`internal/apiserver/event/`) — Event bus with async listeners for auth/user events

### IoC Container

Uses `github.com/barnettZQG/inject` for dependency injection. Beans are registered by name and auto-wired via `Populate()`. Key registration points: `service.InitServiceBean()`, `api.InitAPIBean()`.

### Data Flow

HTTP request → Gin middleware chain → API handler → Domain service → Repository (GORM/MySQL) → Response

Authorization flow: authzserver fetches policies from apiserver via gRPC → caches in Redis → evaluates via Casbin RBAC.

### Authentication

Four strategies in `pkg/middleware/auth/`: JWT, Basic, APIKey, Auto (auto-detects). Middleware chain is configurable via YAML (`server.middlewares`).

### SDK Client (`sdk/`)

Go SDK (replace directive: `github.com/coding-hui/wecoding-sdk-go => ./sdk`) with Kubernetes-style client hierarchy:
```
Clientset → IamClient → APIV1Client (user, policy, role, secretkey, etc.)
                        AuthzV1Client (authz check)
```
REST client uses `gorequest`; auth supports Basic, Bearer token, SecretID/SecretKey.

### Data Initialization

Models self-register via `init()` calling `model.RegisterModel()`. GORM AutoMigrate runs at startup. Data initialization via `service.InitData()` which calls `Init()` on services implementing `DataInit`. System resources and built-in policies are seeded on first run.

## Key Shared Packages

- `pkg/server/` — GenericAPIServer wrapping Gin engine with HTTP/HTTPS lifecycle, graceful shutdown, healthz/metrics/pprof
- `pkg/middleware/` — Gin middlewares: auth (jwt/basic/apikey/auto), cors, request_id, logger, recovery, secure, nocache, dump
- `pkg/container/` — IoC/DI container wrapping barnettZQG/inject
- `pkg/code/` — Error code definitions; define constants in `base.go`/`apiserver.go` with `// ErrXxx - HTTP: description.` comments, then run `make gen`
- `pkg/log/` — Structured logging (zap-based)
- `pkg/options/` — Configuration option structs (mysql, redis, jwt, grpc, secure, auth, etc.) with `AddFlags()`, `Validate()`, `Complete()`
- `pkg/api/proto/apiserver/v1/` — Protobuf definitions + generated gRPC code
- `pkg/app/` — App and Command framework (Cobra-based) with config loading via Viper

## Code Style

- Copyright header required on all source files (template: `hack/boilerplate.txt`), enforced by CI
- golangci-lint v2 with extensive linter set (see `.golangci.yaml`)
- Max function complexity: gocyclo 30; max line length: 180 (golines)
- goimports with local prefix: `github.com/coding-hui/iam`
- JSON tags: camelCase; YAML tags: camelCase; Mapstructure: kebab-case
- Config files: YAML in `configs/`, loaded via Viper; environment variables auto-resolved (VIPER: use `viper.GetString()` with `os.Getenv()` fallback)
- API response helpers: `pkg/api/response.go` provides `OkWithData`, `OkWithMessage`, `FailWithError`, etc.

## API Structure

- Base path: `/api/v1`
- Swagger UI: `/swagger/*any`
- Health/version/metrics/pprof endpoints available
- API groups: Authentication, User, Resource, Role, Organization, Department, ApiKey, Policy, CacheServer, IdentityProvider, Application, EmailTemplate

## CI

GitHub Actions on push to main and PRs: Go 1.24, golangci-lint v2.1, and copyright header check. Release workflow builds multi-arch Docker images (linux/amd64, linux/arm64) on tags.
