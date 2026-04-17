// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Write writes an error response.
func Write(c *gin.Context, err error) {
	WriteWithCode(c, http.StatusInternalServerError, err.Error())
}

// WriteWithCode writes an error response with a specific HTTP status code.
func WriteWithCode(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// BadRequest writes a 400 error response.
func BadRequest(c *gin.Context, message string) {
	WriteWithCode(c, http.StatusBadRequest, message)
}

// Unauthorized writes a 401 error response.
func Unauthorized(c *gin.Context, message string) {
	WriteWithCode(c, http.StatusUnauthorized, message)
}

// Forbidden writes a 403 error response.
func Forbidden(c *gin.Context, message string) {
	WriteWithCode(c, http.StatusForbidden, message)
}

// NotFound writes a 404 error response.
func NotFound(c *gin.Context, message string) {
	WriteWithCode(c, http.StatusNotFound, message)
}

// Conflict writes a 409 error response.
func Conflict(c *gin.Context, message string) {
	WriteWithCode(c, http.StatusConflict, message)
}

// InternalServerError writes a 500 error response.
func InternalServerError(c *gin.Context, message string) {
	WriteWithCode(c, http.StatusInternalServerError, message)
}
