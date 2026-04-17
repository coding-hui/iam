// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Options holds database connection options.
type Options struct {
	MaxIdle     int
	MaxOpen     int
	MaxLifetime time.Duration
}

// Persister defines the interface for database operations.
// It provides transaction support and connection management.
type Persister interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
	Connection(ctx context.Context) *gorm.DB
	Close(ctx context.Context) error
	MigrateUp(ctx context.Context) error
}
