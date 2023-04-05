package mysqldb

import (
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"k8s.io/klog/v2"

	iamv1alpha1 "github.com/coding-hui/api/iam/v1alpha1"
	"github.com/wecoding/iam/pkg/apiserver/domain/model"
	"github.com/wecoding/iam/pkg/apiserver/infrastructure/datastore"
)

// mysqlEmptyDsn msyql empty dsn for create databases
func mysqlEmptyDsn(cfg datastore.Config) string {
	return cfg.URL[:strings.Index(cfg.URL, "/")+1]
}

func createDatabase(cfg datastore.Config) error {
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", cfg.Database)
	db, err := sql.Open("mysql", mysqlEmptyDsn(cfg))
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}

// migrate will migrate the tables to new db
func migrate(client *gorm.DB) {
	models := model.GetRegisterModels()
	for _, k := range models {
		err := client.AutoMigrate(k)
		if err != nil {
			klog.Errorf("migrate db for mysqldb storage err: %v", err)
			continue
		}
		klog.Infof("migrating data for table %v", k.TableName())
	}
	client.AutoMigrate(iamv1alpha1.User{})
}
