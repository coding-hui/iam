// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/pkg/log"
	genericoptions "github.com/coding-hui/iam/pkg/options"
)

// mysqlEmptyDsn msyql empty dsn for create databases.
func mysqlEmptyDsn(opts *genericoptions.MySQLOptions) string {
	host := opts.Host
	idx := strings.Index(opts.Host, "/")
	if idx > 0 {
		host = opts.Host[:idx+1]
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/", opts.Username, opts.Password, host)
}

// createDatabase create database if not exists.
func createDatabase(opts *genericoptions.MySQLOptions) error {
	createSql := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;",
		opts.Database,
	)
	db, err := sql.Open("mysql", mysqlEmptyDsn(opts))
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}

// migrate will migrate the tables to new db.
func migrate(client *gorm.DB) {
	models := model.GetRegisterModels()
	for _, k := range models {
		err := client.AutoMigrate(k)
		if err != nil {
			log.Errorf("migrate db for mysqldb storage err: %v", err)
			continue
		}
		log.Infof("migrating data for table %v", k.TableName())
	}
}
