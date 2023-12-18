// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import "github.com/gin-gonic/gin"

// Defines the key in gin context which represents the owner of the secret.
const (
	UsernameKey  string = "username"
	RequestIDKey string = "requestID"
)

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(RequestIDKey, c.GetString(XRequestIDKey))
		c.Set(UsernameKey, c.GetString(UsernameKey))
		c.Next()
	}
}
