// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coding-hui/iam/pkg/api"
)

// Handler handles HTTP requests for session operations.
type Handler struct {
	manager Manager
}

// NewHandler creates a new session handler.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager: manager}
}

// Create handles POST /api/v1/sessions.
func (h *Handler) Create(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	if req.TTL == 0 {
		req.TTL = 24 * time.Hour
	}

	sess, err := h.manager.CreateSession(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(sess, c)
}

// Get handles GET /api/v1/sessions/:id.
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	sess, err := h.manager.GetSession(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			api.FailWithMessage(err.Error(), c)
			return
		}
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(sess, c)
}

// List handles GET /api/v1/sessions.
func (h *Handler) List(c *gin.Context) {
	identityIDStr := c.GetString("identity_id")
	if identityIDStr == "" {
		api.FailWithMessage("invalid identity_id", c)
		return
	}
	identityID, err := uuid.Parse(identityIDStr)
	if err != nil {
		api.FailWithMessage("invalid identity_id", c)
		return
	}

	sessions, err := h.manager.ListSessions(c.Request.Context(), identityID)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(sessions, int64(len(sessions)), c)
}

// Revoke handles DELETE /api/v1/sessions/:id.
func (h *Handler) Revoke(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	if err := h.manager.RevokeSession(c.Request.Context(), id); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// RevokeAll handles DELETE /api/v1/sessions.
func (h *Handler) RevokeAll(c *gin.Context) {
	identityIDStr := c.GetString("identity_id")
	identityID, err := uuid.Parse(identityIDStr)
	if err != nil {
		api.FailWithMessage("invalid identity_id", c)
		return
	}

	if err := h.manager.RevokeAllSessions(c.Request.Context(), identityID); err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

// Extend handles PATCH /api/v1/sessions/:id.
func (h *Handler) Extend(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.FailWithMessage("invalid id", c)
		return
	}

	var req struct {
		Extend time.Duration `json:"extend"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithMessage("invalid request: "+err.Error(), c)
		return
	}

	if req.Extend == 0 {
		req.Extend = 24 * time.Hour
	}

	sess, err := h.manager.ExtendSession(c.Request.Context(), id, req.Extend)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(sess, c)
}
