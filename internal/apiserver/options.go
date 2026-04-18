// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	cliflag "github.com/coding-hui/common/cli/flag"
	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/config"
	"github.com/coding-hui/iam/pkg/app"
)

// Options embeds config.Config and implements CliOptions, CompleteableOptions, PrintableOptions.
type Options struct {
	config.Config `mapstructure:",squash"`
}

var globalOptions *Options

// GetOptions returns the global options instance.
func GetOptions() *Options {
	return globalOptions
}

// NewOptions creates a new Options instance with default values.
func NewOptions() *Options {
	return &Options{
		Config: config.Config{
			Server: config.ServerConfig{
				Port:    8080,
				Mode:    "debug",
				Healthz: true,
			},
			Database: config.DatabaseConfig{
				Driver:  "sqlite",
				MaxIdle: 10,
				MaxOpen: 100,
			},
			Redis: config.RedisConfig{
				Addr:     "localhost:6379",
				DB:       0,
				PoolSize: 100,
			},
			MFA: config.MFAConfig{
				TOTPIssuer: "iam",
			},
			Audit: config.AuditConfig{
				Enabled: true,
			},
		},
	}
}

// Flags returns the flag sets for all configuration options.
func (o *Options) Flags() cliflag.NamedFlagSets {
	var fss cliflag.NamedFlagSets

	o.addServerFlags(fss.FlagSet("Server"))
	o.addDatabaseFlags(fss.FlagSet("Database"))
	o.addRedisFlags(fss.FlagSet("Redis"))
	o.addAuthzFlags(fss.FlagSet("Authz"))
	o.addJWTFlags(fss.FlagSet("JWT"))
	o.addMFAFlags(fss.FlagSet("MFA"))
	o.addAuditFlags(fss.FlagSet("Audit"))
	o.addNotifyFlags(fss.FlagSet("Notify"))

	return fss
}

// Complete sets default values for options that cannot be expressed as zero values.
func (o *Options) Complete() error {
	if o.Server.Port == 0 {
		o.Server.Port = 8080
	}
	if o.Server.Mode == "" {
		o.Server.Mode = "debug"
	}
	if o.Database.MaxIdle == 0 {
		o.Database.MaxIdle = 10
	}
	if o.Database.MaxOpen == 0 {
		o.Database.MaxOpen = 100
	}
	if o.MFA.TOTPIssuer == "" {
		o.MFA.TOTPIssuer = "iam"
	}
	globalOptions = o
	return nil
}

// Validate validates the required fields.
func (o *Options) Validate() []error {
	errs := []error{}

	if o.Database.DSN == "" {
		errs = append(errs, errors.New("database DSN is required"))
	}
	if o.Server.Port < 1 || o.Server.Port > 65535 {
		errs = append(errs, errors.New("server port must be between 1 and 65535"))
	}

	return errs
}

// String returns a string representation of the configuration.
func (o *Options) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Server: port=%d, mode=%s\n", o.Server.Port, o.Server.Mode))
	sb.WriteString(fmt.Sprintf("Database: driver=%s, dsn=%s\n", o.Database.Driver, o.Database.DSN))
	sb.WriteString(fmt.Sprintf("Redis: addr=%s, db=%d\n", o.Redis.Addr, o.Redis.DB))
	return sb.String()
}

func (o *Options) addServerFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Server.Host, "server.host", o.Server.Host, "Server host address")
	fs.IntVar(&o.Server.Port, "server.port", o.Server.Port, "Server port")
	fs.StringVar(&o.Server.Mode, "server.mode", o.Server.Mode, "Server mode (debug, release, test)")
	fs.BoolVar(&o.Server.Healthz, "server.healthz", o.Server.Healthz, "Enable health check endpoint")
	fs.StringSliceVar(&o.Server.Middlewares, "server.middlewares", o.Server.Middlewares, "Server middlewares")
	fs.StringVar(&o.Server.TLSCert, "server.tls-cert", o.Server.TLSCert, "TLS certificate file path")
	fs.StringVar(&o.Server.TLSKey, "server.tls-key", o.Server.TLSKey, "TLS key file path")
}

func (o *Options) addDatabaseFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Database.Driver, "database.driver", o.Database.Driver, "Database driver (mysql, sqlite)")
	fs.StringVar(&o.Database.DSN, "database.dsn", o.Database.DSN, "Database DSN")
	fs.IntVar(&o.Database.MaxIdle, "database.max-idle", o.Database.MaxIdle, "Max idle connections")
	fs.IntVar(&o.Database.MaxOpen, "database.max-open", o.Database.MaxOpen, "Max open connections")
	fs.DurationVar(&o.Database.MaxLifetime, "database.max-lifetime", o.Database.MaxLifetime, "Max connection lifetime")
}

func (o *Options) addRedisFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Redis.Addr, "redis.addr", o.Redis.Addr, "Redis address")
	fs.StringVar(&o.Redis.Password, "redis.password", o.Redis.Password, "Redis password")
	fs.IntVar(&o.Redis.DB, "redis.db", o.Redis.DB, "Redis database number")
	fs.IntVar(&o.Redis.PoolSize, "redis.pool-size", o.Redis.PoolSize, "Redis pool size")
}

func (o *Options) addAuthzFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Authz.EnforcerPath, "authz.enforcer-path", o.Authz.EnforcerPath, "Casbin enforcer policy path")
}

func (o *Options) addJWTFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.JWT.Secret, "jwt.secret", o.JWT.Secret, "JWT secret")
	fs.DurationVar(&o.JWT.Expiration, "jwt.expiration", o.JWT.Expiration, "JWT expiration time")
	fs.DurationVar(&o.JWT.RefreshExpiration, "jwt.refresh-expiration", o.JWT.RefreshExpiration, "JWT refresh expiration time")
}

func (o *Options) addMFAFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.MFA.TOTPIssuer, "mfa.totp-issuer", o.MFA.TOTPIssuer, "TOTP issuer")
}

func (o *Options) addAuditFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.Audit.Enabled, "audit.enabled", o.Audit.Enabled, "Enable audit logging")
}

func (o *Options) addNotifyFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.Notify.Email.Enabled, "notify.email.enabled", o.Notify.Email.Enabled, "Enable email notifications")
	fs.StringVar(&o.Notify.Email.Host, "notify.email.host", o.Notify.Email.Host, "Email SMTP host")
	fs.IntVar(&o.Notify.Email.Port, "notify.email.port", o.Notify.Email.Port, "Email SMTP port")
	fs.StringVar(&o.Notify.Email.Username, "notify.email.username", o.Notify.Email.Username, "Email SMTP username")
	fs.StringVar(&o.Notify.Email.Password, "notify.email.password", o.Notify.Email.Password, "Email SMTP password")
	fs.StringVar(&o.Notify.Email.From, "notify.email.from", o.Notify.Email.From, "Email sender address")
}

// Ensure Options implements app.CliOptions
var _ app.CliOptions = (*Options)(nil)
var _ app.CompleteableOptions = (*Options)(nil)
var _ app.PrintableOptions = (*Options)(nil)
