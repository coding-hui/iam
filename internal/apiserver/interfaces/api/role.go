// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/apiserver/utils"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

type role struct {
	RoleService service.RoleService `inject:""`
}

// NewRole is the of role.
func NewRole() Interface {
	return &role{}
}

func (r *role) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix+"/roles").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("roles"))
	{
		apiv1.POST("", r.createRole)
		apiv1.PUT("/:instanceId", r.updateRole)
		apiv1.DELETE("/:instanceId", r.deleteRole)
		apiv1.GET("/:instanceId", r.roleCheckFilter, r.detailRole)
		apiv1.GET("", r.listRole)
		apiv1.POST("/:instanceId/assign", r.assignRole)
		apiv1.POST("/batch-assign", r.batchAssignRole)
		apiv1.POST("/:instanceId/revoke", r.revokeRole)
	}
}

//	@Tags			Roles
//	@Summary		CreateRole
//	@Description	Create role
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1alpha1.CreateRoleRequest	true	"role info"
//	@Success		200		{object}	api.Response				"create a new role"
//	@Router			/api/v1/roles [post]
//	@Security		BearerTokenAuth
//
// createRole create new role.
func (r *role) createRole(c *gin.Context) {
	createReq := v1.CreateRoleRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = r.RoleService.CreateRole(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Roles
//	@Summary		UpdateRole
//	@Description	Update role
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1alpha1.UpdateRoleRequest	true	"role info"
//	@Success		200		{object}	api.Response				"update role info"
//	@Router			/api/v1/roles/{instanceId}  [put]
//	@Security		BearerTokenAuth
//
// updateRole update role info.
func (r *role) updateRole(c *gin.Context) {
	updateReq := v1.UpdateRoleRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = r.RoleService.UpdateRole(c.Request.Context(), c.Param("instanceId"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Roles
//	@Summary		DeleteRole
//	@Description	DeleteByInstanceId role
//	@Param			name	path		string			true	"identifier of a role"
//	@Success		200		{object}	api.Response	"delete role"
//	@Router			/api/v1/roles/{instanceId} [delete]
//	@Security		BearerTokenAuth
//
// deleteRole delete role by instanceId.
func (r *role) deleteRole(c *gin.Context) {
	err := r.RoleService.DeleteRoleByInstanceId(c.Request.Context(), c.Param("instanceId"), metav1.DeleteOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Roles
//	@Summary		GetRoleInfo
//	@Description	GetByName role info
//	@Param			name	path		string											true	"identifier of a role"
//	@Success		200		{object}	api.Response{data=v1alpha1.DetailRoleResponse}	"role detail"
//	@Router			/api/v1/roles/{instanceId} [get]
//	@Security		BearerTokenAuth
//
// detailRole get role detail info.
func (r *role) detailRole(c *gin.Context) {
	role := c.Request.Context().Value(&v1.CtxKeyRole).(*model.Role)
	detail, err := r.RoleService.DetailRole(c.Request.Context(), role, metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(detail, c)
}

//	@Tags			Roles
//	@Summary		ListRoles
//	@Description	List role
//	@Param			name	query		string									false	"fuzzy search based on name"
//	@Param			offset	query		int										false	"query the page number"
//	@Param			limit	query		int										false	"query the page size number"
//	@Success		200		{object}	api.Response{data=[]v1alpha1.RoleBase}	"roles"
//	@Router			/api/v1/roles [get]
//	@Security		BearerTokenAuth
//
// listRole list role page.
func (r *role) listRole(c *gin.Context) {
	page, pageSize, err := utils.ExtractPagingParams(c, minPageSize, maxPageSize)
	if err != nil {
		api.Fail(c)
		return
	}
	resp, err := r.RoleService.ListRoles(c.Request.Context(), metav1.ListOptions{
		Limit:         &pageSize,
		Offset:        &page,
		FieldSelector: c.Query("fieldSelector"),
	})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}

//	@Tags			Roles
//	@Summary		AssignRole
//	@Description	Assign role
//	@Param			data	body		v1alpha1.AssignRoleRequest	true	"assign role request"
//	@Success		200		{object}	api.Response				"assign role"
//	@Router			/api/v1/roles/{instanceId}/assign [post]
//	@Security		BearerTokenAuth
//
// assignRole assign role.
func (r *role) assignRole(c *gin.Context) {
	assignReq := v1.AssignRoleRequest{}
	err := c.ShouldBindJSON(&assignReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	assignReq.InstanceID = c.Param("instanceId")
	err = r.RoleService.AssignRole(c.Request.Context(), assignReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Roles
//	@Summary		BatchAssignRole
//	@Description	Batch assign role
//	@Param			data	body		v1alpha1.BatchAssignRoleRequest	true	"batch assign role request"
//	@Success		200		{object}	api.Response					"assign role"
//	@Router			/api/v1/roles/{instanceId}/assign [post]
//	@Security		BearerTokenAuth
//
// batchAssignRole assign role.
func (r *role) batchAssignRole(c *gin.Context) {
	batchAssignReq := v1.BatchAssignRoleRequest{}
	err := c.ShouldBindJSON(&batchAssignReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	err = r.RoleService.BatchAssignRole(c.Request.Context(), batchAssignReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Roles
//	@Summary		RevokeRole
//	@Description	Revoke role
//	@Param			data	body		v1alpha1.RevokeRoleRequest	true	"revoke role request"
//	@Success		200		{object}	api.Response				"revoke role"
//	@Router			/api/v1/roles/{instanceId}/revoke [post]
//	@Security		BearerTokenAuth
//
// revokeRole revoke role.
func (r *role) revokeRole(c *gin.Context) {
	revokeReq := v1.RevokeRoleRequest{}
	err := c.ShouldBindJSON(&revokeReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	revokeReq.InstanceID = c.Param("instanceId")
	err = r.RoleService.RevokeRole(c.Request.Context(), revokeReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

func (r *role) roleCheckFilter(c *gin.Context) {
	role, err := r.RoleService.GetRoleByInstanceId(c.Request.Context(), c.Param("instanceId"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		c.Abort()
		return
	}
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1.CtxKeyRole, role))
	c.Next()
}
