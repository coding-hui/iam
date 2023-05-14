// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"k8s.io/klog/v2"

	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/code"
	genericoptions "github.com/coding-hui/iam/internal/pkg/options"
	"github.com/coding-hui/iam/pkg/db"

	"github.com/coding-hui/common/errors"
)

type mysqldb struct {
	client *gorm.DB
}

// GetMySQLFactory create mysql factory with the given config.
func GetMySQLFactory(_ context.Context, opts *genericoptions.MySQLOptions) (factory repository.Factory, lastErr error) {
	lastErr = createDatabase(opts)
	if lastErr != nil {
		return nil, fmt.Errorf("failed to create database, error: %w", lastErr)
	}

	var dbIns *gorm.DB
	options := &db.Options{
		Host:                  opts.Host,
		Username:              opts.Username,
		Password:              opts.Password,
		Database:              opts.Database,
		MaxIdleConnections:    opts.MaxIdleConnections,
		MaxOpenConnections:    opts.MaxOpenConnections,
		MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
		LogLevel:              opts.LogLevel,
	}
	dbIns, lastErr = db.New(options)

	m := &mysqldb{client: dbIns}
	if m == nil || lastErr != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", m, lastErr)
	}

	migrate(m.client)

	klog.Infof("create mysqldb datastore instance successful")

	return m, lastErr
}

func (m *mysqldb) UserRepository() repository.UserRepository {
	return newUserRepository(m.client)
}

func (m *mysqldb) CasbinRepository() repository.CasbinRepository {
	return newCasbinRepository(m.client)
}

func (m *mysqldb) ResourceRepository() repository.ResourceRepository {
	return newResourceRepository(m.client)
}

func (m *mysqldb) RoleRepository() repository.RoleRepository {
	return newRoleRepository(m.client)
}

func (m *mysqldb) PolicyRepository() repository.PolicyRepository {
	return newPolicyRepository(m.client)
}

func (m *mysqldb) Close() error {
	db, err := m.client.DB()
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return db.Close()
}
