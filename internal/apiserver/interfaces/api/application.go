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

type application struct {
	ApplicationService service.ApplicationService `inject:""`
}

// NewApplication is the of application.
func NewApplication() Interface {
	return &application{}
}

func (a *application) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix+"/applications").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("applications"))
	{
		apiv1.POST("", a.createApplication)
		apiv1.PUT("/:idOrName", a.updateApplication)
		apiv1.DELETE("/:instanceId", a.deleteApplication)
		apiv1.GET("/:idOrName", a.detailApplication)
		apiv1.GET("", a.listApplications)
	}

	publicApi := g.Group(versionPrefix + "/applications/public")
	{
		publicApi.GET("/:idOrName/config", a.detailApplication)
	}
}

//	@Tags			Applications
//	@Summary		CreateApplication
//	@Description	Create a application
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreateApplicationRequest	true	"application"
//	@Success		200		{object}	api.Response				"Create a new application"
//	@Router			/api/v1/applications [post]
//	@Security		BearerTokenAuth
//
// createApplication create a new application.
func (a *application) createApplication(c *gin.Context) {
	createReq := v1.CreateApplicationRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = a.ApplicationService.CreateApplication(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Applications
//	@Summary		UpdateApplication
//	@Description	Update a application
//	@Accept			application/json
//	@Product		application/json
//	@Param			data		body		v1.UpdateApplicationRequest	true	"application"
//	@Param			instanceId	path		string						true	"identifier of a application"
//	@Success		200			{object}	api.Response				"Update application info"
//	@Router			/api/v1/applications/{idOrName} [put]
//	@Security		BearerTokenAuth
//
// updateApplication update application info.
func (a *application) updateApplication(c *gin.Context) {
	updateReq := v1.UpdateApplicationRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = a.ApplicationService.UpdateApplication(c.Request.Context(), c.Param("idOrName"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Applications
//	@Summary		DeleteApplication
//	@Description	Delete application by instanceId
//	@Param			instanceId	path		string			true	"identifier of a application"
//	@Success		200			{object}	api.Response	"Application successfully deleted"
//	@Router			/api/v1/applications/{instanceId} [DELETE]
//	@Security		BearerTokenAuth
//
// deleteApplication delete application by identifier.
func (a *application) deleteApplication(c *gin.Context) {
	err := a.ApplicationService.DeleteApplication(
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

//	@Tags			Applications
//	@Summary		GetApplicationInfo
//	@Description	Get a application by identifier
//	@Param			instanceId	path		string											true	"identifier of a application"
//	@Success		200			{object}	api.Response{data=v1.DetailApplicationResponse}	"Application detail"
//	@Router			/api/v1/applications/{idOrName} [get]
//	@Security		BearerTokenAuth
//
// detailApplication get application detail info.
func (a *application) detailApplication(c *gin.Context) {
	detail, err := a.ApplicationService.GetApplication(c.Request.Context(), c.Param("idOrName"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(detail, c)
}

//	@Tags			ListApplications
//	@Summary		ListApplications
//	@Description	List Applications
//	@Param			offset	query		int										false	"query the page number"
//	@Param			limit	query		int										false	"query the page size number"
//	@Success		200		{object}	api.Response{data=v1.ApplicationList}	"Applications"
//	@Router			/api/v1/applications [get]
//	@Security		BearerTokenAuth
//
// listApplications list Application page.
func (a *application) listApplications(c *gin.Context) {
	var opts metav1.ListOptions
	err := c.ShouldBindQuery(&opts)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	resp, err := a.ApplicationService.ListApplications(c.Request.Context(), opts)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}
