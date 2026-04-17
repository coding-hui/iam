// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/coding-hui/iam/internal/persistence"
)

// Persister implements persistence.Persister using GORM.
type Persister struct {
	db *gorm.DB
}

// NewMySQLPersister creates a new MySQL persister.
// dsn should be a MySQL DSN string like "user:password@tcp(host:port)/database?charset=utf8&parseTime=True&loc=Local"
func NewMySQLPersister(dsn string, opts *persistence.Options) (*Persister, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying DB: %w", err)
	}

	if opts != nil {
		sqlDB.SetMaxOpenConns(opts.MaxOpen)
		sqlDB.SetMaxIdleConns(opts.MaxIdle)
		sqlDB.SetConnMaxLifetime(opts.MaxLifetime)
	}

	return &Persister{db: db}, nil
}

// NewSQLitePersister creates a new SQLite persister.
func NewSQLitePersister(dsn string, opts *persistence.Options) (*Persister, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying DB: %w", err)
	}

	if opts != nil {
		sqlDB.SetMaxOpenConns(opts.MaxOpen)
		sqlDB.SetMaxIdleConns(opts.MaxIdle)
	}

	return &Persister{db: db}, nil
}

// Transaction executes fn within a database transaction.
func (p *Persister) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txContextKey{}, tx))
	})
}

// Connection returns the database connection.
func (p *Persister) Connection(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok {
		return tx
	}
	return p.db.WithContext(ctx)
}

// Close closes the database connection.
func (p *Persister) Close(ctx context.Context) error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

type txContextKey struct{}

// MigrateUp runs database migrations.
func (p *Persister) MigrateUp(ctx context.Context) error {
	models := []any{
		&IdentityModel{},
		&SessionModel{},
		&RoleModel{},
		&PolicyModel{},
		&TokenModel{},
		&AuditEventModel{},
		&SecretKey{},
	}

	return p.db.AutoMigrate(models...)
}

// Ensure Persister implements persistence.Persister.
var _ persistence.Persister = (*Persister)(nil)
