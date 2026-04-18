// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/coding-hui/iam/internal/audit"
	"github.com/coding-hui/iam/internal/authz"
	"github.com/coding-hui/iam/internal/authz/policy"
	"github.com/coding-hui/iam/internal/authz/role"
	"github.com/coding-hui/iam/internal/cache"
	"github.com/coding-hui/iam/internal/config"
	"github.com/coding-hui/iam/internal/identity"
	"github.com/coding-hui/iam/internal/identity/lockout"
	"github.com/coding-hui/iam/internal/identity/session"
	"github.com/coding-hui/iam/internal/identity/token"
	"github.com/coding-hui/iam/internal/persistence/sql"
	"github.com/coding-hui/iam/internal/selfservice"
	"github.com/coding-hui/iam/internal/selfservice/courier"
	"github.com/coding-hui/iam/internal/selfservice/strategies"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// initOnce provides thread-safe lazy initialization.
type initOnce[T any] struct {
	mu  sync.Mutex
	on  sync.Once
	fn  func() T
	val T
}

// Get returns the singleton value, initializing it on first call.
func (o *initOnce[T]) Get() T {
	o.mu.Lock()
	o.on.Do(func() { o.val = o.fn() })
	o.mu.Unlock()
	return o.val
}

// RegistryDefault provides the default implementation of Registry.
// It uses sync.Once for thread-safe lazy initialization of all services.
type RegistryDefault struct {
	config *config.Config

	persister initOnce[*sql.Persister]
	cache     initOnce[Cache]
	logger    *logrus.Logger
	tracer    trace.Tracer
	router    HTTPHandler

	init    sync.Once
	initErr error

	identityPool           initOnce[identity.Pool]
	identityPrivilegedPool initOnce[identity.PrivilegedPool]
	identityManager        initOnce[identity.Manager]
	identityHasher         identity.Hasher

	sessionPool           initOnce[session.Pool]
	sessionPrivilegedPool initOnce[session.PrivilegedPool]
	sessionManager        initOnce[session.Manager]

	tokenPool           initOnce[token.Pool]
	tokenPrivilegedPool initOnce[token.PrivilegedPool]
	tokenManager        initOnce[token.Manager]

	lockoutManager initOnce[lockout.Manager]

	rolePool           initOnce[role.Pool]
	rolePrivilegedPool initOnce[role.PrivilegedPool]
	roleManager        initOnce[role.Manager]

	policyPool           initOnce[policy.Pool]
	policyPrivilegedPool initOnce[policy.PrivilegedPool]
	policyManager        initOnce[policy.Manager]

	authzEngine *authz.Engine

	passwordAuthenticator *strategies.PasswordAuthenticator
	mfaManager            *strategies.ManagerImpl

	selfserviceHandler initOnce[*selfservice.Handler]
	courierInstance    initOnce[courier.Courier]

	auditPool     initOnce[audit.Pool]
	auditRecorder initOnce[audit.Recorder]
	auditManager  initOnce[audit.Manager]
}

// NewRegistry creates a new Registry instance with the given configuration.
func NewRegistry(cfg *config.Config) *RegistryDefault {
	return &RegistryDefault{
		config: cfg,
		logger: newLogger(),
		tracer: otel.Tracer("iam"),
	}
}

// Config returns the application configuration.
func (r *RegistryDefault) Config() *config.Config {
	return r.config
}

// Logger returns the application logger.
func (r *RegistryDefault) Logger() *logrus.Logger {
	return r.logger
}

// Tracer returns the OpenTelemetry tracer.
func (r *RegistryDefault) Tracer() trace.Tracer {
	return r.tracer
}

// Router returns the HTTP router.
func (r *RegistryDefault) Router() HTTPHandler {
	return r.router
}

// Init initializes all lazy-loaded services.
func (r *RegistryDefault) Init(ctx context.Context) error {
	r.init.Do(func() {
		r.initErr = r.initialize(ctx)
	})
	return r.initErr
}

