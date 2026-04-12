// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLiteOptions defines options for sqlite database.
type SQLiteOptions struct {
	Database string
	Logger   logger.Interface
}

// NewForSQLite creates a new gorm db instance with sqlite driver.
func NewForSQLite(opts *SQLiteOptions) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(opts.Database), &gorm.Config{
		Logger: opts.Logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SQLite doesn't have connection pool settings like MySQL
	// but we can set some basic pragmas for better performance
	sqlDB.SetMaxOpenConns(1) // SQLite doesn't support concurrent writes well

	return db, nil
}
