// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

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
	Host            string   `mapstructure:"host"`
	Port            int      `mapstructure:"port"`
	Mode            string   `mapstructure:"mode"` // debug, release, test
	Healthz         bool     `mapstructure:"healthz"`
	Middlewares     []string `mapstructure:"middlewares"`
	StaticLocations string   `mapstructure:"static-locations"`
	TLSCert         string   `mapstructure:"tls-cert"`
	TLSKey          string   `mapstructure:"tls-key"`
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Driver      string        `mapstructure:"driver"`
	DSN         string        `mapstructure:"dsn"`
	MaxIdle     int           `mapstructure:"max-idle"`
	MaxOpen     int           `mapstructure:"max-open"`
	MaxLifetime time.Duration `mapstructure:"max-lifetime"`
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool-size"`
}

// AuthzConfig holds authorization configuration.
type AuthzConfig struct {
	EnforcerPath string `mapstructure:"enforcer-path"`
}

// JWTConfig holds JWT token configuration.
type JWTConfig struct {
	Secret            string        `mapstructure:"secret"`
	Expiration        time.Duration `mapstructure:"expiration"`
	RefreshExpiration time.Duration `mapstructure:"refresh-expiration"`
}

// MFAConfig holds MFA configuration.
type MFAConfig struct {
	TOTPIssuer string `mapstructure:"totp-issuer"`
}

// AuditConfig holds audit configuration.
type AuditConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// NotifyConfig holds notification configuration.
type NotifyConfig struct {
	Email EmailConfig `mapstructure:"email"`
}

// EmailConfig holds email notification configuration.
type EmailConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}
