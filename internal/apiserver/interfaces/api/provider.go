// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

type provider struct {
	ProviderService service.ProviderService `inject:""`
}

// NewProvider is the of provider.
func NewProvider() Interface {
	return &provider{}
}

func (p *provider) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix+"/providers").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("providers"))
	{
		apiv1.POST("", p.createProvider)
		apiv1.PUT("/:instanceId", p.updateProvider)
		apiv1.DELETE("/:instanceId", p.deleteProvider)
		apiv1.GET("/:instanceId", p.detailProvider)
	}
}

//	@Tags			Providers
//	@Summary		CreateProvider
//	@Description	Create a provider
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreateProviderRequest	true	"provider"
//	@Success		200		{object}	api.Response			"Create a nnw provider"
//	@Router			/api/v1/providers [post]
//	@Security		BearerTokenAuth
//
// createProvider create a new provider.
func (p *provider) createProvider(c *gin.Context) {
	createReq := v1.CreateProviderRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = p.ProviderService.CreateProvider(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Providers
//	@Summary		UpdateProvider
//	@Description	Update a provider
//	@Accept			application/json
//	@Product		application/json
//	@Param			data		body		v1.UpdateProviderRequest	true	"provider"
//	@Param			instanceId	path		string					true	"identifier of a provider"
//	@Success		200			{object}	api.Response			"Update provider info"
//	@Router			/api/v1/providers/{instanceId} [put]
//	@Security		BearerTokenAuth
//
// updateProvider update provider info.
func (p *provider) updateProvider(c *gin.Context) {
	updateReq := v1.UpdateProviderRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = p.ProviderService.UpdateProvider(c.Request.Context(), c.Param("instanceId"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Providers
//	@Summary		DeleteProvider
//	@Description	DeleteByInstanceId provider by instanceId
//	@Param			instanceId	path		string			true	"identifier of a provider"
//	@Success		200			{object}	api.Response	"Provider successfully deleted"
//	@Router			/api/v1/providers/{instanceId} [DELETE]
//	@Security		BearerTokenAuth
//
// deleteProvider delete provider by identifier.
func (p *provider) deleteProvider(c *gin.Context) {
	err := p.ProviderService.DeleteProvider(
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

//	@Tags			Providers
//	@Summary		GetProviderInfo
//	@Description	GetByName a provider by name
//	@Param			instanceId	path		string										true	"name of a provider"
//	@Success		200			{object}	api.Response{data=v1.DetailProviderResponse}	"Provider detail"
//	@Router			/api/v1/providers/{instanceId} [get]
//	@Security		BearerTokenAuth
//
// detailProvider get provider detail info.
func (p *provider) detailProvider(c *gin.Context) {
	detail, err := p.ProviderService.GetProvider(c.Request.Context(), c.Param("instanceId"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(detail, c)
}
