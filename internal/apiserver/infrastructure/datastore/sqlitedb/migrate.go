// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/pkg/log"
)

// migrate will migrate the tables to new db.
func migrate(client *gorm.DB) {
	models := model.GetRegisterModels()
	for _, k := range models {
		err := client.AutoMigrate(k)
		if err != nil {
			log.Errorf("migrate db for sqlite storage err: %v", err)
			continue
		}
		log.Infof("migrating data for table %v", k.TableName())
	}
}
