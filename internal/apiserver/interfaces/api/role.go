// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/apiserver/utils"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

type role struct {
	RoleService service.RoleService `inject:""`
}

// NewRole is the of role.
func NewRole() Interface {
	return &role{}
}

func (r *role) RegisterApiGroup(g *gin.Engine) {
	v1 := g.Group(versionPrefix + "/roles").Use(authCheckFilter)
	{
		v1.POST("", r.createRole)
		v1.PUT("/:instanceId", r.updateRole)
		v1.DELETE("/:instanceId", r.deleteRole)
		v1.GET("/:instanceId", r.roleCheckFilter, r.detailRole)
		v1.GET("", r.listRole)
		v1.POST("/:instanceId/assign", r.roleCheckFilter, r.assignRole)
		v1.POST("/:instanceId/revoke", r.roleCheckFilter, r.revokeRole)
	}
}

// createRole create new role.
func (r *role) createRole(c *gin.Context) {
	createReq := v1alpha1.CreateRoleRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	err = r.RoleService.CreateRole(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// updateRole update role info.
func (r *role) updateRole(c *gin.Context) {
	updateReq := v1alpha1.UpdateRoleRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	err = r.RoleService.UpdateRole(c.Request.Context(), c.Param("name"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// deleteRole delete role by roleName.
func (r *role) deleteRole(c *gin.Context) {
	err := r.RoleService.DeleteRole(c.Request.Context(), c.Param("name"), metav1alpha1.DeleteOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// detailRole get role detail info.
func (r *role) detailRole(c *gin.Context) {
	role := c.Request.Context().Value(&v1alpha1.CtxKeyRole).(*model.Role)
	detail, err := r.RoleService.DetailRole(c.Request.Context(), role, metav1alpha1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(detail, c)
}

// listRole list role page.
func (r *role) listRole(c *gin.Context) {
	page, pageSize, err := utils.ExtractPagingParams(c, minPageSize, maxPageSize)
	if err != nil {
		api.Fail(c)
		return
	}
	resp, err := r.RoleService.List(c.Request.Context(), metav1alpha1.ListOptions{
		Limit:  &pageSize,
		Offset: &page,
	})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}

// assignRole assign role.
func (r *role) assignRole(c *gin.Context) {
	assignReq := v1alpha1.AssignRoleRequest{}
	err := c.ShouldBindJSON(&assignReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	role := c.Request.Context().Value(&v1alpha1.CtxKeyRole).(*model.Role)
	err = r.RoleService.AssignRole(c.Request.Context(), role, assignReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// revokeRole revoke role.
func (r *role) revokeRole(c *gin.Context) {
	revokeReq := v1alpha1.RevokeRoleRequest{}
	err := c.ShouldBindJSON(&revokeReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	role := c.Request.Context().Value(&v1alpha1.CtxKeyRole).(*model.Role)
	err = r.RoleService.RevokeRole(c.Request.Context(), role, revokeReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

func (r *role) roleCheckFilter(c *gin.Context) {
	role, err := r.RoleService.GetRole(c.Request.Context(), c.Param("instanceId"), metav1alpha1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		c.Abort()
		return
	}
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1alpha1.CtxKeyRole, role))
	c.Next()
}
