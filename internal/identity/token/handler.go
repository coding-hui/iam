// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := h.manager.CreateToken(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

// Introspect handles POST /api/v1/tokens/introspect.
func (h *Handler) Introspect(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := h.manager.IntrospectToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"active": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active":      true,
		"token_id":    t.ID,
		"identity_id": t.IdentityID,
		"token_type":  t.Type,
		"expires_at":  t.ExpiresAt,
	})
}

// Revoke handles DELETE /api/v1/tokens/:id.
func (h *Handler) Revoke(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.manager.RevokeToken(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