func (r *RegistryDefault) initialize(ctx context.Context) error {
	r.persister = initOnce[*sql.Persister]{
		fn: func() *sql.Persister {
			return r.newPersister()
		},
	}

	r.cache = initOnce[Cache]{
		fn: func() Cache {
			return r.newCache()
		},
	}

	r.identityHasher = identity.NewArgon2idHasher()

	r.identityPool = initOnce[identity.Pool]{
		fn: func() identity.Pool {
			p := r.persister.Get()
			return identity.NewPool(sql.NewIdentityPool(p))
		},
	}

	r.identityPrivilegedPool = initOnce[identity.PrivilegedPool]{
		fn: func() identity.PrivilegedPool {
			p := r.persister.Get()
			return identity.NewPrivilegedPool(sql.NewIdentityPool(p))
		},
	}

	r.identityManager = initOnce[identity.Manager]{
		fn: func() identity.Manager {
			return identity.NewManagerImpl(
				r.identityPool.Get(),
				r.identityPrivilegedPool.Get(),
				r.identityHasher,
			)
		},
	}

	r.sessionPool = initOnce[session.Pool]{
		fn: func() session.Pool {
			p := r.persister.Get()
			return session.NewPool(sql.NewSessionPool(p))
		},
	}

	r.sessionPrivilegedPool = initOnce[session.PrivilegedPool]{
		fn: func() session.PrivilegedPool {
			p := r.persister.Get()
			return session.NewPrivilegedPool(sql.NewSessionPool(p))
		},
	}

	r.sessionManager = initOnce[session.Manager]{
		fn: func() session.Manager {
			return session.NewManagerImpl(
				r.sessionPool.Get(),
				r.sessionPrivilegedPool.Get(),
			)
		},
	}

	r.tokenPool = initOnce[token.Pool]{
		fn: func() token.Pool {
			p := r.persister.Get()
			return token.NewPool(sql.NewTokenPool(p))
		},
	}

	r.tokenPrivilegedPool = initOnce[token.PrivilegedPool]{
		fn: func() token.PrivilegedPool {
			p := r.persister.Get()
			return token.NewPrivilegedPool(sql.NewTokenPool(p))
		},
	}

	r.tokenManager = initOnce[token.Manager]{
		fn: func() token.Manager {
			return token.NewManagerImpl(
				r.tokenPool.Get(),
				r.tokenPrivilegedPool.Get(),
			)
		},
	}

	r.lockoutManager = initOnce[lockout.Manager]{
		fn: func() lockout.Manager {
			return lockout.NewManager(5, 15*60*1000*1000*1000)
		},
	}

	r.rolePool = initOnce[role.Pool]{
		fn: func() role.Pool {
			p := r.persister.Get()
			return role.NewPool(sql.NewRolePool(p))
		},
	}

	r.rolePrivilegedPool = initOnce[role.PrivilegedPool]{
		fn: func() role.PrivilegedPool {
			p := r.persister.Get()
			return role.NewPrivilegedPool(sql.NewRolePool(p))
		},
	}

	r.roleManager = initOnce[role.Manager]{
		fn: func() role.Manager {
			return role.NewManagerImpl(
				r.rolePool.Get(),
				r.rolePrivilegedPool.Get(),
			)
		},
	}

	r.policyPool = initOnce[policy.Pool]{
		fn: func() policy.Pool {
			p := r.persister.Get()
			return policy.NewPool(sql.NewPolicyPool(p))
		},
	}

	r.policyPrivilegedPool = initOnce[policy.PrivilegedPool]{
		fn: func() policy.PrivilegedPool {
			p := r.persister.Get()
			return policy.NewPrivilegedPool(sql.NewPolicyPool(p))
		},
	}

	r.policyManager = initOnce[policy.Manager]{
		fn: func() policy.Manager {
			return policy.NewManagerImpl(
				r.policyPool.Get(),
				r.policyPrivilegedPool.Get(),
			)
		},
	}

	r.authzEngine = authz.NewEngine()

	// Selfservice (L1) - Strategies
	r.passwordAuthenticator = strategies.NewPasswordAuthenticator(
		r.identityPrivilegedPool.Get(),
		r.sessionPrivilegedPool.Get(),
		r.identityHasher,
	)

	r.mfaManager = strategies.NewManagerImpl()

	// Selfservice Handler
	r.selfserviceHandler = initOnce[*selfservice.Handler]{
		fn: func() *selfservice.Handler {
			return selfservice.NewHandler(
				r.passwordAuthenticator,
				r.mfaManager,
			)
		},
	}

	// Courier (L3)
	r.courierInstance = initOnce[courier.Courier]{
		fn: func() courier.Courier {
			return courier.NewCourier()
		},
	}

	// Audit (L3)
	r.auditPool = initOnce[audit.Pool]{
		fn: func() audit.Pool {
			p := r.persister.Get()
			return audit.NewPool(sql.NewAuditPool(p))
		},
	}

	r.auditRecorder = initOnce[audit.Recorder]{
		fn: func() audit.Recorder {
			p := r.persister.Get()
			return audit.NewRecorder(sql.NewAuditPool(p))
		},
	}

	r.auditManager = initOnce[audit.Manager]{
		fn: func() audit.Manager {
			return audit.NewManagerImpl(
				r.auditPool.Get(),
				r.auditRecorder.Get(),
			)
		},
	}

	return nil
}

func (r *RegistryDefault) newPersister() *sql.Persister {
	dbConfig := r.config.Database

	var p *sql.Persister
	var err error

	switch dbConfig.Driver {
	case "mysql":
		p, err = sql.NewMySQLPersister(dbConfig.DSN, nil)
	case "sqlite":
		p, err = sql.NewSQLitePersister(dbConfig.DSN, nil)
	default:
		p, err = sql.NewSQLitePersister(dbConfig.DSN, nil)
	}

	if err != nil {
		panic("failed to initialize persister: " + err.Error())
	}

	return p
}

