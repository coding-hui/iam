// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package selfservice

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/selfservice/strategies"
	"github.com/coding-hui/iam/pkg/api"
)

// Handler handles HTTP requests for selfservice operations.
type Handler struct {
	passwordAuthenticator *strategies.PasswordAuthenticator
	mfaManager            *strategies.ManagerImpl
}

// NewHandler creates a new selfservice handler.
func NewHandler(passwordAuthenticator *strategies.PasswordAuthenticator, mfaManager *strategies.ManagerImpl) *Handler {
	return &Handler{
		passwordAuthenticator: passwordAuthenticator,
		mfaManager:            mfaManager,
	}
}

// Login handles POST /api/v1/login.
func (h *Handler) Login(c *gin.Context) {
	var req strategies.AuthenticateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	req.UserAgent = c.GetHeader("User-Agent")
	req.ClientIP = c.ClientIP()

	resp, err := h.passwordAuthenticator.Authenticate(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(resp, c)
}

// SetupTOTP handles POST /api/v1/mfa/totp/setup.
func (h *Handler) SetupTOTP(c *gin.Context) {
	// TODO: implement
	api.FailWithMessage("not implemented", c)
}

// VerifyTOTP handles POST /api/v1/mfa/totp/verify.
func (h *Handler) VerifyTOTP(c *gin.Context) {
	// TODO: implement
	api.FailWithMessage("not implemented", c)
}

// DisableTOTP handles POST /api/v1/mfa/totp/disable.
func (h *Handler) DisableTOTP(c *gin.Context) {
	// TODO: implement
	api.FailWithMessage("not implemented", c)
}
