// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/pkg/request"
	"github.com/coding-hui/iam/internal/pkg/token"
	"github.com/coding-hui/iam/pkg/api"
	"github.com/coding-hui/iam/pkg/code"
	"github.com/coding-hui/iam/pkg/middleware"
)

// AuthzAudience defines the value of jwt audience field.
const AuthzAudience = "iam.api.wecoding.top"

// JWTStrategy defines jwt bearer authentication strategy.
type JWTStrategy struct {
	verify func(tokenStr string) (*token.VerifiedResponse, error)
}

var _ middleware.AuthStrategy = &JWTStrategy{}

// NewJWTStrategy create jwt bearer strategy with GinJWTMiddleware.
func NewJWTStrategy(verify func(tokenStr string) (*token.VerifiedResponse, error)) JWTStrategy {
	return JWTStrategy{verify: verify}
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
		t, err := j.verify(tokenValue)
		if err != nil {
			api.FailWithErrCode(err, c)
			c.Abort()
			return
		}
		if t.TokenType != token.AccessToken {
			api.FailWithErrCode(errors.WithCode(code.ErrPermissionDenied, "Invalid authorization header"), c)
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(request.WithUser(c.Request.Context(), t.User))

		c.Next()
	}
}
