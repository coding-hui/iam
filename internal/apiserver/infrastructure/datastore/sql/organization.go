// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore"
	"github.com/coding-hui/iam/pkg/code"
)

type orgRepositoryImpl struct {
	client *Client
}

// newOrganizationRepository new Organization Repository.
func newOrganizationRepository(client *Client) repository.OrganizationRepository {
	return &orgRepositoryImpl{client}
}

func (o *orgRepositoryImpl) Create(ctx context.Context, org *model.Organization, opts metav1.CreateOptions) error {
	if old, _ := o.GetByName(ctx, org.Name, metav1.GetOptions{}); old != nil {
		return errors.WithCode(code.ErrOrgAlreadyExist, "Organization %s already exist", org.Name)
	}
	if err := o.client.WithCtx(ctx).Create(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.WithCode(code.ErrOrgAlreadyExist, "%s", err.Error())
		}
		return err
	}

	return nil
}

func (o *orgRepositoryImpl) Update(ctx context.Context, org *model.Organization, opts metav1.UpdateOptions) error {
	err := o.client.WithCtx(ctx).Save(org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrOrgNotFound, "%s", err.Error())
		}

		return err
	}

	return nil
}

func (o *orgRepositoryImpl) BatchUpdate(ctx context.Context, list []*model.Organization, opts metav1.UpdateOptions) error {
	db := o.client.WithCtx(ctx).Model(model.Organization{})
	if len(opts.DryRun) > 0 {
		db.DryRun = true
	}
	needUpdatesColumns := []string{
		"ancestors",
		"parent_id",
		"display_name",
		"website_url",
		"favicon",
		"disabled",
		"description",
	}
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "instance_id"}},
		DoUpdates: clause.AssignmentColumns(needUpdatesColumns),
	}).CreateInBatches(&list, 500).Error
	if err != nil {
		return datastore.NewDBError(err, "failed to batch update organizations.")
	}

	return nil
}

func (o *orgRepositoryImpl) DeleteByName(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	db := o.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	err := db.Where("name = ?", name).Delete(&model.Organization{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrOrgNotFound, "%s", err.Error())
		}

		return err
	}

	return nil
}

func (o *orgRepositoryImpl) DeleteByInstanceId(ctx context.Context, uid string, opts metav1.DeleteOptions) error {
	db := o.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	err := db.Where("instance_id = ?", uid).Delete(&model.Organization{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrOrgNotFound, "%s", err.Error())
		}

		return err
	}

	return nil
}

func (o *orgRepositoryImpl) GetByName(ctx context.Context, name string, opts metav1.GetOptions) (*model.Organization, error) {
	var org *model.Organization
	err := o.client.WithCtx(ctx).Where("name = ?", name).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrOrgNotFound, "%s", err.Error())
		}

		return nil, err
	}

	return org, err
}

func (o *orgRepositoryImpl) GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Organization, error) {
	var org *model.Organization
	err := o.client.WithCtx(ctx).Where("instance_id = ?", instanceId).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrOrgNotFound, "%s", err.Error())
		}

		return nil, err
	}

	return org, err
}

func (o *orgRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) ([]model.Organization, error) {
	var list []model.Organization
	err := o.client.WithCtx(ctx).Model(model.Organization{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Find(&list).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to list organizations")
	}

	return list, err
}

func (o *orgRepositoryImpl) Count(ctx context.Context, opts metav1.ListOptions) (int64, error) {
	var totalCount int64
	err := o.client.WithCtx(ctx).Model(&model.Organization{}).
		Scopes(
			makeCondition(opts),
		).
		Count(&totalCount).Error
	if err != nil {
		return 0, datastore.NewDBError(err, "failed to get organization total")
	}

	return totalCount, nil
}

func (o *orgRepositoryImpl) CountDepartmentByParent(ctx context.Context, parent string, opts metav1.ListOptions) (int64, error) {
	var totalCount int64
	err := o.client.WithCtx(ctx).Model(&model.Organization{}).
		Scopes(
			makeCondition(opts),
		).
		Where("FIND_IN_SET(?, ancestors)", parent).
		Count(&totalCount).Error
	if err != nil {
		return 0, datastore.NewDBError(err, "failed to get department total")
	}

	return totalCount, nil
}

func (o *orgRepositoryImpl) AddDepartmentMembers(ctx context.Context, members []*model.DepartmentMember) error {
	db := o.client.WithCtx(ctx).Model(&model.DepartmentMember{})
	err := db.CreateInBatches(members, 500).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") {
			return errors.WithCode(code.ErrMemberAlreadyInDepartment, "%s", err.Error())
		}
		return datastore.NewDBError(err, "failed to add department members")
	}

	return err
}

func (o *orgRepositoryImpl) RemoveDepartmentMembers(ctx context.Context, members []*model.DepartmentMember) error {
	db := o.client.WithCtx(ctx).Model(&model.DepartmentMember{})
	err := db.Delete(members).Error
	if err != nil {
		return datastore.NewDBError(err, "failed to remove department members")
	}

	return err
}

func (o *orgRepositoryImpl) ListChildDepartments(
	ctx context.Context,
	org string,
	opts metav1.ListOptions,
) ([]model.Organization, error) {
	var child []model.Organization
	err := o.client.WithCtx(ctx).Model(&model.Organization{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Where("FIND_IN_SET(?, ancestors)", org).
		Find(&child).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to list 【org】 child departments")
	}

	return child, nil
}

func (o *orgRepositoryImpl) ListDepartmentMembers(
	ctx context.Context,
	department string,
	opts metav1.ListOptions,
) ([]model.DepartmentMember, error) {
	var members []model.DepartmentMember
	err := o.client.WithCtx(ctx).Model(&model.DepartmentMember{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Where("department_id = ?", department).
		Find(&members).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to list department members")
	}

	return members, nil
}

func (o *orgRepositoryImpl) CountDepartmentMembers(ctx context.Context, department string, opts metav1.ListOptions) (int64, error) {
	var totalCount int64
	err := o.client.WithCtx(ctx).Model(&model.DepartmentMember{}).
		Scopes(
			makeCondition(opts),
		).
		Where("department_id = ?", department).
		Count(&totalCount).Error
	if err != nil {
		return 0, datastore.NewDBError(err, "failed to get department total members")
	}

	return totalCount, nil
}

func (o *orgRepositoryImpl) UpdateIsLeafState(ctx context.Context, orgOrDept string, isLeaf bool) error {
	err := o.client.WithCtx(ctx).Model(&model.Organization{}).
		Where("instance_id = ?", orgOrDept).
		Update("is_leaf", isLeaf).
		Error
	if err != nil {
		return err
	}

	return nil
}
