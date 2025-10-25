// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/pkg/api"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
)

// deviceAuth handles device authorization related API endpoints.
type deviceAuth struct {
	cfg               config.Config
	DeviceAuthService service.DeviceAuthService `inject:""`
}

// NewDeviceAuth creates a new deviceAuth instance.
func NewDeviceAuth(c config.Config) Interface {
	return &deviceAuth{cfg: c}
}

// RegisterApiGroup registers device authorization API endpoints.
func (d *deviceAuth) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix + "/device")
	{
		apiv1.POST("/code", d.createDeviceCode)
		apiv1.POST("/token", d.getDeviceToken)
		apiv1.GET("/authorize", d.deviceAuthorizePage)
		apiv1.POST("/authorize", d.verifyUserAuthorization)
	}
}

// createDeviceCode creates a device authorization code.
//
//	@Summary		Create device authorization code
//	@Description	Create device authorization code for OAuth 2.0 Device Flow
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			data	body		v1.DeviceAuthorizationRequest	true	"Device authorization request"
//	@Success		200		{object}	v1.DeviceAuthorizationResponse
//	@Router			/api/v1/device/code [post]
func (d *deviceAuth) createDeviceCode(c *gin.Context) {
	var req v1.DeviceAuthorizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}

	resp, err := d.DeviceAuthService.CreateDeviceAuthorization(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(resp, c)
}

// getDeviceToken retrieves access token using device code.
//
//	@Summary		Get device token
//	@Description	Get access token using device code
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			data	body		v1.DeviceTokenRequest	true	"Device token request"
//	@Success		200		{object}	v1.DeviceTokenResponse
//	@Router			/api/v1/device/token [post]
func (d *deviceAuth) getDeviceToken(c *gin.Context) {
	var req v1.DeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}

	resp, err := d.DeviceAuthService.GetDeviceToken(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(resp, c)
}

// deviceAuthorizePage renders the device authorization page.
//
//	@Summary		Device authorization page
//	@Description	Show device authorization page for user to enter user code
//	@Tags			Authentication
//	@Accept			html
//	@Produce		html
//	@Param			user_code	query	string	false	"User code"
//	@Router			/api/v1/device/authorize [get]
func (d *deviceAuth) deviceAuthorizePage(c *gin.Context) {
	userCode := c.Query("user_code")

	// Render device authorization page
	c.HTML(http.StatusOK, "device_authorize.html", gin.H{
		"UserCode": userCode,
		"Title":    "Device Authorization",
	})
}

// verifyUserAuthorization verifies user authorization for device flow.
//
//	@Summary		Verify user authorization
//	@Description	Verify user authorization for device flow
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			data	body		v1.VerifyDeviceRequest	true	"Verify device request"
//	@Success		200		{object}	api.Response
//	@Router			/api/v1/device/authorize [post]
func (d *deviceAuth) verifyUserAuthorization(c *gin.Context) {
	var req v1.VerifyDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}

	if err := d.DeviceAuthService.VerifyUserAuthorization(c.Request.Context(), &req); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithMessage("Authorization approved successfully", c)
}
