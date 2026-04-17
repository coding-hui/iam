// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for audit operations.
type Handler struct {
	manager Manager
}

// NewHandler creates a new audit handler.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager: manager}
}

// List handles GET /api/v1/audit/events.
func (h *Handler) List(c *gin.Context) {
	networkIDStr := c.GetString("network_id")
	networkID, err := uuid.Parse(networkIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid network_id"})
		return
	}

	var req struct {
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	events, total, err := h.manager.ListEvents(c.Request.Context(), networkID, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": events,
		"total": total,
	})
}
