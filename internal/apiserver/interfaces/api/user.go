// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/apiserver/utils"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
)

type user struct {
	UserService service.UserService `inject:""`
}

// NewUser is the of user
func NewUser() Interface {
	return &user{}
}

func (u *user) RegisterApiGroup(g *gin.Engine) {
	v1 := g.Group(versionPrefix + "/users").Use(authCheckFilter)
	{
		v1.POST("", u.createUser)
		v1.PUT("/:name", u.updateUser)
		v1.DELETE("/:name", u.deleteUser)
		v1.GET("/:name", u.getUser)
		v1.GET("", u.listUser)
	}
}

// createUser create new user
func (u *user) createUser(c *gin.Context) {
	createReq := v1alpha1.CreateUserRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = u.UserService.Create(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// updateUser update user info
func (u *user) updateUser(c *gin.Context) {
	updateReq := v1alpha1.UpdateUserRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = u.UserService.Update(c.Request.Context(), c.Param("name"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// deleteUser delete user by username
func (u *user) deleteUser(c *gin.Context) {
	err := u.UserService.Delete(c.Request.Context(), c.Param("name"), metav1alpha1.DeleteOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// getUser get user info
func (u *user) getUser(c *gin.Context) {
	user, err := u.UserService.Get(c.Request.Context(), c.Param("name"), metav1alpha1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(user, c)
}

// listUser list user page
func (u *user) listUser(c *gin.Context) {
	page, pageSize, err := utils.ExtractPagingParams(c, minPageSize, maxPageSize)
	if err != nil {
		api.Fail(c)
		return
	}
	resp, err := u.UserService.List(c.Request.Context(), metav1alpha1.ListOptions{
		Limit:  &pageSize,
		Offset: &page,
	})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}
