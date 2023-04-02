package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wecoding/iam/pkg/utils"
	"k8s.io/klog/v2"
)

// requestLog write to console
func requestLog() gin.HandlerFunc {
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
		klog.InfoS("request log",
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
