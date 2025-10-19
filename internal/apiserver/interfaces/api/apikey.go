// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/pkg/api"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

// apiKeyHandler create a apiKey handler used to handle request for apiKey resource.
type apiKey struct {
	ApiKeyService service.ApiKeyService `inject:""`
}

// NewApiKey is the constructor of apiKeyHandler.
func NewApiKey() Interface {
	return &apiKey{}
}

// RegisterRouter register apiKey router to gin router group.
func (a *apiKey) RegisterRouter(g *gin.RouterGroup) {
	apiv1 := g.Group("/v1/apikeys")
	{
		apiv1.POST("", a.createApiKey)
		apiv1.PUT("/:instanceId", a.updateApiKey)
		apiv1.DELETE("/:instanceId", a.deleteApiKey)
		apiv1.GET("/:instanceId", a.getApiKey)
		apiv1.GET("", a.listApiKeys)

		apiv1.PUT("/:instanceId/enable", a.enableApiKey)
		apiv1.PUT("/:instanceId/disable", a.disableApiKey)
		apiv1.POST("/:instanceId/regenerate", a.regenerateSecret)
	}
}

// RegisterApiGroup register apiKey router to gin engine.
func (a *apiKey) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix + "/apikeys")
	a.RegisterRouter(apiv1)
}

//	@Tags			ApiKeys
//	@Summary		RegenerateApiKeySecret
//	@Description	Regenerate API Key secret
//	@Accept			application/json
//	@Product		application/json
//	@Param			instanceId	path		string			true	"identifier of an API Key"
//	@Success		200			{object}	api.Response	"regenerate API Key secret"
//	@Router			/api/v1/apikeys/{instanceId}/regenerate [post]
//	@Security		BearerTokenAuth
//
// regenerateSecret regenerate API Key secret.
func (a *apiKey) regenerateSecret(c *gin.Context) {
	apiKeyResp, err := a.ApiKeyService.RegenerateSecret(c.Request.Context(), c.Param("instanceId"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(apiKeyResp, c)
}

//	@Tags			ApiKeys
//	@Summary		CreateApiKey
//	@Description	Create a new API Key
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreateApiKeyRequest	true	"API Key info"
//	@Success		200		{object}	api.Response			"create a new API Key"
//	@Router			/api/v1/apikeys [post]
//	@Security		BearerTokenAuth
//
// createApiKey create a new API Key.
func (a *apiKey) createApiKey(c *gin.Context) {
	createReq := v1.CreateApiKeyRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}

	apiKeyResp, err := a.ApiKeyService.CreateApiKey(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(apiKeyResp, c)
}

//	@Tags			ApiKeys
//	@Summary		UpdateApiKey
//	@Description	Update API Key
//	@Accept			application/json
//	@Product		application/json
//	@Param			instanceId	path		string					true	"identifier of an API Key"
//	@Param			data		body		v1.UpdateApiKeyRequest	true	"API Key info"
//	@Success		200			{object}	api.Response			"update API Key"
//	@Router			/api/v1/apikeys/{instanceId} [put]
//	@Security		BearerTokenAuth
//
// updateApiKey update API Key.
func (a *apiKey) updateApiKey(c *gin.Context) {
	updateReq := v1.UpdateApiKeyRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}

	apiKeyResp, err := a.ApiKeyService.UpdateApiKey(c.Request.Context(), c.Param("instanceId"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(apiKeyResp, c)
}

//	@Tags			ApiKeys
//	@Summary		DeleteApiKey
//	@Description	Delete API Key
//	@Accept			application/json
//	@Product		application/json
//	@Param			instanceId	path		string			true	"identifier of an API Key"
//	@Success		200			{object}	api.Response	"delete API Key"
//	@Router			/api/v1/apikeys/{instanceId} [delete]
//	@Security		BearerTokenAuth
//
// deleteApiKey delete API Key.
func (a *apiKey) deleteApiKey(c *gin.Context) {
	err := a.ApiKeyService.DeleteApiKey(c.Request.Context(), c.Param("instanceId"), metav1.DeleteOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			ApiKeys
//	@Summary		GetApiKey
//	@Description	Get API Key
//	@Accept			application/json
//	@Product		application/json
//	@Param			instanceId	path		string			true	"identifier of an API Key"
//	@Success		200			{object}	api.Response	"get API Key"
//	@Router			/api/v1/apikeys/{instanceId} [get]
//	@Security		BearerTokenAuth
//
// getApiKey get API Key.
func (a *apiKey) getApiKey(c *gin.Context) {
	apiKeyResp, err := a.ApiKeyService.GetApiKey(c.Request.Context(), c.Param("instanceId"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(apiKeyResp, c)
}

//	@Tags			ApiKeys
//	@Summary		ListApiKeys
//	@Description	List API Keys
//	@Accept			application/json
//	@Product		application/json
//	@Param			query	query		v1.ListApiKeyOptions	false	"query parameters"
//	@Success		200		{object}	api.Response			"list API Keys"
//	@Router			/api/v1/apikeys [get]
//	@Security		BearerTokenAuth
//
// listApiKeys list API Keys.
func (a *apiKey) listApiKeys(c *gin.Context) {
	var opts v1.ListApiKeyOptions
	err := c.ShouldBindQuery(&opts)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}

	apiKeys, err := a.ApiKeyService.ListApiKeys(c.Request.Context(), opts)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(apiKeys, c)
}

//	@Tags			ApiKeys
//	@Summary		EnableApiKey
//	@Description	Enable API Key
//	@Accept			application/json
//	@Product		application/json
//	@Param			instanceId	path		string			true	"identifier of an API Key"
//	@Success		200			{object}	api.Response	"enable API Key"
//	@Router			/api/v1/apikeys/{instanceId}/enable [put]
//	@Security		BearerTokenAuth
//
// enableApiKey enable API Key.
func (a *apiKey) enableApiKey(c *gin.Context) {
	err := a.ApiKeyService.EnableApiKey(c.Request.Context(), c.Param("instanceId"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			ApiKeys
//	@Summary		DisableApiKey
//	@Description	Disable API Key
//	@Accept			application/json
//	@Product		application/json
//	@Param			instanceId	path		string			true	"identifier of an API Key"
//	@Success		200			{object}	api.Response	"disable API Key"
//	@Router			/api/v1/apikeys/{instanceId}/disable [put]
//	@Security		BearerTokenAuth
//
// disableApiKey disable API Key.
func (a *apiKey) disableApiKey(c *gin.Context) {
	err := a.ApiKeyService.DisableApiKey(c.Request.Context(), c.Param("instanceId"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}
