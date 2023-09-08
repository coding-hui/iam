// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/utils/gormutil"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/fields"
	metav1 "github.com/coding-hui/common/meta/v1"
)

type orgRepositoryImpl struct {
	db *gorm.DB
}

// newOrganizationRepository new Organization Repository.
func newOrganizationRepository(db *gorm.DB) repository.OrganizationRepository {
	return &orgRepositoryImpl{db}
}

func (o *orgRepositoryImpl) Create(ctx context.Context, org *model.Organization, opts metav1.CreateOptions) error {
	if old, _ := o.GetByName(ctx, org.Name, metav1.GetOptions{}); old != nil {
		return errors.WithCode(code.ErrOrgAlreadyExist, "Organization %s already exist", org.Name)
	}
	if err := o.db.WithContext(ctx).Create(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.WithCode(code.ErrOrgAlreadyExist, err.Error())
		}
		return err
	}

	return nil
}

func (o *orgRepositoryImpl) Update(ctx context.Context, org *model.Organization, opts metav1.UpdateOptions) error {
	err := o.db.WithContext(ctx).Save(org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrOrgNotFound, err.Error())
		}

		return err
	}

	return nil
}

func (o *orgRepositoryImpl) DeleteByName(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	if opts.Unscoped {
		o.db = o.db.Unscoped()
	}
	err := o.db.WithContext(ctx).Where("name = ?", name).Delete(&model.Organization{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrOrgNotFound, err.Error())
		}

		return err
	}

	return nil
}

func (o *orgRepositoryImpl) DeleteByInstanceId(ctx context.Context, uid string, opts metav1.DeleteOptions) error {
	if opts.Unscoped {
		o.db = o.db.Unscoped()
	}
	err := o.db.WithContext(ctx).Where("instance_id = ?", uid).Delete(&model.Organization{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrOrgNotFound, err.Error())
		}

		return err
	}

	return nil
}

func (o *orgRepositoryImpl) GetByName(ctx context.Context, name string, opts metav1.GetOptions) (*model.Organization, error) {
	var org *model.Organization
	err := o.db.WithContext(ctx).Where("name = ?", name).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrOrgNotFound, err.Error())
		}

		return nil, err
	}

	return org, err
}

func (o *orgRepositoryImpl) GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Organization, error) {
	var org *model.Organization
	err := o.db.WithContext(ctx).Where("instance_id = ?", instanceId).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrOrgNotFound, err.Error())
		}

		return nil, err
	}

	return org, err
}

func (o *orgRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) (*v1.OrganizationList, error) {
	list := &v1.OrganizationList{}

	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	db := o.db.WithContext(ctx).Model(model.Organization{})
	var clauses []clause.Expression
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	clauses = _applyFieldSelector(clauses, selector)
	db.Offset(ol.Offset).
		Limit(ol.Limit).
		Clauses(clauses...).
		Order("id desc").
		Find(&list.Items).
		Offset(-1).
		Limit(-1).
		Count(&list.TotalCount)

	return list, db.Error
}
