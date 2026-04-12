// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sqlitedb

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

type sqlitedb struct {
	sql.Driver
}

// New create sqlite factory with the given config.
func New(ctx context.Context, c config.Config) (factory repository.Factory, lastErr error) {
	var dbIns *gorm.DB
	options := &db.SQLiteOptions{
		Database: c.SQLiteOptions.Database,
	}
	dbIns, lastErr = db.NewForSQLite(options)
	if dbIns == nil {
		return nil, fmt.Errorf("failed to create sqlite datastore instance: %w", lastErr)
	}

	m := &sqlitedb{Driver: sql.NewDriver(dbIns.WithContext(ctx), c)}
	if lastErr != nil {
		return nil, fmt.Errorf("failed to get sqlite store factory, sqliteFactory: %+v, error: %w", m, lastErr)
	}

	migrate(m.Client.WithCtx(ctx))

	log.Infof("create sqlitedb datastore instance successful")

	return m, lastErr
}
