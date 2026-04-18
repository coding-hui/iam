// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/api/middleware"
	"github.com/coding-hui/iam/internal/audit"
	"github.com/coding-hui/iam/internal/authz"
	"github.com/coding-hui/iam/internal/authz/policy"
	"github.com/coding-hui/iam/internal/authz/role"
	"github.com/coding-hui/iam/internal/driver"
	"github.com/coding-hui/iam/internal/identity"
	"github.com/coding-hui/iam/internal/identity/lockout"
	"github.com/coding-hui/iam/internal/identity/session"
	"github.com/coding-hui/iam/internal/identity/token"
	"github.com/coding-hui/iam/internal/selfservice/courier"
)

// NewRouter creates a new Gin router with all middleware and routes configured.
func NewRouter(reg driver.Registry) *gin.Engine {
	r := gin.New()

	// Apply JSON logger config for consistent logging
	r.Use(gin.LoggerWithConfig(middleware.GetLoggerConfig(nil, nil)))

	// Global middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLog())
	r.Use(middleware.GinRecovery())
	r.Use(middleware.Cors())

	// Health check endpoints
	r.GET("/healthz", healthHandler(reg))
	r.GET("/ping", pingHandler)

	// Register API routes
	registerRoutes(r, reg)

	return r
}

func registerRoutes(r *gin.Engine, reg driver.Registry) {
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		identityHandler := identity.NewHandler(reg.IdentityManager())
		v1.POST("/identities", identityHandler.Create)
		v1.GET("/identities", identityHandler.List)
		v1.GET("/identities/:id", identityHandler.Get)
		v1.PATCH("/identities/:id", identityHandler.Update)
		v1.DELETE("/identities/:id", identityHandler.Delete)
		v1.POST("/identities/:id/credentials", identityHandler.AddCredentials)
		v1.DELETE("/identities/:id/credentials/:type", identityHandler.DeleteCredentials)

		sessionHandler := session.NewHandler(reg.SessionManager())
		v1.POST("/sessions", sessionHandler.Create)
		v1.GET("/sessions", sessionHandler.List)
		v1.GET("/sessions/:id", sessionHandler.Get)
		v1.DELETE("/sessions/:id", sessionHandler.Revoke)
		v1.DELETE("/sessions", sessionHandler.RevokeAll)
		v1.PATCH("/sessions/:id", sessionHandler.Extend)

		selfserviceHandler := reg.SelfserviceHandler()
		v1.POST("/login", selfserviceHandler.Login)
		v1.POST("/mfa/totp/setup", selfserviceHandler.SetupTOTP)
		v1.POST("/mfa/totp/verify", selfserviceHandler.VerifyTOTP)
		v1.POST("/mfa/totp/disable", selfserviceHandler.DisableTOTP)

		roleHandler := role.NewHandler(reg.RoleManager())
		v1.POST("/roles", roleHandler.Create)
		v1.GET("/roles", roleHandler.List)
		v1.GET("/roles/:id", roleHandler.Get)
		v1.PATCH("/roles/:id", roleHandler.Update)
		v1.DELETE("/roles/:id", roleHandler.Delete)

		policyHandler := policy.NewHandler(reg.PolicyManager())
		v1.POST("/policies", policyHandler.Create)
		v1.GET("/policies", policyHandler.List)
		v1.GET("/policies/:id", policyHandler.Get)
		v1.PATCH("/policies/:id", policyHandler.Update)
		v1.DELETE("/policies/:id", policyHandler.Delete)

		authzHandler := authz.NewHandler(reg.AuthzEngine())
		v1.POST("/authz/check", authzHandler.Check)

		tokenHandler := token.NewHandler(reg.TokenManager())
		v1.POST("/tokens", tokenHandler.Create)
		v1.POST("/tokens/introspect", tokenHandler.Introspect)
		v1.DELETE("/tokens/:id", tokenHandler.Revoke)

		auditHandler := audit.NewHandler(reg.AuditManager())
		v1.GET("/audit/events", auditHandler.List)

		lockoutHandler := lockout.NewHandler(reg.LockoutManager())
		v1.GET("/lockout/:identifier", lockoutHandler.CheckStatus)
		v1.POST("/lockout/:identifier/unlock", lockoutHandler.Unlock)

		courierHandler := courier.NewHandler(reg.Courier())
		v1.POST("/webhooks", courierHandler.Register)
		v1.DELETE("/webhooks", courierHandler.Unregister)
		v1.POST("/webhooks/events", courierHandler.SendEvent)
	}
}

// healthHandler handles health check requests.
func healthHandler(reg driver.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	}
}

// pingHandler handles ping requests.
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
