# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WeCoding IAM (Identity and Access Management) is a Go-based identity and access management system providing RBAC access control based on Casbin. It uses the Gin web framework, GORM for database abstraction, and Redis for caching.

## Development Commands

```bash
# Build all binaries (apiserver)
make build

# Run tests with coverage
make test
make cover        # with coverage report

# Lint code
make lint

# Format code (gofmt, goimports, golines, swag)
make format

# Run all quality checks (tidy, gen, add-copyright, format, lint, build)
make all

# Generate swagger docs
make swag

# Generate error code files and other generated sources
make gen

# Add/update copyright headers
make add-copyright

# Clean build output
make clean

# Run a single Go test
go test ./internal/... -v -run TestFunctionName

# Run with dlv debugger
DLV=1 make build
```

### Running the API Server

```bash
# Build first
make build

# Run with config file
go run ./cmd/apiserver/main.go -c configs/iam-apiserver.yaml

# Binary location after build: _output/platforms/<os>_<arch>/apiserver
```

### Frontend (web/)

The frontend is a separate UmiJS-based project in `web/`:

```bash
cd web
pnpm install
pnpm start    # Development server on http://localhost:8000
```

## Architecture

### Registry Pattern

The application uses an Ory-inspired Registry pattern (`internal/driver/registry.go`) for dependency injection. All services are lazily initialized using `sync.Once` to ensure thread-safe singleton behavior.

```go
// All services accessed through Registry interface
type Registry interface {
    IdentityManager() identity.Manager
    SessionManager() session.Manager
    RoleManager() role.Manager
    PolicyManager() policy.Manager
    AuthzEngine() *authz.Engine
    // ... other services
}
```

### Layer Separation (L1/L2/L3)

The system is organized into three layers based on the Registry interface:

- **L1 - Core Identity**: `identity`, `session`, `selfservice` - Core authentication and user management
- **L2 - Authorization**: `authz`, `role`, `policy` - RBAC permission management using Casbin
- **L3 - Infrastructure**: `token`, `lockout`, `audit`, `courier` - Supporting services

### Directory Structure

```
cmd/apiserver/          # Main entry point for API server
internal/
  apiserver/            # Server options and Run function
  api/                  # HTTP layer (router.go, middleware/)
  config/               # Configuration structs
  driver/               # Registry pattern implementation
  identity/             # Identity management (pool, manager, handler, session, token, lockout)
  authz/                # Authorization engine
    policy/             # Policy management
    role/               # Role management
  persistence/          # Database abstraction layer
    sql/                 # GORM-based implementation
  cache/                # Redis caching
  selfservice/          # Self-service flows (login, recovery, settings)
    strategies/         # Authentication strategies (password, TOTP, OIDC)
  audit/                # Audit logging
  x/                    # Utilities (env, gormutil, options, version)
pkg/
  app/                  # CLI application framework (app.NewApp)
  shutdown/             # Graceful shutdown manager
api/                    # Generated swagger docs
web/                    # Frontend (separate UmiJS project)
hack/                   # Build scripts and Makefile includes
  makelib/              # Modular Makefile includes (golang.mk, image.mk, etc.)
```

### HTTP API Structure

Routes are registered in `internal/api/router.go` under `/api/v1`:
- `/identities` - Identity management (CRUD + credentials)
- `/sessions` - Session management
- `/roles` - Role management
- `/policies` - Policy management
- `/authz/check` - Authorization check endpoint
- `/tokens` - Token management/introspection
- `/audit/events` - Audit log查询
- `/login`, `/mfa/*` - Self-service authentication flows

### Key Dependencies

- **Web Framework**: github.com/gin-gonic/gin
- **ORM**: gorm.io/gorm with mysql/sqlite drivers
- **Authorization**: github.com/casbin/casbin
- **Cache**: github.com/redis/go-redis/v9
- **CLI**: github.com/spf13/cobra + viper
- **Tracing**: go.opentelemetry.io/otel
- **Common Library**: github.com/coding-hui/common (shared CLI utilities)

### Configuration

Configuration is loaded via Viper from YAML files. Key config sections in `internal/config/config.go`:
- `Server` - HTTP server settings (port, mode, TLS)
- `Database` - GORM connection settings
- `Redis` - Cache connection
- `Authz` - Casbin enforcer path
- `JWT` - Token settings
- `MFA` - TOTP configuration
- `Audit` - Audit logging toggle

### Testing

Tests use standard Go testing with Ginkgo/Gomega for some packages:
- Test files are co-located with `*_test.go` suffix
- Run specific tests: `go test ./internal/x/gormutil/... -v`
- Coverage output: `_output/coverage.html`
