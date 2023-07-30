// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1alpha1"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

// RoleService Role manage api.
type RoleService interface {
	CreateRole(ctx context.Context, req v1alpha1.CreateRoleRequest) error
	UpdateRole(ctx context.Context, roleName string, req v1alpha1.UpdateRoleRequest) error
	DeleteRoleByInstanceId(ctx context.Context, instanceId string, opts metav1alpha1.DeleteOptions) error
	BatchDeleteRoles(ctx context.Context, roleNames []string, opts metav1alpha1.DeleteOptions) error
	GetRoleByInstanceId(ctx context.Context, instanceId string, opts metav1alpha1.GetOptions) (*model.Role, error)
	DetailRole(ctx context.Context, role *model.Role, opts metav1alpha1.GetOptions) (*v1alpha1.DetailRoleResponse, error)
	ListRoles(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.RoleList, error)
	AssignRole(ctx context.Context, role *model.Role, assignReq v1alpha1.AssignRoleRequest) error
	BatchAssignRole(ctx context.Context, assignReq v1alpha1.BatchAssignRoleRequest) error
	RevokeRole(ctx context.Context, role *model.Role, revokeReq v1alpha1.RevokeRoleRequest) error
	AuthorizeRoleResources(ctx context.Context, role *model.Role, authorizeReq v1alpha1.AuthorizeResources) error
	Init(ctx context.Context) error
}

type roleServiceImpl struct {
	Store       repository.Factory `inject:"repository"`
	UserService UserService        `inject:""`
}

// NewRoleService new Role service.
func NewRoleService() RoleService {
	return &roleServiceImpl{}
}

type roleHandler interface {
	assign(ctx context.Context) error
	revoke(ctx context.Context) error
}

type userRoleHandlerImpl struct {
	userService       UserService
	store             repository.Factory
	role              *model.Role
	targetInstanceIds []string
}

func (r *roleServiceImpl) newUserRoleHandlerImpl(role *model.Role, targetInstanceIds []string) (*userRoleHandlerImpl, error) {
	return &userRoleHandlerImpl{
		store:             r.Store,
		userService:       r.UserService,
		role:              role,
		targetInstanceIds: targetInstanceIds,
	}, nil
}

// Init initialize role data.
func (r *roleServiceImpl) Init(ctx context.Context) error {
	// initialize default roles before
	initRoles := []v1alpha1.UserType{v1alpha1.PlatformAdmin, v1alpha1.TenantAdmin, v1alpha1.Default}
	for _, role := range initRoles {
		_, err := r.Store.RoleRepository().GetByName(ctx, role.String(), metav1alpha1.GetOptions{})
		if err != nil && errors.IsCode(err, code.ErrRoleNotFound) {
			createReq := v1alpha1.CreateRoleRequest{
				Owner:       DefaultAdmin,
				Name:        role.String(),
				DisplayName: "",
				Description: "",
			}
			err = r.CreateRole(ctx, createReq)
			if err != nil {
				return errors.WithMessagef(err, "Failed to initialize default role %s", role.String())
			}
			log.Infof("initialize %s role done", role.String())
		}
	}

	// find platform role
	platformRole, err := r.Store.RoleRepository().GetByName(ctx, v1alpha1.PlatformAdmin.String(), metav1alpha1.GetOptions{})
	if err != nil {
		return errors.WithMessagef(err, "Failed to get %s role info.", v1alpha1.PlatformAdmin.String())
	}

	// find default platform admin
	admin, err := r.Store.UserRepository().GetByName(ctx, DefaultAdmin, metav1alpha1.GetOptions{})
	if err != nil || admin == nil {
		return errors.WithMessagef(err, "Failed to get default admin info. Please check UserService is initialized.")
	}

	// assign the platform role to the default administrator
	assignRoleReq := v1alpha1.AssignRoleRequest{
		InstanceID: platformRole.InstanceID,
		Targets:    []string{admin.InstanceID},
	}
	err = r.AssignRole(ctx, platformRole, assignRoleReq)
	if err != nil {
		return errors.WithMessagef(err, "Failed to assign the platform role to the default administrator.")
	}
	log.Infof("assign the platform role to the default administrator done")

	return nil
}

