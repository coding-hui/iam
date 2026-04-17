// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that handles CORS headers.
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = false
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"}
	config.ExposeHeaders = []string{"X-Request-ID"}
	config.AllowCredentials = true
	config.MaxAge = 86400 // 24 hours

	return cors.New(config)
}
