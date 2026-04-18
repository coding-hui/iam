// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lockout

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/pkg/api"
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
		api.FailWithMessage("identifier is required", c)
		return
	}

	locked, remaining, err := h.manager.IsLocked(c.Request.Context(), identifier)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(gin.H{
		"locked":    locked,
		"remaining": remaining.Seconds(),
	}, c)
}

// Unlock handles POST /api/v1/lockout/:identifier/unlock.
func (h *Handler) Unlock(c *gin.Context) {
	identifier := c.Param("identifier")
	if identifier == "" {
		api.FailWithMessage("identifier is required", c)
		return
	}

	if err := h.manager.Unlock(c.Request.Context(), identifier); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}
