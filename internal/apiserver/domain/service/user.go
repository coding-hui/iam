// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"time"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/auth"
)

const (
	// DefaultAdmin default admin username.
	DefaultAdmin string = "ADMIN"
	// DefaultAdminPwd default admin password.
	DefaultAdminPwd string = "WECODING"
	// DefaultAdminUserAlias default admin user alias.
	DefaultAdminUserAlias string = "Administrator"
)

// UserService User manage api.
type UserService interface {
	CreateUser(ctx context.Context, req v1.CreateUserRequest) (*v1.CreateUserResponse, error)
	UpdateUser(ctx context.Context, instanceId string, req v1.UpdateUserRequest) (*v1.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error
	BatchDeleteUsers(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error
	GetUser(ctx context.Context, username string, opts metav1.GetOptions) (*model.User, error)
	GetUserByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.User, error)
	DetailUser(ctx context.Context, user *model.User) (*v1.DetailUserResponse, error)
	ListUsers(ctx context.Context, opts v1.ListUserOptions) (*v1.UserList, error)
	ListUserRoles(ctx context.Context, instanceId string, opts metav1.ListOptions) (*v1.RoleList, error)
	FlushLastLoginTime(ctx context.Context, user *model.User) error
	DisableUser(ctx context.Context, instanceId string) error
	EnableUser(ctx context.Context, instanceId string) error
	Init(ctx context.Context) error
}

type userServiceImpl struct {
	Store       repository.Factory `inject:"repository"`
	RoleService RoleService        `inject:""`
}

// NewUserService new User service.
func NewUserService() UserService {
	return &userServiceImpl{}
}

// Init initialize user data.
func (u *userServiceImpl) Init(ctx context.Context) error {
	_, err := u.GetUser(ctx, DefaultAdmin, metav1.GetOptions{})
	if err != nil && errors.IsCode(err, code.ErrUserNotFound) {
		user := v1.CreateUserRequest{
			Name:     DefaultAdmin,
			Password: DefaultAdminPwd,
			Alias:    DefaultAdminUserAlias,
			UserType: v1.PlatformAdmin.String(),
		}
		_, err = u.CreateUser(ctx, user)
		if err != nil {
			return errors.WithMessagef(err, "Failed to initialize default admin")
		}
		log.Info("initialize default admin done")
	}

	return nil
}

// CreateUser create user.
func (u *userServiceImpl) CreateUser(ctx context.Context, req v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	var external *model.UserExternal
	if req.ExternalUID != "" && req.IdentifyProvider != "" {
		external = &model.UserExternal{
			ExternalUID:      req.ExternalUID,
			IdentifyProvider: req.IdentifyProvider,
		}
	}
	// create user
	user := assembler.ConvertCreateUserReqToUserModel(req, external)
	createUser, err := u.Store.UserRepository().Create(ctx, user, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	// user org association
	var deptMembers []*model.DepartmentMember
	for _, dept := range req.DepartmentIds {
		deptMembers = append(deptMembers, &model.DepartmentMember{
			DepartmentID: dept,
			MemberID:     user.GetInstanceID(),
		})
	}
	err = u.Store.OrganizationRepository().AddDepartmentMembers(ctx, deptMembers)
	if err != nil {
		log.Errorf("Failed to add user to the department: %w", err)
	}

	return &v1.CreateUserResponse{
		UserBase: *assembler.ConvertUserModelToBase(createUser),
	}, nil
}

// UpdateUser update user.
func (u *userServiceImpl) UpdateUser(ctx context.Context, instanceId string, req v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	user, err := u.GetUserByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if req.Alias != "" {
		user.Alias = req.Alias
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Password != "" {
		hash, err := auth.Encrypt(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hash
	}
	if err := u.Store.UserRepository().Update(ctx, user, metav1.UpdateOptions{}); err != nil {
		return nil, err
	}
	base := assembler.ConvertUserModelToBase(user)

	return &v1.UpdateUserResponse{
		UserBase: *base,
	}, nil
}

// DeleteUser delete user.
func (u *userServiceImpl) DeleteUser(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	roles, err := u.ListUserRoles(ctx, instanceId, metav1.ListOptions{})
	if err != nil {
		return err
	}
	batchRevokeRoleReq := v1.BatchRevokeRoleRequest{
		Targets: []string{instanceId},
	}
	for _, r := range roles.Items {
		batchRevokeRoleReq.InstanceIds = append(batchRevokeRoleReq.InstanceIds, r.InstanceID)
	}
	if err := u.RoleService.BatchRevokeRole(ctx, batchRevokeRoleReq); err != nil {
		log.Errorf("failed to delete user [%s] roles: %s", instanceId, err.Error())
		return err
	}
	if err := u.Store.UserRepository().DeleteByInstanceId(ctx, instanceId, opts); err != nil {
		return err
	}

	return nil
}

// BatchDeleteUsers batch delete user.
func (u *userServiceImpl) BatchDeleteUsers(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	if err := u.Store.UserRepository().BatchDelete(ctx, usernames, opts); err != nil {
		return err
	}

	return nil
}

// GetUser get user.
func (u *userServiceImpl) GetUser(ctx context.Context, username string, opts metav1.GetOptions) (*model.User, error) {
	user, err := u.Store.UserRepository().GetByName(ctx, username, opts)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByInstanceId get user by instanceId.
func (u *userServiceImpl) GetUserByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.User, error) {
	user, err := u.Store.UserRepository().GetByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DetailUser return user detail
func (u *userServiceImpl) DetailUser(ctx context.Context, user *model.User) (*v1.DetailUserResponse, error) {
	base := *assembler.ConvertUserModelToBase(user)

	return &v1.DetailUserResponse{
		UserBase: base,
	}, nil
}

// ListUsers list users.
func (u *userServiceImpl) ListUsers(ctx context.Context, opts v1.ListUserOptions) (*v1.UserList, error) {
	userRepo := u.Store.UserRepository()
	var userList []*v1.DetailUserResponse
	users, err := userRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range users {
		userList = append(userList, convertUserModel(&v))
	}
	count, err := userRepo.Count(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &v1.UserList{
		Items: userList,
		ListMeta: metav1.ListMeta{
			TotalCount: count,
		},
	}, nil
}

// ListUserRoles list users.
func (u *userServiceImpl) ListUserRoles(
	ctx context.Context,
	instanceId string,
	listOptions metav1.ListOptions,
) (*v1.RoleList, error) {
	roles, err := u.Store.RoleRepository().ListByUserInstanceId(ctx, instanceId, listOptions)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// FlushLastLoginTime update user login time.
func (u *userServiceImpl) FlushLastLoginTime(ctx context.Context, user *model.User) error {
	now := time.Now()
	user.LastLoginTime = &now

	return u.Store.UserRepository().Update(ctx, user, metav1.UpdateOptions{})
}

// DisableUser disable user
func (u *userServiceImpl) DisableUser(ctx context.Context, instanceId string) error {
	user, err := u.GetUserByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if user.Disabled {
		return errors.WithCode(code.ErrUserAlreadyDisabled, "The user [%s] is already disabled.", user.Name)
	}
	user.Disabled = true

	return u.Store.UserRepository().Update(ctx, user, metav1.UpdateOptions{})
}

// EnableUser disable user
func (u *userServiceImpl) EnableUser(ctx context.Context, instanceId string) error {
	user, err := u.GetUserByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if !user.Disabled {
		return errors.WithCode(code.ErrUserAlreadyEnabled, "The user [%s] is already enabled.", user.Name)
	}
	user.Disabled = false

	return u.Store.UserRepository().Update(ctx, user, metav1.UpdateOptions{})
}

func convertUserModel(user *model.User) *v1.DetailUserResponse {
	return &v1.DetailUserResponse{
		UserBase: *assembler.ConvertUserModelToBase(user),
	}
}
