// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authn

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for authentication operations.
type Handler struct {
	authenticator Authenticator
}

// NewHandler creates a new authn handler.
func NewHandler(authenticator Authenticator) *Handler {
	return &Handler{authenticator: authenticator}
}

// Login handles POST /api/v1/login.
func (h *Handler) Login(c *gin.Context) {
	var req AuthenticateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserAgent = c.GetHeader("User-Agent")
	req.ClientIP = c.ClientIP()

	resp, err := h.authenticator.Authenticate(c.Request.Context(), &req)
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
