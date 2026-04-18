// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lockout

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for lockout operations.
type Handler struct {
	manager Manager
}

// NewHandler creates a new lockout handler.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager: manager}
}

// CheckStatus handles GET /api/v1/lockout/:identifier.
func (h *Handler) CheckStatus(c *gin.Context) {
	identifier := c.Param("identifier")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identifier is required"})
		return
	}

	locked, remaining, err := h.manager.IsLocked(c.Request.Context(), identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"locked":    locked,
		"remaining": remaining.Seconds(),
	})
}

// Unlock handles POST /api/v1/lockout/:identifier/unlock.
func (h *Handler) Unlock(c *gin.Context) {
	identifier := c.Param("identifier")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identifier is required"})
		return
	}

	if err := h.manager.Unlock(c.Request.Context(), identifier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unlocked": true})
}
