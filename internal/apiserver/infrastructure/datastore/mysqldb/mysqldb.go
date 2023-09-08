// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/code"
	pkgUtils "github.com/coding-hui/iam/internal/pkg/utils"
	"github.com/coding-hui/iam/pkg/db"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/fields"
	"github.com/coding-hui/common/selection"
)

type mysqldb struct {
	client *gorm.DB
	cfg    config.Config
}

// GetMySQLFactory create mysql factory with the given config.
func GetMySQLFactory(_ context.Context, c config.Config) (factory repository.Factory, lastErr error) {
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

	m := &mysqldb{client: dbIns, cfg: c}
	if m == nil || lastErr != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", m, lastErr)
	}

	migrate(m.client)

	log.Infof("create mysqldb datastore instance successful")

	return m, lastErr
}

func (m *mysqldb) UserRepository() repository.UserRepository {
	return newUserRepository(m.client)
}

func (m *mysqldb) CasbinRepository() repository.CasbinRepository {
	return newCasbinRepository(m.client, m.cfg.RedisOptions)
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

func (m *mysqldb) OrganizationRepository() repository.OrganizationRepository {
	return newOrganizationRepository(m.client)
}

func (m *mysqldb) Close() error {
	db, err := m.client.DB()
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return db.Close()
}

// _toColumnName converts the field to lowercase underscores
func _toColumnName(columnName string) string {
	if strings.HasPrefix(columnName, "metadata.") {
		columnName = columnName[9:]
	}
	return pkgUtils.CamelToUnderscore(columnName)
}

func _applyFieldSelector(clauses []clause.Expression, selector fields.Selector) []clause.Expression {
	if selector == nil || len(selector.Requirements()) == 0 {
		return nil
	}
	for _, req := range selector.Requirements() {
		field := _toColumnName(req.Field)
		operator := req.Operator
		value := req.Value

		var condition clause.Expression

		switch operator {
		case selection.Equals:
			condition = clause.Eq{
				Column: clause.Column{Name: field},
				Value:  value,
			}
		case selection.NotEquals:
			condition = clause.Neq{
				Column: clause.Column{Name: field},
				Value:  value,
			}
		default:
			// 忽略不支持的操作符
			log.Warnf("Unsupported field selector operator: %s", operator)
		}

		if condition != nil {
			clauses = append(clauses, condition)
		}
	}

	return clauses
}
