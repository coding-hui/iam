// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package driver

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/coding-hui/iam/internal/config"
	"github.com/coding-hui/iam/internal/selfservice"

	"go.opentelemetry.io/otel/trace"

	"github.com/coding-hui/iam/internal/audit"
	"github.com/coding-hui/iam/internal/authz"
	"github.com/coding-hui/iam/internal/authz/policy"
	"github.com/coding-hui/iam/internal/authz/role"
	"github.com/coding-hui/iam/internal/identity"
	"github.com/coding-hui/iam/internal/identity/lockout"
	"github.com/coding-hui/iam/internal/identity/session"
	"github.com/coding-hui/iam/internal/identity/token"
	"github.com/coding-hui/iam/internal/selfservice/courier"
	"github.com/coding-hui/iam/internal/selfservice/strategies"
)

// Registry is the central interface that provides access to all
// application services and infrastructure components.
// It follows the Ory Registry pattern with sync.Once lazy initialization.
type Registry interface {
	// Infrastructure
	Config() *config.Config
	Logger() *logrus.Logger
	Tracer() trace.Tracer

	// Identity (L1)
	IdentityPool() identity.Pool
	PrivilegedIdentityPool() identity.PrivilegedPool
	IdentityManager() identity.Manager
	IdentityHasher() identity.Hasher

	// Session (L1)
	SessionPool() session.Pool
	PrivilegedSessionPool() session.PrivilegedPool
	SessionManager() session.Manager

	// Token (L3)
	TokenPool() token.Pool
	TokenManager() token.Manager

	// Lockout (L3)
	LockoutManager() lockout.Manager

	// Authz (L2)
	AuthzEngine() *authz.Engine
	RolePool() role.Pool
	PrivilegedRolePool() role.PrivilegedPool
	RoleManager() role.Manager
	PolicyPool() policy.Pool
	PrivilegedPolicyPool() policy.PrivilegedPool
	PolicyManager() policy.Manager

	// Selfservice (L1)
	PasswordAuthenticator() *strategies.PasswordAuthenticator
	MFAManager() *strategies.ManagerImpl
	SelfserviceHandler() *selfservice.Handler

	// Courier (L3)
	Courier() courier.Courier

	// Audit (L3)
	AuditPool() audit.Pool
	AuditRecorder() audit.Recorder
	AuditManager() audit.Manager

	// Persistence
	Persister() Persister
	MigrateUp(ctx context.Context) error

	// Cache
	Cache() Cache

	// HTTP
	Router() HTTPHandler

	// Notification (stub)
	NotificationSender() any

	// Init initializes all services
	Init(ctx context.Context) error
}

// HTTPHandler is an interface abstracting HTTP server to avoid import cycle.
type HTTPHandler interface {
	Group(relativePath string, middleware ...any) HTTPGroup
	GET(relativePath string, handlers ...any)
	POST(relativePath string, handlers ...any)
	PUT(relativePath string, handlers ...any)
	PATCH(relativePath string, handlers ...any)
	DELETE(relativePath string, handlers ...any)
	Run(addr ...string) error
}

// HTTPGroup is an interface abstracting route groups.
type HTTPGroup interface {
	Handle(httpMethod, relativePath string, handlers ...any)
	GET(relativePath string, handlers ...any)
	POST(relativePath string, handlers ...any)
	PUT(relativePath string, handlers ...any)
	PATCH(relativePath string, handlers ...any)
	DELETE(relativePath string, handlers ...any)
	Use(middleware ...any)
}

// Persister is the database persistence interface.
type Persister interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
	Close(ctx context.Context) error
}

// Cache is the cache interface.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl int64) error
	Delete(ctx context.Context, key string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Ping(ctx context.Context) error
}
