// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identity

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coding-hui/iam/pkg/api"
)

// Handler handles HTTP requests for identity operations.
type Handler struct {
	manager Manager
}

// NewHandler creates a new identity handler.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager: manager}
}

// Create handles POST /api/v1/identities.
func (h *Handler) Create(c *gin.Context) {
	var req CreateIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	identity, err := h.manager.CreateIdentity(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(identity, c)
}

// Get handles GET /api/v1/identities/:id.
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	identity, err := h.manager.GetIdentity(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrIdentityNotFound) {
			api.FailWithMessage(err.Error(), c)
			return
		}
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(identity, c)
}

// List handles GET /api/v1/identities.
func (h *Handler) List(c *gin.Context) {
	var params ListIdentitiesParams
	if err := c.ShouldBindQuery(&params); err != nil {
		api.FailWithMessage("invalid params: "+err.Error(), c)
		return
	}

	networkIDStr := c.GetString("network_id")
	if networkIDStr == "" {
		networkIDStr = "00000000-0000-0000-0000-000000000000"
	}
	networkID := uuid.MustParse(networkIDStr)

	identities, total, err := h.manager.ListIdentities(c.Request.Context(), networkID, params)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(identities, int64(total), c)
}

// Update handles PATCH /api/v1/identities/:id.
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	var req UpdateIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	identity, err := h.manager.UpdateIdentity(c.Request.Context(), id, &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(identity, c)
}

// Delete handles DELETE /api/v1/identities/:id.
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	if err := h.manager.DeleteIdentity(c.Request.Context(), id); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// AddCredentials handles POST /api/v1/identities/:id/credentials.
func (h *Handler) AddCredentials(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	var req AddCredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	if err := h.manager.AddCredentials(c.Request.Context(), id, &req); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// DeleteCredentials handles DELETE /api/v1/identities/:id/credentials/:type.
func (h *Handler) DeleteCredentials(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	credType := CredentialsType(c.Param("type"))
	if credType != CredentialsTypePassword && credType != CredentialsTypeAPIKey && credType != CredentialsTypeTOTP {
		api.FailWithMessage("invalid credentials type", c)
		return
	}

	if err := h.manager.DeleteCredentials(c.Request.Context(), id, credType); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}
