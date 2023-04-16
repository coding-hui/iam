// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"time"

	"k8s.io/klog/v2"

	iamv1alpha1 "github.com/coding-hui/api/iam/v1alpha1"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/code"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
	"github.com/coding-hui/common/util/auth"
)

const (
	// DefaultAdmin default admin username
	DefaultAdmin string = "ADMIN"
	// DefaultAdminPwd default admin password
	DefaultAdminPwd string = "WECODING"
	// DefaultAdminUserAlias default admin user alias
	DefaultAdminUserAlias string = "Administrator"
)

// UserService User manage api
type UserService interface {
	Create(ctx context.Context, user *iamv1alpha1.User, opts metav1alpha1.CreateOptions) error
	Update(ctx context.Context, user *iamv1alpha1.User, opts metav1alpha1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metav1alpha1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1alpha1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metav1alpha1.GetOptions) (*iamv1alpha1.User, error)
	List(ctx context.Context, opts metav1alpha1.ListOptions) (*iamv1alpha1.UserList, error)
	FlushLastLoginTime(ctx context.Context, user *iamv1alpha1.User) error
	Init(ctx context.Context) error
}

type userServiceImpl struct {
	Store repository.Factory `inject:"repository"`
}

// NewUserService new User service
func NewUserService() UserService {
	return &userServiceImpl{}
}

// Init initialize user data
func (u *userServiceImpl) Init(ctx context.Context) error {
	_, err := u.Get(ctx, DefaultAdmin, metav1alpha1.GetOptions{})
	if err != nil && errors.IsCode(err, code.ErrUserNotFound) {
		user := &iamv1alpha1.User{
			ObjectMeta: metav1alpha1.ObjectMeta{
				Name: DefaultAdmin,
			},
			Password: DefaultAdminPwd,
			Alias:    DefaultAdminUserAlias,
			Disabled: false,
		}
		err = u.Create(ctx, user, metav1alpha1.CreateOptions{})
		if err != nil {
			return errors.WithMessagef(err, "Failed to initialize default admin")
		}
		klog.Info("initialize default admin done")
	}

	return nil
}

// Create create user
func (u *userServiceImpl) Create(ctx context.Context, user *iamv1alpha1.User, opts metav1alpha1.CreateOptions) error {
	user.Password, _ = auth.Encrypt(user.Password)
	if err := u.Store.UserRepository().Create(ctx, user, opts); err != nil {
		return err
	}

	return nil
}

// Update update user
func (u *userServiceImpl) Update(ctx context.Context, user *iamv1alpha1.User, opts metav1alpha1.UpdateOptions) error {
	if err := u.Store.UserRepository().Update(ctx, user, opts); err != nil {
		return err
	}

	return nil
}

// Delete delete user
func (u *userServiceImpl) Delete(ctx context.Context, username string, opts metav1alpha1.DeleteOptions) error {
	if err := u.Store.UserRepository().Delete(ctx, username, opts); err != nil {
		return err
	}

	return nil
}

// DeleteCollection batch delete user
func (u *userServiceImpl) DeleteCollection(ctx context.Context, usernames []string, opts metav1alpha1.DeleteOptions) error {
	if err := u.Store.UserRepository().DeleteCollection(ctx, usernames, opts); err != nil {
		return err
	}

	return nil
}

// Get get user
func (u *userServiceImpl) Get(ctx context.Context, username string, opts metav1alpha1.GetOptions) (*iamv1alpha1.User, error) {
	user, err := u.Store.UserRepository().Get(ctx, username, opts)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// List list users
func (u *userServiceImpl) List(ctx context.Context, listOptions metav1alpha1.ListOptions) (*iamv1alpha1.UserList, error) {
	users, err := u.Store.UserRepository().List(ctx, metav1alpha1.ListOptions{
		Offset: listOptions.Offset,
		Limit:  listOptions.Limit,
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

// FlushLastLoginTime update user login time
func (u *userServiceImpl) FlushLastLoginTime(ctx context.Context, user *iamv1alpha1.User) error {
	now := time.Now()
	user.LastLoginTime = &now

	return u.Update(ctx, user, metav1alpha1.UpdateOptions{})
}
