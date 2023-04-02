package mysqldb

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s.io/klog/v2"

	"github.com/wecoding/iam/pkg/apiserver/infrastructure/datastore"
)

type mysqldb struct {
	client   *gorm.DB
	database string
}

// PrimaryKey primary key
const PrimaryKey = "_pk"

// New new mysqldb datastore instance
func New(ctx context.Context, cfg datastore.Config) (datastore.DataStore, error) {
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

	m := &mysqldb{
		client:   db,
		database: cfg.Database,
	}

	klog.Infof("create mysqldb datastore instance successful")

	return m, nil
}

// Add add data model
func (m *mysqldb) Add(ctx context.Context, entity datastore.Entity) error {
	return nil
}

// BatchAdd batch add entity, this operation has some atomicity.
func (m *mysqldb) BatchAdd(ctx context.Context, entities []datastore.Entity) error {
	notRollback := make(map[string]int)
	for i, saveEntity := range entities {
		if err := m.Add(ctx, saveEntity); err != nil {
			if errors.Is(err, datastore.ErrRecordExist) {
				notRollback[saveEntity.PrimaryKey()] = 1
			}
			for _, deleteEntity := range entities[:i] {
				if _, exit := notRollback[deleteEntity.PrimaryKey()]; !exit {
					if err := m.Delete(ctx, deleteEntity); err != nil {
						if !errors.Is(err, datastore.ErrRecordNotExist) {
							klog.Errorf("rollback delete entity failure %w", err)
						}
					}
				}
			}
			return datastore.NewDBError(fmt.Errorf("save entities occur error, %w", err))
		}
	}
	return nil
}

// Get get data model
func (m *mysqldb) Get(ctx context.Context, entity datastore.Entity) error {
	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}
	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}
	res := m.client.Where(&entity).First(&entity)
	if res.Error != nil {
		return datastore.NewDBError(res.Error)
	}
	return nil
}

// Update update data model
func (m *mysqldb) Update(ctx context.Context, entity datastore.Entity) error {
	return nil
}

// IsExist determine whether data exists.
func (m *mysqldb) IsExist(ctx context.Context, entity datastore.Entity) (bool, error) {
	return true, nil
}

// Delete delete data
func (m *mysqldb) Delete(ctx context.Context, entity datastore.Entity) error {
	return nil
}

// List list entity function
func (m *mysqldb) List(ctx context.Context, entity datastore.Entity, op *datastore.ListOptions) ([]datastore.Entity, error) {
	return nil, nil
}

// Count counts entities
func (m *mysqldb) Count(ctx context.Context, entity datastore.Entity, filterOptions *datastore.FilterOptions) (int64, error) {
	return 0, nil
}
