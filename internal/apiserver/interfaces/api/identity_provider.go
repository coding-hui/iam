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

var idpApiPath = versionPrefix + "/identity_providers"

type identityProvider struct {
	IdentityProviderService service.IdentityProviderService `inject:""`
}

// NewIdentityProvider is the of IdentityProvider.
func NewIdentityProvider() Interface {
	return &identityProvider{}
}

func (p *identityProvider) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(idpApiPath).
		Use(
			autoAuthCheck.AuthFunc(),
			permissionCheckFunc("identity_providers"),
		)
	{
		apiv1.POST("", p.createIdentityProvider)
		apiv1.PUT("/:identifier", p.updateIdentityProvider)
		apiv1.DELETE("/:identifier", p.deleteIdentityProvider)
		apiv1.GET("/:identifier", p.detailIdentityProvider)
	}
}

//	@Tags			Providers
//	@Summary		CreateIdentityProvider
//	@Description	Create a IdentityProvider
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreateProviderRequest	true	"IdentityProvider"
//	@Success		200		{object}	api.Response				"Create a nnw IdentityProvider"
//	@Router			/api/v1/providers [post]
//	@Security		BearerTokenAuth
//
// createIdentityProvider create a new IdentityProvider.
func (p *identityProvider) createIdentityProvider(c *gin.Context) {
	createReq := v1.CreateIdentityProviderRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = p.IdentityProviderService.CreateIdentityProvider(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Providers
//	@Summary		UpdateIdentityProvider
//	@Description	Update a IdentityProvider
//	@Accept			application/json
//	@Product		application/json
//	@Param			data		body		v1.UpdateProviderRequest	true	"IdentityProvider"
//	@Param			instanceId	path		string						true	"identifier of a IdentityProvider"
//	@Success		200			{object}	api.Response				"Update IdentityProvider info"
//	@Router			/api/v1/providers/{identifier} [put]
//	@Security		BearerTokenAuth
//
// updateIdentityProvider update IdentityProvider info.
func (p *identityProvider) updateIdentityProvider(c *gin.Context) {
	updateReq := v1.UpdateIdentityProviderRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = p.IdentityProviderService.UpdateIdentityProvider(c.Request.Context(), c.Param("identifier"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Providers
//	@Summary		DeleteIdentityProvider
//	@Description	Delete IdentityProvider by identifier
//	@Param			instanceId	path		string			true	"identifier of a IdentityProvider"
//	@Success		200			{object}	api.Response	"Provider successfully deleted"
//	@Router			/api/v1/providers/{identifier} [DELETE]
//	@Security		BearerTokenAuth
//
// deleteIdentityProvider delete IdentityProvider by identifier.
func (p *identityProvider) deleteIdentityProvider(c *gin.Context) {
	err := p.IdentityProviderService.DeleteIdentityProvider(
		c.Request.Context(),
		c.Param("identifier"),
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
//	@Description	Get a IdentityProvider by name
//	@Param			instanceId	path		string											true	"identifier of a IdentityProvider"
//	@Success		200			{object}	api.Response{data=v1.DetailProviderResponse}	"Provider detail"
//	@Router			/api/v1/providers/{identifier} [get]
//	@Security		BearerTokenAuth
//
// detailIdentityProvider get IdentityProvider detail info.
func (p *identityProvider) detailIdentityProvider(c *gin.Context) {
	idp, err := p.IdentityProviderService.GetIdentityProvider(c.Request.Context(), c.Param("identifier"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	detail, err := p.IdentityProviderService.DetailIdentityProvider(c.Request.Context(), idp, metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(detail, c)
}
