// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/middleware"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

// BasicStrategy defines Basic authentication strategy.
type BasicStrategy struct {
	authenticate func(username string, password string) (*v1.AuthenticateResponse, error)
}

var _ middleware.AuthStrategy = &BasicStrategy{}

// NewBasicStrategy create basic strategy with compare function.
func NewBasicStrategy(authenticate func(username string, password string) (*v1.AuthenticateResponse, error)) BasicStrategy {
	return BasicStrategy{
		authenticate: authenticate,
	}
}

// AuthFunc defines basic strategy as the gin authentication middleware.
func (b BasicStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			api.FailWithErrCode(errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."), c)
			c.Abort()

			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 {
			api.FailWithErrCode(errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."), c)
			c.Abort()

			return
		}

		resp, err := b.authenticate(pair[0], pair[1])
		if err != nil {
			api.FailWithErrCode(errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."), c)
			c.Abort()

			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1.CtxKeyUserInstanceID, resp.User.InstanceID))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1.CtxKeyUserType, resp.User.UserType))

		c.Next()
	}
}
