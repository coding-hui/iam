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

type resource struct {
	ResourceService service.ResourceService `inject:""`
}

// NewResource is the of resource.
func NewResource() Interface {
	return &resource{}
}

func (r *resource) RegisterApiGroup(g *gin.Engine) {
	v1 := g.Group(versionPrefix+"/resources").Use(authCheckFilter, permissionCheckFilter)
	{
		v1.POST("", r.createResource)
		v1.PUT("/:name", r.updateResource)
		v1.DELETE("/:name", r.deleteResource)
		v1.GET("/:name", r.getResource)
		v1.GET("", r.listResource)
	}
}

// createResource create new resource.
func (r *resource) createResource(c *gin.Context) {
	createReq := v1alpha1.CreateResourceRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	err = r.ResourceService.CreateResource(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// updateResource update resource info.
func (r *resource) updateResource(c *gin.Context) {
	updateReq := v1alpha1.UpdateResourceRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	err = r.ResourceService.UpdateResource(c.Request.Context(), c.Param("name"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// deleteResource delete resource by name.
func (r *resource) deleteResource(c *gin.Context) {
	err := r.ResourceService.DeleteResource(
		c.Request.Context(),
		c.Param("name"),
		metav1alpha1.DeleteOptions{},
	)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// getResource get resource info.
func (r *resource) getResource(c *gin.Context) {
	resource, err := r.ResourceService.GetResource(
		c.Request.Context(),
		c.Param("name"),
		metav1alpha1.GetOptions{},
	)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(resource, c)
}

// listResource list resource page.
func (r *resource) listResource(c *gin.Context) {
	page, pageSize, err := utils.ExtractPagingParams(c, minPageSize, maxPageSize)
	if err != nil {
		api.Fail(c)
		return
	}
	resp, err := r.ResourceService.ListResources(c.Request.Context(), metav1alpha1.ListOptions{
		Limit:  &pageSize,
		Offset: &page,
	})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}
