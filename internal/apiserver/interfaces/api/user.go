// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	iamv1alpha1 "github.com/coding-hui/api/iam/v1alpha1"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/apiserver/utils"
	"github.com/coding-hui/iam/internal/pkg/api"
)

type user struct {
	UserService service.UserService `inject:""`
}

// NewUser is the of user
func NewUser() Interface {
	return &user{}
}

func (u *user) GetApiGroup() InitApiGroup {
	v1 := InitApiGroup{
		BaseUrl: versionPrefix + "/users",
		Filters: gin.HandlersChain{authCheckFilter},
		Apis: []InitApi{
			{
				Method:  POST,
				Path:    "",
				Handler: u.createUser,
			},
			{
				Method:  PUT,
				Path:    "/:name",
				Handler: u.updateUser,
			},
			{
				Method:  DELETE,
				Path:    "/:name",
				Handler: u.deleteUser,
			},
			{
				Method:  GET,
				Path:    "/:name",
				Handler: u.getUser,
			},
			{
				Method:  GET,
				Path:    "",
				Handler: u.listUser,
			},
		},
	}

	return v1
}

// createUser create new user
func (u *user) createUser(c *gin.Context) {
	user := &iamv1alpha1.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		api.Fail(c)
		return
	}
	err = u.UserService.Create(c.Request.Context(), user, metav1alpha1.CreateOptions{})
	if err != nil {
		api.FailWithMessage(err.Error(), c)
		return
	}

	api.Ok(c)
}

// updateUser update user info
func (u *user) updateUser(c *gin.Context) {
	user := &iamv1alpha1.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		api.Fail(c)
		return
	}
	err = u.UserService.Update(c.Request.Context(), user, metav1alpha1.UpdateOptions{})
	if err != nil {
		api.FailWithMessage(err.Error(), c)
		return
	}

	api.Ok(c)
}

// deleteUser delete user by username
func (u *user) deleteUser(c *gin.Context) {
	err := u.UserService.Delete(c.Request.Context(), c.Param("name"), metav1alpha1.DeleteOptions{})
	if err != nil {
		api.FailWithMessage(err.Error(), c)
		return
	}

	api.Ok(c)
}

// getUser get user info
func (u *user) getUser(c *gin.Context) {
	user, err := u.UserService.Get(c.Request.Context(), c.Param("name"), metav1alpha1.GetOptions{})
	if err != nil {
		api.FailWithMessage(err.Error(), c)
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
		api.FailWithMessage(err.Error(), c)
		return
	}
	api.OkWithPage(resp.Items, resp.GetTotalCount(), c)
}
