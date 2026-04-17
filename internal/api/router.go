// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/api/middleware"
	"github.com/coding-hui/iam/internal/audit"
	"github.com/coding-hui/iam/internal/authn"
	"github.com/coding-hui/iam/internal/authz"
	"github.com/coding-hui/iam/internal/driver"
	"github.com/coding-hui/iam/internal/identity"
	"github.com/coding-hui/iam/internal/lockout"
	"github.com/coding-hui/iam/internal/mfa"
	"github.com/coding-hui/iam/internal/policy"
	"github.com/coding-hui/iam/internal/role"
	"github.com/coding-hui/iam/internal/session"
	"github.com/coding-hui/iam/internal/token"
	"github.com/coding-hui/iam/internal/webhook"
)

// NewRouter creates a new Gin router with all middleware and routes configured.
func NewRouter(reg driver.Registry) *gin.Engine {
	// Set Gin mode based on config
	mode := reg.Config().Server.Mode
	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// Global middleware
	logger := reg.Logger()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS())

	// Health check endpoints
	r.GET("/healthz", healthHandler)
	r.GET("/ping", pingHandler)

	// Register API routes
	registerRoutes(r, reg)

	return r
}

func registerRoutes(r *gin.Engine, reg driver.Registry) {
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Identity routes (L1)
		identityHandler := identity.NewHandler(reg.IdentityManager().(identity.Manager))
		v1.POST("/identities", identityHandler.Create)
		v1.GET("/identities", identityHandler.List)
		v1.GET("/identities/:id", identityHandler.Get)
		v1.PATCH("/identities/:id", identityHandler.Update)
		v1.DELETE("/identities/:id", identityHandler.Delete)
		v1.POST("/identities/:id/credentials", identityHandler.AddCredentials)
		v1.DELETE("/identities/:id/credentials/:type", identityHandler.DeleteCredentials)

		// Session routes (L1)
		sessionHandler := session.NewHandler(reg.SessionManager().(session.Manager))
		v1.POST("/sessions", sessionHandler.Create)
		v1.GET("/sessions", sessionHandler.List)
		v1.GET("/sessions/:id", sessionHandler.Get)
		v1.DELETE("/sessions/:id", sessionHandler.Revoke)
		v1.DELETE("/sessions", sessionHandler.RevokeAll)
		v1.PATCH("/sessions/:id", sessionHandler.Extend)

		// Authn routes (L1) - Login
		authnHandler := authn.NewHandler(reg.Authenticator().(authn.Authenticator))
		v1.POST("/login", authnHandler.Login)

		// Role routes (L2)
		roleHandler := role.NewHandler(reg.RoleManager().(role.Manager))
		v1.POST("/roles", roleHandler.Create)
		v1.GET("/roles", roleHandler.List)
		v1.GET("/roles/:id", roleHandler.Get)
		v1.PATCH("/roles/:id", roleHandler.Update)
		v1.DELETE("/roles/:id", roleHandler.Delete)

		// Policy routes (L2)
		policyHandler := policy.NewHandler(reg.PolicyManager().(policy.Manager))
		v1.POST("/policies", policyHandler.Create)
		v1.GET("/policies", policyHandler.List)
		v1.GET("/policies/:id", policyHandler.Get)
		v1.PATCH("/policies/:id", policyHandler.Update)
		v1.DELETE("/policies/:id", policyHandler.Delete)

		// Authz routes (L2)
		authzHandler := authz.NewHandler(reg.AuthzEngine().(*authz.Engine))
		v1.POST("/authz/check", authzHandler.Check)

		// Token routes (L3)
		tokenHandler := token.NewHandler(reg.TokenManager().(token.Manager))
		v1.POST("/tokens", tokenHandler.Create)
		v1.POST("/tokens/introspect", tokenHandler.Introspect)
		v1.DELETE("/tokens/:id", tokenHandler.Revoke)

		// Audit routes (L3)
		auditHandler := audit.NewHandler(reg.AuditManager().(audit.Manager))
		v1.GET("/audit/events", auditHandler.List)

		// MFA routes (L3) - TOTP
		mfaHandler := mfa.NewHandler(reg.MFAManager().(mfa.Manager))
		v1.POST("/mfa/totp/setup", mfaHandler.SetupTOTP)
		v1.POST("/mfa/totp/verify", mfaHandler.VerifyTOTP)
		v1.POST("/mfa/totp/disable", mfaHandler.DisableTOTP)

		// Lockout routes (L3)
		lockoutHandler := lockout.NewHandler(reg.LockoutManager().(lockout.Manager))
		v1.GET("/lockout/:identifier", lockoutHandler.CheckStatus)
		v1.POST("/lockout/:identifier/unlock", lockoutHandler.Unlock)

		// Webhook routes (L3)
		webhookHandler := webhook.NewHandler(reg.WebhookManager().(webhook.Manager))
		v1.POST("/webhooks", webhookHandler.Register)
		v1.DELETE("/webhooks", webhookHandler.Unregister)
		v1.POST("/webhooks/events", webhookHandler.SendEvent)
	}
}

// healthHandler handles health check requests.
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// pingHandler handles ping requests.
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
