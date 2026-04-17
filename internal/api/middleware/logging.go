// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logging returns a middleware that logs HTTP requests.
func Logging(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		entry := logger.WithFields(logrus.Fields{
			"request_id": c.GetString("request_id"),
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"status":     status,
			"latency":    latency.String(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		if status >= 500 {
			entry.Error("Server error")
		} else if status >= 400 {
			entry.Warn("Client error")
		} else {
			entry.Info("Request completed")
		}
	}
}
