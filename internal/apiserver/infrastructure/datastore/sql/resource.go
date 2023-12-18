// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"

	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

type resourceRepositoryImpl struct {
	client *Client
}

// newResourceRepository new User Repository.
func newResourceRepository(client *Client) repository.ResourceRepository {
	return &resourceRepositoryImpl{client}
}

// BatchCreate creates a new resource.
func (r *resourceRepositoryImpl) BatchCreate(ctx context.Context, resources []*model.Resource, opts metav1.CreateOptions) error {
	if err := r.client.WithCtx(ctx).CreateInBatches(&resources, 500).Error; err != nil {
		return err
	}

	return nil
}

// Create creates a new resource.
func (r *resourceRepositoryImpl) Create(ctx context.Context, resource *model.Resource, opts metav1.CreateOptions) error {
	if err := r.client.WithCtx(ctx).Create(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.WithCode(code.ErrResourceAlreadyExist, err.Error())
		}
		return err
	}

	return nil
}

// Update updates an resource information.
func (r *resourceRepositoryImpl) Update(ctx context.Context, resource *model.Resource, opts metav1.UpdateOptions) error {
	db := r.client.WithCtx(ctx)
	err := db.Where("resource_id = ?", resource.ID).Delete(&model.Action{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, "failed to delete resource actions")
	}
	err = db.Save(resource).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteByInstanceId deletes the resource by the resource identifier.
func (r *resourceRepositoryImpl) DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	db := r.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	resource, err := r.GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrResourceNotFound, err.Error())
	}
	err = db.Select("Actions").
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
func (r *resourceRepositoryImpl) DeleteCollection(ctx context.Context, names []string, opts metav1.DeleteOptions) error {
	db := r.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}

	return db.Select("Actions").
		Where("name in (?)", names).
		Delete(&model.Resource{}).
		Error
}

// GetByName get resource.
func (r *resourceRepositoryImpl) GetByName(ctx context.Context, name string, _ metav1.GetOptions) (*model.Resource, error) {
	resource := &model.Resource{}
	if name == "" {
		return nil, errors.WithCode(code.ErrResourceNameIsEmpty, "Resource name is empty")
	}
	err := r.client.WithCtx(ctx).Preload("Actions").Where("name = ?", name).First(&resource).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrResourceNotFound, err.Error())
		}

		return nil, err
	}

	return resource, nil
}

// GetByInstanceId get resource.
func (r *resourceRepositoryImpl) GetByInstanceId(
	ctx context.Context,
	instanceId string,
	_ metav1.GetOptions,
) (*model.Resource, error) {
	resource := &model.Resource{}
	if instanceId == "" {
		return nil, errors.WithCode(code.ErrResourceNameIsEmpty, "Resource instanceId is empty")
	}
	err := r.client.WithCtx(ctx).Preload("Actions").Where("instance_id = ?", instanceId).First(&resource).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrResourceNotFound, err.Error())
		}

		return nil, err
	}

	return resource, nil
}

// List list resources.
func (r *resourceRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) (*v1.ResourceList, error) {
	resources := &[]model.Resource{}
	res := &v1.ResourceList{}
	err := r.client.WithCtx(ctx).Model(model.Resource{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Preload("Actions").
		Order("id desc").
		Find(resources).Offset(-1).Limit(-1).
		Count(&res.TotalCount).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to list resources")
	}
	for _, resource := range *resources {
		res.Items = append(res.Items, assembler.ConvertResourceModelToBase(&resource))
	}

	return res, nil
}
