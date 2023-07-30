// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/pkg/utils"
	"github.com/coding-hui/iam/pkg/log"
)

// RequestLog write to console.
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if (c.Request.Header.Get("Upgrade") == "websocket" && c.Request.Header.Get("Upgrade") == "Upgrade") ||
			c.Request.Method == http.MethodOptions {
			c.Next()
		}
		start := time.Now()

		c.Next()

		// 执行时间
		takeTime := time.Since(start)
		// Query
		query := c.Request.URL.RawQuery
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := utils.ClientIP(c.Request)
		log.Infow("request log",
			"clientIP", utils.Sanitize(clientIP),
			"uri", utils.Sanitize(reqUri),
			"query", query,
			"method", reqMethod,
			"status", statusCode,
			"time", takeTime,
			"error", strings.TrimRight(c.Errors.ByType(gin.ErrorTypePrivate).String(), "\n"),
		)
	}
}
