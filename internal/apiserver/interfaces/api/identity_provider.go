// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/pkg/api"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
)

var idpApiPath = versionPrefix + "/identity-providers"

type identityProvider struct {
	IdentityProviderService service.IdentityProviderService `inject:""`
}

// NewIdentityProvider is the of IdentityProvider.
func NewIdentityProvider() Interface {
	return &identityProvider{}
}

func (i *identityProvider) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(idpApiPath).
		Use(
			autoAuthCheck.AuthFunc(),
			permissionCheckFunc("identity-providers"),
		)
	{
		apiv1.POST("", i.createIdentityProvider)
		apiv1.PUT("/:identifier", i.updateIdentityProvider)
		apiv1.DELETE("/:identifier", i.deleteIdentityProvider)
		apiv1.GET("/:identifier", i.detailIdentityProvider)
		apiv1.GET("", i.listIdentityProviders)
	}
}

//	@Tags			IdentityProviders
//	@Summary		CreateIdentityProvider
//	@Description	Create a IdentityProvider
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreateProviderRequest	true	"IdentityProvider"
//	@Success		200		{object}	api.Response				"Create a nnw IdentityProvider"
//	@Router			/api/v1/identity-providers [post]
//	@Security		BearerTokenAuth
//
// createIdentityProvider create a new IdentityProvider.
func (i *identityProvider) createIdentityProvider(c *gin.Context) {
	createReq := v1.CreateIdentityProviderRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, "%s", errs.ToAggregate().Error()), c)
		return
	}
	err = i.IdentityProviderService.CreateIdentityProvider(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			IdentityProviders
//	@Summary		UpdateIdentityProvider
//	@Description	Update a IdentityProvider
//	@Accept			application/json
//	@Product		application/json
//	@Param			data		body		v1.UpdateProviderRequest	true	"IdentityProvider"
//	@Param			instanceId	path		string						true	"identifier of a IdentityProvider"
//	@Success		200			{object}	api.Response				"Update IdentityProvider info"
//	@Router			/api/v1/identity-providers/{identifier} [put]
//	@Security		BearerTokenAuth
//
// updateIdentityProvider update IdentityProvider info.
func (i *identityProvider) updateIdentityProvider(c *gin.Context) {
	updateReq := v1.UpdateIdentityProviderRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, "%s", errs.ToAggregate().Error()), c)
		return
	}
	err = i.IdentityProviderService.UpdateIdentityProvider(c.Request.Context(), c.Param("identifier"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			IdentityProviders
//	@Summary		DeleteIdentityProvider
//	@Description	Delete IdentityProvider by identifier
//	@Param			instanceId	path		string			true	"identifier of a IdentityProvider"
//	@Success		200			{object}	api.Response	"Provider successfully deleted"
//	@Router			/api/v1/identity-providers/{identifier} [DELETE]
//	@Security		BearerTokenAuth
//
// deleteIdentityProvider delete IdentityProvider by identifier.
func (i *identityProvider) deleteIdentityProvider(c *gin.Context) {
	err := i.IdentityProviderService.DeleteIdentityProvider(
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

//	@Tags			IdentityProviders
//	@Summary		GetProviderInfo
//	@Description	Get a IdentityProvider by name
//	@Param			instanceId	path		string											true	"identifier of a IdentityProvider"
//	@Success		200			{object}	api.Response{data=v1.DetailProviderResponse}	"Provider detail"
//	@Router			/api/v1/identity-providers/{identifier} [get]
//	@Security		BearerTokenAuth
//
// detailIdentityProvider get IdentityProvider detail info.
func (i *identityProvider) detailIdentityProvider(c *gin.Context) {
	idp, err := i.IdentityProviderService.GetIdentityProvider(c.Request.Context(), c.Param("identifier"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	detail, err := i.IdentityProviderService.DetailIdentityProvider(c.Request.Context(), idp, metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(detail, c)
}

//	@Tags			IdentityProviders
//	@Summary		ListIdentityProviders
//	@Description	List IdentityProviders
//	@Param			offset	query		int											false	"query the page number"
//	@Param			limit	query		int											false	"query the page size number"
//	@Success		200		{object}	api.Response{data=v1.IdentityProviderList}	"IdentityProviders"
//	@Router			/api/v1/identity-providers [get]
//	@Security		BearerTokenAuth
//
// listIdentityProviders list IdentityProvider page.
func (i *identityProvider) listIdentityProviders(c *gin.Context) {
	var opts metav1.ListOptions
	err := c.ShouldBindQuery(&opts)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	resp, err := i.IdentityProviderService.ListIdentityProviders(c.Request.Context(), opts)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}
