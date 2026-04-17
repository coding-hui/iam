// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mfa

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for MFA operations.
type Handler struct {
	manager Manager
}

// NewHandler creates a new MFA handler.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager: manager}
}

// SetupTOTP handles POST /api/v1/mfa/totp/setup.
func (h *Handler) SetupTOTP(c *gin.Context) {
	identityIDStr := c.Param("identity_id")
	if identityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identity_id is required"})
		return
	}

	identityID, err := uuid.Parse(identityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id"})
		return
	}

	var req struct {
		Issuer  string `json:"issuer"`
		Account string `json:"account"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Issuer == "" {
		req.Issuer = "IAM"
	}

	totp, err := h.manager.SetupTOTP(c.Request.Context(), identityID, req.Issuer, req.Account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, totp)
}

// VerifyTOTP handles POST /api/v1/mfa/totp/verify.
func (h *Handler) VerifyTOTP(c *gin.Context) {
	var req struct {
		IdentityID string `json:"identity_id"`
		Code       string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	identityID, err := uuid.Parse(req.IdentityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id"})
		return
	}

	if err := h.manager.VerifyTOTP(c.Request.Context(), identityID, req.Code); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verified": true})
}

// DisableTOTP handles POST /api/v1/mfa/totp/disable.
func (h *Handler) DisableTOTP(c *gin.Context) {
	identityIDStr := c.Param("identity_id")
	if identityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identity_id is required"})
		return
	}

	identityID, err := uuid.Parse(identityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id"})
		return
	}

	if err := h.manager.DisableTOTP(c.Request.Context(), identityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"disabled": true})
}
