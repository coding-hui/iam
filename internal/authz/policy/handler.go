// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coding-hui/iam/pkg/api"
)

// Handler handles HTTP requests for policy operations.
type Handler struct {
	manager Manager
}

// NewHandler creates a new policy handler.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager: manager}
}

// Create handles POST /api/v1/policies.
func (h *Handler) Create(c *gin.Context) {
	var req CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	p, err := h.manager.CreatePolicy(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(p, c)
}

// Get handles GET /api/v1/policies/:id.
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	p, err := h.manager.GetPolicy(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPolicyNotFound) {
			api.FailWithMessage(err.Error(), c)
			return
		}
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(p, c)
}

// List handles GET /api/v1/policies.
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

	policies, err := h.manager.ListPolicies(c.Request.Context(), networkID)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(policies, int64(len(policies)), c)
}

// Update handles PATCH /api/v1/policies/:id.
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	var req UpdatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	p, err := h.manager.UpdatePolicy(c.Request.Context(), id, &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(p, c)
}

// Delete handles DELETE /api/v1/policies/:id.
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	if err := h.manager.DeletePolicy(c.Request.Context(), id); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}
