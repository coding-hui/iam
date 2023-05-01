// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"context"

	"gorm.io/gorm"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/fields"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	convert "github.com/coding-hui/iam/internal/apiserver/interfaces/api/convert/v1alpha1"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/utils/gormutil"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
)

type resourceRepositoryImpl struct {
	db *gorm.DB
}

// newResourceRepository new User Repository
func newResourceRepository(db *gorm.DB) repository.ResourceRepository {
	return &resourceRepositoryImpl{db}
}

// BatchCreate creates a new resource.
func (r *resourceRepositoryImpl) BatchCreate(ctx context.Context, resources []*model.Resource, opts metav1alpha1.CreateOptions) error {
	if err := r.db.WithContext(ctx).CreateInBatches(&resources, 500).Error; err != nil {
		return err
	}

	return nil
}

// Create creates a new resource.
func (r *resourceRepositoryImpl) Create(ctx context.Context, resource *model.Resource, opts metav1alpha1.CreateOptions) error {
	if err := r.db.WithContext(ctx).Create(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.WithCode(code.ErrResourceAlreadyExist, err.Error())
		}
		return err
	}

	return nil
}

// Update updates an resource information.
func (r *resourceRepositoryImpl) Update(ctx context.Context, resource *model.Resource, opts metav1alpha1.UpdateOptions) error {
	db := r.db.WithContext(ctx)
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("resource_id = ?", resource.ID).Delete(&model.Action{}).Error
		if err != nil {
			return errors.WithCode(code.ErrDatabase, "failed to delete resource actions")
		}
		err = tx.Save(resource).Error
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the resource by the resource identifier.
func (r *resourceRepositoryImpl) Delete(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		r.db = r.db.Unscoped()
	}
	resource, err := r.Get(ctx, name, metav1alpha1.GetOptions{})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrResourceNotFound, err.Error())
	}
	err = r.db.WithContext(ctx).
		Select("Actions").
		Where("id = ?", resource.ID).
		Delete(&model.Resource{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrResourceNotFound, err.Error())
		}

		return err
	}

	return nil
}

// DeleteCollection batch deletes the resource.
func (r *resourceRepositoryImpl) DeleteCollection(ctx context.Context, names []string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		r.db = r.db.Unscoped()
	}

	return r.db.WithContext(ctx).
		Select("Actions").
		Where("name in (?)", names).
		Delete(&model.Resource{}).
		Error
}

// Get get resource
func (r *resourceRepositoryImpl) Get(ctx context.Context, name string, _ metav1alpha1.GetOptions) (*model.Resource, error) {
	resource := &model.Resource{}
	if name == "" {
		return nil, errors.WithCode(code.ErrResourceNameIsEmpty, "Resource name is empty")
	}
	err := r.db.WithContext(ctx).Preload("Actions").Where("name = ?", name).First(&resource).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrResourceNotFound, err.Error())
		}

		return nil, err
	}

	return resource, nil
}

// List list resources
func (r *resourceRepositoryImpl) List(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.ResourceList, error) {
	resources := &[]model.Resource{}
	res := &v1alpha1.ResourceList{}

	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	db := r.db.WithContext(ctx)
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, _ := selector.RequiresExactMatch("name")
	if username != "" {
		db.Where("name like ?", "%"+username+"%")
	}
	err := db.Model(&model.Resource{}).
		Preload("Actions").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(resources).
		Offset(-1).
		Limit(-1).
		Count(&res.TotalCount).
		Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, "failed to list resources")
	}

	for _, resource := range *resources {
		res.Items = append(res.Items, convert.ConvertResourceModelToBase(&resource))
	}

	return res, nil
}
