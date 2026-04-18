// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coding-hui/iam/pkg/api"
)

// Handler handles HTTP requests for role operations.
type Handler struct {
	manager Manager
}

// NewHandler creates a new role handler.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager: manager}
}

// Create handles POST /api/v1/roles.
func (h *Handler) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	r, err := h.manager.CreateRole(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(r, c)
}

// Get handles GET /api/v1/roles/:id.
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	r, err := h.manager.GetRole(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			api.FailWithMessage(err.Error(), c)
			return
		}
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(r, c)
}

// List handles GET /api/v1/roles.
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

	roles, err := h.manager.ListRoles(c.Request.Context(), networkID)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(roles, int64(len(roles)), c)
}

// Update handles PATCH /api/v1/roles/:id.
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	r, err := h.manager.UpdateRole(c.Request.Context(), id, &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(r, c)
}

// Delete handles DELETE /api/v1/roles/:id.
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	if err := h.manager.DeleteRole(c.Request.Context(), id); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}
