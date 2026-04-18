// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coding-hui/iam/pkg/api"
)

// Handler handles HTTP requests for token operations.
type Handler struct {
	manager Manager
}

// NewHandler creates a new token handler.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager: manager}
}

// Create handles POST /api/v1/tokens.
func (h *Handler) Create(c *gin.Context) {
	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	t, err := h.manager.CreateToken(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(t, c)
}

// Introspect handles POST /api/v1/tokens/introspect.
func (h *Handler) Introspect(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	t, err := h.manager.IntrospectToken(c.Request.Context(), req.Token)
	if err != nil {
		api.OkWithData(gin.H{
			"active": false,
			"error":  err.Error(),
		}, c)
		return
	}

	api.OkWithData(gin.H{
		"active":      true,
		"token_id":    t.ID,
		"identity_id": t.IdentityID,
		"token_type":  t.Type,
		"expires_at":  t.ExpiresAt,
	}, c)
}

// Revoke handles DELETE /api/v1/tokens/:id.
func (h *Handler) Revoke(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	if err := h.manager.RevokeToken(c.Request.Context(), id); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}