// CreateRole create role.
func (r *roleServiceImpl) CreateRole(ctx context.Context, req v1alpha1.CreateRoleRequest) error {
	role := &model.Role{
		ObjectMeta: metav1alpha1.ObjectMeta{
			Name: req.Name,
		},
		Owner:       req.Owner,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Disabled:    false,
	}
	if err := r.Store.RoleRepository().Create(ctx, role, metav1alpha1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

// UpdateRole update role.
func (r *roleServiceImpl) UpdateRole(ctx context.Context, roleName string, req v1alpha1.UpdateRoleRequest) error {
	role, err := r.GetRoleByInstanceId(ctx, roleName, metav1alpha1.GetOptions{})
	if err != nil {
		return err
	}
	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Owner != "" {
		role.Owner = req.Owner
	}
	if err := r.Store.RoleRepository().Update(ctx, role, metav1alpha1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

// DeleteRoleByInstanceId delete role by instanceId.
func (r *roleServiceImpl) DeleteRoleByInstanceId(ctx context.Context, instanceId string, opts metav1alpha1.DeleteOptions) error {
	role, err := r.GetRoleByInstanceId(ctx, instanceId, metav1alpha1.GetOptions{})
	if err != nil {
		return err
	}
	if len(role.Users) > 0 {
		return errors.WithCode(code.ErrRoleHasAssignedUser, "Failed to delete role. The role has been assigned to a user")
	}
	if err := r.Store.RoleRepository().DeleteByInstanceId(ctx, instanceId, opts); err != nil {
		return err
	}

	return nil
}

// BatchDeleteRoles batch delete role.
func (r *roleServiceImpl) BatchDeleteRoles(ctx context.Context, roleNames []string, opts metav1alpha1.DeleteOptions) error {
	if err := r.Store.RoleRepository().DeleteCollection(ctx, roleNames, opts); err != nil {
		return err
	}

	return nil
}

// GetRoleByInstanceId get role by instanceId.
func (r *roleServiceImpl) GetRoleByInstanceId(ctx context.Context, instanceId string, opts metav1alpha1.GetOptions) (*model.Role, error) {
	role, err := r.Store.RoleRepository().GetByInstanceID(ctx, instanceId, opts)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// DetailRole get role detail info.
func (r *roleServiceImpl) DetailRole(
	ctx context.Context,
	role *model.Role,
	_ metav1alpha1.GetOptions,
) (*v1alpha1.DetailRoleResponse, error) {
	base := assembler.ConvertRoleModelToBase(role)
	users := make([]v1alpha1.UserBase, 0, len(role.Users))
	for _, u := range role.Users {
		users = append(users, *assembler.ConvertUserModelToBase(&u))
	}
	detail := &v1alpha1.DetailRoleResponse{
		RoleBase: *base,
		Users:    users,
	}

	return detail, nil
}

// ListRoles list roles.
func (r *roleServiceImpl) ListRoles(ctx context.Context, listOptions metav1alpha1.ListOptions) (*v1alpha1.RoleList, error) {
	roles, err := r.Store.RoleRepository().List(ctx, metav1alpha1.ListOptions{
		Offset: listOptions.Offset,
		Limit:  listOptions.Limit,
	})
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// AssignRole assign roles, which can be users or departments.
func (r *roleServiceImpl) AssignRole(ctx context.Context, role *model.Role, assignReq v1alpha1.AssignRoleRequest) (lastErr error) {
	handlers, err := r.determineRoleHandlerByInstanceId(role, assignReq.Targets)
	if err != nil {
		return errors.WithCode(code.ErrAssignRoleFailed, err.Error())
	}

	for h := range handlers {
		lastErr = h.assign(ctx)
		if lastErr != nil {
			log.Errorf("Failed to assign roles. err: %w", lastErr)
		}
	}

	return lastErr
}

// BatchAssignRole batch assign roles, which can be users or departments.
func (r *roleServiceImpl) BatchAssignRole(ctx context.Context, batchAssignReq v1alpha1.BatchAssignRoleRequest) error {
	instanceIds := batchAssignReq.InstanceIds
	targets := batchAssignReq.Targets
	for _, instanceId := range instanceIds {
		role, err := r.GetRoleByInstanceId(ctx, instanceId, metav1alpha1.GetOptions{})
		if err != nil {
			log.Errorf("Failed to get role details using the id", instanceId)
			return err
		}
		err = r.AssignRole(ctx, role, v1alpha1.AssignRoleRequest{Targets: targets})
		if err != nil {
			log.Errorf("Failed to batch assign roles. err: %w", err)
			return err
		}
	}

	return nil
}

// RevokeRole revoke roles, which can be users or departments.
func (r *roleServiceImpl) RevokeRole(ctx context.Context, role *model.Role, revokeReq v1alpha1.RevokeRoleRequest) (lastErr error) {
	handlers, err := r.determineRoleHandlerByInstanceId(role, revokeReq.Targets)
	if err != nil {
		return errors.WithCode(code.ErrRevokeRoleFailed, err.Error())
	}

	for h := range handlers {
		lastErr = h.revoke(ctx)
		if lastErr != nil {
			log.Errorf("Failed to revoke roles. err: %w", lastErr)
		}
	}

	return lastErr
}

// AuthorizeRoleResources grant role resource permission.
func (r *roleServiceImpl) AuthorizeRoleResources(ctx context.Context, role *model.Role, authorizeReq v1alpha1.AuthorizeResources) error {
	handlers, err := r.determineRoleHandlerByInstanceId(role, authorizeReq.Resources)
	if err != nil {
		return errors.WithCode(code.ErrRevokeRoleFailed, err.Error())
	}

	for h := range handlers {
		err := h.revoke(ctx)
		if err != nil {
			log.Errorf("Failed to revoke roles. err: %w", err)
		}
	}

	return nil
}

func (u *userRoleHandlerImpl) assign(ctx context.Context) error {
	count, err := u.store.RoleRepository().AssignUserRoles(ctx, u.role, u.targetInstanceIds)
	if err != nil {
		return err
	}
	e := u.store.CasbinRepository().SyncedEnforcer()
	for _, userInstanceId := range u.targetInstanceIds {
		_, err = e.AddRoleForUser(userInstanceId, u.role.InstanceID)
		if err != nil {
			log.Errorf("Failed assign role %s to user %s", u.role.Name, userInstanceId)
		}
	}
	log.Infof("AssignRole the %s role to %d users", u.role.Name, count)

	return nil
}

func (u *userRoleHandlerImpl) revoke(ctx context.Context) error {
	count, err := u.store.RoleRepository().RevokeUserRoles(ctx, u.role, u.targetInstanceIds)
	if err != nil {
		return err
	}
	e := u.store.CasbinRepository().SyncedEnforcer()
	for _, userInstanceId := range u.targetInstanceIds {
		_, err = e.DeleteRoleForUser(userInstanceId, u.role.InstanceID)
		if err != nil {
			log.Errorf("Failed to revoke the %s role rights of the %s", u.role.Name, userInstanceId)
		}
	}
	log.Infof("RevokeRole the %s role removes %d users", u.role.Name, count)

	return nil
}

func (r *roleServiceImpl) determineRoleHandlerByInstanceId(role *model.Role, instanceIds []string) (map[roleHandler][]string, error) {
	targets := make(map[string][]string, len(instanceIds))
	handlers := make(map[roleHandler][]string, len(instanceIds))

	for _, t := range instanceIds {
		ri := model.GetResourceIdentifier(t)
		targets[ri] = append(targets[ri], t)
	}

	for ri, ts := range targets {
		switch ri {
		case v1alpha1.UserTarget:
			handler, err := r.newUserRoleHandlerImpl(role, ts)
			if err != nil {
				log.Warnf("failed to new user role handler. err: %w", err)
				continue
			}
			handlers[handler] = ts
		default:
			log.Warnf("Unsupported target %s", ri)
		}
	}

	return handlers, nil
}
