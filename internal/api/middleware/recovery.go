// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery returns a middleware that recovers from panics.
func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.WithFields(logrus.Fields{
					"request_id": c.GetString("request_id"),
					"error":      err,
					"stack":      string(debug.Stack()),
				}).Error("Panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
