// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"context"

	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/utils/gormutil"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/fields"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// newUserRepository new User Repository
func newUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db}
}

// Create creates a new user account.
func (u *userRepositoryImpl) Create(ctx context.Context, user *model.User, opts metav1alpha1.CreateOptions) error {
	if oldUser, _ := u.Get(ctx, user.Name, metav1alpha1.GetOptions{}); oldUser != nil {
		return errors.WithCode(code.ErrUserAlreadyExist, "User %s already exist", user.Name)
	}
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}
		return err
	}

	return nil
}

// Update updates an user account information.
func (u *userRepositoryImpl) Update(ctx context.Context, user *model.User, opts metav1alpha1.UpdateOptions) error {
	if err := u.db.WithContext(ctx).Save(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return err
	}

	return nil
}

// Delete deletes the user by the user identifier.
func (u *userRepositoryImpl) Delete(ctx context.Context, username string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}
	if currentUser := ctx.Value(&v1alpha1.CtxKeyUserName); currentUser != "" && currentUser == username {
		return errors.WithCode(code.ErrDeleteOneself, "User %s failed to be deleted and cannot delete itself", currentUser)
	}
	err := u.db.WithContext(ctx).Where("name = ?", username).Delete(&model.User{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return err
	}

	return nil
}

// DeleteCollection batch deletes the users.
func (u *userRepositoryImpl) DeleteCollection(ctx context.Context, usernames []string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	return u.db.WithContext(ctx).Where("name in (?)", usernames).Delete(&model.User{}).Error
}

// Get get user
func (u *userRepositoryImpl) Get(ctx context.Context, username string, _ metav1alpha1.GetOptions) (*model.User, error) {
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

// List list users
func (u *userRepositoryImpl) List(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.UserList, error) {
	list := &v1alpha1.UserList{}

	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	db := u.db.WithContext(ctx)
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, _ := selector.RequiresExactMatch("name")
	if username != "" {
		db.Where("name like ?", "%"+username+"%")
	}
	db.Model(model.User{}).
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&list.Items).
		Offset(-1).
		Limit(-1).
		Count(&list.TotalCount)

	return list, db.Error
}
