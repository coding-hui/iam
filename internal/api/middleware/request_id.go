// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	// XRequestIDKey defines X-Request-ID key string.
	XRequestIDKey = "X-Request-ID"
)

// requestIDPattern validates that incoming X-Request-ID headers contain only safe characters.
var requestIDPattern = regexp.MustCompile(`^[a-zA-Z0-9\-_]{1,64}$`)

// RequestID is a middleware that injects a 'X-Request-ID' into the context and request/response header of each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(XRequestIDKey)

		if rid == "" || !requestIDPattern.MatchString(rid) {
			rid = uuid.Must(uuid.NewV4()).String()
		}

		c.Request.Header.Set(XRequestIDKey, rid)
		c.Set(XRequestIDKey, rid)
		c.Writer.Header().Set(XRequestIDKey, rid)
		c.Next()
	}
}

// GinLogFormatterParams holds the parameters for GIN JSON log format.
type GinLogFormatterParams struct {
	gin.LogFormatterParams
	RequestID string
}

// GetLoggerConfig return gin.LoggerConfig with JSON formatter.
func GetLoggerConfig(output io.Writer, skipPaths []string) gin.LoggerConfig {
	return gin.LoggerConfig{
		Formatter: GetGinJSONLogFormatter(),
		Output:    output,
		SkipPaths: skipPaths,
	}
}

// GetGinJSONLogFormatter returns a gin.LogFormatter that outputs JSON.
func GetGinJSONLogFormatter() gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		rid := ""
		if v, ok := param.Keys[XRequestIDKey]; ok {
			rid, _ = v.(string)
		}

		if param.Latency > time.Minute {
			param.Latency -= param.Latency % time.Second
		}

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)

		logEntry := map[string]any{
			"time":      param.TimeStamp.Format(time.RFC3339Nano),
			"level":     "info",
			"client_ip": param.ClientIP,
			"latency":   param.Latency.String(),
			"method":    param.Method,
			"path":      param.Path,
			"status":    param.StatusCode,
			"error":     param.ErrorMessage,
		}

		if rid != "" {
			logEntry["request_id"] = rid
		}

		if err := enc.Encode(logEntry); err != nil {
			return ""
		}
		return buf.String()
	}
}

// GetRequestIDFromContext returns 'RequestID' from the given context if present.
func GetRequestIDFromContext(c *gin.Context) string {
	if v, ok := c.Get(XRequestIDKey); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}

	return ""
}

// GetRequestIDFromHeaders returns 'RequestID' from the headers if present.
func GetRequestIDFromHeaders(c *gin.Context) string {
	return c.Request.Header.Get(XRequestIDKey)
}
