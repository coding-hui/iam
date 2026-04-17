// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authz

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for authorization operations.
type Handler struct {
	engine *Engine
}

// NewHandler creates a new authz handler.
func NewHandler(engine *Engine) *Handler {
	return &Handler{engine: engine}
}

// Check handles POST /api/v1/authz/check.
func (h *Handler) Check(c *gin.Context) {
	var req AuthzRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.engine.Authorize(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"decision": resp.Decision,
		"reason":   resp.Reason,
	})
}
