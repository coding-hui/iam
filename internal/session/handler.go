// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TTL == 0 {
		req.TTL = 24 * time.Hour
	}

	sess, err := h.manager.CreateSession(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sess)
}

// Get handles GET /api/v1/sessions/:id.
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sess, err := h.manager.GetSession(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sess)
}

// List handles GET /api/v1/sessions.
func (h *Handler) List(c *gin.Context) {
	identityIDStr := c.GetString("identity_id")
	identityID, err := uuid.Parse(identityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id"})
		return
	}

	sessions, err := h.manager.ListSessions(c.Request.Context(), identityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": sessions,
		"total": len(sessions),
	})
}

// Revoke handles DELETE /api/v1/sessions/:id.
func (h *Handler) Revoke(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.manager.RevokeSession(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RevokeAll handles DELETE /api/v1/sessions.
func (h *Handler) RevokeAll(c *gin.Context) {
	identityIDStr := c.GetString("identity_id")
	identityID, err := uuid.Parse(identityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identity_id"})
		return
	}

	if err := h.manager.RevokeAllSessions(c.Request.Context(), identityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Extend handles PATCH /api/v1/sessions/:id.
func (h *Handler) Extend(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Extend time.Duration `json:"extend"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Extend == 0 {
		req.Extend = 24 * time.Hour
	}

	sess, err := h.manager.ExtendSession(c.Request.Context(), id, req.Extend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sess)
}
