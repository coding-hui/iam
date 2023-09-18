// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/apiserver/utils"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

type organization struct {
	OrganizationService service.OrganizationService `inject:""`
}

// NewOrganization is the of organization.
func NewOrganization() Interface {
	return &organization{}
}

func (o *organization) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix+"/organizations").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("organizations"))
	{
		apiv1.POST("", o.createOrganization)
		apiv1.PUT("/:instanceId", o.updateOrganization)
		apiv1.DELETE("/:instanceId", o.deleteOrganization)
		apiv1.GET("/:instanceId", o.getOrganization)
		apiv1.GET("", o.listOrganization)
		apiv1.GET("/:instanceId/disable", o.disableOrganization)
		apiv1.GET("/:instanceId/enable", o.enableOrganization)
	}
}

//	@Tags			Organizations
//	@Summary		CreateOrganization
//	@Description	Create organization
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreateOrganizationRequest	true	"organization info"
//	@Success		200		{object}	api.Response					"create a nnw organization"
//	@Router			/api/v1/organizations [post]
//	@Security		BearerTokenAuth
//
// createOrganization create a new organization.
func (o *organization) createOrganization(c *gin.Context) {
	createReq := v1.CreateOrganizationRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = o.OrganizationService.CreateOrganization(c.Request.Context(), createReq, metav1.CreateOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Organizations
//	@Summary		UpdateOrganization
//	@Description	Update organization info
//	@Accept			application/json
//	@Product		application/json
//	@Param			instanceId	path		string							true	"identifier of a organization"
//	@Param			data		body		v1.UpdateOrganizationRequest	true	"organization info"
//	@Success		200			{object}	api.Response					"update organization info"
//	@Router			/api/v1/organizations/{instanceId} [put]
//	@Security		BearerTokenAuth
//
// updateOrganization update organization info.
func (o *organization) updateOrganization(c *gin.Context) {
	updateReq := v1.UpdateOrganizationRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = o.OrganizationService.UpdateOrganization(c.Request.Context(), c.Param("instanceId"), updateReq, metav1.UpdateOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Organizations
//	@Summary		DeleteOrganization
//	@Description	Delete organization
//	@Param			instanceId	path		string			true	"identifier of a organization"
//	@Success		200			{object}	api.Response	"delete organization"
//	@Router			/api/v1/organizations/{instanceId} [delete]
//	@Security		BearerTokenAuth
//
// deleteOrganization delete organization by identifier.
func (o *organization) deleteOrganization(c *gin.Context) {
	err := o.OrganizationService.DeleteOrganization(c.Request.Context(), c.Param("instanceId"), metav1.DeleteOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Organizations
//	@Summary		GetOrganizationInfo
//	@Description	Get organization info
//	@Param			instanceId	path		string									true	"identifier of a organization"
//	@Success		200			{object}	api.Response{data=v1.OrganizationBase}	"organization detail"
//	@Router			/api/v1/organizations/{instanceId} [get]
//	@Security		BearerTokenAuth
//
// getOrganization get organization detail.
func (o *organization) getOrganization(c *gin.Context) {
	organization, err := o.OrganizationService.GetOrganization(c.Request.Context(), c.Param("instanceId"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	resp, err := o.OrganizationService.DetailOrganization(c.Request.Context(), organization)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(resp, c)
}

//	@Tags			Organizations
//	@Summary		ListOrganizations
//	@Description	List organizations
//	@Param			offset	query		int										false	"query the page number"
//	@Param			limit	query		int										false	"query the page size number"
//	@Success		200		{object}	api.Response{data=v1.OrganizationList}	"organizations"
//	@Router			/api/v1/organizations [get]
//	@Security		BearerTokenAuth
//
// listOrganization list organizations page.
func (o *organization) listOrganization(c *gin.Context) {
	page, pageSize, err := utils.ExtractPagingParams(c, minPageSize, maxPageSize)
	if err != nil {
		api.Fail(c)
		return
	}
	resp, err := o.OrganizationService.ListOrganizations(c.Request.Context(), metav1.ListOptions{
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

//	@Tags			Organizations
//	@Summary		EnableOrganization
//	@Description	Enable a organization
//	@Param			instanceId	path		string	true	"identifier of a organization"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/organizations/{instanceId}/disable [get]
//	@Security		BearerTokenAuth
//
// getOrganizationRoles get organization roles.
func (o *organization) disableOrganization(c *gin.Context) {
	err := o.OrganizationService.DisableOrganization(c.Request.Context(), c.Param("instanceId"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Organizations
//	@Summary		EnableOrganization
//	@Description	Enable a organization
//	@Param			instanceId	path		string	true	"identifier of a organization"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/organizations/{instanceId}/enable [get]
//	@Security		BearerTokenAuth
//
// getOrganizationRoles get organization roles.
func (o *organization) enableOrganization(c *gin.Context) {
	err := o.OrganizationService.EnableOrganization(c.Request.Context(), c.Param("instanceId"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}
