// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package courier

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for courier operations.
type Handler struct {
	courier Courier
}

// NewHandler creates a new courier handler.
func NewHandler(courier Courier) *Handler {
	return &Handler{courier: courier}
}

// Register handles POST /api/v1/webhooks.
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		EventType string `json:"event_type"`
		TargetURL string `json:"target_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.courier.RegisterHook(c.Request.Context(), req.EventType, req.TargetURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"registered": true})
}

// Unregister handles DELETE /api/v1/webhooks.
func (h *Handler) Unregister(c *gin.Context) {
	var req struct {
		EventType string `json:"event_type"`
		TargetURL string `json:"target_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.courier.UnregisterHook(c.Request.Context(), req.EventType, req.TargetURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unregistered": true})
}

// SendEvent handles POST /api/v1/webhooks/events.
func (h *Handler) SendEvent(c *gin.Context) {
	var req struct {
		EventType string `json:"event_type"`
		Payload   any    `json:"payload"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.courier.SendEvent(c.Request.Context(), req.EventType, req.Payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"sent": true})
}
