// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/pkg/request"
	"github.com/coding-hui/iam/pkg/api"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
	"github.com/coding-hui/iam/pkg/middleware"

	"github.com/coding-hui/common/errors"
)

const (
	authHeaderCount  = 2
	authHeaderBearer = "Bearer"
)

// ApiKeyStrategy defines authentication strategy for API Key authentication.
type ApiKeyStrategy struct {
	apiKeyService service.ApiKeyService
}

var _ middleware.AuthStrategy = &ApiKeyStrategy{}

// NewApiKeyStrategy create API Key strategy.
func NewApiKeyStrategy(apiKeyService service.ApiKeyService) ApiKeyStrategy {
	return ApiKeyStrategy{
		apiKeyService: apiKeyService,
	}
}

// AuthFunc defines API Key strategy as the gin authentication middleware.
func (a ApiKeyStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiKey, apiSecret string

		// Try to get API Key from Authorization header first
		authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(authHeader) == authHeaderCount && authHeader[0] == authHeaderBearer {
			// Extract API Key and Secret from Bearer token
			// Format: Bearer {apiKey}:{apiSecret}
			credentials := strings.SplitN(authHeader[1], ":", 2)
			if len(credentials) == 2 {
				apiKey = credentials[0]
				apiSecret = credentials[1]
			} else {
				// If no colon separator, check if it's a JWT token or malformed API Key
				// JWT tokens typically have 3 parts separated by dots
				jwtParts := strings.Split(authHeader[1], ".")
				if len(jwtParts) != 3 {
					// This might be a malformed API Key, try to use it as API Key only
					// and look for secret in X-API-Secret header
					apiKey = authHeader[1]
					apiSecret = c.Request.Header.Get("X-API-Secret")
				}
			}
		}

		// If not found in Authorization header, try X-API-Key and X-API-Secret headers
		if apiKey == "" {
			apiKey = c.Request.Header.Get("X-API-Key")
			apiSecret = c.Request.Header.Get("X-API-Secret")
		}

		// Validate API Key and Secret
		if apiKey == "" || apiSecret == "" {
			errMsg := "Missing API Key or Secret. "
			if authHeader[0] == authHeaderBearer {
				errMsg += "Expected format: Bearer {apiKey}:{apiSecret}"
			} else {
				errMsg += "Use X-API-Key and X-API-Secret headers or Authorization: Bearer {apiKey}:{apiSecret}"
			}
			api.FailWithErrCode(errors.WithCode(code.ErrApiKeyInvalid, "%s", errMsg), c)
			c.Abort()
			return
		}

		// Validate API Key
		user, _, err := a.apiKeyService.ValidateApiKey(c.Request.Context(), apiKey, apiSecret)
		if err != nil {
			api.FailWithErrCode(err, c)
			c.Abort()
			return
		}

		// Set user information in context
		c.Request = c.Request.WithContext(request.WithUser(c.Request.Context(), v1.UserBase{
			ObjectMeta: user.ObjectMeta,
			Status:     user.Status,
			Alias:      user.Alias,
			Email:      user.Email,
			Phone:      user.Phone,
			Avatar:     user.Avatar,
			UserType:   user.UserType,
			Disabled:   user.Disabled,
		}))
	}
}
