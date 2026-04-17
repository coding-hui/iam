# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WeCoding IAM — a Go-based Identity and Access Management system with RBAC via Casbin, JWT/Basic/SecretKey authentication, and a React frontend (git submodule in `web/`). Module path: `github.com/coding-hui/iam`, requires Go 1.24.

## Common Commands

### Build
```bash
make build                                    # Build all binaries (iam-apiserver, iamctl)
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
go test -race -short -v ./internal/identity/...  # Run single package tests
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

**Two binaries**:

1. **iam-apiserver** (`cmd/iam-apiserver/`) — Main REST API server (Gin) + optional gRPC server. Default port 8080 (HTTP), 8081 (gRPC). Connects to MySQL/SQLite + Redis.
2. **iamctl** (`cmd/iamctl/` / `internal/iamctl/`) — CLI tool with commands: authz, info, jwt, login, new, set, user, validate, version.

**Default credentials**: Admin user `ADMIN/WECODING` seeded on first run.

### L1/L2/L3 Layer Architecture

The internal packages follow a layered architecture:

**L1 - Foundation Layer** (no dependencies on other layers):
- `internal/identity/` — User/identity management (Pool, PrivilegedPool, Manager, Hasher using Argon2id)
- `internal/session/` — Session management
- `internal/authn/` — Authentication (password-based authenticator)

**L2 - Core Business Layer** (depends on L1):
- `internal/role/` — Role management (Pool, PrivilegedPool, Manager)
- `internal/policy/` — Policy management (Pool, PrivilegedPool, Manager)
- `internal/authz/` — Authorization engine (Casbin-based)

**L3 - Extension Layer** (depends on L1 + L2):
- `internal/token/` — Token management (Pool, Manager)
- `internal/mfa/` — MFA/TOTP support
- `internal/audit/` — Audit event recording
- `internal/lockout/` — Account lockout protection
- `internal/webhook/` — Webhook management

**Infrastructure**:
- `internal/api/` — HTTP router and middleware (request_id, logging, recovery, cors)
- `internal/cache/` — Redis + in-memory cache adapter
- `internal/driver/` — Registry pattern with thread-safe lazy initialization
- `internal/persistence/` — GORM-based repository implementations (MySQL/SQLite)
- `internal/apiserver/` — APIServer config, gRPC client

### Pool/Manager Pattern

Each domain (identity, role, policy, etc.) follows a consistent pattern:
- **Pool**: Read-only operations (Get, List, FindCredentialsByIdentifier)
- **PrivilegedPool**: Write operations (Create, Update, Delete) — embedded Pool
- **Manager**: Business logic orchestration, uses both Pool and PrivilegedPool

### Registry (IoC)

The `Registry` interface in `internal/driver/registry.go` provides access to all services via thread-safe lazy initialization using `sync.Once` and the `initOnce[T]` generic wrapper. `RegistryDefault` is the concrete implementation. Config is loaded from YAML via Viper.

### Data Flow

HTTP request → Gin middleware chain → API handler → Domain Manager → Domain Pool → Persister (GORM) → MySQL/SQLite
                                                                                    ↓
                                                                              Redis Cache

Authorization: Policies evaluated via Casbin RBAC using the authz engine in `internal/authz/`.

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

GORM AutoMigrate runs via `persistence.Persister.MigrateUp()` at startup with models: IdentityModel, SessionModel, RoleModel, PolicyModel, TokenModel, AuditEventModel, SecretKey. Data initialization via `Registry.Init()` which calls `Init()` on services implementing `DataInit`. System resources and built-in policies are seeded on first run.

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
