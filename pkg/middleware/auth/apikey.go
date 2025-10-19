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
		var apiKey string

		// Try to get API Key from Authorization header first
		authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(authHeader) == authHeaderCount && authHeader[0] == authHeaderBearer {
			apiKey = authHeader[1]
		}

		// If not found in Authorization header, try X-API-Key header
		if apiKey == "" {
			apiKey = c.Request.Header.Get("X-API-Key")
		}

		// Validate API Key
		if apiKey == "" {
			api.FailWithErrCode(errors.WithCode(code.ErrApiKeyInvalid, "Missing API Key. Use Authorization: Bearer {apiKey} or X-API-Key header"), c)
			c.Abort()
			return
		}

		// Validate API Key
		user, _, err := a.apiKeyService.ValidateApiKey(c.Request.Context(), apiKey)
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