func (r *RegistryDefault) newCache() Cache {
	// Use in-memory cache for simplified setup
	return cache.NewMemoryCache()
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return logger
}

// Persister returns the database persister.
func (r *RegistryDefault) Persister() Persister {
	return r.persister.Get()
}

// Cache returns the cache instance.
func (r *RegistryDefault) Cache() Cache {
	return r.cache.Get()
}

// MigrateUp runs database migrations.
func (r *RegistryDefault) MigrateUp(ctx context.Context) error {
	p := r.persister.Get()
	if p == nil {
		return fmt.Errorf("persister not initialized")
	}
	return p.MigrateUp(ctx)
}

// IdentityPool returns the identity pool.
func (r *RegistryDefault) IdentityPool() identity.Pool {
	return r.identityPool.Get()
}

// PrivilegedIdentityPool returns the privileged identity pool.
func (r *RegistryDefault) PrivilegedIdentityPool() identity.PrivilegedPool {
	return r.identityPrivilegedPool.Get()
}

// IdentityManager returns the identity manager.
func (r *RegistryDefault) IdentityManager() identity.Manager {
	return r.identityManager.Get()
}

// IdentityHasher returns the identity hasher.
func (r *RegistryDefault) IdentityHasher() identity.Hasher {
	return r.identityHasher
}

// SessionPool returns the session pool.
func (r *RegistryDefault) SessionPool() session.Pool {
	return r.sessionPool.Get()
}

// PrivilegedSessionPool returns the privileged session pool.
func (r *RegistryDefault) PrivilegedSessionPool() session.PrivilegedPool {
	return r.sessionPrivilegedPool.Get()
}

// SessionManager returns the session manager.
func (r *RegistryDefault) SessionManager() session.Manager {
	return r.sessionManager.Get()
}

// TokenPool returns the token pool.
func (r *RegistryDefault) TokenPool() token.Pool {
	return r.tokenPool.Get()
}

// TokenManager returns the token manager.
func (r *RegistryDefault) TokenManager() token.Manager {
	return r.tokenManager.Get()
}

// LockoutManager returns the lockout manager.
func (r *RegistryDefault) LockoutManager() lockout.Manager {
	return r.lockoutManager.Get()
}

// RolePool returns the role pool.
func (r *RegistryDefault) RolePool() role.Pool {
	return r.rolePool.Get()
}

// PrivilegedRolePool returns the privileged role pool.
func (r *RegistryDefault) PrivilegedRolePool() role.PrivilegedPool {
	return r.rolePrivilegedPool.Get()
}

// RoleManager returns the role manager.
func (r *RegistryDefault) RoleManager() role.Manager {
	return r.roleManager.Get()
}

// PolicyPool returns the policy pool.
func (r *RegistryDefault) PolicyPool() policy.Pool {
	return r.policyPool.Get()
}

// PrivilegedPolicyPool returns the privileged policy pool.
func (r *RegistryDefault) PrivilegedPolicyPool() policy.PrivilegedPool {
	return r.policyPrivilegedPool.Get()
}

// PolicyManager returns the policy manager.
func (r *RegistryDefault) PolicyManager() policy.Manager {
	return r.policyManager.Get()
}

// AuthzEngine returns the authz engine.
func (r *RegistryDefault) AuthzEngine() *authz.Engine {
	return r.authzEngine
}

// PasswordAuthenticator returns the password authenticator.
func (r *RegistryDefault) PasswordAuthenticator() *strategies.PasswordAuthenticator {
	return r.passwordAuthenticator
}

// MFAManager returns the MFA manager.
func (r *RegistryDefault) MFAManager() *strategies.ManagerImpl {
	return r.mfaManager
}

// SelfserviceHandler returns the selfservice handler.
func (r *RegistryDefault) SelfserviceHandler() *selfservice.Handler {
	return r.selfserviceHandler.Get()
}

// Courier returns the courier.
func (r *RegistryDefault) Courier() courier.Courier {
	return r.courierInstance.Get()
}

// AuditPool returns the audit pool.
func (r *RegistryDefault) AuditPool() audit.Pool {
	return r.auditPool.Get()
}

// AuditRecorder returns the audit recorder.
func (r *RegistryDefault) AuditRecorder() audit.Recorder {
	return r.auditRecorder.Get()
}

// AuditManager returns the audit manager.
func (r *RegistryDefault) AuditManager() audit.Manager {
	return r.auditManager.Get()
}

// NotificationSender returns nil - not yet implemented.
func (r *RegistryDefault) NotificationSender() any { return nil }

// Ensure RegistryDefault implements Registry.
var _ Registry = (*RegistryDefault)(nil)
