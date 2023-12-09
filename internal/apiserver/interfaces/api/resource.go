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
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

type resource struct {
	ResourceService service.ResourceService `inject:""`
}

// NewResource is the of resource.
func NewResource() Interface {
	return &resource{}
}

func (r *resource) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix+"/resources").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("resources"))
	{
		apiv1.POST("", r.createResource)
		apiv1.PUT("/:instanceId", r.updateResource)
		apiv1.DELETE("/:instanceId", r.deleteResource)
		apiv1.GET("/:instanceId", r.resourceCheckFilter, r.detailResource)
		apiv1.GET("", r.listResource)
	}
}

//	@Tags			Resource
//	@Summary		CreateResource
//	@Description	Create resource
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreateResourceRequest	true	"resource info"
//	@Success		200		{object}	api.Response				"create a new resource"
//	@Router			/api/v1/resources [post]
//	@Security		BearerTokenAuth
//
// createResource create new resource.
func (r *resource) createResource(c *gin.Context) {
	createReq := v1.CreateResourceRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = r.ResourceService.CreateResource(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Resource
//	@Summary		UpdateResource
//	@Description	Update resource
//	@Accept			application/json
//	@Product		application/json
//	@Param			data		body		v1.UpdateResourceRequest	true	"resource info"
//	@Param			instanceId	path		string						true	"identifier of a resource"
//	@Success		200			{object}	api.Response				"update resource info"
//	@Router			/api/v1/resources/{instanceId}  [put]
//	@Security		BearerTokenAuth
//
// updateResource update resource info.
func (r *resource) updateResource(c *gin.Context) {
	updateReq := v1.UpdateResourceRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = r.ResourceService.UpdateResource(c.Request.Context(), c.Param("instanceId"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Resource
//	@Summary		DeleteResource
//	@Description	DeleteByInstanceId resource
//	@Param			instanceId	path		string			true	"identifier of a resource"
//	@Success		200			{object}	api.Response	"delete resource"
//	@Router			/api/v1/resources/{instanceId} [delete]
//	@Security		BearerTokenAuth
//
// deleteResource delete resource by name.
func (r *resource) deleteResource(c *gin.Context) {
	err := r.ResourceService.DeleteResource(
		c.Request.Context(),
		c.Param("instanceId"),
		metav1.DeleteOptions{},
	)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Resource
//	@Summary		GetResourceInfo
//	@Description	Get resource info
//	@Param			instanceId	path		string											true	"identifier of a resource"
//	@Success		200			{object}	api.Response{data=v1.DetailResourceResponse}	"resource detail"
//	@Router			/api/v1/resources/{instanceId} [get]
//	@Security		BearerTokenAuth
//
// detailResource get resource info.
func (r *resource) detailResource(c *gin.Context) {
	resource := c.Request.Context().Value(&v1.CtxKeyResource).(*model.Resource)
	detail, err := r.ResourceService.DetailResource(c.Request.Context(), resource, metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(detail, c)
}

//	@Tags			Resource
//	@Summary		ListResources
//	@Description	List resources
//	@Param			name	query		string								false	"fuzzy search based on name"
//	@Param			offset	query		int									false	"query the page number"
//	@Param			limit	query		int									false	"query the page size number"
//	@Success		200		{object}	api.Response{data=v1.ResourceList}	"resources"
//	@Router			/api/v1/resources [get]
//	@Security		BearerTokenAuth
//
// listResource list resource page.
func (r *resource) listResource(c *gin.Context) {
	var opts metav1.ListOptions
	err := c.ShouldBindQuery(&opts)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	resp, err := r.ResourceService.ListResources(c.Request.Context(), opts)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}

func (r *resource) resourceCheckFilter(c *gin.Context) {
	resource, err := r.ResourceService.GetResource(
		c.Request.Context(),
		c.Param("instanceId"),
		metav1.GetOptions{},
	)
	if err != nil {
		api.FailWithErrCode(err, c)
		c.Abort()
		return
	}
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1.CtxKeyResource, resource))
	c.Next()
}
