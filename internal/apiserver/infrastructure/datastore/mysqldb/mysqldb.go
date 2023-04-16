// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"context"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s.io/klog/v2"

	repository2 "github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore"
	"github.com/coding-hui/iam/internal/pkg/code"

	"github.com/coding-hui/common/errors"
)

type mysqldb struct {
	client   *gorm.DB
	database string
}

// GetMySQLFactory create mysql factory with the given config.
func GetMySQLFactory(ctx context.Context, cfg datastore.Config) (repository2.Factory, error) {
	err := createDatabase(cfg)
	if err != nil {
		return nil, err
	}
	mysqlCfg := mysql.Config{
		DSN:                       cfg.URL, // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlCfg))
	if err != nil {
		return nil, err
	}
	db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
	migrate(db)
	m := &mysqldb{
		client:   db.WithContext(ctx),
		database: cfg.Database,
	}
	klog.Infof("create mysqldb datastore instance successful")

	return m, nil
}

func (m *mysqldb) UserRepository() repository2.UserRepository {
	return newUserRepository(m.client)
}

func (m *mysqldb) Close() error {
	db, err := m.client.DB()
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return db.Close()
}
