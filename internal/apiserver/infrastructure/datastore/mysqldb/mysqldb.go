// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore/sql"
	"github.com/coding-hui/iam/pkg/db"
	"github.com/coding-hui/iam/pkg/log"
)

type mysqldb struct {
	sql.Driver
}

// New create mysql factory with the given config.
func New(ctx context.Context, c config.Config) (factory repository.Factory, lastErr error) {
	lastErr = createDatabase(c.MySQLOptions)
	if lastErr != nil {
		return nil, fmt.Errorf("failed to create database, error: %w", lastErr)
	}

	var dbIns *gorm.DB
	options := &db.Options{
		Host:                  c.MySQLOptions.Host,
		Username:              c.MySQLOptions.Username,
		Password:              c.MySQLOptions.Password,
		Database:              c.MySQLOptions.Database,
		MaxIdleConnections:    c.MySQLOptions.MaxIdleConnections,
		MaxOpenConnections:    c.MySQLOptions.MaxOpenConnections,
		MaxConnectionLifeTime: c.MySQLOptions.MaxConnectionLifeTime,
		LogLevel:              c.MySQLOptions.LogLevel,
	}
	dbIns, lastErr = db.New(options)

	m := &mysqldb{
		Driver: sql.Driver{
			Client: dbIns.WithContext(ctx),
			Cfg:    c,
		},
	}
	if m == nil || lastErr != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", m, lastErr)
	}

	migrate(m.Driver.Client)

	log.Infof("create mysqldb datastore instance successful")

	return m, lastErr
}
