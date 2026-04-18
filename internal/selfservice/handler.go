// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package selfservice

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/selfservice/strategies"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserAgent = c.GetHeader("User-Agent")
	req.ClientIP = c.ClientIP()

	resp, err := h.passwordAuthenticator.Authenticate(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":  resp.SessionID,
		"identity_id": resp.IdentityID,
		"expires_at":  resp.ExpiresAt,
	})
}

// SetupTOTP handles POST /api/v1/mfa/totp/setup.
func (h *Handler) SetupTOTP(c *gin.Context) {
	// TODO: implement
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// VerifyTOTP handles POST /api/v1/mfa/totp/verify.
func (h *Handler) VerifyTOTP(c *gin.Context) {
	// TODO: implement
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// DisableTOTP handles POST /api/v1/mfa/totp/disable.
func (h *Handler) DisableTOTP(c *gin.Context) {
	// TODO: implement
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
