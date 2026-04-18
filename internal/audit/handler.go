// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package audit

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coding-hui/iam/pkg/api"
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
	if networkIDStr == "" {
		networkIDStr = "00000000-0000-0000-0000-000000000000"
	}
	networkID, err := uuid.Parse(networkIDStr)
	if err != nil {
		api.FailWithMessage("invalid network_id", c)
		return
	}

	var req struct {
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
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
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(events, int64(total), c)
}
