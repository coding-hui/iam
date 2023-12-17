// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

type user struct {
	UserService service.UserService `inject:""`
}

// NewUser is the of user.
func NewUser() Interface {
	return &user{}
}

func (u *user) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix+"/users").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("users"))
	{
		apiv1.POST("", u.createUser)
		apiv1.PUT("/:instanceId", u.updateUser)
		apiv1.DELETE("/:instanceId", u.deleteUser)
		apiv1.GET("/:instanceId", u.getUser)
		apiv1.GET("", u.listUser)
		apiv1.GET("/:instanceId/roles", u.getUserRoles)
		apiv1.GET("/:instanceId/disable", u.disableUser)
		apiv1.GET("/:instanceId/enable", u.enableUser)
	}
}

//	@Tags			Users
//	@Summary		CreateUser
//	@Description	Create user
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreateUserRequest	true	"user info"
//	@Success		200		{object}	api.Response			"create a nnw user"
//	@Router			/api/v1/users [post]
//	@Security		BearerTokenAuth
//
// createUser create a new user.
func (u *user) createUser(c *gin.Context) {
	createReq := v1.CreateUserRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	user, err := u.UserService.CreateUser(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(user, c)
}

//	@Tags			Users
//	@Summary		UpdateUser
//	@Description	Update user info
//	@Accept			application/json
//	@Product		application/json
//	@Param			instanceId	path		string					true	"identifier of a user"
//	@Param			data		body		v1.UpdateUserRequest	true	"user info"
//	@Success		200			{object}	api.Response			"update user info"
//	@Router			/api/v1/users/{instanceId} [put]
//	@Security		BearerTokenAuth
//
// updateUser update user info.
func (u *user) updateUser(c *gin.Context) {
	updateReq := v1.UpdateUserRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	user, err := u.UserService.UpdateUser(c.Request.Context(), c.Param("instanceId"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(user, c)
}

//	@Tags			Users
//	@Summary		DeleteUser
//	@Description	DeleteByInstanceId user
//	@Param			instanceId	path		string			true	"identifier of a user"
//	@Success		200			{object}	api.Response	"delete user"
//	@Router			/api/v1/users/{instanceId} [delete]
//	@Security		BearerTokenAuth
//
// deleteUser delete user by identifier.
func (u *user) deleteUser(c *gin.Context) {
	err := u.UserService.DeleteUser(c.Request.Context(), c.Param("instanceId"), metav1.DeleteOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Users
//	@Summary		GetUserInfo
//	@Description	GetByName user info
//	@Param			instanceId	path		string							true	"identifier of a user"
//	@Success		200			{object}	api.Response{data=v1.UserBase}	"user detail"
//	@Router			/api/v1/users/{instanceId} [get]
//	@Security		BearerTokenAuth
//
// getUser get user detail.
func (u *user) getUser(c *gin.Context) {
	user, err := u.UserService.GetUserByInstanceId(c.Request.Context(), c.Param("instanceId"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	resp, err := u.UserService.DetailUser(c.Request.Context(), user)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(resp, c)
}

//	@Tags			Users
//	@Summary		ListUsers
//	@Description	List users
//	@Param			name	query		string							false	"fuzzy search based on name"
//	@Param			alias	query		string							false	"fuzzy search based on alias"
//	@Param			email	query		string							false	"fuzzy search based on email"
//	@Param			offset	query		int								false	"query the page number"
//	@Param			limit	query		int								false	"query the page size number"
//	@Success		200		{object}	api.Response{data=v1.UserList}	"users"
//	@Router			/api/v1/users [get]
//	@Security		BearerTokenAuth
//
// listUser list users page.
func (u *user) listUser(c *gin.Context) {
	var opts v1.ListUserOptions
	err := c.ShouldBindQuery(&opts)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	resp, err := u.UserService.ListUsers(c.Request.Context(), opts)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}

//	@Tags			Users
//	@Summary		GetUserRoles
//	@Description	Get user roles
//	@Param			instanceId	path		string							true	"identifier of a user"
//	@Success		200			{object}	api.Response{data=v1.RoleList}	"user roles"
//	@Router			/api/v1/users/{instanceId}/roles [get]
//	@Security		BearerTokenAuth
//
// getUserRoles get user roles.
func (u *user) getUserRoles(c *gin.Context) {
	roles, err := u.UserService.ListUserRoles(c.Request.Context(), c.Param("instanceId"), metav1.ListOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(roles.Items, roles.TotalCount, c)
}

//	@Tags			Users
//	@Summary		EnableUser
//	@Description	Enable a user
//	@Param			instanceId	path		string	true	"identifier of a user"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/users/{instanceId}/disable [get]
//	@Security		BearerTokenAuth
//
// getUserRoles get user roles.
func (u *user) disableUser(c *gin.Context) {
	err := u.UserService.DisableUser(c.Request.Context(), c.Param("instanceId"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Users
//	@Summary		EnableUser
//	@Description	Enable a user
//	@Param			instanceId	path		string	true	"identifier of a user"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/users/{instanceId}/enable [get]
//	@Security		BearerTokenAuth
//
// getUserRoles get user roles.
func (u *user) enableUser(c *gin.Context) {
	err := u.UserService.EnableUser(c.Request.Context(), c.Param("instanceId"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}
