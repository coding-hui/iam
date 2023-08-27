// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/middleware"
	"github.com/coding-hui/iam/internal/pkg/token"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

// JWTStrategy defines jwt bearer authentication strategy.
type JWTStrategy struct {
	signedKey string
}

var _ middleware.AuthStrategy = &JWTStrategy{}

// NewJWTStrategy create jwt bearer strategy with GinJWTMiddleware.
func NewJWTStrategy(signedKey string) JWTStrategy {
	return JWTStrategy{signedKey: signedKey}
}

// AuthFunc defines jwt bearer strategy as the gin authentication middleware.
func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenValue string
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader != "" {
			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				api.FailWithErrCode(errors.WithCode(code.ErrMissingHeader, "The Authorization header was empty"), c)
				c.Abort()
				return
			}
			tokenValue = splitted[1]
		}
		if tokenValue == "" {
			api.FailWithErrCode(errors.WithCode(code.ErrMissingHeader, "The Authorization header was empty"), c)
			c.Abort()
			return
		}
		t, err := token.ParseToken(tokenValue, j.signedKey)
		if err != nil {
			api.FailWithErrCode(err, c)
			c.Abort()
			return
		}
		if t.GrantType != service.GrantTypeAccess {
			api.FailWithErrCode(errors.WithCode(code.ErrPermissionDenied, "Invalid authorization header"), c)
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1.CtxKeyUserInstanceId, t.UserInstanceId))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1.CtxKeyUserType, t.UserType))

		c.Next()
	}
}
