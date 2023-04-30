// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/mycasbin"
)

type casbinRepositoryImpl struct {
	once     sync.Once
	db       *gorm.DB
	enforcer *casbin.SyncedEnforcer
}

// newCasbinRepository new SyncedEnforcer Repository
func newCasbinRepository(db *gorm.DB) repository.CasbinRepository {
	return &casbinRepositoryImpl{db: db}
}

func (c *casbinRepositoryImpl) SyncedEnforcer() *casbin.SyncedEnforcer {
	c.once.Do(func() {
		c.enforcer = mycasbin.Setup(c.db, model.TableNamePrefix, "casbin_rules")
	})
	return c.enforcer
}
