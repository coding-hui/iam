// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

// Middlewares store registered middlewares.
var Middlewares = defaultMiddlewares()

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":  GinRecovery,
		"secure":    Secure,
		"options":   Options,
		"nocache":   NoCache,
		"requestid": RequestID(),
		"dump":      gindump.Dump(),
		"cors":      Cors(),
		"logger":    RequestLog(),
	}
}
