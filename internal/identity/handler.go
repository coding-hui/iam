// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identity

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	identity, err := h.manager.CreateIdentity(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, identity)
}

// Get handles GET /api/v1/identities/:id.
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	identity, err := h.manager.GetIdentity(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrIdentityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, identity)
}

// List handles GET /api/v1/identities.
func (h *Handler) List(c *gin.Context) {
	var params ListIdentitiesParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	networkID := c.GetString("network_id")

	identities, total, err := h.manager.ListIdentities(c.Request.Context(), uuid.MustParse(networkID), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": identities,
		"total": total,
	})
}

// Update handles PATCH /api/v1/identities/:id.
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	identity, err := h.manager.UpdateIdentity(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, identity)
}

// Delete handles DELETE /api/v1/identities/:id.
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.manager.DeleteIdentity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddCredentials handles POST /api/v1/identities/:id/credentials.
func (h *Handler) AddCredentials(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req AddCredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.manager.AddCredentials(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// DeleteCredentials handles DELETE /api/v1/identities/:id/credentials/:type.
func (h *Handler) DeleteCredentials(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	credType := CredentialsType(c.Param("type"))
	if credType != CredentialsTypePassword && credType != CredentialsTypeAPIKey && credType != CredentialsTypeTOTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials type"})
		return
	}

	if err := h.manager.DeleteCredentials(c.Request.Context(), id, credType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
