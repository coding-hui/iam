// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// newUserRepository new User Repository.
func newUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db}
}

// Create creates a new user account.
func (u *userRepositoryImpl) Create(ctx context.Context, user *model.User, opts metav1.CreateOptions) (*model.User, error) {
	if oldUser, _ := u.GetByName(ctx, user.Name, metav1.GetOptions{}); oldUser != nil {
		return nil, errors.WithCode(code.ErrUserAlreadyExist, "User %s already exist", user.Name)
	}
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}
		return nil, err
	}
	if user.External != nil && user.External.ExternalUID != "" && user.External.IdentifyProvider != "" {
		user.External.UserId = user.InstanceID
		externalUser := &model.UserExternal{
			UserId:           user.InstanceID,
			ExternalUID:      user.External.ExternalUID,
			IdentifyProvider: user.External.IdentifyProvider,
		}
		err := u.deleteExternalUser(ctx, user.InstanceID, externalUser.ExternalUID, externalUser.IdentifyProvider)
		if err != nil {
			return nil, err
		}
		err = u.db.WithContext(ctx).Create(&externalUser).Error
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// Update updates an user account information.
func (u *userRepositoryImpl) Update(ctx context.Context, user *model.User, opts metav1.UpdateOptions) error {
	if err := u.db.WithContext(ctx).Save(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return err
	}

	return nil
}

// DeleteByInstanceId deletes the user by the user identifier.
func (u *userRepositoryImpl) DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}
	user, err := u.GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return err
	}
	if currentUser := ctx.Value(&v1.CtxKeyUserInstanceId); currentUser != "" && currentUser == user.InstanceID {
		return errors.WithCode(code.ErrDeleteOneself, "User %s failed to be deleted and cannot delete itself", currentUser)
	}
	err = u.db.WithContext(ctx).Where("instance_id = ?", instanceId).Select(clause.Associations).Delete(&model.User{}).Error
	if err != nil {
		return err
	}

	return nil
}

// BatchDelete batch deletes the users.
func (u *userRepositoryImpl) BatchDelete(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	return u.db.WithContext(ctx).Where("name in (?)", usernames).Delete(&model.User{}).Error
}

// GetByName get user by username.
func (u *userRepositoryImpl) GetByName(ctx context.Context, username string, _ metav1.GetOptions) (*model.User, error) {
	user := &model.User{}
	if username == "" {
		return nil, errors.WithCode(code.ErrUserNameIsEmpty, "Username is empty")
	}
	err := u.db.WithContext(ctx).Where("name = ?", username).First(&user).Error
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
	err := u.db.WithContext(ctx).Where("instance_id = ?", instanceId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, err
	}

	return user, nil
}

// GetByExternalId get user by external identifier.
func (u *userRepositoryImpl) GetByExternalId(
	ctx context.Context,
	externalUid, externalIdp string,
	opts metav1.GetOptions,
) (*model.User, error) {
	externalUser := &model.UserExternal{}
	err := u.db.WithContext(ctx).
		Where("external_uid = ?", externalUid).
		Where("idp = ?", externalIdp).
		First(&externalUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, err
	}
	user, err := u.GetByInstanceId(ctx, externalUser.UserId, opts)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, err
	}

	return user, nil
}

// List list users.
func (u *userRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	list := &v1.UserList{}
	err := u.db.WithContext(ctx).Model(model.User{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Find(&list.Items).Offset(-1).Limit(-1).
		Count(&list.TotalCount).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to list users")
	}

	return list, err
}

func (u *userRepositoryImpl) deleteExternalUser(ctx context.Context, userId, externalUid, idp string) error {
	externalUser := &model.UserExternal{
		UserId:           userId,
		ExternalUID:      externalUid,
		IdentifyProvider: idp,
	}
	return u.db.WithContext(ctx).
		Where("external_uid = ?", externalUid).
		Where("idp = ?", idp).
		Delete(&externalUser).Error
}
