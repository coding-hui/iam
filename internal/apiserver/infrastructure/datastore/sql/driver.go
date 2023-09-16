// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/code"
	pkgUtils "github.com/coding-hui/iam/internal/pkg/utils"
	"github.com/coding-hui/iam/internal/pkg/utils/gormutil"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/fields"
	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/selection"
)

// Driver is a unified implementation of SQL driver of datastore
type Driver struct {
	Client *gorm.DB
	Cfg    config.Config
}

func (m *Driver) UserRepository() repository.UserRepository {
	return newUserRepository(m.Client)
}

func (m *Driver) CasbinRepository() repository.CasbinRepository {
	return newCasbinRepository(m.Client, m.Cfg.RedisOptions)
}

func (m *Driver) ResourceRepository() repository.ResourceRepository {
	return newResourceRepository(m.Client)
}

func (m *Driver) RoleRepository() repository.RoleRepository {
	return newRoleRepository(m.Client)
}

func (m *Driver) PolicyRepository() repository.PolicyRepository {
	return newPolicyRepository(m.Client)
}

func (m *Driver) OrganizationRepository() repository.OrganizationRepository {
	return newOrganizationRepository(m.Client)
}

func (m *Driver) Close() error {
	db, err := m.Client.DB()
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return db.Close()
}

func makeCondition(opts metav1.ListOptions) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var clauses []clause.Expression
		selector, _ := fields.ParseSelector(opts.FieldSelector)
		clauses = _applyFieldSelector(clauses, selector)
		return db.Clauses(clauses...)
	}
}

func paginate(opts metav1.ListOptions) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ol := gormutil.Unpointer(opts.Offset, opts.Limit)
		offset, limit := ol.Offset, ol.Limit
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(limit)
	}
}

// _toColumnName converts keys of the models to lowercase as the column name are in lowercase in the database
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
