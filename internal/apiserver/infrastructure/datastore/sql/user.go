// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore"
	"github.com/coding-hui/iam/internal/pkg/request"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

type userRepositoryImpl struct {
	client *Client
}

// newUserRepository new User Repository.
func newUserRepository(client *Client) repository.UserRepository {
	return &userRepositoryImpl{client}
}

// Create creates a new user account.
func (u *userRepositoryImpl) Create(ctx context.Context, user *model.User, opts metav1.CreateOptions) (*model.User, error) {
	if oldUser, _ := u.GetByName(ctx, user.Name, metav1.GetOptions{}); oldUser != nil {
		return nil, errors.WithCode(code.ErrUserAlreadyExist, "User %s already exist", user.Name)
	}
	if err := u.client.WithCtx(ctx).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}
		return nil, err
	}
	if user.External != nil && user.External.ExternalUID != "" && user.External.IdentifyProvider != "" {
		user.External.UserID = user.InstanceID
		externalUser := &model.UserExternal{
			UserID:           user.InstanceID,
			ExternalUID:      user.External.ExternalUID,
			IdentifyProvider: user.External.IdentifyProvider,
		}
		err := u.deleteExternalUser(ctx, user.InstanceID, externalUser.ExternalUID, externalUser.IdentifyProvider)
		if err != nil {
			return nil, err
		}
		err = u.client.WithCtx(ctx).Create(&externalUser).Error
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// Update updates an user account information.
func (u *userRepositoryImpl) Update(ctx context.Context, user *model.User, opts metav1.UpdateOptions) error {
	if err := u.client.WithCtx(ctx).Save(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return err
	}

	return nil
}

// DeleteByInstanceId deletes the user by the user identifier.
func (u *userRepositoryImpl) DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	db := u.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	user, err := u.GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return err
	}
	currentUser, ok := request.UserFrom(ctx)
	if ok && currentUser.InstanceID == user.InstanceID {
		return errors.WithCode(code.ErrDeleteOneself, "User %s failed to be deleted and cannot delete itself", currentUser.Name)
	}
	err = db.Where("instance_id = ?", instanceId).Select(clause.Associations).Delete(&model.User{}).Error
	if err != nil {
		return err
	}

	return nil
}

// BatchDelete batch deletes the users.
func (u *userRepositoryImpl) BatchDelete(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	db := u.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}

	return db.Where("name in (?)", usernames).Delete(&model.User{}).Error
}

// GetByName get user by username.
func (u *userRepositoryImpl) GetByName(ctx context.Context, username string, _ metav1.GetOptions) (*model.User, error) {
	user := &model.User{}
	if username == "" {
		return nil, errors.WithCode(code.ErrUserNameIsEmpty, "Username is empty")
	}
	err := u.client.WithCtx(ctx).Where("name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, err
	}

	return user, nil
}

// GetByInstanceId get user by instanceId.
func (u *userRepositoryImpl) GetByInstanceId(ctx context.Context, instanceId string, _ metav1.GetOptions) (*model.User, error) {
	user := &model.User{}
	err := u.client.WithCtx(ctx).Where("instance_id = ?", instanceId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, err
	}
	departmentIds, _ := u.getUserDepartments(ctx, instanceId)
	for _, v := range departmentIds {
		user.DepartmentIds = append(user.DepartmentIds, v.DepartmentID)
	}

	return user, nil
}

// GetByNameOrInstanceId get user by name or instanceId.
func (u *userRepositoryImpl) GetByNameOrInstanceId(ctx context.Context, nameOrId string, _ metav1.GetOptions) (*model.User, error) {
	user := &model.User{}
	err := u.client.WithCtx(ctx).Where("instance_id = ? or name = ?", nameOrId, nameOrId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, err
	}
	departmentIds, _ := u.getUserDepartments(ctx, nameOrId)
	for _, v := range departmentIds {
		user.DepartmentIds = append(user.DepartmentIds, v.DepartmentID)
	}

	return user, nil
}

func (u *userRepositoryImpl) FlushLastLoginTime(ctx context.Context, nameOrId string) error {
	err := u.client.WithCtx(ctx).Debug().
		Model(&model.User{}).
		Where("instance_id = ? or name = ?", nameOrId, nameOrId).
		Update("last_login_time", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

// GetByExternalId get user by external identifier.
func (u *userRepositoryImpl) GetByExternalId(
	ctx context.Context,
	externalUid, externalIdp string,
	opts metav1.GetOptions,
) (*model.User, error) {
	externalUser := &model.UserExternal{}
	err := u.client.WithCtx(ctx).
		Where("external_uid = ?", externalUid).
		Where("idp = ?", externalIdp).
		First(&externalUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, err
	}
	user, err := u.GetByInstanceId(ctx, externalUser.UserID, opts)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, err
	}

	return user, nil
}

// List list users.
func (u *userRepositoryImpl) List(ctx context.Context, opts v1.ListUserOptions) ([]model.User, error) {
	var list []model.User
	db := u.client.WithCtx(ctx).Model(model.User{}).
		Scopes(
			makeCondition(opts.ListOptions),
			paginate(opts.ListOptions),
			joinDepartments(opts),
		).
		Order("id desc")
	err := db.Find(&list).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to list users")
	}

	return list, err
}

// Count count users.
func (u *userRepositoryImpl) Count(ctx context.Context, opts v1.ListUserOptions) (int64, error) {
	var totalCount int64
	err := u.client.WithCtx(ctx).Model(&model.User{}).
		Scopes(
			makeCondition(opts.ListOptions),
			joinDepartments(opts),
		).
		Count(&totalCount).Error
	if err != nil {
		return 0, datastore.NewDBError(err, "failed to get user total")
	}

	return totalCount, nil
}

func (u *userRepositoryImpl) deleteExternalUser(ctx context.Context, userId, externalUid, idp string) error {
	externalUser := &model.UserExternal{
		UserID:           userId,
		ExternalUID:      externalUid,
		IdentifyProvider: idp,
	}
	return u.client.WithCtx(ctx).
		Where("external_uid = ?", externalUid).
		Where("idp = ?", idp).
		Delete(&externalUser).Error
}

func (u *userRepositoryImpl) getUserDepartments(ctx context.Context, user string) ([]*model.DepartmentMember, error) {
	var resp []*model.DepartmentMember
	err := u.client.WithCtx(ctx).
		Model(&model.DepartmentMember{}).
		Where("member_id = ?", user).
		Find(&resp).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to get user departments.")
	}
	return resp, err
}

func joinDepartments(opts v1.ListUserOptions) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if opts.DepartmentID != "" {
			db.Joins("INNER JOIN iam_department_member dm ON dm.member_id = iam_user.instance_id").
				Joins("INNER JOIN iam_organization o ON o.instance_id = dm.department_id").
				Where("dm.department_id = ?", opts.DepartmentID)
			if opts.IncludeChildrenDepartments {
				db.Or("FIND_IN_SET(?, o.ancestors)", opts.DepartmentID)
			}
		}
		return db
	}
}
