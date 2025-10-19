// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/pkg/api"
	"github.com/coding-hui/iam/pkg/code"
	"github.com/coding-hui/iam/pkg/middleware"
)

const (
	apiKeyHeader = "X-API-Key"
)

// AutoStrategy defines authentication strategy which can automatically choose between Basic, Bearer, and API Key
// according `Authorization` header or API Key headers.
type AutoStrategy struct {
	basic  middleware.AuthStrategy
	jwt    middleware.AuthStrategy
	apiKey middleware.AuthStrategy
}

var _ middleware.AuthStrategy = &AutoStrategy{}

// NewAutoStrategy create auto strategy with basic strategy and jwt strategy.
func NewAutoStrategy(basic, jwt, apiKey middleware.AuthStrategy) AutoStrategy {
	return AutoStrategy{
		basic:  basic,
		jwt:    jwt,
		apiKey: apiKey,
	}
}

// AuthFunc defines auto strategy as the gin authentication middleware.
func (a AutoStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		operator := middleware.AuthOperator{}

		// Check for API Key authentication first (via headers)
		apiKey := c.Request.Header.Get(apiKeyHeader)

		if apiKey != "" {
			operator.SetStrategy(a.apiKey)
			operator.AuthFunc()(c)
			c.Next()
			return
		}

		// Check for Authorization header
		authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(authHeader) != authHeaderCount {
			api.FailWithErrCode(errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong."), c)
			c.Abort()

			return
		}

		switch authHeader[0] {
		case "Basic":
			operator.SetStrategy(a.basic)
		case authHeaderBearer:
			// Check if Bearer token is an API Key or JWT
			token := authHeader[1]

			// API Key format: sk-{32 hex chars} or any string with colon separator
			if strings.HasPrefix(token, "sk-") || strings.Contains(token, ":") {
				// This is likely an API Key authentication
				operator.SetStrategy(a.apiKey)
			} else {
				// Check if it's a valid JWT (3 parts separated by dots)
				jwtParts := strings.Split(token, ".")
				if len(jwtParts) == 3 {
					// This is a JWT token
					operator.SetStrategy(a.jwt)
				} else {
					// Try API Key strategy as fallback
					operator.SetStrategy(a.apiKey)
				}
			}
		default:
			api.FailWithErrCode(errors.WithCode(code.ErrSignatureInvalid, "unrecognized Authorization header."), c)
			c.Abort()

			return
		}

		operator.AuthFunc()(c)

		c.Next()
	}
}
