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
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/fields"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

type roleRepositoryImpl struct {
	db *gorm.DB
}

// newRoleRepository new Role Repository.
func newRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepositoryImpl{db}
}

// Create creates a new role.
func (u *roleRepositoryImpl) Create(ctx context.Context, role *model.Role, _ metav1alpha1.CreateOptions) error {
	if oldRole, _ := u.GetByName(ctx, role.Name, metav1alpha1.GetOptions{}); oldRole != nil {
		return errors.WithCode(code.ErrRoleAlreadyExist, "Role %s already exist", role.Name)
	}
	if err := u.db.WithContext(ctx).Create(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.WithCode(code.ErrRoleAlreadyExist, "Role %s already exist", role.Name)
		}
		return err
	}

	return nil
}

// Update updates a role information.
func (u *roleRepositoryImpl) Update(ctx context.Context, role *model.Role, _ metav1alpha1.UpdateOptions) error {
	if err := u.db.WithContext(ctx).Save(role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrRoleNotFound, err.Error())
		}

		return err
	}

	return nil
}

// DeleteByInstanceId deletes the role by the role identifier.
func (u *roleRepositoryImpl) DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}
	err := u.db.WithContext(ctx).Where("instance_id = ?", instanceId).Delete(&model.Role{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrRoleNotFound, err.Error())
		}

		return err
	}

	return nil
}

// DeleteCollection batch deletes the roles.
func (u *roleRepositoryImpl) DeleteCollection(ctx context.Context, names []string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	return u.db.WithContext(ctx).Where("name in (?)", names).Delete(&model.Role{}).Error
}

// GetByName get role.
func (u *roleRepositoryImpl) GetByName(ctx context.Context, name string, _ metav1alpha1.GetOptions) (*model.Role, error) {
	role := &model.Role{}
	if name == "" {
		return nil, errors.WithCode(code.ErrRoleNameIsEmpty, "Role name is empty")
	}
	err := u.db.WithContext(ctx).Preload("Users").Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrRoleNotFound, err.Error())
		}

		return nil, err
	}

	return role, nil
}

// GetByInstanceID get role by instanceID.
func (u *roleRepositoryImpl) GetByInstanceID(ctx context.Context, instanceID string, _ metav1alpha1.GetOptions) (*model.Role, error) {
	role := &model.Role{}
	err := u.db.WithContext(ctx).Preload("Users").Where("instance_id = ?", instanceID).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrRoleNotFound, err.Error())
		}

		return nil, err
	}

	return role, nil
}

// List list roles.
func (u *roleRepositoryImpl) List(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.RoleList, error) {
	list := &v1alpha1.RoleList{}

	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	db := u.db.WithContext(ctx)
	var clauses []clause.Expression
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	clauses = _applyFieldSelector(clauses, selector)
	db.Model(&model.Role{}).
		Offset(ol.Offset).
		Limit(ol.Limit).
		Clauses(clauses...).
		Order("id desc").
		Find(&list.Items).
		Offset(-1).
		Limit(-1).
		Count(&list.TotalCount)

	return list, db.Error
}

// ListByUserInstanceId list roles by user instanceId.
func (u *roleRepositoryImpl) ListByUserInstanceId(
	ctx context.Context,
	userInstanceId string,
	opts metav1alpha1.ListOptions,
) (*v1alpha1.RoleList, error) {
	list := &v1alpha1.RoleList{}
	var roleIds []uint64

	ol := gormutil.Unpointer(opts.Offset, opts.Limit)
	db := u.db.WithContext(ctx)
	db.Raw("SELECT role_id FROM iam_user_role WHERE user_instance_id = ?", userInstanceId).Find(&roleIds)
	db.Model(&model.Role{}).
		Where("id in ?", roleIds).
		Offset(ol.Offset).
		Limit(ol.Limit).
		Find(&list.Items).
		Offset(-1).
		Limit(-1).
		Count(&list.TotalCount)

	return list, db.Error
}

// AssignUserRoles assign user roles.
func (u *roleRepositoryImpl) AssignUserRoles(ctx context.Context, role *model.Role, userInstanceIds []string) (int64, error) {
	err := u.db.WithContext(ctx).Model(&model.User{}).Where("instance_id in ?", userInstanceIds).Find(&role.Users).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrAssignRoleFailed, err.Error())
	}
	res := u.db.WithContext(ctx).Model(&role).Save(role)
	if res.Error != nil {
		return 0, errors.WithCode(code.ErrAssignRoleFailed, res.Error.Error())
	}

	return int64(len(role.Users)), nil
}

// RevokeUserRoles revoke user roles.
func (u *roleRepositoryImpl) RevokeUserRoles(ctx context.Context, role *model.Role, userInstanceIds []string) (int64, error) {
	err := u.db.WithContext(ctx).Model(&role).Where("instance_id in ?", userInstanceIds).Association("Users").Find(&role.Users)
	if err != nil {
		return 0, errors.WithCode(code.ErrRevokeRoleFailed, err.Error())
	}
	err = u.db.WithContext(ctx).Model(&role).Association("Users").Delete(&role.Users)
	if err != nil {
		return 0, errors.WithCode(code.ErrRevokeRoleFailed, err.Error())
	}

	return int64(len(userInstanceIds)), nil
}
