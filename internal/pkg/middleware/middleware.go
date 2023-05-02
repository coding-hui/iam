// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import "github.com/gin-gonic/gin"

// Middlewares store registered middlewares.
var Middlewares = defaultMiddlewares()

// InitMiddleware initialize middleware.
func InitMiddleware(r *gin.Engine) {
	// Custom Error Wrapper
	r.Use(GinRecovery)
	// NoCache is a middleware function that appends headers
	r.Use(NoCache)
	// 跨域处理
	r.Use(Options)
	// Secure is a middleware function that appends security
	r.Use(Secure)
	// request log
	r.Use(RequestLog())
}

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":  GinRecovery,
		"secure":    Secure,
		"options":   Options,
		"nocache":   NoCache,
		"requestid": RequestID(),
	}
}
