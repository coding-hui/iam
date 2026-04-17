// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package driver

import (
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Authz    AuthzConfig
	JWT      JWTConfig
	MFA      MFAConfig
	Audit    AuditConfig
	Notify   NotifyConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host            string
	Port            int
	Mode            string // debug, release, test
	Healthz         bool
	Middlewares     []string
	StaticLocations string
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Driver      string // mysql, sqlite
	DSN         string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime time.Duration
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

// AuthzConfig holds authorization configuration.
type AuthzConfig struct {
	EnforcerPath string
}

// JWTConfig holds JWT token configuration.
type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

// MFAConfig holds MFA configuration.
type MFAConfig struct {
	TOTPIssuer string
}

// AuditConfig holds audit configuration.
type AuditConfig struct {
	Enabled bool
}

// NotifyConfig holds notification configuration.
type NotifyConfig struct {
	Email EmailConfig
}

// EmailConfig holds email notification configuration.
type EmailConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Username string
	Password string
	From     string
}
